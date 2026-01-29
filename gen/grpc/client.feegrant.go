package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	feegrant "cosmossdk.io/x/feegrant"

	"github.com/allora-network/allora-sdk-go/config"
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

func (c *FeegrantGRPCClient) Allowance(ctx context.Context, req *feegrant.QueryAllowanceRequest, opts ...config.CallOpt) (*feegrant.QueryAllowanceResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Allowance, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling FeegrantGRPCClient.Allowance")
	}
	return resp, nil
}

func (c *FeegrantGRPCClient) Allowances(ctx context.Context, req *feegrant.QueryAllowancesRequest, opts ...config.CallOpt) (*feegrant.QueryAllowancesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Allowances, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling FeegrantGRPCClient.Allowances")
	}
	return resp, nil
}

func (c *FeegrantGRPCClient) AllowancesByGranter(ctx context.Context, req *feegrant.QueryAllowancesByGranterRequest, opts ...config.CallOpt) (*feegrant.QueryAllowancesByGranterResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.AllowancesByGranter, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling FeegrantGRPCClient.AllowancesByGranter")
	}
	return resp, nil
}
