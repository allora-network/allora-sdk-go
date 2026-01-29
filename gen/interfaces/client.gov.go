package interfaces

import (
	"context"

	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/allora-network/allora-sdk-go/config"
)

type GovClient interface {
	Constitution(ctx context.Context, req *govv1.QueryConstitutionRequest, opts ...config.CallOpt) (*govv1.QueryConstitutionResponse, error)
	Deposit(ctx context.Context, req *govv1.QueryDepositRequest, opts ...config.CallOpt) (*govv1.QueryDepositResponse, error)
	Deposits(ctx context.Context, req *govv1.QueryDepositsRequest, opts ...config.CallOpt) (*govv1.QueryDepositsResponse, error)
	Params(ctx context.Context, req *govv1.QueryParamsRequest, opts ...config.CallOpt) (*govv1.QueryParamsResponse, error)
	Proposal(ctx context.Context, req *govv1.QueryProposalRequest, opts ...config.CallOpt) (*govv1.QueryProposalResponse, error)
	Proposals(ctx context.Context, req *govv1.QueryProposalsRequest, opts ...config.CallOpt) (*govv1.QueryProposalsResponse, error)
	TallyResult(ctx context.Context, req *govv1.QueryTallyResultRequest, opts ...config.CallOpt) (*govv1.QueryTallyResultResponse, error)
	Vote(ctx context.Context, req *govv1.QueryVoteRequest, opts ...config.CallOpt) (*govv1.QueryVoteResponse, error)
	Votes(ctx context.Context, req *govv1.QueryVotesRequest, opts ...config.CallOpt) (*govv1.QueryVotesResponse, error)
}
