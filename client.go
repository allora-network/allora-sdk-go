package allorasdk

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/grpc"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/gen/rest"
	"github.com/allora-network/allora-sdk-go/gen/wrapper"
	"github.com/allora-network/allora-sdk-go/pool"
)

// Client is the Allora Network client that provides access to all query services.
// It manages a pool of underlying gRPC and REST clients with automatic load balancing and failover.
type Client struct {
	poolManager *pool.ClientPoolManager
	logger      zerolog.Logger
	config      *config.ClientConfig
	*wrapper.WrapperClient
}

// NewClient creates a new Allora Network client
func NewClient(cfg *config.ClientConfig, logger zerolog.Logger) (*Client, error) {
	if cfg == nil {
		cfg = config.DefaultClientConfig()
	}

	if len(cfg.Endpoints) == 0 {
		return nil, fmt.Errorf("at least one endpoint must be specified")
	}

	clients := make([]interfaces.Client, 0, len(cfg.Endpoints))
	for _, endpoint := range cfg.Endpoints {
		var client interfaces.Client
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

	return &Client{
		WrapperClient: wrapper.NewWrapperClient(poolManager, clientLogger),
		poolManager:   poolManager,
		logger:        clientLogger,
		config:        cfg,
	}, nil
}

func (c *Client) Close() {
	c.logger.Info().Msg("shutting down Allora client")
	c.poolManager.Close()
}

// GetHealthStatus returns the current health status of all clients in the pool
func (c *Client) GetHealthStatus() map[string]any {
	return c.poolManager.GetHealthStatus()
}
