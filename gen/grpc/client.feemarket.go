package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// FeemarketGRPCClient provides gRPC access to the feemarket module
type FeemarketGRPCClient struct {
	client feemarkettypes.QueryClient
	logger zerolog.Logger
}

var _ interfaces.FeemarketClient = (*FeemarketGRPCClient)(nil)

// NewFeemarketGRPCClient creates a new feemarket REST client
func NewFeemarketGRPCClient(conn *grpc.ClientConn, logger zerolog.Logger) *FeemarketGRPCClient {
	return &FeemarketGRPCClient{
		client: feemarkettypes.NewQueryClient(conn),
		logger: logger.With().Str("module", "feemarket").Str("protocol", "grpc").Logger(),
	}
}

func (c *FeemarketGRPCClient) GasPrice(ctx context.Context, req *feemarkettypes.GasPriceRequest, opts ...config.CallOpt) (*feemarkettypes.GasPriceResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GasPrice, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling FeemarketGRPCClient.GasPrice")
	}
	return resp, nil
}

func (c *FeemarketGRPCClient) GasPrices(ctx context.Context, req *feemarkettypes.GasPricesRequest, opts ...config.CallOpt) (*feemarkettypes.GasPricesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GasPrices, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling FeemarketGRPCClient.GasPrices")
	}
	return resp, nil
}

func (c *FeemarketGRPCClient) Params(ctx context.Context, req *feemarkettypes.ParamsRequest, opts ...config.CallOpt) (*feemarkettypes.ParamsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Params, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling FeemarketGRPCClient.Params")
	}
	return resp, nil
}

func (c *FeemarketGRPCClient) State(ctx context.Context, req *feemarkettypes.StateRequest, opts ...config.CallOpt) (*feemarkettypes.StateResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.State, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling FeemarketGRPCClient.State")
	}
	return resp, nil
}
