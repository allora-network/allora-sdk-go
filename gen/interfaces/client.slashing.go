package interfaces

import (
	"context"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	"github.com/allora-network/allora-sdk-go/config"
)

type SlashingClient interface {
	Params(ctx context.Context, req *slashingtypes.QueryParamsRequest, opts ...config.CallOpt) (*slashingtypes.QueryParamsResponse, error)
	SigningInfo(ctx context.Context, req *slashingtypes.QuerySigningInfoRequest, opts ...config.CallOpt) (*slashingtypes.QuerySigningInfoResponse, error)
	SigningInfos(ctx context.Context, req *slashingtypes.QuerySigningInfosRequest, opts ...config.CallOpt) (*slashingtypes.QuerySigningInfosResponse, error)
}
