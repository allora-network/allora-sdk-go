package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/allora-network/allora-sdk-go/config"
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

func (c *GovGRPCClient) Constitution(ctx context.Context, req *govv1.QueryConstitutionRequest, opts ...config.CallOpt) (*govv1.QueryConstitutionResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Constitution, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling GovGRPCClient.Constitution")
	}
	return resp, nil
}

func (c *GovGRPCClient) Deposit(ctx context.Context, req *govv1.QueryDepositRequest, opts ...config.CallOpt) (*govv1.QueryDepositResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Deposit, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling GovGRPCClient.Deposit")
	}
	return resp, nil
}

func (c *GovGRPCClient) Deposits(ctx context.Context, req *govv1.QueryDepositsRequest, opts ...config.CallOpt) (*govv1.QueryDepositsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Deposits, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling GovGRPCClient.Deposits")
	}
	return resp, nil
}

func (c *GovGRPCClient) Params(ctx context.Context, req *govv1.QueryParamsRequest, opts ...config.CallOpt) (*govv1.QueryParamsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Params, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling GovGRPCClient.Params")
	}
	return resp, nil
}

func (c *GovGRPCClient) Proposal(ctx context.Context, req *govv1.QueryProposalRequest, opts ...config.CallOpt) (*govv1.QueryProposalResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Proposal, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling GovGRPCClient.Proposal")
	}
	return resp, nil
}

func (c *GovGRPCClient) Proposals(ctx context.Context, req *govv1.QueryProposalsRequest, opts ...config.CallOpt) (*govv1.QueryProposalsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Proposals, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling GovGRPCClient.Proposals")
	}
	return resp, nil
}

func (c *GovGRPCClient) TallyResult(ctx context.Context, req *govv1.QueryTallyResultRequest, opts ...config.CallOpt) (*govv1.QueryTallyResultResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.TallyResult, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling GovGRPCClient.TallyResult")
	}
	return resp, nil
}

func (c *GovGRPCClient) Vote(ctx context.Context, req *govv1.QueryVoteRequest, opts ...config.CallOpt) (*govv1.QueryVoteResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Vote, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling GovGRPCClient.Vote")
	}
	return resp, nil
}

func (c *GovGRPCClient) Votes(ctx context.Context, req *govv1.QueryVotesRequest, opts ...config.CallOpt) (*govv1.QueryVotesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Votes, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling GovGRPCClient.Votes")
	}
	return resp, nil
}
