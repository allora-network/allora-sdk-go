package interfaces

import (
	"context"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type AuthClient interface {
	Accounts(ctx context.Context, req *authtypes.QueryAccountsRequest) (*authtypes.QueryAccountsResponse, error)
	Account(ctx context.Context, req *authtypes.QueryAccountRequest) (*authtypes.QueryAccountResponse, error)
	AccountAddressByID(ctx context.Context, req *authtypes.QueryAccountAddressByIDRequest) (*authtypes.QueryAccountAddressByIDResponse, error)
	Params(ctx context.Context, req *authtypes.QueryParamsRequest) (*authtypes.QueryParamsResponse, error)
	ModuleAccounts(ctx context.Context, req *authtypes.QueryModuleAccountsRequest) (*authtypes.QueryModuleAccountsResponse, error)
	ModuleAccountByName(ctx context.Context, req *authtypes.QueryModuleAccountByNameRequest) (*authtypes.QueryModuleAccountByNameResponse, error)
	Bech32Prefix(ctx context.Context, req *authtypes.Bech32PrefixRequest) (*authtypes.Bech32PrefixResponse, error)
	AddressBytesToString(ctx context.Context, req *authtypes.AddressBytesToStringRequest) (*authtypes.AddressBytesToStringResponse, error)
	AddressStringToBytes(ctx context.Context, req *authtypes.AddressStringToBytesRequest) (*authtypes.AddressStringToBytesResponse, error)
	AccountInfo(ctx context.Context, req *authtypes.QueryAccountInfoRequest) (*authtypes.QueryAccountInfoResponse, error)
}
