package wrapper

import (
	"context"

	proposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
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

func (c *ParamsClientWrapper) Params(ctx context.Context, req *proposal.QueryParamsRequest) (*proposal.QueryParamsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*proposal.QueryParamsResponse, error) {
		return client.Params().Params(ctx, req)
	})
}

func (c *ParamsClientWrapper) Subspaces(ctx context.Context, req *proposal.QuerySubspacesRequest) (*proposal.QuerySubspacesResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*proposal.QuerySubspacesResponse, error) {
		return client.Params().Subspaces(ctx, req)
	})
}
