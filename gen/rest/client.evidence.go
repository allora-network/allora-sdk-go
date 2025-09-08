package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	evidencetypes "cosmossdk.io/x/evidence/types"

	"github.com/allora-network/allora-sdk-go/config"
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

func (c *EvidenceRESTClient) Evidence(ctx context.Context, req *evidencetypes.QueryEvidenceRequest, opts ...config.CallOpt) (*evidencetypes.QueryEvidenceResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &evidencetypes.QueryEvidenceResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/evidence/v1beta1/evidence/{hash}",
		[]string{"hash"}, []string{"evidence_hash"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EvidenceRESTClient.Evidence")
	}
	return resp, nil
}

func (c *EvidenceRESTClient) AllEvidence(ctx context.Context, req *evidencetypes.QueryAllEvidenceRequest, opts ...config.CallOpt) (*evidencetypes.QueryAllEvidenceResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &evidencetypes.QueryAllEvidenceResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/evidence/v1beta1/evidence",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EvidenceRESTClient.AllEvidence")
	}
	return resp, nil
}
