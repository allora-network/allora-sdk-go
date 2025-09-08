package interfaces

import (
	"context"

	tx "github.com/cosmos/cosmos-sdk/types/tx"
)

type TxClient interface {
	Simulate(ctx context.Context, req *tx.SimulateRequest) (*tx.SimulateResponse, error)
	GetTx(ctx context.Context, req *tx.GetTxRequest) (*tx.GetTxResponse, error)
	BroadcastTx(ctx context.Context, req *tx.BroadcastTxRequest) (*tx.BroadcastTxResponse, error)
	GetTxsEvent(ctx context.Context, req *tx.GetTxsEventRequest) (*tx.GetTxsEventResponse, error)
	GetBlockWithTxs(ctx context.Context, req *tx.GetBlockWithTxsRequest) (*tx.GetBlockWithTxsResponse, error)
	TxDecode(ctx context.Context, req *tx.TxDecodeRequest) (*tx.TxDecodeResponse, error)
	TxEncode(ctx context.Context, req *tx.TxEncodeRequest) (*tx.TxEncodeResponse, error)
	TxEncodeAmino(ctx context.Context, req *tx.TxEncodeAminoRequest) (*tx.TxEncodeAminoResponse, error)
	TxDecodeAmino(ctx context.Context, req *tx.TxDecodeAminoRequest) (*tx.TxDecodeAminoResponse, error)
}
