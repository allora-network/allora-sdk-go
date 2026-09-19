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

// TestRESTClientCloseReleasesIdleConnections pins down that a client which owns
// its transport hands the sockets back on Close rather than leaving them idling
// until the pool's idle timeout expires.
func TestRESTClientCloseReleasesIdleConnections(t *testing.T) {
	server, closedConns := newConnClosingServer(t, okHandler)

	client := NewRESTClient(server.URL, zerolog.Nop())
	err := client.core.executeRequest(context.Background(), http.MethodGet, "/health", nil, nil, nil, nil, 0)
	require.NoError(t, err)
	require.Zero(t, closedConns(), "the connection should be idling in the pool after a completed request")

	require.NoError(t, client.Close())
	require.Eventually(t, func() bool { return closedConns() > 0 }, time.Second, 5*time.Millisecond)
}

// TestRESTClientCloseLeavesBorrowedTransportAlone covers the other side of that
// rule: a transport handed to the client may be shared, so closing one client
// must not tear down connections another user still wants.
func TestRESTClientCloseLeavesBorrowedTransportAlone(t *testing.T) {
	server, acceptedConns := newConnCountingServer(t, okHandler)

	shared := &http.Transport{}
	t.Cleanup(shared.CloseIdleConnections)

	client := NewRESTClient(server.URL, zerolog.Nop(), WithTransport(shared))
	err := client.core.executeRequest(context.Background(), http.MethodGet, "/health", nil, nil, nil, nil, 0)
	require.NoError(t, err)
	require.NoError(t, client.Close())

	// The pooled connection survives, so the next user of the shared transport
	// reaches the server without a second handshake.
	other := NewRESTClient(server.URL, zerolog.Nop(), WithTransport(shared))
	err = other.core.executeRequest(context.Background(), http.MethodGet, "/health", nil, nil, nil, nil, 0)
	require.NoError(t, err)

	require.Equal(t, 1, acceptedConns())
}

// TestRESTClientCoreBoundsBodyDrainOnStalledPeer covers an endpoint that answers
// with an error, starts the body and then stops sending. Draining that body is
// worth a moment (it buys the connection back) but not the client timeout: the
// caller needs the error promptly so a pool can try another endpoint.
func TestRESTClientCoreBoundsBodyDrainOnStalledPeer(t *testing.T) {
	release := make(chan struct{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"upstream failure`))
		w.(http.Flusher).Flush()
		<-release
	}))
	// Cleanups run last registered first, so the handler is released before the
	// server waits for it.
	t.Cleanup(server.Close)
	t.Cleanup(func() { close(release) })

	core := NewRESTClientCore(server.URL, zerolog.Nop())
	t.Cleanup(func() { _ = core.Close() })

	start := time.Now()
	err := core.executeRequest(context.Background(), http.MethodGet, "/health", nil, nil, nil, nil, 0)
	elapsed := time.Since(start)

	require.Error(t, err)
	require.Contains(t, err.Error(), "HTTP error 500")
	require.Less(t, elapsed, time.Second, "the drain must not wait out the client timeout")
}

// TestRESTClientCoreDrainDoesNotSlowSuccessfulRequests guards the other side of
// the bound: a body the client has already read to EOF should cost nothing.
func TestRESTClientCoreDrainDoesNotSlowSuccessfulRequests(t *testing.T) {
	server, _ := newConnCountingServer(t, okHandler)

	core := NewRESTClientCore(server.URL, zerolog.Nop())
	t.Cleanup(func() { _ = core.Close() })

	start := time.Now()
	for i := 0; i < 20; i++ {
		err := core.executeRequest(context.Background(), http.MethodGet, "/health", nil, nil, nil, nil, 0)
		require.NoError(t, err)
	}
	require.Less(t, time.Since(start), maxDrainDuration, "successful requests should not pay the drain bound")
}

// newConnClosingServer starts a test server and returns a function reporting how
// many TCP connections it has seen closed.
func newConnClosingServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, func() int) {
	t.Helper()

	var closed atomic.Int64

	server := httptest.NewUnstartedServer(handler)
	server.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateClosed {
			closed.Add(1)
		}
	}
	server.Start()
	t.Cleanup(server.Close)

	return server, func() int { return int(closed.Load()) }
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
