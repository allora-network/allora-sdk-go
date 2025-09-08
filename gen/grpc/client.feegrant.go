package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	feegrant "cosmossdk.io/x/feegrant"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// FeegrantGRPCClient provides gRPC access to the feegrant module
type FeegrantGRPCClient struct {
	client feegrant.QueryClient
	logger zerolog.Logger
}

var _ interfaces.FeegrantClient = (*FeegrantGRPCClient)(nil)

// NewFeegrantGRPCClient creates a new feegrant REST client
func NewFeegrantGRPCClient(conn *grpc.ClientConn, logger zerolog.Logger) *FeegrantGRPCClient {
	return &FeegrantGRPCClient{
		client: feegrant.NewQueryClient(conn),
		logger: logger.With().Str("module", "feegrant").Str("protocol", "grpc").Logger(),
	}
}

func (c *FeegrantGRPCClient) Allowance(ctx context.Context, req *feegrant.QueryAllowanceRequest) (*feegrant.QueryAllowanceResponse, error) {
	resp, err := c.client.Allowance(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Allowance")
}

func (c *FeegrantGRPCClient) Allowances(ctx context.Context, req *feegrant.QueryAllowancesRequest) (*feegrant.QueryAllowancesResponse, error) {
	resp, err := c.client.Allowances(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Allowances")
}

func (c *FeegrantGRPCClient) AllowancesByGranter(ctx context.Context, req *feegrant.QueryAllowancesByGranterRequest) (*feegrant.QueryAllowancesByGranterResponse, error) {
	resp, err := c.client.AllowancesByGranter(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.AllowancesByGranter")
}
