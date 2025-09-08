package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/allora-network/allora-sdk-go/config"

	authz "github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// AuthzGRPCClient provides gRPC access to the authz module
type AuthzGRPCClient struct {
	client authz.QueryClient
	logger zerolog.Logger
}

var _ interfaces.AuthzClient = (*AuthzGRPCClient)(nil)

// NewAuthzGRPCClient creates a new authz REST client
func NewAuthzGRPCClient(conn *grpc.ClientConn, logger zerolog.Logger) *AuthzGRPCClient {
	return &AuthzGRPCClient{
		client: authz.NewQueryClient(conn),
		logger: logger.With().Str("module", "authz").Str("protocol", "grpc").Logger(),
	}
}

func (c *AuthzGRPCClient) Grants(ctx context.Context, req *authz.QueryGrantsRequest, opts ...config.CallOpt) (*authz.QueryGrantsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Grants, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling AuthzGRPCClient.Grants")
	}
	return resp, nil
}

func (c *AuthzGRPCClient) GranterGrants(ctx context.Context, req *authz.QueryGranterGrantsRequest, opts ...config.CallOpt) (*authz.QueryGranterGrantsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GranterGrants, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling AuthzGRPCClient.GranterGrants")
	}
	return resp, nil
}

func (c *AuthzGRPCClient) GranteeGrants(ctx context.Context, req *authz.QueryGranteeGrantsRequest, opts ...config.CallOpt) (*authz.QueryGranteeGrantsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GranteeGrants, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling AuthzGRPCClient.GranteeGrants")
	}
	return resp, nil
}
