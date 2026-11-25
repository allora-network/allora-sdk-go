package httpstatus

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestErrorReadsGatewayStatus verifies that a numeric gateway code and message
// take precedence over an ambiguous HTTP status.
func TestErrorReadsGatewayStatus(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Status:     "429 Too Many Requests",
		Body: io.NopCloser(strings.NewReader(
			`{"code":8,"message":"grpc: received message larger than max (5 vs. 4)"}`)),
	}

	err := Error(resp, 1024, time.Second)

	require.Equal(t, codes.ResourceExhausted, status.Code(err))
	require.Contains(t, status.Convert(err).Message(), "larger than max")
}

// TestErrorFallsBack verifies the HTTP mapping is used only when no valid
// gateway status is present.
func TestErrorFallsBack(t *testing.T) {
	for _, tc := range []struct {
		name       string
		body       string
		httpStatus int
		wantCode   codes.Code
	}{
		{"bare proxy rate limit", "", http.StatusTooManyRequests, codes.ResourceExhausted},
		{"malformed gateway status", `{"code":`, http.StatusNotImplemented, codes.Unimplemented},
		{"ambiguous status", "", http.StatusBadRequest, codes.Unknown},
		{"code out of range", fmt.Sprintf(`{"code":%d}`, uint32(codes.Unauthenticated)+1), http.StatusTooManyRequests, codes.ResourceExhausted},
		{"success code in error response", `{"code":0,"message":"not an error"}`, http.StatusBadRequest, codes.Unknown},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: tc.httpStatus,
				Status:     http.StatusText(tc.httpStatus),
				Body:       io.NopCloser(strings.NewReader(tc.body)),
			}

			err := Error(resp, 1024, time.Second)

			require.Error(t, err)
			require.Equal(t, tc.wantCode, status.Code(err))
			require.Contains(t, err.Error(), "HTTP error")
		})
	}
}

// TestErrorRestoresBodyForDrain verifies that reading the status does not
// bypass the caller's body-drain accounting.
func TestErrorRestoresBodyForDrain(t *testing.T) {
	body := `{"code":14,"message":"down"}` + strings.Repeat("x", 1024)
	resp := &http.Response{
		StatusCode: http.StatusServiceUnavailable,
		Status:     "503 Service Unavailable",
		Body:       io.NopCloser(strings.NewReader(body)),
	}

	_ = Error(resp, 16, time.Second)
	remaining, err := io.ReadAll(resp.Body)

	require.NoError(t, err)
	require.Equal(t, body[:16], string(remaining))
}
