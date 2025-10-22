package cosmosrpc

import (
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/gen/wrapper"
	"github.com/allora-network/allora-sdk-go/pool"
	"github.com/rs/zerolog"
)

type Client = interfaces.CosmosClient

type ClientPool interface {
	interfaces.CosmosClientPool
	GetHealthStatus() map[string]any
}

type clientPool struct {
	*wrapper.WrapperClient
	poolManager *pool.ClientPoolManager[Client]
}

var _ ClientPool = (*clientPool)(nil)

func NewClientPoolManager(clients []Client, logger zerolog.Logger) *clientPool {
	mgr := pool.NewClientPoolManager(clients, logger)
	return &clientPool{
		WrapperClient: wrapper.NewWrapperClient(mgr, logger),
		poolManager:   pool.NewClientPoolManager(clients, logger),
	}
}

func (p *clientPool) Close() error {
	p.poolManager.Close()
	return nil
}

func (p *clientPool) GetHealthStatus() map[string]any {
	return p.poolManager.GetHealthStatus()
}
