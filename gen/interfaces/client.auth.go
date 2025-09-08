package interfaces

import (
	"context"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/allora-network/allora-sdk-go/config"
)

type AuthClient interface {
	Accounts(ctx context.Context, req *authtypes.QueryAccountsRequest, opts ...config.CallOpt) (*authtypes.QueryAccountsResponse, error)
	Account(ctx context.Context, req *authtypes.QueryAccountRequest, opts ...config.CallOpt) (*authtypes.QueryAccountResponse, error)
	AccountAddressByID(ctx context.Context, req *authtypes.QueryAccountAddressByIDRequest, opts ...config.CallOpt) (*authtypes.QueryAccountAddressByIDResponse, error)
	Params(ctx context.Context, req *authtypes.QueryParamsRequest, opts ...config.CallOpt) (*authtypes.QueryParamsResponse, error)
	ModuleAccounts(ctx context.Context, req *authtypes.QueryModuleAccountsRequest, opts ...config.CallOpt) (*authtypes.QueryModuleAccountsResponse, error)
	ModuleAccountByName(ctx context.Context, req *authtypes.QueryModuleAccountByNameRequest, opts ...config.CallOpt) (*authtypes.QueryModuleAccountByNameResponse, error)
	Bech32Prefix(ctx context.Context, req *authtypes.Bech32PrefixRequest, opts ...config.CallOpt) (*authtypes.Bech32PrefixResponse, error)
	AddressBytesToString(ctx context.Context, req *authtypes.AddressBytesToStringRequest, opts ...config.CallOpt) (*authtypes.AddressBytesToStringResponse, error)
	AddressStringToBytes(ctx context.Context, req *authtypes.AddressStringToBytesRequest, opts ...config.CallOpt) (*authtypes.AddressStringToBytesResponse, error)
	AccountInfo(ctx context.Context, req *authtypes.QueryAccountInfoRequest, opts ...config.CallOpt) (*authtypes.QueryAccountInfoResponse, error)
}
