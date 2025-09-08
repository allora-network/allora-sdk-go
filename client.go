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
)

type Client interface {
	interfaces.Client
	Subscribe(mb *butils.Mailbox[ctypes.TMEventData], query string)
	GetHealthStatus() map[string]any
}

// Client is the Allora Network client that provides access to all query services.
// It manages a pool of underlying gRPC and REST clients with automatic load balancing and failover.
type client struct {
	poolManager *pool.ClientPoolManager
	logger      zerolog.Logger
	config      *config.ClientConfig

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

	clients := make([]pool.Client, 0, len(cfg.Endpoints))
	for _, endpoint := range cfg.Endpoints {
		var client pool.Client
		var err error

		switch endpoint.Protocol {
		case config.ProtocolGRPC:
			client, err = grpc.NewGRPCClient(endpoint, logger)
		case config.ProtocolREST:
			client = rest.NewRESTClient(endpoint.URL, logger)
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

		if err != nil {
			// Log the error but continue with other endpoints
			logger.Error().
				Str("endpoint", endpoint.URL).
				Str("protocol", string(endpoint.Protocol)).
				Err(err).
				Msg("failed to create client for endpoint")
			continue
		}
		clients = append(clients, client)
	}

	if len(clients) == 0 {
		return nil, fmt.Errorf("failed to create any clients from the provided endpoints")
	}

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
	}, nil
}

func (c *client) Close() error {
	c.logger.Info().Msg("shutting down Allora client")
	if c.CometRPCWebsocket != nil {
		c.CometRPCWebsocket.Close()
	}
	c.poolManager.Close()
	return nil
}

// GetHealthStatus returns the current health status of all clients in the pool
func (c *client) GetHealthStatus() map[string]any {
	return c.poolManager.GetHealthStatus()
}

var SetMetricsPrefix = metrics.SetPrefix
