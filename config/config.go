package config

import (
	"time"
)

// ClientConfig represents the configuration for the Allora client
type ClientConfig struct {
	// Endpoints is a list of RPC endpoints to connect to
	Endpoints []EndpointConfig

	// Timeout for individual requests
	RequestTimeout time.Duration

	// Connection timeout
	ConnectionTimeout time.Duration
}

type Protocol string

const (
	ProtocolGRPC Protocol = "grpc"
	ProtocolREST Protocol = "rest"
)

// EndpointConfig represents a single endpoint configuration
type EndpointConfig struct {
	URL      string
	Protocol Protocol
}

// DefaultClientConfig returns a default client configuration
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		RequestTimeout:    30 * time.Second,
		ConnectionTimeout: 10 * time.Second,
	}
}
