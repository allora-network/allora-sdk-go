package tmrpc

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
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

	transport, ok := client.Transport.(*http.Transport)
	require.True(t, ok)

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

// TestHTTPClientReusesConnections drives a fixed number of workers through many
// sequential RPCs each. A connection can only be reused once it is idle, so the
// total number of accepted connections should stay close to the worker count;
// with the standard library's default of two idle connections per host, most
// calls would instead open a new one.
func TestHTTPClientReusesConnections(t *testing.T) {
	server, acceptedConns, _ := newJSONRPCServer(t)

	client, err := NewHTTPClient(server.URL, server.URL+"/websocket", 5*time.Second, zerolog.Nop())
	require.NoError(t, err)
	t.Cleanup(func() {
		// The error is ignored: the websocket half of the client was never
		// started, and reporting that is not what this test is about.
		_ = client.Close()
	})

	const (
		workers   = 8
		perWorker = 20
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var (
		mu   sync.Mutex
		errs []error
		wg   sync.WaitGroup
	)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < perWorker; j++ {
				_, err := client.Status(ctx)
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	require.Len(t, errs, workers*perWorker)
	for _, err := range errs {
		require.NoError(t, err)
	}

	t.Logf("accepted %d connections for %d RPCs", acceptedConns(), workers*perWorker)
	require.LessOrEqual(t, acceptedConns(), workers*3)
}

// TestHTTPClientCloseReleasesIdleConnections pins down that Close hands the
// pooled sockets back instead of leaving them idling until the pool's idle
// timeout expires. Stopping the websocket client alone does not touch them.
func TestHTTPClientCloseReleasesIdleConnections(t *testing.T) {
	server, _, closedConns := newJSONRPCServer(t)

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

// newJSONRPCServer starts a test server answering CometBFT JSON-RPC calls with
// an empty result, and returns functions reporting how many TCP connections it
// has accepted and how many it has seen closed.
func newJSONRPCServer(t *testing.T) (*httptest.Server, func() int, func() int) {
	t.Helper()

	var conns, closed atomic.Int64

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

		// The client rejects a response whose id does not match its request's.
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":` + string(request.ID) + `,"result":{}}`))
	}))
	server.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		switch state {
		case http.StateNew:
			conns.Add(1)
		case http.StateClosed:
			closed.Add(1)
		}
	}
	server.Start()
	t.Cleanup(server.Close)

	return server,
		func() int { return int(conns.Load()) },
		func() int { return int(closed.Load()) }
}
