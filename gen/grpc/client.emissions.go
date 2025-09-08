package grpc

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

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

func (c *EmissionsGRPCClient) GetParams(ctx context.Context, req *emissionstypes.GetParamsRequest) (*emissionstypes.GetParamsResponse, error) {
	resp, err := c.client.GetParams(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetParams")
}

func (c *EmissionsGRPCClient) GetNextTopicId(ctx context.Context, req *emissionstypes.GetNextTopicIdRequest) (*emissionstypes.GetNextTopicIdResponse, error) {
	resp, err := c.client.GetNextTopicId(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetNextTopicId")
}

func (c *EmissionsGRPCClient) GetTopic(ctx context.Context, req *emissionstypes.GetTopicRequest) (*emissionstypes.GetTopicResponse, error) {
	resp, err := c.client.GetTopic(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetTopic")
}

func (c *EmissionsGRPCClient) GetWorkerLatestInferenceByTopicId(ctx context.Context, req *emissionstypes.GetWorkerLatestInferenceByTopicIdRequest) (*emissionstypes.GetWorkerLatestInferenceByTopicIdResponse, error) {
	resp, err := c.client.GetWorkerLatestInferenceByTopicId(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetWorkerLatestInferenceByTopicId")
}

func (c *EmissionsGRPCClient) GetInferencesAtBlock(ctx context.Context, req *emissionstypes.GetInferencesAtBlockRequest) (*emissionstypes.GetInferencesAtBlockResponse, error) {
	resp, err := c.client.GetInferencesAtBlock(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetInferencesAtBlock")
}

func (c *EmissionsGRPCClient) GetLatestTopicInferences(ctx context.Context, req *emissionstypes.GetLatestTopicInferencesRequest) (*emissionstypes.GetLatestTopicInferencesResponse, error) {
	resp, err := c.client.GetLatestTopicInferences(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetLatestTopicInferences")
}

func (c *EmissionsGRPCClient) GetForecastsAtBlock(ctx context.Context, req *emissionstypes.GetForecastsAtBlockRequest) (*emissionstypes.GetForecastsAtBlockResponse, error) {
	resp, err := c.client.GetForecastsAtBlock(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetForecastsAtBlock")
}

func (c *EmissionsGRPCClient) GetNetworkLossBundleAtBlock(ctx context.Context, req *emissionstypes.GetNetworkLossBundleAtBlockRequest) (*emissionstypes.GetNetworkLossBundleAtBlockResponse, error) {
	resp, err := c.client.GetNetworkLossBundleAtBlock(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetNetworkLossBundleAtBlock")
}

func (c *EmissionsGRPCClient) GetTotalStake(ctx context.Context, req *emissionstypes.GetTotalStakeRequest) (*emissionstypes.GetTotalStakeResponse, error) {
	resp, err := c.client.GetTotalStake(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetTotalStake")
}

func (c *EmissionsGRPCClient) GetReputerStakeInTopic(ctx context.Context, req *emissionstypes.GetReputerStakeInTopicRequest) (*emissionstypes.GetReputerStakeInTopicResponse, error) {
	resp, err := c.client.GetReputerStakeInTopic(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetReputerStakeInTopic")
}

func (c *EmissionsGRPCClient) GetMultiReputerStakeInTopic(ctx context.Context, req *emissionstypes.GetMultiReputerStakeInTopicRequest) (*emissionstypes.GetMultiReputerStakeInTopicResponse, error) {
	resp, err := c.client.GetMultiReputerStakeInTopic(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetMultiReputerStakeInTopic")
}

func (c *EmissionsGRPCClient) GetStakeFromReputerInTopicInSelf(ctx context.Context, req *emissionstypes.GetStakeFromReputerInTopicInSelfRequest) (*emissionstypes.GetStakeFromReputerInTopicInSelfResponse, error) {
	resp, err := c.client.GetStakeFromReputerInTopicInSelf(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetStakeFromReputerInTopicInSelf")
}

func (c *EmissionsGRPCClient) GetDelegateStakeInTopicInReputer(ctx context.Context, req *emissionstypes.GetDelegateStakeInTopicInReputerRequest) (*emissionstypes.GetDelegateStakeInTopicInReputerResponse, error) {
	resp, err := c.client.GetDelegateStakeInTopicInReputer(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetDelegateStakeInTopicInReputer")
}

func (c *EmissionsGRPCClient) GetStakeFromDelegatorInTopicInReputer(ctx context.Context, req *emissionstypes.GetStakeFromDelegatorInTopicInReputerRequest) (*emissionstypes.GetStakeFromDelegatorInTopicInReputerResponse, error) {
	resp, err := c.client.GetStakeFromDelegatorInTopicInReputer(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetStakeFromDelegatorInTopicInReputer")
}

func (c *EmissionsGRPCClient) GetStakeFromDelegatorInTopic(ctx context.Context, req *emissionstypes.GetStakeFromDelegatorInTopicRequest) (*emissionstypes.GetStakeFromDelegatorInTopicResponse, error) {
	resp, err := c.client.GetStakeFromDelegatorInTopic(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetStakeFromDelegatorInTopic")
}

func (c *EmissionsGRPCClient) GetTopicStake(ctx context.Context, req *emissionstypes.GetTopicStakeRequest) (*emissionstypes.GetTopicStakeResponse, error) {
	resp, err := c.client.GetTopicStake(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetTopicStake")
}

func (c *EmissionsGRPCClient) GetStakeRemovalsUpUntilBlock(ctx context.Context, req *emissionstypes.GetStakeRemovalsUpUntilBlockRequest) (*emissionstypes.GetStakeRemovalsUpUntilBlockResponse, error) {
	resp, err := c.client.GetStakeRemovalsUpUntilBlock(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetStakeRemovalsUpUntilBlock")
}

func (c *EmissionsGRPCClient) GetDelegateStakeRemovalsUpUntilBlock(ctx context.Context, req *emissionstypes.GetDelegateStakeRemovalsUpUntilBlockRequest) (*emissionstypes.GetDelegateStakeRemovalsUpUntilBlockResponse, error) {
	resp, err := c.client.GetDelegateStakeRemovalsUpUntilBlock(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetDelegateStakeRemovalsUpUntilBlock")
}

func (c *EmissionsGRPCClient) GetStakeRemovalInfo(ctx context.Context, req *emissionstypes.GetStakeRemovalInfoRequest) (*emissionstypes.GetStakeRemovalInfoResponse, error) {
	resp, err := c.client.GetStakeRemovalInfo(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetStakeRemovalInfo")
}

func (c *EmissionsGRPCClient) GetDelegateStakeRemovalInfo(ctx context.Context, req *emissionstypes.GetDelegateStakeRemovalInfoRequest) (*emissionstypes.GetDelegateStakeRemovalInfoResponse, error) {
	resp, err := c.client.GetDelegateStakeRemovalInfo(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetDelegateStakeRemovalInfo")
}

func (c *EmissionsGRPCClient) GetWorkerNodeInfo(ctx context.Context, req *emissionstypes.GetWorkerNodeInfoRequest) (*emissionstypes.GetWorkerNodeInfoResponse, error) {
	resp, err := c.client.GetWorkerNodeInfo(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetWorkerNodeInfo")
}

func (c *EmissionsGRPCClient) GetReputerNodeInfo(ctx context.Context, req *emissionstypes.GetReputerNodeInfoRequest) (*emissionstypes.GetReputerNodeInfoResponse, error) {
	resp, err := c.client.GetReputerNodeInfo(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetReputerNodeInfo")
}

func (c *EmissionsGRPCClient) IsWorkerRegisteredInTopicId(ctx context.Context, req *emissionstypes.IsWorkerRegisteredInTopicIdRequest) (*emissionstypes.IsWorkerRegisteredInTopicIdResponse, error) {
	resp, err := c.client.IsWorkerRegisteredInTopicId(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.IsWorkerRegisteredInTopicId")
}

func (c *EmissionsGRPCClient) IsReputerRegisteredInTopicId(ctx context.Context, req *emissionstypes.IsReputerRegisteredInTopicIdRequest) (*emissionstypes.IsReputerRegisteredInTopicIdResponse, error) {
	resp, err := c.client.IsReputerRegisteredInTopicId(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.IsReputerRegisteredInTopicId")
}

func (c *EmissionsGRPCClient) GetNetworkInferencesAtBlock(ctx context.Context, req *emissionstypes.GetNetworkInferencesAtBlockRequest) (*emissionstypes.GetNetworkInferencesAtBlockResponse, error) {
	resp, err := c.client.GetNetworkInferencesAtBlock(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetNetworkInferencesAtBlock")
}

func (c *EmissionsGRPCClient) GetNetworkInferencesAtBlockOutlierResistant(ctx context.Context, req *emissionstypes.GetNetworkInferencesAtBlockOutlierResistantRequest) (*emissionstypes.GetNetworkInferencesAtBlockOutlierResistantResponse, error) {
	resp, err := c.client.GetNetworkInferencesAtBlockOutlierResistant(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetNetworkInferencesAtBlockOutlierResistant")
}

func (c *EmissionsGRPCClient) GetLatestNetworkInferences(ctx context.Context, req *emissionstypes.GetLatestNetworkInferencesRequest) (*emissionstypes.GetLatestNetworkInferencesResponse, error) {
	resp, err := c.client.GetLatestNetworkInferences(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetLatestNetworkInferences")
}

func (c *EmissionsGRPCClient) GetLatestNetworkInferencesOutlierResistant(ctx context.Context, req *emissionstypes.GetLatestNetworkInferencesOutlierResistantRequest) (*emissionstypes.GetLatestNetworkInferencesOutlierResistantResponse, error) {
	resp, err := c.client.GetLatestNetworkInferencesOutlierResistant(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetLatestNetworkInferencesOutlierResistant")
}

func (c *EmissionsGRPCClient) IsWorkerNonceUnfulfilled(ctx context.Context, req *emissionstypes.IsWorkerNonceUnfulfilledRequest) (*emissionstypes.IsWorkerNonceUnfulfilledResponse, error) {
	resp, err := c.client.IsWorkerNonceUnfulfilled(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.IsWorkerNonceUnfulfilled")
}

func (c *EmissionsGRPCClient) IsReputerNonceUnfulfilled(ctx context.Context, req *emissionstypes.IsReputerNonceUnfulfilledRequest) (*emissionstypes.IsReputerNonceUnfulfilledResponse, error) {
	resp, err := c.client.IsReputerNonceUnfulfilled(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.IsReputerNonceUnfulfilled")
}

func (c *EmissionsGRPCClient) GetUnfulfilledWorkerNonces(ctx context.Context, req *emissionstypes.GetUnfulfilledWorkerNoncesRequest) (*emissionstypes.GetUnfulfilledWorkerNoncesResponse, error) {
	resp, err := c.client.GetUnfulfilledWorkerNonces(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetUnfulfilledWorkerNonces")
}

func (c *EmissionsGRPCClient) GetUnfulfilledReputerNonces(ctx context.Context, req *emissionstypes.GetUnfulfilledReputerNoncesRequest) (*emissionstypes.GetUnfulfilledReputerNoncesResponse, error) {
	resp, err := c.client.GetUnfulfilledReputerNonces(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetUnfulfilledReputerNonces")
}

func (c *EmissionsGRPCClient) GetInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetInfererNetworkRegretRequest) (*emissionstypes.GetInfererNetworkRegretResponse, error) {
	resp, err := c.client.GetInfererNetworkRegret(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetInfererNetworkRegret")
}

func (c *EmissionsGRPCClient) GetForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetForecasterNetworkRegretRequest) (*emissionstypes.GetForecasterNetworkRegretResponse, error) {
	resp, err := c.client.GetForecasterNetworkRegret(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetForecasterNetworkRegret")
}

func (c *EmissionsGRPCClient) GetOneInForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetOneInForecasterNetworkRegretRequest) (*emissionstypes.GetOneInForecasterNetworkRegretResponse, error) {
	resp, err := c.client.GetOneInForecasterNetworkRegret(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetOneInForecasterNetworkRegret")
}

func (c *EmissionsGRPCClient) IsWhitelistAdmin(ctx context.Context, req *emissionstypes.IsWhitelistAdminRequest) (*emissionstypes.IsWhitelistAdminResponse, error) {
	resp, err := c.client.IsWhitelistAdmin(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.IsWhitelistAdmin")
}

func (c *EmissionsGRPCClient) GetTopicLastWorkerCommitInfo(ctx context.Context, req *emissionstypes.GetTopicLastWorkerCommitInfoRequest) (*emissionstypes.GetTopicLastWorkerCommitInfoResponse, error) {
	resp, err := c.client.GetTopicLastWorkerCommitInfo(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetTopicLastWorkerCommitInfo")
}

func (c *EmissionsGRPCClient) GetTopicLastReputerCommitInfo(ctx context.Context, req *emissionstypes.GetTopicLastReputerCommitInfoRequest) (*emissionstypes.GetTopicLastReputerCommitInfoResponse, error) {
	resp, err := c.client.GetTopicLastReputerCommitInfo(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetTopicLastReputerCommitInfo")
}

func (c *EmissionsGRPCClient) GetTopicRewardNonce(ctx context.Context, req *emissionstypes.GetTopicRewardNonceRequest) (*emissionstypes.GetTopicRewardNonceResponse, error) {
	resp, err := c.client.GetTopicRewardNonce(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetTopicRewardNonce")
}

func (c *EmissionsGRPCClient) GetReputerLossBundlesAtBlock(ctx context.Context, req *emissionstypes.GetReputerLossBundlesAtBlockRequest) (*emissionstypes.GetReputerLossBundlesAtBlockResponse, error) {
	resp, err := c.client.GetReputerLossBundlesAtBlock(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetReputerLossBundlesAtBlock")
}

func (c *EmissionsGRPCClient) GetStakeReputerAuthority(ctx context.Context, req *emissionstypes.GetStakeReputerAuthorityRequest) (*emissionstypes.GetStakeReputerAuthorityResponse, error) {
	resp, err := c.client.GetStakeReputerAuthority(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetStakeReputerAuthority")
}

func (c *EmissionsGRPCClient) GetDelegateStakePlacement(ctx context.Context, req *emissionstypes.GetDelegateStakePlacementRequest) (*emissionstypes.GetDelegateStakePlacementResponse, error) {
	resp, err := c.client.GetDelegateStakePlacement(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetDelegateStakePlacement")
}

func (c *EmissionsGRPCClient) GetDelegateStakeUponReputer(ctx context.Context, req *emissionstypes.GetDelegateStakeUponReputerRequest) (*emissionstypes.GetDelegateStakeUponReputerResponse, error) {
	resp, err := c.client.GetDelegateStakeUponReputer(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetDelegateStakeUponReputer")
}

func (c *EmissionsGRPCClient) GetDelegateRewardPerShare(ctx context.Context, req *emissionstypes.GetDelegateRewardPerShareRequest) (*emissionstypes.GetDelegateRewardPerShareResponse, error) {
	resp, err := c.client.GetDelegateRewardPerShare(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetDelegateRewardPerShare")
}

func (c *EmissionsGRPCClient) GetStakeRemovalForReputerAndTopicId(ctx context.Context, req *emissionstypes.GetStakeRemovalForReputerAndTopicIdRequest) (*emissionstypes.GetStakeRemovalForReputerAndTopicIdResponse, error) {
	resp, err := c.client.GetStakeRemovalForReputerAndTopicId(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetStakeRemovalForReputerAndTopicId")
}

func (c *EmissionsGRPCClient) GetDelegateStakeRemoval(ctx context.Context, req *emissionstypes.GetDelegateStakeRemovalRequest) (*emissionstypes.GetDelegateStakeRemovalResponse, error) {
	resp, err := c.client.GetDelegateStakeRemoval(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetDelegateStakeRemoval")
}

func (c *EmissionsGRPCClient) GetPreviousTopicWeight(ctx context.Context, req *emissionstypes.GetPreviousTopicWeightRequest) (*emissionstypes.GetPreviousTopicWeightResponse, error) {
	resp, err := c.client.GetPreviousTopicWeight(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetPreviousTopicWeight")
}

func (c *EmissionsGRPCClient) GetTotalSumPreviousTopicWeights(ctx context.Context, req *emissionstypes.GetTotalSumPreviousTopicWeightsRequest) (*emissionstypes.GetTotalSumPreviousTopicWeightsResponse, error) {
	resp, err := c.client.GetTotalSumPreviousTopicWeights(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetTotalSumPreviousTopicWeights")
}

func (c *EmissionsGRPCClient) TopicExists(ctx context.Context, req *emissionstypes.TopicExistsRequest) (*emissionstypes.TopicExistsResponse, error) {
	resp, err := c.client.TopicExists(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.TopicExists")
}

func (c *EmissionsGRPCClient) IsTopicActive(ctx context.Context, req *emissionstypes.IsTopicActiveRequest) (*emissionstypes.IsTopicActiveResponse, error) {
	resp, err := c.client.IsTopicActive(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.IsTopicActive")
}

func (c *EmissionsGRPCClient) GetTopicFeeRevenue(ctx context.Context, req *emissionstypes.GetTopicFeeRevenueRequest) (*emissionstypes.GetTopicFeeRevenueResponse, error) {
	resp, err := c.client.GetTopicFeeRevenue(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetTopicFeeRevenue")
}

func (c *EmissionsGRPCClient) GetInfererScoreEma(ctx context.Context, req *emissionstypes.GetInfererScoreEmaRequest) (*emissionstypes.GetInfererScoreEmaResponse, error) {
	resp, err := c.client.GetInfererScoreEma(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetInfererScoreEma")
}

func (c *EmissionsGRPCClient) GetForecasterScoreEma(ctx context.Context, req *emissionstypes.GetForecasterScoreEmaRequest) (*emissionstypes.GetForecasterScoreEmaResponse, error) {
	resp, err := c.client.GetForecasterScoreEma(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetForecasterScoreEma")
}

func (c *EmissionsGRPCClient) GetReputerScoreEma(ctx context.Context, req *emissionstypes.GetReputerScoreEmaRequest) (*emissionstypes.GetReputerScoreEmaResponse, error) {
	resp, err := c.client.GetReputerScoreEma(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetReputerScoreEma")
}

func (c *EmissionsGRPCClient) GetInferenceScoresUntilBlock(ctx context.Context, req *emissionstypes.GetInferenceScoresUntilBlockRequest) (*emissionstypes.GetInferenceScoresUntilBlockResponse, error) {
	resp, err := c.client.GetInferenceScoresUntilBlock(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetInferenceScoresUntilBlock")
}

func (c *EmissionsGRPCClient) GetPreviousTopicQuantileForecasterScoreEma(ctx context.Context, req *emissionstypes.GetPreviousTopicQuantileForecasterScoreEmaRequest) (*emissionstypes.GetPreviousTopicQuantileForecasterScoreEmaResponse, error) {
	resp, err := c.client.GetPreviousTopicQuantileForecasterScoreEma(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetPreviousTopicQuantileForecasterScoreEma")
}

func (c *EmissionsGRPCClient) GetPreviousTopicQuantileInfererScoreEma(ctx context.Context, req *emissionstypes.GetPreviousTopicQuantileInfererScoreEmaRequest) (*emissionstypes.GetPreviousTopicQuantileInfererScoreEmaResponse, error) {
	resp, err := c.client.GetPreviousTopicQuantileInfererScoreEma(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetPreviousTopicQuantileInfererScoreEma")
}

func (c *EmissionsGRPCClient) GetPreviousTopicQuantileReputerScoreEma(ctx context.Context, req *emissionstypes.GetPreviousTopicQuantileReputerScoreEmaRequest) (*emissionstypes.GetPreviousTopicQuantileReputerScoreEmaResponse, error) {
	resp, err := c.client.GetPreviousTopicQuantileReputerScoreEma(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetPreviousTopicQuantileReputerScoreEma")
}

func (c *EmissionsGRPCClient) GetWorkerInferenceScoresAtBlock(ctx context.Context, req *emissionstypes.GetWorkerInferenceScoresAtBlockRequest) (*emissionstypes.GetWorkerInferenceScoresAtBlockResponse, error) {
	resp, err := c.client.GetWorkerInferenceScoresAtBlock(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetWorkerInferenceScoresAtBlock")
}

func (c *EmissionsGRPCClient) GetCurrentLowestInfererScore(ctx context.Context, req *emissionstypes.GetCurrentLowestInfererScoreRequest) (*emissionstypes.GetCurrentLowestInfererScoreResponse, error) {
	resp, err := c.client.GetCurrentLowestInfererScore(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetCurrentLowestInfererScore")
}

func (c *EmissionsGRPCClient) GetForecastScoresUntilBlock(ctx context.Context, req *emissionstypes.GetForecastScoresUntilBlockRequest) (*emissionstypes.GetForecastScoresUntilBlockResponse, error) {
	resp, err := c.client.GetForecastScoresUntilBlock(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetForecastScoresUntilBlock")
}

func (c *EmissionsGRPCClient) GetWorkerForecastScoresAtBlock(ctx context.Context, req *emissionstypes.GetWorkerForecastScoresAtBlockRequest) (*emissionstypes.GetWorkerForecastScoresAtBlockResponse, error) {
	resp, err := c.client.GetWorkerForecastScoresAtBlock(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetWorkerForecastScoresAtBlock")
}

func (c *EmissionsGRPCClient) GetCurrentLowestForecasterScore(ctx context.Context, req *emissionstypes.GetCurrentLowestForecasterScoreRequest) (*emissionstypes.GetCurrentLowestForecasterScoreResponse, error) {
	resp, err := c.client.GetCurrentLowestForecasterScore(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetCurrentLowestForecasterScore")
}

func (c *EmissionsGRPCClient) GetReputersScoresAtBlock(ctx context.Context, req *emissionstypes.GetReputersScoresAtBlockRequest) (*emissionstypes.GetReputersScoresAtBlockResponse, error) {
	resp, err := c.client.GetReputersScoresAtBlock(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetReputersScoresAtBlock")
}

func (c *EmissionsGRPCClient) GetCurrentLowestReputerScore(ctx context.Context, req *emissionstypes.GetCurrentLowestReputerScoreRequest) (*emissionstypes.GetCurrentLowestReputerScoreResponse, error) {
	resp, err := c.client.GetCurrentLowestReputerScore(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetCurrentLowestReputerScore")
}

func (c *EmissionsGRPCClient) GetListeningCoefficient(ctx context.Context, req *emissionstypes.GetListeningCoefficientRequest) (*emissionstypes.GetListeningCoefficientResponse, error) {
	resp, err := c.client.GetListeningCoefficient(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetListeningCoefficient")
}

func (c *EmissionsGRPCClient) GetPreviousReputerRewardFraction(ctx context.Context, req *emissionstypes.GetPreviousReputerRewardFractionRequest) (*emissionstypes.GetPreviousReputerRewardFractionResponse, error) {
	resp, err := c.client.GetPreviousReputerRewardFraction(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetPreviousReputerRewardFraction")
}

func (c *EmissionsGRPCClient) GetPreviousInferenceRewardFraction(ctx context.Context, req *emissionstypes.GetPreviousInferenceRewardFractionRequest) (*emissionstypes.GetPreviousInferenceRewardFractionResponse, error) {
	resp, err := c.client.GetPreviousInferenceRewardFraction(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetPreviousInferenceRewardFraction")
}

func (c *EmissionsGRPCClient) GetPreviousForecastRewardFraction(ctx context.Context, req *emissionstypes.GetPreviousForecastRewardFractionRequest) (*emissionstypes.GetPreviousForecastRewardFractionResponse, error) {
	resp, err := c.client.GetPreviousForecastRewardFraction(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetPreviousForecastRewardFraction")
}

func (c *EmissionsGRPCClient) GetPreviousPercentageRewardToStakedReputers(ctx context.Context, req *emissionstypes.GetPreviousPercentageRewardToStakedReputersRequest) (*emissionstypes.GetPreviousPercentageRewardToStakedReputersResponse, error) {
	resp, err := c.client.GetPreviousPercentageRewardToStakedReputers(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetPreviousPercentageRewardToStakedReputers")
}

func (c *EmissionsGRPCClient) GetTotalRewardToDistribute(ctx context.Context, req *emissionstypes.GetTotalRewardToDistributeRequest) (*emissionstypes.GetTotalRewardToDistributeResponse, error) {
	resp, err := c.client.GetTotalRewardToDistribute(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetTotalRewardToDistribute")
}

func (c *EmissionsGRPCClient) GetNaiveInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetNaiveInfererNetworkRegretRequest) (*emissionstypes.GetNaiveInfererNetworkRegretResponse, error) {
	resp, err := c.client.GetNaiveInfererNetworkRegret(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetNaiveInfererNetworkRegret")
}

func (c *EmissionsGRPCClient) GetOneOutInfererInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutInfererInfererNetworkRegretRequest) (*emissionstypes.GetOneOutInfererInfererNetworkRegretResponse, error) {
	resp, err := c.client.GetOneOutInfererInfererNetworkRegret(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetOneOutInfererInfererNetworkRegret")
}

func (c *EmissionsGRPCClient) GetOneOutInfererForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutInfererForecasterNetworkRegretRequest) (*emissionstypes.GetOneOutInfererForecasterNetworkRegretResponse, error) {
	resp, err := c.client.GetOneOutInfererForecasterNetworkRegret(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetOneOutInfererForecasterNetworkRegret")
}

func (c *EmissionsGRPCClient) GetOneOutForecasterInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutForecasterInfererNetworkRegretRequest) (*emissionstypes.GetOneOutForecasterInfererNetworkRegretResponse, error) {
	resp, err := c.client.GetOneOutForecasterInfererNetworkRegret(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetOneOutForecasterInfererNetworkRegret")
}

func (c *EmissionsGRPCClient) GetOneOutForecasterForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutForecasterForecasterNetworkRegretRequest) (*emissionstypes.GetOneOutForecasterForecasterNetworkRegretResponse, error) {
	resp, err := c.client.GetOneOutForecasterForecasterNetworkRegret(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetOneOutForecasterForecasterNetworkRegret")
}

func (c *EmissionsGRPCClient) GetActiveTopicsAtBlock(ctx context.Context, req *emissionstypes.GetActiveTopicsAtBlockRequest) (*emissionstypes.GetActiveTopicsAtBlockResponse, error) {
	resp, err := c.client.GetActiveTopicsAtBlock(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetActiveTopicsAtBlock")
}

func (c *EmissionsGRPCClient) GetNextChurningBlockByTopicId(ctx context.Context, req *emissionstypes.GetNextChurningBlockByTopicIdRequest) (*emissionstypes.GetNextChurningBlockByTopicIdResponse, error) {
	resp, err := c.client.GetNextChurningBlockByTopicId(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetNextChurningBlockByTopicId")
}

func (c *EmissionsGRPCClient) GetCountInfererInclusionsInTopic(ctx context.Context, req *emissionstypes.GetCountInfererInclusionsInTopicRequest) (*emissionstypes.GetCountInfererInclusionsInTopicResponse, error) {
	resp, err := c.client.GetCountInfererInclusionsInTopic(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetCountInfererInclusionsInTopic")
}

func (c *EmissionsGRPCClient) GetCountForecasterInclusionsInTopic(ctx context.Context, req *emissionstypes.GetCountForecasterInclusionsInTopicRequest) (*emissionstypes.GetCountForecasterInclusionsInTopicResponse, error) {
	resp, err := c.client.GetCountForecasterInclusionsInTopic(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetCountForecasterInclusionsInTopic")
}

func (c *EmissionsGRPCClient) GetActiveReputersForTopic(ctx context.Context, req *emissionstypes.GetActiveReputersForTopicRequest) (*emissionstypes.GetActiveReputersForTopicResponse, error) {
	resp, err := c.client.GetActiveReputersForTopic(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetActiveReputersForTopic")
}

func (c *EmissionsGRPCClient) GetActiveForecastersForTopic(ctx context.Context, req *emissionstypes.GetActiveForecastersForTopicRequest) (*emissionstypes.GetActiveForecastersForTopicResponse, error) {
	resp, err := c.client.GetActiveForecastersForTopic(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetActiveForecastersForTopic")
}

func (c *EmissionsGRPCClient) GetActiveInferersForTopic(ctx context.Context, req *emissionstypes.GetActiveInferersForTopicRequest) (*emissionstypes.GetActiveInferersForTopicResponse, error) {
	resp, err := c.client.GetActiveInferersForTopic(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetActiveInferersForTopic")
}

func (c *EmissionsGRPCClient) IsWhitelistedGlobalWorker(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalWorkerRequest) (*emissionstypes.IsWhitelistedGlobalWorkerResponse, error) {
	resp, err := c.client.IsWhitelistedGlobalWorker(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.IsWhitelistedGlobalWorker")
}

func (c *EmissionsGRPCClient) IsWhitelistedGlobalReputer(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalReputerRequest) (*emissionstypes.IsWhitelistedGlobalReputerResponse, error) {
	resp, err := c.client.IsWhitelistedGlobalReputer(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.IsWhitelistedGlobalReputer")
}

func (c *EmissionsGRPCClient) IsWhitelistedGlobalAdmin(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalAdminRequest) (*emissionstypes.IsWhitelistedGlobalAdminResponse, error) {
	resp, err := c.client.IsWhitelistedGlobalAdmin(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.IsWhitelistedGlobalAdmin")
}

func (c *EmissionsGRPCClient) IsTopicWorkerWhitelistEnabled(ctx context.Context, req *emissionstypes.IsTopicWorkerWhitelistEnabledRequest) (*emissionstypes.IsTopicWorkerWhitelistEnabledResponse, error) {
	resp, err := c.client.IsTopicWorkerWhitelistEnabled(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.IsTopicWorkerWhitelistEnabled")
}

func (c *EmissionsGRPCClient) IsTopicReputerWhitelistEnabled(ctx context.Context, req *emissionstypes.IsTopicReputerWhitelistEnabledRequest) (*emissionstypes.IsTopicReputerWhitelistEnabledResponse, error) {
	resp, err := c.client.IsTopicReputerWhitelistEnabled(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.IsTopicReputerWhitelistEnabled")
}

func (c *EmissionsGRPCClient) IsWhitelistedTopicCreator(ctx context.Context, req *emissionstypes.IsWhitelistedTopicCreatorRequest) (*emissionstypes.IsWhitelistedTopicCreatorResponse, error) {
	resp, err := c.client.IsWhitelistedTopicCreator(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.IsWhitelistedTopicCreator")
}

func (c *EmissionsGRPCClient) IsWhitelistedGlobalActor(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalActorRequest) (*emissionstypes.IsWhitelistedGlobalActorResponse, error) {
	resp, err := c.client.IsWhitelistedGlobalActor(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.IsWhitelistedGlobalActor")
}

func (c *EmissionsGRPCClient) IsWhitelistedTopicWorker(ctx context.Context, req *emissionstypes.IsWhitelistedTopicWorkerRequest) (*emissionstypes.IsWhitelistedTopicWorkerResponse, error) {
	resp, err := c.client.IsWhitelistedTopicWorker(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.IsWhitelistedTopicWorker")
}

func (c *EmissionsGRPCClient) IsWhitelistedTopicReputer(ctx context.Context, req *emissionstypes.IsWhitelistedTopicReputerRequest) (*emissionstypes.IsWhitelistedTopicReputerResponse, error) {
	resp, err := c.client.IsWhitelistedTopicReputer(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.IsWhitelistedTopicReputer")
}

func (c *EmissionsGRPCClient) CanUpdateAllGlobalWhitelists(ctx context.Context, req *emissionstypes.CanUpdateAllGlobalWhitelistsRequest) (*emissionstypes.CanUpdateAllGlobalWhitelistsResponse, error) {
	resp, err := c.client.CanUpdateAllGlobalWhitelists(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.CanUpdateAllGlobalWhitelists")
}

func (c *EmissionsGRPCClient) CanUpdateGlobalWorkerWhitelist(ctx context.Context, req *emissionstypes.CanUpdateGlobalWorkerWhitelistRequest) (*emissionstypes.CanUpdateGlobalWorkerWhitelistResponse, error) {
	resp, err := c.client.CanUpdateGlobalWorkerWhitelist(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.CanUpdateGlobalWorkerWhitelist")
}

func (c *EmissionsGRPCClient) CanUpdateGlobalReputerWhitelist(ctx context.Context, req *emissionstypes.CanUpdateGlobalReputerWhitelistRequest) (*emissionstypes.CanUpdateGlobalReputerWhitelistResponse, error) {
	resp, err := c.client.CanUpdateGlobalReputerWhitelist(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.CanUpdateGlobalReputerWhitelist")
}

func (c *EmissionsGRPCClient) CanUpdateParams(ctx context.Context, req *emissionstypes.CanUpdateParamsRequest) (*emissionstypes.CanUpdateParamsResponse, error) {
	resp, err := c.client.CanUpdateParams(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.CanUpdateParams")
}

func (c *EmissionsGRPCClient) CanUpdateTopicWhitelist(ctx context.Context, req *emissionstypes.CanUpdateTopicWhitelistRequest) (*emissionstypes.CanUpdateTopicWhitelistResponse, error) {
	resp, err := c.client.CanUpdateTopicWhitelist(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.CanUpdateTopicWhitelist")
}

func (c *EmissionsGRPCClient) CanCreateTopic(ctx context.Context, req *emissionstypes.CanCreateTopicRequest) (*emissionstypes.CanCreateTopicResponse, error) {
	resp, err := c.client.CanCreateTopic(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.CanCreateTopic")
}

func (c *EmissionsGRPCClient) CanSubmitWorkerPayload(ctx context.Context, req *emissionstypes.CanSubmitWorkerPayloadRequest) (*emissionstypes.CanSubmitWorkerPayloadResponse, error) {
	resp, err := c.client.CanSubmitWorkerPayload(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.CanSubmitWorkerPayload")
}

func (c *EmissionsGRPCClient) CanSubmitReputerPayload(ctx context.Context, req *emissionstypes.CanSubmitReputerPayloadRequest) (*emissionstypes.CanSubmitReputerPayloadResponse, error) {
	resp, err := c.client.CanSubmitReputerPayload(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.CanSubmitReputerPayload")
}

func (c *EmissionsGRPCClient) GetTopicInitialInfererEmaScore(ctx context.Context, req *emissionstypes.GetTopicInitialInfererEmaScoreRequest) (*emissionstypes.GetTopicInitialInfererEmaScoreResponse, error) {
	resp, err := c.client.GetTopicInitialInfererEmaScore(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetTopicInitialInfererEmaScore")
}

func (c *EmissionsGRPCClient) GetTopicInitialForecasterEmaScore(ctx context.Context, req *emissionstypes.GetTopicInitialForecasterEmaScoreRequest) (*emissionstypes.GetTopicInitialForecasterEmaScoreResponse, error) {
	resp, err := c.client.GetTopicInitialForecasterEmaScore(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetTopicInitialForecasterEmaScore")
}

func (c *EmissionsGRPCClient) GetTopicInitialReputerEmaScore(ctx context.Context, req *emissionstypes.GetTopicInitialReputerEmaScoreRequest) (*emissionstypes.GetTopicInitialReputerEmaScoreResponse, error) {
	resp, err := c.client.GetTopicInitialReputerEmaScore(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetTopicInitialReputerEmaScore")
}

func (c *EmissionsGRPCClient) GetLatestRegretStdNorm(ctx context.Context, req *emissionstypes.GetLatestRegretStdNormRequest) (*emissionstypes.GetLatestRegretStdNormResponse, error) {
	resp, err := c.client.GetLatestRegretStdNorm(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetLatestRegretStdNorm")
}

func (c *EmissionsGRPCClient) GetLatestInfererWeight(ctx context.Context, req *emissionstypes.GetLatestInfererWeightRequest) (*emissionstypes.GetLatestInfererWeightResponse, error) {
	resp, err := c.client.GetLatestInfererWeight(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetLatestInfererWeight")
}

func (c *EmissionsGRPCClient) GetLatestForecasterWeight(ctx context.Context, req *emissionstypes.GetLatestForecasterWeightRequest) (*emissionstypes.GetLatestForecasterWeightResponse, error) {
	resp, err := c.client.GetLatestForecasterWeight(ctx, req)
	return resp, errors.Wrap(err, "while calling AuthGRPCClient.GetLatestForecasterWeight")
}
