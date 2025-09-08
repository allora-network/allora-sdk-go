package interfaces

import (
	"context"

	minttypes "github.com/allora-network/allora-chain/x/mint/types"
)

type MintClient interface {
	Params(ctx context.Context, req *minttypes.QueryServiceParamsRequest) (*minttypes.QueryServiceParamsResponse, error)
	Inflation(ctx context.Context, req *minttypes.QueryServiceInflationRequest) (*minttypes.QueryServiceInflationResponse, error)
	EmissionInfo(ctx context.Context, req *minttypes.QueryServiceEmissionInfoRequest) (*minttypes.QueryServiceEmissionInfoResponse, error)
}
