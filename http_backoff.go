package allora

import (
	"math"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
)

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
			return nil, req.Context().Err()
		}

		shouldRetry := false
		var retryAfter time.Duration

		if err != nil {
			shouldRetry = isRetryableError(req, err)
		} else if resp != nil && isRetryableStatus(resp.StatusCode) {
			shouldRetry = true
			if ra := resp.Header.Get("Retry-After"); ra != "" {
				if secs, perr := strconv.Atoi(ra); perr == nil && secs >= 0 {
					retryAfter = time.Duration(secs) * time.Second
				}
			}
		}

		if shouldRetry && !(isIdempotent(req.Method) || canReplayBody(req)) {
			shouldRetry = false
		}

		if !shouldRetry {
			if resp != nil && resp.Body != nil {
				return resp, nil
			}
			return nil, err
		}

		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}

		sleep := nextDelay
		if retryAfter > 0 {
			sleep = retryAfter
		}
		if jitterFrac > 0 {
			jf := (rand.Float64()*2 - 1) * jitterFrac
			factor := math.Max(0, 1+jf)
			sleep = time.Duration(float64(sleep) * factor)
		}
		if sleep > maxDelay {
			sleep = maxDelay
		}

		bt.logger.Warn().
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
	if req.Context().Err() != nil {
		return false
	}
	if ne, ok := err.(net.Error); ok {
		if ne.Timeout() {
			return true
		}
		type temporary interface{ Temporary() bool }
		if t, ok := any(ne).(temporary); ok && t.Temporary() {
			return true
		}
	}
	return true
}
