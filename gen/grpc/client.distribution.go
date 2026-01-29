package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// DistributionGRPCClient provides gRPC access to the distribution module
type DistributionGRPCClient struct {
	client distributiontypes.QueryClient
	logger zerolog.Logger
}

var _ interfaces.DistributionClient = (*DistributionGRPCClient)(nil)

// NewDistributionGRPCClient creates a new distribution REST client
func NewDistributionGRPCClient(conn *grpc.ClientConn, logger zerolog.Logger) *DistributionGRPCClient {
	return &DistributionGRPCClient{
		client: distributiontypes.NewQueryClient(conn),
		logger: logger.With().Str("module", "distribution").Str("protocol", "grpc").Logger(),
	}
}

func (c *DistributionGRPCClient) CommunityPool(ctx context.Context, req *distributiontypes.QueryCommunityPoolRequest, opts ...config.CallOpt) (*distributiontypes.QueryCommunityPoolResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.CommunityPool, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling DistributionGRPCClient.CommunityPool")
	}
	return resp, nil
}

func (c *DistributionGRPCClient) DelegationRewards(ctx context.Context, req *distributiontypes.QueryDelegationRewardsRequest, opts ...config.CallOpt) (*distributiontypes.QueryDelegationRewardsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.DelegationRewards, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling DistributionGRPCClient.DelegationRewards")
	}
	return resp, nil
}

func (c *DistributionGRPCClient) DelegationTotalRewards(ctx context.Context, req *distributiontypes.QueryDelegationTotalRewardsRequest, opts ...config.CallOpt) (*distributiontypes.QueryDelegationTotalRewardsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.DelegationTotalRewards, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling DistributionGRPCClient.DelegationTotalRewards")
	}
	return resp, nil
}

func (c *DistributionGRPCClient) DelegatorValidators(ctx context.Context, req *distributiontypes.QueryDelegatorValidatorsRequest, opts ...config.CallOpt) (*distributiontypes.QueryDelegatorValidatorsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.DelegatorValidators, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling DistributionGRPCClient.DelegatorValidators")
	}
	return resp, nil
}

func (c *DistributionGRPCClient) DelegatorWithdrawAddress(ctx context.Context, req *distributiontypes.QueryDelegatorWithdrawAddressRequest, opts ...config.CallOpt) (*distributiontypes.QueryDelegatorWithdrawAddressResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.DelegatorWithdrawAddress, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling DistributionGRPCClient.DelegatorWithdrawAddress")
	}
	return resp, nil
}

func (c *DistributionGRPCClient) Params(ctx context.Context, req *distributiontypes.QueryParamsRequest, opts ...config.CallOpt) (*distributiontypes.QueryParamsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Params, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling DistributionGRPCClient.Params")
	}
	return resp, nil
}

func (c *DistributionGRPCClient) ValidatorCommission(ctx context.Context, req *distributiontypes.QueryValidatorCommissionRequest, opts ...config.CallOpt) (*distributiontypes.QueryValidatorCommissionResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.ValidatorCommission, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling DistributionGRPCClient.ValidatorCommission")
	}
	return resp, nil
}

func (c *DistributionGRPCClient) ValidatorDistributionInfo(ctx context.Context, req *distributiontypes.QueryValidatorDistributionInfoRequest, opts ...config.CallOpt) (*distributiontypes.QueryValidatorDistributionInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.ValidatorDistributionInfo, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling DistributionGRPCClient.ValidatorDistributionInfo")
	}
	return resp, nil
}

func (c *DistributionGRPCClient) ValidatorOutstandingRewards(ctx context.Context, req *distributiontypes.QueryValidatorOutstandingRewardsRequest, opts ...config.CallOpt) (*distributiontypes.QueryValidatorOutstandingRewardsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.ValidatorOutstandingRewards, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling DistributionGRPCClient.ValidatorOutstandingRewards")
	}
	return resp, nil
}

func (c *DistributionGRPCClient) ValidatorSlashes(ctx context.Context, req *distributiontypes.QueryValidatorSlashesRequest, opts ...config.CallOpt) (*distributiontypes.QueryValidatorSlashesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.ValidatorSlashes, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling DistributionGRPCClient.ValidatorSlashes")
	}
	return resp, nil
}
