package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"

	"github.com/allora-network/allora-sdk-go/config"
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

func (c *TendermintGRPCClient) ABCIQuery(ctx context.Context, req *cmtservice.ABCIQueryRequest, opts ...config.CallOpt) (*cmtservice.ABCIQueryResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.ABCIQuery, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling TendermintGRPCClient.ABCIQuery")
	}
	return resp, nil
}

func (c *TendermintGRPCClient) GetBlockByHeight(ctx context.Context, req *cmtservice.GetBlockByHeightRequest, opts ...config.CallOpt) (*cmtservice.GetBlockByHeightResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetBlockByHeight, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling TendermintGRPCClient.GetBlockByHeight")
	}
	return resp, nil
}

func (c *TendermintGRPCClient) GetLatestBlock(ctx context.Context, req *cmtservice.GetLatestBlockRequest, opts ...config.CallOpt) (*cmtservice.GetLatestBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetLatestBlock, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling TendermintGRPCClient.GetLatestBlock")
	}
	return resp, nil
}

func (c *TendermintGRPCClient) GetLatestValidatorSet(ctx context.Context, req *cmtservice.GetLatestValidatorSetRequest, opts ...config.CallOpt) (*cmtservice.GetLatestValidatorSetResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetLatestValidatorSet, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling TendermintGRPCClient.GetLatestValidatorSet")
	}
	return resp, nil
}

func (c *TendermintGRPCClient) GetNodeInfo(ctx context.Context, req *cmtservice.GetNodeInfoRequest, opts ...config.CallOpt) (*cmtservice.GetNodeInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetNodeInfo, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling TendermintGRPCClient.GetNodeInfo")
	}
	return resp, nil
}

func (c *TendermintGRPCClient) GetSyncing(ctx context.Context, req *cmtservice.GetSyncingRequest, opts ...config.CallOpt) (*cmtservice.GetSyncingResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetSyncing, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling TendermintGRPCClient.GetSyncing")
	}
	return resp, nil
}

func (c *TendermintGRPCClient) GetValidatorSetByHeight(ctx context.Context, req *cmtservice.GetValidatorSetByHeightRequest, opts ...config.CallOpt) (*cmtservice.GetValidatorSetByHeightResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetValidatorSetByHeight, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling TendermintGRPCClient.GetValidatorSetByHeight")
	}
	return resp, nil
}
