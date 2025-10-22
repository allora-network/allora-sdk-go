package wrapper

import (
	"context"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/pool"
)

// StakingClientWrapper wraps the staking module with pool management and retry logic
type StakingClientWrapper struct {
	poolManager *pool.ClientPoolManager[interfaces.CosmosClient]
	logger      zerolog.Logger
}

// NewStakingClientWrapper creates a new staking client wrapper
func NewStakingClientWrapper(poolManager *pool.ClientPoolManager[interfaces.CosmosClient], logger zerolog.Logger) *StakingClientWrapper {
	return &StakingClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "staking").Logger(),
	}
}

func (c *StakingClientWrapper) Validators(ctx context.Context, req *stakingtypes.QueryValidatorsRequest, opts ...config.CallOpt) (*stakingtypes.QueryValidatorsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*stakingtypes.QueryValidatorsResponse, error) {
		return client.Staking().Validators(ctx, req, opts...)
	})
}

func (c *StakingClientWrapper) Validator(ctx context.Context, req *stakingtypes.QueryValidatorRequest, opts ...config.CallOpt) (*stakingtypes.QueryValidatorResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*stakingtypes.QueryValidatorResponse, error) {
		return client.Staking().Validator(ctx, req, opts...)
	})
}

func (c *StakingClientWrapper) ValidatorDelegations(ctx context.Context, req *stakingtypes.QueryValidatorDelegationsRequest, opts ...config.CallOpt) (*stakingtypes.QueryValidatorDelegationsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*stakingtypes.QueryValidatorDelegationsResponse, error) {
		return client.Staking().ValidatorDelegations(ctx, req, opts...)
	})
}

func (c *StakingClientWrapper) ValidatorUnbondingDelegations(ctx context.Context, req *stakingtypes.QueryValidatorUnbondingDelegationsRequest, opts ...config.CallOpt) (*stakingtypes.QueryValidatorUnbondingDelegationsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*stakingtypes.QueryValidatorUnbondingDelegationsResponse, error) {
		return client.Staking().ValidatorUnbondingDelegations(ctx, req, opts...)
	})
}

func (c *StakingClientWrapper) Delegation(ctx context.Context, req *stakingtypes.QueryDelegationRequest, opts ...config.CallOpt) (*stakingtypes.QueryDelegationResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*stakingtypes.QueryDelegationResponse, error) {
		return client.Staking().Delegation(ctx, req, opts...)
	})
}

func (c *StakingClientWrapper) UnbondingDelegation(ctx context.Context, req *stakingtypes.QueryUnbondingDelegationRequest, opts ...config.CallOpt) (*stakingtypes.QueryUnbondingDelegationResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*stakingtypes.QueryUnbondingDelegationResponse, error) {
		return client.Staking().UnbondingDelegation(ctx, req, opts...)
	})
}

func (c *StakingClientWrapper) DelegatorDelegations(ctx context.Context, req *stakingtypes.QueryDelegatorDelegationsRequest, opts ...config.CallOpt) (*stakingtypes.QueryDelegatorDelegationsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*stakingtypes.QueryDelegatorDelegationsResponse, error) {
		return client.Staking().DelegatorDelegations(ctx, req, opts...)
	})
}

func (c *StakingClientWrapper) DelegatorUnbondingDelegations(ctx context.Context, req *stakingtypes.QueryDelegatorUnbondingDelegationsRequest, opts ...config.CallOpt) (*stakingtypes.QueryDelegatorUnbondingDelegationsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*stakingtypes.QueryDelegatorUnbondingDelegationsResponse, error) {
		return client.Staking().DelegatorUnbondingDelegations(ctx, req, opts...)
	})
}

func (c *StakingClientWrapper) Redelegations(ctx context.Context, req *stakingtypes.QueryRedelegationsRequest, opts ...config.CallOpt) (*stakingtypes.QueryRedelegationsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*stakingtypes.QueryRedelegationsResponse, error) {
		return client.Staking().Redelegations(ctx, req, opts...)
	})
}

func (c *StakingClientWrapper) DelegatorValidators(ctx context.Context, req *stakingtypes.QueryDelegatorValidatorsRequest, opts ...config.CallOpt) (*stakingtypes.QueryDelegatorValidatorsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*stakingtypes.QueryDelegatorValidatorsResponse, error) {
		return client.Staking().DelegatorValidators(ctx, req, opts...)
	})
}

func (c *StakingClientWrapper) DelegatorValidator(ctx context.Context, req *stakingtypes.QueryDelegatorValidatorRequest, opts ...config.CallOpt) (*stakingtypes.QueryDelegatorValidatorResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*stakingtypes.QueryDelegatorValidatorResponse, error) {
		return client.Staking().DelegatorValidator(ctx, req, opts...)
	})
}

func (c *StakingClientWrapper) HistoricalInfo(ctx context.Context, req *stakingtypes.QueryHistoricalInfoRequest, opts ...config.CallOpt) (*stakingtypes.QueryHistoricalInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*stakingtypes.QueryHistoricalInfoResponse, error) {
		return client.Staking().HistoricalInfo(ctx, req, opts...)
	})
}

func (c *StakingClientWrapper) Pool(ctx context.Context, req *stakingtypes.QueryPoolRequest, opts ...config.CallOpt) (*stakingtypes.QueryPoolResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*stakingtypes.QueryPoolResponse, error) {
		return client.Staking().Pool(ctx, req, opts...)
	})
}

func (c *StakingClientWrapper) Params(ctx context.Context, req *stakingtypes.QueryParamsRequest, opts ...config.CallOpt) (*stakingtypes.QueryParamsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*stakingtypes.QueryParamsResponse, error) {
		return client.Staking().Params(ctx, req, opts...)
	})
}
