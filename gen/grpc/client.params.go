package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	proposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"

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

func (c *ParamsGRPCClient) Params(ctx context.Context, req *proposal.QueryParamsRequest) (*proposal.QueryParamsResponse, error) {
	resp, err := c.client.Params(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Params")
}

func (c *ParamsGRPCClient) Subspaces(ctx context.Context, req *proposal.QuerySubspacesRequest) (*proposal.QuerySubspacesResponse, error) {
	resp, err := c.client.Subspaces(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Subspaces")
}
