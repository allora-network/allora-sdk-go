package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	feegrant "cosmossdk.io/x/feegrant"

	"github.com/allora-network/allora-sdk-go/config"
)

// FeegrantRESTClient provides REST access to the feegrant module
type FeegrantRESTClient struct {
	*RESTClientCore
	logger zerolog.Logger
}

// NewFeegrantRESTClient creates a new feegrant REST client
func NewFeegrantRESTClient(core *RESTClientCore, logger zerolog.Logger) *FeegrantRESTClient {
	return &FeegrantRESTClient{
		RESTClientCore: core,
		logger:         logger.With().Str("module", "feegrant").Str("protocol", "rest").Logger(),
	}
}

func (c *FeegrantRESTClient) Allowance(ctx context.Context, req *feegrant.QueryAllowanceRequest, opts ...config.CallOpt) (*feegrant.QueryAllowanceResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &feegrant.QueryAllowanceResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/feegrant/v1beta1/allowance/{granter}/{grantee}",
		[]string{"granter", "grantee"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling FeegrantRESTClient.Allowance")
	}
	return resp, nil
}

func (c *FeegrantRESTClient) Allowances(ctx context.Context, req *feegrant.QueryAllowancesRequest, opts ...config.CallOpt) (*feegrant.QueryAllowancesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &feegrant.QueryAllowancesResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/feegrant/v1beta1/allowances/{grantee}",
		[]string{"grantee"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling FeegrantRESTClient.Allowances")
	}
	return resp, nil
}

func (c *FeegrantRESTClient) AllowancesByGranter(ctx context.Context, req *feegrant.QueryAllowancesByGranterRequest, opts ...config.CallOpt) (*feegrant.QueryAllowancesByGranterResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &feegrant.QueryAllowancesByGranterResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/feegrant/v1beta1/issued/{granter}",
		[]string{"granter"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling FeegrantRESTClient.AllowancesByGranter")
	}
	return resp, nil
}
