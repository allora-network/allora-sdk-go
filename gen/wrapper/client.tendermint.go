package wrapper

import (
	"context"

	cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/pool"
)

// TendermintClientWrapper wraps the tendermint module with pool management and retry logic
type TendermintClientWrapper struct {
	poolManager *pool.ClientPoolManager
	logger      zerolog.Logger
}

// NewTendermintClientWrapper creates a new tendermint client wrapper
func NewTendermintClientWrapper(poolManager *pool.ClientPoolManager, logger zerolog.Logger) *TendermintClientWrapper {
	return &TendermintClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "tendermint").Logger(),
	}
}

func (c *TendermintClientWrapper) GetNodeInfo(ctx context.Context, req *cmtservice.GetNodeInfoRequest, opts ...config.CallOpt) (*cmtservice.GetNodeInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*cmtservice.GetNodeInfoResponse, error) {
		return client.Tendermint().GetNodeInfo(ctx, req, opts...)
	})
}

func (c *TendermintClientWrapper) GetSyncing(ctx context.Context, req *cmtservice.GetSyncingRequest, opts ...config.CallOpt) (*cmtservice.GetSyncingResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*cmtservice.GetSyncingResponse, error) {
		return client.Tendermint().GetSyncing(ctx, req, opts...)
	})
}

func (c *TendermintClientWrapper) GetLatestBlock(ctx context.Context, req *cmtservice.GetLatestBlockRequest, opts ...config.CallOpt) (*cmtservice.GetLatestBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*cmtservice.GetLatestBlockResponse, error) {
		return client.Tendermint().GetLatestBlock(ctx, req, opts...)
	})
}

func (c *TendermintClientWrapper) GetBlockByHeight(ctx context.Context, req *cmtservice.GetBlockByHeightRequest, opts ...config.CallOpt) (*cmtservice.GetBlockByHeightResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*cmtservice.GetBlockByHeightResponse, error) {
		return client.Tendermint().GetBlockByHeight(ctx, req, opts...)
	})
}

func (c *TendermintClientWrapper) GetLatestValidatorSet(ctx context.Context, req *cmtservice.GetLatestValidatorSetRequest, opts ...config.CallOpt) (*cmtservice.GetLatestValidatorSetResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*cmtservice.GetLatestValidatorSetResponse, error) {
		return client.Tendermint().GetLatestValidatorSet(ctx, req, opts...)
	})
}

func (c *TendermintClientWrapper) GetValidatorSetByHeight(ctx context.Context, req *cmtservice.GetValidatorSetByHeightRequest, opts ...config.CallOpt) (*cmtservice.GetValidatorSetByHeightResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*cmtservice.GetValidatorSetByHeightResponse, error) {
		return client.Tendermint().GetValidatorSetByHeight(ctx, req, opts...)
	})
}

func (c *TendermintClientWrapper) ABCIQuery(ctx context.Context, req *cmtservice.ABCIQueryRequest, opts ...config.CallOpt) (*cmtservice.ABCIQueryResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*cmtservice.ABCIQueryResponse, error) {
		return client.Tendermint().ABCIQuery(ctx, req, opts...)
	})
}
