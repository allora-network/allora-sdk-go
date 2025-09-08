package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	node "github.com/cosmos/cosmos-sdk/client/grpc/node"

	"github.com/allora-network/allora-sdk-go/config"
)

// NodeRESTClient provides REST access to the node module
type NodeRESTClient struct {
	*RESTClientCore
	logger zerolog.Logger
}

// NewNodeRESTClient creates a new node REST client
func NewNodeRESTClient(core *RESTClientCore, logger zerolog.Logger) *NodeRESTClient {
	return &NodeRESTClient{
		RESTClientCore: core,
		logger:         logger.With().Str("module", "node").Str("protocol", "rest").Logger(),
	}
}

func (c *NodeRESTClient) Config(ctx context.Context, req *node.ConfigRequest, opts ...config.CallOpt) (*node.ConfigResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &node.ConfigResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/base/node/v1beta1/config",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling NodeRESTClient.Config")
	}
	return resp, nil
}

func (c *NodeRESTClient) Status(ctx context.Context, req *node.StatusRequest, opts ...config.CallOpt) (*node.StatusResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &node.StatusResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/base/node/v1beta1/status",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling NodeRESTClient.Status")
	}
	return resp, nil
}
