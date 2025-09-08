package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

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

func (c *NodeGRPCClient) Config(ctx context.Context, req *node.ConfigRequest) (*node.ConfigResponse, error) {
	resp, err := c.client.Config(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Config")
}

func (c *NodeGRPCClient) Status(ctx context.Context, req *node.StatusRequest) (*node.StatusResponse, error) {
	resp, err := c.client.Status(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Status")
}
