package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// GovGRPCClient provides gRPC access to the gov module
type GovGRPCClient struct {
	client govv1.QueryClient
	logger zerolog.Logger
}

var _ interfaces.GovClient = (*GovGRPCClient)(nil)

// NewGovGRPCClient creates a new gov REST client
func NewGovGRPCClient(conn *grpc.ClientConn, logger zerolog.Logger) *GovGRPCClient {
	return &GovGRPCClient{
		client: govv1.NewQueryClient(conn),
		logger: logger.With().Str("module", "gov").Str("protocol", "grpc").Logger(),
	}
}

func (c *GovGRPCClient) Constitution(ctx context.Context, req *govv1.QueryConstitutionRequest) (*govv1.QueryConstitutionResponse, error) {
	resp, err := c.client.Constitution(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Constitution")
}

func (c *GovGRPCClient) Proposal(ctx context.Context, req *govv1.QueryProposalRequest) (*govv1.QueryProposalResponse, error) {
	resp, err := c.client.Proposal(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Proposal")
}

func (c *GovGRPCClient) Proposals(ctx context.Context, req *govv1.QueryProposalsRequest) (*govv1.QueryProposalsResponse, error) {
	resp, err := c.client.Proposals(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Proposals")
}

func (c *GovGRPCClient) Vote(ctx context.Context, req *govv1.QueryVoteRequest) (*govv1.QueryVoteResponse, error) {
	resp, err := c.client.Vote(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Vote")
}

func (c *GovGRPCClient) Votes(ctx context.Context, req *govv1.QueryVotesRequest) (*govv1.QueryVotesResponse, error) {
	resp, err := c.client.Votes(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Votes")
}

func (c *GovGRPCClient) Params(ctx context.Context, req *govv1.QueryParamsRequest) (*govv1.QueryParamsResponse, error) {
	resp, err := c.client.Params(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Params")
}

func (c *GovGRPCClient) Deposit(ctx context.Context, req *govv1.QueryDepositRequest) (*govv1.QueryDepositResponse, error) {
	resp, err := c.client.Deposit(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Deposit")
}

func (c *GovGRPCClient) Deposits(ctx context.Context, req *govv1.QueryDepositsRequest) (*govv1.QueryDepositsResponse, error) {
	resp, err := c.client.Deposits(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Deposits")
}

func (c *GovGRPCClient) TallyResult(ctx context.Context, req *govv1.QueryTallyResultRequest) (*govv1.QueryTallyResultResponse, error) {
	resp, err := c.client.TallyResult(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.TallyResult")
}
