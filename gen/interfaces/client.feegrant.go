package interfaces

import (
	"context"

	feegrant "cosmossdk.io/x/feegrant"
)

type FeegrantClient interface {
	Allowance(ctx context.Context, req *feegrant.QueryAllowanceRequest) (*feegrant.QueryAllowanceResponse, error)
	Allowances(ctx context.Context, req *feegrant.QueryAllowancesRequest) (*feegrant.QueryAllowancesResponse, error)
	AllowancesByGranter(ctx context.Context, req *feegrant.QueryAllowancesByGranterRequest) (*feegrant.QueryAllowancesByGranterResponse, error)
}
