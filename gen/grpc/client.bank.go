package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

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

func (c *BankGRPCClient) Balance(ctx context.Context, req *banktypes.QueryBalanceRequest) (*banktypes.QueryBalanceResponse, error) {
	resp, err := c.client.Balance(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Balance")
}

func (c *BankGRPCClient) AllBalances(ctx context.Context, req *banktypes.QueryAllBalancesRequest) (*banktypes.QueryAllBalancesResponse, error) {
	resp, err := c.client.AllBalances(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.AllBalances")
}

func (c *BankGRPCClient) SpendableBalances(ctx context.Context, req *banktypes.QuerySpendableBalancesRequest) (*banktypes.QuerySpendableBalancesResponse, error) {
	resp, err := c.client.SpendableBalances(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.SpendableBalances")
}

func (c *BankGRPCClient) SpendableBalanceByDenom(ctx context.Context, req *banktypes.QuerySpendableBalanceByDenomRequest) (*banktypes.QuerySpendableBalanceByDenomResponse, error) {
	resp, err := c.client.SpendableBalanceByDenom(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.SpendableBalanceByDenom")
}

func (c *BankGRPCClient) TotalSupply(ctx context.Context, req *banktypes.QueryTotalSupplyRequest) (*banktypes.QueryTotalSupplyResponse, error) {
	resp, err := c.client.TotalSupply(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.TotalSupply")
}

func (c *BankGRPCClient) SupplyOf(ctx context.Context, req *banktypes.QuerySupplyOfRequest) (*banktypes.QuerySupplyOfResponse, error) {
	resp, err := c.client.SupplyOf(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.SupplyOf")
}

func (c *BankGRPCClient) Params(ctx context.Context, req *banktypes.QueryParamsRequest) (*banktypes.QueryParamsResponse, error) {
	resp, err := c.client.Params(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.Params")
}

func (c *BankGRPCClient) DenomsMetadata(ctx context.Context, req *banktypes.QueryDenomsMetadataRequest) (*banktypes.QueryDenomsMetadataResponse, error) {
	resp, err := c.client.DenomsMetadata(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.DenomsMetadata")
}

func (c *BankGRPCClient) DenomMetadata(ctx context.Context, req *banktypes.QueryDenomMetadataRequest) (*banktypes.QueryDenomMetadataResponse, error) {
	resp, err := c.client.DenomMetadata(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.DenomMetadata")
}

func (c *BankGRPCClient) DenomMetadataByQueryString(ctx context.Context, req *banktypes.QueryDenomMetadataByQueryStringRequest) (*banktypes.QueryDenomMetadataByQueryStringResponse, error) {
	resp, err := c.client.DenomMetadataByQueryString(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.DenomMetadataByQueryString")
}

func (c *BankGRPCClient) DenomOwners(ctx context.Context, req *banktypes.QueryDenomOwnersRequest) (*banktypes.QueryDenomOwnersResponse, error) {
	resp, err := c.client.DenomOwners(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.DenomOwners")
}

func (c *BankGRPCClient) DenomOwnersByQuery(ctx context.Context, req *banktypes.QueryDenomOwnersByQueryRequest) (*banktypes.QueryDenomOwnersByQueryResponse, error) {
	resp, err := c.client.DenomOwnersByQuery(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.DenomOwnersByQuery")
}

func (c *BankGRPCClient) SendEnabled(ctx context.Context, req *banktypes.QuerySendEnabledRequest) (*banktypes.QuerySendEnabledResponse, error) {
	resp, err := c.client.SendEnabled(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.SendEnabled")
}
