package cosmospool_test

import (
	"context"
	"errors"
	"testing"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	sdktx "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/txsend"
	"github.com/allora-network/allora-sdk-go/txsend/cosmospool"
	"github.com/rs/zerolog"
)

// fakeTxClient implements interfaces.TxClient. BroadcastTx and GetTx are
// meaningful; all other methods panic (not called in these tests).
type fakeTxClient struct {
	broadcastFn func(ctx context.Context, req *txtypes.BroadcastTxRequest, opts ...config.CallOpt) (*txtypes.BroadcastTxResponse, error)
	getTxFn     func(ctx context.Context, req *txtypes.GetTxRequest, opts ...config.CallOpt) (*txtypes.GetTxResponse, error)
}

func (f *fakeTxClient) BroadcastTx(ctx context.Context, req *txtypes.BroadcastTxRequest, opts ...config.CallOpt) (*txtypes.BroadcastTxResponse, error) {
	return f.broadcastFn(ctx, req, opts...)
}

func (f *fakeTxClient) GetTx(ctx context.Context, req *txtypes.GetTxRequest, opts ...config.CallOpt) (*txtypes.GetTxResponse, error) {
	return f.getTxFn(ctx, req, opts...)
}

func (f *fakeTxClient) GetBlockWithTxs(ctx context.Context, req *txtypes.GetBlockWithTxsRequest, opts ...config.CallOpt) (*txtypes.GetBlockWithTxsResponse, error) {
	panic("GetBlockWithTxs not implemented in fakeTxClient")
}

func (f *fakeTxClient) GetTxsEvent(ctx context.Context, req *txtypes.GetTxsEventRequest, opts ...config.CallOpt) (*txtypes.GetTxsEventResponse, error) {
	panic("GetTxsEvent not implemented in fakeTxClient")
}

func (f *fakeTxClient) Simulate(ctx context.Context, req *txtypes.SimulateRequest, opts ...config.CallOpt) (*txtypes.SimulateResponse, error) {
	panic("Simulate not implemented in fakeTxClient")
}

func (f *fakeTxClient) TxDecode(ctx context.Context, req *txtypes.TxDecodeRequest, opts ...config.CallOpt) (*txtypes.TxDecodeResponse, error) {
	panic("TxDecode not implemented in fakeTxClient")
}

func (f *fakeTxClient) TxDecodeAmino(ctx context.Context, req *txtypes.TxDecodeAminoRequest, opts ...config.CallOpt) (*txtypes.TxDecodeAminoResponse, error) {
	panic("TxDecodeAmino not implemented in fakeTxClient")
}

func (f *fakeTxClient) TxEncode(ctx context.Context, req *txtypes.TxEncodeRequest, opts ...config.CallOpt) (*txtypes.TxEncodeResponse, error) {
	panic("TxEncode not implemented in fakeTxClient")
}

func (f *fakeTxClient) TxEncodeAmino(ctx context.Context, req *txtypes.TxEncodeAminoRequest, opts ...config.CallOpt) (*txtypes.TxEncodeAminoResponse, error) {
	panic("TxEncodeAmino not implemented in fakeTxClient")
}

var _ interfaces.TxClient = (*fakeTxClient)(nil)

// stubPool satisfies txsend.CosmosTxPool so we can construct a Broadcaster
// without a real cosmosrpc.ClientPool.
type stubPool struct {
	txClient interfaces.TxClient
}

func (s stubPool) Tx() interfaces.TxClient     { return s.txClient }
func (s stubPool) Auth() interfaces.AuthClient { return nil }

// fakeClock is a Clock whose After returns a channel the test controls.
// Callers send on ch to simulate a tick.
type fakeClock struct {
	ch <-chan time.Time
}

func (f fakeClock) Now() time.Time                        { return time.Time{} }
func (f fakeClock) After(d time.Duration) <-chan time.Time { return f.ch }

// ---------------------------------------------------------------------------
// Broadcast tests
// ---------------------------------------------------------------------------

func TestBroadcast_SyncSuccess(t *testing.T) {
	fc := &fakeTxClient{
		broadcastFn: func(ctx context.Context, req *txtypes.BroadcastTxRequest, opts ...config.CallOpt) (*txtypes.BroadcastTxResponse, error) {
			require.Equal(t, txtypes.BroadcastMode_BROADCAST_MODE_SYNC, req.Mode)
			require.Equal(t, []byte("signed"), req.TxBytes)
			return &txtypes.BroadcastTxResponse{
				TxResponse: &sdktx.TxResponse{
					TxHash:    "ABC123",
					Code:      0,
					Codespace: "",
					RawLog:    "",
				},
			}, nil
		},
	}

	b := cosmospool.New(stubPool{txClient: fc}, zerolog.Nop())
	result, err := b.Broadcast(context.Background(), []byte("signed"), txsend.BroadcastModeSync)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "ABC123", result.TxHash)
	require.Equal(t, uint32(0), result.Code)
}

func TestBroadcast_NonZeroCode(t *testing.T) {
	// Non-zero Code is NOT an error — result returned with Code set, err nil.
	fc := &fakeTxClient{
		broadcastFn: func(ctx context.Context, req *txtypes.BroadcastTxRequest, opts ...config.CallOpt) (*txtypes.BroadcastTxResponse, error) {
			return &txtypes.BroadcastTxResponse{
				TxResponse: &sdktx.TxResponse{
					TxHash:    "ABC123",
					Code:      5,
					Codespace: "sdk",
					RawLog:    "out of gas",
				},
			}, nil
		},
	}

	b := cosmospool.New(stubPool{txClient: fc}, zerolog.Nop())
	result, err := b.Broadcast(context.Background(), []byte("signed"), txsend.BroadcastModeSync)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, uint32(5), result.Code)
	require.Equal(t, "sdk", result.Codespace)
	require.Equal(t, "out of gas", result.RawLog)
}

func TestBroadcast_TransportError(t *testing.T) {
	transportErr := errors.New("gRPC connection refused")
	fc := &fakeTxClient{
		broadcastFn: func(ctx context.Context, req *txtypes.BroadcastTxRequest, opts ...config.CallOpt) (*txtypes.BroadcastTxResponse, error) {
			return nil, transportErr
		},
	}

	b := cosmospool.New(stubPool{txClient: fc}, zerolog.Nop())
	result, err := b.Broadcast(context.Background(), []byte("signed"), txsend.BroadcastModeSync)

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "broadcasting tx")
}

func TestBroadcast_UnknownMode(t *testing.T) {
	fc := &fakeTxClient{
		broadcastFn: func(ctx context.Context, req *txtypes.BroadcastTxRequest, opts ...config.CallOpt) (*txtypes.BroadcastTxResponse, error) {
			t.Fatal("BroadcastTx should not be called for unknown mode")
			return nil, nil
		},
	}

	b := cosmospool.New(stubPool{txClient: fc}, zerolog.Nop())
	result, err := b.Broadcast(context.Background(), []byte("signed"), txsend.BroadcastMode(99))

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "unknown broadcast mode")
}

// ---------------------------------------------------------------------------
// WaitForTx tests
// ---------------------------------------------------------------------------

func TestWaitForTx_NotFoundThenSuccess(t *testing.T) {
	callCount := 0
	fc := &fakeTxClient{
		getTxFn: func(ctx context.Context, req *txtypes.GetTxRequest, opts ...config.CallOpt) (*txtypes.GetTxResponse, error) {
			callCount++
			require.Equal(t, "ABC123", req.Hash)
			if callCount <= 2 {
				// NotFound: return error, no response
				return nil, errors.New("tx not found")
			}
			// Success on third call
			return &txtypes.GetTxResponse{
				TxResponse: &sdktx.TxResponse{
					TxHash:    "ABC123",
					Height:    42,
					Code:      0,
					Codespace: "",
					Data:      "data",
					RawLog:    "ok",
					GasWanted: 100000,
					GasUsed:   50000,
					Timestamp: "2026-01-01T00:00:00Z",
					Events: []abci.Event{
						{
							Type: "message",
							Attributes: []abci.EventAttribute{
								{Key: "action", Value: "send"},
							},
						},
					},
				},
			}, nil
		},
	}

	tickCh := make(chan time.Time)
	fclock := fakeClock{ch: tickCh}
	b := cosmospool.New(stubPool{txClient: fc}, zerolog.Nop(),
		cosmospool.WithClock(fclock),
		cosmospool.WithPollInterval(1*time.Millisecond),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start WaitForTx in a goroutine
	type outcome struct {
		result *txsend.TxResult
		err    error
	}
	chOut := make(chan outcome, 1)
	go func() {
		r, e := b.WaitForTx(ctx, "ABC123")
		chOut <- outcome{r, e}
	}()

	// Drive 2 ticks (NotFound both times)
	tickCh <- time.Now()
	tickCh <- time.Now()

	// Drive 3rd tick (success)
	tickCh <- time.Now()

	out := <-chOut
	require.NoError(t, out.err)
	require.NotNil(t, out.result)
	require.Equal(t, "ABC123", out.result.TxHash)
	require.Equal(t, int64(42), out.result.Height)
	require.Equal(t, uint32(0), out.result.Code)
	require.Equal(t, "data", out.result.Data)
	require.Equal(t, int64(100000), out.result.GasWanted)
	require.Equal(t, int64(50000), out.result.GasUsed)
	require.Equal(t, "2026-01-01T00:00:00Z", out.result.Timestamp)
	require.Len(t, out.result.Events, 1)
	require.Equal(t, "message", out.result.Events[0].Type)
	require.Len(t, out.result.Events[0].Attributes, 1)
	require.Equal(t, "action", out.result.Events[0].Attributes[0].Key)
	require.Equal(t, "send", out.result.Events[0].Attributes[0].Value)
}

func TestWaitForTx_CtxCancelled(t *testing.T) {
	fc := &fakeTxClient{
		getTxFn: func(ctx context.Context, req *txtypes.GetTxRequest, opts ...config.CallOpt) (*txtypes.GetTxResponse, error) {
			return nil, errors.New("tx not found")
		},
	}

	tickCh := make(chan time.Time)
	fclock := fakeClock{ch: tickCh}
	b := cosmospool.New(stubPool{txClient: fc}, zerolog.Nop(),
		cosmospool.WithClock(fclock),
		cosmospool.WithPollInterval(1*time.Millisecond),
	)

	ctx, cancel := context.WithCancel(context.Background())

	// Start WaitForTx
	chOut := make(chan error, 1)
	go func() {
		_, e := b.WaitForTx(ctx, "ABC123")
		chOut <- e
	}()

	// Drive one tick so the select resets
	tickCh <- time.Now()

	// Cancel the context: the next iteration through the loop hits ctx.Done()
	cancel()

	err := <-chOut
	require.Error(t, err)
	require.Contains(t, err.Error(), "waiting for tx ABC123")
	require.Contains(t, err.Error(), "canceled")
}

func TestWaitForTx_CommittedButFailed(t *testing.T) {
	// Committed-but-failed: Code != 0, err nil.
	fc := &fakeTxClient{
		getTxFn: func(ctx context.Context, req *txtypes.GetTxRequest, opts ...config.CallOpt) (*txtypes.GetTxResponse, error) {
			return &txtypes.GetTxResponse{
				TxResponse: &sdktx.TxResponse{
					TxHash:    "ABC123",
					Height:    42,
					Code:      11,
					Codespace: "sdk",
					RawLog:    "failed to execute",
				},
			}, nil
		},
	}

	tickCh := make(chan time.Time)
	fclock := fakeClock{ch: tickCh}
	b := cosmospool.New(stubPool{txClient: fc}, zerolog.Nop(),
		cosmospool.WithClock(fclock),
		cosmospool.WithPollInterval(1*time.Millisecond),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chOut := make(chan struct {
		result *txsend.TxResult
		err    error
	}, 1)
	go func() {
		r, e := b.WaitForTx(ctx, "ABC123")
		chOut <- struct {
			result *txsend.TxResult
			err    error
		}{r, e}
	}()

	tickCh <- time.Now()

	out := <-chOut
	require.NoError(t, out.err)
	require.NotNil(t, out.result)
	require.Equal(t, uint32(11), out.result.Code)
	require.Equal(t, "sdk", out.result.Codespace)
	require.Equal(t, "failed to execute", out.result.RawLog)
}
