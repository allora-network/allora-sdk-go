package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"

	"github.com/allora-network/allora-sdk-go/config"
)

// FeemarketRESTClient provides REST access to the feemarket module
type FeemarketRESTClient struct {
	*RESTClientCore
	logger zerolog.Logger
}

// NewFeemarketRESTClient creates a new feemarket REST client
func NewFeemarketRESTClient(core *RESTClientCore, logger zerolog.Logger) *FeemarketRESTClient {
	return &FeemarketRESTClient{
		RESTClientCore: core,
		logger:         logger.With().Str("module", "feemarket").Str("protocol", "rest").Logger(),
	}
}

func (c *FeemarketRESTClient) GasPrice(ctx context.Context, req *feemarkettypes.GasPriceRequest, opts ...config.CallOpt) (*feemarkettypes.GasPriceResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &feemarkettypes.GasPriceResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/feemarket/v1/gas_price/{denom}",
		[]string{"denom"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling FeemarketRESTClient.GasPrice")
	}
	return resp, nil
}

func (c *FeemarketRESTClient) GasPrices(ctx context.Context, req *feemarkettypes.GasPricesRequest, opts ...config.CallOpt) (*feemarkettypes.GasPricesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &feemarkettypes.GasPricesResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/feemarket/v1/gas_prices",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling FeemarketRESTClient.GasPrices")
	}
	return resp, nil
}

func (c *FeemarketRESTClient) Params(ctx context.Context, req *feemarkettypes.ParamsRequest, opts ...config.CallOpt) (*feemarkettypes.ParamsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &feemarkettypes.ParamsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/feemarket/v1/params",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling FeemarketRESTClient.Params")
	}
	return resp, nil
}

func (c *FeemarketRESTClient) State(ctx context.Context, req *feemarkettypes.StateRequest, opts ...config.CallOpt) (*feemarkettypes.StateResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &feemarkettypes.StateResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/feemarket/v1/state",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling FeemarketRESTClient.State")
	}
	return resp, nil
}
