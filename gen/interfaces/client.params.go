package interfaces

import (
	"context"

	proposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"

	"github.com/allora-network/allora-sdk-go/config"
)

type ParamsClient interface {
	Params(ctx context.Context, req *proposal.QueryParamsRequest, opts ...config.CallOpt) (*proposal.QueryParamsResponse, error)
	Subspaces(ctx context.Context, req *proposal.QuerySubspacesRequest, opts ...config.CallOpt) (*proposal.QuerySubspacesResponse, error)
}
