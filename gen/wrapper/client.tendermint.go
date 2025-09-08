package wrapper

import (
	"context"

	cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
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

func (c *TendermintClientWrapper) GetNodeInfo(ctx context.Context, req *cmtservice.GetNodeInfoRequest) (*cmtservice.GetNodeInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*cmtservice.GetNodeInfoResponse, error) {
		return client.Tendermint().GetNodeInfo(ctx, req)
	})
}

func (c *TendermintClientWrapper) GetSyncing(ctx context.Context, req *cmtservice.GetSyncingRequest) (*cmtservice.GetSyncingResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*cmtservice.GetSyncingResponse, error) {
		return client.Tendermint().GetSyncing(ctx, req)
	})
}

func (c *TendermintClientWrapper) GetLatestBlock(ctx context.Context, req *cmtservice.GetLatestBlockRequest) (*cmtservice.GetLatestBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*cmtservice.GetLatestBlockResponse, error) {
		return client.Tendermint().GetLatestBlock(ctx, req)
	})
}

func (c *TendermintClientWrapper) GetBlockByHeight(ctx context.Context, req *cmtservice.GetBlockByHeightRequest) (*cmtservice.GetBlockByHeightResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*cmtservice.GetBlockByHeightResponse, error) {
		return client.Tendermint().GetBlockByHeight(ctx, req)
	})
}

func (c *TendermintClientWrapper) GetLatestValidatorSet(ctx context.Context, req *cmtservice.GetLatestValidatorSetRequest) (*cmtservice.GetLatestValidatorSetResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*cmtservice.GetLatestValidatorSetResponse, error) {
		return client.Tendermint().GetLatestValidatorSet(ctx, req)
	})
}

func (c *TendermintClientWrapper) GetValidatorSetByHeight(ctx context.Context, req *cmtservice.GetValidatorSetByHeightRequest) (*cmtservice.GetValidatorSetByHeightResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*cmtservice.GetValidatorSetByHeightResponse, error) {
		return client.Tendermint().GetValidatorSetByHeight(ctx, req)
	})
}

func (c *TendermintClientWrapper) ABCIQuery(ctx context.Context, req *cmtservice.ABCIQueryRequest) (*cmtservice.ABCIQueryResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*cmtservice.ABCIQueryResponse, error) {
		return client.Tendermint().ABCIQuery(ctx, req)
	})
}
