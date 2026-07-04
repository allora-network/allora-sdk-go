package allora

import (
	"fmt"
	"sync"

	butils "github.com/brynbellomy/go-utils"
	ctypes "github.com/cometbft/cometbft/types"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/cosmosrpc"
	"github.com/allora-network/allora-sdk-go/gen/grpc"
	"github.com/allora-network/allora-sdk-go/gen/rest"
	"github.com/allora-network/allora-sdk-go/metrics"
	"github.com/allora-network/allora-sdk-go/tmrpc"
	"github.com/allora-network/allora-sdk-go/txsend/cosmospool"
)

// Client is the top-level Allora Network client. It is the single handle
// callers hold: it exposes the pooled query surfaces (Cosmos, Tendermint) and
// the transaction-send surface (Tx) through one constructed value.
type Client interface {
	// Cosmos returns the pooled Cosmos gRPC/REST client used for queries
	// (bank, auth, staking, etc.) and as the backing pool for the Tx sender.
	Cosmos() cosmosrpc.ClientPool

	// Tendermint returns the pooled Tendermint RPC client used for block,
	// validator, and subscription queries.
	Tendermint() tmrpc.ClientPool

	// Subscribe registers a subscription for Tendermint events matching query,
	// delivering matching events into mb. It is a thin wrapper over the
	// websocket pool.
	Subscribe(mb *butils.Mailbox[ctypes.TMEventData], query string)

	// Tx returns the Sender that drives a message set through the full send
	// lifecycle (account discovery, build, sign, gas estimation, broadcast,
	// confirmation). It is backed by a cosmospool.Broadcaster constructed over
	// the client's Cosmos pool — the same pooled, health-tracked client the
	// query surface uses — so sending reuses the existing connection and retry
	// machinery rather than opening a new path. The Sender is constructed once
	// on first call and cached for the lifetime of the client; it is safe to
	// call Tx() concurrently and to retain the returned Sender.
	Tx() Sender
}

type client struct {
	config         *config.ClientConfig
	cosmosPool     cosmosrpc.ClientPool
	tendermintPool tmrpc.ClientPool
	websocketPool  tmrpc.WebsocketPool
	logger         zerolog.Logger

	// sender is lazily constructed on the first Tx() call so that a client
	// constructed only for queries never pays the cosmospool wiring cost. It is
	// guarded by senderMu; once set it is never replaced.
	sender   Sender
	senderMu sync.Mutex
}

var _ Client = (*client)(nil)

func NewClient(cfg *config.ClientConfig, logger zerolog.Logger) (*client, error) {
	if cfg == nil {
		cfg = config.DefaultClientConfig()
	}

	if len(cfg.Endpoints) == 0 {
		return nil, fmt.Errorf("at least one endpoint must be specified")
	}

	cosmosClients := []cosmosrpc.Client{}
	tmRPCClients := []tmrpc.Client{}
	websockets := []tmrpc.Websocket{}
	for _, endpoint := range cfg.Endpoints {
		switch endpoint.Protocol {
		case config.ProtocolGRPC:
			client, err := grpc.NewGRPCClient(endpoint, logger)
			if err != nil {
				logger.Error().
					Str("endpoint", endpoint.URL).
					Str("protocol", string(endpoint.Protocol)).
					Err(err).
					Msg("failed to create client for endpoint")
				continue
			}
			cosmosClients = append(cosmosClients, client)
		case config.ProtocolREST:
			client := rest.NewRESTClient(endpoint.URL, logger)
			cosmosClients = append(cosmosClients, client)
		case config.ProtocolTendermintRPC:
			if endpoint.URL != "" {
				tmClient, err := tmrpc.NewHTTPClient(endpoint.URL, endpoint.WebsocketURL, cfg.RequestTimeout, logger)
				if err != nil {
					logger.Error().
						Str("endpoint", endpoint.URL).
						Str("protocol", string(endpoint.Protocol)).
						Err(err).
						Msg("failed to create tendermint rpc client for endpoint")
					continue
				}
				tmRPCClients = append(tmRPCClients, tmClient)
			}
			if endpoint.WebsocketURL != "" {
				ws := tmrpc.NewTendermintWebsocket(endpoint.WebsocketURL, logger)
				websockets = append(websockets, ws)
			}
		case "":
			logger.Error().
				Str("endpoint", endpoint.URL).
				Msg("no protocol specified")
		default:
			logger.Error().
				Str("endpoint", endpoint.URL).
				Str("protocol", string(endpoint.Protocol)).
				Msgf("unsupported protocol")
			continue
		}
	}

	// Validate that at least one client was successfully created
	if len(cosmosClients) == 0 && len(tmRPCClients) == 0 {
		return nil, fmt.Errorf("failed to create any valid clients from the provided endpoints")
	}

	logger = logger.With().
		Str("component", "allora_client").
		Logger()

	return &client{
		cosmosPool:     cosmosrpc.NewClientPool(cosmosClients, logger),
		tendermintPool: tmrpc.NewClientPool(tmRPCClients, logger),
		websocketPool:  tmrpc.NewWebsocketPool(websockets),
		logger:         logger,
		config:         cfg,
	}, nil
}

func (c *client) Close() error {
	c.logger.Info().Msg("shutting down Allora client")
	c.cosmosPool.Close()
	c.tendermintPool.Close()
	c.websocketPool.Close()
	return nil
}

func (c *client) GetHealthStatus() map[string]any {
	return map[string]any{
		"cosmos":     c.cosmosPool.GetHealthStatus(),
		"tendermint": c.tendermintPool.GetHealthStatus(),
	}
}

func (c *client) Cosmos() cosmosrpc.ClientPool {
	return c.cosmosPool
}

func (c *client) Tendermint() tmrpc.ClientPool {
	return c.tendermintPool
}

func (c *client) Subscribe(mb *butils.Mailbox[ctypes.TMEventData], query string) {
	c.websocketPool.Subscribe(mb, query)
}

// Tx returns the Sender backed by a cosmospool.Broadcaster constructed over
// the client's Cosmos pool. It is constructed lazily on first call and cached;
// subsequent calls return the same Sender. Safe for concurrent use.
func (c *client) Tx() Sender {
	c.senderMu.Lock()
	defer c.senderMu.Unlock()
	if c.sender == nil {
		opts := []cosmospool.Option{}
		// Only override the broadcaster default when the caller explicitly set a
		// gas adjustment; a zero value keeps the broadcaster's default (1.5).
		if c.config.GasAdjustment > 0 {
			opts = append(opts, cosmospool.WithGasAdjustment(c.config.GasAdjustment))
		}
		broadcaster := cosmospool.New(c.cosmosPool, c.logger, opts...)
		c.sender = NewSender(broadcaster, c.logger)
	}
	return c.sender
}

var SetMetricsPrefix = metrics.SetPrefix
