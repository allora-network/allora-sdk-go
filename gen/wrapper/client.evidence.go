package wrapper

import (
	"context"

	evidencetypes "cosmossdk.io/x/evidence/types"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/pool"
)

// EvidenceClientWrapper wraps the evidence module with pool management and retry logic
type EvidenceClientWrapper struct {
	poolManager *pool.ClientPoolManager
	logger      zerolog.Logger
}

// NewEvidenceClientWrapper creates a new evidence client wrapper
func NewEvidenceClientWrapper(poolManager *pool.ClientPoolManager, logger zerolog.Logger) *EvidenceClientWrapper {
	return &EvidenceClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "evidence").Logger(),
	}
}

func (c *EvidenceClientWrapper) Evidence(ctx context.Context, req *evidencetypes.QueryEvidenceRequest) (*evidencetypes.QueryEvidenceResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*evidencetypes.QueryEvidenceResponse, error) {
		return client.Evidence().Evidence(ctx, req)
	})
}

func (c *EvidenceClientWrapper) AllEvidence(ctx context.Context, req *evidencetypes.QueryAllEvidenceRequest) (*evidencetypes.QueryAllEvidenceResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*evidencetypes.QueryAllEvidenceResponse, error) {
		return client.Evidence().AllEvidence(ctx, req)
	})
}
