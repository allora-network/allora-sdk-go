package interfaces

import (
	"context"

	evidencetypes "cosmossdk.io/x/evidence/types"

	"github.com/allora-network/allora-sdk-go/config"
)

type EvidenceClient interface {
	Evidence(ctx context.Context, req *evidencetypes.QueryEvidenceRequest, opts ...config.CallOpt) (*evidencetypes.QueryEvidenceResponse, error)
	AllEvidence(ctx context.Context, req *evidencetypes.QueryAllEvidenceRequest, opts ...config.CallOpt) (*evidencetypes.QueryAllEvidenceResponse, error)
}
