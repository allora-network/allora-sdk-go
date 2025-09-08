package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
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

func (c *SlashingRESTClient) Params(ctx context.Context, req *slashingtypes.QueryParamsRequest) (*slashingtypes.QueryParamsResponse, error) {
	resp := &slashingtypes.QueryParamsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/slashing/v1beta1/params",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling SlashingRESTClient.Params")
}

func (c *SlashingRESTClient) SigningInfo(ctx context.Context, req *slashingtypes.QuerySigningInfoRequest) (*slashingtypes.QuerySigningInfoResponse, error) {
	resp := &slashingtypes.QuerySigningInfoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/slashing/v1beta1/signing_infos/{cons_address}",
		[]string{"cons_address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling SlashingRESTClient.SigningInfo")
}

func (c *SlashingRESTClient) SigningInfos(ctx context.Context, req *slashingtypes.QuerySigningInfosRequest) (*slashingtypes.QuerySigningInfosResponse, error) {
	resp := &slashingtypes.QuerySigningInfosResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/slashing/v1beta1/signing_infos",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling SlashingRESTClient.SigningInfos")
}
