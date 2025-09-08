package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
)

// TendermintRESTClient provides REST access to the tendermint module
type TendermintRESTClient struct {
	*RESTClientCore
	logger zerolog.Logger
}

// NewTendermintRESTClient creates a new tendermint REST client
func NewTendermintRESTClient(core *RESTClientCore, logger zerolog.Logger) *TendermintRESTClient {
	return &TendermintRESTClient{
		RESTClientCore: core,
		logger:         logger.With().Str("module", "tendermint").Str("protocol", "rest").Logger(),
	}
}

func (c *TendermintRESTClient) GetNodeInfo(ctx context.Context, req *cmtservice.GetNodeInfoRequest) (*cmtservice.GetNodeInfoResponse, error) {
	resp := &cmtservice.GetNodeInfoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/base/tendermint/v1beta1/node_info",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling TendermintRESTClient.GetNodeInfo")
}

func (c *TendermintRESTClient) GetSyncing(ctx context.Context, req *cmtservice.GetSyncingRequest) (*cmtservice.GetSyncingResponse, error) {
	resp := &cmtservice.GetSyncingResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/base/tendermint/v1beta1/syncing",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling TendermintRESTClient.GetSyncing")
}

func (c *TendermintRESTClient) GetLatestBlock(ctx context.Context, req *cmtservice.GetLatestBlockRequest) (*cmtservice.GetLatestBlockResponse, error) {
	resp := &cmtservice.GetLatestBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/base/tendermint/v1beta1/blocks/latest",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling TendermintRESTClient.GetLatestBlock")
}

func (c *TendermintRESTClient) GetBlockByHeight(ctx context.Context, req *cmtservice.GetBlockByHeightRequest) (*cmtservice.GetBlockByHeightResponse, error) {
	resp := &cmtservice.GetBlockByHeightResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/base/tendermint/v1beta1/blocks/{height}",
		[]string{"height"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling TendermintRESTClient.GetBlockByHeight")
}

func (c *TendermintRESTClient) GetLatestValidatorSet(ctx context.Context, req *cmtservice.GetLatestValidatorSetRequest) (*cmtservice.GetLatestValidatorSetResponse, error) {
	resp := &cmtservice.GetLatestValidatorSetResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/base/tendermint/v1beta1/validatorsets/latest",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling TendermintRESTClient.GetLatestValidatorSet")
}

func (c *TendermintRESTClient) GetValidatorSetByHeight(ctx context.Context, req *cmtservice.GetValidatorSetByHeightRequest) (*cmtservice.GetValidatorSetByHeightResponse, error) {
	resp := &cmtservice.GetValidatorSetByHeightResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/base/tendermint/v1beta1/validatorsets/{height}",
		[]string{"height"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling TendermintRESTClient.GetValidatorSetByHeight")
}

func (c *TendermintRESTClient) ABCIQuery(ctx context.Context, req *cmtservice.ABCIQueryRequest) (*cmtservice.ABCIQueryResponse, error) {
	resp := &cmtservice.ABCIQueryResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/base/tendermint/v1beta1/abci_query",
		nil, []string{"data", "path", "height", "prove"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling TendermintRESTClient.ABCIQuery")
}
