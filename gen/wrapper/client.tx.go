package wrapper

import (
	"context"

	tx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/pool"
)

// TxClientWrapper wraps the tx module with pool management and retry logic
type TxClientWrapper struct {
	poolManager *pool.ClientPoolManager[interfaces.CosmosClient]
	logger      zerolog.Logger
}

// NewTxClientWrapper creates a new tx client wrapper
func NewTxClientWrapper(poolManager *pool.ClientPoolManager[interfaces.CosmosClient], logger zerolog.Logger) *TxClientWrapper {
	return &TxClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "tx").Logger(),
	}
}

func (c *TxClientWrapper) Simulate(ctx context.Context, req *tx.SimulateRequest, opts ...config.CallOpt) (*tx.SimulateResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*tx.SimulateResponse, error) {
		return client.Tx().Simulate(ctx, req, opts...)
	})
}

func (c *TxClientWrapper) GetTx(ctx context.Context, req *tx.GetTxRequest, opts ...config.CallOpt) (*tx.GetTxResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*tx.GetTxResponse, error) {
		return client.Tx().GetTx(ctx, req, opts...)
	})
}

func (c *TxClientWrapper) BroadcastTx(ctx context.Context, req *tx.BroadcastTxRequest, opts ...config.CallOpt) (*tx.BroadcastTxResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*tx.BroadcastTxResponse, error) {
		return client.Tx().BroadcastTx(ctx, req, opts...)
	})
}

func (c *TxClientWrapper) GetTxsEvent(ctx context.Context, req *tx.GetTxsEventRequest, opts ...config.CallOpt) (*tx.GetTxsEventResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*tx.GetTxsEventResponse, error) {
		return client.Tx().GetTxsEvent(ctx, req, opts...)
	})
}

func (c *TxClientWrapper) GetBlockWithTxs(ctx context.Context, req *tx.GetBlockWithTxsRequest, opts ...config.CallOpt) (*tx.GetBlockWithTxsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*tx.GetBlockWithTxsResponse, error) {
		return client.Tx().GetBlockWithTxs(ctx, req, opts...)
	})
}

func (c *TxClientWrapper) TxDecode(ctx context.Context, req *tx.TxDecodeRequest, opts ...config.CallOpt) (*tx.TxDecodeResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*tx.TxDecodeResponse, error) {
		return client.Tx().TxDecode(ctx, req, opts...)
	})
}

func (c *TxClientWrapper) TxEncode(ctx context.Context, req *tx.TxEncodeRequest, opts ...config.CallOpt) (*tx.TxEncodeResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*tx.TxEncodeResponse, error) {
		return client.Tx().TxEncode(ctx, req, opts...)
	})
}

func (c *TxClientWrapper) TxEncodeAmino(ctx context.Context, req *tx.TxEncodeAminoRequest, opts ...config.CallOpt) (*tx.TxEncodeAminoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*tx.TxEncodeAminoResponse, error) {
		return client.Tx().TxEncodeAmino(ctx, req, opts...)
	})
}

func (c *TxClientWrapper) TxDecodeAmino(ctx context.Context, req *tx.TxDecodeAminoRequest, opts ...config.CallOpt) (*tx.TxDecodeAminoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*tx.TxDecodeAminoResponse, error) {
		return client.Tx().TxDecodeAmino(ctx, req, opts...)
	})
}
