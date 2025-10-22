package tmrpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	tmhttp "github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/pool"
)

// Client exposes the subset of CometBFT RPC functionality required by the SDK.
// Implementations should be safe for concurrent use.
type Client interface {
	pool.PoolParticipant
	BlockResults(ctx context.Context, height *int64) (*coretypes.ResultBlockResults, error)
	Block(ctx context.Context, height *int64) (*coretypes.ResultBlock, error)
	Commit(ctx context.Context, height *int64) (*coretypes.ResultCommit, error)
	ABCIQuery(ctx context.Context, path string, data []byte) (*coretypes.ResultABCIQuery, error)
	Status(ctx context.Context) (*coretypes.ResultStatus, error)
}

// NewHTTPClient constructs a Tendermint RPC client backed by CometBFT's HTTP implementation.
// The remote must include the scheme, e.g. http://node:26657.
func NewHTTPClient(remote, wsURL string, timeout time.Duration, logger zerolog.Logger) (Client, error) {
	var (
		tmClient *tmhttp.HTTP
		err      error
	)

	if timeout > 0 {
		seconds := uint((timeout + time.Second - 1) / time.Second)
		tmClient, err = tmhttp.NewWithTimeout(remote, wsURL, seconds)
	} else {
		tmClient, err = tmhttp.New(remote, wsURL)
	}
	if err != nil {
		return nil, fmt.Errorf("tmrpc: failed to create HTTP client for %s: %w", remote, err)
	}

	return &httpClient{
		remote: remote,
		http:   tmClient,
		logger: logger.With().Str("component", "tmrpc_client").Str("remote", remote).Logger(),
	}, nil
}

type httpClient struct {
	remote string
	http   *tmhttp.HTTP
	logger zerolog.Logger
}

func (c *httpClient) HealthCheck(ctx context.Context) error {
	_, err := c.Status(ctx)
	return err
}

func (c *httpClient) GetEndpointURL() string {
	return c.remote
}

func (c *httpClient) GetProtocol() config.Protocol {
	return config.ProtocolTendermintRPC
}

func (c *httpClient) BlockResults(ctx context.Context, height *int64) (*coretypes.ResultBlockResults, error) {
	return c.http.BlockResults(ctx, height)
}

func (c *httpClient) Block(ctx context.Context, height *int64) (*coretypes.ResultBlock, error) {
	return c.http.Block(ctx, height)
}

func (c *httpClient) Commit(ctx context.Context, height *int64) (*coretypes.ResultCommit, error) {
	return c.http.Commit(ctx, height)
}

func (c *httpClient) ABCIQuery(ctx context.Context, path string, data []byte) (*coretypes.ResultABCIQuery, error) {
	return c.http.ABCIQuery(ctx, path, data)
}

func (c *httpClient) Status(ctx context.Context) (*coretypes.ResultStatus, error) {
	return c.http.Status(ctx)
}

func (c *httpClient) Close() error {
	if err := c.http.Stop(); err != nil && !errors.Is(err, context.Canceled) {
		c.logger.Debug().Err(err).Msg("error stopping tendermint rpc client")
		return err
	}
	return nil
}

type ClientPool interface {
	Close()
	GetHealthStatus() map[string]any

	BlockResults(ctx context.Context, height *int64) (*coretypes.ResultBlockResults, error)
	Block(ctx context.Context, height *int64) (*coretypes.ResultBlock, error)
	Commit(ctx context.Context, height *int64) (*coretypes.ResultCommit, error)
	ABCIQuery(ctx context.Context, path string, data []byte) (*coretypes.ResultABCIQuery, error)
	Status(ctx context.Context) (*coretypes.ResultStatus, error)
	HealthCheck(ctx context.Context) error
}

type clientPool struct {
	logger      zerolog.Logger
	poolManager *pool.ClientPoolManager[Client]
}

// NewPool constructs a pool from the provided clients. The pool will panic if
// no clients are supplied.
func NewClientPool(clients []Client, logger zerolog.Logger) *clientPool {
	poolLogger := logger.With().Str("component", "tmrpc_pool").Logger()
	return &clientPool{
		logger:      poolLogger,
		poolManager: pool.NewClientPoolManager(clients, poolLogger),
	}
}

func (p *clientPool) Close() {
	p.poolManager.Close()
}

func (p *clientPool) GetHealthStatus() map[string]any {
	return p.poolManager.GetHealthStatus()
}

func (p *clientPool) BlockResults(ctx context.Context, height *int64) (*coretypes.ResultBlockResults, error) {
	return pool.ExecuteWithRetry(ctx, p.poolManager, &p.logger, func(c Client) (*coretypes.ResultBlockResults, error) {
		return c.BlockResults(ctx, height)
	})
}

func (p *clientPool) Block(ctx context.Context, height *int64) (*coretypes.ResultBlock, error) {
	return pool.ExecuteWithRetry(ctx, p.poolManager, &p.logger, func(c Client) (*coretypes.ResultBlock, error) {
		return c.Block(ctx, height)
	})
}

func (p *clientPool) Commit(ctx context.Context, height *int64) (*coretypes.ResultCommit, error) {
	return pool.ExecuteWithRetry(ctx, p.poolManager, &p.logger, func(c Client) (*coretypes.ResultCommit, error) {
		return c.Commit(ctx, height)
	})
}

func (p *clientPool) ABCIQuery(ctx context.Context, path string, data []byte) (*coretypes.ResultABCIQuery, error) {
	return pool.ExecuteWithRetry(ctx, p.poolManager, &p.logger, func(c Client) (*coretypes.ResultABCIQuery, error) {
		return c.ABCIQuery(ctx, path, data)
	})
}

func (p *clientPool) Status(ctx context.Context) (*coretypes.ResultStatus, error) {
	return pool.ExecuteWithRetry(ctx, p.poolManager, &p.logger, func(c Client) (*coretypes.ResultStatus, error) {
		return c.Status(ctx)
	})
}

func (p *clientPool) HealthCheck(ctx context.Context) error {
	_, err := pool.ExecuteWithRetry(ctx, p.poolManager, &p.logger, func(c Client) (struct{}, error) {
		return struct{}{}, c.HealthCheck(ctx)
	})
	return err
}
