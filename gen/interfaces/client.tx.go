package interfaces

import (
	"context"

	tx "github.com/cosmos/cosmos-sdk/types/tx"

	"github.com/allora-network/allora-sdk-go/config"
)

type TxClient interface {
	Simulate(ctx context.Context, req *tx.SimulateRequest, opts ...config.CallOpt) (*tx.SimulateResponse, error)
	GetTx(ctx context.Context, req *tx.GetTxRequest, opts ...config.CallOpt) (*tx.GetTxResponse, error)
	BroadcastTx(ctx context.Context, req *tx.BroadcastTxRequest, opts ...config.CallOpt) (*tx.BroadcastTxResponse, error)
	GetTxsEvent(ctx context.Context, req *tx.GetTxsEventRequest, opts ...config.CallOpt) (*tx.GetTxsEventResponse, error)
	GetBlockWithTxs(ctx context.Context, req *tx.GetBlockWithTxsRequest, opts ...config.CallOpt) (*tx.GetBlockWithTxsResponse, error)
	TxDecode(ctx context.Context, req *tx.TxDecodeRequest, opts ...config.CallOpt) (*tx.TxDecodeResponse, error)
	TxEncode(ctx context.Context, req *tx.TxEncodeRequest, opts ...config.CallOpt) (*tx.TxEncodeResponse, error)
	TxEncodeAmino(ctx context.Context, req *tx.TxEncodeAminoRequest, opts ...config.CallOpt) (*tx.TxEncodeAminoResponse, error)
	TxDecodeAmino(ctx context.Context, req *tx.TxDecodeAminoRequest, opts ...config.CallOpt) (*tx.TxDecodeAminoResponse, error)
}
