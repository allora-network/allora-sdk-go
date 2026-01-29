package wrapper

import (
	"context"

	emissionstypes "github.com/allora-network/allora-chain/x/emissions/types"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/pool"
)

// EmissionsClientWrapper wraps the emissions module with pool management and retry logic
type EmissionsClientWrapper struct {
	poolManager *pool.ClientPoolManager[interfaces.CosmosClient]
	logger      zerolog.Logger
}

// NewEmissionsClientWrapper creates a new emissions client wrapper
func NewEmissionsClientWrapper(poolManager *pool.ClientPoolManager[interfaces.CosmosClient], logger zerolog.Logger) *EmissionsClientWrapper {
	return &EmissionsClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "emissions").Logger(),
	}
}

func (c *EmissionsClientWrapper) CanCreateTopic(ctx context.Context, req *emissionstypes.CanCreateTopicRequest, opts ...config.CallOpt) (*emissionstypes.CanCreateTopicResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.CanCreateTopicResponse, error) {
		return client.Emissions().CanCreateTopic(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) CanSubmitReputerPayload(ctx context.Context, req *emissionstypes.CanSubmitReputerPayloadRequest, opts ...config.CallOpt) (*emissionstypes.CanSubmitReputerPayloadResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.CanSubmitReputerPayloadResponse, error) {
		return client.Emissions().CanSubmitReputerPayload(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) CanSubmitWorkerPayload(ctx context.Context, req *emissionstypes.CanSubmitWorkerPayloadRequest, opts ...config.CallOpt) (*emissionstypes.CanSubmitWorkerPayloadResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.CanSubmitWorkerPayloadResponse, error) {
		return client.Emissions().CanSubmitWorkerPayload(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) CanUpdateAllGlobalWhitelists(ctx context.Context, req *emissionstypes.CanUpdateAllGlobalWhitelistsRequest, opts ...config.CallOpt) (*emissionstypes.CanUpdateAllGlobalWhitelistsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.CanUpdateAllGlobalWhitelistsResponse, error) {
		return client.Emissions().CanUpdateAllGlobalWhitelists(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) CanUpdateGlobalReputerWhitelist(ctx context.Context, req *emissionstypes.CanUpdateGlobalReputerWhitelistRequest, opts ...config.CallOpt) (*emissionstypes.CanUpdateGlobalReputerWhitelistResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.CanUpdateGlobalReputerWhitelistResponse, error) {
		return client.Emissions().CanUpdateGlobalReputerWhitelist(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) CanUpdateGlobalWorkerWhitelist(ctx context.Context, req *emissionstypes.CanUpdateGlobalWorkerWhitelistRequest, opts ...config.CallOpt) (*emissionstypes.CanUpdateGlobalWorkerWhitelistResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.CanUpdateGlobalWorkerWhitelistResponse, error) {
		return client.Emissions().CanUpdateGlobalWorkerWhitelist(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) CanUpdateParams(ctx context.Context, req *emissionstypes.CanUpdateParamsRequest, opts ...config.CallOpt) (*emissionstypes.CanUpdateParamsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.CanUpdateParamsResponse, error) {
		return client.Emissions().CanUpdateParams(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) CanUpdateTopicWhitelist(ctx context.Context, req *emissionstypes.CanUpdateTopicWhitelistRequest, opts ...config.CallOpt) (*emissionstypes.CanUpdateTopicWhitelistResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.CanUpdateTopicWhitelistResponse, error) {
		return client.Emissions().CanUpdateTopicWhitelist(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetActiveTopicsAtBlock(ctx context.Context, req *emissionstypes.GetActiveTopicsAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetActiveTopicsAtBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetActiveTopicsAtBlockResponse, error) {
		return client.Emissions().GetActiveTopicsAtBlock(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetCountForecasterInclusionsInTopic(ctx context.Context, req *emissionstypes.GetCountForecasterInclusionsInTopicRequest, opts ...config.CallOpt) (*emissionstypes.GetCountForecasterInclusionsInTopicResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetCountForecasterInclusionsInTopicResponse, error) {
		return client.Emissions().GetCountForecasterInclusionsInTopic(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetCountInfererInclusionsInTopic(ctx context.Context, req *emissionstypes.GetCountInfererInclusionsInTopicRequest, opts ...config.CallOpt) (*emissionstypes.GetCountInfererInclusionsInTopicResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetCountInfererInclusionsInTopicResponse, error) {
		return client.Emissions().GetCountInfererInclusionsInTopic(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetCurrentLowestForecasterScore(ctx context.Context, req *emissionstypes.GetCurrentLowestForecasterScoreRequest, opts ...config.CallOpt) (*emissionstypes.GetCurrentLowestForecasterScoreResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetCurrentLowestForecasterScoreResponse, error) {
		return client.Emissions().GetCurrentLowestForecasterScore(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetCurrentLowestInfererScore(ctx context.Context, req *emissionstypes.GetCurrentLowestInfererScoreRequest, opts ...config.CallOpt) (*emissionstypes.GetCurrentLowestInfererScoreResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetCurrentLowestInfererScoreResponse, error) {
		return client.Emissions().GetCurrentLowestInfererScore(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetCurrentLowestReputerScore(ctx context.Context, req *emissionstypes.GetCurrentLowestReputerScoreRequest, opts ...config.CallOpt) (*emissionstypes.GetCurrentLowestReputerScoreResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetCurrentLowestReputerScoreResponse, error) {
		return client.Emissions().GetCurrentLowestReputerScore(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetDelegateRewardPerShare(ctx context.Context, req *emissionstypes.GetDelegateRewardPerShareRequest, opts ...config.CallOpt) (*emissionstypes.GetDelegateRewardPerShareResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetDelegateRewardPerShareResponse, error) {
		return client.Emissions().GetDelegateRewardPerShare(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetDelegateStakeInTopicInReputer(ctx context.Context, req *emissionstypes.GetDelegateStakeInTopicInReputerRequest, opts ...config.CallOpt) (*emissionstypes.GetDelegateStakeInTopicInReputerResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetDelegateStakeInTopicInReputerResponse, error) {
		return client.Emissions().GetDelegateStakeInTopicInReputer(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetDelegateStakePlacement(ctx context.Context, req *emissionstypes.GetDelegateStakePlacementRequest, opts ...config.CallOpt) (*emissionstypes.GetDelegateStakePlacementResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetDelegateStakePlacementResponse, error) {
		return client.Emissions().GetDelegateStakePlacement(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetDelegateStakeRemoval(ctx context.Context, req *emissionstypes.GetDelegateStakeRemovalRequest, opts ...config.CallOpt) (*emissionstypes.GetDelegateStakeRemovalResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetDelegateStakeRemovalResponse, error) {
		return client.Emissions().GetDelegateStakeRemoval(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetDelegateStakeRemovalInfo(ctx context.Context, req *emissionstypes.GetDelegateStakeRemovalInfoRequest, opts ...config.CallOpt) (*emissionstypes.GetDelegateStakeRemovalInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetDelegateStakeRemovalInfoResponse, error) {
		return client.Emissions().GetDelegateStakeRemovalInfo(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetDelegateStakeRemovalsUpUntilBlock(ctx context.Context, req *emissionstypes.GetDelegateStakeRemovalsUpUntilBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetDelegateStakeRemovalsUpUntilBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetDelegateStakeRemovalsUpUntilBlockResponse, error) {
		return client.Emissions().GetDelegateStakeRemovalsUpUntilBlock(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetDelegateStakeUponReputer(ctx context.Context, req *emissionstypes.GetDelegateStakeUponReputerRequest, opts ...config.CallOpt) (*emissionstypes.GetDelegateStakeUponReputerResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetDelegateStakeUponReputerResponse, error) {
		return client.Emissions().GetDelegateStakeUponReputer(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetForecastScoresUntilBlock(ctx context.Context, req *emissionstypes.GetForecastScoresUntilBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetForecastScoresUntilBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetForecastScoresUntilBlockResponse, error) {
		return client.Emissions().GetForecastScoresUntilBlock(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetForecasterNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetForecasterNetworkRegretResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetForecasterNetworkRegretResponse, error) {
		return client.Emissions().GetForecasterNetworkRegret(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetForecasterScoreEma(ctx context.Context, req *emissionstypes.GetForecasterScoreEmaRequest, opts ...config.CallOpt) (*emissionstypes.GetForecasterScoreEmaResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetForecasterScoreEmaResponse, error) {
		return client.Emissions().GetForecasterScoreEma(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetForecastsAtBlock(ctx context.Context, req *emissionstypes.GetForecastsAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetForecastsAtBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetForecastsAtBlockResponse, error) {
		return client.Emissions().GetForecastsAtBlock(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetInferenceScoresUntilBlock(ctx context.Context, req *emissionstypes.GetInferenceScoresUntilBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetInferenceScoresUntilBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetInferenceScoresUntilBlockResponse, error) {
		return client.Emissions().GetInferenceScoresUntilBlock(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetInferencesAtBlock(ctx context.Context, req *emissionstypes.GetInferencesAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetInferencesAtBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetInferencesAtBlockResponse, error) {
		return client.Emissions().GetInferencesAtBlock(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetInfererNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetInfererNetworkRegretResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetInfererNetworkRegretResponse, error) {
		return client.Emissions().GetInfererNetworkRegret(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetInfererScoreEma(ctx context.Context, req *emissionstypes.GetInfererScoreEmaRequest, opts ...config.CallOpt) (*emissionstypes.GetInfererScoreEmaResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetInfererScoreEmaResponse, error) {
		return client.Emissions().GetInfererScoreEma(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetLatestForecasterWeight(ctx context.Context, req *emissionstypes.GetLatestForecasterWeightRequest, opts ...config.CallOpt) (*emissionstypes.GetLatestForecasterWeightResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetLatestForecasterWeightResponse, error) {
		return client.Emissions().GetLatestForecasterWeight(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetLatestInfererWeight(ctx context.Context, req *emissionstypes.GetLatestInfererWeightRequest, opts ...config.CallOpt) (*emissionstypes.GetLatestInfererWeightResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetLatestInfererWeightResponse, error) {
		return client.Emissions().GetLatestInfererWeight(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetLatestNetworkInferences(ctx context.Context, req *emissionstypes.GetLatestNetworkInferencesRequest, opts ...config.CallOpt) (*emissionstypes.GetLatestNetworkInferencesResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetLatestNetworkInferencesResponse, error) {
		return client.Emissions().GetLatestNetworkInferences(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetLatestNetworkInferencesOutlierResistant(ctx context.Context, req *emissionstypes.GetLatestNetworkInferencesOutlierResistantRequest, opts ...config.CallOpt) (*emissionstypes.GetLatestNetworkInferencesOutlierResistantResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetLatestNetworkInferencesOutlierResistantResponse, error) {
		return client.Emissions().GetLatestNetworkInferencesOutlierResistant(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetLatestRegretStdNorm(ctx context.Context, req *emissionstypes.GetLatestRegretStdNormRequest, opts ...config.CallOpt) (*emissionstypes.GetLatestRegretStdNormResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetLatestRegretStdNormResponse, error) {
		return client.Emissions().GetLatestRegretStdNorm(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetLatestTopicInferences(ctx context.Context, req *emissionstypes.GetLatestTopicInferencesRequest, opts ...config.CallOpt) (*emissionstypes.GetLatestTopicInferencesResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetLatestTopicInferencesResponse, error) {
		return client.Emissions().GetLatestTopicInferences(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetListeningCoefficient(ctx context.Context, req *emissionstypes.GetListeningCoefficientRequest, opts ...config.CallOpt) (*emissionstypes.GetListeningCoefficientResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetListeningCoefficientResponse, error) {
		return client.Emissions().GetListeningCoefficient(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetMultiReputerStakeInTopic(ctx context.Context, req *emissionstypes.GetMultiReputerStakeInTopicRequest, opts ...config.CallOpt) (*emissionstypes.GetMultiReputerStakeInTopicResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetMultiReputerStakeInTopicResponse, error) {
		return client.Emissions().GetMultiReputerStakeInTopic(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetNaiveInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetNaiveInfererNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetNaiveInfererNetworkRegretResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetNaiveInfererNetworkRegretResponse, error) {
		return client.Emissions().GetNaiveInfererNetworkRegret(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetNetworkInferencesAtBlock(ctx context.Context, req *emissionstypes.GetNetworkInferencesAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetNetworkInferencesAtBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetNetworkInferencesAtBlockResponse, error) {
		return client.Emissions().GetNetworkInferencesAtBlock(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetNetworkInferencesAtBlockOutlierResistant(ctx context.Context, req *emissionstypes.GetNetworkInferencesAtBlockOutlierResistantRequest, opts ...config.CallOpt) (*emissionstypes.GetNetworkInferencesAtBlockOutlierResistantResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetNetworkInferencesAtBlockOutlierResistantResponse, error) {
		return client.Emissions().GetNetworkInferencesAtBlockOutlierResistant(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetNetworkLossBundleAtBlock(ctx context.Context, req *emissionstypes.GetNetworkLossBundleAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetNetworkLossBundleAtBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetNetworkLossBundleAtBlockResponse, error) {
		return client.Emissions().GetNetworkLossBundleAtBlock(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetNextChurningBlockByTopicId(ctx context.Context, req *emissionstypes.GetNextChurningBlockByTopicIdRequest, opts ...config.CallOpt) (*emissionstypes.GetNextChurningBlockByTopicIdResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetNextChurningBlockByTopicIdResponse, error) {
		return client.Emissions().GetNextChurningBlockByTopicId(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetNextTopicId(ctx context.Context, req *emissionstypes.GetNextTopicIdRequest, opts ...config.CallOpt) (*emissionstypes.GetNextTopicIdResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetNextTopicIdResponse, error) {
		return client.Emissions().GetNextTopicId(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetOneInForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetOneInForecasterNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetOneInForecasterNetworkRegretResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetOneInForecasterNetworkRegretResponse, error) {
		return client.Emissions().GetOneInForecasterNetworkRegret(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetOneOutForecasterForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutForecasterForecasterNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetOneOutForecasterForecasterNetworkRegretResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetOneOutForecasterForecasterNetworkRegretResponse, error) {
		return client.Emissions().GetOneOutForecasterForecasterNetworkRegret(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetOneOutForecasterInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutForecasterInfererNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetOneOutForecasterInfererNetworkRegretResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetOneOutForecasterInfererNetworkRegretResponse, error) {
		return client.Emissions().GetOneOutForecasterInfererNetworkRegret(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetOneOutInfererForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutInfererForecasterNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetOneOutInfererForecasterNetworkRegretResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetOneOutInfererForecasterNetworkRegretResponse, error) {
		return client.Emissions().GetOneOutInfererForecasterNetworkRegret(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetOneOutInfererInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutInfererInfererNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetOneOutInfererInfererNetworkRegretResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetOneOutInfererInfererNetworkRegretResponse, error) {
		return client.Emissions().GetOneOutInfererInfererNetworkRegret(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetParams(ctx context.Context, req *emissionstypes.GetParamsRequest, opts ...config.CallOpt) (*emissionstypes.GetParamsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetParamsResponse, error) {
		return client.Emissions().GetParams(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetPreviousForecastRewardFraction(ctx context.Context, req *emissionstypes.GetPreviousForecastRewardFractionRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousForecastRewardFractionResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetPreviousForecastRewardFractionResponse, error) {
		return client.Emissions().GetPreviousForecastRewardFraction(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetPreviousInferenceRewardFraction(ctx context.Context, req *emissionstypes.GetPreviousInferenceRewardFractionRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousInferenceRewardFractionResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetPreviousInferenceRewardFractionResponse, error) {
		return client.Emissions().GetPreviousInferenceRewardFraction(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetPreviousPercentageRewardToStakedReputers(ctx context.Context, req *emissionstypes.GetPreviousPercentageRewardToStakedReputersRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousPercentageRewardToStakedReputersResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetPreviousPercentageRewardToStakedReputersResponse, error) {
		return client.Emissions().GetPreviousPercentageRewardToStakedReputers(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetPreviousReputerRewardFraction(ctx context.Context, req *emissionstypes.GetPreviousReputerRewardFractionRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousReputerRewardFractionResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetPreviousReputerRewardFractionResponse, error) {
		return client.Emissions().GetPreviousReputerRewardFraction(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetPreviousTopicQuantileForecasterScoreEma(ctx context.Context, req *emissionstypes.GetPreviousTopicQuantileForecasterScoreEmaRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousTopicQuantileForecasterScoreEmaResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetPreviousTopicQuantileForecasterScoreEmaResponse, error) {
		return client.Emissions().GetPreviousTopicQuantileForecasterScoreEma(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetPreviousTopicQuantileInfererScoreEma(ctx context.Context, req *emissionstypes.GetPreviousTopicQuantileInfererScoreEmaRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousTopicQuantileInfererScoreEmaResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetPreviousTopicQuantileInfererScoreEmaResponse, error) {
		return client.Emissions().GetPreviousTopicQuantileInfererScoreEma(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetPreviousTopicQuantileReputerScoreEma(ctx context.Context, req *emissionstypes.GetPreviousTopicQuantileReputerScoreEmaRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousTopicQuantileReputerScoreEmaResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetPreviousTopicQuantileReputerScoreEmaResponse, error) {
		return client.Emissions().GetPreviousTopicQuantileReputerScoreEma(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetPreviousTopicWeight(ctx context.Context, req *emissionstypes.GetPreviousTopicWeightRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousTopicWeightResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetPreviousTopicWeightResponse, error) {
		return client.Emissions().GetPreviousTopicWeight(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetReputerLossBundlesAtBlock(ctx context.Context, req *emissionstypes.GetReputerLossBundlesAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetReputerLossBundlesAtBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetReputerLossBundlesAtBlockResponse, error) {
		return client.Emissions().GetReputerLossBundlesAtBlock(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetReputerNodeInfo(ctx context.Context, req *emissionstypes.GetReputerNodeInfoRequest, opts ...config.CallOpt) (*emissionstypes.GetReputerNodeInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetReputerNodeInfoResponse, error) {
		return client.Emissions().GetReputerNodeInfo(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetReputerScoreEma(ctx context.Context, req *emissionstypes.GetReputerScoreEmaRequest, opts ...config.CallOpt) (*emissionstypes.GetReputerScoreEmaResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetReputerScoreEmaResponse, error) {
		return client.Emissions().GetReputerScoreEma(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetReputerStakeInTopic(ctx context.Context, req *emissionstypes.GetReputerStakeInTopicRequest, opts ...config.CallOpt) (*emissionstypes.GetReputerStakeInTopicResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetReputerStakeInTopicResponse, error) {
		return client.Emissions().GetReputerStakeInTopic(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetReputerSubmissionWindowStatus(ctx context.Context, req *emissionstypes.GetReputerSubmissionWindowStatusRequest, opts ...config.CallOpt) (*emissionstypes.GetReputerSubmissionWindowStatusResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetReputerSubmissionWindowStatusResponse, error) {
		return client.Emissions().GetReputerSubmissionWindowStatus(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetReputersScoresAtBlock(ctx context.Context, req *emissionstypes.GetReputersScoresAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetReputersScoresAtBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetReputersScoresAtBlockResponse, error) {
		return client.Emissions().GetReputersScoresAtBlock(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetStakeFromDelegatorInTopic(ctx context.Context, req *emissionstypes.GetStakeFromDelegatorInTopicRequest, opts ...config.CallOpt) (*emissionstypes.GetStakeFromDelegatorInTopicResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetStakeFromDelegatorInTopicResponse, error) {
		return client.Emissions().GetStakeFromDelegatorInTopic(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetStakeFromDelegatorInTopicInReputer(ctx context.Context, req *emissionstypes.GetStakeFromDelegatorInTopicInReputerRequest, opts ...config.CallOpt) (*emissionstypes.GetStakeFromDelegatorInTopicInReputerResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetStakeFromDelegatorInTopicInReputerResponse, error) {
		return client.Emissions().GetStakeFromDelegatorInTopicInReputer(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetStakeFromReputerInTopicInSelf(ctx context.Context, req *emissionstypes.GetStakeFromReputerInTopicInSelfRequest, opts ...config.CallOpt) (*emissionstypes.GetStakeFromReputerInTopicInSelfResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetStakeFromReputerInTopicInSelfResponse, error) {
		return client.Emissions().GetStakeFromReputerInTopicInSelf(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetStakeRemovalForReputerAndTopicId(ctx context.Context, req *emissionstypes.GetStakeRemovalForReputerAndTopicIdRequest, opts ...config.CallOpt) (*emissionstypes.GetStakeRemovalForReputerAndTopicIdResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetStakeRemovalForReputerAndTopicIdResponse, error) {
		return client.Emissions().GetStakeRemovalForReputerAndTopicId(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetStakeRemovalInfo(ctx context.Context, req *emissionstypes.GetStakeRemovalInfoRequest, opts ...config.CallOpt) (*emissionstypes.GetStakeRemovalInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetStakeRemovalInfoResponse, error) {
		return client.Emissions().GetStakeRemovalInfo(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetStakeRemovalsUpUntilBlock(ctx context.Context, req *emissionstypes.GetStakeRemovalsUpUntilBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetStakeRemovalsUpUntilBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetStakeRemovalsUpUntilBlockResponse, error) {
		return client.Emissions().GetStakeRemovalsUpUntilBlock(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetStakeReputerAuthority(ctx context.Context, req *emissionstypes.GetStakeReputerAuthorityRequest, opts ...config.CallOpt) (*emissionstypes.GetStakeReputerAuthorityResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetStakeReputerAuthorityResponse, error) {
		return client.Emissions().GetStakeReputerAuthority(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetTopic(ctx context.Context, req *emissionstypes.GetTopicRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetTopicResponse, error) {
		return client.Emissions().GetTopic(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetTopicFeeRevenue(ctx context.Context, req *emissionstypes.GetTopicFeeRevenueRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicFeeRevenueResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetTopicFeeRevenueResponse, error) {
		return client.Emissions().GetTopicFeeRevenue(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetTopicInitialForecasterEmaScore(ctx context.Context, req *emissionstypes.GetTopicInitialForecasterEmaScoreRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicInitialForecasterEmaScoreResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetTopicInitialForecasterEmaScoreResponse, error) {
		return client.Emissions().GetTopicInitialForecasterEmaScore(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetTopicInitialInfererEmaScore(ctx context.Context, req *emissionstypes.GetTopicInitialInfererEmaScoreRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicInitialInfererEmaScoreResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetTopicInitialInfererEmaScoreResponse, error) {
		return client.Emissions().GetTopicInitialInfererEmaScore(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetTopicInitialReputerEmaScore(ctx context.Context, req *emissionstypes.GetTopicInitialReputerEmaScoreRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicInitialReputerEmaScoreResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetTopicInitialReputerEmaScoreResponse, error) {
		return client.Emissions().GetTopicInitialReputerEmaScore(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetTopicLastReputerCommitInfo(ctx context.Context, req *emissionstypes.GetTopicLastReputerCommitInfoRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicLastReputerCommitInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetTopicLastReputerCommitInfoResponse, error) {
		return client.Emissions().GetTopicLastReputerCommitInfo(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetTopicLastWorkerCommitInfo(ctx context.Context, req *emissionstypes.GetTopicLastWorkerCommitInfoRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicLastWorkerCommitInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetTopicLastWorkerCommitInfoResponse, error) {
		return client.Emissions().GetTopicLastWorkerCommitInfo(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetTopicRewardNonce(ctx context.Context, req *emissionstypes.GetTopicRewardNonceRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicRewardNonceResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetTopicRewardNonceResponse, error) {
		return client.Emissions().GetTopicRewardNonce(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetTopicStake(ctx context.Context, req *emissionstypes.GetTopicStakeRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicStakeResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetTopicStakeResponse, error) {
		return client.Emissions().GetTopicStake(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetTotalRewardToDistribute(ctx context.Context, req *emissionstypes.GetTotalRewardToDistributeRequest, opts ...config.CallOpt) (*emissionstypes.GetTotalRewardToDistributeResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetTotalRewardToDistributeResponse, error) {
		return client.Emissions().GetTotalRewardToDistribute(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetTotalStake(ctx context.Context, req *emissionstypes.GetTotalStakeRequest, opts ...config.CallOpt) (*emissionstypes.GetTotalStakeResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetTotalStakeResponse, error) {
		return client.Emissions().GetTotalStake(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetTotalSumPreviousTopicWeights(ctx context.Context, req *emissionstypes.GetTotalSumPreviousTopicWeightsRequest, opts ...config.CallOpt) (*emissionstypes.GetTotalSumPreviousTopicWeightsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetTotalSumPreviousTopicWeightsResponse, error) {
		return client.Emissions().GetTotalSumPreviousTopicWeights(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetUnfulfilledReputerNonces(ctx context.Context, req *emissionstypes.GetUnfulfilledReputerNoncesRequest, opts ...config.CallOpt) (*emissionstypes.GetUnfulfilledReputerNoncesResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetUnfulfilledReputerNoncesResponse, error) {
		return client.Emissions().GetUnfulfilledReputerNonces(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetUnfulfilledWorkerNonces(ctx context.Context, req *emissionstypes.GetUnfulfilledWorkerNoncesRequest, opts ...config.CallOpt) (*emissionstypes.GetUnfulfilledWorkerNoncesResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetUnfulfilledWorkerNoncesResponse, error) {
		return client.Emissions().GetUnfulfilledWorkerNonces(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetWorkerForecastScoresAtBlock(ctx context.Context, req *emissionstypes.GetWorkerForecastScoresAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetWorkerForecastScoresAtBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetWorkerForecastScoresAtBlockResponse, error) {
		return client.Emissions().GetWorkerForecastScoresAtBlock(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetWorkerInferenceScoresAtBlock(ctx context.Context, req *emissionstypes.GetWorkerInferenceScoresAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetWorkerInferenceScoresAtBlockResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetWorkerInferenceScoresAtBlockResponse, error) {
		return client.Emissions().GetWorkerInferenceScoresAtBlock(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetWorkerLatestInferenceByTopicId(ctx context.Context, req *emissionstypes.GetWorkerLatestInferenceByTopicIdRequest, opts ...config.CallOpt) (*emissionstypes.GetWorkerLatestInferenceByTopicIdResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetWorkerLatestInferenceByTopicIdResponse, error) {
		return client.Emissions().GetWorkerLatestInferenceByTopicId(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetWorkerNodeInfo(ctx context.Context, req *emissionstypes.GetWorkerNodeInfoRequest, opts ...config.CallOpt) (*emissionstypes.GetWorkerNodeInfoResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetWorkerNodeInfoResponse, error) {
		return client.Emissions().GetWorkerNodeInfo(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) GetWorkerSubmissionWindowStatus(ctx context.Context, req *emissionstypes.GetWorkerSubmissionWindowStatusRequest, opts ...config.CallOpt) (*emissionstypes.GetWorkerSubmissionWindowStatusResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.GetWorkerSubmissionWindowStatusResponse, error) {
		return client.Emissions().GetWorkerSubmissionWindowStatus(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) IsReputerNonceUnfulfilled(ctx context.Context, req *emissionstypes.IsReputerNonceUnfulfilledRequest, opts ...config.CallOpt) (*emissionstypes.IsReputerNonceUnfulfilledResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.IsReputerNonceUnfulfilledResponse, error) {
		return client.Emissions().IsReputerNonceUnfulfilled(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) IsReputerRegisteredInTopicId(ctx context.Context, req *emissionstypes.IsReputerRegisteredInTopicIdRequest, opts ...config.CallOpt) (*emissionstypes.IsReputerRegisteredInTopicIdResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.IsReputerRegisteredInTopicIdResponse, error) {
		return client.Emissions().IsReputerRegisteredInTopicId(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) IsTopicActive(ctx context.Context, req *emissionstypes.IsTopicActiveRequest, opts ...config.CallOpt) (*emissionstypes.IsTopicActiveResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.IsTopicActiveResponse, error) {
		return client.Emissions().IsTopicActive(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) IsTopicReputerWhitelistEnabled(ctx context.Context, req *emissionstypes.IsTopicReputerWhitelistEnabledRequest, opts ...config.CallOpt) (*emissionstypes.IsTopicReputerWhitelistEnabledResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.IsTopicReputerWhitelistEnabledResponse, error) {
		return client.Emissions().IsTopicReputerWhitelistEnabled(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) IsTopicWorkerWhitelistEnabled(ctx context.Context, req *emissionstypes.IsTopicWorkerWhitelistEnabledRequest, opts ...config.CallOpt) (*emissionstypes.IsTopicWorkerWhitelistEnabledResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.IsTopicWorkerWhitelistEnabledResponse, error) {
		return client.Emissions().IsTopicWorkerWhitelistEnabled(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) IsWhitelistAdmin(ctx context.Context, req *emissionstypes.IsWhitelistAdminRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistAdminResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.IsWhitelistAdminResponse, error) {
		return client.Emissions().IsWhitelistAdmin(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) IsWhitelistedGlobalActor(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalActorRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistedGlobalActorResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.IsWhitelistedGlobalActorResponse, error) {
		return client.Emissions().IsWhitelistedGlobalActor(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) IsWhitelistedGlobalAdmin(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalAdminRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistedGlobalAdminResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.IsWhitelistedGlobalAdminResponse, error) {
		return client.Emissions().IsWhitelistedGlobalAdmin(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) IsWhitelistedGlobalReputer(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalReputerRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistedGlobalReputerResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.IsWhitelistedGlobalReputerResponse, error) {
		return client.Emissions().IsWhitelistedGlobalReputer(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) IsWhitelistedGlobalWorker(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalWorkerRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistedGlobalWorkerResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.IsWhitelistedGlobalWorkerResponse, error) {
		return client.Emissions().IsWhitelistedGlobalWorker(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) IsWhitelistedTopicCreator(ctx context.Context, req *emissionstypes.IsWhitelistedTopicCreatorRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistedTopicCreatorResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.IsWhitelistedTopicCreatorResponse, error) {
		return client.Emissions().IsWhitelistedTopicCreator(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) IsWhitelistedTopicReputer(ctx context.Context, req *emissionstypes.IsWhitelistedTopicReputerRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistedTopicReputerResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.IsWhitelistedTopicReputerResponse, error) {
		return client.Emissions().IsWhitelistedTopicReputer(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) IsWhitelistedTopicWorker(ctx context.Context, req *emissionstypes.IsWhitelistedTopicWorkerRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistedTopicWorkerResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.IsWhitelistedTopicWorkerResponse, error) {
		return client.Emissions().IsWhitelistedTopicWorker(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) IsWorkerNonceUnfulfilled(ctx context.Context, req *emissionstypes.IsWorkerNonceUnfulfilledRequest, opts ...config.CallOpt) (*emissionstypes.IsWorkerNonceUnfulfilledResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.IsWorkerNonceUnfulfilledResponse, error) {
		return client.Emissions().IsWorkerNonceUnfulfilled(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) IsWorkerRegisteredInTopicId(ctx context.Context, req *emissionstypes.IsWorkerRegisteredInTopicIdRequest, opts ...config.CallOpt) (*emissionstypes.IsWorkerRegisteredInTopicIdResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.IsWorkerRegisteredInTopicIdResponse, error) {
		return client.Emissions().IsWorkerRegisteredInTopicId(ctx, req, opts...)
	})
}

func (c *EmissionsClientWrapper) TopicExists(ctx context.Context, req *emissionstypes.TopicExistsRequest, opts ...config.CallOpt) (*emissionstypes.TopicExistsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*emissionstypes.TopicExistsResponse, error) {
		return client.Emissions().TopicExists(ctx, req, opts...)
	})
}
