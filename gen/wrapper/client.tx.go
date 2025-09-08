package wrapper

import (
	"context"

	tx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/pool"
)

// TxClientWrapper wraps the tx module with pool management and retry logic
type TxClientWrapper struct {
	poolManager *pool.ClientPoolManager
	logger      zerolog.Logger
}

// NewTxClientWrapper creates a new tx client wrapper
func NewTxClientWrapper(poolManager *pool.ClientPoolManager, logger zerolog.Logger) *TxClientWrapper {
	return &TxClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "tx").Logger(),
	}
}

func (c *TxClientWrapper) Simulate(ctx context.Context, req *tx.SimulateRequest) (*tx.SimulateResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*tx.SimulateResponse, error) {
		return client.Tx().Simulate(ctx, req)
	})
}

func (c *TxClientWrapper) GetTx(ctx context.Context, req *tx.GetTxRequest) (*tx.GetTxResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*tx.GetTxResponse, error) {
		return client.Tx().GetTx(ctx, req)
	})
}

func (c *TxClientWrapper) BroadcastTx(ctx context.Context, req *tx.BroadcastTxRequest) (*tx.BroadcastTxResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*tx.BroadcastTxResponse, error) {
		return client.Tx().BroadcastTx(ctx, req)
	})
}

func (c *TxClientWrapper) GetTxsEvent(ctx context.Context, req *tx.GetTxsEventRequest) (*tx.GetTxsEventResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*tx.GetTxsEventResponse, error) {
		return client.Tx().GetTxsEvent(ctx, req)
	})
}

func (c *TxClientWrapper) GetBlockWithTxs(ctx context.Context, req *tx.GetBlockWithTxsRequest) (*tx.GetBlockWithTxsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*tx.GetBlockWithTxsResponse, error) {
		return client.Tx().GetBlockWithTxs(ctx, req)
	})
}

func (c *TxClientWrapper) TxDecode(ctx context.Context, req *tx.TxDecodeRequest) (*tx.TxDecodeResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*tx.TxDecodeResponse, error) {
		return client.Tx().TxDecode(ctx, req)
	})
}

func (c *TxClientWrapper) TxEncode(ctx context.Context, req *tx.TxEncodeRequest) (*tx.TxEncodeResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*tx.TxEncodeResponse, error) {
		return client.Tx().TxEncode(ctx, req)
	})
}

func (c *TxClientWrapper) TxEncodeAmino(ctx context.Context, req *tx.TxEncodeAminoRequest) (*tx.TxEncodeAminoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*tx.TxEncodeAminoResponse, error) {
		return client.Tx().TxEncodeAmino(ctx, req)
	})
}

func (c *TxClientWrapper) TxDecodeAmino(ctx context.Context, req *tx.TxDecodeAminoRequest) (*tx.TxDecodeAminoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*tx.TxDecodeAminoResponse, error) {
		return client.Tx().TxDecodeAmino(ctx, req)
	})
}
