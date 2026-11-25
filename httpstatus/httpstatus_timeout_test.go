package httpstatus

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestErrorClientBodyReadTimeoutIsTransportError exercises REST through its
// actual transport. A response status that arrived before a body read failed
// must not become the node's definitive query answer.
func TestErrorClientBodyReadTimeoutIsTransportError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", "200")
		w.WriteHeader(http.StatusNotImplemented)
		_, _ = w.Write([]byte(`{"code":12,`))
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
		<-r.Context().Done()
	}))
	t.Cleanup(server.Close)
	client := &http.Client{Timeout: 20 * time.Millisecond}

	resp, err := client.Get(server.URL)
	require.NoError(t, err)
	t.Cleanup(func() { _ = resp.Body.Close() })

	err = Error(resp, 1024, time.Second)

	require.Error(t, err)
	require.NotEqual(t, codes.Unimplemented, status.Code(err), "got %v", err)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

// TestErrorDrainTimeoutRemainsTransportError verifies REST's bounded-drain
// timeout also prevents the HTTP status line from becoming a definitive code.
func TestErrorDrainTimeoutRemainsTransportError(t *testing.T) {
	body := newRacingBody()
	resp := &http.Response{
		StatusCode: http.StatusNotImplemented,
		Status:     "501 Not Implemented",
		Body:       body,
	}

	errCh := make(chan error, 1)
	go func() { errCh <- Error(resp, 1024, time.Millisecond) }()
	<-body.closeStart
	close(body.allowClose)
	err := <-errCh

	require.Error(t, err)
	require.NotEqual(t, codes.Unimplemented, status.Code(err), "got %v", err)
	require.ErrorIs(t, err, io.EOF)
	require.Equal(t, "read HTTP error response body: EOF", err.Error())
	require.NoError(t, resp.Body.Close())
}

// TestErrorMalformedTransportErrorPrecedesHTTPStatus verifies a failed JSON
// unmarshal is preserved as a read error instead of falling back to the status.
func TestErrorMalformedTransportErrorPrecedesHTTPStatus(t *testing.T) {
	readErr := context.DeadlineExceeded
	resp := &http.Response{
		StatusCode: http.StatusNotImplemented,
		Status:     "501 Not Implemented",
		Body:       &errorReadCloser{err: readErr},
	}

	err := Error(resp, 1024, time.Second)

	require.Error(t, err)
	require.NotEqual(t, codes.Unimplemented, status.Code(err))
	require.ErrorIs(t, err, readErr)
	require.NoError(t, resp.Body.Close())
}

type errorReadCloser struct {
	err error
}

func (b *errorReadCloser) Read([]byte) (int, error) { return 0, b.err }
func (b *errorReadCloser) Close() error             { return nil }
