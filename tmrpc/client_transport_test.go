package tmrpc

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestNewPooledHTTPClientConfiguresConnectionPool(t *testing.T) {
	client, err := newPooledHTTPClient("http://127.0.0.1:26657", 7*time.Second)
	require.NoError(t, err)
	require.Equal(t, 7*time.Second, client.Timeout)

	recorder, ok := client.Transport.(*httpStatusTransport)
	require.True(t, ok)
	transport := recorder.next

	// The intended sizes are spelled out rather than read back from the
	// constants that configure them, so that shrinking the pool towards the
	// standard library defaults fails here instead of silently redefining what
	// this test asserts.
	require.Equal(t, 512, transport.MaxIdleConns)
	require.Equal(t, 256, transport.MaxIdleConnsPerHost)
	require.Equal(t, 90*time.Second, transport.IdleConnTimeout)
}

func TestNewPooledHTTPClientLeavesTimeoutUnsetWhenNonPositive(t *testing.T) {
	client, err := newPooledHTTPClient("http://127.0.0.1:26657", 0)
	require.NoError(t, err)
	require.Zero(t, client.Timeout)
}

func TestNewPooledHTTPClientRejectsInvalidRemote(t *testing.T) {
	_, err := newPooledHTTPClient("::not-a-url::", time.Second)
	require.Error(t, err)
}

// TestHTTPClientReusesConnections checks that sequential RPCs reuse a pooled
// connection after the initial dial.
func TestHTTPClientReusesConnections(t *testing.T) {
	server, _ := newJSONRPCServer(t)

	client, err := NewHTTPClient(server.URL, server.URL+"/websocket", 5*time.Second, zerolog.Nop())
	require.NoError(t, err)
	t.Cleanup(func() {
		// The error is ignored: the websocket half of the client was never
		// started, and reporting that is not what this test is about.
		_ = client.Close()
	})

	const requests = 20

	result := executeTracedStatus(context.Background(), client)
	require.NoError(t, result.err)
	require.False(t, result.reused, "the first request should establish a connection")

	for i := 1; i < requests; i++ {
		result = executeTracedStatus(context.Background(), client)
		require.NoError(t, result.err)
		require.True(t, result.reused, "request did not reuse the pooled connection")
	}
}

// TestHTTPClientReusesConnectionsConcurrently checks that all concurrent
// connections return to the idle pool and serve a subsequent burst.
func TestHTTPClientReusesConnectionsConcurrently(t *testing.T) {
	const (
		workers = 8
		timeout = 5 * time.Second
	)

	var activeBarrier atomic.Pointer[requestBarrier]
	firstBarrier := newRequestBarrier(workers)
	activeBarrier.Store(firstBarrier)
	t.Cleanup(firstBarrier.release)

	server, _ := newJSONRPCServer(t, func(r *http.Request) {
		activeBarrier.Load().hold(r)
	})
	client, err := NewHTTPClient(server.URL, server.URL+"/websocket", 0, zerolog.Nop())
	require.NoError(t, err)
	t.Cleanup(func() {
		// The error is ignored: the websocket half of the client was never
		// started, and reporting that is not what this test is about.
		_ = client.Close()
	})

	idleConnections := make(chan struct{})
	var idleCount atomic.Int64
	trace := &httptrace.ClientTrace{
		PutIdleConn: func(error) {
			if idleCount.Add(1) == workers {
				close(idleConnections)
			}
		},
	}
	executeBurst := func(barrier *requestBarrier) []statusResult {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		results := make(chan statusResult, workers)
		for range workers {
			go func() {
				results <- executeTracedStatus(ctx, client, trace)
			}()
		}
		waitForSignal(t, barrier.arrived, timeout, "requests to reach the server")
		barrier.release()
		return waitForStatusResults(t, results, workers, timeout)
	}

	results := executeBurst(firstBarrier)
	for _, result := range results {
		require.NoError(t, result.err)
	}
	waitForSignal(t, idleConnections, timeout, "connections to return to the idle pool")

	secondBarrier := newRequestBarrier(workers)
	t.Cleanup(secondBarrier.release)
	activeBarrier.Store(secondBarrier)
	results = executeBurst(secondBarrier)
	for i, result := range results {
		require.NoError(t, result.err)
		require.True(t, result.reused, "second-burst request %d did not reuse a pooled connection", i)
	}
}

type requestBarrier struct {
	arrived     chan struct{}
	released    chan struct{}
	remaining   atomic.Int64
	releaseOnce sync.Once
}

func newRequestBarrier(workers int) *requestBarrier {
	barrier := &requestBarrier{
		arrived:  make(chan struct{}),
		released: make(chan struct{}),
	}
	barrier.remaining.Store(int64(workers))
	return barrier
}

func (b *requestBarrier) hold(r *http.Request) {
	if b.remaining.Add(-1) == 0 {
		close(b.arrived)
	}
	select {
	case <-b.released:
	case <-r.Context().Done():
	}
}

func (b *requestBarrier) release() {
	b.releaseOnce.Do(func() { close(b.released) })
}

func waitForStatusResults(t *testing.T, results <-chan statusResult, workers int, timeout time.Duration) []statusResult {
	t.Helper()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	collected := make([]statusResult, 0, workers)
	for len(collected) < workers {
		select {
		case result := <-results:
			collected = append(collected, result)
		case <-timer.C:
			t.Fatal("timed out waiting for request results")
		}
	}
	return collected
}

func waitForSignal(t *testing.T, signal <-chan struct{}, timeout time.Duration, description string) {
	t.Helper()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-signal:
	case <-timer.C:
		t.Fatalf("timed out waiting for %s", description)
	}
}

// TestHTTPClientCloseReleasesIdleConnections pins down that Close hands the
// pooled sockets back instead of leaving them idling until the pool's idle
// timeout expires. Stopping the websocket client alone does not touch them.
func TestHTTPClientCloseReleasesIdleConnections(t *testing.T) {
	server, closedConns := newJSONRPCServer(t)

	client, err := NewHTTPClient(server.URL, server.URL+"/websocket", 5*time.Second, zerolog.Nop())
	require.NoError(t, err)

	_, err = client.Status(context.Background())
	require.NoError(t, err)
	require.Zero(t, closedConns(), "the connection should be idling in the pool after a completed RPC")

	// The returned error is ignored: the websocket half of the client was never
	// started, and the pool has to be released either way.
	_ = client.Close()
	require.Eventually(t, func() bool { return closedConns() > 0 }, time.Second, 5*time.Millisecond)
}

type statusResult struct {
	err    error
	reused bool
}

func executeTracedStatus(ctx context.Context, client Client, traces ...*httptrace.ClientTrace) statusResult {
	result := statusResult{}
	trace := &httptrace.ClientTrace{
		GotConn: func(info httptrace.GotConnInfo) {
			result.reused = info.Reused
		},
	}
	if len(traces) > 0 {
		trace = &httptrace.ClientTrace{
			GotConn:     trace.GotConn,
			PutIdleConn: traces[0].PutIdleConn,
		}
	}
	ctx = httptrace.WithClientTrace(ctx, trace)
	_, result.err = client.Status(ctx)
	return result
}

// newJSONRPCServer starts a test server answering CometBFT JSON-RPC calls with
// an empty result and returns a function reporting how many TCP connections it
// has seen closed.
func newJSONRPCServer(t *testing.T, hooks ...func(*http.Request)) (*httptest.Server, func() int) {
	t.Helper()

	var closed atomic.Int64

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var request struct {
			ID json.RawMessage `json:"id"`
		}
		if err := json.Unmarshal(body, &request); err != nil || len(request.ID) == 0 {
			http.Error(w, "malformed JSON-RPC request", http.StatusBadRequest)
			return
		}

		for _, hook := range hooks {
			hook(r)
		}

		// The client rejects a response whose id does not match its request's.
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":` + string(request.ID) + `,"result":{}}`))
	}))
	server.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateClosed {
			closed.Add(1)
		}
	}
	server.Start()
	t.Cleanup(server.Close)

	return server, func() int { return int(closed.Load()) }
}
