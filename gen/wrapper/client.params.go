package wrapper

import (
	"context"

	proposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/pool"
)

// ParamsClientWrapper wraps the params module with pool management and retry logic
type ParamsClientWrapper struct {
	poolManager *pool.ClientPoolManager
	logger      zerolog.Logger
}

// NewParamsClientWrapper creates a new params client wrapper
func NewParamsClientWrapper(poolManager *pool.ClientPoolManager, logger zerolog.Logger) *ParamsClientWrapper {
	return &ParamsClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "params").Logger(),
	}
}

func (c *ParamsClientWrapper) Params(ctx context.Context, req *proposal.QueryParamsRequest, opts ...config.CallOpt) (*proposal.QueryParamsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*proposal.QueryParamsResponse, error) {
		return client.Params().Params(ctx, req, opts...)
	})
}

func (c *ParamsClientWrapper) Subspaces(ctx context.Context, req *proposal.QuerySubspacesRequest, opts ...config.CallOpt) (*proposal.QuerySubspacesResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*proposal.QuerySubspacesResponse, error) {
		return client.Params().Subspaces(ctx, req, opts...)
	})
}
