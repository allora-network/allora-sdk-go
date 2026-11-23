package tmrpc

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"sync/atomic"

	"google.golang.org/grpc/status"

	"github.com/allora-network/allora-sdk-go/httpstatus"
)

type httpStatusRecorderKey struct{}

// httpStatusRecorder carries one RPC's HTTP status from the transport to the
// client method. CometBFT has no typed status field in its returned error, so
// this typed request-local record avoids extracting a status from error text.
type httpStatusRecorder struct {
	ctx        context.Context
	status     atomic.Int32
	bodyResult atomic.Pointer[httpBodyReadResult]
}

type httpBodyReadResult struct {
	err error
}

func withHTTPStatusRecorder(ctx context.Context) (context.Context, *httpStatusRecorder) {
	recorder := &httpStatusRecorder{ctx: ctx}
	return context.WithValue(ctx, httpStatusRecorderKey{}, recorder), recorder
}

func (r *httpStatusRecorder) error(err error) error {
	if err == nil {
		return nil
	}
	if ctxErr := r.ctx.Err(); ctxErr != nil {
		return ctxErr
	}
	if isHTTPTransportError(err) {
		return err
	}
	if result := r.bodyResult.Load(); result != nil && result.err != nil && result.err != io.EOF {
		return result.err
	}
	httpStatus := int(r.status.Load())
	if code, ok := httpstatus.GRPCCode(httpStatus); ok {
		return status.Errorf(code, "CometBFT HTTP response %d: %s", httpStatus, http.StatusText(httpStatus))
	}
	return err
}

// isHTTPTransportError distinguishes failures while sending or reading HTTP
// from an HTTP response whose body contains an application error. A recorded
// status must not mask a transport fault with a definitive gRPC code.
func isHTTPTransportError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, os.ErrDeadlineExceeded) ||
		errors.Is(err, http.ErrBodyReadAfterClose) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr)
}

// httpStatusTransport records the HTTP status and tracks body reads for the
// client method. It otherwise forwards each response unchanged.
type httpStatusTransport struct {
	next *http.Transport
}

// CloseIdleConnections releases connections held by the pooled transport.
func (t *httpStatusTransport) CloseIdleConnections() {
	t.next.CloseIdleConnections()
}

func (t *httpStatusTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.next.RoundTrip(req)
	if recorder, ok := req.Context().Value(httpStatusRecorderKey{}).(*httpStatusRecorder); ok &&
		resp != nil {
		recorder.status.Store(int32(resp.StatusCode))
		resp.Body = &httpBodyRecorder{ReadCloser: resp.Body, recorder: recorder}
	}
	return resp, err
}

// httpBodyRecorder exposes body-read failures to the status recorder.
type httpBodyRecorder struct {
	io.ReadCloser
	recorder *httpStatusRecorder
}

func (b *httpBodyRecorder) Read(dst []byte) (int, error) {
	n, err := b.ReadCloser.Read(dst)
	b.recorder.bodyResult.Store(&httpBodyReadResult{err: err})
	return n, err
}
