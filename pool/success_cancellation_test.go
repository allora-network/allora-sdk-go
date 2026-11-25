package pool

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// TestExecuteWithRetryRecordsSuccessBeforeCancellation verifies that a
// completed RPC updates health and clears expired backoff before cancellation
// reaches the caller.
func TestExecuteWithRetryRecordsSuccessBeforeCancellation(t *testing.T) {
	logger := zerolog.Nop()
	client := &fakeParticipant{url: "grpc://a:1"}
	mgr := NewClientPoolManager[*fakeParticipant]([]*fakeParticipant{client}, logger, 0)
	defer mgr.Close()
	mgr.backoff[client.url] = &backoffState{until: time.Now().Add(-time.Hour)}

	ctx, cancel := context.WithCancel(context.Background())
	result, err := ExecuteWithRetry(ctx, mgr, &logger,
		func(_ context.Context, _ *fakeParticipant) (string, error) {
			cancel()
			return "completed", nil
		})

	require.Empty(t, result)
	require.True(t, errors.Is(err, context.Canceled))
	health := mgr.GetHealthStatus()
	active := health["active"].([]map[string]any)
	require.Len(t, active, 1)
	require.Zero(t, mgr.IsClientInBackoff(client), "the completed RPC must clear expired backoff")
	mgr.backMu.Lock()
	_, backoffRecorded := mgr.backoff[client.url]
	mgr.backMu.Unlock()
	require.False(t, backoffRecorded, "expired backoff must be removed on success")
}
