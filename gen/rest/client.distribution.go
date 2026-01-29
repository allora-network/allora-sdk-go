package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/allora-network/allora-sdk-go/config"
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

func (c *DistributionRESTClient) CommunityPool(ctx context.Context, req *distributiontypes.QueryCommunityPoolRequest, opts ...config.CallOpt) (*distributiontypes.QueryCommunityPoolResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &distributiontypes.QueryCommunityPoolResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/distribution/v1beta1/community_pool",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling DistributionRESTClient.CommunityPool")
	}
	return resp, nil
}

func (c *DistributionRESTClient) DelegationRewards(ctx context.Context, req *distributiontypes.QueryDelegationRewardsRequest, opts ...config.CallOpt) (*distributiontypes.QueryDelegationRewardsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &distributiontypes.QueryDelegationRewardsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/distribution/v1beta1/delegators/{delegator_address}/rewards/{validator_address}",
		[]string{"delegator_address", "validator_address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling DistributionRESTClient.DelegationRewards")
	}
	return resp, nil
}

func (c *DistributionRESTClient) DelegationTotalRewards(ctx context.Context, req *distributiontypes.QueryDelegationTotalRewardsRequest, opts ...config.CallOpt) (*distributiontypes.QueryDelegationTotalRewardsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &distributiontypes.QueryDelegationTotalRewardsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/distribution/v1beta1/delegators/{delegator_address}/rewards",
		[]string{"delegator_address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling DistributionRESTClient.DelegationTotalRewards")
	}
	return resp, nil
}

func (c *DistributionRESTClient) DelegatorValidators(ctx context.Context, req *distributiontypes.QueryDelegatorValidatorsRequest, opts ...config.CallOpt) (*distributiontypes.QueryDelegatorValidatorsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &distributiontypes.QueryDelegatorValidatorsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/distribution/v1beta1/delegators/{delegator_address}/validators",
		[]string{"delegator_address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling DistributionRESTClient.DelegatorValidators")
	}
	return resp, nil
}

func (c *DistributionRESTClient) DelegatorWithdrawAddress(ctx context.Context, req *distributiontypes.QueryDelegatorWithdrawAddressRequest, opts ...config.CallOpt) (*distributiontypes.QueryDelegatorWithdrawAddressResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &distributiontypes.QueryDelegatorWithdrawAddressResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/distribution/v1beta1/delegators/{delegator_address}/withdraw_address",
		[]string{"delegator_address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling DistributionRESTClient.DelegatorWithdrawAddress")
	}
	return resp, nil
}

func (c *DistributionRESTClient) Params(ctx context.Context, req *distributiontypes.QueryParamsRequest, opts ...config.CallOpt) (*distributiontypes.QueryParamsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &distributiontypes.QueryParamsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/distribution/v1beta1/params",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling DistributionRESTClient.Params")
	}
	return resp, nil
}

func (c *DistributionRESTClient) ValidatorCommission(ctx context.Context, req *distributiontypes.QueryValidatorCommissionRequest, opts ...config.CallOpt) (*distributiontypes.QueryValidatorCommissionResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &distributiontypes.QueryValidatorCommissionResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/distribution/v1beta1/validators/{validator_address}/commission",
		[]string{"validator_address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling DistributionRESTClient.ValidatorCommission")
	}
	return resp, nil
}

func (c *DistributionRESTClient) ValidatorDistributionInfo(ctx context.Context, req *distributiontypes.QueryValidatorDistributionInfoRequest, opts ...config.CallOpt) (*distributiontypes.QueryValidatorDistributionInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &distributiontypes.QueryValidatorDistributionInfoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/distribution/v1beta1/validators/{validator_address}",
		[]string{"validator_address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling DistributionRESTClient.ValidatorDistributionInfo")
	}
	return resp, nil
}

func (c *DistributionRESTClient) ValidatorOutstandingRewards(ctx context.Context, req *distributiontypes.QueryValidatorOutstandingRewardsRequest, opts ...config.CallOpt) (*distributiontypes.QueryValidatorOutstandingRewardsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &distributiontypes.QueryValidatorOutstandingRewardsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/distribution/v1beta1/validators/{validator_address}/outstanding_rewards",
		[]string{"validator_address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling DistributionRESTClient.ValidatorOutstandingRewards")
	}
	return resp, nil
}

func (c *DistributionRESTClient) ValidatorSlashes(ctx context.Context, req *distributiontypes.QueryValidatorSlashesRequest, opts ...config.CallOpt) (*distributiontypes.QueryValidatorSlashesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &distributiontypes.QueryValidatorSlashesResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/distribution/v1beta1/validators/{validator_address}/slashes",
		[]string{"validator_address"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse", "starting_height", "ending_height"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling DistributionRESTClient.ValidatorSlashes")
	}
	return resp, nil
}
