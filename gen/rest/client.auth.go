package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AuthRESTClient provides REST access to the auth module
type AuthRESTClient struct {
	*RESTClientCore
	logger zerolog.Logger
}

// NewAuthRESTClient creates a new auth REST client
func NewAuthRESTClient(core *RESTClientCore, logger zerolog.Logger) *AuthRESTClient {
	return &AuthRESTClient{
		RESTClientCore: core,
		logger:         logger.With().Str("module", "auth").Str("protocol", "rest").Logger(),
	}
}

func (c *AuthRESTClient) Accounts(ctx context.Context, req *authtypes.QueryAccountsRequest) (*authtypes.QueryAccountsResponse, error) {
	resp := &authtypes.QueryAccountsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/auth/v1beta1/accounts",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling AuthRESTClient.Accounts")
}

func (c *AuthRESTClient) Account(ctx context.Context, req *authtypes.QueryAccountRequest) (*authtypes.QueryAccountResponse, error) {
	resp := &authtypes.QueryAccountResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/auth/v1beta1/accounts/{address}",
		[]string{"address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling AuthRESTClient.Account")
}

func (c *AuthRESTClient) AccountAddressByID(ctx context.Context, req *authtypes.QueryAccountAddressByIDRequest) (*authtypes.QueryAccountAddressByIDResponse, error) {
	resp := &authtypes.QueryAccountAddressByIDResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/auth/v1beta1/address_by_id/{id}",
		[]string{"id"}, []string{"account_id"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling AuthRESTClient.AccountAddressByID")
}

func (c *AuthRESTClient) Params(ctx context.Context, req *authtypes.QueryParamsRequest) (*authtypes.QueryParamsResponse, error) {
	resp := &authtypes.QueryParamsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/auth/v1beta1/params",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling AuthRESTClient.Params")
}

func (c *AuthRESTClient) ModuleAccounts(ctx context.Context, req *authtypes.QueryModuleAccountsRequest) (*authtypes.QueryModuleAccountsResponse, error) {
	resp := &authtypes.QueryModuleAccountsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/auth/v1beta1/module_accounts",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling AuthRESTClient.ModuleAccounts")
}

func (c *AuthRESTClient) ModuleAccountByName(ctx context.Context, req *authtypes.QueryModuleAccountByNameRequest) (*authtypes.QueryModuleAccountByNameResponse, error) {
	resp := &authtypes.QueryModuleAccountByNameResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/auth/v1beta1/module_accounts/{name}",
		[]string{"name"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling AuthRESTClient.ModuleAccountByName")
}

func (c *AuthRESTClient) Bech32Prefix(ctx context.Context, req *authtypes.Bech32PrefixRequest) (*authtypes.Bech32PrefixResponse, error) {
	resp := &authtypes.Bech32PrefixResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/auth/v1beta1/bech32",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling AuthRESTClient.Bech32Prefix")
}

func (c *AuthRESTClient) AddressBytesToString(ctx context.Context, req *authtypes.AddressBytesToStringRequest) (*authtypes.AddressBytesToStringResponse, error) {
	resp := &authtypes.AddressBytesToStringResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/auth/v1beta1/bech32/{address_bytes}",
		[]string{"address_bytes"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling AuthRESTClient.AddressBytesToString")
}

func (c *AuthRESTClient) AddressStringToBytes(ctx context.Context, req *authtypes.AddressStringToBytesRequest) (*authtypes.AddressStringToBytesResponse, error) {
	resp := &authtypes.AddressStringToBytesResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/auth/v1beta1/bech32/{address_string}",
		[]string{"address_string"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling AuthRESTClient.AddressStringToBytes")
}

func (c *AuthRESTClient) AccountInfo(ctx context.Context, req *authtypes.QueryAccountInfoRequest) (*authtypes.QueryAccountInfoResponse, error) {
	resp := &authtypes.QueryAccountInfoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/auth/v1beta1/account_info/{address}",
		[]string{"address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling AuthRESTClient.AccountInfo")
}
