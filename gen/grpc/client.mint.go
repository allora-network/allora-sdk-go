package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	minttypes "github.com/allora-network/allora-chain/x/mint/types"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// MintGRPCClient provides gRPC access to the mint module
type MintGRPCClient struct {
	client minttypes.QueryServiceClient
	logger zerolog.Logger
}

var _ interfaces.MintClient = (*MintGRPCClient)(nil)

// NewMintGRPCClient creates a new mint REST client
func NewMintGRPCClient(conn *grpc.ClientConn, logger zerolog.Logger) *MintGRPCClient {
	return &MintGRPCClient{
		client: minttypes.NewQueryServiceClient(conn),
		logger: logger.With().Str("module", "mint").Str("protocol", "grpc").Logger(),
	}
}

func (c *MintGRPCClient) EmissionInfo(ctx context.Context, req *minttypes.QueryServiceEmissionInfoRequest, opts ...config.CallOpt) (*minttypes.QueryServiceEmissionInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.EmissionInfo, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling MintGRPCClient.EmissionInfo")
	}
	return resp, nil
}

func (c *MintGRPCClient) Inflation(ctx context.Context, req *minttypes.QueryServiceInflationRequest, opts ...config.CallOpt) (*minttypes.QueryServiceInflationResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Inflation, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling MintGRPCClient.Inflation")
	}
	return resp, nil
}

func (c *MintGRPCClient) Params(ctx context.Context, req *minttypes.QueryServiceParamsRequest, opts ...config.CallOpt) (*minttypes.QueryServiceParamsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Params, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling MintGRPCClient.Params")
	}
	return resp, nil
}
