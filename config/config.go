package config

import (
	"time"
)

type ClientConfig struct {
	Endpoints         []EndpointConfig
	RequestTimeout    time.Duration
	ConnectionTimeout time.Duration

	// GasAdjustment is the multiplier applied to simulated gas usage to derive
	// the gas limit for transactions sent via Client.Tx(). It defaults to the
	// broadcaster's default (1.5) when zero; set it to override (e.g. 1.2 to
	// match a node's configured minimum-gas-adjustment). Values below 1.0 are
	// clamped to 1.0 by the broadcaster.
	GasAdjustment float64
}

type EndpointConfig struct {
	URL          string
	WebsocketURL string
	Protocol     Protocol
}

type Protocol string

const (
	ProtocolGRPC          Protocol = "grpc"
	ProtocolREST          Protocol = "rest"
	ProtocolTendermintRPC Protocol = "tendermint_rpc"
)

func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		RequestTimeout:    30 * time.Second,
		ConnectionTimeout: 10 * time.Second,
	}
}
