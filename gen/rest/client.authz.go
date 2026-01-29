package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	authz "github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/allora-network/allora-sdk-go/config"
)

// AuthzRESTClient provides REST access to the authz module
type AuthzRESTClient struct {
	*RESTClientCore
	logger zerolog.Logger
}

// NewAuthzRESTClient creates a new authz REST client
func NewAuthzRESTClient(core *RESTClientCore, logger zerolog.Logger) *AuthzRESTClient {
	return &AuthzRESTClient{
		RESTClientCore: core,
		logger:         logger.With().Str("module", "authz").Str("protocol", "rest").Logger(),
	}
}

func (c *AuthzRESTClient) GranteeGrants(ctx context.Context, req *authz.QueryGranteeGrantsRequest, opts ...config.CallOpt) (*authz.QueryGranteeGrantsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &authz.QueryGranteeGrantsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/authz/v1beta1/grants/grantee/{grantee}",
		[]string{"grantee"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling AuthzRESTClient.GranteeGrants")
	}
	return resp, nil
}

func (c *AuthzRESTClient) GranterGrants(ctx context.Context, req *authz.QueryGranterGrantsRequest, opts ...config.CallOpt) (*authz.QueryGranterGrantsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &authz.QueryGranterGrantsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/authz/v1beta1/grants/granter/{granter}",
		[]string{"granter"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling AuthzRESTClient.GranterGrants")
	}
	return resp, nil
}

func (c *AuthzRESTClient) Grants(ctx context.Context, req *authz.QueryGrantsRequest, opts ...config.CallOpt) (*authz.QueryGrantsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &authz.QueryGrantsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/authz/v1beta1/grants",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse", "granter", "grantee", "msg_type_url"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling AuthzRESTClient.Grants")
	}
	return resp, nil
}
