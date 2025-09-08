package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/allora-network/allora-sdk-go/config"

	emissionstypes "github.com/allora-network/allora-chain/x/emissions/types"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// EmissionsGRPCClient provides gRPC access to the emissions module
type EmissionsGRPCClient struct {
	client emissionstypes.QueryServiceClient
	logger zerolog.Logger
}

var _ interfaces.EmissionsClient = (*EmissionsGRPCClient)(nil)

// NewEmissionsGRPCClient creates a new emissions REST client
func NewEmissionsGRPCClient(conn *grpc.ClientConn, logger zerolog.Logger) *EmissionsGRPCClient {
	return &EmissionsGRPCClient{
		client: emissionstypes.NewQueryServiceClient(conn),
		logger: logger.With().Str("module", "emissions").Str("protocol", "grpc").Logger(),
	}
}

func (c *EmissionsGRPCClient) GetParams(ctx context.Context, req *emissionstypes.GetParamsRequest, opts ...config.CallOpt) (*emissionstypes.GetParamsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetParams, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetParams")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetNextTopicId(ctx context.Context, req *emissionstypes.GetNextTopicIdRequest, opts ...config.CallOpt) (*emissionstypes.GetNextTopicIdResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetNextTopicId, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetNextTopicId")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetTopic(ctx context.Context, req *emissionstypes.GetTopicRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetTopic, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetTopic")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetWorkerLatestInferenceByTopicId(ctx context.Context, req *emissionstypes.GetWorkerLatestInferenceByTopicIdRequest, opts ...config.CallOpt) (*emissionstypes.GetWorkerLatestInferenceByTopicIdResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetWorkerLatestInferenceByTopicId, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetWorkerLatestInferenceByTopicId")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetInferencesAtBlock(ctx context.Context, req *emissionstypes.GetInferencesAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetInferencesAtBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetInferencesAtBlock, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetInferencesAtBlock")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetLatestTopicInferences(ctx context.Context, req *emissionstypes.GetLatestTopicInferencesRequest, opts ...config.CallOpt) (*emissionstypes.GetLatestTopicInferencesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetLatestTopicInferences, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetLatestTopicInferences")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetForecastsAtBlock(ctx context.Context, req *emissionstypes.GetForecastsAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetForecastsAtBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetForecastsAtBlock, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetForecastsAtBlock")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetNetworkLossBundleAtBlock(ctx context.Context, req *emissionstypes.GetNetworkLossBundleAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetNetworkLossBundleAtBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetNetworkLossBundleAtBlock, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetNetworkLossBundleAtBlock")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetTotalStake(ctx context.Context, req *emissionstypes.GetTotalStakeRequest, opts ...config.CallOpt) (*emissionstypes.GetTotalStakeResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetTotalStake, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetTotalStake")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetReputerStakeInTopic(ctx context.Context, req *emissionstypes.GetReputerStakeInTopicRequest, opts ...config.CallOpt) (*emissionstypes.GetReputerStakeInTopicResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetReputerStakeInTopic, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetReputerStakeInTopic")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetMultiReputerStakeInTopic(ctx context.Context, req *emissionstypes.GetMultiReputerStakeInTopicRequest, opts ...config.CallOpt) (*emissionstypes.GetMultiReputerStakeInTopicResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetMultiReputerStakeInTopic, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetMultiReputerStakeInTopic")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetStakeFromReputerInTopicInSelf(ctx context.Context, req *emissionstypes.GetStakeFromReputerInTopicInSelfRequest, opts ...config.CallOpt) (*emissionstypes.GetStakeFromReputerInTopicInSelfResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetStakeFromReputerInTopicInSelf, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetStakeFromReputerInTopicInSelf")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetDelegateStakeInTopicInReputer(ctx context.Context, req *emissionstypes.GetDelegateStakeInTopicInReputerRequest, opts ...config.CallOpt) (*emissionstypes.GetDelegateStakeInTopicInReputerResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetDelegateStakeInTopicInReputer, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetDelegateStakeInTopicInReputer")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetStakeFromDelegatorInTopicInReputer(ctx context.Context, req *emissionstypes.GetStakeFromDelegatorInTopicInReputerRequest, opts ...config.CallOpt) (*emissionstypes.GetStakeFromDelegatorInTopicInReputerResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetStakeFromDelegatorInTopicInReputer, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetStakeFromDelegatorInTopicInReputer")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetStakeFromDelegatorInTopic(ctx context.Context, req *emissionstypes.GetStakeFromDelegatorInTopicRequest, opts ...config.CallOpt) (*emissionstypes.GetStakeFromDelegatorInTopicResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetStakeFromDelegatorInTopic, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetStakeFromDelegatorInTopic")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetTopicStake(ctx context.Context, req *emissionstypes.GetTopicStakeRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicStakeResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetTopicStake, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetTopicStake")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetStakeRemovalsUpUntilBlock(ctx context.Context, req *emissionstypes.GetStakeRemovalsUpUntilBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetStakeRemovalsUpUntilBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetStakeRemovalsUpUntilBlock, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetStakeRemovalsUpUntilBlock")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetDelegateStakeRemovalsUpUntilBlock(ctx context.Context, req *emissionstypes.GetDelegateStakeRemovalsUpUntilBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetDelegateStakeRemovalsUpUntilBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetDelegateStakeRemovalsUpUntilBlock, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetDelegateStakeRemovalsUpUntilBlock")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetStakeRemovalInfo(ctx context.Context, req *emissionstypes.GetStakeRemovalInfoRequest, opts ...config.CallOpt) (*emissionstypes.GetStakeRemovalInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetStakeRemovalInfo, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetStakeRemovalInfo")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetDelegateStakeRemovalInfo(ctx context.Context, req *emissionstypes.GetDelegateStakeRemovalInfoRequest, opts ...config.CallOpt) (*emissionstypes.GetDelegateStakeRemovalInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetDelegateStakeRemovalInfo, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetDelegateStakeRemovalInfo")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetWorkerNodeInfo(ctx context.Context, req *emissionstypes.GetWorkerNodeInfoRequest, opts ...config.CallOpt) (*emissionstypes.GetWorkerNodeInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetWorkerNodeInfo, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetWorkerNodeInfo")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetReputerNodeInfo(ctx context.Context, req *emissionstypes.GetReputerNodeInfoRequest, opts ...config.CallOpt) (*emissionstypes.GetReputerNodeInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetReputerNodeInfo, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetReputerNodeInfo")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) IsWorkerRegisteredInTopicId(ctx context.Context, req *emissionstypes.IsWorkerRegisteredInTopicIdRequest, opts ...config.CallOpt) (*emissionstypes.IsWorkerRegisteredInTopicIdResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.IsWorkerRegisteredInTopicId, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.IsWorkerRegisteredInTopicId")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) IsReputerRegisteredInTopicId(ctx context.Context, req *emissionstypes.IsReputerRegisteredInTopicIdRequest, opts ...config.CallOpt) (*emissionstypes.IsReputerRegisteredInTopicIdResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.IsReputerRegisteredInTopicId, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.IsReputerRegisteredInTopicId")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetNetworkInferencesAtBlock(ctx context.Context, req *emissionstypes.GetNetworkInferencesAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetNetworkInferencesAtBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetNetworkInferencesAtBlock, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetNetworkInferencesAtBlock")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetNetworkInferencesAtBlockOutlierResistant(ctx context.Context, req *emissionstypes.GetNetworkInferencesAtBlockOutlierResistantRequest, opts ...config.CallOpt) (*emissionstypes.GetNetworkInferencesAtBlockOutlierResistantResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetNetworkInferencesAtBlockOutlierResistant, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetNetworkInferencesAtBlockOutlierResistant")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetLatestNetworkInferences(ctx context.Context, req *emissionstypes.GetLatestNetworkInferencesRequest, opts ...config.CallOpt) (*emissionstypes.GetLatestNetworkInferencesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetLatestNetworkInferences, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetLatestNetworkInferences")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetLatestNetworkInferencesOutlierResistant(ctx context.Context, req *emissionstypes.GetLatestNetworkInferencesOutlierResistantRequest, opts ...config.CallOpt) (*emissionstypes.GetLatestNetworkInferencesOutlierResistantResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetLatestNetworkInferencesOutlierResistant, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetLatestNetworkInferencesOutlierResistant")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) IsWorkerNonceUnfulfilled(ctx context.Context, req *emissionstypes.IsWorkerNonceUnfulfilledRequest, opts ...config.CallOpt) (*emissionstypes.IsWorkerNonceUnfulfilledResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.IsWorkerNonceUnfulfilled, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.IsWorkerNonceUnfulfilled")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) IsReputerNonceUnfulfilled(ctx context.Context, req *emissionstypes.IsReputerNonceUnfulfilledRequest, opts ...config.CallOpt) (*emissionstypes.IsReputerNonceUnfulfilledResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.IsReputerNonceUnfulfilled, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.IsReputerNonceUnfulfilled")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetUnfulfilledWorkerNonces(ctx context.Context, req *emissionstypes.GetUnfulfilledWorkerNoncesRequest, opts ...config.CallOpt) (*emissionstypes.GetUnfulfilledWorkerNoncesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetUnfulfilledWorkerNonces, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetUnfulfilledWorkerNonces")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetUnfulfilledReputerNonces(ctx context.Context, req *emissionstypes.GetUnfulfilledReputerNoncesRequest, opts ...config.CallOpt) (*emissionstypes.GetUnfulfilledReputerNoncesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetUnfulfilledReputerNonces, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetUnfulfilledReputerNonces")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetInfererNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetInfererNetworkRegretResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetInfererNetworkRegret, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetInfererNetworkRegret")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetForecasterNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetForecasterNetworkRegretResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetForecasterNetworkRegret, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetForecasterNetworkRegret")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetOneInForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetOneInForecasterNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetOneInForecasterNetworkRegretResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetOneInForecasterNetworkRegret, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetOneInForecasterNetworkRegret")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) IsWhitelistAdmin(ctx context.Context, req *emissionstypes.IsWhitelistAdminRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistAdminResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.IsWhitelistAdmin, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.IsWhitelistAdmin")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetTopicLastWorkerCommitInfo(ctx context.Context, req *emissionstypes.GetTopicLastWorkerCommitInfoRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicLastWorkerCommitInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetTopicLastWorkerCommitInfo, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetTopicLastWorkerCommitInfo")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetTopicLastReputerCommitInfo(ctx context.Context, req *emissionstypes.GetTopicLastReputerCommitInfoRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicLastReputerCommitInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetTopicLastReputerCommitInfo, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetTopicLastReputerCommitInfo")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetTopicRewardNonce(ctx context.Context, req *emissionstypes.GetTopicRewardNonceRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicRewardNonceResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetTopicRewardNonce, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetTopicRewardNonce")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetReputerLossBundlesAtBlock(ctx context.Context, req *emissionstypes.GetReputerLossBundlesAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetReputerLossBundlesAtBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetReputerLossBundlesAtBlock, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetReputerLossBundlesAtBlock")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetStakeReputerAuthority(ctx context.Context, req *emissionstypes.GetStakeReputerAuthorityRequest, opts ...config.CallOpt) (*emissionstypes.GetStakeReputerAuthorityResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetStakeReputerAuthority, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetStakeReputerAuthority")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetDelegateStakePlacement(ctx context.Context, req *emissionstypes.GetDelegateStakePlacementRequest, opts ...config.CallOpt) (*emissionstypes.GetDelegateStakePlacementResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetDelegateStakePlacement, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetDelegateStakePlacement")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetDelegateStakeUponReputer(ctx context.Context, req *emissionstypes.GetDelegateStakeUponReputerRequest, opts ...config.CallOpt) (*emissionstypes.GetDelegateStakeUponReputerResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetDelegateStakeUponReputer, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetDelegateStakeUponReputer")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetDelegateRewardPerShare(ctx context.Context, req *emissionstypes.GetDelegateRewardPerShareRequest, opts ...config.CallOpt) (*emissionstypes.GetDelegateRewardPerShareResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetDelegateRewardPerShare, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetDelegateRewardPerShare")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetStakeRemovalForReputerAndTopicId(ctx context.Context, req *emissionstypes.GetStakeRemovalForReputerAndTopicIdRequest, opts ...config.CallOpt) (*emissionstypes.GetStakeRemovalForReputerAndTopicIdResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetStakeRemovalForReputerAndTopicId, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetStakeRemovalForReputerAndTopicId")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetDelegateStakeRemoval(ctx context.Context, req *emissionstypes.GetDelegateStakeRemovalRequest, opts ...config.CallOpt) (*emissionstypes.GetDelegateStakeRemovalResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetDelegateStakeRemoval, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetDelegateStakeRemoval")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetPreviousTopicWeight(ctx context.Context, req *emissionstypes.GetPreviousTopicWeightRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousTopicWeightResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetPreviousTopicWeight, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetPreviousTopicWeight")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetTotalSumPreviousTopicWeights(ctx context.Context, req *emissionstypes.GetTotalSumPreviousTopicWeightsRequest, opts ...config.CallOpt) (*emissionstypes.GetTotalSumPreviousTopicWeightsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetTotalSumPreviousTopicWeights, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetTotalSumPreviousTopicWeights")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) TopicExists(ctx context.Context, req *emissionstypes.TopicExistsRequest, opts ...config.CallOpt) (*emissionstypes.TopicExistsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.TopicExists, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.TopicExists")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) IsTopicActive(ctx context.Context, req *emissionstypes.IsTopicActiveRequest, opts ...config.CallOpt) (*emissionstypes.IsTopicActiveResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.IsTopicActive, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.IsTopicActive")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetTopicFeeRevenue(ctx context.Context, req *emissionstypes.GetTopicFeeRevenueRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicFeeRevenueResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetTopicFeeRevenue, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetTopicFeeRevenue")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetInfererScoreEma(ctx context.Context, req *emissionstypes.GetInfererScoreEmaRequest, opts ...config.CallOpt) (*emissionstypes.GetInfererScoreEmaResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetInfererScoreEma, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetInfererScoreEma")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetForecasterScoreEma(ctx context.Context, req *emissionstypes.GetForecasterScoreEmaRequest, opts ...config.CallOpt) (*emissionstypes.GetForecasterScoreEmaResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetForecasterScoreEma, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetForecasterScoreEma")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetReputerScoreEma(ctx context.Context, req *emissionstypes.GetReputerScoreEmaRequest, opts ...config.CallOpt) (*emissionstypes.GetReputerScoreEmaResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetReputerScoreEma, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetReputerScoreEma")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetInferenceScoresUntilBlock(ctx context.Context, req *emissionstypes.GetInferenceScoresUntilBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetInferenceScoresUntilBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetInferenceScoresUntilBlock, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetInferenceScoresUntilBlock")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetPreviousTopicQuantileForecasterScoreEma(ctx context.Context, req *emissionstypes.GetPreviousTopicQuantileForecasterScoreEmaRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousTopicQuantileForecasterScoreEmaResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetPreviousTopicQuantileForecasterScoreEma, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetPreviousTopicQuantileForecasterScoreEma")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetPreviousTopicQuantileInfererScoreEma(ctx context.Context, req *emissionstypes.GetPreviousTopicQuantileInfererScoreEmaRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousTopicQuantileInfererScoreEmaResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetPreviousTopicQuantileInfererScoreEma, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetPreviousTopicQuantileInfererScoreEma")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetPreviousTopicQuantileReputerScoreEma(ctx context.Context, req *emissionstypes.GetPreviousTopicQuantileReputerScoreEmaRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousTopicQuantileReputerScoreEmaResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetPreviousTopicQuantileReputerScoreEma, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetPreviousTopicQuantileReputerScoreEma")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetWorkerInferenceScoresAtBlock(ctx context.Context, req *emissionstypes.GetWorkerInferenceScoresAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetWorkerInferenceScoresAtBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetWorkerInferenceScoresAtBlock, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetWorkerInferenceScoresAtBlock")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetCurrentLowestInfererScore(ctx context.Context, req *emissionstypes.GetCurrentLowestInfererScoreRequest, opts ...config.CallOpt) (*emissionstypes.GetCurrentLowestInfererScoreResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetCurrentLowestInfererScore, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetCurrentLowestInfererScore")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetForecastScoresUntilBlock(ctx context.Context, req *emissionstypes.GetForecastScoresUntilBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetForecastScoresUntilBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetForecastScoresUntilBlock, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetForecastScoresUntilBlock")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetWorkerForecastScoresAtBlock(ctx context.Context, req *emissionstypes.GetWorkerForecastScoresAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetWorkerForecastScoresAtBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetWorkerForecastScoresAtBlock, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetWorkerForecastScoresAtBlock")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetCurrentLowestForecasterScore(ctx context.Context, req *emissionstypes.GetCurrentLowestForecasterScoreRequest, opts ...config.CallOpt) (*emissionstypes.GetCurrentLowestForecasterScoreResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetCurrentLowestForecasterScore, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetCurrentLowestForecasterScore")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetReputersScoresAtBlock(ctx context.Context, req *emissionstypes.GetReputersScoresAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetReputersScoresAtBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetReputersScoresAtBlock, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetReputersScoresAtBlock")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetCurrentLowestReputerScore(ctx context.Context, req *emissionstypes.GetCurrentLowestReputerScoreRequest, opts ...config.CallOpt) (*emissionstypes.GetCurrentLowestReputerScoreResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetCurrentLowestReputerScore, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetCurrentLowestReputerScore")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetListeningCoefficient(ctx context.Context, req *emissionstypes.GetListeningCoefficientRequest, opts ...config.CallOpt) (*emissionstypes.GetListeningCoefficientResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetListeningCoefficient, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetListeningCoefficient")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetPreviousReputerRewardFraction(ctx context.Context, req *emissionstypes.GetPreviousReputerRewardFractionRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousReputerRewardFractionResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetPreviousReputerRewardFraction, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetPreviousReputerRewardFraction")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetPreviousInferenceRewardFraction(ctx context.Context, req *emissionstypes.GetPreviousInferenceRewardFractionRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousInferenceRewardFractionResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetPreviousInferenceRewardFraction, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetPreviousInferenceRewardFraction")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetPreviousForecastRewardFraction(ctx context.Context, req *emissionstypes.GetPreviousForecastRewardFractionRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousForecastRewardFractionResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetPreviousForecastRewardFraction, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetPreviousForecastRewardFraction")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetPreviousPercentageRewardToStakedReputers(ctx context.Context, req *emissionstypes.GetPreviousPercentageRewardToStakedReputersRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousPercentageRewardToStakedReputersResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetPreviousPercentageRewardToStakedReputers, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetPreviousPercentageRewardToStakedReputers")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetTotalRewardToDistribute(ctx context.Context, req *emissionstypes.GetTotalRewardToDistributeRequest, opts ...config.CallOpt) (*emissionstypes.GetTotalRewardToDistributeResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetTotalRewardToDistribute, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetTotalRewardToDistribute")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetNaiveInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetNaiveInfererNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetNaiveInfererNetworkRegretResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetNaiveInfererNetworkRegret, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetNaiveInfererNetworkRegret")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetOneOutInfererInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutInfererInfererNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetOneOutInfererInfererNetworkRegretResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetOneOutInfererInfererNetworkRegret, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetOneOutInfererInfererNetworkRegret")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetOneOutInfererForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutInfererForecasterNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetOneOutInfererForecasterNetworkRegretResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetOneOutInfererForecasterNetworkRegret, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetOneOutInfererForecasterNetworkRegret")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetOneOutForecasterInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutForecasterInfererNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetOneOutForecasterInfererNetworkRegretResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetOneOutForecasterInfererNetworkRegret, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetOneOutForecasterInfererNetworkRegret")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetOneOutForecasterForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutForecasterForecasterNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetOneOutForecasterForecasterNetworkRegretResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetOneOutForecasterForecasterNetworkRegret, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetOneOutForecasterForecasterNetworkRegret")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetActiveTopicsAtBlock(ctx context.Context, req *emissionstypes.GetActiveTopicsAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetActiveTopicsAtBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetActiveTopicsAtBlock, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetActiveTopicsAtBlock")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetNextChurningBlockByTopicId(ctx context.Context, req *emissionstypes.GetNextChurningBlockByTopicIdRequest, opts ...config.CallOpt) (*emissionstypes.GetNextChurningBlockByTopicIdResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetNextChurningBlockByTopicId, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetNextChurningBlockByTopicId")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetCountInfererInclusionsInTopic(ctx context.Context, req *emissionstypes.GetCountInfererInclusionsInTopicRequest, opts ...config.CallOpt) (*emissionstypes.GetCountInfererInclusionsInTopicResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetCountInfererInclusionsInTopic, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetCountInfererInclusionsInTopic")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetCountForecasterInclusionsInTopic(ctx context.Context, req *emissionstypes.GetCountForecasterInclusionsInTopicRequest, opts ...config.CallOpt) (*emissionstypes.GetCountForecasterInclusionsInTopicResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetCountForecasterInclusionsInTopic, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetCountForecasterInclusionsInTopic")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) IsWhitelistedGlobalWorker(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalWorkerRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistedGlobalWorkerResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.IsWhitelistedGlobalWorker, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.IsWhitelistedGlobalWorker")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) IsWhitelistedGlobalReputer(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalReputerRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistedGlobalReputerResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.IsWhitelistedGlobalReputer, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.IsWhitelistedGlobalReputer")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) IsWhitelistedGlobalAdmin(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalAdminRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistedGlobalAdminResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.IsWhitelistedGlobalAdmin, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.IsWhitelistedGlobalAdmin")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) IsTopicWorkerWhitelistEnabled(ctx context.Context, req *emissionstypes.IsTopicWorkerWhitelistEnabledRequest, opts ...config.CallOpt) (*emissionstypes.IsTopicWorkerWhitelistEnabledResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.IsTopicWorkerWhitelistEnabled, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.IsTopicWorkerWhitelistEnabled")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) IsTopicReputerWhitelistEnabled(ctx context.Context, req *emissionstypes.IsTopicReputerWhitelistEnabledRequest, opts ...config.CallOpt) (*emissionstypes.IsTopicReputerWhitelistEnabledResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.IsTopicReputerWhitelistEnabled, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.IsTopicReputerWhitelistEnabled")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) IsWhitelistedTopicCreator(ctx context.Context, req *emissionstypes.IsWhitelistedTopicCreatorRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistedTopicCreatorResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.IsWhitelistedTopicCreator, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.IsWhitelistedTopicCreator")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) IsWhitelistedGlobalActor(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalActorRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistedGlobalActorResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.IsWhitelistedGlobalActor, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.IsWhitelistedGlobalActor")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) IsWhitelistedTopicWorker(ctx context.Context, req *emissionstypes.IsWhitelistedTopicWorkerRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistedTopicWorkerResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.IsWhitelistedTopicWorker, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.IsWhitelistedTopicWorker")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) IsWhitelistedTopicReputer(ctx context.Context, req *emissionstypes.IsWhitelistedTopicReputerRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistedTopicReputerResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.IsWhitelistedTopicReputer, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.IsWhitelistedTopicReputer")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) CanUpdateAllGlobalWhitelists(ctx context.Context, req *emissionstypes.CanUpdateAllGlobalWhitelistsRequest, opts ...config.CallOpt) (*emissionstypes.CanUpdateAllGlobalWhitelistsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.CanUpdateAllGlobalWhitelists, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.CanUpdateAllGlobalWhitelists")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) CanUpdateGlobalWorkerWhitelist(ctx context.Context, req *emissionstypes.CanUpdateGlobalWorkerWhitelistRequest, opts ...config.CallOpt) (*emissionstypes.CanUpdateGlobalWorkerWhitelistResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.CanUpdateGlobalWorkerWhitelist, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.CanUpdateGlobalWorkerWhitelist")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) CanUpdateGlobalReputerWhitelist(ctx context.Context, req *emissionstypes.CanUpdateGlobalReputerWhitelistRequest, opts ...config.CallOpt) (*emissionstypes.CanUpdateGlobalReputerWhitelistResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.CanUpdateGlobalReputerWhitelist, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.CanUpdateGlobalReputerWhitelist")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) CanUpdateParams(ctx context.Context, req *emissionstypes.CanUpdateParamsRequest, opts ...config.CallOpt) (*emissionstypes.CanUpdateParamsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.CanUpdateParams, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.CanUpdateParams")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) CanUpdateTopicWhitelist(ctx context.Context, req *emissionstypes.CanUpdateTopicWhitelistRequest, opts ...config.CallOpt) (*emissionstypes.CanUpdateTopicWhitelistResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.CanUpdateTopicWhitelist, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.CanUpdateTopicWhitelist")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) CanCreateTopic(ctx context.Context, req *emissionstypes.CanCreateTopicRequest, opts ...config.CallOpt) (*emissionstypes.CanCreateTopicResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.CanCreateTopic, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.CanCreateTopic")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) CanSubmitWorkerPayload(ctx context.Context, req *emissionstypes.CanSubmitWorkerPayloadRequest, opts ...config.CallOpt) (*emissionstypes.CanSubmitWorkerPayloadResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.CanSubmitWorkerPayload, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.CanSubmitWorkerPayload")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) CanSubmitReputerPayload(ctx context.Context, req *emissionstypes.CanSubmitReputerPayloadRequest, opts ...config.CallOpt) (*emissionstypes.CanSubmitReputerPayloadResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.CanSubmitReputerPayload, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.CanSubmitReputerPayload")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetTopicInitialInfererEmaScore(ctx context.Context, req *emissionstypes.GetTopicInitialInfererEmaScoreRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicInitialInfererEmaScoreResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetTopicInitialInfererEmaScore, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetTopicInitialInfererEmaScore")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetTopicInitialForecasterEmaScore(ctx context.Context, req *emissionstypes.GetTopicInitialForecasterEmaScoreRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicInitialForecasterEmaScoreResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetTopicInitialForecasterEmaScore, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetTopicInitialForecasterEmaScore")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetTopicInitialReputerEmaScore(ctx context.Context, req *emissionstypes.GetTopicInitialReputerEmaScoreRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicInitialReputerEmaScoreResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetTopicInitialReputerEmaScore, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetTopicInitialReputerEmaScore")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetLatestRegretStdNorm(ctx context.Context, req *emissionstypes.GetLatestRegretStdNormRequest, opts ...config.CallOpt) (*emissionstypes.GetLatestRegretStdNormResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetLatestRegretStdNorm, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetLatestRegretStdNorm")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetLatestInfererWeight(ctx context.Context, req *emissionstypes.GetLatestInfererWeightRequest, opts ...config.CallOpt) (*emissionstypes.GetLatestInfererWeightResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetLatestInfererWeight, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetLatestInfererWeight")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetLatestForecasterWeight(ctx context.Context, req *emissionstypes.GetLatestForecasterWeightRequest, opts ...config.CallOpt) (*emissionstypes.GetLatestForecasterWeightResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetLatestForecasterWeight, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetLatestForecasterWeight")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetWorkerSubmissionWindowStatus(ctx context.Context, req *emissionstypes.GetWorkerSubmissionWindowStatusRequest, opts ...config.CallOpt) (*emissionstypes.GetWorkerSubmissionWindowStatusResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetWorkerSubmissionWindowStatus, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetWorkerSubmissionWindowStatus")
	}
	return resp, nil
}

func (c *EmissionsGRPCClient) GetReputerSubmissionWindowStatus(ctx context.Context, req *emissionstypes.GetReputerSubmissionWindowStatusRequest, opts ...config.CallOpt) (*emissionstypes.GetReputerSubmissionWindowStatusResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp, err := queryWithHeight(ctx, callOpts.Height, c.client.GetReputerSubmissionWindowStatus, req)
	if err != nil {
		return resp, errors.Wrap(err, "while calling EmissionsGRPCClient.GetReputerSubmissionWindowStatus")
	}
	return resp, nil
}
