package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/allora-network/allora-sdk-go/config"
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

func (c *GovRESTClient) Constitution(ctx context.Context, req *govv1.QueryConstitutionRequest, opts ...config.CallOpt) (*govv1.QueryConstitutionResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &govv1.QueryConstitutionResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/gov/v1/constitution",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling GovRESTClient.Constitution")
	}
	return resp, nil
}

func (c *GovRESTClient) Deposit(ctx context.Context, req *govv1.QueryDepositRequest, opts ...config.CallOpt) (*govv1.QueryDepositResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &govv1.QueryDepositResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/gov/v1/proposals/{proposal_id}/deposits/{depositor}",
		[]string{"proposal_id", "depositor"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling GovRESTClient.Deposit")
	}
	return resp, nil
}

func (c *GovRESTClient) Deposits(ctx context.Context, req *govv1.QueryDepositsRequest, opts ...config.CallOpt) (*govv1.QueryDepositsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &govv1.QueryDepositsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/gov/v1/proposals/{proposal_id}/deposits",
		[]string{"proposal_id"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling GovRESTClient.Deposits")
	}
	return resp, nil
}

func (c *GovRESTClient) Params(ctx context.Context, req *govv1.QueryParamsRequest, opts ...config.CallOpt) (*govv1.QueryParamsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &govv1.QueryParamsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/gov/v1/params/{params_type}",
		[]string{"params_type"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling GovRESTClient.Params")
	}
	return resp, nil
}

func (c *GovRESTClient) Proposal(ctx context.Context, req *govv1.QueryProposalRequest, opts ...config.CallOpt) (*govv1.QueryProposalResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &govv1.QueryProposalResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/gov/v1/proposals/{proposal_id}",
		[]string{"proposal_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling GovRESTClient.Proposal")
	}
	return resp, nil
}

func (c *GovRESTClient) Proposals(ctx context.Context, req *govv1.QueryProposalsRequest, opts ...config.CallOpt) (*govv1.QueryProposalsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &govv1.QueryProposalsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/gov/v1/proposals",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse", "proposal_status", "voter", "depositor"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling GovRESTClient.Proposals")
	}
	return resp, nil
}

func (c *GovRESTClient) TallyResult(ctx context.Context, req *govv1.QueryTallyResultRequest, opts ...config.CallOpt) (*govv1.QueryTallyResultResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &govv1.QueryTallyResultResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/gov/v1/proposals/{proposal_id}/tally",
		[]string{"proposal_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling GovRESTClient.TallyResult")
	}
	return resp, nil
}

func (c *GovRESTClient) Vote(ctx context.Context, req *govv1.QueryVoteRequest, opts ...config.CallOpt) (*govv1.QueryVoteResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &govv1.QueryVoteResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/gov/v1/proposals/{proposal_id}/votes/{voter}",
		[]string{"proposal_id", "voter"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling GovRESTClient.Vote")
	}
	return resp, nil
}

func (c *GovRESTClient) Votes(ctx context.Context, req *govv1.QueryVotesRequest, opts ...config.CallOpt) (*govv1.QueryVotesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &govv1.QueryVotesResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/gov/v1/proposals/{proposal_id}/votes",
		[]string{"proposal_id"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling GovRESTClient.Votes")
	}
	return resp, nil
}
