package wrapper

import (
	"context"

	feegrant "cosmossdk.io/x/feegrant"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/config"
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

func (c *FeegrantClientWrapper) Allowance(ctx context.Context, req *feegrant.QueryAllowanceRequest, opts ...config.CallOpt) (*feegrant.QueryAllowanceResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*feegrant.QueryAllowanceResponse, error) {
		return client.Feegrant().Allowance(ctx, req, opts...)
	})
}

func (c *FeegrantClientWrapper) Allowances(ctx context.Context, req *feegrant.QueryAllowancesRequest, opts ...config.CallOpt) (*feegrant.QueryAllowancesResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*feegrant.QueryAllowancesResponse, error) {
		return client.Feegrant().Allowances(ctx, req, opts...)
	})
}

func (c *FeegrantClientWrapper) AllowancesByGranter(ctx context.Context, req *feegrant.QueryAllowancesByGranterRequest, opts ...config.CallOpt) (*feegrant.QueryAllowancesByGranterResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*feegrant.QueryAllowancesByGranterResponse, error) {
		return client.Feegrant().AllowancesByGranter(ctx, req, opts...)
	})
}
