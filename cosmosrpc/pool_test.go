package cosmosrpc

import (
	"context"
	"errors"
	"testing"

	node "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// fakeNodeClient fails Config calls so the pool records a health failure.
type fakeNodeClient struct {
	interfaces.NodeClient
	err error
}

func (f *fakeNodeClient) Config(ctx context.Context, req *node.ConfigRequest, opts ...config.CallOpt) (*node.ConfigResponse, error) {
	return nil, f.err
}

// fakeCosmosClient is a minimal CosmosClient whose Node() can be forced to fail.
type fakeCosmosClient struct {
	interfaces.CosmosClient
	url  string
	node interfaces.NodeClient
}

func (f *fakeCosmosClient) Close() error                          { return nil }
func (f *fakeCosmosClient) GetEndpointURL() string                { return f.url }
func (f *fakeCosmosClient) GetProtocol() config.Protocol          { return config.ProtocolGRPC }
func (f *fakeCosmosClient) HealthCheck(ctx context.Context) error { return nil }
func (f *fakeCosmosClient) Node() interfaces.NodeClient           { return f.node }

// TestNewClientPoolSharesManager pins the invariant that the ClientPoolManager
// executing requests is the SAME one reported by GetHealthStatus. Regression:
// NewClientPool built two managers, so a failure observed on the executing
// manager never showed up in GetHealthStatus (dashboards lied).
func TestNewClientPoolSharesManager(t *testing.T) {
	logger := zerolog.Nop()
	failNode := &fakeNodeClient{err: errors.New("boom")}
	clients := []Client{
		&fakeCosmosClient{url: "grpc://a:9090", node: failNode},
		&fakeCosmosClient{url: "grpc://b:9090", node: failNode},
	}

	p := NewClientPool(clients, logger)
	defer func() { require.NoError(t, p.Close()) }()

	// Drive the EXECUTING path: a failing Config call routes through the
	// wrapper's pool manager, which reports the failure and cools the client.
	for i := 0; i < 3; i++ {
		_, _ = p.Node().Config(context.Background(), &node.ConfigRequest{})
	}

	status := p.GetHealthStatus()
	require.Greater(t, status["cooling_clients"], 0,
		"GetHealthStatus must observe the same manager that executes requests; "+
			"a split-brain pool reports 0 cooling clients here")
}
