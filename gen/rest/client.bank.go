package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// BankRESTClient provides REST access to the bank module
type BankRESTClient struct {
	*RESTClientCore
	logger zerolog.Logger
}

// NewBankRESTClient creates a new bank REST client
func NewBankRESTClient(core *RESTClientCore, logger zerolog.Logger) *BankRESTClient {
	return &BankRESTClient{
		RESTClientCore: core,
		logger:         logger.With().Str("module", "bank").Str("protocol", "rest").Logger(),
	}
}

func (c *BankRESTClient) Balance(ctx context.Context, req *banktypes.QueryBalanceRequest) (*banktypes.QueryBalanceResponse, error) {
	resp := &banktypes.QueryBalanceResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/balances/{address}/by_denom",
		[]string{"address"}, []string{"denom"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling BankRESTClient.Balance")
}

func (c *BankRESTClient) AllBalances(ctx context.Context, req *banktypes.QueryAllBalancesRequest) (*banktypes.QueryAllBalancesResponse, error) {
	resp := &banktypes.QueryAllBalancesResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/balances/{address}",
		[]string{"address"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse", "resolve_denom"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling BankRESTClient.AllBalances")
}

func (c *BankRESTClient) SpendableBalances(ctx context.Context, req *banktypes.QuerySpendableBalancesRequest) (*banktypes.QuerySpendableBalancesResponse, error) {
	resp := &banktypes.QuerySpendableBalancesResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/spendable_balances/{address}",
		[]string{"address"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling BankRESTClient.SpendableBalances")
}

func (c *BankRESTClient) SpendableBalanceByDenom(ctx context.Context, req *banktypes.QuerySpendableBalanceByDenomRequest) (*banktypes.QuerySpendableBalanceByDenomResponse, error) {
	resp := &banktypes.QuerySpendableBalanceByDenomResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/spendable_balances/{address}/by_denom",
		[]string{"address"}, []string{"denom"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling BankRESTClient.SpendableBalanceByDenom")
}

func (c *BankRESTClient) TotalSupply(ctx context.Context, req *banktypes.QueryTotalSupplyRequest) (*banktypes.QueryTotalSupplyResponse, error) {
	resp := &banktypes.QueryTotalSupplyResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/supply",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling BankRESTClient.TotalSupply")
}

func (c *BankRESTClient) SupplyOf(ctx context.Context, req *banktypes.QuerySupplyOfRequest) (*banktypes.QuerySupplyOfResponse, error) {
	resp := &banktypes.QuerySupplyOfResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/supply/by_denom",
		nil, []string{"denom"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling BankRESTClient.SupplyOf")
}

func (c *BankRESTClient) Params(ctx context.Context, req *banktypes.QueryParamsRequest) (*banktypes.QueryParamsResponse, error) {
	resp := &banktypes.QueryParamsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/params",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling BankRESTClient.Params")
}

func (c *BankRESTClient) DenomsMetadata(ctx context.Context, req *banktypes.QueryDenomsMetadataRequest) (*banktypes.QueryDenomsMetadataResponse, error) {
	resp := &banktypes.QueryDenomsMetadataResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/denoms_metadata",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling BankRESTClient.DenomsMetadata")
}

func (c *BankRESTClient) DenomMetadata(ctx context.Context, req *banktypes.QueryDenomMetadataRequest) (*banktypes.QueryDenomMetadataResponse, error) {
	resp := &banktypes.QueryDenomMetadataResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/denoms_metadata/{denom}",
		[]string{"denom"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling BankRESTClient.DenomMetadata")
}

func (c *BankRESTClient) DenomMetadataByQueryString(ctx context.Context, req *banktypes.QueryDenomMetadataByQueryStringRequest) (*banktypes.QueryDenomMetadataByQueryStringResponse, error) {
	resp := &banktypes.QueryDenomMetadataByQueryStringResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/denoms_metadata_by_query_string",
		nil, []string{"denom"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling BankRESTClient.DenomMetadataByQueryString")
}

func (c *BankRESTClient) DenomOwners(ctx context.Context, req *banktypes.QueryDenomOwnersRequest) (*banktypes.QueryDenomOwnersResponse, error) {
	resp := &banktypes.QueryDenomOwnersResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/denom_owners/{denom}",
		[]string{"denom"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling BankRESTClient.DenomOwners")
}

func (c *BankRESTClient) DenomOwnersByQuery(ctx context.Context, req *banktypes.QueryDenomOwnersByQueryRequest) (*banktypes.QueryDenomOwnersByQueryResponse, error) {
	resp := &banktypes.QueryDenomOwnersByQueryResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/denom_owners_by_query",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse", "denom"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling BankRESTClient.DenomOwnersByQuery")
}

func (c *BankRESTClient) SendEnabled(ctx context.Context, req *banktypes.QuerySendEnabledRequest) (*banktypes.QuerySendEnabledResponse, error) {
	resp := &banktypes.QuerySendEnabledResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/send_enabled",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse", "denoms"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling BankRESTClient.SendEnabled")
}
