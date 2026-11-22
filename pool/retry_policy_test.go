package pool

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// retryCase is one row of the retry policy table: the error an endpoint
// returns, the decision the policy must reach, and the observable consequences
// of that decision.
type retryCase struct {
	name string
	err  error

	wantDecision retryDecision
	// wantAttempts is how many endpoints ExecuteWithRetry must call.
	wantAttempts int
	// wantCooled is how many endpoints must end up in the cooling pool.
	wantCooled int
	// wantBackoff is whether the endpoint that answered must be in backoff.
	wantBackoff bool
}

// TestClassifyRetryTable pins the decision for every gRPC status code, plus the
// non-status and context cases. Each row also asserts the observable effect on
// the pool, so a change to the classifier that does not reach ExecuteWithRetry
// fails here.
func TestClassifyRetryTable(t *testing.T) {
	cases := []retryCase{
		// State-dependent answers: another healthy endpoint may hold the answer.
		{"InvalidArgument", status.Error(codes.InvalidArgument, "vote not found in this state"), failOverOnly, 2, 0, false},
		{"NotFound", status.Error(codes.NotFound, "topic 5 not found"), failOverOnly, 2, 0, false},
		{"AlreadyExists", status.Error(codes.AlreadyExists, "exists"), failOverOnly, 2, 0, false},
		{"OutOfRange", status.Error(codes.OutOfRange, "out of range"), failOverOnly, 2, 0, false},
		// Chain upgrades are coordinated, so unsupported methods are terminal.
		{"Unimplemented", status.Error(codes.Unimplemented, "unknown method"), returnToCaller, 1, 0, false},

		// Healthy endpoint that may not be able to answer this request.
		{"Unknown", status.Error(codes.Unknown, "codespace sdk code 18: invalid request: failed to load state at height 5"), failOverOnly, 2, 0, false},
		{"Canceled", status.Error(codes.Canceled, "canceled by peer"), failOverOnly, 2, 0, false},

		// Endpoint at fault.
		{"DeadlineExceeded", status.Error(codes.DeadlineExceeded, "timeout"), failOverAndPenalise, 2, 2, true},
		{"PermissionDenied", status.Error(codes.PermissionDenied, "forbidden"), failOverAndPenalise, 2, 2, true},
		{"Unauthenticated", status.Error(codes.Unauthenticated, "unauthorized"), failOverAndPenalise, 2, 2, true},
		{"ResourceExhausted", status.Error(codes.ResourceExhausted, "rate limit exceeded"), failOverAndPenalise, 2, 2, true},
		{"FailedPrecondition", status.Error(codes.FailedPrecondition, "not ready"), failOverAndPenalise, 2, 2, true},
		{"Aborted", status.Error(codes.Aborted, "conflict"), failOverAndPenalise, 2, 2, true},
		{"Internal", status.Error(codes.Internal, "panic"), failOverAndPenalise, 2, 2, true},
		{"Unavailable", status.Error(codes.Unavailable, "connection refused"), failOverAndPenalise, 2, 2, true},
		{"DataLoss", status.Error(codes.DataLoss, "corrupt"), failOverAndPenalise, 2, 2, true},

		// A payload over the message size limit is not a rate limit: no
		// endpoint will accept it.
		{"ResourceExhausted/message too large", status.Error(codes.ResourceExhausted,
			"grpc: received message larger than max (5000000 vs. 4194304)"), returnToCaller, 1, 0, false},

		// Not a gRPC status: a local or transport fault.
		{"plain error", errors.New("dial tcp: connection refused"), failOverAndPenalise, 2, 2, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.wantDecision, classifyRetry(tc.err, context.Background()),
				"classifyRetry decision for %v", tc.err)

			logger := zerolog.Nop()
			clients := []*fakeParticipant{{url: "grpc://a:1"}, {url: "grpc://b:1"}}
			mgr := NewClientPoolManager[*fakeParticipant](clients, logger)
			defer mgr.Close()
			mgr.currentIndex = 0

			var calls int
			_, err := ExecuteWithRetry(context.Background(), mgr, &logger,
				func(_ context.Context, _ *fakeParticipant) (struct{}, error) {
					calls++
					return struct{}{}, tc.err
				})

			require.Error(t, err)
			require.Equal(t, tc.wantAttempts, calls, "attempts for %v", tc.err)

			health := mgr.GetHealthStatus()
			require.Equal(t, tc.wantCooled, health["cooling_clients"], "cooling pool for %v", tc.err)
			require.Equal(t, 2-tc.wantCooled, health["active_clients"], "active pool for %v", tc.err)

			for _, c := range clients {
				if tc.wantBackoff {
					require.NotZero(t, mgr.IsClientInBackoff(c), "endpoint %s must be backed off for %v", c.url, tc.err)
				} else {
					require.Zero(t, mgr.IsClientInBackoff(c), "endpoint %s must not be backed off for %v", c.url, tc.err)
				}
			}
		})
	}
}

// TestClassifyRetryCallerCancellation pins that a caller who cancels gets the
// cancellation back at once, whatever the endpoint reported. Retrying would
// ignore an explicit request to stop.
func TestClassifyRetryCallerCancellation(t *testing.T) {
	for _, err := range []error{
		status.Error(codes.Unavailable, "down"),
		status.Error(codes.DeadlineExceeded, "timeout"),
		status.Error(codes.Internal, "panic"),
		context.Canceled,
	} {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		require.Equal(t, returnToCaller, classifyRetry(err, ctx),
			"a cancelled caller must not be retried, got %v", err)
	}
}

// TestClassifyRetryAttemptDeadlineIsEndpointFault pins the other half of the
// context rule: with a live caller context, DeadlineExceeded can only come from
// the per-attempt deadline, so the endpoint that missed it is at fault.
func TestClassifyRetryAttemptDeadlineIsEndpointFault(t *testing.T) {
	require.Equal(t, failOverAndPenalise,
		classifyRetry(context.DeadlineExceeded, context.Background()))
	require.Equal(t, failOverAndPenalise,
		classifyRetry(status.Error(codes.DeadlineExceeded, "timeout"), context.Background()))
}

// TestExecuteWithRetryCallerCancelDoesNotCoolEndpoint pins that a caller
// cancellation leaves the pool untouched: the endpoint did nothing wrong.
func TestExecuteWithRetryCallerCancelDoesNotCoolEndpoint(t *testing.T) {
	logger := zerolog.Nop()
	clients := []*fakeParticipant{{url: "grpc://a:1"}, {url: "grpc://b:1"}}
	mgr := NewClientPoolManager[*fakeParticipant](clients, logger)
	defer mgr.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var calls int
	_, err := ExecuteWithRetry(ctx, mgr, &logger,
		func(attemptCtx context.Context, _ *fakeParticipant) (struct{}, error) {
			calls++
			return struct{}{}, attemptCtx.Err()
		})

	require.Error(t, err)
	require.Equal(t, 0, calls, "a cancelled caller must not start an RPC")
	require.Equal(t, 2, mgr.GetHealthStatus()["active_clients"],
		"a cancelled caller must not cool a healthy endpoint")
}

// TestExecuteWithRetryFailOverOnlyKeepsEndpointsActive runs both start orders.
// The first endpoint cannot serve the pruned height; the second can. Both stay
// active because the failure carries no health penalty.
func TestExecuteWithRetryFailOverOnlyKeepsEndpointsActive(t *testing.T) {
	for _, first := range []int{0, 1} {
		t.Run(fmt.Sprintf("start_%d", first), func(t *testing.T) {
			testFailOverOnlyKeepsEndpointsActive(t, first)
		})
	}
}

func testFailOverOnlyKeepsEndpointsActive(t *testing.T, first int) {
	logger := zerolog.Nop()
	clients := []*fakeParticipant{{url: "grpc://a:1"}, {url: "grpc://b:1"}}
	mgr := NewClientPoolManager[*fakeParticipant](clients, logger)
	defer mgr.Close()
	mgr.currentIndex = first

	var calls int
	result, err := ExecuteWithRetry(context.Background(), mgr, &logger,
		func(_ context.Context, c *fakeParticipant) (string, error) {
			calls++
			if calls == 1 {
				return "", status.Error(codes.Unknown,
					"codespace sdk code 18: invalid request: failed to load state at height 5")
			}
			return "ok", nil
		})

	require.NoError(t, err)
	require.Equal(t, "ok", result)
	require.Equal(t, 2, calls, "the second endpoint must be tried")

	health := mgr.GetHealthStatus()
	require.Equal(t, 2, health["active_clients"],
		"an endpoint that pruned the state is healthy and must stay active")
	require.Equal(t, 0, health["cooling_clients"])
	for _, c := range clients {
		require.Zero(t, mgr.IsClientInBackoff(c),
			"an endpoint that pruned the state must not be backed off")
	}
}

// TestExecuteWithRetryDefinitiveAnswerNamesEndpoint pins that the caller can
// tell which endpoint gave the definitive answer, and that the gRPC code
// survives the wrap.
func TestExecuteWithRetryDefinitiveAnswerNamesEndpoint(t *testing.T) {
	logger := zerolog.Nop()
	clients := []*fakeParticipant{{url: "grpc://a:1"}, {url: "grpc://b:1"}}
	mgr := NewClientPoolManager[*fakeParticipant](clients, logger)
	defer mgr.Close()

	_, err := ExecuteWithRetry(context.Background(), mgr, &logger,
		func(_ context.Context, _ *fakeParticipant) (struct{}, error) {
			return struct{}{}, status.Error(codes.InvalidArgument, "empty request")
		})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err),
		"the caller must see the code the endpoint reported, got %v", err)
	require.Contains(t, err.Error(), "grpc://",
		"the error must name the endpoint that answered, got %v", err)
}

// TestExpDelayGrowsWithAttempt pins that the backoff a penalised endpoint
// receives grows with each failure, so a persistently broken endpoint is tried
// less often. The jitter is a fraction of the delay, so each delay is compared
// against the previous one with the jitter band removed.
func TestExpDelayGrowsWithAttempt(t *testing.T) {
	logger := zerolog.Nop()
	mgr := NewClientPoolManager[*fakeParticipant]([]*fakeParticipant{{url: "grpc://a:1"}}, logger)
	defer mgr.Close()

	// jitterFrac is 0.1, so a delay is within ±10% of its nominal value.
	const jitterBand = 0.1
	prev := mgr.ExpDelay(0)
	require.Positive(t, prev)
	for attempt := 1; attempt <= 4; attempt++ {
		next := mgr.ExpDelay(attempt)
		require.Greater(t, float64(next), float64(prev)*(1+jitterBand),
			"attempt %d must back off for longer than attempt %d", attempt, attempt-1)
		prev = next
	}
	require.LessOrEqual(t, prev, mgr.maxDelay,
		"the backoff must stay within the configured maximum")
}

// TestApplyBackoffPenaltyIsIdempotentWithinWindow pins that a second penalty
// while the first is still running does not extend it, so a burst of failures
// against one endpoint does not lock it out for longer than one backoff.
func TestApplyBackoffPenaltyIsIdempotentWithinWindow(t *testing.T) {
	logger := zerolog.Nop()
	clients := []*fakeParticipant{{url: "grpc://a:1"}}
	mgr := NewClientPoolManager[*fakeParticipant](clients, logger)
	defer mgr.Close()

	first := mgr.ApplyBackoffPenalty(clients[0], mgr.ExpDelay(0))
	second := mgr.ApplyBackoffPenalty(clients[0], mgr.ExpDelay(3))

	require.Positive(t, first)
	require.LessOrEqual(t, second, first,
		"a penalty inside an active backoff window must not extend it")
}
