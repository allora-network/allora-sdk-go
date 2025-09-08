package wrapper

import (
	"context"

	emissionstypes "github.com/allora-network/allora-chain/x/emissions/types"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/pool"
)

// EmissionsClientWrapper wraps the emissions module with pool management and retry logic
type EmissionsClientWrapper struct {
	poolManager *pool.ClientPoolManager
	logger      zerolog.Logger
}

// NewEmissionsClientWrapper creates a new emissions client wrapper
func NewEmissionsClientWrapper(poolManager *pool.ClientPoolManager, logger zerolog.Logger) *EmissionsClientWrapper {
	return &EmissionsClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "emissions").Logger(),
	}
}

func (c *EmissionsClientWrapper) GetParams(ctx context.Context, req *emissionstypes.GetParamsRequest) (*emissionstypes.GetParamsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetParamsResponse, error) {
		return client.Emissions().GetParams(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetNextTopicId(ctx context.Context, req *emissionstypes.GetNextTopicIdRequest) (*emissionstypes.GetNextTopicIdResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetNextTopicIdResponse, error) {
		return client.Emissions().GetNextTopicId(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetTopic(ctx context.Context, req *emissionstypes.GetTopicRequest) (*emissionstypes.GetTopicResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetTopicResponse, error) {
		return client.Emissions().GetTopic(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetWorkerLatestInferenceByTopicId(ctx context.Context, req *emissionstypes.GetWorkerLatestInferenceByTopicIdRequest) (*emissionstypes.GetWorkerLatestInferenceByTopicIdResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetWorkerLatestInferenceByTopicIdResponse, error) {
		return client.Emissions().GetWorkerLatestInferenceByTopicId(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetInferencesAtBlock(ctx context.Context, req *emissionstypes.GetInferencesAtBlockRequest) (*emissionstypes.GetInferencesAtBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetInferencesAtBlockResponse, error) {
		return client.Emissions().GetInferencesAtBlock(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetLatestTopicInferences(ctx context.Context, req *emissionstypes.GetLatestTopicInferencesRequest) (*emissionstypes.GetLatestTopicInferencesResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetLatestTopicInferencesResponse, error) {
		return client.Emissions().GetLatestTopicInferences(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetForecastsAtBlock(ctx context.Context, req *emissionstypes.GetForecastsAtBlockRequest) (*emissionstypes.GetForecastsAtBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetForecastsAtBlockResponse, error) {
		return client.Emissions().GetForecastsAtBlock(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetNetworkLossBundleAtBlock(ctx context.Context, req *emissionstypes.GetNetworkLossBundleAtBlockRequest) (*emissionstypes.GetNetworkLossBundleAtBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetNetworkLossBundleAtBlockResponse, error) {
		return client.Emissions().GetNetworkLossBundleAtBlock(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetTotalStake(ctx context.Context, req *emissionstypes.GetTotalStakeRequest) (*emissionstypes.GetTotalStakeResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetTotalStakeResponse, error) {
		return client.Emissions().GetTotalStake(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetReputerStakeInTopic(ctx context.Context, req *emissionstypes.GetReputerStakeInTopicRequest) (*emissionstypes.GetReputerStakeInTopicResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetReputerStakeInTopicResponse, error) {
		return client.Emissions().GetReputerStakeInTopic(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetMultiReputerStakeInTopic(ctx context.Context, req *emissionstypes.GetMultiReputerStakeInTopicRequest) (*emissionstypes.GetMultiReputerStakeInTopicResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetMultiReputerStakeInTopicResponse, error) {
		return client.Emissions().GetMultiReputerStakeInTopic(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetStakeFromReputerInTopicInSelf(ctx context.Context, req *emissionstypes.GetStakeFromReputerInTopicInSelfRequest) (*emissionstypes.GetStakeFromReputerInTopicInSelfResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetStakeFromReputerInTopicInSelfResponse, error) {
		return client.Emissions().GetStakeFromReputerInTopicInSelf(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetDelegateStakeInTopicInReputer(ctx context.Context, req *emissionstypes.GetDelegateStakeInTopicInReputerRequest) (*emissionstypes.GetDelegateStakeInTopicInReputerResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetDelegateStakeInTopicInReputerResponse, error) {
		return client.Emissions().GetDelegateStakeInTopicInReputer(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetStakeFromDelegatorInTopicInReputer(ctx context.Context, req *emissionstypes.GetStakeFromDelegatorInTopicInReputerRequest) (*emissionstypes.GetStakeFromDelegatorInTopicInReputerResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetStakeFromDelegatorInTopicInReputerResponse, error) {
		return client.Emissions().GetStakeFromDelegatorInTopicInReputer(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetStakeFromDelegatorInTopic(ctx context.Context, req *emissionstypes.GetStakeFromDelegatorInTopicRequest) (*emissionstypes.GetStakeFromDelegatorInTopicResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetStakeFromDelegatorInTopicResponse, error) {
		return client.Emissions().GetStakeFromDelegatorInTopic(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetTopicStake(ctx context.Context, req *emissionstypes.GetTopicStakeRequest) (*emissionstypes.GetTopicStakeResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetTopicStakeResponse, error) {
		return client.Emissions().GetTopicStake(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetStakeRemovalsUpUntilBlock(ctx context.Context, req *emissionstypes.GetStakeRemovalsUpUntilBlockRequest) (*emissionstypes.GetStakeRemovalsUpUntilBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetStakeRemovalsUpUntilBlockResponse, error) {
		return client.Emissions().GetStakeRemovalsUpUntilBlock(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetDelegateStakeRemovalsUpUntilBlock(ctx context.Context, req *emissionstypes.GetDelegateStakeRemovalsUpUntilBlockRequest) (*emissionstypes.GetDelegateStakeRemovalsUpUntilBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetDelegateStakeRemovalsUpUntilBlockResponse, error) {
		return client.Emissions().GetDelegateStakeRemovalsUpUntilBlock(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetStakeRemovalInfo(ctx context.Context, req *emissionstypes.GetStakeRemovalInfoRequest) (*emissionstypes.GetStakeRemovalInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetStakeRemovalInfoResponse, error) {
		return client.Emissions().GetStakeRemovalInfo(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetDelegateStakeRemovalInfo(ctx context.Context, req *emissionstypes.GetDelegateStakeRemovalInfoRequest) (*emissionstypes.GetDelegateStakeRemovalInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetDelegateStakeRemovalInfoResponse, error) {
		return client.Emissions().GetDelegateStakeRemovalInfo(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetWorkerNodeInfo(ctx context.Context, req *emissionstypes.GetWorkerNodeInfoRequest) (*emissionstypes.GetWorkerNodeInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetWorkerNodeInfoResponse, error) {
		return client.Emissions().GetWorkerNodeInfo(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetReputerNodeInfo(ctx context.Context, req *emissionstypes.GetReputerNodeInfoRequest) (*emissionstypes.GetReputerNodeInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetReputerNodeInfoResponse, error) {
		return client.Emissions().GetReputerNodeInfo(ctx, req)
	})
}

func (c *EmissionsClientWrapper) IsWorkerRegisteredInTopicId(ctx context.Context, req *emissionstypes.IsWorkerRegisteredInTopicIdRequest) (*emissionstypes.IsWorkerRegisteredInTopicIdResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.IsWorkerRegisteredInTopicIdResponse, error) {
		return client.Emissions().IsWorkerRegisteredInTopicId(ctx, req)
	})
}

func (c *EmissionsClientWrapper) IsReputerRegisteredInTopicId(ctx context.Context, req *emissionstypes.IsReputerRegisteredInTopicIdRequest) (*emissionstypes.IsReputerRegisteredInTopicIdResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.IsReputerRegisteredInTopicIdResponse, error) {
		return client.Emissions().IsReputerRegisteredInTopicId(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetNetworkInferencesAtBlock(ctx context.Context, req *emissionstypes.GetNetworkInferencesAtBlockRequest) (*emissionstypes.GetNetworkInferencesAtBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetNetworkInferencesAtBlockResponse, error) {
		return client.Emissions().GetNetworkInferencesAtBlock(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetNetworkInferencesAtBlockOutlierResistant(ctx context.Context, req *emissionstypes.GetNetworkInferencesAtBlockOutlierResistantRequest) (*emissionstypes.GetNetworkInferencesAtBlockOutlierResistantResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetNetworkInferencesAtBlockOutlierResistantResponse, error) {
		return client.Emissions().GetNetworkInferencesAtBlockOutlierResistant(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetLatestNetworkInferences(ctx context.Context, req *emissionstypes.GetLatestNetworkInferencesRequest) (*emissionstypes.GetLatestNetworkInferencesResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetLatestNetworkInferencesResponse, error) {
		return client.Emissions().GetLatestNetworkInferences(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetLatestNetworkInferencesOutlierResistant(ctx context.Context, req *emissionstypes.GetLatestNetworkInferencesOutlierResistantRequest) (*emissionstypes.GetLatestNetworkInferencesOutlierResistantResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetLatestNetworkInferencesOutlierResistantResponse, error) {
		return client.Emissions().GetLatestNetworkInferencesOutlierResistant(ctx, req)
	})
}

func (c *EmissionsClientWrapper) IsWorkerNonceUnfulfilled(ctx context.Context, req *emissionstypes.IsWorkerNonceUnfulfilledRequest) (*emissionstypes.IsWorkerNonceUnfulfilledResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.IsWorkerNonceUnfulfilledResponse, error) {
		return client.Emissions().IsWorkerNonceUnfulfilled(ctx, req)
	})
}

func (c *EmissionsClientWrapper) IsReputerNonceUnfulfilled(ctx context.Context, req *emissionstypes.IsReputerNonceUnfulfilledRequest) (*emissionstypes.IsReputerNonceUnfulfilledResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.IsReputerNonceUnfulfilledResponse, error) {
		return client.Emissions().IsReputerNonceUnfulfilled(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetUnfulfilledWorkerNonces(ctx context.Context, req *emissionstypes.GetUnfulfilledWorkerNoncesRequest) (*emissionstypes.GetUnfulfilledWorkerNoncesResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetUnfulfilledWorkerNoncesResponse, error) {
		return client.Emissions().GetUnfulfilledWorkerNonces(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetUnfulfilledReputerNonces(ctx context.Context, req *emissionstypes.GetUnfulfilledReputerNoncesRequest) (*emissionstypes.GetUnfulfilledReputerNoncesResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetUnfulfilledReputerNoncesResponse, error) {
		return client.Emissions().GetUnfulfilledReputerNonces(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetInfererNetworkRegretRequest) (*emissionstypes.GetInfererNetworkRegretResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetInfererNetworkRegretResponse, error) {
		return client.Emissions().GetInfererNetworkRegret(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetForecasterNetworkRegretRequest) (*emissionstypes.GetForecasterNetworkRegretResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetForecasterNetworkRegretResponse, error) {
		return client.Emissions().GetForecasterNetworkRegret(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetOneInForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetOneInForecasterNetworkRegretRequest) (*emissionstypes.GetOneInForecasterNetworkRegretResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetOneInForecasterNetworkRegretResponse, error) {
		return client.Emissions().GetOneInForecasterNetworkRegret(ctx, req)
	})
}

func (c *EmissionsClientWrapper) IsWhitelistAdmin(ctx context.Context, req *emissionstypes.IsWhitelistAdminRequest) (*emissionstypes.IsWhitelistAdminResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.IsWhitelistAdminResponse, error) {
		return client.Emissions().IsWhitelistAdmin(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetTopicLastWorkerCommitInfo(ctx context.Context, req *emissionstypes.GetTopicLastWorkerCommitInfoRequest) (*emissionstypes.GetTopicLastWorkerCommitInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetTopicLastWorkerCommitInfoResponse, error) {
		return client.Emissions().GetTopicLastWorkerCommitInfo(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetTopicLastReputerCommitInfo(ctx context.Context, req *emissionstypes.GetTopicLastReputerCommitInfoRequest) (*emissionstypes.GetTopicLastReputerCommitInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetTopicLastReputerCommitInfoResponse, error) {
		return client.Emissions().GetTopicLastReputerCommitInfo(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetTopicRewardNonce(ctx context.Context, req *emissionstypes.GetTopicRewardNonceRequest) (*emissionstypes.GetTopicRewardNonceResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetTopicRewardNonceResponse, error) {
		return client.Emissions().GetTopicRewardNonce(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetReputerLossBundlesAtBlock(ctx context.Context, req *emissionstypes.GetReputerLossBundlesAtBlockRequest) (*emissionstypes.GetReputerLossBundlesAtBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetReputerLossBundlesAtBlockResponse, error) {
		return client.Emissions().GetReputerLossBundlesAtBlock(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetStakeReputerAuthority(ctx context.Context, req *emissionstypes.GetStakeReputerAuthorityRequest) (*emissionstypes.GetStakeReputerAuthorityResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetStakeReputerAuthorityResponse, error) {
		return client.Emissions().GetStakeReputerAuthority(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetDelegateStakePlacement(ctx context.Context, req *emissionstypes.GetDelegateStakePlacementRequest) (*emissionstypes.GetDelegateStakePlacementResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetDelegateStakePlacementResponse, error) {
		return client.Emissions().GetDelegateStakePlacement(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetDelegateStakeUponReputer(ctx context.Context, req *emissionstypes.GetDelegateStakeUponReputerRequest) (*emissionstypes.GetDelegateStakeUponReputerResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetDelegateStakeUponReputerResponse, error) {
		return client.Emissions().GetDelegateStakeUponReputer(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetDelegateRewardPerShare(ctx context.Context, req *emissionstypes.GetDelegateRewardPerShareRequest) (*emissionstypes.GetDelegateRewardPerShareResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetDelegateRewardPerShareResponse, error) {
		return client.Emissions().GetDelegateRewardPerShare(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetStakeRemovalForReputerAndTopicId(ctx context.Context, req *emissionstypes.GetStakeRemovalForReputerAndTopicIdRequest) (*emissionstypes.GetStakeRemovalForReputerAndTopicIdResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetStakeRemovalForReputerAndTopicIdResponse, error) {
		return client.Emissions().GetStakeRemovalForReputerAndTopicId(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetDelegateStakeRemoval(ctx context.Context, req *emissionstypes.GetDelegateStakeRemovalRequest) (*emissionstypes.GetDelegateStakeRemovalResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetDelegateStakeRemovalResponse, error) {
		return client.Emissions().GetDelegateStakeRemoval(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetPreviousTopicWeight(ctx context.Context, req *emissionstypes.GetPreviousTopicWeightRequest) (*emissionstypes.GetPreviousTopicWeightResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetPreviousTopicWeightResponse, error) {
		return client.Emissions().GetPreviousTopicWeight(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetTotalSumPreviousTopicWeights(ctx context.Context, req *emissionstypes.GetTotalSumPreviousTopicWeightsRequest) (*emissionstypes.GetTotalSumPreviousTopicWeightsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetTotalSumPreviousTopicWeightsResponse, error) {
		return client.Emissions().GetTotalSumPreviousTopicWeights(ctx, req)
	})
}

func (c *EmissionsClientWrapper) TopicExists(ctx context.Context, req *emissionstypes.TopicExistsRequest) (*emissionstypes.TopicExistsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.TopicExistsResponse, error) {
		return client.Emissions().TopicExists(ctx, req)
	})
}

func (c *EmissionsClientWrapper) IsTopicActive(ctx context.Context, req *emissionstypes.IsTopicActiveRequest) (*emissionstypes.IsTopicActiveResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.IsTopicActiveResponse, error) {
		return client.Emissions().IsTopicActive(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetTopicFeeRevenue(ctx context.Context, req *emissionstypes.GetTopicFeeRevenueRequest) (*emissionstypes.GetTopicFeeRevenueResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetTopicFeeRevenueResponse, error) {
		return client.Emissions().GetTopicFeeRevenue(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetInfererScoreEma(ctx context.Context, req *emissionstypes.GetInfererScoreEmaRequest) (*emissionstypes.GetInfererScoreEmaResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetInfererScoreEmaResponse, error) {
		return client.Emissions().GetInfererScoreEma(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetForecasterScoreEma(ctx context.Context, req *emissionstypes.GetForecasterScoreEmaRequest) (*emissionstypes.GetForecasterScoreEmaResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetForecasterScoreEmaResponse, error) {
		return client.Emissions().GetForecasterScoreEma(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetReputerScoreEma(ctx context.Context, req *emissionstypes.GetReputerScoreEmaRequest) (*emissionstypes.GetReputerScoreEmaResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetReputerScoreEmaResponse, error) {
		return client.Emissions().GetReputerScoreEma(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetInferenceScoresUntilBlock(ctx context.Context, req *emissionstypes.GetInferenceScoresUntilBlockRequest) (*emissionstypes.GetInferenceScoresUntilBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetInferenceScoresUntilBlockResponse, error) {
		return client.Emissions().GetInferenceScoresUntilBlock(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetPreviousTopicQuantileForecasterScoreEma(ctx context.Context, req *emissionstypes.GetPreviousTopicQuantileForecasterScoreEmaRequest) (*emissionstypes.GetPreviousTopicQuantileForecasterScoreEmaResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetPreviousTopicQuantileForecasterScoreEmaResponse, error) {
		return client.Emissions().GetPreviousTopicQuantileForecasterScoreEma(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetPreviousTopicQuantileInfererScoreEma(ctx context.Context, req *emissionstypes.GetPreviousTopicQuantileInfererScoreEmaRequest) (*emissionstypes.GetPreviousTopicQuantileInfererScoreEmaResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetPreviousTopicQuantileInfererScoreEmaResponse, error) {
		return client.Emissions().GetPreviousTopicQuantileInfererScoreEma(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetPreviousTopicQuantileReputerScoreEma(ctx context.Context, req *emissionstypes.GetPreviousTopicQuantileReputerScoreEmaRequest) (*emissionstypes.GetPreviousTopicQuantileReputerScoreEmaResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetPreviousTopicQuantileReputerScoreEmaResponse, error) {
		return client.Emissions().GetPreviousTopicQuantileReputerScoreEma(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetWorkerInferenceScoresAtBlock(ctx context.Context, req *emissionstypes.GetWorkerInferenceScoresAtBlockRequest) (*emissionstypes.GetWorkerInferenceScoresAtBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetWorkerInferenceScoresAtBlockResponse, error) {
		return client.Emissions().GetWorkerInferenceScoresAtBlock(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetCurrentLowestInfererScore(ctx context.Context, req *emissionstypes.GetCurrentLowestInfererScoreRequest) (*emissionstypes.GetCurrentLowestInfererScoreResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetCurrentLowestInfererScoreResponse, error) {
		return client.Emissions().GetCurrentLowestInfererScore(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetForecastScoresUntilBlock(ctx context.Context, req *emissionstypes.GetForecastScoresUntilBlockRequest) (*emissionstypes.GetForecastScoresUntilBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetForecastScoresUntilBlockResponse, error) {
		return client.Emissions().GetForecastScoresUntilBlock(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetWorkerForecastScoresAtBlock(ctx context.Context, req *emissionstypes.GetWorkerForecastScoresAtBlockRequest) (*emissionstypes.GetWorkerForecastScoresAtBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetWorkerForecastScoresAtBlockResponse, error) {
		return client.Emissions().GetWorkerForecastScoresAtBlock(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetCurrentLowestForecasterScore(ctx context.Context, req *emissionstypes.GetCurrentLowestForecasterScoreRequest) (*emissionstypes.GetCurrentLowestForecasterScoreResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetCurrentLowestForecasterScoreResponse, error) {
		return client.Emissions().GetCurrentLowestForecasterScore(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetReputersScoresAtBlock(ctx context.Context, req *emissionstypes.GetReputersScoresAtBlockRequest) (*emissionstypes.GetReputersScoresAtBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetReputersScoresAtBlockResponse, error) {
		return client.Emissions().GetReputersScoresAtBlock(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetCurrentLowestReputerScore(ctx context.Context, req *emissionstypes.GetCurrentLowestReputerScoreRequest) (*emissionstypes.GetCurrentLowestReputerScoreResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetCurrentLowestReputerScoreResponse, error) {
		return client.Emissions().GetCurrentLowestReputerScore(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetListeningCoefficient(ctx context.Context, req *emissionstypes.GetListeningCoefficientRequest) (*emissionstypes.GetListeningCoefficientResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetListeningCoefficientResponse, error) {
		return client.Emissions().GetListeningCoefficient(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetPreviousReputerRewardFraction(ctx context.Context, req *emissionstypes.GetPreviousReputerRewardFractionRequest) (*emissionstypes.GetPreviousReputerRewardFractionResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetPreviousReputerRewardFractionResponse, error) {
		return client.Emissions().GetPreviousReputerRewardFraction(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetPreviousInferenceRewardFraction(ctx context.Context, req *emissionstypes.GetPreviousInferenceRewardFractionRequest) (*emissionstypes.GetPreviousInferenceRewardFractionResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetPreviousInferenceRewardFractionResponse, error) {
		return client.Emissions().GetPreviousInferenceRewardFraction(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetPreviousForecastRewardFraction(ctx context.Context, req *emissionstypes.GetPreviousForecastRewardFractionRequest) (*emissionstypes.GetPreviousForecastRewardFractionResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetPreviousForecastRewardFractionResponse, error) {
		return client.Emissions().GetPreviousForecastRewardFraction(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetPreviousPercentageRewardToStakedReputers(ctx context.Context, req *emissionstypes.GetPreviousPercentageRewardToStakedReputersRequest) (*emissionstypes.GetPreviousPercentageRewardToStakedReputersResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetPreviousPercentageRewardToStakedReputersResponse, error) {
		return client.Emissions().GetPreviousPercentageRewardToStakedReputers(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetTotalRewardToDistribute(ctx context.Context, req *emissionstypes.GetTotalRewardToDistributeRequest) (*emissionstypes.GetTotalRewardToDistributeResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetTotalRewardToDistributeResponse, error) {
		return client.Emissions().GetTotalRewardToDistribute(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetNaiveInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetNaiveInfererNetworkRegretRequest) (*emissionstypes.GetNaiveInfererNetworkRegretResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetNaiveInfererNetworkRegretResponse, error) {
		return client.Emissions().GetNaiveInfererNetworkRegret(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetOneOutInfererInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutInfererInfererNetworkRegretRequest) (*emissionstypes.GetOneOutInfererInfererNetworkRegretResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetOneOutInfererInfererNetworkRegretResponse, error) {
		return client.Emissions().GetOneOutInfererInfererNetworkRegret(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetOneOutInfererForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutInfererForecasterNetworkRegretRequest) (*emissionstypes.GetOneOutInfererForecasterNetworkRegretResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetOneOutInfererForecasterNetworkRegretResponse, error) {
		return client.Emissions().GetOneOutInfererForecasterNetworkRegret(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetOneOutForecasterInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutForecasterInfererNetworkRegretRequest) (*emissionstypes.GetOneOutForecasterInfererNetworkRegretResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetOneOutForecasterInfererNetworkRegretResponse, error) {
		return client.Emissions().GetOneOutForecasterInfererNetworkRegret(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetOneOutForecasterForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutForecasterForecasterNetworkRegretRequest) (*emissionstypes.GetOneOutForecasterForecasterNetworkRegretResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetOneOutForecasterForecasterNetworkRegretResponse, error) {
		return client.Emissions().GetOneOutForecasterForecasterNetworkRegret(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetActiveTopicsAtBlock(ctx context.Context, req *emissionstypes.GetActiveTopicsAtBlockRequest) (*emissionstypes.GetActiveTopicsAtBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetActiveTopicsAtBlockResponse, error) {
		return client.Emissions().GetActiveTopicsAtBlock(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetNextChurningBlockByTopicId(ctx context.Context, req *emissionstypes.GetNextChurningBlockByTopicIdRequest) (*emissionstypes.GetNextChurningBlockByTopicIdResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetNextChurningBlockByTopicIdResponse, error) {
		return client.Emissions().GetNextChurningBlockByTopicId(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetCountInfererInclusionsInTopic(ctx context.Context, req *emissionstypes.GetCountInfererInclusionsInTopicRequest) (*emissionstypes.GetCountInfererInclusionsInTopicResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetCountInfererInclusionsInTopicResponse, error) {
		return client.Emissions().GetCountInfererInclusionsInTopic(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetCountForecasterInclusionsInTopic(ctx context.Context, req *emissionstypes.GetCountForecasterInclusionsInTopicRequest) (*emissionstypes.GetCountForecasterInclusionsInTopicResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetCountForecasterInclusionsInTopicResponse, error) {
		return client.Emissions().GetCountForecasterInclusionsInTopic(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetActiveReputersForTopic(ctx context.Context, req *emissionstypes.GetActiveReputersForTopicRequest) (*emissionstypes.GetActiveReputersForTopicResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetActiveReputersForTopicResponse, error) {
		return client.Emissions().GetActiveReputersForTopic(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetActiveForecastersForTopic(ctx context.Context, req *emissionstypes.GetActiveForecastersForTopicRequest) (*emissionstypes.GetActiveForecastersForTopicResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetActiveForecastersForTopicResponse, error) {
		return client.Emissions().GetActiveForecastersForTopic(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetActiveInferersForTopic(ctx context.Context, req *emissionstypes.GetActiveInferersForTopicRequest) (*emissionstypes.GetActiveInferersForTopicResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetActiveInferersForTopicResponse, error) {
		return client.Emissions().GetActiveInferersForTopic(ctx, req)
	})
}

func (c *EmissionsClientWrapper) IsWhitelistedGlobalWorker(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalWorkerRequest) (*emissionstypes.IsWhitelistedGlobalWorkerResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.IsWhitelistedGlobalWorkerResponse, error) {
		return client.Emissions().IsWhitelistedGlobalWorker(ctx, req)
	})
}

func (c *EmissionsClientWrapper) IsWhitelistedGlobalReputer(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalReputerRequest) (*emissionstypes.IsWhitelistedGlobalReputerResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.IsWhitelistedGlobalReputerResponse, error) {
		return client.Emissions().IsWhitelistedGlobalReputer(ctx, req)
	})
}

func (c *EmissionsClientWrapper) IsWhitelistedGlobalAdmin(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalAdminRequest) (*emissionstypes.IsWhitelistedGlobalAdminResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.IsWhitelistedGlobalAdminResponse, error) {
		return client.Emissions().IsWhitelistedGlobalAdmin(ctx, req)
	})
}

func (c *EmissionsClientWrapper) IsTopicWorkerWhitelistEnabled(ctx context.Context, req *emissionstypes.IsTopicWorkerWhitelistEnabledRequest) (*emissionstypes.IsTopicWorkerWhitelistEnabledResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.IsTopicWorkerWhitelistEnabledResponse, error) {
		return client.Emissions().IsTopicWorkerWhitelistEnabled(ctx, req)
	})
}

func (c *EmissionsClientWrapper) IsTopicReputerWhitelistEnabled(ctx context.Context, req *emissionstypes.IsTopicReputerWhitelistEnabledRequest) (*emissionstypes.IsTopicReputerWhitelistEnabledResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.IsTopicReputerWhitelistEnabledResponse, error) {
		return client.Emissions().IsTopicReputerWhitelistEnabled(ctx, req)
	})
}

func (c *EmissionsClientWrapper) IsWhitelistedTopicCreator(ctx context.Context, req *emissionstypes.IsWhitelistedTopicCreatorRequest) (*emissionstypes.IsWhitelistedTopicCreatorResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.IsWhitelistedTopicCreatorResponse, error) {
		return client.Emissions().IsWhitelistedTopicCreator(ctx, req)
	})
}

func (c *EmissionsClientWrapper) IsWhitelistedGlobalActor(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalActorRequest) (*emissionstypes.IsWhitelistedGlobalActorResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.IsWhitelistedGlobalActorResponse, error) {
		return client.Emissions().IsWhitelistedGlobalActor(ctx, req)
	})
}

func (c *EmissionsClientWrapper) IsWhitelistedTopicWorker(ctx context.Context, req *emissionstypes.IsWhitelistedTopicWorkerRequest) (*emissionstypes.IsWhitelistedTopicWorkerResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.IsWhitelistedTopicWorkerResponse, error) {
		return client.Emissions().IsWhitelistedTopicWorker(ctx, req)
	})
}

func (c *EmissionsClientWrapper) IsWhitelistedTopicReputer(ctx context.Context, req *emissionstypes.IsWhitelistedTopicReputerRequest) (*emissionstypes.IsWhitelistedTopicReputerResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.IsWhitelistedTopicReputerResponse, error) {
		return client.Emissions().IsWhitelistedTopicReputer(ctx, req)
	})
}

func (c *EmissionsClientWrapper) CanUpdateAllGlobalWhitelists(ctx context.Context, req *emissionstypes.CanUpdateAllGlobalWhitelistsRequest) (*emissionstypes.CanUpdateAllGlobalWhitelistsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.CanUpdateAllGlobalWhitelistsResponse, error) {
		return client.Emissions().CanUpdateAllGlobalWhitelists(ctx, req)
	})
}

func (c *EmissionsClientWrapper) CanUpdateGlobalWorkerWhitelist(ctx context.Context, req *emissionstypes.CanUpdateGlobalWorkerWhitelistRequest) (*emissionstypes.CanUpdateGlobalWorkerWhitelistResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.CanUpdateGlobalWorkerWhitelistResponse, error) {
		return client.Emissions().CanUpdateGlobalWorkerWhitelist(ctx, req)
	})
}

func (c *EmissionsClientWrapper) CanUpdateGlobalReputerWhitelist(ctx context.Context, req *emissionstypes.CanUpdateGlobalReputerWhitelistRequest) (*emissionstypes.CanUpdateGlobalReputerWhitelistResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.CanUpdateGlobalReputerWhitelistResponse, error) {
		return client.Emissions().CanUpdateGlobalReputerWhitelist(ctx, req)
	})
}

func (c *EmissionsClientWrapper) CanUpdateParams(ctx context.Context, req *emissionstypes.CanUpdateParamsRequest) (*emissionstypes.CanUpdateParamsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.CanUpdateParamsResponse, error) {
		return client.Emissions().CanUpdateParams(ctx, req)
	})
}

func (c *EmissionsClientWrapper) CanUpdateTopicWhitelist(ctx context.Context, req *emissionstypes.CanUpdateTopicWhitelistRequest) (*emissionstypes.CanUpdateTopicWhitelistResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.CanUpdateTopicWhitelistResponse, error) {
		return client.Emissions().CanUpdateTopicWhitelist(ctx, req)
	})
}

func (c *EmissionsClientWrapper) CanCreateTopic(ctx context.Context, req *emissionstypes.CanCreateTopicRequest) (*emissionstypes.CanCreateTopicResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.CanCreateTopicResponse, error) {
		return client.Emissions().CanCreateTopic(ctx, req)
	})
}

func (c *EmissionsClientWrapper) CanSubmitWorkerPayload(ctx context.Context, req *emissionstypes.CanSubmitWorkerPayloadRequest) (*emissionstypes.CanSubmitWorkerPayloadResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.CanSubmitWorkerPayloadResponse, error) {
		return client.Emissions().CanSubmitWorkerPayload(ctx, req)
	})
}

func (c *EmissionsClientWrapper) CanSubmitReputerPayload(ctx context.Context, req *emissionstypes.CanSubmitReputerPayloadRequest) (*emissionstypes.CanSubmitReputerPayloadResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.CanSubmitReputerPayloadResponse, error) {
		return client.Emissions().CanSubmitReputerPayload(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetTopicInitialInfererEmaScore(ctx context.Context, req *emissionstypes.GetTopicInitialInfererEmaScoreRequest) (*emissionstypes.GetTopicInitialInfererEmaScoreResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetTopicInitialInfererEmaScoreResponse, error) {
		return client.Emissions().GetTopicInitialInfererEmaScore(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetTopicInitialForecasterEmaScore(ctx context.Context, req *emissionstypes.GetTopicInitialForecasterEmaScoreRequest) (*emissionstypes.GetTopicInitialForecasterEmaScoreResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetTopicInitialForecasterEmaScoreResponse, error) {
		return client.Emissions().GetTopicInitialForecasterEmaScore(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetTopicInitialReputerEmaScore(ctx context.Context, req *emissionstypes.GetTopicInitialReputerEmaScoreRequest) (*emissionstypes.GetTopicInitialReputerEmaScoreResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetTopicInitialReputerEmaScoreResponse, error) {
		return client.Emissions().GetTopicInitialReputerEmaScore(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetLatestRegretStdNorm(ctx context.Context, req *emissionstypes.GetLatestRegretStdNormRequest) (*emissionstypes.GetLatestRegretStdNormResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetLatestRegretStdNormResponse, error) {
		return client.Emissions().GetLatestRegretStdNorm(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetLatestInfererWeight(ctx context.Context, req *emissionstypes.GetLatestInfererWeightRequest) (*emissionstypes.GetLatestInfererWeightResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetLatestInfererWeightResponse, error) {
		return client.Emissions().GetLatestInfererWeight(ctx, req)
	})
}

func (c *EmissionsClientWrapper) GetLatestForecasterWeight(ctx context.Context, req *emissionstypes.GetLatestForecasterWeightRequest) (*emissionstypes.GetLatestForecasterWeightResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.Client) (*emissionstypes.GetLatestForecasterWeightResponse, error) {
		return client.Emissions().GetLatestForecasterWeight(ctx, req)
	})
}
