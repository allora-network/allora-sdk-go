package httpstatus

import (
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type trackedBody struct {
	io.Reader
	closed atomic.Int32
}

func (b *trackedBody) Close() error {
	b.closed.Add(1)
	return nil
}

// TestErrorClosingResponseClosesOriginalBody verifies that replacing the body
// for the drain keeps ownership of the network body.
func TestErrorClosingResponseClosesOriginalBody(t *testing.T) {
	original := &trackedBody{Reader: strings.NewReader(`{"code":5,"message":"absent"}`)}
	resp := &http.Response{
		StatusCode: http.StatusNotFound,
		Status:     "404 Not Found",
		Body:       original,
	}

	_ = Error(resp, 1024, time.Second)
	require.NoError(t, resp.Body.Close())
	require.Equal(t, int32(1), original.closed.Load())
}

// racingBody blocks its read until the timeout callback closes it. The timer
// races the response-body replacement during many iterations of this test.
type racingBody struct {
	started    chan struct{}
	once       sync.Once
	done       chan struct{}
	closeStart chan struct{}
	closed     sync.Once
	allowClose chan struct{}
}

func newRacingBody() *racingBody {
	return &racingBody{
		started:    make(chan struct{}),
		done:       make(chan struct{}),
		closeStart: make(chan struct{}),
		allowClose: make(chan struct{}),
	}
}

func (b *racingBody) Read([]byte) (int, error) {
	b.once.Do(func() { close(b.started) })
	<-b.done
	return 0, io.EOF
}

func (b *racingBody) Close() error {
	b.closed.Do(func() {
		close(b.closeStart)
		<-b.allowClose
		close(b.done)
	})
	return nil
}

// TestErrorTimerAndBodyCloseAreSafe runs under -race to check the timer's
// access to the response body. The timeout must close the same body regardless
// of how the restored prefix wrapper is composed.
func TestErrorTimerAndBodyCloseAreSafe(t *testing.T) {
	for i := 0; i < 100; i++ {
		original := newRacingBody()
		resp := &http.Response{
			StatusCode: http.StatusServiceUnavailable,
			Status:     "503 Service Unavailable",
			Body:       original,
		}
		done := make(chan struct{})
		go func() {
			_ = Error(resp, 1024, time.Millisecond)
			close(done)
		}()
		<-original.closeStart
		close(original.allowClose)
		<-done
		require.NoError(t, resp.Body.Close())
	}
}
