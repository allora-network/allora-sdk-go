package wrapper

import (
	"context"

	node "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/pool"
)

// NodeClientWrapper wraps the node module with pool management and retry logic
type NodeClientWrapper struct {
	poolManager *pool.ClientPoolManager
	logger      zerolog.Logger
}

// NewNodeClientWrapper creates a new node client wrapper
func NewNodeClientWrapper(poolManager *pool.ClientPoolManager, logger zerolog.Logger) *NodeClientWrapper {
	return &NodeClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "node").Logger(),
	}
}

func (c *NodeClientWrapper) Config(ctx context.Context, req *node.ConfigRequest, opts ...config.CallOpt) (*node.ConfigResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*node.ConfigResponse, error) {
		return client.Node().Config(ctx, req, opts...)
	})
}

func (c *NodeClientWrapper) Status(ctx context.Context, req *node.StatusRequest, opts ...config.CallOpt) (*node.StatusResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*node.StatusResponse, error) {
		return client.Node().Status(ctx, req, opts...)
	})
}
