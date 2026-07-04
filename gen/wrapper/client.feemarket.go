package wrapper

import (
	"context"

	"github.com/rs/zerolog"
	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/pool"
)

// FeemarketClientWrapper wraps the feemarket module with pool management and retry logic
type FeemarketClientWrapper struct {
	poolManager *pool.ClientPoolManager[interfaces.CosmosClient]
	logger      zerolog.Logger
}

// NewFeemarketClientWrapper creates a new feemarket client wrapper
func NewFeemarketClientWrapper(poolManager *pool.ClientPoolManager[interfaces.CosmosClient], logger zerolog.Logger) *FeemarketClientWrapper {
	return &FeemarketClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "feemarket").Logger(),
	}
}

func (c *FeemarketClientWrapper) GasPrice(ctx context.Context, req *feemarkettypes.GasPriceRequest, opts ...config.CallOpt) (*feemarkettypes.GasPriceResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*feemarkettypes.GasPriceResponse, error) {
		return client.Feemarket().GasPrice(ctx, req, opts...)
	})
}

func (c *FeemarketClientWrapper) GasPrices(ctx context.Context, req *feemarkettypes.GasPricesRequest, opts ...config.CallOpt) (*feemarkettypes.GasPricesResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*feemarkettypes.GasPricesResponse, error) {
		return client.Feemarket().GasPrices(ctx, req, opts...)
	})
}

func (c *FeemarketClientWrapper) Params(ctx context.Context, req *feemarkettypes.ParamsRequest, opts ...config.CallOpt) (*feemarkettypes.ParamsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*feemarkettypes.ParamsResponse, error) {
		return client.Feemarket().Params(ctx, req, opts...)
	})
}

func (c *FeemarketClientWrapper) State(ctx context.Context, req *feemarkettypes.StateRequest, opts ...config.CallOpt) (*feemarkettypes.StateResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*feemarkettypes.StateResponse, error) {
		return client.Feemarket().State(ctx, req, opts...)
	})
}
