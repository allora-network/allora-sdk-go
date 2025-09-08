package wrapper

import (
	"context"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/pool"
)

// StakingClientWrapper wraps the staking module with pool management and retry logic
type StakingClientWrapper struct {
	poolManager *pool.ClientPoolManager
	logger      zerolog.Logger
}

// NewStakingClientWrapper creates a new staking client wrapper
func NewStakingClientWrapper(poolManager *pool.ClientPoolManager, logger zerolog.Logger) *StakingClientWrapper {
	return &StakingClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "staking").Logger(),
	}
}

func (c *StakingClientWrapper) Validators(ctx context.Context, req *stakingtypes.QueryValidatorsRequest) (*stakingtypes.QueryValidatorsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*stakingtypes.QueryValidatorsResponse, error) {
		return client.Staking().Validators(ctx, req)
	})
}

func (c *StakingClientWrapper) Validator(ctx context.Context, req *stakingtypes.QueryValidatorRequest) (*stakingtypes.QueryValidatorResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*stakingtypes.QueryValidatorResponse, error) {
		return client.Staking().Validator(ctx, req)
	})
}

func (c *StakingClientWrapper) ValidatorDelegations(ctx context.Context, req *stakingtypes.QueryValidatorDelegationsRequest) (*stakingtypes.QueryValidatorDelegationsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*stakingtypes.QueryValidatorDelegationsResponse, error) {
		return client.Staking().ValidatorDelegations(ctx, req)
	})
}

func (c *StakingClientWrapper) ValidatorUnbondingDelegations(ctx context.Context, req *stakingtypes.QueryValidatorUnbondingDelegationsRequest) (*stakingtypes.QueryValidatorUnbondingDelegationsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*stakingtypes.QueryValidatorUnbondingDelegationsResponse, error) {
		return client.Staking().ValidatorUnbondingDelegations(ctx, req)
	})
}

func (c *StakingClientWrapper) Delegation(ctx context.Context, req *stakingtypes.QueryDelegationRequest) (*stakingtypes.QueryDelegationResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*stakingtypes.QueryDelegationResponse, error) {
		return client.Staking().Delegation(ctx, req)
	})
}

func (c *StakingClientWrapper) UnbondingDelegation(ctx context.Context, req *stakingtypes.QueryUnbondingDelegationRequest) (*stakingtypes.QueryUnbondingDelegationResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*stakingtypes.QueryUnbondingDelegationResponse, error) {
		return client.Staking().UnbondingDelegation(ctx, req)
	})
}

func (c *StakingClientWrapper) DelegatorDelegations(ctx context.Context, req *stakingtypes.QueryDelegatorDelegationsRequest) (*stakingtypes.QueryDelegatorDelegationsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*stakingtypes.QueryDelegatorDelegationsResponse, error) {
		return client.Staking().DelegatorDelegations(ctx, req)
	})
}

func (c *StakingClientWrapper) DelegatorUnbondingDelegations(ctx context.Context, req *stakingtypes.QueryDelegatorUnbondingDelegationsRequest) (*stakingtypes.QueryDelegatorUnbondingDelegationsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*stakingtypes.QueryDelegatorUnbondingDelegationsResponse, error) {
		return client.Staking().DelegatorUnbondingDelegations(ctx, req)
	})
}

func (c *StakingClientWrapper) Redelegations(ctx context.Context, req *stakingtypes.QueryRedelegationsRequest) (*stakingtypes.QueryRedelegationsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*stakingtypes.QueryRedelegationsResponse, error) {
		return client.Staking().Redelegations(ctx, req)
	})
}

func (c *StakingClientWrapper) DelegatorValidators(ctx context.Context, req *stakingtypes.QueryDelegatorValidatorsRequest) (*stakingtypes.QueryDelegatorValidatorsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*stakingtypes.QueryDelegatorValidatorsResponse, error) {
		return client.Staking().DelegatorValidators(ctx, req)
	})
}

func (c *StakingClientWrapper) DelegatorValidator(ctx context.Context, req *stakingtypes.QueryDelegatorValidatorRequest) (*stakingtypes.QueryDelegatorValidatorResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*stakingtypes.QueryDelegatorValidatorResponse, error) {
		return client.Staking().DelegatorValidator(ctx, req)
	})
}

func (c *StakingClientWrapper) HistoricalInfo(ctx context.Context, req *stakingtypes.QueryHistoricalInfoRequest) (*stakingtypes.QueryHistoricalInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*stakingtypes.QueryHistoricalInfoResponse, error) {
		return client.Staking().HistoricalInfo(ctx, req)
	})
}

func (c *StakingClientWrapper) Pool(ctx context.Context, req *stakingtypes.QueryPoolRequest) (*stakingtypes.QueryPoolResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*stakingtypes.QueryPoolResponse, error) {
		return client.Staking().Pool(ctx, req)
	})
}

func (c *StakingClientWrapper) Params(ctx context.Context, req *stakingtypes.QueryParamsRequest) (*stakingtypes.QueryParamsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*stakingtypes.QueryParamsResponse, error) {
		return client.Staking().Params(ctx, req)
	})
}
