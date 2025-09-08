package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// TendermintGRPCClient provides gRPC access to the tendermint module
type TendermintGRPCClient struct {
	client cmtservice.ServiceClient
	logger zerolog.Logger
}

var _ interfaces.TendermintClient = (*TendermintGRPCClient)(nil)

// NewTendermintGRPCClient creates a new tendermint REST client
func NewTendermintGRPCClient(conn *grpc.ClientConn, logger zerolog.Logger) *TendermintGRPCClient {
	return &TendermintGRPCClient{
		client: cmtservice.NewServiceClient(conn),
		logger: logger.With().Str("module", "tendermint").Str("protocol", "grpc").Logger(),
	}
}

func (c *TendermintGRPCClient) GetNodeInfo(ctx context.Context, req *cmtservice.GetNodeInfoRequest) (*cmtservice.GetNodeInfoResponse, error) {
	resp, err := c.client.GetNodeInfo(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetNodeInfo")
}

func (c *TendermintGRPCClient) GetSyncing(ctx context.Context, req *cmtservice.GetSyncingRequest) (*cmtservice.GetSyncingResponse, error) {
	resp, err := c.client.GetSyncing(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetSyncing")
}

func (c *TendermintGRPCClient) GetLatestBlock(ctx context.Context, req *cmtservice.GetLatestBlockRequest) (*cmtservice.GetLatestBlockResponse, error) {
	resp, err := c.client.GetLatestBlock(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetLatestBlock")
}

func (c *TendermintGRPCClient) GetBlockByHeight(ctx context.Context, req *cmtservice.GetBlockByHeightRequest) (*cmtservice.GetBlockByHeightResponse, error) {
	resp, err := c.client.GetBlockByHeight(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetBlockByHeight")
}

func (c *TendermintGRPCClient) GetLatestValidatorSet(ctx context.Context, req *cmtservice.GetLatestValidatorSetRequest) (*cmtservice.GetLatestValidatorSetResponse, error) {
	resp, err := c.client.GetLatestValidatorSet(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetLatestValidatorSet")
}

func (c *TendermintGRPCClient) GetValidatorSetByHeight(ctx context.Context, req *cmtservice.GetValidatorSetByHeightRequest) (*cmtservice.GetValidatorSetByHeightResponse, error) {
	resp, err := c.client.GetValidatorSetByHeight(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetValidatorSetByHeight")
}

func (c *TendermintGRPCClient) ABCIQuery(ctx context.Context, req *cmtservice.ABCIQueryRequest) (*cmtservice.ABCIQueryResponse, error) {
	resp, err := c.client.ABCIQuery(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.ABCIQuery")
}
