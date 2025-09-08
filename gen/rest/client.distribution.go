package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// DistributionRESTClient provides REST access to the distribution module
type DistributionRESTClient struct {
	*RESTClientCore
	logger zerolog.Logger
}

// NewDistributionRESTClient creates a new distribution REST client
func NewDistributionRESTClient(core *RESTClientCore, logger zerolog.Logger) *DistributionRESTClient {
	return &DistributionRESTClient{
		RESTClientCore: core,
		logger:         logger.With().Str("module", "distribution").Str("protocol", "rest").Logger(),
	}
}

func (c *DistributionRESTClient) Params(ctx context.Context, req *distributiontypes.QueryParamsRequest) (*distributiontypes.QueryParamsResponse, error) {
	resp := &distributiontypes.QueryParamsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/distribution/v1beta1/params",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling DistributionRESTClient.Params")
}

func (c *DistributionRESTClient) ValidatorDistributionInfo(ctx context.Context, req *distributiontypes.QueryValidatorDistributionInfoRequest) (*distributiontypes.QueryValidatorDistributionInfoResponse, error) {
	resp := &distributiontypes.QueryValidatorDistributionInfoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/distribution/v1beta1/validators/{validator_address}",
		[]string{"validator_address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling DistributionRESTClient.ValidatorDistributionInfo")
}

func (c *DistributionRESTClient) ValidatorOutstandingRewards(ctx context.Context, req *distributiontypes.QueryValidatorOutstandingRewardsRequest) (*distributiontypes.QueryValidatorOutstandingRewardsResponse, error) {
	resp := &distributiontypes.QueryValidatorOutstandingRewardsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/distribution/v1beta1/validators/{validator_address}/outstanding_rewards",
		[]string{"validator_address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling DistributionRESTClient.ValidatorOutstandingRewards")
}

func (c *DistributionRESTClient) ValidatorCommission(ctx context.Context, req *distributiontypes.QueryValidatorCommissionRequest) (*distributiontypes.QueryValidatorCommissionResponse, error) {
	resp := &distributiontypes.QueryValidatorCommissionResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/distribution/v1beta1/validators/{validator_address}/commission",
		[]string{"validator_address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling DistributionRESTClient.ValidatorCommission")
}

func (c *DistributionRESTClient) ValidatorSlashes(ctx context.Context, req *distributiontypes.QueryValidatorSlashesRequest) (*distributiontypes.QueryValidatorSlashesResponse, error) {
	resp := &distributiontypes.QueryValidatorSlashesResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/distribution/v1beta1/validators/{validator_address}/slashes",
		[]string{"validator_address"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse", "starting_height", "ending_height"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling DistributionRESTClient.ValidatorSlashes")
}

func (c *DistributionRESTClient) DelegationRewards(ctx context.Context, req *distributiontypes.QueryDelegationRewardsRequest) (*distributiontypes.QueryDelegationRewardsResponse, error) {
	resp := &distributiontypes.QueryDelegationRewardsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/distribution/v1beta1/delegators/{delegator_address}/rewards/{validator_address}",
		[]string{"delegator_address", "validator_address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling DistributionRESTClient.DelegationRewards")
}

func (c *DistributionRESTClient) DelegationTotalRewards(ctx context.Context, req *distributiontypes.QueryDelegationTotalRewardsRequest) (*distributiontypes.QueryDelegationTotalRewardsResponse, error) {
	resp := &distributiontypes.QueryDelegationTotalRewardsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/distribution/v1beta1/delegators/{delegator_address}/rewards",
		[]string{"delegator_address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling DistributionRESTClient.DelegationTotalRewards")
}

func (c *DistributionRESTClient) DelegatorValidators(ctx context.Context, req *distributiontypes.QueryDelegatorValidatorsRequest) (*distributiontypes.QueryDelegatorValidatorsResponse, error) {
	resp := &distributiontypes.QueryDelegatorValidatorsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/distribution/v1beta1/delegators/{delegator_address}/validators",
		[]string{"delegator_address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling DistributionRESTClient.DelegatorValidators")
}

func (c *DistributionRESTClient) DelegatorWithdrawAddress(ctx context.Context, req *distributiontypes.QueryDelegatorWithdrawAddressRequest) (*distributiontypes.QueryDelegatorWithdrawAddressResponse, error) {
	resp := &distributiontypes.QueryDelegatorWithdrawAddressResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/distribution/v1beta1/delegators/{delegator_address}/withdraw_address",
		[]string{"delegator_address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling DistributionRESTClient.DelegatorWithdrawAddress")
}

func (c *DistributionRESTClient) CommunityPool(ctx context.Context, req *distributiontypes.QueryCommunityPoolRequest) (*distributiontypes.QueryCommunityPoolResponse, error) {
	resp := &distributiontypes.QueryCommunityPoolResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/distribution/v1beta1/community_pool",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling DistributionRESTClient.CommunityPool")
}
