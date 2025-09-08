package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/allora-network/allora-sdk-go/config"

	tx "github.com/cosmos/cosmos-sdk/types/tx"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// TxGRPCClient provides gRPC access to the tx module
type TxGRPCClient struct {
	client tx.ServiceClient
	logger zerolog.Logger
}

var _ interfaces.TxClient = (*TxGRPCClient)(nil)

// NewTxGRPCClient creates a new tx REST client
func NewTxGRPCClient(conn *grpc.ClientConn, logger zerolog.Logger) *TxGRPCClient {
	return &TxGRPCClient{
		client: tx.NewServiceClient(conn),
		logger: logger.With().Str("module", "tx").Str("protocol", "grpc").Logger(),
	}
}

func (c *TxGRPCClient) Simulate(ctx context.Context, req *tx.SimulateRequest, opts ...config.CallOpt) (*tx.SimulateResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Simulate, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling TxGRPCClient.Simulate")
	}
	return resp, nil
}

func (c *TxGRPCClient) GetTx(ctx context.Context, req *tx.GetTxRequest, opts ...config.CallOpt) (*tx.GetTxResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetTx, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling TxGRPCClient.GetTx")
	}
	return resp, nil
}

func (c *TxGRPCClient) BroadcastTx(ctx context.Context, req *tx.BroadcastTxRequest, opts ...config.CallOpt) (*tx.BroadcastTxResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.BroadcastTx, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling TxGRPCClient.BroadcastTx")
	}
	return resp, nil
}

func (c *TxGRPCClient) GetTxsEvent(ctx context.Context, req *tx.GetTxsEventRequest, opts ...config.CallOpt) (*tx.GetTxsEventResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetTxsEvent, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling TxGRPCClient.GetTxsEvent")
	}
	return resp, nil
}

func (c *TxGRPCClient) GetBlockWithTxs(ctx context.Context, req *tx.GetBlockWithTxsRequest, opts ...config.CallOpt) (*tx.GetBlockWithTxsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetBlockWithTxs, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling TxGRPCClient.GetBlockWithTxs")
	}
	return resp, nil
}

func (c *TxGRPCClient) TxDecode(ctx context.Context, req *tx.TxDecodeRequest, opts ...config.CallOpt) (*tx.TxDecodeResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.TxDecode, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling TxGRPCClient.TxDecode")
	}
	return resp, nil
}

func (c *TxGRPCClient) TxEncode(ctx context.Context, req *tx.TxEncodeRequest, opts ...config.CallOpt) (*tx.TxEncodeResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.TxEncode, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling TxGRPCClient.TxEncode")
	}
	return resp, nil
}

func (c *TxGRPCClient) TxEncodeAmino(ctx context.Context, req *tx.TxEncodeAminoRequest, opts ...config.CallOpt) (*tx.TxEncodeAminoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.TxEncodeAmino, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling TxGRPCClient.TxEncodeAmino")
	}
	return resp, nil
}

func (c *TxGRPCClient) TxDecodeAmino(ctx context.Context, req *tx.TxDecodeAminoRequest, opts ...config.CallOpt) (*tx.TxDecodeAminoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.TxDecodeAmino, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling TxGRPCClient.TxDecodeAmino")
	}
	return resp, nil
}
