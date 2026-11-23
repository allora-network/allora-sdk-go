package tmrpc

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// TestNewHTTPClientCancellationOutranksHTTPStatus exercises a caller that ends
// while CometBFT is reading a response body. The observed HTTP status must not
// hide the caller's cancellation or deadline.
func TestNewHTTPClientCancellationOutranksHTTPStatus(t *testing.T) {
	for _, tc := range []struct {
		name     string
		deadline time.Duration
		cancel   bool
		wantErr  error
	}{
		{"cancellation", 0, true, context.Canceled},
		{"deadline", 20 * time.Millisecond, false, context.DeadlineExceeded},
	} {
		t.Run(tc.name, func(t *testing.T) {
			started := make(chan struct{})
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotImplemented)
				if flusher, ok := w.(http.Flusher); ok {
					flusher.Flush()
				}
				close(started)
				<-r.Context().Done()
			}))
			t.Cleanup(server.Close)

			client, err := NewHTTPClient(server.URL, server.URL+"/websocket", time.Second, zerolog.Nop())
			require.NoError(t, err)
			t.Cleanup(func() { _ = client.Close() })

			ctx := context.Background()
			cancel := func() {}
			if tc.cancel {
				ctx, cancel = context.WithCancel(ctx)
			} else {
				ctx, cancel = context.WithTimeout(ctx, tc.deadline)
			}
			defer cancel()

			if tc.cancel {
				go func() {
					<-started
					cancel()
				}()
			}
			_, err = client.Status(ctx)

			require.Error(t, err)
			require.True(t, errors.Is(err, tc.wantErr), "got %v", err)
		})
	}
}
