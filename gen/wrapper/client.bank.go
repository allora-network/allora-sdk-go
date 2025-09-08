package wrapper

import (
	"context"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/pool"
)

// BankClientWrapper wraps the bank module with pool management and retry logic
type BankClientWrapper struct {
	poolManager *pool.ClientPoolManager
	logger      zerolog.Logger
}

// NewBankClientWrapper creates a new bank client wrapper
func NewBankClientWrapper(poolManager *pool.ClientPoolManager, logger zerolog.Logger) *BankClientWrapper {
	return &BankClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "bank").Logger(),
	}
}

func (c *BankClientWrapper) Balance(ctx context.Context, req *banktypes.QueryBalanceRequest) (*banktypes.QueryBalanceResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*banktypes.QueryBalanceResponse, error) {
		return client.Bank().Balance(ctx, req)
	})
}

func (c *BankClientWrapper) AllBalances(ctx context.Context, req *banktypes.QueryAllBalancesRequest) (*banktypes.QueryAllBalancesResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*banktypes.QueryAllBalancesResponse, error) {
		return client.Bank().AllBalances(ctx, req)
	})
}

func (c *BankClientWrapper) SpendableBalances(ctx context.Context, req *banktypes.QuerySpendableBalancesRequest) (*banktypes.QuerySpendableBalancesResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*banktypes.QuerySpendableBalancesResponse, error) {
		return client.Bank().SpendableBalances(ctx, req)
	})
}

func (c *BankClientWrapper) SpendableBalanceByDenom(ctx context.Context, req *banktypes.QuerySpendableBalanceByDenomRequest) (*banktypes.QuerySpendableBalanceByDenomResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*banktypes.QuerySpendableBalanceByDenomResponse, error) {
		return client.Bank().SpendableBalanceByDenom(ctx, req)
	})
}

func (c *BankClientWrapper) TotalSupply(ctx context.Context, req *banktypes.QueryTotalSupplyRequest) (*banktypes.QueryTotalSupplyResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*banktypes.QueryTotalSupplyResponse, error) {
		return client.Bank().TotalSupply(ctx, req)
	})
}

func (c *BankClientWrapper) SupplyOf(ctx context.Context, req *banktypes.QuerySupplyOfRequest) (*banktypes.QuerySupplyOfResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*banktypes.QuerySupplyOfResponse, error) {
		return client.Bank().SupplyOf(ctx, req)
	})
}

func (c *BankClientWrapper) Params(ctx context.Context, req *banktypes.QueryParamsRequest) (*banktypes.QueryParamsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*banktypes.QueryParamsResponse, error) {
		return client.Bank().Params(ctx, req)
	})
}

func (c *BankClientWrapper) DenomsMetadata(ctx context.Context, req *banktypes.QueryDenomsMetadataRequest) (*banktypes.QueryDenomsMetadataResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*banktypes.QueryDenomsMetadataResponse, error) {
		return client.Bank().DenomsMetadata(ctx, req)
	})
}

func (c *BankClientWrapper) DenomMetadata(ctx context.Context, req *banktypes.QueryDenomMetadataRequest) (*banktypes.QueryDenomMetadataResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*banktypes.QueryDenomMetadataResponse, error) {
		return client.Bank().DenomMetadata(ctx, req)
	})
}

func (c *BankClientWrapper) DenomMetadataByQueryString(ctx context.Context, req *banktypes.QueryDenomMetadataByQueryStringRequest) (*banktypes.QueryDenomMetadataByQueryStringResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*banktypes.QueryDenomMetadataByQueryStringResponse, error) {
		return client.Bank().DenomMetadataByQueryString(ctx, req)
	})
}

func (c *BankClientWrapper) DenomOwners(ctx context.Context, req *banktypes.QueryDenomOwnersRequest) (*banktypes.QueryDenomOwnersResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*banktypes.QueryDenomOwnersResponse, error) {
		return client.Bank().DenomOwners(ctx, req)
	})
}

func (c *BankClientWrapper) DenomOwnersByQuery(ctx context.Context, req *banktypes.QueryDenomOwnersByQueryRequest) (*banktypes.QueryDenomOwnersByQueryResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*banktypes.QueryDenomOwnersByQueryResponse, error) {
		return client.Bank().DenomOwnersByQuery(ctx, req)
	})
}

func (c *BankClientWrapper) SendEnabled(ctx context.Context, req *banktypes.QuerySendEnabledRequest) (*banktypes.QuerySendEnabledResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*banktypes.QuerySendEnabledResponse, error) {
		return client.Bank().SendEnabled(ctx, req)
	})
}
