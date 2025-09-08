package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	node "github.com/cosmos/cosmos-sdk/client/grpc/node"
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

func (c *NodeRESTClient) Config(ctx context.Context, req *node.ConfigRequest) (*node.ConfigResponse, error) {
	resp := &node.ConfigResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/base/node/v1beta1/config",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling NodeRESTClient.Config")
}

func (c *NodeRESTClient) Status(ctx context.Context, req *node.StatusRequest) (*node.StatusResponse, error) {
	resp := &node.StatusResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/base/node/v1beta1/status",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling NodeRESTClient.Status")
}
