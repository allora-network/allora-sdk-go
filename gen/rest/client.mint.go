package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	minttypes "github.com/allora-network/allora-chain/x/mint/types"

	"github.com/allora-network/allora-sdk-go/config"
)

// MintRESTClient provides REST access to the mint module
type MintRESTClient struct {
	*RESTClientCore
	logger zerolog.Logger
}

// NewMintRESTClient creates a new mint REST client
func NewMintRESTClient(core *RESTClientCore, logger zerolog.Logger) *MintRESTClient {
	return &MintRESTClient{
		RESTClientCore: core,
		logger:         logger.With().Str("module", "mint").Str("protocol", "rest").Logger(),
	}
}

func (c *MintRESTClient) Params(ctx context.Context, req *minttypes.QueryServiceParamsRequest, opts ...config.CallOpt) (*minttypes.QueryServiceParamsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &minttypes.QueryServiceParamsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/mint/v5/params",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling MintRESTClient.Params")
	}
	return resp, nil
}

func (c *MintRESTClient) Inflation(ctx context.Context, req *minttypes.QueryServiceInflationRequest, opts ...config.CallOpt) (*minttypes.QueryServiceInflationResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &minttypes.QueryServiceInflationResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/mint/v5/inflation",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling MintRESTClient.Inflation")
	}
	return resp, nil
}

func (c *MintRESTClient) EmissionInfo(ctx context.Context, req *minttypes.QueryServiceEmissionInfoRequest, opts ...config.CallOpt) (*minttypes.QueryServiceEmissionInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &minttypes.QueryServiceEmissionInfoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/mint/v5/emission_info",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling MintRESTClient.EmissionInfo")
	}
	return resp, nil
}
