package cosmospool_test

import (
	"context"
	"testing"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/txsend"
	"github.com/allora-network/allora-sdk-go/txsend/cosmospool"
	"github.com/brynbellomy/go-utils/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// fakeTxClient is a minimal interfaces.TxClient whose only meaningful method
// is Simulate; the rest return errors. It lets the gas tests drive the
// EstimateGas path without a generated mock.
type fakeTxClient struct {
	resp    *txtypes.SimulateResponse
	err     error
	called  int
	lastReq *txtypes.SimulateRequest
}

func (f *fakeTxClient) Simulate(_ context.Context, req *txtypes.SimulateRequest, _ ...config.CallOpt) (*txtypes.SimulateResponse, error) {
	f.called++
	f.lastReq = req
	return f.resp, f.err
}
func (fakeTxClient) BroadcastTx(context.Context, *txtypes.BroadcastTxRequest, ...config.CallOpt) (*txtypes.BroadcastTxResponse, error) {
	return nil, errors.New("not implemented")
}
func (fakeTxClient) GetBlockWithTxs(context.Context, *txtypes.GetBlockWithTxsRequest, ...config.CallOpt) (*txtypes.GetBlockWithTxsResponse, error) {
	return nil, errors.New("not implemented")
}
func (fakeTxClient) GetTx(context.Context, *txtypes.GetTxRequest, ...config.CallOpt) (*txtypes.GetTxResponse, error) {
	return nil, errors.New("not implemented")
}
func (fakeTxClient) GetTxsEvent(context.Context, *txtypes.GetTxsEventRequest, ...config.CallOpt) (*txtypes.GetTxsEventResponse, error) {
	return nil, errors.New("not implemented")
}
func (fakeTxClient) TxDecode(context.Context, *txtypes.TxDecodeRequest, ...config.CallOpt) (*txtypes.TxDecodeResponse, error) {
	return nil, errors.New("not implemented")
}
func (fakeTxClient) TxDecodeAmino(context.Context, *txtypes.TxDecodeAminoRequest, ...config.CallOpt) (*txtypes.TxDecodeAminoResponse, error) {
	return nil, errors.New("not implemented")
}
func (fakeTxClient) TxEncode(context.Context, *txtypes.TxEncodeRequest, ...config.CallOpt) (*txtypes.TxEncodeResponse, error) {
	return nil, errors.New("not implemented")
}
func (fakeTxClient) TxEncodeAmino(context.Context, *txtypes.TxEncodeAminoRequest, ...config.CallOpt) (*txtypes.TxEncodeAminoResponse, error) {
	return nil, errors.New("not implemented")
}

// Compile-time assertion that fakeTxClient satisfies the interface.
var _ interfaces.TxClient = (*fakeTxClient)(nil)

// fakePool is a minimal txsend.CosmosTxPool returning the injected TxClient.
type fakePool struct{ tx interfaces.TxClient }

func (p fakePool) Tx() interfaces.TxClient   { return p.tx }
func (fakePool) Auth() interfaces.AuthClient { return nil }

var _ txsend.CosmosTxPool = (*fakePool)(nil)

func TestEstimateGas(t *testing.T) {
	unsignedTx := []byte("deadbeef")

	tests := []struct {
		name          string
		gasUsed       uint64 // 0 means simulate returns an error
		simulateErr   error
		adjustment    float64 // 0 means use default
		wantGas       uint64
		wantErr       bool
		wantCallCount int
		ctxCancelled  bool
		ctxDeadline   bool
	}{
		{
			name:          "success: GasUsed=100000 returns raw 100000 (adjustment applied by caller)",
			gasUsed:       100000,
			adjustment:    1.5,
			wantGas:       100000,
			wantCallCount: 1,
		},
		{
			name:          "success with default adjustment",
			gasUsed:       100000,
			adjustment:    0, // default 1.5, but EstimateGas returns raw regardless
			wantGas:       100000,
			wantCallCount: 1,
		},
		{
			name:          "simulate error → fallback returned (no error)",
			gasUsed:       0,
			simulateErr:   errors.New("rpc error: simulate unavailable"),
			adjustment:    1.5,
			wantGas:       cosmospool.FallbackGasEstimate,
			wantErr:       false,
			wantCallCount: 1,
		},
		{
			name:          "ceil edge: GasUsed=100001*1.5 → 150002",
			gasUsed:       100001,
			adjustment:    1.5,
			wantGas:       100001,
			wantCallCount: 1,
		},
		{
			name:          "nil GasInfo → fallback returned (no error)",
			gasUsed:       0, // unused; resp constructed with nil GasInfo below
			simulateErr:   nil,
			adjustment:    1.5,
			wantGas:       cosmospool.FallbackGasEstimate,
			wantCallCount: 1,
		},
		{
			name:          "ctx cancelled returns ctx error (NOT fallback)",
			gasUsed:       0,
			simulateErr:   context.Canceled,
			adjustment:    1.5,
			wantGas:       0,
			wantErr:       true,
			wantCallCount: 1,
			ctxCancelled:  true,
		},
		{
			name:          "ctx deadline exceeded returns ctx error (NOT fallback)",
			gasUsed:       0,
			simulateErr:   context.DeadlineExceeded,
			adjustment:    1.5,
			wantGas:       0,
			wantErr:       true,
			wantCallCount: 1,
			ctxDeadline:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp *txtypes.SimulateResponse
			if tt.simulateErr != nil {
				resp = nil
			} else if tt.name == "nil GasInfo → fallback returned (no error)" {
				resp = &txtypes.SimulateResponse{GasInfo: nil}
			} else {
				resp = &txtypes.SimulateResponse{
					GasInfo: &types.GasInfo{GasUsed: tt.gasUsed},
				}
			}

			ftc := &fakeTxClient{resp: resp, err: tt.simulateErr}
			pool := fakePool{tx: ftc}

			opts := []cosmospool.Option{}
			if tt.adjustment > 0 {
				opts = append(opts, cosmospool.WithGasAdjustment(tt.adjustment))
			}
			b := cosmospool.New(pool, zerolog.Nop(), opts...)

			ctx := context.Background()
			if tt.ctxCancelled {
				c, cancel := context.WithCancel(ctx)
				cancel()
				ctx = c
			} else if tt.ctxDeadline {
				c, cancel := context.WithTimeout(ctx, 0)
				defer cancel()
				ctx = c
			}
			gas, err := b.EstimateGas(ctx, unsignedTx)
			if tt.wantErr {
				require.Error(t, err)
				require.Zero(t, gas)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantGas, gas)
			}
			require.Equal(t, tt.wantCallCount, ftc.called, "Simulate call count")
			if tt.wantCallCount > 0 && ftc.lastReq != nil {
				require.Equal(t, unsignedTx, ftc.lastReq.TxBytes, "Simulate should receive the raw tx bytes")
			}
		})
	}
}

func TestWithGasAdjustment(t *testing.T) {
	// WithGasAdjustment configures the broadcaster's adjustment field (used for
	// the fallback path and as the documented default). EstimateGas itself now
	// returns RAW simulated gas (the caller applies the multiplier), so these
	// cases assert the configured field value via a helper rather than the
	// EstimateGas return.
	t.Run("default is 1.5", func(t *testing.T) {
		b := cosmospool.New(fakePool{tx: &fakeTxClient{}}, zerolog.Nop())
		require.Equal(t, cosmospool.DefaultGasAdjustment, b.GasAdjustment())
	})

	t.Run("zero adjustment is ignored (keeps default)", func(t *testing.T) {
		b := cosmospool.New(fakePool{tx: &fakeTxClient{}}, zerolog.Nop(), cosmospool.WithGasAdjustment(0))
		require.Equal(t, cosmospool.DefaultGasAdjustment, b.GasAdjustment())
	})

	t.Run("custom adjustment", func(t *testing.T) {
		b := cosmospool.New(fakePool{tx: &fakeTxClient{}}, zerolog.Nop(), cosmospool.WithGasAdjustment(2.0))
		require.Equal(t, 2.0, b.GasAdjustment())
	})

	t.Run("values below 1.0 are clamped to 1.0", func(t *testing.T) {
		b := cosmospool.New(fakePool{tx: &fakeTxClient{}}, zerolog.Nop(), cosmospool.WithGasAdjustment(0.5))
		require.Equal(t, 1.0, b.GasAdjustment())
	})
}

func TestGetGasPrice(t *testing.T) {
	tests := []struct {
		name string
		tier cosmospool.FeeTier
		want sdkmath.LegacyDec
	}{
		{"low", cosmospool.FeeTierLow, sdkmath.LegacyMustNewDecFromStr("0.001")},
		{"medium", cosmospool.FeeTierMedium, sdkmath.LegacyMustNewDecFromStr("0.01")},
		{"high", cosmospool.FeeTierHigh, sdkmath.LegacyMustNewDecFromStr("0.025")},
		{"unknown falls back to medium", cosmospool.FeeTier("bogus"), sdkmath.LegacyMustNewDecFromStr("0.01")},
		{"empty falls back to medium", cosmospool.FeeTier(""), sdkmath.LegacyMustNewDecFromStr("0.01")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cosmospool.GetGasPrice(tt.tier)
			require.True(t, got.Equal(tt.want), "GetGasPrice(%q) = %s, want %s", tt.tier, got.String(), tt.want.String())
		})
	}
}
