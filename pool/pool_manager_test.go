package pool

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/allora-network/allora-sdk-go/config"
)

// fakeParticipant is a minimal PoolParticipant for pool tests.
type fakeParticipant struct {
	url string
}

func (f *fakeParticipant) Close() error                          { return nil }
func (f *fakeParticipant) GetEndpointURL() string                { return f.url }
func (f *fakeParticipant) GetProtocol() config.Protocol          { return config.ProtocolGRPC }
func (f *fakeParticipant) HealthCheck(ctx context.Context) error { return nil }

// TestExecuteWithRetryPerAttemptDeadline pins that each attempt receives its
// own context deadline, so a slow first endpoint cannot burn the caller's
// context and starve the remaining endpoints (the Oct 15 "both endpoints dead
// simultaneously" cascade). Mutation: sharing the caller's ctx across attempts
// makes the second attempt inherit an expired deadline and this test fails.
func TestExecuteWithRetryPerAttemptDeadline(t *testing.T) {
	logger := zerolog.Nop()
	clients := []*fakeParticipant{{url: "grpc://a:1"}, {url: "grpc://b:1"}}
	mgr := NewClientPoolManager[*fakeParticipant](clients, logger)
	defer mgr.Close()
	mgr.SetRequestTimeout(60 * time.Millisecond)

	// Caller ctx outlives a single attempt timeout; if the attempt deadline
	// were the caller's (shared), the first slow attempt would not expire on
	// its own and the second attempt would see the budget already consumed.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var attemptDeadlines []time.Duration
	sawDeadlineExceeded := false

	_, err := ExecuteWithRetry(ctx, mgr, &logger,
		func(attemptCtx context.Context, c *fakeParticipant) (struct{}, error) {
			// Record how long THIS attempt is allowed to run.
			if dl, ok := attemptCtx.Deadline(); ok {
				attemptDeadlines = append(attemptDeadlines, time.Until(dl))
			}
			// Block on the attempt ctx: a per-attempt deadline fires quickly;
			// a shared caller deadline (5s) would not, failing the test bound.
			select {
			case <-attemptCtx.Done():
				if errors.Is(attemptCtx.Err(), context.DeadlineExceeded) {
					sawDeadlineExceeded = true
				}
				return struct{}{}, attemptCtx.Err()
			case <-time.After(2 * time.Second):
				return struct{}{}, nil
			}
		})

	require.Error(t, err)
	require.True(t, sawDeadlineExceeded,
		"an attempt must be terminated by its own per-attempt deadline")
	require.NotEmpty(t, attemptDeadlines)
	for i, d := range attemptDeadlines {
		require.Lessf(t, d, time.Second,
			"attempt %d must have a short per-attempt deadline, got %v (shared caller ctx?)", i, d)
	}
}

// TestExecuteWithRetryCallerCancelPropagates ensures the per-attempt timeout
// never swallows an explicit caller cancellation.
func TestExecuteWithRetryCallerCancelPropagates(t *testing.T) {
	logger := zerolog.Nop()
	clients := []*fakeParticipant{{url: "grpc://a:1"}}
	mgr := NewClientPoolManager[*fakeParticipant](clients, logger)
	defer mgr.Close()
	mgr.SetRequestTimeout(10 * time.Second) // long; caller cancel must win

	ctx, cancel := context.WithCancel(context.Background())

	started := make(chan struct{})
	go func() {
		<-started
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	_, err := ExecuteWithRetry(ctx, mgr, &logger,
		func(attemptCtx context.Context, c *fakeParticipant) (struct{}, error) {
			close(started)
			<-attemptCtx.Done()
			return struct{}{}, attemptCtx.Err()
		})

	require.Error(t, err)
	require.True(t, errors.Is(err, context.Canceled),
		"caller cancellation must propagate, got %v", err)
}

// TestExecuteWithRetryUnimplementedIsTerminal pins that codes.Unimplemented
// ends the request on the first endpoint. The endpoint answered; it simply does
// not serve the method, so the pool must not cool it, must not back it off, and
// must not ask the other endpoints the same question.
func TestExecuteWithRetryUnimplementedIsTerminal(t *testing.T) {
	logger := zerolog.Nop()
	clients := []*fakeParticipant{{url: "grpc://a:1"}, {url: "grpc://b:1"}}
	mgr := NewClientPoolManager[*fakeParticipant](clients, logger)
	defer mgr.Close()

	var calls int
	_, err := ExecuteWithRetry(context.Background(), mgr, &logger,
		func(_ context.Context, c *fakeParticipant) (struct{}, error) {
			calls++
			return struct{}{}, status.Error(codes.Unimplemented, "unknown method")
		})

	require.Error(t, err)
	require.Equal(t, codes.Unimplemented, status.Code(err),
		"the caller must be able to tell an unimplemented method from a broken endpoint, got %v", err)
	require.Equal(t, 1, calls, "an unimplemented method must not be retried on the other endpoints")

	health := mgr.GetHealthStatus()
	require.Equal(t, 2, health["active_clients"], "an unimplemented method must not cool a healthy endpoint")
	require.Equal(t, 0, health["cooling_clients"])
	for _, c := range clients {
		require.Zero(t, mgr.IsClientInBackoff(c), "an unimplemented method must not back an endpoint off")
	}
}

// TestExecuteWithRetryOtherCodesStillFailOver guards the other side of that
// rule: every other error class keeps its current behaviour, so a failing
// endpoint is still penalised and the next one is tried.
func TestExecuteWithRetryOtherCodesStillFailOver(t *testing.T) {
	for _, code := range []codes.Code{
		codes.Unavailable,
		codes.Internal,
		codes.DeadlineExceeded,
		codes.InvalidArgument,
		codes.NotFound,
		codes.Unknown,
	} {
		t.Run(code.String(), func(t *testing.T) {
			logger := zerolog.Nop()
			clients := []*fakeParticipant{{url: "grpc://a:1"}, {url: "grpc://b:1"}}
			mgr := NewClientPoolManager[*fakeParticipant](clients, logger)
			defer mgr.Close()

			var calls int
			_, err := ExecuteWithRetry(context.Background(), mgr, &logger,
				func(_ context.Context, c *fakeParticipant) (struct{}, error) {
					calls++
					return struct{}{}, status.Error(code, "boom")
				})

			require.Error(t, err)
			require.Equal(t, 2, calls, "%v must still be retried on the other endpoint", code)
			require.Equal(t, 0, mgr.GetHealthStatus()["active_clients"],
				"%v must still cool the failing endpoints", code)
		})
	}
}
