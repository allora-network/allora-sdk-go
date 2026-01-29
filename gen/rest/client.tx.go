package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	tx "github.com/cosmos/cosmos-sdk/types/tx"

	"github.com/allora-network/allora-sdk-go/config"
)

// TxRESTClient provides REST access to the tx module
type TxRESTClient struct {
	*RESTClientCore
	logger zerolog.Logger
}

// NewTxRESTClient creates a new tx REST client
func NewTxRESTClient(core *RESTClientCore, logger zerolog.Logger) *TxRESTClient {
	return &TxRESTClient{
		RESTClientCore: core,
		logger:         logger.With().Str("module", "tx").Str("protocol", "rest").Logger(),
	}
}

func (c *TxRESTClient) BroadcastTx(ctx context.Context, req *tx.BroadcastTxRequest, opts ...config.CallOpt) (*tx.BroadcastTxResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &tx.BroadcastTxResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"POST", "/cosmos/tx/v1beta1/txs",
		nil, []string{"tx_bytes", "mode"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling TxRESTClient.BroadcastTx")
	}
	return resp, nil
}

func (c *TxRESTClient) GetBlockWithTxs(ctx context.Context, req *tx.GetBlockWithTxsRequest, opts ...config.CallOpt) (*tx.GetBlockWithTxsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &tx.GetBlockWithTxsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/tx/v1beta1/txs/block/{height}",
		[]string{"height"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling TxRESTClient.GetBlockWithTxs")
	}
	return resp, nil
}

func (c *TxRESTClient) GetTx(ctx context.Context, req *tx.GetTxRequest, opts ...config.CallOpt) (*tx.GetTxResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &tx.GetTxResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/tx/v1beta1/txs/{hash}",
		[]string{"hash"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling TxRESTClient.GetTx")
	}
	return resp, nil
}

func (c *TxRESTClient) GetTxsEvent(ctx context.Context, req *tx.GetTxsEventRequest, opts ...config.CallOpt) (*tx.GetTxsEventResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &tx.GetTxsEventResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/tx/v1beta1/txs",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse", "events", "order_by", "page", "limit", "query"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling TxRESTClient.GetTxsEvent")
	}
	return resp, nil
}

func (c *TxRESTClient) Simulate(ctx context.Context, req *tx.SimulateRequest, opts ...config.CallOpt) (*tx.SimulateResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &tx.SimulateResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"POST", "/cosmos/tx/v1beta1/simulate",
		nil, []string{"tx_bytes"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling TxRESTClient.Simulate")
	}
	return resp, nil
}

func (c *TxRESTClient) TxDecode(ctx context.Context, req *tx.TxDecodeRequest, opts ...config.CallOpt) (*tx.TxDecodeResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &tx.TxDecodeResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"POST", "/cosmos/tx/v1beta1/decode",
		nil, []string{"tx_bytes"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling TxRESTClient.TxDecode")
	}
	return resp, nil
}

func (c *TxRESTClient) TxDecodeAmino(ctx context.Context, req *tx.TxDecodeAminoRequest, opts ...config.CallOpt) (*tx.TxDecodeAminoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &tx.TxDecodeAminoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"POST", "/cosmos/tx/v1beta1/decode/amino",
		nil, []string{"amino_binary"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling TxRESTClient.TxDecodeAmino")
	}
	return resp, nil
}

func (c *TxRESTClient) TxEncode(ctx context.Context, req *tx.TxEncodeRequest, opts ...config.CallOpt) (*tx.TxEncodeResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &tx.TxEncodeResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"POST", "/cosmos/tx/v1beta1/encode",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling TxRESTClient.TxEncode")
	}
	return resp, nil
}

func (c *TxRESTClient) TxEncodeAmino(ctx context.Context, req *tx.TxEncodeAminoRequest, opts ...config.CallOpt) (*tx.TxEncodeAminoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &tx.TxEncodeAminoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"POST", "/cosmos/tx/v1beta1/encode/amino",
		nil, []string{"amino_json"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling TxRESTClient.TxEncodeAmino")
	}
	return resp, nil
}
