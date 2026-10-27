package allora

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/stretchr/testify/require"
)

func testBackoff(rt http.RoundTripper) *BackoffTransport {
	return &BackoffTransport{
		roundTripper: rt,
		baseDelay:    time.Millisecond,
		maxDelay:     5 * time.Millisecond,
		jitterFrac:   0,
		maxAttempts:  3,
	}
}

func TestBackoffRetries503ThenSucceeds(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := hits.Add(1)
		if n < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	client := &http.Client{Transport: testBackoff(http.DefaultTransport)}
	resp, err := client.Get(srv.URL)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, int32(3), hits.Load())
}

func TestBackoffRetries429(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := hits.Add(1)
		if n == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := &http.Client{Transport: testBackoff(http.DefaultTransport)}
	resp, err := client.Get(srv.URL)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, int32(2), hits.Load())
}

func TestBackoffNoRetryOn400(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()

	client := &http.Client{Transport: testBackoff(http.DefaultTransport)}
	resp, err := client.Get(srv.URL)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	require.Equal(t, int32(1), hits.Load())
}

func TestBackoffNoRetryNonIdempotent(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	client := &http.Client{Transport: testBackoff(http.DefaultTransport)}
	resp, err := client.Post(srv.URL, "text/plain", strings.NewReader("body"))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	require.Equal(t, int32(1), hits.Load())
}

func TestBackoffRespectsContextCancel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	bt := testBackoff(http.DefaultTransport)
	bt.baseDelay = 50 * time.Millisecond
	bt.maxDelay = 50 * time.Millisecond
	client := &http.Client{Transport: bt}

	ctx, cancel := context.WithCancel(context.Background())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL, nil)
	require.NoError(t, err)

	go func() {
		time.Sleep(5 * time.Millisecond)
		cancel()
	}()

	_, err = client.Do(req)
	require.ErrorIs(t, err, context.Canceled)
}

func TestBackoffReturnsLastResponseAfterMaxAttempts(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("down"))
	}))
	defer srv.Close()

	client := &http.Client{Transport: testBackoff(http.DefaultTransport)}
	resp, err := client.Get(srv.URL)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, "down", string(body))
	require.Equal(t, int32(3), hits.Load())
}

func TestClassifyTransportError(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://example.invalid", nil)
	require.NoError(t, err)
	require.Nil(t, classifyTransportError(req, nil))

	classified := classifyTransportError(req, io.EOF)
	require.True(t, errors.IsRetryable(classified))
	require.Equal(t, errors.FaultUpstream, errors.GetFault(classified))

	classified = classifyTransportError(req, syscall.ECONNREFUSED)
	require.True(t, errors.IsRetryable(classified))

	classified = classifyTransportError(req, syscall.ECONNRESET)
	require.True(t, errors.IsRetryable(classified))
}
