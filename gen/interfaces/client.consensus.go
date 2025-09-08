package interfaces

import (
	"context"

	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
)

type ConsensusClient interface {
	Params(ctx context.Context, req *consensustypes.QueryParamsRequest) (*consensustypes.QueryParamsResponse, error)
}
