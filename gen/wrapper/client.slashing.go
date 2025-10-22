package wrapper

import (
	"context"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/pool"
)

// SlashingClientWrapper wraps the slashing module with pool management and retry logic
type SlashingClientWrapper struct {
	poolManager *pool.ClientPoolManager[interfaces.CosmosClient]
	logger      zerolog.Logger
}

// NewSlashingClientWrapper creates a new slashing client wrapper
func NewSlashingClientWrapper(poolManager *pool.ClientPoolManager[interfaces.CosmosClient], logger zerolog.Logger) *SlashingClientWrapper {
	return &SlashingClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "slashing").Logger(),
	}
}

func (c *SlashingClientWrapper) Params(ctx context.Context, req *slashingtypes.QueryParamsRequest, opts ...config.CallOpt) (*slashingtypes.QueryParamsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*slashingtypes.QueryParamsResponse, error) {
		return client.Slashing().Params(ctx, req, opts...)
	})
}

func (c *SlashingClientWrapper) SigningInfo(ctx context.Context, req *slashingtypes.QuerySigningInfoRequest, opts ...config.CallOpt) (*slashingtypes.QuerySigningInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*slashingtypes.QuerySigningInfoResponse, error) {
		return client.Slashing().SigningInfo(ctx, req, opts...)
	})
}

func (c *SlashingClientWrapper) SigningInfos(ctx context.Context, req *slashingtypes.QuerySigningInfosRequest, opts ...config.CallOpt) (*slashingtypes.QuerySigningInfosResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*slashingtypes.QuerySigningInfosResponse, error) {
		return client.Slashing().SigningInfos(ctx, req, opts...)
	})
}
