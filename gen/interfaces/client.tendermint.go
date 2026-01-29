package interfaces

import (
	"context"

	cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"

	"github.com/allora-network/allora-sdk-go/config"
)

type TendermintClient interface {
	ABCIQuery(ctx context.Context, req *cmtservice.ABCIQueryRequest, opts ...config.CallOpt) (*cmtservice.ABCIQueryResponse, error)
	GetBlockByHeight(ctx context.Context, req *cmtservice.GetBlockByHeightRequest, opts ...config.CallOpt) (*cmtservice.GetBlockByHeightResponse, error)
	GetLatestBlock(ctx context.Context, req *cmtservice.GetLatestBlockRequest, opts ...config.CallOpt) (*cmtservice.GetLatestBlockResponse, error)
	GetLatestValidatorSet(ctx context.Context, req *cmtservice.GetLatestValidatorSetRequest, opts ...config.CallOpt) (*cmtservice.GetLatestValidatorSetResponse, error)
	GetNodeInfo(ctx context.Context, req *cmtservice.GetNodeInfoRequest, opts ...config.CallOpt) (*cmtservice.GetNodeInfoResponse, error)
	GetSyncing(ctx context.Context, req *cmtservice.GetSyncingRequest, opts ...config.CallOpt) (*cmtservice.GetSyncingResponse, error)
	GetValidatorSetByHeight(ctx context.Context, req *cmtservice.GetValidatorSetByHeightRequest, opts ...config.CallOpt) (*cmtservice.GetValidatorSetByHeightResponse, error)
}
