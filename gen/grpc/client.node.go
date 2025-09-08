package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/allora-network/allora-sdk-go/config"

	node "github.com/cosmos/cosmos-sdk/client/grpc/node"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// NodeGRPCClient provides gRPC access to the node module
type NodeGRPCClient struct {
	client node.ServiceClient
	logger zerolog.Logger
}

var _ interfaces.NodeClient = (*NodeGRPCClient)(nil)

// NewNodeGRPCClient creates a new node REST client
func NewNodeGRPCClient(conn *grpc.ClientConn, logger zerolog.Logger) *NodeGRPCClient {
	return &NodeGRPCClient{
		client: node.NewServiceClient(conn),
		logger: logger.With().Str("module", "node").Str("protocol", "grpc").Logger(),
	}
}

func (c *NodeGRPCClient) Config(ctx context.Context, req *node.ConfigRequest, opts ...config.CallOpt) (*node.ConfigResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Config, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling NodeGRPCClient.Config")
	}
	return resp, nil
}

func (c *NodeGRPCClient) Status(ctx context.Context, req *node.StatusRequest, opts ...config.CallOpt) (*node.StatusResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Status, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling NodeGRPCClient.Status")
	}
	return resp, nil
}
