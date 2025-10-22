package allora

import (
	"fmt"

	butils "github.com/brynbellomy/go-utils"
	ctypes "github.com/cometbft/cometbft/types"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/cosmosrpc"
	"github.com/allora-network/allora-sdk-go/gen/grpc"
	"github.com/allora-network/allora-sdk-go/gen/rest"
	"github.com/allora-network/allora-sdk-go/metrics"
	"github.com/allora-network/allora-sdk-go/tmrpc"
)

type Client interface {
	Cosmos() cosmosrpc.ClientPool
	Tendermint() tmrpc.ClientPool
	Subscribe(mb *butils.Mailbox[ctypes.TMEventData], query string)
}

// Client is the Allora Network client that provides access to all query services.
// It manages a pool of underlying gRPC and REST clients with automatic load balancing and failover.
type client struct {
	config         *config.ClientConfig
	cosmosPool     cosmosrpc.ClientPool
	tendermintPool tmrpc.ClientPool
	websocketPool  tmrpc.WebsocketPool
	logger         zerolog.Logger
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

var SetMetricsPrefix = metrics.SetPrefix
