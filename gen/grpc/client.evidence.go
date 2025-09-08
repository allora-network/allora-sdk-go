package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/allora-network/allora-sdk-go/config"

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

func (c *EvidenceGRPCClient) Evidence(ctx context.Context, req *evidencetypes.QueryEvidenceRequest, opts ...config.CallOpt) (*evidencetypes.QueryEvidenceResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Evidence, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EvidenceGRPCClient.Evidence")
	}
	return resp, nil
}

func (c *EvidenceGRPCClient) AllEvidence(ctx context.Context, req *evidencetypes.QueryAllEvidenceRequest, opts ...config.CallOpt) (*evidencetypes.QueryAllEvidenceResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.AllEvidence, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EvidenceGRPCClient.AllEvidence")
	}
	return resp, nil
}
