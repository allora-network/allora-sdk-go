package allora

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/metrics"
)

// Network is the gateway path segment (not a Cosmos chain id).
type Network string

const (
	NetworkTestnet  Network = "testnet"
	NetworkMainnet  Network = "mainnet"
	NetworkOffchain Network = "offchain"

	defaultAPINetwork = NetworkMainnet

	// DefaultAPIBaseURL is the Forge API Gateway origin used when WithBaseURL is omitted.
	DefaultAPIBaseURL = "https://api.developer.allora.network"
	// DefaultDataAPIBaseURL is the Forge Data (Atlas) root used when WithDataBaseURL is omitted.
	DefaultDataAPIBaseURL = "https://forge-data.allora.run/api"

	maxAPIResponseBytes = 32 << 20 // 32 MiB
)

// APIClient is the public Forge HTTP surface (gateway topics/inferences + Atlas candles).
type APIClient interface {
	GetTopics() iter.Seq2[*Topic, error]
	GetTopic(topicID uint64) (*Topic, error)
	GetNetworkInference(ctx context.Context, topicID uint64, opts ...InferenceOption) (*NetworkInference, error)
	IterCandles(ctx context.Context, q CandleQuery) iter.Seq2[[]Candle, error]
}

type apiClient struct {
	url        string // API Gateway origin (defaults to DefaultAPIBaseURL)
	dataURL    string // Forge Data / Atlas origin including /api (defaults to DefaultDataAPIBaseURL)
	apiKey     string
	network    Network
	httpClient *http.Client
	logger     zerolog.Logger
}

// NewAPIClient constructs a Forge API client.
// apiKey is sent as Bearer (gateway) and X-API-Key (Atlas); the server validates it.
// Gateway and data base URLs default to DefaultAPIBaseURL and DefaultDataAPIBaseURL;
// override with WithBaseURL / WithDataBaseURL (e.g. local or stage).
func NewAPIClient(apiKey string, opts ...APIClientOption) (*apiClient, error) {
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return nil, errors.New("ALLORA_API_KEY is required")
	}

	client := &apiClient{
		url:     DefaultAPIBaseURL,
		dataURL: DefaultDataAPIBaseURL,
		apiKey:  apiKey,
		network: defaultAPINetwork,
		httpClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: http.DefaultTransport,
		},
		logger: zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(client)
	}

	switch client.network {
	case NetworkTestnet, NetworkMainnet, NetworkOffchain:
	default:
		return nil, errors.Errorf("unsupported network %q (want testnet|mainnet|offchain)", client.network)
	}

	// Refuse cross-host redirects so custom auth headers (esp. X-API-Key) are
	// never forwarded to an unexpected origin. net/http only strips Authorization.
	client.httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		if len(via) == 0 {
			return nil
		}
		origin := via[0].URL
		if req.URL.Scheme != origin.Scheme || req.URL.Host != origin.Host {
			return errors.Errorf("refusing cross-host redirect to %s", req.URL.Host)
		}
		return nil
	}
	return client, nil
}

type APIClientOption func(*apiClient)

// WithBaseURL overrides the API Gateway origin (default: DefaultAPIBaseURL).
func WithBaseURL(baseURL string) APIClientOption {
	return func(c *apiClient) {
		c.url = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	}
}

// WithDataBaseURL overrides the Forge Data (Atlas) API root (default: DefaultDataAPIBaseURL).
// Must include the /api prefix (e.g. https://host/api).
func WithDataBaseURL(baseURL string) APIClientOption {
	return func(c *apiClient) {
		c.dataURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	}
}

// WithNetwork sets the gateway network path segment (testnet|mainnet|offchain).
// Defaults to mainnet when omitted.
func WithNetwork(network Network) APIClientOption {
	return func(c *apiClient) {
		c.network = network
	}
}

// WithHTTPClient replaces the underlying *http.Client with a shallow copy of
// hc so the caller's shared client is not mutated (CheckRedirect / Transport /
// Timeout are adjusted on the copy). Compose with WithTimeout / WithBackoff /
// WithDefaultBackoff as needed (order matters).
func WithHTTPClient(hc *http.Client) APIClientOption {
	return func(c *apiClient) {
		if hc != nil {
			cp := *hc
			c.httpClient = &cp
		}
	}
}

// WithBackoff configures roundTripper / max / jitter for exponential delays.
func WithBackoff(baseDelay, maxDelay time.Duration, jitter float64) APIClientOption {
	return func(c *apiClient) {
		t := BackoffTransport{
			roundTripper: c.httpClient.Transport,
			baseDelay:    baseDelay,
			maxDelay:     maxDelay,
			jitterFrac:   jitter,
			maxAttempts:  defaultMaxAttempts,
			logger:       &c.logger,
		}
		c.httpClient.Transport = &t
	}
}

// WithDefaultBackoff configures the HTTP client transport with default
// exponential backoff parameters.
func WithDefaultBackoff() APIClientOption {
	return func(c *apiClient) {
		t := BackoffTransport{
			roundTripper: c.httpClient.Transport,
			baseDelay:    defaultBaseDelay,
			maxDelay:     defaultMaxDelay,
			jitterFrac:   defaultJitterFrac,
			maxAttempts:  defaultMaxAttempts,
			logger:       &c.logger,
		}
		c.httpClient.Transport = &t
	}
}

// WithMetrics wraps the transport with a Prometheus metrics collector
func WithMetrics() APIClientOption {
	return func(c *apiClient) {
		c.httpClient.Transport = metrics.NewMetricsRoundTripper(c.httpClient.Transport)
	}
}

// WithTimeout sets the timeout for the HTTP client
func WithTimeout(d time.Duration) APIClientOption {
	return func(c *apiClient) {
		c.httpClient.Timeout = d
	}
}

// WithLogger sets the logger. Safe to call after WithDefaultBackoff / WithBackoff
// because those options store a pointer to the client's logger field.
func WithLogger(logger zerolog.Logger) APIClientOption {
	return func(c *apiClient) {
		c.logger = logger
	}
}

// APIError is a non-2xx HTTP response from Forge Gateway or Atlas.
type APIError struct {
	StatusCode int
	Message    string
	Body       string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("API error: status=%d message=%s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("API error: status=%d", e.StatusCode)
}

// ErrTopicNotFound is returned by GetTopic when the id is absent from the catalog.
var ErrTopicNotFound = errors.Sentinel("topic not found", errors.ErrNotFound)

// ErrDatasetNotFound is returned when a ticker/dataset cannot be resolved in Atlas.
var ErrDatasetNotFound = errors.Sentinel("dataset not found", errors.ErrNotFound)

// ErrAmbiguousTicker is returned when a ticker matches multiple Atlas datasets
// and no unique 1-minute preference can be applied.
var ErrAmbiguousTicker = errors.Sentinel("ambiguous ticker", errors.ErrBadRequest)

type gatewayErrorBody struct {
	Error string `json:"error"`
}

type authStyle int

const (
	authGatewayBearer authStyle = iota
	authDataAPIKey
)

func (c *apiClient) resolveURL(base, path string, queryParams map[string]string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	*u = *u.JoinPath(strings.TrimPrefix(path, "/"))
	query := u.Query()
	for k, v := range queryParams {
		query.Set(k, v)
	}
	u.RawQuery = query.Encode()
	return u.String(), nil
}

func (c *apiClient) gatewayURL(path string, queryParams map[string]string) (string, error) {
	return c.resolveURL(c.url, path, queryParams)
}

func (c *apiClient) dataURLPath(path string, queryParams map[string]string) (string, error) {
	return c.resolveURL(c.dataURL, path, queryParams)
}

func (c *apiClient) setAuth(req *http.Request, style authStyle) {
	req.Header.Set("Accept", "application/json")
	switch style {
	case authDataAPIKey:
		req.Header.Set("X-API-Key", c.apiKey)
	default:
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
}

func (c *apiClient) doJSON(ctx context.Context, methodName, method, reqURL string, dest any) error {
	return c.doJSONWithAuth(ctx, methodName, method, reqURL, c.url, dest, authGatewayBearer)
}

func (c *apiClient) doDataJSON(ctx context.Context, methodName, method, reqURL string, dest any) error {
	return c.doJSONWithAuth(ctx, methodName, method, reqURL, c.dataURL, dest, authDataAPIKey)
}

func (c *apiClient) doJSONWithAuth(
	ctx context.Context,
	methodName, method, reqURL, endpoint string,
	dest any,
	style authStyle,
) (err error) {
	start := time.Now()
	outcome := "success"
	defer func() {
		c.recordAPIMetrics(endpoint, methodName, outcome, start)
	}()

	req, err := http.NewRequestWithContext(ctx, method, reqURL, nil)
	if err != nil {
		outcome = "request_build_error"
		return errors.WithMessage(err, "failed to create request")
	}
	c.setAuth(req, style)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		err = classifyTransportError(req, err)
		outcome = classifyAPIError(err)
		return errors.WithMessage(err, "failed to execute request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxAPIResponseBytes+1))
	if err != nil {
		outcome = "read_error"
		return errors.WithMessage(err, "failed to read response body")
	}
	if int64(len(body)) > maxAPIResponseBytes {
		outcome = "response_too_large"
		return errors.Errorf("response body exceeds %d bytes", maxAPIResponseBytes)
	}

	if resp.StatusCode != http.StatusOK {
		outcome = "http_" + strconv.Itoa(resp.StatusCode)
		var ge gatewayErrorBody
		_ = json.Unmarshal(body, &ge)
		apiErr := &APIError{
			StatusCode: resp.StatusCode,
			Message:    ge.Error,
			Body:       string(body),
		}
		// Debug: callers / higher methods may log the returned error at Error.
		c.logger.Debug().
			Err(apiErr).
			Int("status", resp.StatusCode).
			Str("endpoint", endpoint).
			Str("method", methodName).
			Str("url", reqURL).
			Msg("forge api request failed")

		retry := errors.NonRetryable
		if isRetryableStatus(resp.StatusCode) {
			retry = errors.Retryable
		}
		fault := errors.FaultUpstream
		if resp.StatusCode < 500 {
			fault = errors.FaultCaller
		}
		return errors.Wrap(
			apiErr,
			"forge api request failed",
			errors.StatusCode(resp.StatusCode), fault, retry,
			"endpoint", endpoint, "method", methodName,
		)
	}

	if dest == nil {
		return nil
	}
	if err := json.Unmarshal(body, dest); err != nil {
		outcome = "decode_error"
		return errors.WithMessage(err, "failed to decode response")
	}
	return nil
}

func (c *apiClient) recordAPIMetrics(endpoint, method, outcome string, start time.Time) {
	metrics.ObserveRPCRequest("http_api", endpoint, "api", method, outcome, 1, time.Since(start))
	metrics.ObserveRPCAttempts("http_api", "api", method, outcome, 1)
}

func classifyAPIError(err error) string {
	switch {
	case err == nil:
		return "success"
	case errors.Is(err, context.Canceled):
		return "context_canceled"
	case errors.Is(err, context.DeadlineExceeded):
		return "context_deadline"
	case errors.IsRetryable(err):
		return "retryable_" + errors.GetFault(err).String()
	default:
		return "error_" + errors.GetFault(err).String()
	}
}
