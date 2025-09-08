package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

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

func (c *TxGRPCClient) Simulate(ctx context.Context, req *tx.SimulateRequest) (*tx.SimulateResponse, error) {
	resp, err := c.client.Simulate(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Simulate")
}

func (c *TxGRPCClient) GetTx(ctx context.Context, req *tx.GetTxRequest) (*tx.GetTxResponse, error) {
	resp, err := c.client.GetTx(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetTx")
}

func (c *TxGRPCClient) BroadcastTx(ctx context.Context, req *tx.BroadcastTxRequest) (*tx.BroadcastTxResponse, error) {
	resp, err := c.client.BroadcastTx(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.BroadcastTx")
}

func (c *TxGRPCClient) GetTxsEvent(ctx context.Context, req *tx.GetTxsEventRequest) (*tx.GetTxsEventResponse, error) {
	resp, err := c.client.GetTxsEvent(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetTxsEvent")
}

func (c *TxGRPCClient) GetBlockWithTxs(ctx context.Context, req *tx.GetBlockWithTxsRequest) (*tx.GetBlockWithTxsResponse, error) {
	resp, err := c.client.GetBlockWithTxs(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetBlockWithTxs")
}

func (c *TxGRPCClient) TxDecode(ctx context.Context, req *tx.TxDecodeRequest) (*tx.TxDecodeResponse, error) {
	resp, err := c.client.TxDecode(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.TxDecode")
}

func (c *TxGRPCClient) TxEncode(ctx context.Context, req *tx.TxEncodeRequest) (*tx.TxEncodeResponse, error) {
	resp, err := c.client.TxEncode(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.TxEncode")
}

func (c *TxGRPCClient) TxEncodeAmino(ctx context.Context, req *tx.TxEncodeAminoRequest) (*tx.TxEncodeAminoResponse, error) {
	resp, err := c.client.TxEncodeAmino(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.TxEncodeAmino")
}

func (c *TxGRPCClient) TxDecodeAmino(ctx context.Context, req *tx.TxDecodeAminoRequest) (*tx.TxDecodeAminoResponse, error) {
	resp, err := c.client.TxDecodeAmino(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.TxDecodeAmino")
}
