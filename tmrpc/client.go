package tmrpc

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	tmhttp "github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/rs/zerolog"
)

// Client exposes the subset of CometBFT RPC functionality required by the SDK.
// Implementations should be safe for concurrent use.
type Client interface {
	BlockResults(ctx context.Context, height *int64) (*coretypes.ResultBlockResults, error)
	Block(ctx context.Context, height *int64) (*coretypes.ResultBlock, error)
	Commit(ctx context.Context, height *int64) (*coretypes.ResultCommit, error)
	ABCIQuery(ctx context.Context, path string, data []byte) (*coretypes.ResultABCIQuery, error)
	Status(ctx context.Context) (*coretypes.ResultStatus, error)
	Close() error
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

func (c *httpClient) endpoint() string {
	return c.remote
}

type pooledClient struct {
	Client
	endpoint string
}

// Pool provides simple round-robin load balancing with retry semantics across
// multiple Tendermint RPC endpoints.
type Pool struct {
	logger  zerolog.Logger
	clients []pooledClient
	mu      sync.Mutex
	nextIdx int
}

// NewPool constructs a pool from the provided clients. The pool will panic if
// no clients are supplied.
func NewPool(clients []Client, logger zerolog.Logger) *Pool {
	if len(clients) == 0 {
		panic("tmrpc: NewPool requires at least one client")
	}

	pooled := make([]pooledClient, len(clients))
	for i, client := range clients {
		endpoint := ""
		if provider, ok := client.(interface{ endpoint() string }); ok {
			endpoint = provider.endpoint()
		}
		pooled[i] = pooledClient{Client: client, endpoint: endpoint}
	}

	return &Pool{
		logger:  logger.With().Str("component", "tmrpc_pool").Logger(),
		clients: pooled,
	}
}

func (p *Pool) Close() error {
	var errs error
	for _, client := range p.clients {
		if err := client.Close(); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func (p *Pool) BlockResults(ctx context.Context, height *int64) (*coretypes.ResultBlockResults, error) {
	return doWithRetry(ctx, p, "BlockResults", func(c Client) (*coretypes.ResultBlockResults, error) {
		return c.BlockResults(ctx, height)
	})
}

func (p *Pool) Block(ctx context.Context, height *int64) (*coretypes.ResultBlock, error) {
	return doWithRetry(ctx, p, "Block", func(c Client) (*coretypes.ResultBlock, error) {
		return c.Block(ctx, height)
	})
}

func (p *Pool) Commit(ctx context.Context, height *int64) (*coretypes.ResultCommit, error) {
	return doWithRetry(ctx, p, "Commit", func(c Client) (*coretypes.ResultCommit, error) {
		return c.Commit(ctx, height)
	})
}

func (p *Pool) ABCIQuery(ctx context.Context, path string, data []byte) (*coretypes.ResultABCIQuery, error) {
	return doWithRetry(ctx, p, "ABCIQuery", func(c Client) (*coretypes.ResultABCIQuery, error) {
		return c.ABCIQuery(ctx, path, data)
	})
}

func (p *Pool) Status(ctx context.Context) (*coretypes.ResultStatus, error) {
	return doWithRetry(ctx, p, "Status", func(c Client) (*coretypes.ResultStatus, error) {
		return c.Status(ctx)
	})
}

func doWithRetry[R any](ctx context.Context, p *Pool, op string, fn func(Client) (R, error)) (R, error) {
	var zero R
	if len(p.clients) == 0 {
		return zero, fmt.Errorf("tmrpc: no clients configured for %s", op)
	}

	var lastErr error
	attempts := len(p.clients)
	for attempt := 0; attempt < attempts; attempt++ {
		select {
		case <-ctx.Done():
			return zero, ctx.Err()
		default:
		}

		client := p.next()
		result, err := fn(client.Client)
		if err != nil {
			lastErr = err
			p.logger.Debug().
				Str("endpoint", client.endpoint).
				Int("attempt", attempt+1).
				Err(err).
				Msgf("tendermint rpc %s failed", op)
			continue
		}
		return result, nil
	}

	if lastErr != nil {
		return zero, fmt.Errorf("tmrpc: all clients failed for %s, last error: %w", op, lastErr)
	}
	return zero, fmt.Errorf("tmrpc: no clients available for %s", op)
}

func (p *Pool) next() pooledClient {
	p.mu.Lock()
	defer p.mu.Unlock()

	client := p.clients[p.nextIdx]
	p.nextIdx = (p.nextIdx + 1) % len(p.clients)
	return client
}

var _ Client = (*Pool)(nil)
