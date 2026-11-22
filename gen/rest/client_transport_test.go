package rest

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
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

// A typed nil is not equal to a nil interface, so without an explicit check it
// would be installed as the transport and panic on the first request.
func TestWithTransportIgnoresTypedNil(t *testing.T) {
	var typedNil *http.Transport

	core := NewRESTClientCore("http://example.invalid", zerolog.Nop(), WithTransport(typedNil))

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

// TestRESTClientCoreReusesConnections checks that sequential requests reuse a
// pooled connection after the initial dial.
func TestRESTClientCoreReusesConnections(t *testing.T) {
	server, _ := newConnCountingServer(t, okHandler)
	core := NewRESTClientCore(server.URL, zerolog.Nop())
	t.Cleanup(func() { _ = core.Close() })

	const requests = 20

	result := executeTracedRequest(context.Background(), core)
	require.NoError(t, result.err)
	require.False(t, result.reused, "the first request should establish a connection")

	for i := 1; i < requests; i++ {
		result = executeTracedRequest(context.Background(), core)
		require.NoError(t, result.err)
		require.True(t, result.reused, "request did not reuse the pooled connection")
	}
}

// TestRESTClientCoreReusesConnectionsOnErrorResponses covers the path that
// returns before the response body is decoded: unless the body is drained, the
// transport discards the connection and every error costs a new handshake.
//
// The reuse assertions assume the server delivers the whole body within
// maxDrainDuration and that the body fits in maxDrainBytes. A peer that stalls
// longer loses the connection, which is the behaviour
// TestRESTClientCoreBoundsBodyDrainOnStalledPeer pins down. httptest writes the
// body in one call, so the drain always reaches EOF here.
func TestRESTClientCoreReusesConnectionsOnErrorResponses(t *testing.T) {
	// net/http recovers a connection whose body the caller closed early only by
	// draining it itself, and it gives up past maxPostCloseReadBytes (256 KiB)
	// or maxPostCloseReadTime (50 ms). A body over that limit costs the
	// connection unless the caller reads it to EOF, which is what drainAndClose
	// does. httptest sends this body chunked, so ContentLength is -1 and the
	// byte limit, not the length check, is what applies.
	errorBody := []byte(`{"code":2,"message":"` + strings.Repeat("upstream failure ", 32768) + `"}`)

	server, _ := newConnCountingServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(errorBody)
	})
	core := NewRESTClientCore(server.URL, zerolog.Nop())
	t.Cleanup(func() { _ = core.Close() })

	const requests = 20

	result := executeTracedRequest(context.Background(), core)
	require.Error(t, result.err)
	require.False(t, result.reused, "the first request should establish a connection")

	for i := 1; i < requests; i++ {
		result = executeTracedRequest(context.Background(), core)
		require.Error(t, result.err)
		require.True(t, result.reused, "error response body did not return the connection to the pool")
	}
}

// TestRESTClientCoreReusesConnectionsUnderConcurrency checks that a burst of
// concurrent requests leaves the pool holding a connection for each of them, so
// that a second burst reuses all of them instead of dialing again.
//
// The server holds each burst until all N requests are in flight, so the N
// connections exist at the same time and none can be handed to a later request
// in the same burst. Every connection is therefore idle in the pool when the
// burst ends, and the second burst reuses all N only if the pool can hold N of
// them: MaxIdleConnsPerHost must be at least N. Below that the transport evicts
// the surplus (tryPutIdleConn returns errTooManyIdleHost) and those requests
// dial again. With the standard library's default of two, six of the eight
// connections are dropped between bursts.
func TestRESTClientCoreReusesConnectionsUnderConcurrency(t *testing.T) {
	const burst = 8

	first := newBurstBarrier(burst)
	second := newBurstBarrier(burst)

	var current atomic.Pointer[burstBarrier]
	current.Store(first)

	server, _ := newConnCountingServer(t, func(w http.ResponseWriter, _ *http.Request) {
		current.Load().wait()
		_, _ = w.Write([]byte(`{}`))
	})
	// Cleanups run last registered first, so this runs after the client is
	// closed and before the server is, unblocking any handler still waiting on
	// a barrier that never filled.
	t.Cleanup(func() {
		first.release()
		second.release()
	})

	core := NewRESTClientCore(server.URL, zerolog.Nop())
	t.Cleanup(func() { _ = core.Close() })

	// The pool starts empty, so the first burst dials every connection it holds.
	for i, info := range runBurst(t, core, burst, first) {
		require.Falsef(t, info.Reused, "request %d in the first burst reused a connection from an empty pool", i)
	}

	current.Store(second)
	for i, info := range runBurst(t, core, burst, second) {
		require.Truef(t, info.Reused, "request %d in the second burst dialed a new connection", i)
	}
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

type requestResult struct {
	err    error
	reused bool
}

// burstBarrier holds a burst of requests at the server until all of them have
// arrived, then releases them together. Holding the whole burst in flight at
// once is what forces the connection pool to hold one connection per request.
type burstBarrier struct {
	size     int
	mu       sync.Mutex
	arrived  int
	closed   bool
	released chan struct{}
}

func newBurstBarrier(size int) *burstBarrier {
	return &burstBarrier{size: size, released: make(chan struct{})}
}

// wait blocks until size requests have called it, then releases all of them.
func (b *burstBarrier) wait() {
	b.mu.Lock()
	b.arrived++
	if b.arrived == b.size {
		b.closeLocked()
	}
	b.mu.Unlock()

	<-b.released
}

// release unblocks every waiter without waiting for the burst to fill, so a
// request that never arrives cannot leave the server's handlers stuck.
func (b *burstBarrier) release() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.closeLocked()
}

func (b *burstBarrier) closeLocked() {
	if !b.closed {
		b.closed = true
		close(b.released)
	}
}

// runBurst issues size concurrent requests through core and returns the
// connection info each request received. The requests are held at the server
// until all of them are in flight, so the pool must hold size connections at
// once. A request that fails releases the barrier, so the failure is reported
// instead of the rest of the burst waiting out the context deadline.
func runBurst(t *testing.T, core *RESTClientCore, size int, barrier *burstBarrier) []httptrace.GotConnInfo {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	var (
		wg    sync.WaitGroup
		infos = make([]httptrace.GotConnInfo, size)
		errs  = make([]error, size)
	)

	wg.Add(size)
	for i := 0; i < size; i++ {
		go func() {
			defer wg.Done()

			trace := &httptrace.ClientTrace{
				GotConn: func(info httptrace.GotConnInfo) { infos[i] = info },
			}
			errs[i] = core.executeRequest(
				httptrace.WithClientTrace(ctx, trace),
				http.MethodGet, "/health", nil, nil, nil, nil, 0,
			)
			if errs[i] != nil {
				barrier.release()
			}
		}()
	}
	wg.Wait()

	for i, err := range errs {
		require.NoErrorf(t, err, "request %d in the burst failed", i)
	}
	return infos
}

func executeTracedRequest(ctx context.Context, core *RESTClientCore) requestResult {
	result := requestResult{}
	trace := &httptrace.ClientTrace{
		GotConn: func(info httptrace.GotConnInfo) {
			result.reused = info.Reused
		},
	}
	ctx = httptrace.WithClientTrace(ctx, trace)
	result.err = core.executeRequest(ctx, http.MethodGet, "/health", nil, nil, nil, nil, 0)
	return result
}
