package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	evidencetypes "cosmossdk.io/x/evidence/types"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// EvidenceGRPCClient provides gRPC access to the evidence module
type EvidenceGRPCClient struct {
	client evidencetypes.QueryClient
	logger zerolog.Logger
}

var _ interfaces.EvidenceClient = (*EvidenceGRPCClient)(nil)

// NewEvidenceGRPCClient creates a new evidence REST client
func NewEvidenceGRPCClient(conn *grpc.ClientConn, logger zerolog.Logger) *EvidenceGRPCClient {
	return &EvidenceGRPCClient{
		client: evidencetypes.NewQueryClient(conn),
		logger: logger.With().Str("module", "evidence").Str("protocol", "grpc").Logger(),
	}
}

func (c *EvidenceGRPCClient) Evidence(ctx context.Context, req *evidencetypes.QueryEvidenceRequest) (*evidencetypes.QueryEvidenceResponse, error) {
	resp, err := c.client.Evidence(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Evidence")
}

func (c *EvidenceGRPCClient) AllEvidence(ctx context.Context, req *evidencetypes.QueryAllEvidenceRequest) (*evidencetypes.QueryAllEvidenceResponse, error) {
	resp, err := c.client.AllEvidence(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.AllEvidence")
}
