package interfaces

import (
	"context"

	cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
)

type TendermintClient interface {
	GetNodeInfo(ctx context.Context, req *cmtservice.GetNodeInfoRequest) (*cmtservice.GetNodeInfoResponse, error)
	GetSyncing(ctx context.Context, req *cmtservice.GetSyncingRequest) (*cmtservice.GetSyncingResponse, error)
	GetLatestBlock(ctx context.Context, req *cmtservice.GetLatestBlockRequest) (*cmtservice.GetLatestBlockResponse, error)
	GetBlockByHeight(ctx context.Context, req *cmtservice.GetBlockByHeightRequest) (*cmtservice.GetBlockByHeightResponse, error)
	GetLatestValidatorSet(ctx context.Context, req *cmtservice.GetLatestValidatorSetRequest) (*cmtservice.GetLatestValidatorSetResponse, error)
	GetValidatorSetByHeight(ctx context.Context, req *cmtservice.GetValidatorSetByHeightRequest) (*cmtservice.GetValidatorSetByHeightResponse, error)
	ABCIQuery(ctx context.Context, req *cmtservice.ABCIQueryRequest) (*cmtservice.ABCIQueryResponse, error)
}
