package tmrpc

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestNewPooledHTTPClientReadTimeoutIsTransportError covers an HTTP response
// whose status line says 501 before its body stalls. The client timeout occurs
// while reading, so the caller must not receive the status-line code.
func TestNewPooledHTTPClientReadTimeoutIsTransportError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", "200")
		w.WriteHeader(http.StatusNotImplemented)
		_, _ = w.Write([]byte(`{"jsonrpc":"2.0",`))
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
		<-r.Context().Done()
	}))
	t.Cleanup(server.Close)

	client, err := NewHTTPClient(server.URL, server.URL+"/websocket", 20*time.Millisecond, zerolog.Nop())
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	_, err = client.Status(context.Background())

	require.Error(t, err)
	require.NotEqual(t, codes.Unimplemented, status.Code(err), "transport error was mapped to a definitive status: %v", err)
	var netErr net.Error
	require.ErrorAs(t, err, &netErr)
	require.True(t, netErr.Timeout())
}

// TestHTTPStatusRecorderPreservesTransportErrors covers errors CometBFT wraps
// during JSON body decoding. Such failures say nothing about the endpoint's
// query-method support and must continue through the endpoint-fault policy.
func TestHTTPStatusRecorderPreservesTransportErrors(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  error
	}{
		{"network error", &net.DNSError{IsTimeout: true}},
		{"caller deadline wrapped by transport", os.ErrDeadlineExceeded},
		{"body closed after read timer", http.ErrBodyReadAfterClose},
	} {
		t.Run(tc.name, func(t *testing.T) {
			recorder := &httpStatusRecorder{ctx: context.Background()}
			recorder.status.Store(http.StatusNotImplemented)
			err := recorder.error(tc.err)
			require.Exactly(t, tc.err, err)
			require.NotEqual(t, codes.Unimplemented, status.Code(err))
			require.False(t, errors.Is(err, context.Canceled))
		})
	}
}
