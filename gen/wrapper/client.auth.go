package wrapper

import (
	"context"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/pool"
)

// AuthClientWrapper wraps the auth module with pool management and retry logic
type AuthClientWrapper struct {
	poolManager *pool.ClientPoolManager
	logger      zerolog.Logger
}

// NewAuthClientWrapper creates a new auth client wrapper
func NewAuthClientWrapper(poolManager *pool.ClientPoolManager, logger zerolog.Logger) *AuthClientWrapper {
	return &AuthClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "auth").Logger(),
	}
}

func (c *AuthClientWrapper) Accounts(ctx context.Context, req *authtypes.QueryAccountsRequest) (*authtypes.QueryAccountsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*authtypes.QueryAccountsResponse, error) {
		return client.Auth().Accounts(ctx, req)
	})
}

func (c *AuthClientWrapper) Account(ctx context.Context, req *authtypes.QueryAccountRequest) (*authtypes.QueryAccountResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*authtypes.QueryAccountResponse, error) {
		return client.Auth().Account(ctx, req)
	})
}

func (c *AuthClientWrapper) AccountAddressByID(ctx context.Context, req *authtypes.QueryAccountAddressByIDRequest) (*authtypes.QueryAccountAddressByIDResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*authtypes.QueryAccountAddressByIDResponse, error) {
		return client.Auth().AccountAddressByID(ctx, req)
	})
}

func (c *AuthClientWrapper) Params(ctx context.Context, req *authtypes.QueryParamsRequest) (*authtypes.QueryParamsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*authtypes.QueryParamsResponse, error) {
		return client.Auth().Params(ctx, req)
	})
}

func (c *AuthClientWrapper) ModuleAccounts(ctx context.Context, req *authtypes.QueryModuleAccountsRequest) (*authtypes.QueryModuleAccountsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*authtypes.QueryModuleAccountsResponse, error) {
		return client.Auth().ModuleAccounts(ctx, req)
	})
}

func (c *AuthClientWrapper) ModuleAccountByName(ctx context.Context, req *authtypes.QueryModuleAccountByNameRequest) (*authtypes.QueryModuleAccountByNameResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*authtypes.QueryModuleAccountByNameResponse, error) {
		return client.Auth().ModuleAccountByName(ctx, req)
	})
}

func (c *AuthClientWrapper) Bech32Prefix(ctx context.Context, req *authtypes.Bech32PrefixRequest) (*authtypes.Bech32PrefixResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*authtypes.Bech32PrefixResponse, error) {
		return client.Auth().Bech32Prefix(ctx, req)
	})
}

func (c *AuthClientWrapper) AddressBytesToString(ctx context.Context, req *authtypes.AddressBytesToStringRequest) (*authtypes.AddressBytesToStringResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*authtypes.AddressBytesToStringResponse, error) {
		return client.Auth().AddressBytesToString(ctx, req)
	})
}

func (c *AuthClientWrapper) AddressStringToBytes(ctx context.Context, req *authtypes.AddressStringToBytesRequest) (*authtypes.AddressStringToBytesResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*authtypes.AddressStringToBytesResponse, error) {
		return client.Auth().AddressStringToBytes(ctx, req)
	})
}

func (c *AuthClientWrapper) AccountInfo(ctx context.Context, req *authtypes.QueryAccountInfoRequest) (*authtypes.QueryAccountInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*authtypes.QueryAccountInfoResponse, error) {
		return client.Auth().AccountInfo(ctx, req)
	})
}
