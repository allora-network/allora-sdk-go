package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

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

func (c *AuthzGRPCClient) Grants(ctx context.Context, req *authz.QueryGrantsRequest) (*authz.QueryGrantsResponse, error) {
	resp, err := c.client.Grants(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Grants")
}

func (c *AuthzGRPCClient) GranterGrants(ctx context.Context, req *authz.QueryGranterGrantsRequest) (*authz.QueryGranterGrantsResponse, error) {
	resp, err := c.client.GranterGrants(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GranterGrants")
}

func (c *AuthzGRPCClient) GranteeGrants(ctx context.Context, req *authz.QueryGranteeGrantsRequest) (*authz.QueryGranteeGrantsResponse, error) {
	resp, err := c.client.GranteeGrants(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GranteeGrants")
}
