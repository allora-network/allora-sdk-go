package interfaces

import (
	"context"

	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

type DistributionClient interface {
	Params(ctx context.Context, req *distributiontypes.QueryParamsRequest) (*distributiontypes.QueryParamsResponse, error)
	ValidatorDistributionInfo(ctx context.Context, req *distributiontypes.QueryValidatorDistributionInfoRequest) (*distributiontypes.QueryValidatorDistributionInfoResponse, error)
	ValidatorOutstandingRewards(ctx context.Context, req *distributiontypes.QueryValidatorOutstandingRewardsRequest) (*distributiontypes.QueryValidatorOutstandingRewardsResponse, error)
	ValidatorCommission(ctx context.Context, req *distributiontypes.QueryValidatorCommissionRequest) (*distributiontypes.QueryValidatorCommissionResponse, error)
	ValidatorSlashes(ctx context.Context, req *distributiontypes.QueryValidatorSlashesRequest) (*distributiontypes.QueryValidatorSlashesResponse, error)
	DelegationRewards(ctx context.Context, req *distributiontypes.QueryDelegationRewardsRequest) (*distributiontypes.QueryDelegationRewardsResponse, error)
	DelegationTotalRewards(ctx context.Context, req *distributiontypes.QueryDelegationTotalRewardsRequest) (*distributiontypes.QueryDelegationTotalRewardsResponse, error)
	DelegatorValidators(ctx context.Context, req *distributiontypes.QueryDelegatorValidatorsRequest) (*distributiontypes.QueryDelegatorValidatorsResponse, error)
	DelegatorWithdrawAddress(ctx context.Context, req *distributiontypes.QueryDelegatorWithdrawAddressRequest) (*distributiontypes.QueryDelegatorWithdrawAddressResponse, error)
	CommunityPool(ctx context.Context, req *distributiontypes.QueryCommunityPoolRequest) (*distributiontypes.QueryCommunityPoolResponse, error)
}
