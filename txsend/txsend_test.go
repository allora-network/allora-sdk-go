package txsend_test

import (
	"context"
	"testing"

	"github.com/allora-network/allora-sdk-go/cosmosrpc"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/txsend"
	"github.com/allora-network/allora-sdk-go/txsend/cosmospool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// This compile-time assertion lives in an external _test package (not in
// production txsend or cosmospool) deliberately: txsend.CosmosTxPool is the
// narrow interface txsend depends on, and we must verify that the concrete
// cosmosrpc.ClientPool satisfies it WITHOUT introducing a production import
// cycle. txsend itself is forbidden from importing cosmosrpc (cosmosrpc pulls
// in gen/wrapper -> gen/interfaces and the wiring layer sits in the top-level
// allora package; a txsend -> cosmosrpc edge would close a cycle). The _test
// package can import both without contaminating production, so the interface
// satisfaction is checked here.
var _ txsend.CosmosTxPool = (cosmosrpc.ClientPool)(nil)

// TestNewPanicsOnNilPool asserts the invariant: a Broadcaster without a pool is
// a wiring bug, not a recoverable runtime error.
func TestNewPanicsOnNilPool(t *testing.T) {
	require.Panics(t, func() {
		cosmospool.New(nil, zerolog.Nop())
	}, "New must panic when pool is nil")
}

// TestStubMethodsReturnNotImplemented asserts the skeleton's two remaining
// stub methods compile and return the "not implemented" error for their beads,
// guarding against an accidental real implementation landing before its bead.
func TestStubMethodsReturnNotImplemented(t *testing.T) {
	b := cosmospool.New(stubPool{}, zerolog.Nop())

	_, _, err := b.AccountInfo(context.Background(), "allo1...")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not implemented: bead asg-pvd.3")

	_, err = b.EstimateGas(context.Background(), []byte{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "not implemented: bead asg-pvd.4")
}

// stubPool satisfies txsend.CosmosTxPool minimally so the constructor's non-nil
// path and the stub methods can be exercised without a real cosmos client. The
// returned clients are nil interface values; the stub broadcaster methods never
// call them, so this is safe.
type stubPool struct{}

func (stubPool) Tx() interfaces.TxClient   { return nil }
func (stubPool) Auth() interfaces.AuthClient { return nil }
