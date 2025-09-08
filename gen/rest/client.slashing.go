package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	"github.com/allora-network/allora-sdk-go/config"
)

// SlashingRESTClient provides REST access to the slashing module
type SlashingRESTClient struct {
	*RESTClientCore
	logger zerolog.Logger
}

// NewSlashingRESTClient creates a new slashing REST client
func NewSlashingRESTClient(core *RESTClientCore, logger zerolog.Logger) *SlashingRESTClient {
	return &SlashingRESTClient{
		RESTClientCore: core,
		logger:         logger.With().Str("module", "slashing").Str("protocol", "rest").Logger(),
	}
}

func (c *SlashingRESTClient) Params(ctx context.Context, req *slashingtypes.QueryParamsRequest, opts ...config.CallOpt) (*slashingtypes.QueryParamsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &slashingtypes.QueryParamsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/slashing/v1beta1/params",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling SlashingRESTClient.Params")
	}
	return resp, nil
}

func (c *SlashingRESTClient) SigningInfo(ctx context.Context, req *slashingtypes.QuerySigningInfoRequest, opts ...config.CallOpt) (*slashingtypes.QuerySigningInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &slashingtypes.QuerySigningInfoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/slashing/v1beta1/signing_infos/{cons_address}",
		[]string{"cons_address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling SlashingRESTClient.SigningInfo")
	}
	return resp, nil
}

func (c *SlashingRESTClient) SigningInfos(ctx context.Context, req *slashingtypes.QuerySigningInfosRequest, opts ...config.CallOpt) (*slashingtypes.QuerySigningInfosResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &slashingtypes.QuerySigningInfosResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/slashing/v1beta1/signing_infos",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling SlashingRESTClient.SigningInfos")
	}
	return resp, nil
}
