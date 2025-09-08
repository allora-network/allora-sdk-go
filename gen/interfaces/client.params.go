package interfaces

import (
	"context"

	proposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
)

type ParamsClient interface {
	Params(ctx context.Context, req *proposal.QueryParamsRequest) (*proposal.QueryParamsResponse, error)
	Subspaces(ctx context.Context, req *proposal.QuerySubspacesRequest) (*proposal.QuerySubspacesResponse, error)
}
