package interfaces

import (
	"context"

	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

type GovClient interface {
	Constitution(ctx context.Context, req *govv1.QueryConstitutionRequest) (*govv1.QueryConstitutionResponse, error)
	Proposal(ctx context.Context, req *govv1.QueryProposalRequest) (*govv1.QueryProposalResponse, error)
	Proposals(ctx context.Context, req *govv1.QueryProposalsRequest) (*govv1.QueryProposalsResponse, error)
	Vote(ctx context.Context, req *govv1.QueryVoteRequest) (*govv1.QueryVoteResponse, error)
	Votes(ctx context.Context, req *govv1.QueryVotesRequest) (*govv1.QueryVotesResponse, error)
	Params(ctx context.Context, req *govv1.QueryParamsRequest) (*govv1.QueryParamsResponse, error)
	Deposit(ctx context.Context, req *govv1.QueryDepositRequest) (*govv1.QueryDepositResponse, error)
	Deposits(ctx context.Context, req *govv1.QueryDepositsRequest) (*govv1.QueryDepositsResponse, error)
	TallyResult(ctx context.Context, req *govv1.QueryTallyResultRequest) (*govv1.QueryTallyResultResponse, error)
}
