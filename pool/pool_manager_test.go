package pool

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

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
