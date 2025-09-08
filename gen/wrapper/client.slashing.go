package wrapper

import (
	"context"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/pool"
)

// SlashingClientWrapper wraps the slashing module with pool management and retry logic
type SlashingClientWrapper struct {
	poolManager *pool.ClientPoolManager
	logger      zerolog.Logger
}

// NewSlashingClientWrapper creates a new slashing client wrapper
func NewSlashingClientWrapper(poolManager *pool.ClientPoolManager, logger zerolog.Logger) *SlashingClientWrapper {
	return &SlashingClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "slashing").Logger(),
	}
}

func (c *SlashingClientWrapper) Params(ctx context.Context, req *slashingtypes.QueryParamsRequest) (*slashingtypes.QueryParamsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*slashingtypes.QueryParamsResponse, error) {
		return client.Slashing().Params(ctx, req)
	})
}

func (c *SlashingClientWrapper) SigningInfo(ctx context.Context, req *slashingtypes.QuerySigningInfoRequest) (*slashingtypes.QuerySigningInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*slashingtypes.QuerySigningInfoResponse, error) {
		return client.Slashing().SigningInfo(ctx, req)
	})
}

func (c *SlashingClientWrapper) SigningInfos(ctx context.Context, req *slashingtypes.QuerySigningInfosRequest) (*slashingtypes.QuerySigningInfosResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*slashingtypes.QuerySigningInfosResponse, error) {
		return client.Slashing().SigningInfos(ctx, req)
	})
}
