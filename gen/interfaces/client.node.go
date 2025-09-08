package interfaces

import (
	"context"

	node "github.com/cosmos/cosmos-sdk/client/grpc/node"

	"github.com/allora-network/allora-sdk-go/config"
)

type NodeClient interface {
	Config(ctx context.Context, req *node.ConfigRequest, opts ...config.CallOpt) (*node.ConfigResponse, error)
	Status(ctx context.Context, req *node.StatusRequest, opts ...config.CallOpt) (*node.StatusResponse, error)
}
