package interfaces

import (
	"context"

	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"

	"github.com/allora-network/allora-sdk-go/config"
)

type ConsensusClient interface {
	Params(ctx context.Context, req *consensustypes.QueryParamsRequest, opts ...config.CallOpt) (*consensustypes.QueryParamsResponse, error)
}
