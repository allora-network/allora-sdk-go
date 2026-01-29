package wrapper

import (
	"context"

	minttypes "github.com/allora-network/allora-chain/x/mint/types"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/pool"
)

// MintClientWrapper wraps the mint module with pool management and retry logic
type MintClientWrapper struct {
	poolManager *pool.ClientPoolManager[interfaces.CosmosClient]
	logger      zerolog.Logger
}

// NewMintClientWrapper creates a new mint client wrapper
func NewMintClientWrapper(poolManager *pool.ClientPoolManager[interfaces.CosmosClient], logger zerolog.Logger) *MintClientWrapper {
	return &MintClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "mint").Logger(),
	}
}

func (c *MintClientWrapper) EmissionInfo(ctx context.Context, req *minttypes.QueryServiceEmissionInfoRequest, opts ...config.CallOpt) (*minttypes.QueryServiceEmissionInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*minttypes.QueryServiceEmissionInfoResponse, error) {
		return client.Mint().EmissionInfo(ctx, req, opts...)
	})
}

func (c *MintClientWrapper) Inflation(ctx context.Context, req *minttypes.QueryServiceInflationRequest, opts ...config.CallOpt) (*minttypes.QueryServiceInflationResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*minttypes.QueryServiceInflationResponse, error) {
		return client.Mint().Inflation(ctx, req, opts...)
	})
}

func (c *MintClientWrapper) Params(ctx context.Context, req *minttypes.QueryServiceParamsRequest, opts ...config.CallOpt) (*minttypes.QueryServiceParamsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*minttypes.QueryServiceParamsResponse, error) {
		return client.Mint().Params(ctx, req, opts...)
	})
}
