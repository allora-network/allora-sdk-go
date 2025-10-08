package allora

import (
	"context"
	"encoding/json"
	"iter"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	stderrors "errors"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/metrics"
)

type APIClient interface {
	GetTopics() iter.Seq2[*Topic, error]
	GetTopic(topicID uint64) (*Topic, error)
	GetOHLCData(ctx context.Context, ticker string, fromDate string) iter.Seq2[*OHLCResponse, error]
}

type apiClient struct {
	url        string
	apiKey     string
	httpClient *http.Client
	logger     zerolog.Logger
}

func NewAPIClient(apiKey string, opts ...APIClientOption) *apiClient {
	if apiKey == "" {
		apiKey = "UP-8cbc632a67a84ac1b4078661"
	}

	client := &apiClient{
		url:    "https://api.allora.network/v2",
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout:   10 * time.Second,
			Transport: http.DefaultTransport,
		},
		logger: zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(client)
	}
	return client
}

type APIClientOption func(*apiClient)

// WithBackoff configures roundTripper / max / jitter for exponential delays.
func WithBackoff(baseDelay, maxDelay time.Duration, jitter float64) APIClientOption {
	return func(c *apiClient) {
		t := BackoffTransport{
			roundTripper: c.httpClient.Transport,
			baseDelay:    baseDelay,
			maxDelay:     maxDelay,
			jitterFrac:   jitter,
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

type APIResponse[T any] struct {
	RequestID string `json:"request_id"`
	Status    bool   `json:"status"`
	Data      T      `json:"data"`
}

func (c *apiClient) getURL(path string, queryParams map[string]string) (string, error) {
	u, err := url.Parse(c.url)
	if err != nil {
		return "", err
	}
	u.Path = path
	query := u.Query()
	for k, v := range queryParams {
		query.Set(k, v)
	}
	u.RawQuery = query.Encode()
	return u.String(), nil
}

type GetTopicsResponse struct {
	Topics            []Topic `json:"topics"`
	ContinuationToken *string `json:"continuation_token"`
}

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
	TotalStakedAllo           float64    `json:"total_staked_allo"`
	TotalEmissionsAllo        float64    `json:"total_emissions_allo"`
	IsActive                  bool       `json:"is_active"`
	IsEndorsed                bool       `json:"is_endorsed"`
	ForgeCompetitionID        *int       `json:"forge_competition_id"`
	ForgeCompetitionStartDate *time.Time `json:"forge_competition_start_date"`
	ForgeCompetitionEndDate   *time.Time `json:"forge_competition_end_date"`
	LatestNetworkInference    *Inference `json:"latest_network_inference"`
	UpdatedAt                 time.Time  `json:"updated_at"`
}

type Inference struct {
	CombinedValue string `json:"combined_value"`
	NaiveValue    string `json:"naive_value"`
	Timestamp     int64  `json:"timestamp"` // epoch millis
}

func (c *apiClient) GetTopics() iter.Seq2[*Topic, error] {
	return func(yield func(*Topic, error) bool) {
		var token *string
		for {
			topics, err := c.getTopics(token)
			if err != nil {
				c.logger.Error().Err(err).Msg("failed to get topics")
			}

			for _, topic := range topics.Topics {
				if !yield(&topic, err) {
					return
				}
			}
			if topics.ContinuationToken == nil {
				return
			}
			token = topics.ContinuationToken
		}
	}
}

func (c *apiClient) getTopics(token *string) (_ *GetTopicsResponse, err error) {
	methodName := "GetTopics"
	start := time.Now()
	outcome := "success"
	defer func() {
		c.recordAPIMetrics(methodName, outcome, start)
	}()

	query := map[string]string{}
	if token != nil {
		query["continuation_token"] = *token
	}

	url, err := c.getURL("/v2/allora/allora-testnet-1/topics", query)
	if err != nil {
		outcome = "url_parse_error"
		return nil, errors.WithMessage(err, "could not parse url")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		outcome = "request_build_error"
		return nil, errors.WithMessage(err, "failed to create request")
	}

	req.Header.Set("x-api-key", c.apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		outcome = classifyAPIError(err)
		return nil, errors.WithMessage(err, "failed to execute request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		outcome = "http_" + strconv.Itoa(resp.StatusCode)
		return nil, errors.Errorf("GetTopics unexpected status code: %d", resp.StatusCode)
	}

	var topicsResp APIResponse[GetTopicsResponse]
	if err := json.NewDecoder(resp.Body).Decode(&topicsResp); err != nil {
		outcome = "decode_error"
		return nil, errors.WithMessage(err, "failed to decode response")
	}
	return &topicsResp.Data, nil
}

func (c *apiClient) GetTopic(topicID uint64) (_ *Topic, err error) {
	methodName := "GetTopic"
	start := time.Now()
	outcome := "success"
	defer func() {
		c.recordAPIMetrics(methodName, outcome, start)
	}()

	url, err := c.getURL("/v2/allora/allora-testnet-1/topics/"+strconv.FormatUint(topicID, 10), nil)
	if err != nil {
		outcome = "url_parse_error"
		return nil, errors.WithMessage(err, "could not parse url")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		outcome = "request_build_error"
		return nil, errors.WithMessage(err, "failed to create request")
	}

	req.Header.Set("x-api-key", c.apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		outcome = classifyAPIError(err)
		return nil, errors.WithMessage(err, "failed to execute request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		outcome = "http_" + strconv.Itoa(resp.StatusCode)
		return nil, errors.Errorf("GetTopic unexpected status code: %d", resp.StatusCode)
	}

	var topicResp APIResponse[Topic]
	if err := json.NewDecoder(resp.Body).Decode(&topicResp); err != nil {
		outcome = "decode_error"
		return nil, errors.WithMessage(err, "failed to decode response")
	}
	return &topicResp.Data, nil
}

type OHLCResponse struct {
	Data              []OHLCData `json:"data"`
	ContinuationToken *string    `json:"continuation_token"`
}

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

func (c *apiClient) GetOHLCData(ctx context.Context, ticker string, fromDate string) iter.Seq2[*OHLCResponse, error] {
	return func(yield func(*OHLCResponse, error) bool) {
		var token *string
		for {
			ohlcData, err := c.getOHLCData(ctx, ticker, fromDate, token)
			if err != nil {
				c.logger.Error().Err(err).Msg("failed to get OHLC data")
			}

			if !yield(ohlcData, err) {
				return
			} else if ohlcData.ContinuationToken == nil {
				return
			}
			token = ohlcData.ContinuationToken
		}
	}
}

func (c *apiClient) getOHLCData(ctx context.Context, ticker string, fromDate string, token *string) (_ *OHLCResponse, err error) {
	methodName := "GetOHLCData"
	start := time.Now()
	outcome := "success"
	defer func() {
		c.recordAPIMetrics(methodName, outcome, start)
	}()

	query := map[string]string{
		"tickers":   ticker,
		"from_date": fromDate,
	}
	if token != nil {
		query["continuation_token"] = *token
	}

	url, err := c.getURL("/v2/allora/market-data/ohlc", query)
	if err != nil {
		outcome = "url_parse_error"
		return nil, errors.WithMessage(err, "could not parse url")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		outcome = "request_build_error"
		return nil, errors.WithMessage(err, "failed to create request")
	}

	req.Header.Set("x-api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		outcome = classifyAPIError(err)
		return nil, errors.WithMessage(err, "failed to execute request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		outcome = "http_" + strconv.Itoa(resp.StatusCode)
		return nil, errors.Errorf("GetOHLCData unexpected status code: %d", resp.StatusCode)
	}

	var ohlcResp APIResponse[OHLCResponse]
	if err := json.NewDecoder(resp.Body).Decode(&ohlcResp); err != nil {
		outcome = "decode_error"
		return nil, errors.WithMessage(err, "failed to decode response")
	}

	// check if the response is empty
	if len(ohlcResp.Data.Data) == 0 {
		outcome = "empty_response"
		return nil, errors.WithMessagef(err, "no data found for ticker %s", ticker)
	}
	return &ohlcResp.Data, nil
}

func (c *apiClient) recordAPIMetrics(method, outcome string, start time.Time) {
	metrics.ObserveRPCRequest("http_api", c.url, "api", method, outcome, 1, time.Since(start))
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

type BackoffTransport struct {
	roundTripper http.RoundTripper
	baseDelay    time.Duration
	maxDelay     time.Duration
	jitterFrac   float64
	logger       zerolog.Logger
}

const (
	defaultBaseDelay  = 1 * time.Second
	defaultMaxDelay   = 30 * time.Second
	defaultJitterFrac = 0.1
)

func NewBackoffTransport(roundTripper http.RoundTripper) http.RoundTripper {
	return &BackoffTransport{
		roundTripper: roundTripper,
		baseDelay:    defaultBaseDelay,
		maxDelay:     defaultMaxDelay,
		jitterFrac:   defaultJitterFrac,
		logger:       zerolog.Nop(),
	}
}

func (bt *BackoffTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := bt.roundTripper
	if base == nil {
		base = http.DefaultTransport
	}

	baseDelay := bt.baseDelay
	if baseDelay <= 0 {
		baseDelay = defaultBaseDelay
	}
	maxDelay := bt.maxDelay
	if maxDelay <= 0 {
		maxDelay = defaultMaxDelay
	}
	jitterFrac := bt.jitterFrac
	if jitterFrac < 0 {
		jitterFrac = 0
	}

	nextDelay := baseDelay
	attempt := 0

	for {
		// For attempts > 0, ensure body is resettable before we try again
		if attempt > 0 {
			if !canReplayBody(req) {
				return nil, errors.Errorf("cannot retry non-replayable request body")
			}
			if err := resetBody(req); err != nil {
				return nil, err
			}
		}

		resp, err := base.RoundTrip(req)
		if err == nil && (resp == nil || !isRetryableStatus(resp.StatusCode)) {
			return resp, nil
		}

		if req.Context().Err() != nil {
			if resp != nil && resp.Body != nil {
				_ = resp.Body.Close()
			}
			if req.Context().Err() != nil {
				return nil, req.Context().Err()
			}
			return resp, err
		}

		shouldRetry := false
		var retryAfter time.Duration

		if err != nil {
			shouldRetry = isRetryableError(req, err)
		} else if resp != nil {
			if isRetryableStatus(resp.StatusCode) {
				shouldRetry = true
				// honor Retry-After if provided and parseable
				if ra := resp.Header.Get("Retry-After"); ra != "" {
					if secs, perr := strconv.Atoi(ra); perr == nil && secs >= 0 {
						retryAfter = time.Duration(secs) * time.Second
					}
				}
			}
		}

		// only retry for idempotent or replayable-body requests
		if shouldRetry && !(isIdempotent(req.Method) || canReplayBody(req)) {
			shouldRetry = false
		}

		if !shouldRetry {
			if resp != nil && resp.Body != nil {
				// Return the response to caller (non-retryable status)
				return resp, nil
			}
			// Return the transport error
			return nil, err
		}

		// Close response body before retrying to avoid leaks
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}

		// Compute backoff with jitter, honoring Retry-After if specified
		sleep := nextDelay
		if retryAfter > 0 {
			sleep = retryAfter
		}
		// Apply jitter: multiply by a factor in [1-j, 1+j]
		if jitterFrac > 0 {
			jf := (rand.Float64()*2 - 1) * jitterFrac // [-jitter, +jitter]
			// guard against negative
			factor := math.Max(0, 1+jf)
			sleep = time.Duration(float64(sleep) * factor)
		}

		// Cap the sleep at maxDelay
		if sleep > maxDelay {
			sleep = maxDelay
		}

		// Log retry attempt
		bt.logger.Warn().
			Str("method", req.Method).
			Str("url", req.URL.String()).
			Int("attempt", attempt+1).
			Dur("backoff", sleep).
			Msg("retrying HTTP request")

		// Sleep respecting context
		timer := time.NewTimer(sleep)
		select {
		case <-timer.C:
			// proceed
		case <-req.Context().Done():
			if !timer.Stop() {
				<-timer.C
			}
			return nil, req.Context().Err()
		}

		// Exponential growth for next iteration
		nextDelay = nextDelay * 2
		if nextDelay > maxDelay {
			nextDelay = maxDelay
		}
		attempt++
	}
}

func isIdempotent(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
		return true
	default:
		return false
	}
}

func canReplayBody(r *http.Request) bool {
	if r.Body == nil {
		return true
	}
	return r.GetBody != nil
}

func resetBody(r *http.Request) error {
	if r.Body == nil {
		return nil
	}
	if r.GetBody == nil {
		return errors.Errorf("request body is not replayable for retries")
	}
	// Close old body before replacing
	_ = r.Body.Close()
	nb, err := r.GetBody()
	if err != nil {
		return err
	}
	r.Body = nb
	return nil
}

func isRetryableStatus(code int) bool {
	if code == http.StatusTooManyRequests || code == http.StatusRequestTimeout {
		return true
	}
	return code >= 500
}

func isRetryableError(req *http.Request, err error) bool {
	if err == nil {
		return false
	}
	// Respect request context cancellation/deadline
	if req.Context().Err() != nil {
		return false
	}
	if ne, ok := err.(net.Error); ok {
		if ne.Timeout() {
			return true
		}
		// Many transient network errors implement Temporary()
		// (Deprecated in Go 1.18+, but still implemented by some types)
		type temporary interface{ Temporary() bool }
		if t, ok := any(ne).(temporary); ok && t.Temporary() {
			return true
		}
	}
	// Fallback: treat other transport-level errors as retryable
	return true
}
