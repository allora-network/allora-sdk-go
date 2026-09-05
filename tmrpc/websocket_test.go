package tmrpc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestReconnectDelayBackoffAndJitter pins that reconnect delays grow
// exponentially, are capped at wsReconnectMaxDelay, and include jitter (the
// previous implementation used a fixed 5s loop with no backoff or jitter).
func TestReconnectDelayBackoffAndJitter(t *testing.T) {
	// Exponential growth in the low range (base 500ms, doubling).
	prev := time.Duration(0)
	for attempt := 0; attempt < 5; attempt++ {
		// Average several samples to smooth jitter for the monotonic check.
		var sum time.Duration
		const n = 50
		for i := 0; i < n; i++ {
			sum += reconnectDelay(attempt)
		}
		avg := sum / n
		require.Greater(t, avg, prev, "average delay must grow with attempt %d", attempt)
		prev = avg
	}

	// Jitter present: not every sample equals the deterministic base.
	d0 := reconnectDelay(0)
	seenDifferent := false
	for i := 0; i < 20; i++ {
		if reconnectDelay(0) != d0 {
			seenDifferent = true
			break
		}
	}
	require.True(t, seenDifferent, "reconnect delay must include jitter")

	// Capped at max delay (jitter can exceed max by up to wsJitterFrac).
	maxAllowed := time.Duration(float64(wsReconnectMaxDelay) * (1 + wsJitterFrac))
	for attempt := 10; attempt < 40; attempt++ {
		require.LessOrEqual(t, reconnectDelay(attempt), maxAllowed,
			"delay for high attempt %d must be capped near wsReconnectMaxDelay", attempt)
	}
}
