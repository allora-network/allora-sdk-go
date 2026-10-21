package allora

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	stderrors "errors"

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
	forgeAPIKeyPrefix = "forge_sk_"

	// DefaultAPIBaseURL is the Forge API Gateway origin used when WithBaseURL is omitted.
	DefaultAPIBaseURL = "https://api.developer.allora.network"
	// DefaultDataAPIBaseURL is the Forge Data (Atlas) root used when WithDataBaseURL is omitted.
	// Staging host by design until a production data default is published.
	DefaultDataAPIBaseURL = "https://forge-data.stage.allora.run/api"

	maxAPIResponseBytes = 32 << 20 // 32 MiB
)

// APIClient is the compatible public surface from the previous SDK.
// Extended Forge/Atlas helpers live on *apiClient but are not part of this interface.
type APIClient interface {
	GetTopics() iter.Seq2[*Topic, error]
	GetTopic(topicID uint64) (*Topic, error)
	GetOHLCData(ctx context.Context, ticker string, fromDate string) iter.Seq2[*OHLCResponse, error]
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
// apiKey must be a Forge secret key (forge_sk_…).
// Gateway and data base URLs default to DefaultAPIBaseURL and DefaultDataAPIBaseURL;
// override with WithBaseURL / WithDataBaseURL (e.g. local or alternate stage).
func NewAPIClient(apiKey string, opts ...APIClientOption) (*apiClient, error) {
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return nil, errors.New("ALLORA_API_KEY is required (forge_sk_…)")
	}
	if !strings.HasPrefix(apiKey, forgeAPIKeyPrefix) {
		return nil, errors.Errorf("API key must start with %s", forgeAPIKeyPrefix)
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
func WithNetwork(network Network) APIClientOption {
	return func(c *apiClient) {
		c.network = network
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
			logger:       c.logger,
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
			logger:       c.logger,
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

// WithLogger sets the logger
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
var ErrTopicNotFound = stderrors.New("topic not found")

// ErrDatasetNotFound is returned when a ticker/dataset cannot be resolved in Atlas.
var ErrDatasetNotFound = stderrors.New("dataset not found")

// ErrAmbiguousTicker is returned when a ticker matches multiple Atlas datasets
// and no unique 1-minute preference can be applied.
var ErrAmbiguousTicker = stderrors.New("ambiguous ticker")

type gatewayErrorBody struct {
	Error string `json:"error"`
}

// Topic is the previous SDK topic shape (compat).
// Only fields the gateway actually supplies are populated; other fields stay zero.
type Topic struct {
	TopicID                   int        `json:"topic_id"`
	TopicName                 string     `json:"topic_name"`
	Description               *string    `json:"description"`
	EpochLength               int        `json:"epoch_length"`
	GroundTruthLag            int        `json:"ground_truth_lag"`
	LossMethod                string     `json:"loss_method"`
	WorkerSubmissionWindow    int        `json:"worker_submission_window"`
	WorkerCount               int        `json:"worker_count"`
	ReputerCount              int        `json:"reputer_count"`
	TotalStakedAllo           int        `json:"total_staked_allo"`
	TotalEmissionsAllo        int        `json:"total_emissions_allo"`
	IsActive                  bool       `json:"is_active"`
	IsEndorsed                bool       `json:"is_endorsed"`
	ForgeCompetitionID        *int       `json:"forge_competition_id"`
	ForgeCompetitionStartDate *time.Time `json:"forge_competition_start_date"`
	ForgeCompetitionEndDate   *time.Time `json:"forge_competition_end_date"`
	LatestNetworkInference    *Inference `json:"latest_network_inference"`
	UpdatedAt                 time.Time  `json:"updated_at"`
}

// Inference is the previous SDK prediction shape (compat).
type Inference struct {
	CombinedValue string `json:"combined_value"`
	NaiveValue    string `json:"naive_value"`
	Timestamp     int64  `json:"timestamp"` // epoch millis
}

// OHLCResponse is the previous SDK candles page shape (compat).
type OHLCResponse struct {
	Data              []OHLCData `json:"data"`
	ContinuationToken *string    `json:"continuation_token"`
}

// OHLCData is one candle in the previous SDK shape (compat).
type OHLCData struct {
	Ticker         string `json:"ticker"`
	ExchangeCode   string `json:"exchange_code"`
	Date           string `json:"date"`
	Open           string `json:"open"`
	High           string `json:"high"`
	Low            string `json:"low"`
	Close          string `json:"close"`
	TradesDone     int    `json:"trades_done"`
	Volume         string `json:"volume"`
	VolumeNotional string `json:"volume_notional"`
}

// APIResponse is the generic envelope returned by the previous SDK's Upshot
// endpoints (request_id / status / data). It is preserved for callers that
// reference it directly; the Forge gateway/Atlas backends do not use it.
type APIResponse[T any] struct {
	RequestID string `json:"request_id"`
	Status    bool   `json:"status"`
	Data      T      `json:"data"`
}

// GetTopicsResponse is the page shape returned by the previous SDK's
// GetTopics helper. Kept for backward compatibility; the Forge gateway
// returns a different envelope that is decoded internally.
type GetTopicsResponse struct {
	Topics            []Topic `json:"topics"`
	ContinuationToken *string `json:"continuation_token"`
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
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    ge.Error,
			Body:       string(body),
		}
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
	case stderrors.Is(err, context.Canceled):
		return "context_canceled"
	case stderrors.Is(err, context.DeadlineExceeded):
		return "context_deadline"
	}

	var netErr net.Error
	if stderrors.As(err, &netErr) {
		if netErr.Timeout() {
			return "network_timeout"
		}
		if netErr.Temporary() {
			return "network_temporary"
		}
		return "network_error"
	}

	return "error"
}
