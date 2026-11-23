package cosmosrpc

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	emissionstypes "github.com/allora-network/allora-chain/x/emissions/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/allora-network/allora-sdk-go/config"
	sdkgrpc "github.com/allora-network/allora-sdk-go/gen/grpc"
	"github.com/allora-network/allora-sdk-go/gen/rest"
)

// notImplementedRESTEndpoint starts a REST endpoint that answers every request
// the way a node answers a query method it does not serve: HTTP 501 with the
// grpc-gateway error body. It returns the endpoint URL and a counter of the
// requests it received.
func notImplementedRESTEndpoint(t *testing.T) (string, *atomic.Int64) {
	t.Helper()

	var hits atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/emissions/") {
			hits.Add(1)
		}
		w.WriteHeader(http.StatusNotImplemented)
		_, _ = w.Write([]byte(`{"code":12,"message":"Not Implemented"}`))
	}))
	t.Cleanup(server.Close)

	return server.URL, &hits
}

// notImplementedGRPCEndpoint starts a gRPC endpoint that answers every method
// with codes.Unimplemented, which is what a node returns for a method its
// registered service does not implement. It returns the endpoint URL and a
// counter of the calls it received.
func notImplementedGRPCEndpoint(t *testing.T) (string, *atomic.Int64) {
	t.Helper()

	var hits atomic.Int64
	server := grpcpkg.NewServer(grpcpkg.UnknownServiceHandler(
		func(any, grpcpkg.ServerStream) error {
			hits.Add(1)
			return status.Error(codes.Unimplemented, "unknown service")
		}))

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	go func() { _ = server.Serve(lis) }()
	t.Cleanup(server.Stop)

	return "grpc://" + lis.Addr().String(), &hits
}

// TestPoolTreatsRESTNotImplementedAsTerminal pins that a REST endpoint which
// answers 501 is not a broken endpoint. The node is up and has told the caller
// that it does not serve this method, so the pool must return that answer at
// once: no health penalty, no backoff, no rotation through the other endpoints,
// and an error whose code is Unimplemented rather than Unknown.
func TestPoolTreatsRESTNotImplementedAsTerminal(t *testing.T) {
	logger := zerolog.Nop()
	urlA, hitsA := notImplementedRESTEndpoint(t)
	urlB, hitsB := notImplementedRESTEndpoint(t)

	p := NewClientPool([]Client{
		rest.NewRESTClient(urlA, logger),
		rest.NewRESTClient(urlB, logger),
	}, logger, 0)
	t.Cleanup(func() { _ = p.Close() })

	_, err := p.Emissions().IsTopicWorkerWhitelistEnabled(context.Background(),
		&emissionstypes.IsTopicWorkerWhitelistEnabledRequest{TopicId: 1})

	require.Error(t, err)
	require.Equal(t, codes.Unimplemented, status.Code(err),
		"a 501 must reach the caller as Unimplemented, got %v", err)
	require.Contains(t, err.Error(), "IsTopicWorkerWhitelistEnabled",
		"the error must name the method that is not implemented")

	require.Equal(t, int64(1), hitsA.Load()+hitsB.Load(),
		"an unimplemented method must not be retried on the other endpoints")

	hitURL := urlA
	if hitsB.Load() > 0 {
		hitURL = urlB
	}
	require.Contains(t, err.Error(), hitURL,
		"the error must name the endpoint that does not implement the method")

	health := p.GetHealthStatus()
	require.Equal(t, 2, health["active_clients"],
		"a 501 says nothing about endpoint health; both endpoints must stay active")
	require.Equal(t, 0, health["cooling_clients"],
		"a 501 must not cool a healthy endpoint")
}

// TestPoolTreatsGRPCUnimplementedAsTerminal is the gRPC half of the same
// contract: codes.Unimplemented is a definitive answer from a healthy node.
func TestPoolTreatsGRPCUnimplementedAsTerminal(t *testing.T) {
	logger := zerolog.Nop()
	urlA, hitsA := notImplementedGRPCEndpoint(t)
	urlB, hitsB := notImplementedGRPCEndpoint(t)

	clientA, err := sdkgrpc.NewGRPCClient(config.EndpointConfig{URL: urlA, Protocol: config.ProtocolGRPC}, logger)
	require.NoError(t, err)
	clientB, err := sdkgrpc.NewGRPCClient(config.EndpointConfig{URL: urlB, Protocol: config.ProtocolGRPC}, logger)
	require.NoError(t, err)

	p := NewClientPool([]Client{clientA, clientB}, logger, 0)
	t.Cleanup(func() { _ = p.Close() })

	_, err = p.Emissions().IsTopicWorkerWhitelistEnabled(context.Background(),
		&emissionstypes.IsTopicWorkerWhitelistEnabledRequest{TopicId: 1})

	require.Error(t, err)
	require.Equal(t, codes.Unimplemented, status.Code(err),
		"Unimplemented must reach the caller unchanged, got %v", err)
	require.Contains(t, err.Error(), "IsTopicWorkerWhitelistEnabled")

	require.Equal(t, int64(1), hitsA.Load()+hitsB.Load(),
		"an unimplemented method must not be retried on the other endpoints")

	hitURL := urlA
	if hitsB.Load() > 0 {
		hitURL = urlB
	}
	require.Contains(t, err.Error(), hitURL)

	health := p.GetHealthStatus()
	require.Equal(t, 2, health["active_clients"])
	require.Equal(t, 0, health["cooling_clients"])
}

// TestPoolStillFailsOverOnServerError guards the other side of the change: a
// 500 is a real endpoint failure, so the pool must keep penalising the endpoint
// and trying the others.
func TestPoolStillFailsOverOnServerError(t *testing.T) {
	logger := zerolog.Nop()

	var hits atomic.Int64
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/emissions/") {
			hits.Add(1)
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"code":13,"message":"Internal"}`))
	})
	serverA := httptest.NewServer(handler)
	t.Cleanup(serverA.Close)
	serverB := httptest.NewServer(handler)
	t.Cleanup(serverB.Close)

	p := NewClientPool([]Client{
		rest.NewRESTClient(serverA.URL, logger),
		rest.NewRESTClient(serverB.URL, logger),
	}, logger, 0)
	t.Cleanup(func() { _ = p.Close() })

	_, err := p.Emissions().IsTopicWorkerWhitelistEnabled(context.Background(),
		&emissionstypes.IsTopicWorkerWhitelistEnabledRequest{TopicId: 1})

	require.Error(t, err)
	require.NotEqual(t, codes.Unimplemented, status.Code(err),
		"a 500 is not an unimplemented method, got %v", err)
	require.Equal(t, int64(2), hits.Load(),
		"a 500 must still be retried on the other endpoint")

	health := p.GetHealthStatus()
	require.Equal(t, 0, health["active_clients"],
		"a 500 must still cool the failing endpoints")
	require.Equal(t, 2, health["cooling_clients"])
}
