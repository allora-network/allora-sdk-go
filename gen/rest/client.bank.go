package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/allora-network/allora-sdk-go/config"
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

func (c *BankRESTClient) AllBalances(ctx context.Context, req *banktypes.QueryAllBalancesRequest, opts ...config.CallOpt) (*banktypes.QueryAllBalancesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &banktypes.QueryAllBalancesResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/balances/{address}",
		[]string{"address"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse", "resolve_denom"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling BankRESTClient.AllBalances")
	}
	return resp, nil
}

func (c *BankRESTClient) Balance(ctx context.Context, req *banktypes.QueryBalanceRequest, opts ...config.CallOpt) (*banktypes.QueryBalanceResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &banktypes.QueryBalanceResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/balances/{address}/by_denom",
		[]string{"address"}, []string{"denom"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling BankRESTClient.Balance")
	}
	return resp, nil
}

func (c *BankRESTClient) DenomMetadata(ctx context.Context, req *banktypes.QueryDenomMetadataRequest, opts ...config.CallOpt) (*banktypes.QueryDenomMetadataResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &banktypes.QueryDenomMetadataResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/denoms_metadata/{denom}",
		[]string{"denom"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling BankRESTClient.DenomMetadata")
	}
	return resp, nil
}

func (c *BankRESTClient) DenomMetadataByQueryString(ctx context.Context, req *banktypes.QueryDenomMetadataByQueryStringRequest, opts ...config.CallOpt) (*banktypes.QueryDenomMetadataByQueryStringResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &banktypes.QueryDenomMetadataByQueryStringResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/denoms_metadata_by_query_string",
		nil, []string{"denom"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling BankRESTClient.DenomMetadataByQueryString")
	}
	return resp, nil
}

func (c *BankRESTClient) DenomOwners(ctx context.Context, req *banktypes.QueryDenomOwnersRequest, opts ...config.CallOpt) (*banktypes.QueryDenomOwnersResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &banktypes.QueryDenomOwnersResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/denom_owners/{denom}",
		[]string{"denom"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling BankRESTClient.DenomOwners")
	}
	return resp, nil
}

func (c *BankRESTClient) DenomOwnersByQuery(ctx context.Context, req *banktypes.QueryDenomOwnersByQueryRequest, opts ...config.CallOpt) (*banktypes.QueryDenomOwnersByQueryResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &banktypes.QueryDenomOwnersByQueryResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/denom_owners_by_query",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse", "denom"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling BankRESTClient.DenomOwnersByQuery")
	}
	return resp, nil
}

func (c *BankRESTClient) DenomsMetadata(ctx context.Context, req *banktypes.QueryDenomsMetadataRequest, opts ...config.CallOpt) (*banktypes.QueryDenomsMetadataResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &banktypes.QueryDenomsMetadataResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/denoms_metadata",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling BankRESTClient.DenomsMetadata")
	}
	return resp, nil
}

func (c *BankRESTClient) Params(ctx context.Context, req *banktypes.QueryParamsRequest, opts ...config.CallOpt) (*banktypes.QueryParamsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &banktypes.QueryParamsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/params",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling BankRESTClient.Params")
	}
	return resp, nil
}

func (c *BankRESTClient) SendEnabled(ctx context.Context, req *banktypes.QuerySendEnabledRequest, opts ...config.CallOpt) (*banktypes.QuerySendEnabledResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &banktypes.QuerySendEnabledResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/send_enabled",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse", "denoms"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling BankRESTClient.SendEnabled")
	}
	return resp, nil
}

func (c *BankRESTClient) SpendableBalanceByDenom(ctx context.Context, req *banktypes.QuerySpendableBalanceByDenomRequest, opts ...config.CallOpt) (*banktypes.QuerySpendableBalanceByDenomResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &banktypes.QuerySpendableBalanceByDenomResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/spendable_balances/{address}/by_denom",
		[]string{"address"}, []string{"denom"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling BankRESTClient.SpendableBalanceByDenom")
	}
	return resp, nil
}

func (c *BankRESTClient) SpendableBalances(ctx context.Context, req *banktypes.QuerySpendableBalancesRequest, opts ...config.CallOpt) (*banktypes.QuerySpendableBalancesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &banktypes.QuerySpendableBalancesResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/spendable_balances/{address}",
		[]string{"address"}, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling BankRESTClient.SpendableBalances")
	}
	return resp, nil
}

func (c *BankRESTClient) SupplyOf(ctx context.Context, req *banktypes.QuerySupplyOfRequest, opts ...config.CallOpt) (*banktypes.QuerySupplyOfResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &banktypes.QuerySupplyOfResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/supply/by_denom",
		nil, []string{"denom"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling BankRESTClient.SupplyOf")
	}
	return resp, nil
}

func (c *BankRESTClient) TotalSupply(ctx context.Context, req *banktypes.QueryTotalSupplyRequest, opts ...config.CallOpt) (*banktypes.QueryTotalSupplyResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &banktypes.QueryTotalSupplyResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/bank/v1beta1/supply",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling BankRESTClient.TotalSupply")
	}
	return resp, nil
}
