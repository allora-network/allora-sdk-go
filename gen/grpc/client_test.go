package grpc

import (
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/allora-network/allora-sdk-go/config"
)

// TestNewGRPCClientDialErrorWrapsCause pins that a dial failure wraps the
// underlying cause so errors.Is classification works for callers. Regression:
// errors.Errorf (Sprintf-based) cannot wrap %w, producing "%!w(...)" and
// breaking errors.Is. Dial an unroutable address; the 15s internal dial
// deadline (or an immediate connection refusal) yields an error that must NOT
// contain the %!w artifact and must wrap its cause.
func TestNewGRPCClientDialErrorWrapsCause(t *testing.T) {
	// 192.0.2.0/24 is TEST-NET-1 (RFC 5737), guaranteed unroutable.
	cfg := config.EndpointConfig{
		URL:      "grpc://192.0.2.1:9090",
		Protocol: config.ProtocolGRPC,
	}

	start := time.Now()
	_, err := NewGRPCClient(cfg, zerolog.Nop())
	require.Error(t, err)
	require.Less(t, time.Since(start), 20*time.Second, "dial should respect its internal timeout")

	// The error must be a real wrapped error, not the %!w artifact.
	require.NotContains(t, err.Error(), "%!w",
		"error must wrap its cause with %w, not Sprintf-format it")
	// Unwrapping must reach a non-nil cause (errors.Is traversable).
	require.NotNil(t, errors.Unwrap(err), "dial error must wrap its underlying cause")
}
