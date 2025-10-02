package interfaces

import (
	"context"

	authz "github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/allora-network/allora-sdk-go/config"
)

type AuthzClient interface {
	Grants(ctx context.Context, req *authz.QueryGrantsRequest, opts ...config.CallOpt) (*authz.QueryGrantsResponse, error)
	GranterGrants(ctx context.Context, req *authz.QueryGranterGrantsRequest, opts ...config.CallOpt) (*authz.QueryGranterGrantsResponse, error)
	GranteeGrants(ctx context.Context, req *authz.QueryGranteeGrantsRequest, opts ...config.CallOpt) (*authz.QueryGranteeGrantsResponse, error)
}
