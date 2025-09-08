package wrapper

import (
	"context"

	minttypes "github.com/allora-network/allora-chain/x/mint/types"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/pool"
)

// MintClientWrapper wraps the mint module with pool management and retry logic
type MintClientWrapper struct {
	poolManager *pool.ClientPoolManager
	logger      zerolog.Logger
}

// NewMintClientWrapper creates a new mint client wrapper
func NewMintClientWrapper(poolManager *pool.ClientPoolManager, logger zerolog.Logger) *MintClientWrapper {
	return &MintClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "mint").Logger(),
	}
}

func (c *MintClientWrapper) Params(ctx context.Context, req *minttypes.QueryServiceParamsRequest) (*minttypes.QueryServiceParamsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*minttypes.QueryServiceParamsResponse, error) {
		return client.Mint().Params(ctx, req)
	})
}

func (c *MintClientWrapper) Inflation(ctx context.Context, req *minttypes.QueryServiceInflationRequest) (*minttypes.QueryServiceInflationResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*minttypes.QueryServiceInflationResponse, error) {
		return client.Mint().Inflation(ctx, req)
	})
}

func (c *MintClientWrapper) EmissionInfo(ctx context.Context, req *minttypes.QueryServiceEmissionInfoRequest) (*minttypes.QueryServiceEmissionInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*minttypes.QueryServiceEmissionInfoResponse, error) {
		return client.Mint().EmissionInfo(ctx, req)
	})
}
