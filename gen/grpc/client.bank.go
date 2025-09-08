package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/allora-network/allora-sdk-go/config"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// BankGRPCClient provides gRPC access to the bank module
type BankGRPCClient struct {
	client banktypes.QueryClient
	logger zerolog.Logger
}

var _ interfaces.BankClient = (*BankGRPCClient)(nil)

// NewBankGRPCClient creates a new bank REST client
func NewBankGRPCClient(conn *grpc.ClientConn, logger zerolog.Logger) *BankGRPCClient {
	return &BankGRPCClient{
		client: banktypes.NewQueryClient(conn),
		logger: logger.With().Str("module", "bank").Str("protocol", "grpc").Logger(),
	}
}

func (c *BankGRPCClient) Balance(ctx context.Context, req *banktypes.QueryBalanceRequest, opts ...config.CallOpt) (*banktypes.QueryBalanceResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Balance, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling BankGRPCClient.Balance")
	}
	return resp, nil
}

func (c *BankGRPCClient) AllBalances(ctx context.Context, req *banktypes.QueryAllBalancesRequest, opts ...config.CallOpt) (*banktypes.QueryAllBalancesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.AllBalances, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling BankGRPCClient.AllBalances")
	}
	return resp, nil
}

func (c *BankGRPCClient) SpendableBalances(ctx context.Context, req *banktypes.QuerySpendableBalancesRequest, opts ...config.CallOpt) (*banktypes.QuerySpendableBalancesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.SpendableBalances, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling BankGRPCClient.SpendableBalances")
	}
	return resp, nil
}

func (c *BankGRPCClient) SpendableBalanceByDenom(ctx context.Context, req *banktypes.QuerySpendableBalanceByDenomRequest, opts ...config.CallOpt) (*banktypes.QuerySpendableBalanceByDenomResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.SpendableBalanceByDenom, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling BankGRPCClient.SpendableBalanceByDenom")
	}
	return resp, nil
}

func (c *BankGRPCClient) TotalSupply(ctx context.Context, req *banktypes.QueryTotalSupplyRequest, opts ...config.CallOpt) (*banktypes.QueryTotalSupplyResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.TotalSupply, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling BankGRPCClient.TotalSupply")
	}
	return resp, nil
}

func (c *BankGRPCClient) SupplyOf(ctx context.Context, req *banktypes.QuerySupplyOfRequest, opts ...config.CallOpt) (*banktypes.QuerySupplyOfResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.SupplyOf, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling BankGRPCClient.SupplyOf")
	}
	return resp, nil
}

func (c *BankGRPCClient) Params(ctx context.Context, req *banktypes.QueryParamsRequest, opts ...config.CallOpt) (*banktypes.QueryParamsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.Params, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling BankGRPCClient.Params")
	}
	return resp, nil
}

func (c *BankGRPCClient) DenomMetadata(ctx context.Context, req *banktypes.QueryDenomMetadataRequest, opts ...config.CallOpt) (*banktypes.QueryDenomMetadataResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.DenomMetadata, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling BankGRPCClient.DenomMetadata")
	}
	return resp, nil
}

func (c *BankGRPCClient) DenomMetadataByQueryString(ctx context.Context, req *banktypes.QueryDenomMetadataByQueryStringRequest, opts ...config.CallOpt) (*banktypes.QueryDenomMetadataByQueryStringResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.DenomMetadataByQueryString, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling BankGRPCClient.DenomMetadataByQueryString")
	}
	return resp, nil
}

func (c *BankGRPCClient) DenomsMetadata(ctx context.Context, req *banktypes.QueryDenomsMetadataRequest, opts ...config.CallOpt) (*banktypes.QueryDenomsMetadataResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.DenomsMetadata, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling BankGRPCClient.DenomsMetadata")
	}
	return resp, nil
}

func (c *BankGRPCClient) DenomOwners(ctx context.Context, req *banktypes.QueryDenomOwnersRequest, opts ...config.CallOpt) (*banktypes.QueryDenomOwnersResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.DenomOwners, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling BankGRPCClient.DenomOwners")
	}
	return resp, nil
}

func (c *BankGRPCClient) DenomOwnersByQuery(ctx context.Context, req *banktypes.QueryDenomOwnersByQueryRequest, opts ...config.CallOpt) (*banktypes.QueryDenomOwnersByQueryResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.DenomOwnersByQuery, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling BankGRPCClient.DenomOwnersByQuery")
	}
	return resp, nil
}

func (c *BankGRPCClient) SendEnabled(ctx context.Context, req *banktypes.QuerySendEnabledRequest, opts ...config.CallOpt) (*banktypes.QuerySendEnabledResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.SendEnabled, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling BankGRPCClient.SendEnabled")
	}
	return resp, nil
}
