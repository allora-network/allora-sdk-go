package interfaces

import (
	"context"

	minttypes "github.com/allora-network/allora-chain/x/mint/types"

	"github.com/allora-network/allora-sdk-go/config"
)

type MintClient interface {
	EmissionInfo(ctx context.Context, req *minttypes.QueryServiceEmissionInfoRequest, opts ...config.CallOpt) (*minttypes.QueryServiceEmissionInfoResponse, error)
	Inflation(ctx context.Context, req *minttypes.QueryServiceInflationRequest, opts ...config.CallOpt) (*minttypes.QueryServiceInflationResponse, error)
	Params(ctx context.Context, req *minttypes.QueryServiceParamsRequest, opts ...config.CallOpt) (*minttypes.QueryServiceParamsResponse, error)
}
