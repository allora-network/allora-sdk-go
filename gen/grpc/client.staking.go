package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// StakingGRPCClient provides gRPC access to the staking module
type StakingGRPCClient struct {
	client stakingtypes.QueryClient
	logger zerolog.Logger
}

var _ interfaces.StakingClient = (*StakingGRPCClient)(nil)

// NewStakingGRPCClient creates a new staking REST client
func NewStakingGRPCClient(conn *grpc.ClientConn, logger zerolog.Logger) *StakingGRPCClient {
	return &StakingGRPCClient{
		client: stakingtypes.NewQueryClient(conn),
		logger: logger.With().Str("module", "staking").Str("protocol", "grpc").Logger(),
	}
}

func (c *StakingGRPCClient) Delegation(ctx context.Context, req *stakingtypes.QueryDelegationRequest, opts ...config.CallOpt) (*stakingtypes.QueryDelegationResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Delegation, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling StakingGRPCClient.Delegation")
	}
	return resp, nil
}

func (c *StakingGRPCClient) DelegatorDelegations(ctx context.Context, req *stakingtypes.QueryDelegatorDelegationsRequest, opts ...config.CallOpt) (*stakingtypes.QueryDelegatorDelegationsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.DelegatorDelegations, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling StakingGRPCClient.DelegatorDelegations")
	}
	return resp, nil
}

func (c *StakingGRPCClient) DelegatorUnbondingDelegations(ctx context.Context, req *stakingtypes.QueryDelegatorUnbondingDelegationsRequest, opts ...config.CallOpt) (*stakingtypes.QueryDelegatorUnbondingDelegationsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.DelegatorUnbondingDelegations, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling StakingGRPCClient.DelegatorUnbondingDelegations")
	}
	return resp, nil
}

func (c *StakingGRPCClient) DelegatorValidator(ctx context.Context, req *stakingtypes.QueryDelegatorValidatorRequest, opts ...config.CallOpt) (*stakingtypes.QueryDelegatorValidatorResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.DelegatorValidator, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling StakingGRPCClient.DelegatorValidator")
	}
	return resp, nil
}

func (c *StakingGRPCClient) DelegatorValidators(ctx context.Context, req *stakingtypes.QueryDelegatorValidatorsRequest, opts ...config.CallOpt) (*stakingtypes.QueryDelegatorValidatorsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.DelegatorValidators, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling StakingGRPCClient.DelegatorValidators")
	}
	return resp, nil
}

func (c *StakingGRPCClient) HistoricalInfo(ctx context.Context, req *stakingtypes.QueryHistoricalInfoRequest, opts ...config.CallOpt) (*stakingtypes.QueryHistoricalInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.HistoricalInfo, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling StakingGRPCClient.HistoricalInfo")
	}
	return resp, nil
}

func (c *StakingGRPCClient) Params(ctx context.Context, req *stakingtypes.QueryParamsRequest, opts ...config.CallOpt) (*stakingtypes.QueryParamsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Params, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling StakingGRPCClient.Params")
	}
	return resp, nil
}

func (c *StakingGRPCClient) Pool(ctx context.Context, req *stakingtypes.QueryPoolRequest, opts ...config.CallOpt) (*stakingtypes.QueryPoolResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Pool, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling StakingGRPCClient.Pool")
	}
	return resp, nil
}

func (c *StakingGRPCClient) Redelegations(ctx context.Context, req *stakingtypes.QueryRedelegationsRequest, opts ...config.CallOpt) (*stakingtypes.QueryRedelegationsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Redelegations, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling StakingGRPCClient.Redelegations")
	}
	return resp, nil
}

func (c *StakingGRPCClient) UnbondingDelegation(ctx context.Context, req *stakingtypes.QueryUnbondingDelegationRequest, opts ...config.CallOpt) (*stakingtypes.QueryUnbondingDelegationResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.UnbondingDelegation, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling StakingGRPCClient.UnbondingDelegation")
	}
	return resp, nil
}

func (c *StakingGRPCClient) Validator(ctx context.Context, req *stakingtypes.QueryValidatorRequest, opts ...config.CallOpt) (*stakingtypes.QueryValidatorResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Validator, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling StakingGRPCClient.Validator")
	}
	return resp, nil
}

func (c *StakingGRPCClient) ValidatorDelegations(ctx context.Context, req *stakingtypes.QueryValidatorDelegationsRequest, opts ...config.CallOpt) (*stakingtypes.QueryValidatorDelegationsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.ValidatorDelegations, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling StakingGRPCClient.ValidatorDelegations")
	}
	return resp, nil
}

func (c *StakingGRPCClient) ValidatorUnbondingDelegations(ctx context.Context, req *stakingtypes.QueryValidatorUnbondingDelegationsRequest, opts ...config.CallOpt) (*stakingtypes.QueryValidatorUnbondingDelegationsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.ValidatorUnbondingDelegations, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling StakingGRPCClient.ValidatorUnbondingDelegations")
	}
	return resp, nil
}

func (c *StakingGRPCClient) Validators(ctx context.Context, req *stakingtypes.QueryValidatorsRequest, opts ...config.CallOpt) (*stakingtypes.QueryValidatorsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Validators, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling StakingGRPCClient.Validators")
	}
	return resp, nil
}
