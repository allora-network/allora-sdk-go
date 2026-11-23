// Package httpstatus translates HTTP responses into gRPC status errors.
package httpstatus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCCode returns the gRPC code for an HTTP status when grpc-gateway maps
// exactly one code to it and proxies do not commonly reuse its meaning.
func GRPCCode(httpStatus int) (codes.Code, bool) {
	switch httpStatus {
	case http.StatusUnauthorized:
		return codes.Unauthenticated, true
	case http.StatusForbidden:
		return codes.PermissionDenied, true
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted, true
	case http.StatusNotImplemented:
		return codes.Unimplemented, true
	case http.StatusServiceUnavailable:
		return codes.Unavailable, true
	case http.StatusGatewayTimeout:
		return codes.DeadlineExceeded, true
	default:
		return codes.Unknown, false
	}
}

// Error reads a bounded grpc-gateway status body and returns its gRPC status.
// Failed reads and drain timeouts return the transport error. A missing or
// invalid body uses the unambiguous HTTP-status mapping, or a plain error for
// an ambiguous status.
func Error(resp *http.Response, maxBytes int64, maxDuration time.Duration) error {
	msg := fmt.Sprintf("HTTP error %d: %s", resp.StatusCode, resp.Status)

	originalBody := resp.Body
	timerFired := atomic.Bool{}
	timer := time.AfterFunc(maxDuration, func() {
		timerFired.Store(true)
		_ = originalBody.Close()
	})
	body, err := io.ReadAll(io.LimitReader(originalBody, maxBytes))
	timedOut := timerFired.Load()
	timer.Stop()
	restoreBody(resp, originalBody, body)
	if err != nil || timedOut {
		if err == nil {
			err = io.EOF
		}
		return fmt.Errorf("read HTTP error response body: %w", err)
	}
	{
		var gatewayStatus struct {
			Code    *uint32 `json:"code"`
			Message string  `json:"message"`
		}
		if json.Unmarshal(body, &gatewayStatus) == nil && gatewayStatus.Code != nil &&
			*gatewayStatus.Code > uint32(codes.OK) && *gatewayStatus.Code <= uint32(codes.Unauthenticated) {
			return status.Error(codes.Code(*gatewayStatus.Code), gatewayStatus.Message)
		}
	}

	if code, ok := GRPCCode(resp.StatusCode); ok {
		return status.Error(code, msg)
	}
	return fmt.Errorf("%s", msg)
}

type prefixReadCloser struct {
	io.Reader
	closeOriginal io.Closer
}

func (b *prefixReadCloser) Close() error {
	return b.closeOriginal.Close()
}

func restoreBody(resp *http.Response, original io.ReadCloser, body []byte) {
	resp.Body = &prefixReadCloser{
		Reader:        bytes.NewReader(body),
		closeOriginal: original,
	}
}
