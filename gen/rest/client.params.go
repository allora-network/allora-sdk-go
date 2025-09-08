package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	proposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
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

func (c *ParamsRESTClient) Params(ctx context.Context, req *proposal.QueryParamsRequest) (*proposal.QueryParamsResponse, error) {
	resp := &proposal.QueryParamsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/params/v1beta1/params",
		nil, []string{"subspace", "key"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling ParamsRESTClient.Params")
}

func (c *ParamsRESTClient) Subspaces(ctx context.Context, req *proposal.QuerySubspacesRequest) (*proposal.QuerySubspacesResponse, error) {
	resp := &proposal.QuerySubspacesResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/params/v1beta1/subspaces",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling ParamsRESTClient.Subspaces")
}
