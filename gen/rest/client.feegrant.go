package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	feegrant "cosmossdk.io/x/feegrant"
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

func (c *FeegrantRESTClient) Allowance(ctx context.Context, req *feegrant.QueryAllowanceRequest) (*feegrant.QueryAllowanceResponse, error) {
	resp := &feegrant.QueryAllowanceResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/feegrant/v1beta1/allowance/{granter}/{grantee}",
		[]string{"granter", "grantee"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling FeegrantRESTClient.Allowance")
}

func (c *FeegrantRESTClient) Allowances(ctx context.Context, req *feegrant.QueryAllowancesRequest) (*feegrant.QueryAllowancesResponse, error) {
	resp := &feegrant.QueryAllowancesResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/feegrant/v1beta1/allowances/{grantee}",
		[]string{"grantee"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling FeegrantRESTClient.Allowances")
}

func (c *FeegrantRESTClient) AllowancesByGranter(ctx context.Context, req *feegrant.QueryAllowancesByGranterRequest) (*feegrant.QueryAllowancesByGranterResponse, error) {
	resp := &feegrant.QueryAllowancesByGranterResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/feegrant/v1beta1/issued/{granter}",
		[]string{"granter"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling FeegrantRESTClient.AllowancesByGranter")
}
