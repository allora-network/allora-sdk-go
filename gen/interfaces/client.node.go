package interfaces

import (
	"context"

	node "github.com/cosmos/cosmos-sdk/client/grpc/node"
)

type NodeClient interface {
	Config(ctx context.Context, req *node.ConfigRequest) (*node.ConfigResponse, error)
	Status(ctx context.Context, req *node.StatusRequest) (*node.StatusResponse, error)
}
