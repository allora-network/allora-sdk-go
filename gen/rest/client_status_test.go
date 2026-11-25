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

// TestExecuteRequestMapsNotImplementedToUnimplemented pins that a 501 from a
// grpc-gateway node reaches the caller as codes.Unimplemented. The node reports
// an unimplemented method as HTTP 501 with a body of {"code":12,...}; without
// the translation the caller sees an opaque string and cannot tell "this node
// does not serve this method" from "this node is broken".
func TestExecuteRequestMapsNotImplementedToUnimplemented(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		_, _ = w.Write([]byte(`{"code":12,"message":"Not Implemented"}`))
	}))
	t.Cleanup(server.Close)

	core := NewRESTClientCore(server.URL, zerolog.Nop())
	t.Cleanup(func() { _ = core.Close() })

	err := core.executeRequest(context.Background(), http.MethodGet, "/emissions/v10/is_topic_forecaster_whitelist_enabled/1", nil, nil, nil, nil, 0)

	require.Error(t, err)
	require.Equal(t, codes.Unimplemented, status.Code(err), "got %v", err)
	require.Contains(t, err.Error(), "HTTP error 501", "the HTTP status must stay in the message")
}

// TestExecuteRequestLeavesOtherStatusesUncoded pins that every other non-200
// status keeps its current shape: a plain error with the status in the message
// and no gRPC code attached.
func TestExecuteRequestLeavesOtherStatusesUncoded(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
	}{
		{"bad request", http.StatusBadRequest},
		{"not found", http.StatusNotFound},
		{"internal server error", http.StatusInternalServerError},
		{"service unavailable", http.StatusServiceUnavailable},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.status)
			}))
			t.Cleanup(server.Close)

			core := NewRESTClientCore(server.URL, zerolog.Nop())
			t.Cleanup(func() { _ = core.Close() })

			err := core.executeRequest(context.Background(), http.MethodGet, "/health", nil, nil, nil, nil, 0)

			require.Error(t, err)
			require.Equal(t, codes.Unknown, status.Code(err),
				"only 501 carries a gRPC code, got %v", err)
			require.Contains(t, err.Error(), "HTTP error")
		})
	}
}
