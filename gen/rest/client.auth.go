package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/allora-network/allora-sdk-go/config"
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

func (c *AuthRESTClient) Account(ctx context.Context, req *authtypes.QueryAccountRequest, opts ...config.CallOpt) (*authtypes.QueryAccountResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &authtypes.QueryAccountResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/auth/v1beta1/accounts/{address}",
		[]string{"address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling AuthRESTClient.Account")
	}
	return resp, nil
}

func (c *AuthRESTClient) AccountAddressByID(ctx context.Context, req *authtypes.QueryAccountAddressByIDRequest, opts ...config.CallOpt) (*authtypes.QueryAccountAddressByIDResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &authtypes.QueryAccountAddressByIDResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/auth/v1beta1/address_by_id/{id}",
		[]string{"id"}, []string{"account_id"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling AuthRESTClient.AccountAddressByID")
	}
	return resp, nil
}

func (c *AuthRESTClient) AccountInfo(ctx context.Context, req *authtypes.QueryAccountInfoRequest, opts ...config.CallOpt) (*authtypes.QueryAccountInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &authtypes.QueryAccountInfoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/auth/v1beta1/account_info/{address}",
		[]string{"address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling AuthRESTClient.AccountInfo")
	}
	return resp, nil
}

func (c *AuthRESTClient) Accounts(ctx context.Context, req *authtypes.QueryAccountsRequest, opts ...config.CallOpt) (*authtypes.QueryAccountsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &authtypes.QueryAccountsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/auth/v1beta1/accounts",
		nil, []string{"pagination.key", "pagination.offset", "pagination.limit", "pagination.count_total", "pagination.reverse"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling AuthRESTClient.Accounts")
	}
	return resp, nil
}

func (c *AuthRESTClient) AddressBytesToString(ctx context.Context, req *authtypes.AddressBytesToStringRequest, opts ...config.CallOpt) (*authtypes.AddressBytesToStringResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &authtypes.AddressBytesToStringResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/auth/v1beta1/bech32/{address_bytes}",
		[]string{"address_bytes"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling AuthRESTClient.AddressBytesToString")
	}
	return resp, nil
}

func (c *AuthRESTClient) AddressStringToBytes(ctx context.Context, req *authtypes.AddressStringToBytesRequest, opts ...config.CallOpt) (*authtypes.AddressStringToBytesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &authtypes.AddressStringToBytesResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/auth/v1beta1/bech32/{address_string}",
		[]string{"address_string"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling AuthRESTClient.AddressStringToBytes")
	}
	return resp, nil
}

func (c *AuthRESTClient) Bech32Prefix(ctx context.Context, req *authtypes.Bech32PrefixRequest, opts ...config.CallOpt) (*authtypes.Bech32PrefixResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &authtypes.Bech32PrefixResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/auth/v1beta1/bech32",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling AuthRESTClient.Bech32Prefix")
	}
	return resp, nil
}

func (c *AuthRESTClient) ModuleAccountByName(ctx context.Context, req *authtypes.QueryModuleAccountByNameRequest, opts ...config.CallOpt) (*authtypes.QueryModuleAccountByNameResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &authtypes.QueryModuleAccountByNameResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/auth/v1beta1/module_accounts/{name}",
		[]string{"name"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling AuthRESTClient.ModuleAccountByName")
	}
	return resp, nil
}

func (c *AuthRESTClient) ModuleAccounts(ctx context.Context, req *authtypes.QueryModuleAccountsRequest, opts ...config.CallOpt) (*authtypes.QueryModuleAccountsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &authtypes.QueryModuleAccountsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/auth/v1beta1/module_accounts",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling AuthRESTClient.ModuleAccounts")
	}
	return resp, nil
}

func (c *AuthRESTClient) Params(ctx context.Context, req *authtypes.QueryParamsRequest, opts ...config.CallOpt) (*authtypes.QueryParamsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &authtypes.QueryParamsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/cosmos/auth/v1beta1/params",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling AuthRESTClient.Params")
	}
	return resp, nil
}
