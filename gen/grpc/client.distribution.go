package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

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

func (c *DistributionGRPCClient) Params(ctx context.Context, req *distributiontypes.QueryParamsRequest) (*distributiontypes.QueryParamsResponse, error) {
	resp, err := c.client.Params(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Params")
}

func (c *DistributionGRPCClient) ValidatorDistributionInfo(ctx context.Context, req *distributiontypes.QueryValidatorDistributionInfoRequest) (*distributiontypes.QueryValidatorDistributionInfoResponse, error) {
	resp, err := c.client.ValidatorDistributionInfo(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.ValidatorDistributionInfo")
}

func (c *DistributionGRPCClient) ValidatorOutstandingRewards(ctx context.Context, req *distributiontypes.QueryValidatorOutstandingRewardsRequest) (*distributiontypes.QueryValidatorOutstandingRewardsResponse, error) {
	resp, err := c.client.ValidatorOutstandingRewards(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.ValidatorOutstandingRewards")
}

func (c *DistributionGRPCClient) ValidatorCommission(ctx context.Context, req *distributiontypes.QueryValidatorCommissionRequest) (*distributiontypes.QueryValidatorCommissionResponse, error) {
	resp, err := c.client.ValidatorCommission(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.ValidatorCommission")
}

func (c *DistributionGRPCClient) ValidatorSlashes(ctx context.Context, req *distributiontypes.QueryValidatorSlashesRequest) (*distributiontypes.QueryValidatorSlashesResponse, error) {
	resp, err := c.client.ValidatorSlashes(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.ValidatorSlashes")
}

func (c *DistributionGRPCClient) DelegationRewards(ctx context.Context, req *distributiontypes.QueryDelegationRewardsRequest) (*distributiontypes.QueryDelegationRewardsResponse, error) {
	resp, err := c.client.DelegationRewards(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.DelegationRewards")
}

func (c *DistributionGRPCClient) DelegationTotalRewards(ctx context.Context, req *distributiontypes.QueryDelegationTotalRewardsRequest) (*distributiontypes.QueryDelegationTotalRewardsResponse, error) {
	resp, err := c.client.DelegationTotalRewards(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.DelegationTotalRewards")
}

func (c *DistributionGRPCClient) DelegatorValidators(ctx context.Context, req *distributiontypes.QueryDelegatorValidatorsRequest) (*distributiontypes.QueryDelegatorValidatorsResponse, error) {
	resp, err := c.client.DelegatorValidators(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.DelegatorValidators")
}

func (c *DistributionGRPCClient) DelegatorWithdrawAddress(ctx context.Context, req *distributiontypes.QueryDelegatorWithdrawAddressRequest) (*distributiontypes.QueryDelegatorWithdrawAddressResponse, error) {
	resp, err := c.client.DelegatorWithdrawAddress(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.DelegatorWithdrawAddress")
}

func (c *DistributionGRPCClient) CommunityPool(ctx context.Context, req *distributiontypes.QueryCommunityPoolRequest) (*distributiontypes.QueryCommunityPoolResponse, error) {
	resp, err := c.client.CommunityPool(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.CommunityPool")
}
