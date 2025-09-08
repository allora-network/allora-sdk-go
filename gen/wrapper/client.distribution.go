package wrapper

import (
	"context"

	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
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

func (c *DistributionClientWrapper) Params(ctx context.Context, req *distributiontypes.QueryParamsRequest) (*distributiontypes.QueryParamsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*distributiontypes.QueryParamsResponse, error) {
		return client.Distribution().Params(ctx, req)
	})
}

func (c *DistributionClientWrapper) ValidatorDistributionInfo(ctx context.Context, req *distributiontypes.QueryValidatorDistributionInfoRequest) (*distributiontypes.QueryValidatorDistributionInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*distributiontypes.QueryValidatorDistributionInfoResponse, error) {
		return client.Distribution().ValidatorDistributionInfo(ctx, req)
	})
}

func (c *DistributionClientWrapper) ValidatorOutstandingRewards(ctx context.Context, req *distributiontypes.QueryValidatorOutstandingRewardsRequest) (*distributiontypes.QueryValidatorOutstandingRewardsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*distributiontypes.QueryValidatorOutstandingRewardsResponse, error) {
		return client.Distribution().ValidatorOutstandingRewards(ctx, req)
	})
}

func (c *DistributionClientWrapper) ValidatorCommission(ctx context.Context, req *distributiontypes.QueryValidatorCommissionRequest) (*distributiontypes.QueryValidatorCommissionResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*distributiontypes.QueryValidatorCommissionResponse, error) {
		return client.Distribution().ValidatorCommission(ctx, req)
	})
}

func (c *DistributionClientWrapper) ValidatorSlashes(ctx context.Context, req *distributiontypes.QueryValidatorSlashesRequest) (*distributiontypes.QueryValidatorSlashesResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*distributiontypes.QueryValidatorSlashesResponse, error) {
		return client.Distribution().ValidatorSlashes(ctx, req)
	})
}

func (c *DistributionClientWrapper) DelegationRewards(ctx context.Context, req *distributiontypes.QueryDelegationRewardsRequest) (*distributiontypes.QueryDelegationRewardsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*distributiontypes.QueryDelegationRewardsResponse, error) {
		return client.Distribution().DelegationRewards(ctx, req)
	})
}

func (c *DistributionClientWrapper) DelegationTotalRewards(ctx context.Context, req *distributiontypes.QueryDelegationTotalRewardsRequest) (*distributiontypes.QueryDelegationTotalRewardsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*distributiontypes.QueryDelegationTotalRewardsResponse, error) {
		return client.Distribution().DelegationTotalRewards(ctx, req)
	})
}

func (c *DistributionClientWrapper) DelegatorValidators(ctx context.Context, req *distributiontypes.QueryDelegatorValidatorsRequest) (*distributiontypes.QueryDelegatorValidatorsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*distributiontypes.QueryDelegatorValidatorsResponse, error) {
		return client.Distribution().DelegatorValidators(ctx, req)
	})
}

func (c *DistributionClientWrapper) DelegatorWithdrawAddress(ctx context.Context, req *distributiontypes.QueryDelegatorWithdrawAddressRequest) (*distributiontypes.QueryDelegatorWithdrawAddressResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*distributiontypes.QueryDelegatorWithdrawAddressResponse, error) {
		return client.Distribution().DelegatorWithdrawAddress(ctx, req)
	})
}

func (c *DistributionClientWrapper) CommunityPool(ctx context.Context, req *distributiontypes.QueryCommunityPoolRequest) (*distributiontypes.QueryCommunityPoolResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*distributiontypes.QueryCommunityPoolResponse, error) {
		return client.Distribution().CommunityPool(ctx, req)
	})
}
