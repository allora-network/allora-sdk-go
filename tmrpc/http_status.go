package tmrpc

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"

	"google.golang.org/grpc/status"

	"github.com/allora-network/allora-sdk-go/httpstatus"
)

type httpStatusRecorderKey struct{}

// httpStatusRecorder carries one RPC's HTTP status from the transport to the
// client method. CometBFT has no typed status field in its returned error, so
// this typed request-local record avoids extracting a status from error text.
type httpStatusRecorder struct {
	ctx    context.Context
	status atomic.Int32
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
	if errors.Is(err, context.Canceled) {
		return context.Canceled
	}
	httpStatus := int(r.status.Load())
	if code, ok := httpstatus.GRPCCode(httpStatus); ok {
		return status.Errorf(code, "CometBFT HTTP response %d: %s", httpStatus, http.StatusText(httpStatus))
	}
	return err
}

// httpStatusTransport records the concrete net/http status for the RPC method
// to classify. It leaves the response and CometBFT's error untouched.
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
	}
	return resp, err
}
