package allora

import (
	"io"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
)

// BackoffTransport retries idempotent requests with exponential backoff.
// logger is a pointer so WithLogger applied after WithDefaultBackoff still works.
type BackoffTransport struct {
	roundTripper http.RoundTripper
	baseDelay    time.Duration
	maxDelay     time.Duration
	jitterFrac   float64
	maxAttempts  int
	logger       *zerolog.Logger
}

const (
	defaultBaseDelay   = 1 * time.Second
	defaultMaxDelay    = 30 * time.Second
	defaultJitterFrac  = 0.1
	defaultMaxAttempts = 5
)

func NewBackoffTransport(roundTripper http.RoundTripper) http.RoundTripper {
	nop := zerolog.Nop()
	return &BackoffTransport{
		roundTripper: roundTripper,
		baseDelay:    defaultBaseDelay,
		maxDelay:     defaultMaxDelay,
		jitterFrac:   defaultJitterFrac,
		maxAttempts:  defaultMaxAttempts,
		logger:       &nop,
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
	maxAttempts := bt.maxAttempts
	if maxAttempts <= 0 {
		maxAttempts = defaultMaxAttempts
	}

	nextDelay := baseDelay

	for attempt := 0; attempt < maxAttempts; attempt++ {
		attemptReq := req.Clone(req.Context())
		if attempt > 0 {
			if err := resetBody(attemptReq); err != nil {
				return nil, err
			}
		}

		resp, err := base.RoundTrip(attemptReq)
		if err != nil {
			err = classifyTransportError(req, err)
		}
		if err == nil && (resp == nil || !isRetryableStatus(resp.StatusCode)) {
			return resp, nil
		}

		if req.Context().Err() != nil {
			if resp != nil && resp.Body != nil {
				_ = resp.Body.Close()
			}
			return nil, req.Context().Err()
		}

		shouldRetry := false
		var retryAfter time.Duration

		if err != nil {
			shouldRetry = errors.IsRetryable(err)
		} else if resp != nil && isRetryableStatus(resp.StatusCode) {
			shouldRetry = true
			retryAfter = parseRetryAfter(resp.Header.Get("Retry-After"))
		}

		// Only retry idempotent methods with a replayable body.
		// Annotate non-idempotent veto so callers are not told a POST is safe to retry.
		if shouldRetry && (!isIdempotent(req.Method) || !canReplayBody(req)) {
			shouldRetry = false
			if err != nil {
				err = errors.WithMetadata(err, errors.NonRetryable)
			}
		}

		if !shouldRetry {
			if resp != nil {
				return resp, nil
			}
			return nil, err
		}

		// Last attempt: return the response/error so callers see status + body.
		if attempt+1 >= maxAttempts {
			if resp != nil {
				return resp, nil
			}
			return nil, err
		}

		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}

		sleep := nextDelay
		honorServerDelay := retryAfter > 0
		if honorServerDelay {
			sleep = retryAfter
		}
		if jitterFrac > 0 {
			jf := (rand.Float64()*2 - 1) * jitterFrac
			factor := math.Max(0, 1+jf)
			sleep = time.Duration(float64(sleep) * factor)
		}
		// Cap exponential backoff; do not cap a server-requested Retry-After.
		if !honorServerDelay && sleep > maxDelay {
			sleep = maxDelay
		}

		log := zerolog.Nop()
		if bt.logger != nil {
			log = *bt.logger
		}
		log.Warn().
			Str("method", req.Method).
			Str("url", req.URL.String()).
			Int("attempt", attempt+1).
			Dur("backoff", sleep).
			Msg("retrying HTTP request")

		timer := time.NewTimer(sleep)
		select {
		case <-timer.C:
		case <-req.Context().Done():
			if !timer.Stop() {
				<-timer.C
			}
			return nil, req.Context().Err()
		}

		nextDelay *= 2
		if nextDelay > maxDelay {
			nextDelay = maxDelay
		}
	}

	return nil, errors.New("backoff: exhausted attempts")
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
	_ = r.Body.Close()
	nb, err := r.GetBody()
	if err != nil {
		return err
	}
	r.Body = nb
	return nil
}

// isRetryableStatus reports HTTP statuses the transport may retry.
// Includes 408, 429, and all 5xx (including 503).
func isRetryableStatus(code int) bool {
	if code == http.StatusTooManyRequests || code == http.StatusRequestTimeout {
		return true
	}
	return code >= 500
}

// classifyTransportError annotates an error returned by a round trip of req
// with its fault and retryability so callers can use errors.IsRetryable /
// errors.GetFault without re-deriving the answer.
func classifyTransportError(req *http.Request, err error) error {
	if err == nil {
		return nil
	}
	fault := errors.FaultUpstream
	retry := errors.NonRetryable
	switch {
	case req.Context().Err() != nil:
		// Caller aborted (cancel or request-context deadline): not upstream.
		fault = errors.FaultCaller
	case isTransient(err):
		retry = errors.Retryable
	case isCallerTransportError(err):
		fault = errors.FaultCaller
	}
	host := ""
	if req.URL != nil {
		host = req.URL.Host
	}
	return errors.Wrap(err, "forge http request failed",
		fault, retry,
		"method", req.Method, "host", host,
	)
}

func isTransient(err error) bool {
	if errors.OneOf(err, syscall.ECONNREFUSED, syscall.ECONNRESET,
		io.EOF, io.ErrUnexpectedEOF, net.ErrClosed) {
		return true
	}
	ne, ok := errors.AsType[net.Error](err)
	return ok && ne.Timeout()
}

// isCallerTransportError reports config/request-construction failures (bad
// base URL scheme, unparseable URL) rather than upstream outages.
func isCallerTransportError(err error) bool {
	msg := err.Error()
	if strings.Contains(msg, "unsupported protocol scheme") ||
		strings.Contains(msg, "invalid URL") {
		return true
	}
	ue, ok := errors.AsType[*url.Error](err)
	return ok && ue.Op == "parse"
}

func parseRetryAfter(ra string) time.Duration {
	ra = strings.TrimSpace(ra)
	if ra == "" {
		return 0
	}
	if secs, err := strconv.Atoi(ra); err == nil && secs >= 0 {
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(ra); err == nil {
		d := time.Until(t)
		if d < 0 {
			return 0
		}
		return d
	}
	return 0
}
