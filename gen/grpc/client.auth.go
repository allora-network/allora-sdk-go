package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

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

func (c *AuthGRPCClient) Accounts(ctx context.Context, req *authtypes.QueryAccountsRequest) (*authtypes.QueryAccountsResponse, error) {
	resp, err := c.client.Accounts(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Accounts")
}

func (c *AuthGRPCClient) Account(ctx context.Context, req *authtypes.QueryAccountRequest) (*authtypes.QueryAccountResponse, error) {
	resp, err := c.client.Account(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Account")
}

func (c *AuthGRPCClient) AccountAddressByID(ctx context.Context, req *authtypes.QueryAccountAddressByIDRequest) (*authtypes.QueryAccountAddressByIDResponse, error) {
	resp, err := c.client.AccountAddressByID(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.AccountAddressByID")
}

func (c *AuthGRPCClient) Params(ctx context.Context, req *authtypes.QueryParamsRequest) (*authtypes.QueryParamsResponse, error) {
	resp, err := c.client.Params(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Params")
}

func (c *AuthGRPCClient) ModuleAccounts(ctx context.Context, req *authtypes.QueryModuleAccountsRequest) (*authtypes.QueryModuleAccountsResponse, error) {
	resp, err := c.client.ModuleAccounts(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.ModuleAccounts")
}

func (c *AuthGRPCClient) ModuleAccountByName(ctx context.Context, req *authtypes.QueryModuleAccountByNameRequest) (*authtypes.QueryModuleAccountByNameResponse, error) {
	resp, err := c.client.ModuleAccountByName(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.ModuleAccountByName")
}

func (c *AuthGRPCClient) Bech32Prefix(ctx context.Context, req *authtypes.Bech32PrefixRequest) (*authtypes.Bech32PrefixResponse, error) {
	resp, err := c.client.Bech32Prefix(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Bech32Prefix")
}

func (c *AuthGRPCClient) AddressBytesToString(ctx context.Context, req *authtypes.AddressBytesToStringRequest) (*authtypes.AddressBytesToStringResponse, error) {
	resp, err := c.client.AddressBytesToString(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.AddressBytesToString")
}

func (c *AuthGRPCClient) AddressStringToBytes(ctx context.Context, req *authtypes.AddressStringToBytesRequest) (*authtypes.AddressStringToBytesResponse, error) {
	resp, err := c.client.AddressStringToBytes(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.AddressStringToBytes")
}

func (c *AuthGRPCClient) AccountInfo(ctx context.Context, req *authtypes.QueryAccountInfoRequest) (*authtypes.QueryAccountInfoResponse, error) {
	resp, err := c.client.AccountInfo(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.AccountInfo")
}
