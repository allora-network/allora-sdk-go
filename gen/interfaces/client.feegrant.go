package interfaces

import (
	"context"

	feegrant "cosmossdk.io/x/feegrant"

	"github.com/allora-network/allora-sdk-go/config"
)

type FeegrantClient interface {
	Allowance(ctx context.Context, req *feegrant.QueryAllowanceRequest, opts ...config.CallOpt) (*feegrant.QueryAllowanceResponse, error)
	Allowances(ctx context.Context, req *feegrant.QueryAllowancesRequest, opts ...config.CallOpt) (*feegrant.QueryAllowancesResponse, error)
	AllowancesByGranter(ctx context.Context, req *feegrant.QueryAllowancesByGranterRequest, opts ...config.CallOpt) (*feegrant.QueryAllowancesByGranterResponse, error)
}
