package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
)

// ConsensusRESTClient provides REST access to the consensus module
type ConsensusRESTClient struct {
	*RESTClientCore
	logger zerolog.Logger
}

// NewConsensusRESTClient creates a new consensus REST client
func NewConsensusRESTClient(core *RESTClientCore, logger zerolog.Logger) *ConsensusRESTClient {
	return &ConsensusRESTClient{
		RESTClientCore: core,
		logger:         logger.With().Str("module", "consensus").Str("protocol", "rest").Logger(),
	}
}

func (c *ConsensusRESTClient) Params(ctx context.Context, req *consensustypes.QueryParamsRequest) (*consensustypes.QueryParamsResponse, error) {
	resp := &consensustypes.QueryParamsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/consensus/v1/params",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling ConsensusRESTClient.Params")
}
