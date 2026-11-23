package pool

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// TestNewClientPoolManagerStartIndex verifies the documented construction-time
// option chooses the first endpoint without changing its default production
// selection path for callers that omit the argument.
func TestNewClientPoolManagerStartIndex(t *testing.T) {
	clients := []*fakeParticipant{{url: "grpc://a:1"}, {url: "grpc://b:1"}}
	for start, endpoint := range []string{"grpc://a:1", "grpc://b:1"} {
		mgr := NewClientPoolManager[*fakeParticipant](clients, zerolog.Nop(), start)
		client, available := mgr.GetClientWithBackoff()
		require.True(t, available)
		require.Equal(t, endpoint, client.url)
		mgr.Close()
	}
	mgr := NewClientPoolManager[*fakeParticipant](clients, zerolog.Nop(), 1)
	mgr.Close()

	require.Panics(t, func() {
		NewClientPoolManager[*fakeParticipant](clients, zerolog.Nop(), -1)
	})
	require.Panics(t, func() {
		NewClientPoolManager[*fakeParticipant](clients, zerolog.Nop(), len(clients))
	})
	require.Panics(t, func() {
		NewClientPoolManager[*fakeParticipant](clients, zerolog.Nop(), 0, 1)
	})
}
