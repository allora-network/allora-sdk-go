package allora

import (
	stderrors "errors"
	"io"
	"math"
	"math/rand"
	"net"
	"net/http"
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
			shouldRetry = isRetryableError(req, err)
		} else if resp != nil && isRetryableStatus(resp.StatusCode) {
			shouldRetry = true
			retryAfter = parseRetryAfter(resp.Header.Get("Retry-After"))
		}

		// Only retry idempotent methods with a replayable body.
		if shouldRetry && (!isIdempotent(req.Method) || !canReplayBody(req)) {
			shouldRetry = false
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

func isRetryableStatus(code int) bool {
	if code == http.StatusTooManyRequests || code == http.StatusRequestTimeout {
		return true
	}
	return code >= 500
}

// isRetryableError reports transient transport failures. Permanent errors
// (TLS handshake, bad scheme, etc.) are not retried. Connection refused/reset
// and EOF are treated as transient (backend restarts / rolling deploys).
func isRetryableError(req *http.Request, err error) bool {
	if err == nil {
		return false
	}
	if req.Context().Err() != nil {
		return false
	}
	if stderrors.Is(err, syscall.ECONNREFUSED) ||
		stderrors.Is(err, syscall.ECONNRESET) ||
		stderrors.Is(err, io.EOF) ||
		stderrors.Is(err, io.ErrUnexpectedEOF) ||
		stderrors.Is(err, net.ErrClosed) {
		return true
	}
	var ne net.Error
	if stderrors.As(err, &ne) {
		if ne.Timeout() {
			return true
		}
		type temporary interface{ Temporary() bool }
		if t, ok := any(ne).(temporary); ok && t.Temporary() {
			return true
		}
	}
	return false
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
