package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	evidencetypes "cosmossdk.io/x/evidence/types"
)

// EvidenceRESTClient provides REST access to the evidence module
type EvidenceRESTClient struct {
	*RESTClientCore
	logger zerolog.Logger
}

// NewEvidenceRESTClient creates a new evidence REST client
func NewEvidenceRESTClient(core *RESTClientCore, logger zerolog.Logger) *EvidenceRESTClient {
	return &EvidenceRESTClient{
		RESTClientCore: core,
		logger:         logger.With().Str("module", "evidence").Str("protocol", "rest").Logger(),
	}
}

func (c *EvidenceRESTClient) Evidence(ctx context.Context, req *evidencetypes.QueryEvidenceRequest) (*evidencetypes.QueryEvidenceResponse, error) {
	resp := &evidencetypes.QueryEvidenceResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/evidence/v1beta1/evidence/{hash}",
		[]string{"hash"}, []string{"evidence_hash"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EvidenceRESTClient.Evidence")
}

func (c *EvidenceRESTClient) AllEvidence(ctx context.Context, req *evidencetypes.QueryAllEvidenceRequest) (*evidencetypes.QueryAllEvidenceResponse, error) {
	resp := &evidencetypes.QueryAllEvidenceResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/evidence/v1beta1/evidence",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EvidenceRESTClient.AllEvidence")
}
