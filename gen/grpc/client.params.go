package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	proposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// ParamsGRPCClient provides gRPC access to the params module
type ParamsGRPCClient struct {
	client proposal.QueryClient
	logger zerolog.Logger
}

var _ interfaces.ParamsClient = (*ParamsGRPCClient)(nil)

// NewParamsGRPCClient creates a new params REST client
func NewParamsGRPCClient(conn *grpc.ClientConn, logger zerolog.Logger) *ParamsGRPCClient {
	return &ParamsGRPCClient{
		client: proposal.NewQueryClient(conn),
		logger: logger.With().Str("module", "params").Str("protocol", "grpc").Logger(),
	}
}

func (c *ParamsGRPCClient) Params(ctx context.Context, req *proposal.QueryParamsRequest, opts ...config.CallOpt) (*proposal.QueryParamsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Params, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling ParamsGRPCClient.Params")
	}
	return resp, nil
}

func (c *ParamsGRPCClient) Subspaces(ctx context.Context, req *proposal.QuerySubspacesRequest, opts ...config.CallOpt) (*proposal.QuerySubspacesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Subspaces, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling ParamsGRPCClient.Subspaces")
	}
	return resp, nil
}
