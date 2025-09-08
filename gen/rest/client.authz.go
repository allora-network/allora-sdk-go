package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	authz "github.com/cosmos/cosmos-sdk/x/authz"
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

func (c *AuthzRESTClient) Grants(ctx context.Context, req *authz.QueryGrantsRequest) (*authz.QueryGrantsResponse, error) {
	resp := &authz.QueryGrantsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/authz/v1beta1/grants",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse", "granter", "grantee", "msg_type_url"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling AuthzRESTClient.Grants")
}

func (c *AuthzRESTClient) GranterGrants(ctx context.Context, req *authz.QueryGranterGrantsRequest) (*authz.QueryGranterGrantsResponse, error) {
	resp := &authz.QueryGranterGrantsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/authz/v1beta1/grants/granter/{granter}",
		[]string{"granter"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling AuthzRESTClient.GranterGrants")
}

func (c *AuthzRESTClient) GranteeGrants(ctx context.Context, req *authz.QueryGranteeGrantsRequest) (*authz.QueryGranteeGrantsResponse, error) {
	resp := &authz.QueryGranteeGrantsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/authz/v1beta1/grants/grantee/{grantee}",
		[]string{"grantee"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling AuthzRESTClient.GranteeGrants")
}
