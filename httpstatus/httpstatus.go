// Package httpstatus translates HTTP responses into gRPC status errors.
package httpstatus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
// A missing, invalid, or unreadable body falls back to the unambiguous
// HTTP-status mapping, or to a plain error for an ambiguous status.
func Error(resp *http.Response, maxBytes int64, maxDuration time.Duration) error {
	msg := fmt.Sprintf("HTTP error %d: %s", resp.StatusCode, resp.Status)

	originalBody := resp.Body
	timer := time.AfterFunc(maxDuration, func() { _ = originalBody.Close() })
	body, err := io.ReadAll(io.LimitReader(originalBody, maxBytes))
	timer.Stop()
	restoreBody(resp, originalBody, body)
	if err == nil {
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
