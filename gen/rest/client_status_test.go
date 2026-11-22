package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestExecuteRequestMapsUnambiguousStatuses pins the HTTP-to-gRPC translation
// for every status grpc-gateway derives from exactly one gRPC code. The pool
// decides what to do with an error from its gRPC code, so a status that reaches
// the caller as Unknown is treated as an endpoint fault; the statuses here must
// arrive with the code the node meant.
func TestExecuteRequestMapsUnambiguousStatuses(t *testing.T) {
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
			}))
			t.Cleanup(server.Close)

			core := NewRESTClientCore(server.URL, zerolog.Nop())
			t.Cleanup(func() { _ = core.Close() })

			err := core.executeRequest(context.Background(), http.MethodGet, "/health", nil, nil, nil, nil, 0)

			require.Error(t, err)
			require.Equal(t, tc.wantCode, status.Code(err), "got %v", err)
			require.Contains(t, err.Error(), "HTTP error",
				"the HTTP status must stay in the message, got %v", err)
		})
	}
}

// TestExecuteRequestLeavesAmbiguousStatusesUncoded pins that a status
// grpc-gateway derives from more than one gRPC code, or that a proxy reuses for
// its own meaning, keeps its plain-error shape. The pool treats an uncoded
// error as an endpoint fault, which is the safe reading for a 500, a 502, or a
// 404 from a node whose routes sit under another prefix.
func TestExecuteRequestLeavesAmbiguousStatusesUncoded(t *testing.T) {
	for _, tc := range []struct {
		name       string
		httpStatus int
	}{
		{"bad request", http.StatusBadRequest},
		{"not found", http.StatusNotFound},
		{"request timeout", http.StatusRequestTimeout},
		{"conflict", http.StatusConflict},
		{"internal server error", http.StatusInternalServerError},
		{"bad gateway", http.StatusBadGateway},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.httpStatus)
			}))
			t.Cleanup(server.Close)

			core := NewRESTClientCore(server.URL, zerolog.Nop())
			t.Cleanup(func() { _ = core.Close() })

			err := core.executeRequest(context.Background(), http.MethodGet, "/health", nil, nil, nil, nil, 0)

			require.Error(t, err)
			require.Equal(t, codes.Unknown, status.Code(err),
				"an ambiguous status must not carry a gRPC code, got %v", err)
			require.Contains(t, err.Error(), "HTTP error")
		})
	}
}

// TestGrpcCodeForHTTPStatusCoversEveryStatus pins the mapping table itself, so
// a status added to or removed from it is a visible change.
func TestGrpcCodeForHTTPStatusCoversEveryStatus(t *testing.T) {
	mapped := map[int]codes.Code{
		http.StatusUnauthorized:       codes.Unauthenticated,
		http.StatusForbidden:          codes.PermissionDenied,
		http.StatusTooManyRequests:    codes.ResourceExhausted,
		http.StatusNotImplemented:     codes.Unimplemented,
		http.StatusServiceUnavailable: codes.Unavailable,
		http.StatusGatewayTimeout:     codes.DeadlineExceeded,
	}
	for httpStatus, wantCode := range mapped {
		code, ok := grpcCodeForHTTPStatus(httpStatus)
		require.True(t, ok, "HTTP %d must be mapped", httpStatus)
		require.Equal(t, wantCode, code, "HTTP %d", httpStatus)
	}

	for _, httpStatus := range []int{
		http.StatusOK,
		http.StatusBadRequest,
		http.StatusNotFound,
		http.StatusRequestTimeout,
		http.StatusConflict,
		http.StatusInternalServerError,
		http.StatusBadGateway,
	} {
		_, ok := grpcCodeForHTTPStatus(httpStatus)
		require.False(t, ok, "HTTP %d must stay unmapped", httpStatus)
	}
}
