package interfaces

import (
	"context"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/allora-network/allora-sdk-go/config"
)

type StakingClient interface {
	Validators(ctx context.Context, req *stakingtypes.QueryValidatorsRequest, opts ...config.CallOpt) (*stakingtypes.QueryValidatorsResponse, error)
	Validator(ctx context.Context, req *stakingtypes.QueryValidatorRequest, opts ...config.CallOpt) (*stakingtypes.QueryValidatorResponse, error)
	ValidatorDelegations(ctx context.Context, req *stakingtypes.QueryValidatorDelegationsRequest, opts ...config.CallOpt) (*stakingtypes.QueryValidatorDelegationsResponse, error)
	ValidatorUnbondingDelegations(ctx context.Context, req *stakingtypes.QueryValidatorUnbondingDelegationsRequest, opts ...config.CallOpt) (*stakingtypes.QueryValidatorUnbondingDelegationsResponse, error)
	Delegation(ctx context.Context, req *stakingtypes.QueryDelegationRequest, opts ...config.CallOpt) (*stakingtypes.QueryDelegationResponse, error)
	UnbondingDelegation(ctx context.Context, req *stakingtypes.QueryUnbondingDelegationRequest, opts ...config.CallOpt) (*stakingtypes.QueryUnbondingDelegationResponse, error)
	DelegatorDelegations(ctx context.Context, req *stakingtypes.QueryDelegatorDelegationsRequest, opts ...config.CallOpt) (*stakingtypes.QueryDelegatorDelegationsResponse, error)
	DelegatorUnbondingDelegations(ctx context.Context, req *stakingtypes.QueryDelegatorUnbondingDelegationsRequest, opts ...config.CallOpt) (*stakingtypes.QueryDelegatorUnbondingDelegationsResponse, error)
	Redelegations(ctx context.Context, req *stakingtypes.QueryRedelegationsRequest, opts ...config.CallOpt) (*stakingtypes.QueryRedelegationsResponse, error)
	DelegatorValidators(ctx context.Context, req *stakingtypes.QueryDelegatorValidatorsRequest, opts ...config.CallOpt) (*stakingtypes.QueryDelegatorValidatorsResponse, error)
	DelegatorValidator(ctx context.Context, req *stakingtypes.QueryDelegatorValidatorRequest, opts ...config.CallOpt) (*stakingtypes.QueryDelegatorValidatorResponse, error)
	HistoricalInfo(ctx context.Context, req *stakingtypes.QueryHistoricalInfoRequest, opts ...config.CallOpt) (*stakingtypes.QueryHistoricalInfoResponse, error)
	Pool(ctx context.Context, req *stakingtypes.QueryPoolRequest, opts ...config.CallOpt) (*stakingtypes.QueryPoolResponse, error)
	Params(ctx context.Context, req *stakingtypes.QueryParamsRequest, opts ...config.CallOpt) (*stakingtypes.QueryParamsResponse, error)
}
