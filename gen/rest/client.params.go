package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	proposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"

	"github.com/allora-network/allora-sdk-go/config"
)

// ParamsRESTClient provides REST access to the params module
type ParamsRESTClient struct {
	*RESTClientCore
	logger zerolog.Logger
}

// NewParamsRESTClient creates a new params REST client
func NewParamsRESTClient(core *RESTClientCore, logger zerolog.Logger) *ParamsRESTClient {
	return &ParamsRESTClient{
		RESTClientCore: core,
		logger:         logger.With().Str("module", "params").Str("protocol", "rest").Logger(),
	}
}

func (c *ParamsRESTClient) Params(ctx context.Context, req *proposal.QueryParamsRequest, opts ...config.CallOpt) (*proposal.QueryParamsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &proposal.QueryParamsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/params/v1beta1/params",
		nil, []string{"subspace", "key"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling ParamsRESTClient.Params")
	}
	return resp, nil
}

func (c *ParamsRESTClient) Subspaces(ctx context.Context, req *proposal.QuerySubspacesRequest, opts ...config.CallOpt) (*proposal.QuerySubspacesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &proposal.QuerySubspacesResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/params/v1beta1/subspaces",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling ParamsRESTClient.Subspaces")
	}
	return resp, nil
}
