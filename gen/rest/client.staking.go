package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// StakingRESTClient provides REST access to the staking module
type StakingRESTClient struct {
	*RESTClientCore
	logger zerolog.Logger
}

// NewStakingRESTClient creates a new staking REST client
func NewStakingRESTClient(core *RESTClientCore, logger zerolog.Logger) *StakingRESTClient {
	return &StakingRESTClient{
		RESTClientCore: core,
		logger:         logger.With().Str("module", "staking").Str("protocol", "rest").Logger(),
	}
}

func (c *StakingRESTClient) Validators(ctx context.Context, req *stakingtypes.QueryValidatorsRequest) (*stakingtypes.QueryValidatorsResponse, error) {
	resp := &stakingtypes.QueryValidatorsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/validators",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse", "status"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling StakingRESTClient.Validators")
}

func (c *StakingRESTClient) Validator(ctx context.Context, req *stakingtypes.QueryValidatorRequest) (*stakingtypes.QueryValidatorResponse, error) {
	resp := &stakingtypes.QueryValidatorResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/validators/{validator_addr}",
		[]string{"validator_addr"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling StakingRESTClient.Validator")
}

func (c *StakingRESTClient) ValidatorDelegations(ctx context.Context, req *stakingtypes.QueryValidatorDelegationsRequest) (*stakingtypes.QueryValidatorDelegationsResponse, error) {
	resp := &stakingtypes.QueryValidatorDelegationsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/validators/{validator_addr}/delegations",
		[]string{"validator_addr"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling StakingRESTClient.ValidatorDelegations")
}

func (c *StakingRESTClient) ValidatorUnbondingDelegations(ctx context.Context, req *stakingtypes.QueryValidatorUnbondingDelegationsRequest) (*stakingtypes.QueryValidatorUnbondingDelegationsResponse, error) {
	resp := &stakingtypes.QueryValidatorUnbondingDelegationsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/validators/{validator_addr}/unbonding_delegations",
		[]string{"validator_addr"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling StakingRESTClient.ValidatorUnbondingDelegations")
}

func (c *StakingRESTClient) Delegation(ctx context.Context, req *stakingtypes.QueryDelegationRequest) (*stakingtypes.QueryDelegationResponse, error) {
	resp := &stakingtypes.QueryDelegationResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/validators/{validator_addr}/delegations/{delegator_addr}",
		[]string{"validator_addr", "delegator_addr"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling StakingRESTClient.Delegation")
}

func (c *StakingRESTClient) UnbondingDelegation(ctx context.Context, req *stakingtypes.QueryUnbondingDelegationRequest) (*stakingtypes.QueryUnbondingDelegationResponse, error) {
	resp := &stakingtypes.QueryUnbondingDelegationResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/validators/{validator_addr}/delegations/{delegator_addr}/unbonding_delegation",
		[]string{"validator_addr", "delegator_addr"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling StakingRESTClient.UnbondingDelegation")
}

func (c *StakingRESTClient) DelegatorDelegations(ctx context.Context, req *stakingtypes.QueryDelegatorDelegationsRequest) (*stakingtypes.QueryDelegatorDelegationsResponse, error) {
	resp := &stakingtypes.QueryDelegatorDelegationsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/delegations/{delegator_addr}",
		[]string{"delegator_addr"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling StakingRESTClient.DelegatorDelegations")
}

func (c *StakingRESTClient) DelegatorUnbondingDelegations(ctx context.Context, req *stakingtypes.QueryDelegatorUnbondingDelegationsRequest) (*stakingtypes.QueryDelegatorUnbondingDelegationsResponse, error) {
	resp := &stakingtypes.QueryDelegatorUnbondingDelegationsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/delegators/{delegator_addr}/unbonding_delegations",
		[]string{"delegator_addr"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling StakingRESTClient.DelegatorUnbondingDelegations")
}

func (c *StakingRESTClient) Redelegations(ctx context.Context, req *stakingtypes.QueryRedelegationsRequest) (*stakingtypes.QueryRedelegationsResponse, error) {
	resp := &stakingtypes.QueryRedelegationsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/delegators/{delegator_addr}/redelegations",
		[]string{"delegator_addr"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse", "src_validator_addr", "dst_validator_addr"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling StakingRESTClient.Redelegations")
}

func (c *StakingRESTClient) DelegatorValidators(ctx context.Context, req *stakingtypes.QueryDelegatorValidatorsRequest) (*stakingtypes.QueryDelegatorValidatorsResponse, error) {
	resp := &stakingtypes.QueryDelegatorValidatorsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/delegators/{delegator_addr}/validators",
		[]string{"delegator_addr"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling StakingRESTClient.DelegatorValidators")
}

func (c *StakingRESTClient) DelegatorValidator(ctx context.Context, req *stakingtypes.QueryDelegatorValidatorRequest) (*stakingtypes.QueryDelegatorValidatorResponse, error) {
	resp := &stakingtypes.QueryDelegatorValidatorResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/delegators/{delegator_addr}/validators/{validator_addr}",
		[]string{"delegator_addr", "validator_addr"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling StakingRESTClient.DelegatorValidator")
}

func (c *StakingRESTClient) HistoricalInfo(ctx context.Context, req *stakingtypes.QueryHistoricalInfoRequest) (*stakingtypes.QueryHistoricalInfoResponse, error) {
	resp := &stakingtypes.QueryHistoricalInfoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/historical_info/{height}",
		[]string{"height"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling StakingRESTClient.HistoricalInfo")
}

func (c *StakingRESTClient) Pool(ctx context.Context, req *stakingtypes.QueryPoolRequest) (*stakingtypes.QueryPoolResponse, error) {
	resp := &stakingtypes.QueryPoolResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/pool",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling StakingRESTClient.Pool")
}

func (c *StakingRESTClient) Params(ctx context.Context, req *stakingtypes.QueryParamsRequest) (*stakingtypes.QueryParamsResponse, error) {
	resp := &stakingtypes.QueryParamsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/params",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling StakingRESTClient.Params")
}
