package wrapper

import (
	"context"

	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/pool"
)

// GovClientWrapper wraps the gov module with pool management and retry logic
type GovClientWrapper struct {
	poolManager *pool.ClientPoolManager
	logger      zerolog.Logger
}

// NewGovClientWrapper creates a new gov client wrapper
func NewGovClientWrapper(poolManager *pool.ClientPoolManager, logger zerolog.Logger) *GovClientWrapper {
	return &GovClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "gov").Logger(),
	}
}

func (c *GovClientWrapper) Constitution(ctx context.Context, req *govv1.QueryConstitutionRequest) (*govv1.QueryConstitutionResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*govv1.QueryConstitutionResponse, error) {
		return client.Gov().Constitution(ctx, req)
	})
}

func (c *GovClientWrapper) Proposal(ctx context.Context, req *govv1.QueryProposalRequest) (*govv1.QueryProposalResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*govv1.QueryProposalResponse, error) {
		return client.Gov().Proposal(ctx, req)
	})
}

func (c *GovClientWrapper) Proposals(ctx context.Context, req *govv1.QueryProposalsRequest) (*govv1.QueryProposalsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*govv1.QueryProposalsResponse, error) {
		return client.Gov().Proposals(ctx, req)
	})
}

func (c *GovClientWrapper) Vote(ctx context.Context, req *govv1.QueryVoteRequest) (*govv1.QueryVoteResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*govv1.QueryVoteResponse, error) {
		return client.Gov().Vote(ctx, req)
	})
}

func (c *GovClientWrapper) Votes(ctx context.Context, req *govv1.QueryVotesRequest) (*govv1.QueryVotesResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*govv1.QueryVotesResponse, error) {
		return client.Gov().Votes(ctx, req)
	})
}

func (c *GovClientWrapper) Params(ctx context.Context, req *govv1.QueryParamsRequest) (*govv1.QueryParamsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*govv1.QueryParamsResponse, error) {
		return client.Gov().Params(ctx, req)
	})
}

func (c *GovClientWrapper) Deposit(ctx context.Context, req *govv1.QueryDepositRequest) (*govv1.QueryDepositResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*govv1.QueryDepositResponse, error) {
		return client.Gov().Deposit(ctx, req)
	})
}

func (c *GovClientWrapper) Deposits(ctx context.Context, req *govv1.QueryDepositsRequest) (*govv1.QueryDepositsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*govv1.QueryDepositsResponse, error) {
		return client.Gov().Deposits(ctx, req)
	})
}

func (c *GovClientWrapper) TallyResult(ctx context.Context, req *govv1.QueryTallyResultRequest) (*govv1.QueryTallyResultResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*govv1.QueryTallyResultResponse, error) {
		return client.Gov().TallyResult(ctx, req)
	})
}
