package wrapper

import (
	"context"

	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/pool"
)

// DistributionClientWrapper wraps the distribution module with pool management and retry logic
type DistributionClientWrapper struct {
	poolManager *pool.ClientPoolManager
	logger      zerolog.Logger
}

// NewDistributionClientWrapper creates a new distribution client wrapper
func NewDistributionClientWrapper(poolManager *pool.ClientPoolManager, logger zerolog.Logger) *DistributionClientWrapper {
	return &DistributionClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "distribution").Logger(),
	}
}

func (c *DistributionClientWrapper) Params(ctx context.Context, req *distributiontypes.QueryParamsRequest, opts ...config.CallOpt) (*distributiontypes.QueryParamsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*distributiontypes.QueryParamsResponse, error) {
		return client.Distribution().Params(ctx, req, opts...)
	})
}

func (c *DistributionClientWrapper) ValidatorDistributionInfo(ctx context.Context, req *distributiontypes.QueryValidatorDistributionInfoRequest, opts ...config.CallOpt) (*distributiontypes.QueryValidatorDistributionInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*distributiontypes.QueryValidatorDistributionInfoResponse, error) {
		return client.Distribution().ValidatorDistributionInfo(ctx, req, opts...)
	})
}

func (c *DistributionClientWrapper) ValidatorOutstandingRewards(ctx context.Context, req *distributiontypes.QueryValidatorOutstandingRewardsRequest, opts ...config.CallOpt) (*distributiontypes.QueryValidatorOutstandingRewardsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*distributiontypes.QueryValidatorOutstandingRewardsResponse, error) {
		return client.Distribution().ValidatorOutstandingRewards(ctx, req, opts...)
	})
}

func (c *DistributionClientWrapper) ValidatorCommission(ctx context.Context, req *distributiontypes.QueryValidatorCommissionRequest, opts ...config.CallOpt) (*distributiontypes.QueryValidatorCommissionResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*distributiontypes.QueryValidatorCommissionResponse, error) {
		return client.Distribution().ValidatorCommission(ctx, req, opts...)
	})
}

func (c *DistributionClientWrapper) ValidatorSlashes(ctx context.Context, req *distributiontypes.QueryValidatorSlashesRequest, opts ...config.CallOpt) (*distributiontypes.QueryValidatorSlashesResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*distributiontypes.QueryValidatorSlashesResponse, error) {
		return client.Distribution().ValidatorSlashes(ctx, req, opts...)
	})
}

func (c *DistributionClientWrapper) DelegationRewards(ctx context.Context, req *distributiontypes.QueryDelegationRewardsRequest, opts ...config.CallOpt) (*distributiontypes.QueryDelegationRewardsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*distributiontypes.QueryDelegationRewardsResponse, error) {
		return client.Distribution().DelegationRewards(ctx, req, opts...)
	})
}

func (c *DistributionClientWrapper) DelegationTotalRewards(ctx context.Context, req *distributiontypes.QueryDelegationTotalRewardsRequest, opts ...config.CallOpt) (*distributiontypes.QueryDelegationTotalRewardsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*distributiontypes.QueryDelegationTotalRewardsResponse, error) {
		return client.Distribution().DelegationTotalRewards(ctx, req, opts...)
	})
}

func (c *DistributionClientWrapper) DelegatorValidators(ctx context.Context, req *distributiontypes.QueryDelegatorValidatorsRequest, opts ...config.CallOpt) (*distributiontypes.QueryDelegatorValidatorsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*distributiontypes.QueryDelegatorValidatorsResponse, error) {
		return client.Distribution().DelegatorValidators(ctx, req, opts...)
	})
}

func (c *DistributionClientWrapper) DelegatorWithdrawAddress(ctx context.Context, req *distributiontypes.QueryDelegatorWithdrawAddressRequest, opts ...config.CallOpt) (*distributiontypes.QueryDelegatorWithdrawAddressResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*distributiontypes.QueryDelegatorWithdrawAddressResponse, error) {
		return client.Distribution().DelegatorWithdrawAddress(ctx, req, opts...)
	})
}

func (c *DistributionClientWrapper) CommunityPool(ctx context.Context, req *distributiontypes.QueryCommunityPoolRequest, opts ...config.CallOpt) (*distributiontypes.QueryCommunityPoolResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*distributiontypes.QueryCommunityPoolResponse, error) {
		return client.Distribution().CommunityPool(ctx, req, opts...)
	})
}
