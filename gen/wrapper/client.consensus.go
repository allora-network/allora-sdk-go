package wrapper

import (
	"context"

	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/pool"
)

// ConsensusClientWrapper wraps the consensus module with pool management and retry logic
type ConsensusClientWrapper struct {
	poolManager *pool.ClientPoolManager
	logger      zerolog.Logger
}

// NewConsensusClientWrapper creates a new consensus client wrapper
func NewConsensusClientWrapper(poolManager *pool.ClientPoolManager, logger zerolog.Logger) *ConsensusClientWrapper {
	return &ConsensusClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "consensus").Logger(),
	}
}

func (c *ConsensusClientWrapper) Params(ctx context.Context, req *consensustypes.QueryParamsRequest) (*consensustypes.QueryParamsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*consensustypes.QueryParamsResponse, error) {
		return client.Consensus().Params(ctx, req)
	})
}
