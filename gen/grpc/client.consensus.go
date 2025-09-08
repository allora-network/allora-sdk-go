package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/allora-network/allora-sdk-go/config"

	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// ConsensusGRPCClient provides gRPC access to the consensus module
type ConsensusGRPCClient struct {
	client consensustypes.QueryClient
	logger zerolog.Logger
}

var _ interfaces.ConsensusClient = (*ConsensusGRPCClient)(nil)

// NewConsensusGRPCClient creates a new consensus REST client
func NewConsensusGRPCClient(conn *grpc.ClientConn, logger zerolog.Logger) *ConsensusGRPCClient {
	return &ConsensusGRPCClient{
		client: consensustypes.NewQueryClient(conn),
		logger: logger.With().Str("module", "consensus").Str("protocol", "grpc").Logger(),
	}
}

func (c *ConsensusGRPCClient) Params(ctx context.Context, req *consensustypes.QueryParamsRequest, opts ...config.CallOpt) (*consensustypes.QueryParamsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Params, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling ConsensusGRPCClient.Params")
	}
	return resp, nil
}
