package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"

	"github.com/allora-network/allora-sdk-go/config"
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

func (c *TendermintRESTClient) ABCIQuery(ctx context.Context, req *cmtservice.ABCIQueryRequest, opts ...config.CallOpt) (*cmtservice.ABCIQueryResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &cmtservice.ABCIQueryResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/base/tendermint/v1beta1/abci_query",
		nil, []string{"data", "path", "height", "prove"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling TendermintRESTClient.ABCIQuery")
	}
	return resp, nil
}

func (c *TendermintRESTClient) GetBlockByHeight(ctx context.Context, req *cmtservice.GetBlockByHeightRequest, opts ...config.CallOpt) (*cmtservice.GetBlockByHeightResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &cmtservice.GetBlockByHeightResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/base/tendermint/v1beta1/blocks/{height}",
		[]string{"height"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling TendermintRESTClient.GetBlockByHeight")
	}
	return resp, nil
}

func (c *TendermintRESTClient) GetLatestBlock(ctx context.Context, req *cmtservice.GetLatestBlockRequest, opts ...config.CallOpt) (*cmtservice.GetLatestBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &cmtservice.GetLatestBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/base/tendermint/v1beta1/blocks/latest",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling TendermintRESTClient.GetLatestBlock")
	}
	return resp, nil
}

func (c *TendermintRESTClient) GetLatestValidatorSet(ctx context.Context, req *cmtservice.GetLatestValidatorSetRequest, opts ...config.CallOpt) (*cmtservice.GetLatestValidatorSetResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &cmtservice.GetLatestValidatorSetResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/base/tendermint/v1beta1/validatorsets/latest",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling TendermintRESTClient.GetLatestValidatorSet")
	}
	return resp, nil
}

func (c *TendermintRESTClient) GetNodeInfo(ctx context.Context, req *cmtservice.GetNodeInfoRequest, opts ...config.CallOpt) (*cmtservice.GetNodeInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &cmtservice.GetNodeInfoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/base/tendermint/v1beta1/node_info",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling TendermintRESTClient.GetNodeInfo")
	}
	return resp, nil
}

func (c *TendermintRESTClient) GetSyncing(ctx context.Context, req *cmtservice.GetSyncingRequest, opts ...config.CallOpt) (*cmtservice.GetSyncingResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &cmtservice.GetSyncingResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/base/tendermint/v1beta1/syncing",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling TendermintRESTClient.GetSyncing")
	}
	return resp, nil
}

func (c *TendermintRESTClient) GetValidatorSetByHeight(ctx context.Context, req *cmtservice.GetValidatorSetByHeightRequest, opts ...config.CallOpt) (*cmtservice.GetValidatorSetByHeightResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &cmtservice.GetValidatorSetByHeightResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/base/tendermint/v1beta1/validatorsets/{height}",
		[]string{"height"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling TendermintRESTClient.GetValidatorSetByHeight")
	}
	return resp, nil
}
