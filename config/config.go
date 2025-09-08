package config

import (
	"time"
)

type ClientConfig struct {
	Endpoints         []EndpointConfig
	WebsocketEndpoint string
	RequestTimeout    time.Duration
	ConnectionTimeout time.Duration
}

type EndpointConfig struct {
	URL      string
	Protocol Protocol
}

type Protocol string

const (
	ProtocolGRPC Protocol = "grpc"
	ProtocolREST Protocol = "rest"
)

func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		RequestTimeout:    30 * time.Second,
		ConnectionTimeout: 10 * time.Second,
	}
}
