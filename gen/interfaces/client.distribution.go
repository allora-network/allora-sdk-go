package interfaces

import (
	"context"

	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/allora-network/allora-sdk-go/config"
)

type DistributionClient interface {
	CommunityPool(ctx context.Context, req *distributiontypes.QueryCommunityPoolRequest, opts ...config.CallOpt) (*distributiontypes.QueryCommunityPoolResponse, error)
	DelegationRewards(ctx context.Context, req *distributiontypes.QueryDelegationRewardsRequest, opts ...config.CallOpt) (*distributiontypes.QueryDelegationRewardsResponse, error)
	DelegationTotalRewards(ctx context.Context, req *distributiontypes.QueryDelegationTotalRewardsRequest, opts ...config.CallOpt) (*distributiontypes.QueryDelegationTotalRewardsResponse, error)
	DelegatorValidators(ctx context.Context, req *distributiontypes.QueryDelegatorValidatorsRequest, opts ...config.CallOpt) (*distributiontypes.QueryDelegatorValidatorsResponse, error)
	DelegatorWithdrawAddress(ctx context.Context, req *distributiontypes.QueryDelegatorWithdrawAddressRequest, opts ...config.CallOpt) (*distributiontypes.QueryDelegatorWithdrawAddressResponse, error)
	Params(ctx context.Context, req *distributiontypes.QueryParamsRequest, opts ...config.CallOpt) (*distributiontypes.QueryParamsResponse, error)
	ValidatorCommission(ctx context.Context, req *distributiontypes.QueryValidatorCommissionRequest, opts ...config.CallOpt) (*distributiontypes.QueryValidatorCommissionResponse, error)
	ValidatorDistributionInfo(ctx context.Context, req *distributiontypes.QueryValidatorDistributionInfoRequest, opts ...config.CallOpt) (*distributiontypes.QueryValidatorDistributionInfoResponse, error)
	ValidatorOutstandingRewards(ctx context.Context, req *distributiontypes.QueryValidatorOutstandingRewardsRequest, opts ...config.CallOpt) (*distributiontypes.QueryValidatorOutstandingRewardsResponse, error)
	ValidatorSlashes(ctx context.Context, req *distributiontypes.QueryValidatorSlashesRequest, opts ...config.CallOpt) (*distributiontypes.QueryValidatorSlashesResponse, error)
}
