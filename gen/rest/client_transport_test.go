package rest

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestNewRESTClientCoreConfiguresConnectionPool(t *testing.T) {
	core := NewRESTClientCore("http://example.invalid", zerolog.Nop())

	require.Equal(t, defaultMaxIdleConns, core.transport.MaxIdleConns)
	require.Equal(t, defaultMaxIdleConnsPerHost, core.transport.MaxIdleConnsPerHost)
	require.Equal(t, defaultIdleConnTimeout, core.transport.IdleConnTimeout)
	require.Equal(t, defaultRequestTimeout, core.httpClient.Timeout)
	require.Same(t, core.transport, core.httpClient.Transport)
}

// countingRoundTripper stands in for the instrumentation wrappers callers
// install as http.DefaultTransport. It is deliberately not an *http.Transport,
// so it cannot be cloned or pool-tuned.
type countingRoundTripper struct {
	base  http.RoundTripper
	calls atomic.Int64
}

func (rt *countingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.calls.Add(1)
	return rt.base.RoundTrip(req)
}

// useDefaultTransport installs rt as http.DefaultTransport for the duration of
// the test. Tests that call it must not run in parallel with anything issuing
// HTTP requests.
func useDefaultTransport(t *testing.T, rt http.RoundTripper) {
	t.Helper()

	previous := http.DefaultTransport
	http.DefaultTransport = rt
	t.Cleanup(func() { http.DefaultTransport = previous })
}

func okHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(`{}`))
}

// TestNewRESTClientCoreKeepsDefaultTransportItCannotClone covers the process
// that has replaced http.DefaultTransport with a RoundTripper of its own:
// building a bare transport instead would take that wrapper off every REST
// request. Pool tuning is given up in exchange, so the client must also not
// claim ownership of the borrowed round tripper.
func TestNewRESTClientCoreKeepsDefaultTransportItCannotClone(t *testing.T) {
	server, _ := newConnCountingServer(t, okHandler)

	wrapper := &countingRoundTripper{base: http.DefaultTransport}
	useDefaultTransport(t, wrapper)

	core := NewRESTClientCore(server.URL, zerolog.Nop(), WithConnectionTimeout(time.Second))
	require.Nil(t, core.transport, "an uncloneable round tripper is borrowed, not owned")
	require.Same(t, wrapper, core.httpClient.Transport)

	err := core.executeRequest(context.Background(), http.MethodGet, "/health", nil, nil, nil, nil, 0)
	require.NoError(t, err)
	require.Equal(t, int64(1), wrapper.calls.Load())
}

// TestWithMetricsWrapsBorrowedRoundTripper checks that metrics collection still
// composes over a round tripper the client does not own.
func TestWithMetricsWrapsBorrowedRoundTripper(t *testing.T) {
	server, _ := newConnCountingServer(t, okHandler)

	wrapper := &countingRoundTripper{base: http.DefaultTransport}
	useDefaultTransport(t, wrapper)

	core := NewRESTClientCore(server.URL, zerolog.Nop(), WithMetrics())
	require.Nil(t, core.transport)
	require.NotSame(t, wrapper, core.httpClient.Transport, "the metrics collector should sit in front of the wrapper")

	err := core.executeRequest(context.Background(), http.MethodGet, "/health", nil, nil, nil, nil, 0)
	require.NoError(t, err)
	require.Equal(t, int64(1), wrapper.calls.Load())
}

// TestWithMetricsKeepsOwnedTransport checks the same composition in the usual
// case, where the client did build and tune its own transport.
func TestWithMetricsKeepsOwnedTransport(t *testing.T) {
	server, _ := newConnCountingServer(t, okHandler)

	core := NewRESTClientCore(server.URL, zerolog.Nop(), WithMetrics())
	require.NotNil(t, core.transport)
	require.Equal(t, defaultMaxIdleConnsPerHost, core.transport.MaxIdleConnsPerHost)
	require.NotSame(t, core.transport, core.httpClient.Transport)

	err := core.executeRequest(context.Background(), http.MethodGet, "/health", nil, nil, nil, nil, 0)
	require.NoError(t, err)
}

func TestWithTransportUsesSuppliedRoundTripper(t *testing.T) {
	server, _ := newConnCountingServer(t, okHandler)

	supplied := &countingRoundTripper{base: http.DefaultTransport}

	core := NewRESTClientCore(server.URL, zerolog.Nop(), WithTransport(supplied), WithConnectionTimeout(time.Second))
	require.Nil(t, core.transport, "a supplied round tripper is borrowed, not owned")
	require.Same(t, supplied, core.httpClient.Transport)

	err := core.executeRequest(context.Background(), http.MethodGet, "/health", nil, nil, nil, nil, 0)
	require.NoError(t, err)
	require.Equal(t, int64(1), supplied.calls.Load())
}

func TestWithTransportIgnoresNil(t *testing.T) {
	core := NewRESTClientCore("http://example.invalid", zerolog.Nop(), WithTransport(nil))

	require.NotNil(t, core.transport)
	require.Same(t, core.transport, core.httpClient.Transport)
}

func TestWithConnectionTimeoutIgnoresNonPositiveValues(t *testing.T) {
	core := NewRESTClientCore("http://example.invalid", zerolog.Nop())
	dialBefore := core.transport.DialContext

	WithConnectionTimeout(0)(core)
	require.Equal(t, dialBefore == nil, core.transport.DialContext == nil)
	require.NotNil(t, core.transport.DialContext)
}

// TestRESTClientCoreReusesConnections drives a fixed number of workers through
// many sequential requests each. A connection can only be reused once it is
// idle, so the total number of accepted connections should stay close to the
// worker count; with the standard library's default of two idle connections per
// host, most requests would instead open a new one.
func TestRESTClientCoreReusesConnections(t *testing.T) {
	server, acceptedConns := newConnCountingServer(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	})

	core := NewRESTClientCore(server.URL, zerolog.Nop())

	const (
		workers   = 8
		perWorker = 20
	)
	errs := driveRequests(core, workers, perWorker)
	require.Len(t, errs, workers*perWorker)
	for _, err := range errs {
		require.NoError(t, err)
	}

	require.LessOrEqual(t, acceptedConns(), workers*3)
}

// TestRESTClientCoreReusesConnectionsOnErrorResponses covers the path that
// returns before the response body is decoded: unless the body is drained, the
// transport discards the connection and every error costs a new handshake.
func TestRESTClientCoreReusesConnectionsOnErrorResponses(t *testing.T) {
	// A gateway in front of the node can answer with a sizeable error page; the
	// body has to be consumed before the connection can be pooled again.
	errorBody := []byte(`{"code":2,"message":"` + strings.Repeat("upstream failure ", 4096) + `"}`)

	server, acceptedConns := newConnCountingServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(errorBody)
	})

	core := NewRESTClientCore(server.URL, zerolog.Nop())

	const (
		workers   = 8
		perWorker = 20
	)
	errs := driveRequests(core, workers, perWorker)
	require.Len(t, errs, workers*perWorker)
	for _, err := range errs {
		require.Error(t, err)
	}

	t.Logf("accepted %d connections for %d error responses", acceptedConns(), workers*perWorker)
	require.LessOrEqual(t, acceptedConns(), workers*3)
}

// newConnCountingServer starts a test server and returns a function reporting
// how many TCP connections it has accepted.
func newConnCountingServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, func() int) {
	t.Helper()

	var (
		mu    sync.Mutex
		conns int
	)

	server := httptest.NewUnstartedServer(handler)
	server.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state != http.StateNew {
			return
		}
		mu.Lock()
		conns++
		mu.Unlock()
	}
	server.Start()
	t.Cleanup(server.Close)

	return server, func() int {
		mu.Lock()
		defer mu.Unlock()
		return conns
	}
}

// driveRequests runs `workers` goroutines, each issuing `perWorker` sequential
// requests, and returns every result. Errors are collected rather than asserted
// so that assertions stay on the test's own goroutine.
func driveRequests(core *RESTClientCore, workers, perWorker int) []error {
	var (
		mu   sync.Mutex
		errs []error
		wg   sync.WaitGroup
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < perWorker; j++ {
				err := core.executeRequest(ctx, http.MethodGet, "/health", nil, nil, nil, nil, 0)
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	return errs
}
