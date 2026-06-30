package txsend_test

import (
	"context"
	"testing"

	"github.com/allora-network/allora-sdk-go/cosmosrpc"
	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/txsend"
	"github.com/allora-network/allora-sdk-go/txsend/cosmospool"
	"github.com/brynbellomy/go-utils/errors"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
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

// TestStubMethodsReturnNotImplemented asserts the skeleton's four methods
// compile and return the "not implemented" error for their respective beads,
// guarding against an accidental real implementation landing before its bead.
func TestStubMethodsReturnNotImplemented(t *testing.T) {
	b := cosmospool.New(stubPool{}, zerolog.Nop())

	_, _, err := b.AccountInfo(context.Background(), "allo1...")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not implemented: bead asg-pvd.3")

	// EstimateGas is now implemented (asg-pvd.4): with a nil Tx() pool it
	// recovers via the fallback path, returning FallbackGasEstimate and no error.
	gas, err := b.EstimateGas(context.Background(), []byte{})
	require.NoError(t, err)
	require.NotZero(t, gas)

	_, err = b.Broadcast(context.Background(), []byte{}, txsend.BroadcastModeSync)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not implemented: bead asg-pvd.5")

	_, err = b.WaitForTx(context.Background(), "DEADBEEF")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not implemented: bead asg-pvd.5")
}

// stubPool satisfies txsend.CosmosTxPool minimally so the constructor's non-nil
// path and the stub methods can be exercised without a real cosmos client. Auth
// is nil (only AccountInfo touches it, still a stub). Tx returns a
// stubTxClient whose Simulate always errors, so EstimateGas's fallback path
// returns the fallback gas without panicking on a nil interface.
type stubPool struct{}

func (stubPool) Tx() interfaces.TxClient   { return stubTxClient{} }
func (stubPool) Auth() interfaces.AuthClient { return nil }

// stubTxClient satisfies interfaces.TxClient with every method returning the
// zero value / an error. Only Simulate is meaningful (EstimateGas calls it);
// the rest exist so the interface is satisfied.
type stubTxClient struct{}

func (stubTxClient) Simulate(_ context.Context, _ *txtypes.SimulateRequest, _ ...config.CallOpt) (*txtypes.SimulateResponse, error) {
	return nil, errors.New("stub: simulate unavailable")
}
func (stubTxClient) BroadcastTx(_ context.Context, _ *txtypes.BroadcastTxRequest, _ ...config.CallOpt) (*txtypes.BroadcastTxResponse, error) {
	return nil, errors.New("not implemented")
}
func (stubTxClient) GetBlockWithTxs(_ context.Context, _ *txtypes.GetBlockWithTxsRequest, _ ...config.CallOpt) (*txtypes.GetBlockWithTxsResponse, error) {
	return nil, errors.New("not implemented")
}
func (stubTxClient) GetTx(_ context.Context, _ *txtypes.GetTxRequest, _ ...config.CallOpt) (*txtypes.GetTxResponse, error) {
	return nil, errors.New("not implemented")
}
func (stubTxClient) GetTxsEvent(_ context.Context, _ *txtypes.GetTxsEventRequest, _ ...config.CallOpt) (*txtypes.GetTxsEventResponse, error) {
	return nil, errors.New("not implemented")
}
func (stubTxClient) TxDecode(_ context.Context, _ *txtypes.TxDecodeRequest, _ ...config.CallOpt) (*txtypes.TxDecodeResponse, error) {
	return nil, errors.New("not implemented")
}
func (stubTxClient) TxDecodeAmino(_ context.Context, _ *txtypes.TxDecodeAminoRequest, _ ...config.CallOpt) (*txtypes.TxDecodeAminoResponse, error) {
	return nil, errors.New("not implemented")
}
func (stubTxClient) TxEncode(_ context.Context, _ *txtypes.TxEncodeRequest, _ ...config.CallOpt) (*txtypes.TxEncodeResponse, error) {
	return nil, errors.New("not implemented")
}
func (stubTxClient) TxEncodeAmino(_ context.Context, _ *txtypes.TxEncodeAminoRequest, _ ...config.CallOpt) (*txtypes.TxEncodeAminoResponse, error) {
	return nil, errors.New("not implemented")
}
