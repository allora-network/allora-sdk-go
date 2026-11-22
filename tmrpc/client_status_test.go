package tmrpc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestNewHTTPClientMapsHTTPStatusRuns verifies that status codes rejected by
// CometBFT reach the caller with the corresponding gRPC code.
func TestNewHTTPClientMapsHTTPStatusRuns(t *testing.T) {
	for _, tc := range []struct {
		name       string
		httpStatus int
		wantCode   codes.Code
	}{
		{"unauthorized", http.StatusUnauthorized, codes.Unauthenticated},
		{"forbidden", http.StatusForbidden, codes.PermissionDenied},
		{"too many requests", http.StatusTooManyRequests, codes.ResourceExhausted},
		{"not implemented", http.StatusNotImplemented, codes.Unimplemented},
		{"service unavailable", http.StatusServiceUnavailable, codes.Unavailable},
		{"gateway timeout", http.StatusGatewayTimeout, codes.DeadlineExceeded},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.httpStatus)
				_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":-1,"result":{}}`))
			}))
			t.Cleanup(server.Close)

			client, err := NewHTTPClient(server.URL, server.URL+"/websocket", time.Second, zerolog.Nop())
			require.NoError(t, err)
			t.Cleanup(func() { _ = client.Close() })

			_, err = client.Status(context.Background())

			require.Error(t, err)
			require.Equal(t, tc.wantCode, status.Code(err), "got %v", err)
		})
	}
}
