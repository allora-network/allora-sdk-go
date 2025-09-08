package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

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

func (c *StakingGRPCClient) Validators(ctx context.Context, req *stakingtypes.QueryValidatorsRequest) (*stakingtypes.QueryValidatorsResponse, error) {
	resp, err := c.client.Validators(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Validators")
}

func (c *StakingGRPCClient) Validator(ctx context.Context, req *stakingtypes.QueryValidatorRequest) (*stakingtypes.QueryValidatorResponse, error) {
	resp, err := c.client.Validator(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Validator")
}

func (c *StakingGRPCClient) ValidatorDelegations(ctx context.Context, req *stakingtypes.QueryValidatorDelegationsRequest) (*stakingtypes.QueryValidatorDelegationsResponse, error) {
	resp, err := c.client.ValidatorDelegations(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.ValidatorDelegations")
}

func (c *StakingGRPCClient) ValidatorUnbondingDelegations(ctx context.Context, req *stakingtypes.QueryValidatorUnbondingDelegationsRequest) (*stakingtypes.QueryValidatorUnbondingDelegationsResponse, error) {
	resp, err := c.client.ValidatorUnbondingDelegations(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.ValidatorUnbondingDelegations")
}

func (c *StakingGRPCClient) Delegation(ctx context.Context, req *stakingtypes.QueryDelegationRequest) (*stakingtypes.QueryDelegationResponse, error) {
	resp, err := c.client.Delegation(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Delegation")
}

func (c *StakingGRPCClient) UnbondingDelegation(ctx context.Context, req *stakingtypes.QueryUnbondingDelegationRequest) (*stakingtypes.QueryUnbondingDelegationResponse, error) {
	resp, err := c.client.UnbondingDelegation(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.UnbondingDelegation")
}

func (c *StakingGRPCClient) DelegatorDelegations(ctx context.Context, req *stakingtypes.QueryDelegatorDelegationsRequest) (*stakingtypes.QueryDelegatorDelegationsResponse, error) {
	resp, err := c.client.DelegatorDelegations(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.DelegatorDelegations")
}

func (c *StakingGRPCClient) DelegatorUnbondingDelegations(ctx context.Context, req *stakingtypes.QueryDelegatorUnbondingDelegationsRequest) (*stakingtypes.QueryDelegatorUnbondingDelegationsResponse, error) {
	resp, err := c.client.DelegatorUnbondingDelegations(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.DelegatorUnbondingDelegations")
}

func (c *StakingGRPCClient) Redelegations(ctx context.Context, req *stakingtypes.QueryRedelegationsRequest) (*stakingtypes.QueryRedelegationsResponse, error) {
	resp, err := c.client.Redelegations(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Redelegations")
}

func (c *StakingGRPCClient) DelegatorValidators(ctx context.Context, req *stakingtypes.QueryDelegatorValidatorsRequest) (*stakingtypes.QueryDelegatorValidatorsResponse, error) {
	resp, err := c.client.DelegatorValidators(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.DelegatorValidators")
}

func (c *StakingGRPCClient) DelegatorValidator(ctx context.Context, req *stakingtypes.QueryDelegatorValidatorRequest) (*stakingtypes.QueryDelegatorValidatorResponse, error) {
	resp, err := c.client.DelegatorValidator(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.DelegatorValidator")
}

func (c *StakingGRPCClient) HistoricalInfo(ctx context.Context, req *stakingtypes.QueryHistoricalInfoRequest) (*stakingtypes.QueryHistoricalInfoResponse, error) {
	resp, err := c.client.HistoricalInfo(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.HistoricalInfo")
}

func (c *StakingGRPCClient) Pool(ctx context.Context, req *stakingtypes.QueryPoolRequest) (*stakingtypes.QueryPoolResponse, error) {
	resp, err := c.client.Pool(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Pool")
}

func (c *StakingGRPCClient) Params(ctx context.Context, req *stakingtypes.QueryParamsRequest) (*stakingtypes.QueryParamsResponse, error) {
	resp, err := c.client.Params(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Params")
}
