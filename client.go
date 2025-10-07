package allora

import (
	"fmt"
	"strings"

	butils "github.com/brynbellomy/go-utils"
	ctypes "github.com/cometbft/cometbft/types"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/grpc"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/gen/rest"
	"github.com/allora-network/allora-sdk-go/gen/wrapper"
	"github.com/allora-network/allora-sdk-go/metrics"
	"github.com/allora-network/allora-sdk-go/pool"
	"github.com/allora-network/allora-sdk-go/tmrpc"
)

type Client interface {
	interfaces.Client
	Subscribe(mb *butils.Mailbox[ctypes.TMEventData], query string)
	GetHealthStatus() map[string]any
	TendermintRPC() tmrpc.Client
}

// Client is the Allora Network client that provides access to all query services.
// It manages a pool of underlying gRPC and REST clients with automatic load balancing and failover.
type client struct {
	poolManager *pool.ClientPoolManager
	logger      zerolog.Logger
	config      *config.ClientConfig
	tmRPC       tmrpc.Client

	*wrapper.WrapperClient
	*CometRPCWebsocket
}

func NewClient(cfg *config.ClientConfig, logger zerolog.Logger) (*client, error) {
	if cfg == nil {
		cfg = config.DefaultClientConfig()
	}

	if len(cfg.Endpoints) == 0 {
		return nil, fmt.Errorf("at least one endpoint must be specified")
	}

	clients := []pool.Client{}
	tmRPCClients := []tmrpc.Client{}
	for _, endpoint := range cfg.Endpoints {
		fmt.Println("PROTO:", endpoint.Protocol)
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
			clients = append(clients, client)
		case config.ProtocolREST:
			client := rest.NewRESTClient(endpoint.URL, logger)
			clients = append(clients, client)
		case config.ProtocolTendermintRPC:
			tmClient, err := tmrpc.NewHTTPClient(endpoint.URL, cfg.WebsocketEndpoint, cfg.RequestTimeout, logger)
			if err != nil {
				logger.Error().
					Str("endpoint", endpoint.URL).
					Str("protocol", string(endpoint.Protocol)).
					Err(err).
					Msg("failed to create tendermint rpc client for endpoint")
				continue
			}
			tmRPCClients = append(tmRPCClients, tmClient)
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

	if len(clients) == 0 {
		return nil, fmt.Errorf("failed to create any clients from the provided endpoints")
	}

	var tmPool tmrpc.Client
	if len(tmRPCClients) > 0 {
		tmPool = tmrpc.NewPool(tmRPCClients, logger)
	}
	fmt.Println("LEN TM RPC", len(tmRPCClients))

	poolManager := pool.NewClientPoolManager(clients, logger)
	clientLogger := logger.With().Str("component", "allora_client").Logger()

	var ws *CometRPCWebsocket
	if strings.TrimSpace(cfg.WebsocketEndpoint) != "" {
		ws = NewCometRPCWebsocket(strings.TrimSpace(cfg.WebsocketEndpoint), logger)
	}

	return &client{
		WrapperClient:     wrapper.NewWrapperClient(poolManager, clientLogger),
		CometRPCWebsocket: ws,
		poolManager:       poolManager,
		logger:            clientLogger,
		config:            cfg,
		tmRPC:             tmPool,
	}, nil
}

func (c *client) Close() error {
	c.logger.Info().Msg("shutting down Allora client")
	if c.CometRPCWebsocket != nil {
		c.CometRPCWebsocket.Close()
	}
	if c.tmRPC != nil {
		if err := c.tmRPC.Close(); err != nil {
			c.logger.Error().Err(err).Msg("error closing tendermint rpc clients")
		}
	}
	c.poolManager.Close()
	return nil
}

// GetHealthStatus returns the current health status of all clients in the pool
func (c *client) GetHealthStatus() map[string]any {
	return c.poolManager.GetHealthStatus()
}

func (c *client) TendermintRPC() tmrpc.Client {
	return c.tmRPC
}

var SetMetricsPrefix = metrics.SetPrefix
