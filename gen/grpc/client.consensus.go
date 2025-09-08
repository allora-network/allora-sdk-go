package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

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

func (c *ConsensusGRPCClient) Params(ctx context.Context, req *consensustypes.QueryParamsRequest) (*consensustypes.QueryParamsResponse, error) {
	resp, err := c.client.Params(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Params")
}
