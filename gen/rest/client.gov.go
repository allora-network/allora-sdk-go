package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

// GovRESTClient provides REST access to the gov module
type GovRESTClient struct {
	*RESTClientCore
	logger zerolog.Logger
}

// NewGovRESTClient creates a new gov REST client
func NewGovRESTClient(core *RESTClientCore, logger zerolog.Logger) *GovRESTClient {
	return &GovRESTClient{
		RESTClientCore: core,
		logger:         logger.With().Str("module", "gov").Str("protocol", "rest").Logger(),
	}
}

func (c *GovRESTClient) Constitution(ctx context.Context, req *govv1.QueryConstitutionRequest) (*govv1.QueryConstitutionResponse, error) {
	resp := &govv1.QueryConstitutionResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/gov/v1/constitution",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling GovRESTClient.Constitution")
}

func (c *GovRESTClient) Proposal(ctx context.Context, req *govv1.QueryProposalRequest) (*govv1.QueryProposalResponse, error) {
	resp := &govv1.QueryProposalResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/gov/v1/proposals/{proposal_id}",
		[]string{"proposal_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling GovRESTClient.Proposal")
}

func (c *GovRESTClient) Proposals(ctx context.Context, req *govv1.QueryProposalsRequest) (*govv1.QueryProposalsResponse, error) {
	resp := &govv1.QueryProposalsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/gov/v1/proposals",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse", "proposal_status", "voter", "depositor"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling GovRESTClient.Proposals")
}

func (c *GovRESTClient) Vote(ctx context.Context, req *govv1.QueryVoteRequest) (*govv1.QueryVoteResponse, error) {
	resp := &govv1.QueryVoteResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/gov/v1/proposals/{proposal_id}/votes/{voter}",
		[]string{"proposal_id", "voter"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling GovRESTClient.Vote")
}

func (c *GovRESTClient) Votes(ctx context.Context, req *govv1.QueryVotesRequest) (*govv1.QueryVotesResponse, error) {
	resp := &govv1.QueryVotesResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/gov/v1/proposals/{proposal_id}/votes",
		[]string{"proposal_id"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling GovRESTClient.Votes")
}

func (c *GovRESTClient) Params(ctx context.Context, req *govv1.QueryParamsRequest) (*govv1.QueryParamsResponse, error) {
	resp := &govv1.QueryParamsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/gov/v1/params/{params_type}",
		[]string{"params_type"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling GovRESTClient.Params")
}

func (c *GovRESTClient) Deposit(ctx context.Context, req *govv1.QueryDepositRequest) (*govv1.QueryDepositResponse, error) {
	resp := &govv1.QueryDepositResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/gov/v1/proposals/{proposal_id}/deposits/{depositor}",
		[]string{"proposal_id", "depositor"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling GovRESTClient.Deposit")
}

func (c *GovRESTClient) Deposits(ctx context.Context, req *govv1.QueryDepositsRequest) (*govv1.QueryDepositsResponse, error) {
	resp := &govv1.QueryDepositsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/gov/v1/proposals/{proposal_id}/deposits",
		[]string{"proposal_id"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling GovRESTClient.Deposits")
}

func (c *GovRESTClient) TallyResult(ctx context.Context, req *govv1.QueryTallyResultRequest) (*govv1.QueryTallyResultResponse, error) {
	resp := &govv1.QueryTallyResultResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/gov/v1/proposals/{proposal_id}/tally",
		[]string{"proposal_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling GovRESTClient.TallyResult")
}
