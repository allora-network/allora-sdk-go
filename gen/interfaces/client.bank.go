package interfaces

import (
	"context"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/allora-network/allora-sdk-go/config"
)

type BankClient interface {
	Balance(ctx context.Context, req *banktypes.QueryBalanceRequest, opts ...config.CallOpt) (*banktypes.QueryBalanceResponse, error)
	AllBalances(ctx context.Context, req *banktypes.QueryAllBalancesRequest, opts ...config.CallOpt) (*banktypes.QueryAllBalancesResponse, error)
	SpendableBalances(ctx context.Context, req *banktypes.QuerySpendableBalancesRequest, opts ...config.CallOpt) (*banktypes.QuerySpendableBalancesResponse, error)
	SpendableBalanceByDenom(ctx context.Context, req *banktypes.QuerySpendableBalanceByDenomRequest, opts ...config.CallOpt) (*banktypes.QuerySpendableBalanceByDenomResponse, error)
	TotalSupply(ctx context.Context, req *banktypes.QueryTotalSupplyRequest, opts ...config.CallOpt) (*banktypes.QueryTotalSupplyResponse, error)
	SupplyOf(ctx context.Context, req *banktypes.QuerySupplyOfRequest, opts ...config.CallOpt) (*banktypes.QuerySupplyOfResponse, error)
	Params(ctx context.Context, req *banktypes.QueryParamsRequest, opts ...config.CallOpt) (*banktypes.QueryParamsResponse, error)
	DenomMetadata(ctx context.Context, req *banktypes.QueryDenomMetadataRequest, opts ...config.CallOpt) (*banktypes.QueryDenomMetadataResponse, error)
	DenomMetadataByQueryString(ctx context.Context, req *banktypes.QueryDenomMetadataByQueryStringRequest, opts ...config.CallOpt) (*banktypes.QueryDenomMetadataByQueryStringResponse, error)
	DenomsMetadata(ctx context.Context, req *banktypes.QueryDenomsMetadataRequest, opts ...config.CallOpt) (*banktypes.QueryDenomsMetadataResponse, error)
	DenomOwners(ctx context.Context, req *banktypes.QueryDenomOwnersRequest, opts ...config.CallOpt) (*banktypes.QueryDenomOwnersResponse, error)
	DenomOwnersByQuery(ctx context.Context, req *banktypes.QueryDenomOwnersByQueryRequest, opts ...config.CallOpt) (*banktypes.QueryDenomOwnersByQueryResponse, error)
	SendEnabled(ctx context.Context, req *banktypes.QuerySendEnabledRequest, opts ...config.CallOpt) (*banktypes.QuerySendEnabledResponse, error)
}
