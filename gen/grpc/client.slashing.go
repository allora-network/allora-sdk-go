package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/allora-network/allora-sdk-go/config"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// SlashingGRPCClient provides gRPC access to the slashing module
type SlashingGRPCClient struct {
	client slashingtypes.QueryClient
	logger zerolog.Logger
}

var _ interfaces.SlashingClient = (*SlashingGRPCClient)(nil)

// NewSlashingGRPCClient creates a new slashing REST client
func NewSlashingGRPCClient(conn *grpc.ClientConn, logger zerolog.Logger) *SlashingGRPCClient {
	return &SlashingGRPCClient{
		client: slashingtypes.NewQueryClient(conn),
		logger: logger.With().Str("module", "slashing").Str("protocol", "grpc").Logger(),
	}
}

func (c *SlashingGRPCClient) Params(ctx context.Context, req *slashingtypes.QueryParamsRequest, opts ...config.CallOpt) (*slashingtypes.QueryParamsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Params, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling SlashingGRPCClient.Params")
	}
	return resp, nil
}

func (c *SlashingGRPCClient) SigningInfo(ctx context.Context, req *slashingtypes.QuerySigningInfoRequest, opts ...config.CallOpt) (*slashingtypes.QuerySigningInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.SigningInfo, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling SlashingGRPCClient.SigningInfo")
	}
	return resp, nil
}

func (c *SlashingGRPCClient) SigningInfos(ctx context.Context, req *slashingtypes.QuerySigningInfosRequest, opts ...config.CallOpt) (*slashingtypes.QuerySigningInfosResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.SigningInfos, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling SlashingGRPCClient.SigningInfos")
	}
	return resp, nil
}
