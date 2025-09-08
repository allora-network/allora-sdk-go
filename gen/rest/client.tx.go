package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	tx "github.com/cosmos/cosmos-sdk/types/tx"
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

func (c *TxRESTClient) Simulate(ctx context.Context, req *tx.SimulateRequest) (*tx.SimulateResponse, error) {
	resp := &tx.SimulateResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"POST", "/cosmos/tx/v1beta1/simulate",
		nil, []string{"tx_bytes"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling TxRESTClient.Simulate")
}

func (c *TxRESTClient) GetTx(ctx context.Context, req *tx.GetTxRequest) (*tx.GetTxResponse, error) {
	resp := &tx.GetTxResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/tx/v1beta1/txs/{hash}",
		[]string{"hash"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling TxRESTClient.GetTx")
}

func (c *TxRESTClient) BroadcastTx(ctx context.Context, req *tx.BroadcastTxRequest) (*tx.BroadcastTxResponse, error) {
	resp := &tx.BroadcastTxResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"POST", "/cosmos/tx/v1beta1/txs",
		nil, []string{"tx_bytes", "mode"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling TxRESTClient.BroadcastTx")
}

func (c *TxRESTClient) GetTxsEvent(ctx context.Context, req *tx.GetTxsEventRequest) (*tx.GetTxsEventResponse, error) {
	resp := &tx.GetTxsEventResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/tx/v1beta1/txs",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse", "events", "order_by", "page", "limit", "query"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling TxRESTClient.GetTxsEvent")
}

func (c *TxRESTClient) GetBlockWithTxs(ctx context.Context, req *tx.GetBlockWithTxsRequest) (*tx.GetBlockWithTxsResponse, error) {
	resp := &tx.GetBlockWithTxsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/tx/v1beta1/txs/block/{height}",
		[]string{"height"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling TxRESTClient.GetBlockWithTxs")
}

func (c *TxRESTClient) TxDecode(ctx context.Context, req *tx.TxDecodeRequest) (*tx.TxDecodeResponse, error) {
	resp := &tx.TxDecodeResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"POST", "/cosmos/tx/v1beta1/decode",
		nil, []string{"tx_bytes"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling TxRESTClient.TxDecode")
}

func (c *TxRESTClient) TxEncode(ctx context.Context, req *tx.TxEncodeRequest) (*tx.TxEncodeResponse, error) {
	resp := &tx.TxEncodeResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"POST", "/cosmos/tx/v1beta1/encode",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling TxRESTClient.TxEncode")
}

func (c *TxRESTClient) TxEncodeAmino(ctx context.Context, req *tx.TxEncodeAminoRequest) (*tx.TxEncodeAminoResponse, error) {
	resp := &tx.TxEncodeAminoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"POST", "/cosmos/tx/v1beta1/encode/amino",
		nil, []string{"amino_json"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling TxRESTClient.TxEncodeAmino")
}

func (c *TxRESTClient) TxDecodeAmino(ctx context.Context, req *tx.TxDecodeAminoRequest) (*tx.TxDecodeAminoResponse, error) {
	resp := &tx.TxDecodeAminoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"POST", "/cosmos/tx/v1beta1/decode/amino",
		nil, []string{"amino_binary"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling TxRESTClient.TxDecodeAmino")
}
