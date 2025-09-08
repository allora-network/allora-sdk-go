package wrapper

import (
	"context"

	feegrant "cosmossdk.io/x/feegrant"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/pool"
)

// FeegrantClientWrapper wraps the feegrant module with pool management and retry logic
type FeegrantClientWrapper struct {
	poolManager *pool.ClientPoolManager
	logger      zerolog.Logger
}

// NewFeegrantClientWrapper creates a new feegrant client wrapper
func NewFeegrantClientWrapper(poolManager *pool.ClientPoolManager, logger zerolog.Logger) *FeegrantClientWrapper {
	return &FeegrantClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "feegrant").Logger(),
	}
}

func (c *FeegrantClientWrapper) Allowance(ctx context.Context, req *feegrant.QueryAllowanceRequest) (*feegrant.QueryAllowanceResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*feegrant.QueryAllowanceResponse, error) {
		return client.Feegrant().Allowance(ctx, req)
	})
}

func (c *FeegrantClientWrapper) Allowances(ctx context.Context, req *feegrant.QueryAllowancesRequest) (*feegrant.QueryAllowancesResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*feegrant.QueryAllowancesResponse, error) {
		return client.Feegrant().Allowances(ctx, req)
	})
}

func (c *FeegrantClientWrapper) AllowancesByGranter(ctx context.Context, req *feegrant.QueryAllowancesByGranterRequest) (*feegrant.QueryAllowancesByGranterResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*feegrant.QueryAllowancesByGranterResponse, error) {
		return client.Feegrant().AllowancesByGranter(ctx, req)
	})
}
