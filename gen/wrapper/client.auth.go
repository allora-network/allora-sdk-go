package wrapper

import (
	"context"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/pool"
)

// AuthClientWrapper wraps the auth module with pool management and retry logic
type AuthClientWrapper struct {
	poolManager *pool.ClientPoolManager[interfaces.CosmosClient]
	logger      zerolog.Logger
}

// NewAuthClientWrapper creates a new auth client wrapper
func NewAuthClientWrapper(poolManager *pool.ClientPoolManager[interfaces.CosmosClient], logger zerolog.Logger) *AuthClientWrapper {
	return &AuthClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "auth").Logger(),
	}
}

func (c *AuthClientWrapper) Accounts(ctx context.Context, req *authtypes.QueryAccountsRequest, opts ...config.CallOpt) (*authtypes.QueryAccountsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*authtypes.QueryAccountsResponse, error) {
		return client.Auth().Accounts(ctx, req, opts...)
	})
}

func (c *AuthClientWrapper) Account(ctx context.Context, req *authtypes.QueryAccountRequest, opts ...config.CallOpt) (*authtypes.QueryAccountResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*authtypes.QueryAccountResponse, error) {
		return client.Auth().Account(ctx, req, opts...)
	})
}

func (c *AuthClientWrapper) AccountAddressByID(ctx context.Context, req *authtypes.QueryAccountAddressByIDRequest, opts ...config.CallOpt) (*authtypes.QueryAccountAddressByIDResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*authtypes.QueryAccountAddressByIDResponse, error) {
		return client.Auth().AccountAddressByID(ctx, req, opts...)
	})
}

func (c *AuthClientWrapper) Params(ctx context.Context, req *authtypes.QueryParamsRequest, opts ...config.CallOpt) (*authtypes.QueryParamsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*authtypes.QueryParamsResponse, error) {
		return client.Auth().Params(ctx, req, opts...)
	})
}

func (c *AuthClientWrapper) ModuleAccounts(ctx context.Context, req *authtypes.QueryModuleAccountsRequest, opts ...config.CallOpt) (*authtypes.QueryModuleAccountsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*authtypes.QueryModuleAccountsResponse, error) {
		return client.Auth().ModuleAccounts(ctx, req, opts...)
	})
}

func (c *AuthClientWrapper) ModuleAccountByName(ctx context.Context, req *authtypes.QueryModuleAccountByNameRequest, opts ...config.CallOpt) (*authtypes.QueryModuleAccountByNameResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*authtypes.QueryModuleAccountByNameResponse, error) {
		return client.Auth().ModuleAccountByName(ctx, req, opts...)
	})
}

func (c *AuthClientWrapper) Bech32Prefix(ctx context.Context, req *authtypes.Bech32PrefixRequest, opts ...config.CallOpt) (*authtypes.Bech32PrefixResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*authtypes.Bech32PrefixResponse, error) {
		return client.Auth().Bech32Prefix(ctx, req, opts...)
	})
}

func (c *AuthClientWrapper) AddressBytesToString(ctx context.Context, req *authtypes.AddressBytesToStringRequest, opts ...config.CallOpt) (*authtypes.AddressBytesToStringResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*authtypes.AddressBytesToStringResponse, error) {
		return client.Auth().AddressBytesToString(ctx, req, opts...)
	})
}

func (c *AuthClientWrapper) AddressStringToBytes(ctx context.Context, req *authtypes.AddressStringToBytesRequest, opts ...config.CallOpt) (*authtypes.AddressStringToBytesResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*authtypes.AddressStringToBytesResponse, error) {
		return client.Auth().AddressStringToBytes(ctx, req, opts...)
	})
}

func (c *AuthClientWrapper) AccountInfo(ctx context.Context, req *authtypes.QueryAccountInfoRequest, opts ...config.CallOpt) (*authtypes.QueryAccountInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*authtypes.QueryAccountInfoResponse, error) {
		return client.Auth().AccountInfo(ctx, req, opts...)
	})
}
