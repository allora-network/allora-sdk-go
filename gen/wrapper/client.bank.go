package wrapper

import (
	"context"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/pool"
)

// BankClientWrapper wraps the bank module with pool management and retry logic
type BankClientWrapper struct {
	poolManager *pool.ClientPoolManager[interfaces.CosmosClient]
	logger      zerolog.Logger
}

// NewBankClientWrapper creates a new bank client wrapper
func NewBankClientWrapper(poolManager *pool.ClientPoolManager[interfaces.CosmosClient], logger zerolog.Logger) *BankClientWrapper {
	return &BankClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "bank").Logger(),
	}
}

func (c *BankClientWrapper) AllBalances(ctx context.Context, req *banktypes.QueryAllBalancesRequest, opts ...config.CallOpt) (*banktypes.QueryAllBalancesResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*banktypes.QueryAllBalancesResponse, error) {
		return client.Bank().AllBalances(ctx, req, opts...)
	})
}

func (c *BankClientWrapper) Balance(ctx context.Context, req *banktypes.QueryBalanceRequest, opts ...config.CallOpt) (*banktypes.QueryBalanceResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*banktypes.QueryBalanceResponse, error) {
		return client.Bank().Balance(ctx, req, opts...)
	})
}

func (c *BankClientWrapper) DenomMetadata(ctx context.Context, req *banktypes.QueryDenomMetadataRequest, opts ...config.CallOpt) (*banktypes.QueryDenomMetadataResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*banktypes.QueryDenomMetadataResponse, error) {
		return client.Bank().DenomMetadata(ctx, req, opts...)
	})
}

func (c *BankClientWrapper) DenomMetadataByQueryString(ctx context.Context, req *banktypes.QueryDenomMetadataByQueryStringRequest, opts ...config.CallOpt) (*banktypes.QueryDenomMetadataByQueryStringResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*banktypes.QueryDenomMetadataByQueryStringResponse, error) {
		return client.Bank().DenomMetadataByQueryString(ctx, req, opts...)
	})
}

func (c *BankClientWrapper) DenomOwners(ctx context.Context, req *banktypes.QueryDenomOwnersRequest, opts ...config.CallOpt) (*banktypes.QueryDenomOwnersResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*banktypes.QueryDenomOwnersResponse, error) {
		return client.Bank().DenomOwners(ctx, req, opts...)
	})
}

func (c *BankClientWrapper) DenomOwnersByQuery(ctx context.Context, req *banktypes.QueryDenomOwnersByQueryRequest, opts ...config.CallOpt) (*banktypes.QueryDenomOwnersByQueryResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*banktypes.QueryDenomOwnersByQueryResponse, error) {
		return client.Bank().DenomOwnersByQuery(ctx, req, opts...)
	})
}

func (c *BankClientWrapper) DenomsMetadata(ctx context.Context, req *banktypes.QueryDenomsMetadataRequest, opts ...config.CallOpt) (*banktypes.QueryDenomsMetadataResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*banktypes.QueryDenomsMetadataResponse, error) {
		return client.Bank().DenomsMetadata(ctx, req, opts...)
	})
}

func (c *BankClientWrapper) Params(ctx context.Context, req *banktypes.QueryParamsRequest, opts ...config.CallOpt) (*banktypes.QueryParamsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*banktypes.QueryParamsResponse, error) {
		return client.Bank().Params(ctx, req, opts...)
	})
}

func (c *BankClientWrapper) SendEnabled(ctx context.Context, req *banktypes.QuerySendEnabledRequest, opts ...config.CallOpt) (*banktypes.QuerySendEnabledResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*banktypes.QuerySendEnabledResponse, error) {
		return client.Bank().SendEnabled(ctx, req, opts...)
	})
}

func (c *BankClientWrapper) SpendableBalanceByDenom(ctx context.Context, req *banktypes.QuerySpendableBalanceByDenomRequest, opts ...config.CallOpt) (*banktypes.QuerySpendableBalanceByDenomResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*banktypes.QuerySpendableBalanceByDenomResponse, error) {
		return client.Bank().SpendableBalanceByDenom(ctx, req, opts...)
	})
}

func (c *BankClientWrapper) SpendableBalances(ctx context.Context, req *banktypes.QuerySpendableBalancesRequest, opts ...config.CallOpt) (*banktypes.QuerySpendableBalancesResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*banktypes.QuerySpendableBalancesResponse, error) {
		return client.Bank().SpendableBalances(ctx, req, opts...)
	})
}

func (c *BankClientWrapper) SupplyOf(ctx context.Context, req *banktypes.QuerySupplyOfRequest, opts ...config.CallOpt) (*banktypes.QuerySupplyOfResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*banktypes.QuerySupplyOfResponse, error) {
		return client.Bank().SupplyOf(ctx, req, opts...)
	})
}

func (c *BankClientWrapper) TotalSupply(ctx context.Context, req *banktypes.QueryTotalSupplyRequest, opts ...config.CallOpt) (*banktypes.QueryTotalSupplyResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*banktypes.QueryTotalSupplyResponse, error) {
		return client.Bank().TotalSupply(ctx, req, opts...)
	})
}
