package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/allora-network/allora-sdk-go/config"
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

func (c *StakingRESTClient) Delegation(ctx context.Context, req *stakingtypes.QueryDelegationRequest, opts ...config.CallOpt) (*stakingtypes.QueryDelegationResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &stakingtypes.QueryDelegationResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/validators/{validator_addr}/delegations/{delegator_addr}",
		[]string{"validator_addr", "delegator_addr"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling StakingRESTClient.Delegation")
	}
	return resp, nil
}

func (c *StakingRESTClient) DelegatorDelegations(ctx context.Context, req *stakingtypes.QueryDelegatorDelegationsRequest, opts ...config.CallOpt) (*stakingtypes.QueryDelegatorDelegationsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &stakingtypes.QueryDelegatorDelegationsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/delegations/{delegator_addr}",
		[]string{"delegator_addr"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling StakingRESTClient.DelegatorDelegations")
	}
	return resp, nil
}

func (c *StakingRESTClient) DelegatorUnbondingDelegations(ctx context.Context, req *stakingtypes.QueryDelegatorUnbondingDelegationsRequest, opts ...config.CallOpt) (*stakingtypes.QueryDelegatorUnbondingDelegationsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &stakingtypes.QueryDelegatorUnbondingDelegationsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/delegators/{delegator_addr}/unbonding_delegations",
		[]string{"delegator_addr"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling StakingRESTClient.DelegatorUnbondingDelegations")
	}
	return resp, nil
}

func (c *StakingRESTClient) DelegatorValidator(ctx context.Context, req *stakingtypes.QueryDelegatorValidatorRequest, opts ...config.CallOpt) (*stakingtypes.QueryDelegatorValidatorResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &stakingtypes.QueryDelegatorValidatorResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/delegators/{delegator_addr}/validators/{validator_addr}",
		[]string{"delegator_addr", "validator_addr"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling StakingRESTClient.DelegatorValidator")
	}
	return resp, nil
}

func (c *StakingRESTClient) DelegatorValidators(ctx context.Context, req *stakingtypes.QueryDelegatorValidatorsRequest, opts ...config.CallOpt) (*stakingtypes.QueryDelegatorValidatorsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &stakingtypes.QueryDelegatorValidatorsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/delegators/{delegator_addr}/validators",
		[]string{"delegator_addr"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling StakingRESTClient.DelegatorValidators")
	}
	return resp, nil
}

func (c *StakingRESTClient) HistoricalInfo(ctx context.Context, req *stakingtypes.QueryHistoricalInfoRequest, opts ...config.CallOpt) (*stakingtypes.QueryHistoricalInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &stakingtypes.QueryHistoricalInfoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/historical_info/{height}",
		[]string{"height"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling StakingRESTClient.HistoricalInfo")
	}
	return resp, nil
}

func (c *StakingRESTClient) Params(ctx context.Context, req *stakingtypes.QueryParamsRequest, opts ...config.CallOpt) (*stakingtypes.QueryParamsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &stakingtypes.QueryParamsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/params",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling StakingRESTClient.Params")
	}
	return resp, nil
}

func (c *StakingRESTClient) Pool(ctx context.Context, req *stakingtypes.QueryPoolRequest, opts ...config.CallOpt) (*stakingtypes.QueryPoolResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &stakingtypes.QueryPoolResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/pool",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling StakingRESTClient.Pool")
	}
	return resp, nil
}

func (c *StakingRESTClient) Redelegations(ctx context.Context, req *stakingtypes.QueryRedelegationsRequest, opts ...config.CallOpt) (*stakingtypes.QueryRedelegationsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &stakingtypes.QueryRedelegationsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/delegators/{delegator_addr}/redelegations",
		[]string{"delegator_addr"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse", "src_validator_addr", "dst_validator_addr"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling StakingRESTClient.Redelegations")
	}
	return resp, nil
}

func (c *StakingRESTClient) UnbondingDelegation(ctx context.Context, req *stakingtypes.QueryUnbondingDelegationRequest, opts ...config.CallOpt) (*stakingtypes.QueryUnbondingDelegationResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &stakingtypes.QueryUnbondingDelegationResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/validators/{validator_addr}/delegations/{delegator_addr}/unbonding_delegation",
		[]string{"validator_addr", "delegator_addr"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling StakingRESTClient.UnbondingDelegation")
	}
	return resp, nil
}

func (c *StakingRESTClient) Validator(ctx context.Context, req *stakingtypes.QueryValidatorRequest, opts ...config.CallOpt) (*stakingtypes.QueryValidatorResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &stakingtypes.QueryValidatorResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/validators/{validator_addr}",
		[]string{"validator_addr"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling StakingRESTClient.Validator")
	}
	return resp, nil
}

func (c *StakingRESTClient) ValidatorDelegations(ctx context.Context, req *stakingtypes.QueryValidatorDelegationsRequest, opts ...config.CallOpt) (*stakingtypes.QueryValidatorDelegationsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &stakingtypes.QueryValidatorDelegationsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/validators/{validator_addr}/delegations",
		[]string{"validator_addr"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling StakingRESTClient.ValidatorDelegations")
	}
	return resp, nil
}

func (c *StakingRESTClient) ValidatorUnbondingDelegations(ctx context.Context, req *stakingtypes.QueryValidatorUnbondingDelegationsRequest, opts ...config.CallOpt) (*stakingtypes.QueryValidatorUnbondingDelegationsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &stakingtypes.QueryValidatorUnbondingDelegationsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/validators/{validator_addr}/unbonding_delegations",
		[]string{"validator_addr"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling StakingRESTClient.ValidatorUnbondingDelegations")
	}
	return resp, nil
}

func (c *StakingRESTClient) Validators(ctx context.Context, req *stakingtypes.QueryValidatorsRequest, opts ...config.CallOpt) (*stakingtypes.QueryValidatorsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &stakingtypes.QueryValidatorsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/staking/v1beta1/validators",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse", "status"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling StakingRESTClient.Validators")
	}
	return resp, nil
}
