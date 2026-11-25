package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestExecuteRequestReadsGatewayStatusBody verifies that the response body
// takes precedence over the less specific HTTP status.
func TestExecuteRequestReadsGatewayStatusBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code":5,"message":"mapped state does not exist"}`))
	}))
	t.Cleanup(server.Close)

	core := NewRESTClientCore(server.URL, zerolog.Nop())
	t.Cleanup(func() { _ = core.Close() })
	start := time.Now()

	err := core.executeRequest(context.Background(), http.MethodGet, "/health", nil, nil, nil, nil, 0)

	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err), "got %v", err)
	require.Equal(t, "mapped state does not exist", status.Convert(err).Message())
	require.Less(t, time.Since(start), 100*time.Millisecond)
}

// TestExecuteRequestRejectsInvalidGatewayStatusBody verifies that a body
// without a valid numeric status code falls back to the HTTP status line.
func TestExecuteRequestRejectsInvalidGatewayStatusBody(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
	}{
		{"message only", `{"message":"proxy generated"}`},
		{"string code", `{"code":"8","message":"rate limited"}`},
		{"invalid code", `{"code":99,"message":"unknown code"}`},
		{"malformed json", `{"code":5,`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(tc.body))
			}))
			t.Cleanup(server.Close)

			core := NewRESTClientCore(server.URL, zerolog.Nop())
			t.Cleanup(func() { _ = core.Close() })

			err := core.executeRequest(context.Background(), http.MethodGet, "/health", nil, nil, nil, nil, 0)

			require.Error(t, err)
			require.Equal(t, codes.Unknown, status.Code(err), "got %v", err)
			require.Contains(t, err.Error(), "HTTP error 404")
		})
	}
}

// TestExecuteRequestBoundsGatewayStatusBody prevents a large error body from
// being buffered in full. Truncated JSON falls back to the HTTP status mapping.
func TestExecuteRequestBoundsGatewayStatusBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"code":14,"message":"` + strings.Repeat("x", maxDrainBytes) + `"}`))
	}))
	t.Cleanup(server.Close)

	core := NewRESTClientCore(server.URL, zerolog.Nop())
	t.Cleanup(func() { _ = core.Close() })

	err := core.executeRequest(context.Background(), http.MethodGet, "/health", nil, nil, nil, nil, 0)

	require.Error(t, err)
	require.Equal(t, codes.Unknown, status.Code(err), "the capped body must be treated as incomplete")
}
