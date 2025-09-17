package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	minttypes "github.com/allora-network/allora-chain/x/mint/types"

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

func (c *MintGRPCClient) Params(ctx context.Context, req *minttypes.QueryServiceParamsRequest) (*minttypes.QueryServiceParamsResponse, error) {
	resp, err := c.client.Params(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Params")
}

func (c *MintGRPCClient) Inflation(ctx context.Context, req *minttypes.QueryServiceInflationRequest) (*minttypes.QueryServiceInflationResponse, error) {
	resp, err := c.client.Inflation(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Inflation")
}

func (c *MintGRPCClient) EmissionInfo(ctx context.Context, req *minttypes.QueryServiceEmissionInfoRequest) (*minttypes.QueryServiceEmissionInfoResponse, error) {
	resp, err := c.client.EmissionInfo(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.EmissionInfo")
}
