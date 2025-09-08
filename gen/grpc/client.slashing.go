package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

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

func (c *SlashingGRPCClient) Params(ctx context.Context, req *slashingtypes.QueryParamsRequest) (*slashingtypes.QueryParamsResponse, error) {
	resp, err := c.client.Params(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Params")
}

func (c *SlashingGRPCClient) SigningInfo(ctx context.Context, req *slashingtypes.QuerySigningInfoRequest) (*slashingtypes.QuerySigningInfoResponse, error) {
	resp, err := c.client.SigningInfo(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.SigningInfo")
}

func (c *SlashingGRPCClient) SigningInfos(ctx context.Context, req *slashingtypes.QuerySigningInfosRequest) (*slashingtypes.QuerySigningInfosResponse, error) {
	resp, err := c.client.SigningInfos(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.SigningInfos")
}
