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

// TestExecuteWithRetryDefinitiveErrorsFailOverWithoutPenalty covers healthy
// nodes that give different answers for the same query. A missing answer from
// the first endpoint must not hide an answer held by the second.
func TestExecuteWithRetryDefinitiveErrorsFailOverWithoutPenalty(t *testing.T) {
	for _, tc := range []struct {
		name string
		code codes.Code
	}{
		{"not found", codes.NotFound},
		{"invalid argument", codes.InvalidArgument},
		{"already exists", codes.AlreadyExists},
		{"out of range", codes.OutOfRange},
	} {
		t.Run(tc.name, func(t *testing.T) {
			logger := zerolog.Nop()
			clients := []*fakeParticipant{{url: "grpc://a:1"}, {url: "grpc://b:1"}}
			mgr := NewClientPoolManager[*fakeParticipant](clients, logger)
			defer mgr.Close()
			mgr.currentIndex = 0

			var callsA, callsB int
			result, err := ExecuteWithRetry(context.Background(), mgr, &logger,
				func(_ context.Context, c *fakeParticipant) (string, error) {
					if c.url == "grpc://a:1" {
						callsA++
						return "", status.Error(tc.code, "not in this node's state")
					}
					callsB++
					return "held by b", nil
				})

			require.NoError(t, err)
			require.Equal(t, "held by b", result)
			require.Equal(t, 1, callsA)
			require.Equal(t, 1, callsB)
			require.Equal(t, 2, mgr.GetHealthStatus()["active_clients"])
			require.Equal(t, 0, mgr.GetHealthStatus()["cooling_clients"])
			for _, client := range clients {
				require.Zero(t, mgr.IsClientInBackoff(client))
			}
		})
	}
}

// TestQueryFailsOnEachEndpoint covers two endpoints that report missing state.
// The last answer keeps its code, operation name and endpoint. Both stay active
// because neither error is a health penalty.
func TestQueryFailsOnEachEndpoint(t *testing.T) {
	logger := zerolog.Nop()
	clients := []*fakeParticipant{{url: "grpc://a:1"}, {url: "grpc://b:1"}}
	mgr := NewClientPoolManager[*fakeParticipant](clients, logger)
	defer mgr.Close()
	mgr.currentIndex = 0

	var calls int
	_, err := ExecuteWithRetry(context.Background(), mgr, &logger,
		func(_ context.Context, c *fakeParticipant) (struct{}, error) {
			calls++
			return struct{}{}, status.Errorf(codes.NotFound, "state absent on %s", c.url)
		})

	require.Error(t, err)
	require.Equal(t, 2, calls)
	require.Equal(t, codes.NotFound, status.Code(err), "got %v", err)
	require.Contains(t, err.Error(), "Query")
	require.Contains(t, err.Error(), "failed on grpc://b:1")
	require.Contains(t, err.Error(), "grpc://b:1")
	require.Contains(t, err.Error(), "state absent on grpc://b:1")
	require.NotContains(t, err.Error(), "all clients failed")
	require.Equal(t, 2, mgr.GetHealthStatus()["active_clients"])
	require.Equal(t, 0, mgr.GetHealthStatus()["cooling_clients"])
}

// TestExecuteWithRetryCallerCancellationWinsErrorRace cancels the caller while
// the operation returns another error. Both the immediate and exhausted return
// paths must preserve the caller's cancellation.
func TestExecuteWithRetryCallerCancellationWinsErrorRace(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  error
	}{
		{"immediate_unimplemented", status.Error(codes.Unimplemented, "unknown method")},
		{"exhausted_fail_over_only", status.Error(codes.NotFound, "state absent")},
		{"exhausted_penalise", status.Error(codes.Unavailable, "connection failed")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			logger := zerolog.Nop()
			mgr := NewClientPoolManager[*fakeParticipant](
				[]*fakeParticipant{{url: "grpc://a:1"}}, logger)
			defer mgr.Close()

			ctx, cancel := context.WithCancel(context.Background())
			_, err := ExecuteWithRetry(ctx, mgr, &logger,
				func(_ context.Context, _ *fakeParticipant) (struct{}, error) {
					cancel()
					return struct{}{}, tc.err
				})

			require.Error(t, err)
			require.True(t, errors.Is(err, context.Canceled),
				"caller cancellation must not be masked by %v", err)
		})
	}

	for _, withClient := range []bool{false, true} {
		t.Run(fmt.Sprintf("already_cancelled_with_client_%t", withClient), func(t *testing.T) {
			logger := zerolog.Nop()
			clients := []*fakeParticipant{}
			if withClient {
				clients = append(clients, &fakeParticipant{url: "grpc://a:1"})
			}
			mgr := NewClientPoolManager[*fakeParticipant](clients, logger)
			defer mgr.Close()

			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			_, err := ExecuteWithRetry(ctx, mgr, &logger,
				func(_ context.Context, _ *fakeParticipant) (struct{}, error) {
					t.Fatal("a cancelled caller must not start an RPC")
					return struct{}{}, nil
				})
			require.ErrorIs(t, err, context.Canceled)
		})
	}
}
