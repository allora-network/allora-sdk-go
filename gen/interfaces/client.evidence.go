package interfaces

import (
	"context"

	evidencetypes "cosmossdk.io/x/evidence/types"
)

type EvidenceClient interface {
	Evidence(ctx context.Context, req *evidencetypes.QueryEvidenceRequest) (*evidencetypes.QueryEvidenceResponse, error)
	AllEvidence(ctx context.Context, req *evidencetypes.QueryAllEvidenceRequest) (*evidencetypes.QueryAllEvidenceResponse, error)
}
