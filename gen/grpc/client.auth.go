package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// AuthGRPCClient provides gRPC access to the auth module
type AuthGRPCClient struct {
	client authtypes.QueryClient
	logger zerolog.Logger
}

var _ interfaces.AuthClient = (*AuthGRPCClient)(nil)

// NewAuthGRPCClient creates a new auth REST client
func NewAuthGRPCClient(conn *grpc.ClientConn, logger zerolog.Logger) *AuthGRPCClient {
	return &AuthGRPCClient{
		client: authtypes.NewQueryClient(conn),
		logger: logger.With().Str("module", "auth").Str("protocol", "grpc").Logger(),
	}
}

func (c *AuthGRPCClient) Account(ctx context.Context, req *authtypes.QueryAccountRequest, opts ...config.CallOpt) (*authtypes.QueryAccountResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Account, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling AuthGRPCClient.Account")
	}
	return resp, nil
}

func (c *AuthGRPCClient) AccountAddressByID(ctx context.Context, req *authtypes.QueryAccountAddressByIDRequest, opts ...config.CallOpt) (*authtypes.QueryAccountAddressByIDResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.AccountAddressByID, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling AuthGRPCClient.AccountAddressByID")
	}
	return resp, nil
}

func (c *AuthGRPCClient) AccountInfo(ctx context.Context, req *authtypes.QueryAccountInfoRequest, opts ...config.CallOpt) (*authtypes.QueryAccountInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.AccountInfo, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling AuthGRPCClient.AccountInfo")
	}
	return resp, nil
}

func (c *AuthGRPCClient) Accounts(ctx context.Context, req *authtypes.QueryAccountsRequest, opts ...config.CallOpt) (*authtypes.QueryAccountsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Accounts, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling AuthGRPCClient.Accounts")
	}
	return resp, nil
}

func (c *AuthGRPCClient) AddressBytesToString(ctx context.Context, req *authtypes.AddressBytesToStringRequest, opts ...config.CallOpt) (*authtypes.AddressBytesToStringResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.AddressBytesToString, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling AuthGRPCClient.AddressBytesToString")
	}
	return resp, nil
}

func (c *AuthGRPCClient) AddressStringToBytes(ctx context.Context, req *authtypes.AddressStringToBytesRequest, opts ...config.CallOpt) (*authtypes.AddressStringToBytesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.AddressStringToBytes, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling AuthGRPCClient.AddressStringToBytes")
	}
	return resp, nil
}

func (c *AuthGRPCClient) Bech32Prefix(ctx context.Context, req *authtypes.Bech32PrefixRequest, opts ...config.CallOpt) (*authtypes.Bech32PrefixResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Bech32Prefix, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling AuthGRPCClient.Bech32Prefix")
	}
	return resp, nil
}

func (c *AuthGRPCClient) ModuleAccountByName(ctx context.Context, req *authtypes.QueryModuleAccountByNameRequest, opts ...config.CallOpt) (*authtypes.QueryModuleAccountByNameResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.ModuleAccountByName, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling AuthGRPCClient.ModuleAccountByName")
	}
	return resp, nil
}

func (c *AuthGRPCClient) ModuleAccounts(ctx context.Context, req *authtypes.QueryModuleAccountsRequest, opts ...config.CallOpt) (*authtypes.QueryModuleAccountsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.ModuleAccounts, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling AuthGRPCClient.ModuleAccounts")
	}
	return resp, nil
}

func (c *AuthGRPCClient) Params(ctx context.Context, req *authtypes.QueryParamsRequest, opts ...config.CallOpt) (*authtypes.QueryParamsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Params, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling AuthGRPCClient.Params")
	}
	return resp, nil
}
