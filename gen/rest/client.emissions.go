package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	emissionstypes "github.com/allora-network/allora-chain/x/emissions/types"
)

// EmissionsRESTClient provides REST access to the emissions module
type EmissionsRESTClient struct {
	*RESTClientCore
	logger zerolog.Logger
}

// NewEmissionsRESTClient creates a new emissions REST client
func NewEmissionsRESTClient(core *RESTClientCore, logger zerolog.Logger) *EmissionsRESTClient {
	return &EmissionsRESTClient{
		RESTClientCore: core,
		logger:         logger.With().Str("module", "emissions").Str("protocol", "rest").Logger(),
	}
}

func (c *EmissionsRESTClient) GetParams(ctx context.Context, req *emissionstypes.GetParamsRequest) (*emissionstypes.GetParamsResponse, error) {
	resp := &emissionstypes.GetParamsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/params",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetParams")
}

func (c *EmissionsRESTClient) GetNextTopicId(ctx context.Context, req *emissionstypes.GetNextTopicIdRequest) (*emissionstypes.GetNextTopicIdResponse, error) {
	resp := &emissionstypes.GetNextTopicIdResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/next_topic_id",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetNextTopicId")
}

func (c *EmissionsRESTClient) GetTopic(ctx context.Context, req *emissionstypes.GetTopicRequest) (*emissionstypes.GetTopicResponse, error) {
	resp := &emissionstypes.GetTopicResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/topics/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetTopic")
}

func (c *EmissionsRESTClient) GetWorkerLatestInferenceByTopicId(ctx context.Context, req *emissionstypes.GetWorkerLatestInferenceByTopicIdRequest) (*emissionstypes.GetWorkerLatestInferenceByTopicIdResponse, error) {
	resp := &emissionstypes.GetWorkerLatestInferenceByTopicIdResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/topics/{topic_id}/workers/{worker_address}/latest_inference",
		[]string{"topic_id", "worker_address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetWorkerLatestInferenceByTopicId")
}

func (c *EmissionsRESTClient) GetInferencesAtBlock(ctx context.Context, req *emissionstypes.GetInferencesAtBlockRequest) (*emissionstypes.GetInferencesAtBlockResponse, error) {
	resp := &emissionstypes.GetInferencesAtBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/inferences/{topic_id}/{block_height}",
		[]string{"topic_id", "block_height"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetInferencesAtBlock")
}

func (c *EmissionsRESTClient) GetLatestTopicInferences(ctx context.Context, req *emissionstypes.GetLatestTopicInferencesRequest) (*emissionstypes.GetLatestTopicInferencesResponse, error) {
	resp := &emissionstypes.GetLatestTopicInferencesResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/latest_inferences/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetLatestTopicInferences")
}

func (c *EmissionsRESTClient) GetForecastsAtBlock(ctx context.Context, req *emissionstypes.GetForecastsAtBlockRequest) (*emissionstypes.GetForecastsAtBlockResponse, error) {
	resp := &emissionstypes.GetForecastsAtBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/forecasts/{topic_id}/{block_height}",
		[]string{"topic_id", "block_height"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetForecastsAtBlock")
}

func (c *EmissionsRESTClient) GetNetworkLossBundleAtBlock(ctx context.Context, req *emissionstypes.GetNetworkLossBundleAtBlockRequest) (*emissionstypes.GetNetworkLossBundleAtBlockResponse, error) {
	resp := &emissionstypes.GetNetworkLossBundleAtBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/network_loss/{topic_id}/{block_height}",
		[]string{"topic_id", "block_height"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetNetworkLossBundleAtBlock")
}

func (c *EmissionsRESTClient) GetTotalStake(ctx context.Context, req *emissionstypes.GetTotalStakeRequest) (*emissionstypes.GetTotalStakeResponse, error) {
	resp := &emissionstypes.GetTotalStakeResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/total_stake",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetTotalStake")
}

func (c *EmissionsRESTClient) GetReputerStakeInTopic(ctx context.Context, req *emissionstypes.GetReputerStakeInTopicRequest) (*emissionstypes.GetReputerStakeInTopicResponse, error) {
	resp := &emissionstypes.GetReputerStakeInTopicResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/reputer_stake/{address}/{topic_id}",
		[]string{"address", "topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetReputerStakeInTopic")
}

func (c *EmissionsRESTClient) GetMultiReputerStakeInTopic(ctx context.Context, req *emissionstypes.GetMultiReputerStakeInTopicRequest) (*emissionstypes.GetMultiReputerStakeInTopicResponse, error) {
	resp := &emissionstypes.GetMultiReputerStakeInTopicResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/reputers_stakes/{topic_id}",
		[]string{"topic_id"}, []string{"addresses"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetMultiReputerStakeInTopic")
}

func (c *EmissionsRESTClient) GetStakeFromReputerInTopicInSelf(ctx context.Context, req *emissionstypes.GetStakeFromReputerInTopicInSelfRequest) (*emissionstypes.GetStakeFromReputerInTopicInSelfResponse, error) {
	resp := &emissionstypes.GetStakeFromReputerInTopicInSelfResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/reputer_stake_self/{reputer_address}/{topic_id}",
		[]string{"reputer_address", "topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetStakeFromReputerInTopicInSelf")
}

func (c *EmissionsRESTClient) GetDelegateStakeInTopicInReputer(ctx context.Context, req *emissionstypes.GetDelegateStakeInTopicInReputerRequest) (*emissionstypes.GetDelegateStakeInTopicInReputerResponse, error) {
	resp := &emissionstypes.GetDelegateStakeInTopicInReputerResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/reputer_delegate_stake/{reputer_address}/{topic_id}",
		[]string{"reputer_address", "topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetDelegateStakeInTopicInReputer")
}

func (c *EmissionsRESTClient) GetStakeFromDelegatorInTopicInReputer(ctx context.Context, req *emissionstypes.GetStakeFromDelegatorInTopicInReputerRequest) (*emissionstypes.GetStakeFromDelegatorInTopicInReputerResponse, error) {
	resp := &emissionstypes.GetStakeFromDelegatorInTopicInReputerResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/delegate_stake/{delegator_address}/{reputer_address}/{topic_id}",
		[]string{"delegator_address", "reputer_address", "topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetStakeFromDelegatorInTopicInReputer")
}

func (c *EmissionsRESTClient) GetStakeFromDelegatorInTopic(ctx context.Context, req *emissionstypes.GetStakeFromDelegatorInTopicRequest) (*emissionstypes.GetStakeFromDelegatorInTopicResponse, error) {
	resp := &emissionstypes.GetStakeFromDelegatorInTopicResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/delegate_stake/{delegator_address}/{topic_id}",
		[]string{"delegator_address", "topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetStakeFromDelegatorInTopic")
}

func (c *EmissionsRESTClient) GetTopicStake(ctx context.Context, req *emissionstypes.GetTopicStakeRequest) (*emissionstypes.GetTopicStakeResponse, error) {
	resp := &emissionstypes.GetTopicStakeResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/stake/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetTopicStake")
}

func (c *EmissionsRESTClient) GetStakeRemovalsUpUntilBlock(ctx context.Context, req *emissionstypes.GetStakeRemovalsUpUntilBlockRequest) (*emissionstypes.GetStakeRemovalsUpUntilBlockResponse, error) {
	resp := &emissionstypes.GetStakeRemovalsUpUntilBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/stake_removals/{block_height}",
		[]string{"block_height"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetStakeRemovalsUpUntilBlock")
}

func (c *EmissionsRESTClient) GetDelegateStakeRemovalsUpUntilBlock(ctx context.Context, req *emissionstypes.GetDelegateStakeRemovalsUpUntilBlockRequest) (*emissionstypes.GetDelegateStakeRemovalsUpUntilBlockResponse, error) {
	resp := &emissionstypes.GetDelegateStakeRemovalsUpUntilBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/delegate_stake_removals/{block_height}",
		[]string{"block_height"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetDelegateStakeRemovalsUpUntilBlock")
}

func (c *EmissionsRESTClient) GetStakeRemovalInfo(ctx context.Context, req *emissionstypes.GetStakeRemovalInfoRequest) (*emissionstypes.GetStakeRemovalInfoResponse, error) {
	resp := &emissionstypes.GetStakeRemovalInfoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/stake_removal/{topic_id}/{reputer}",
		[]string{"topic_id", "reputer"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetStakeRemovalInfo")
}

func (c *EmissionsRESTClient) GetDelegateStakeRemovalInfo(ctx context.Context, req *emissionstypes.GetDelegateStakeRemovalInfoRequest) (*emissionstypes.GetDelegateStakeRemovalInfoResponse, error) {
	resp := &emissionstypes.GetDelegateStakeRemovalInfoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/delegate_stake_removal/{topic_id}/{delegator}/{reputer}",
		[]string{"topic_id", "delegator", "reputer"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetDelegateStakeRemovalInfo")
}

func (c *EmissionsRESTClient) GetWorkerNodeInfo(ctx context.Context, req *emissionstypes.GetWorkerNodeInfoRequest) (*emissionstypes.GetWorkerNodeInfoResponse, error) {
	resp := &emissionstypes.GetWorkerNodeInfoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/worker/{address}",
		[]string{"address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetWorkerNodeInfo")
}

func (c *EmissionsRESTClient) GetReputerNodeInfo(ctx context.Context, req *emissionstypes.GetReputerNodeInfoRequest) (*emissionstypes.GetReputerNodeInfoResponse, error) {
	resp := &emissionstypes.GetReputerNodeInfoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/reputer/{address}",
		[]string{"address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetReputerNodeInfo")
}

func (c *EmissionsRESTClient) IsWorkerRegisteredInTopicId(ctx context.Context, req *emissionstypes.IsWorkerRegisteredInTopicIdRequest) (*emissionstypes.IsWorkerRegisteredInTopicIdResponse, error) {
	resp := &emissionstypes.IsWorkerRegisteredInTopicIdResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/worker_registered/{topic_id}/{address}",
		[]string{"topic_id", "address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.IsWorkerRegisteredInTopicId")
}

func (c *EmissionsRESTClient) IsReputerRegisteredInTopicId(ctx context.Context, req *emissionstypes.IsReputerRegisteredInTopicIdRequest) (*emissionstypes.IsReputerRegisteredInTopicIdResponse, error) {
	resp := &emissionstypes.IsReputerRegisteredInTopicIdResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/reputer_registered/{topic_id}/{address}",
		[]string{"topic_id", "address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.IsReputerRegisteredInTopicId")
}

func (c *EmissionsRESTClient) GetNetworkInferencesAtBlock(ctx context.Context, req *emissionstypes.GetNetworkInferencesAtBlockRequest) (*emissionstypes.GetNetworkInferencesAtBlockResponse, error) {
	resp := &emissionstypes.GetNetworkInferencesAtBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/network_inferences/{topic_id}/last_inference/{block_height_last_inference}",
		[]string{"topic_id", "block_height_last_inference"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetNetworkInferencesAtBlock")
}

func (c *EmissionsRESTClient) GetNetworkInferencesAtBlockOutlierResistant(ctx context.Context, req *emissionstypes.GetNetworkInferencesAtBlockOutlierResistantRequest) (*emissionstypes.GetNetworkInferencesAtBlockOutlierResistantResponse, error) {
	resp := &emissionstypes.GetNetworkInferencesAtBlockOutlierResistantResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/network_inferences_outlier_resistant/{topic_id}/last_inference/{block_height_last_inference}",
		[]string{"topic_id", "block_height_last_inference"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetNetworkInferencesAtBlockOutlierResistant")
}

func (c *EmissionsRESTClient) GetLatestNetworkInferences(ctx context.Context, req *emissionstypes.GetLatestNetworkInferencesRequest) (*emissionstypes.GetLatestNetworkInferencesResponse, error) {
	resp := &emissionstypes.GetLatestNetworkInferencesResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/latest_network_inferences/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetLatestNetworkInferences")
}

func (c *EmissionsRESTClient) GetLatestNetworkInferencesOutlierResistant(ctx context.Context, req *emissionstypes.GetLatestNetworkInferencesOutlierResistantRequest) (*emissionstypes.GetLatestNetworkInferencesOutlierResistantResponse, error) {
	resp := &emissionstypes.GetLatestNetworkInferencesOutlierResistantResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/latest_network_inferences_outlier_resistant/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetLatestNetworkInferencesOutlierResistant")
}

func (c *EmissionsRESTClient) IsWorkerNonceUnfulfilled(ctx context.Context, req *emissionstypes.IsWorkerNonceUnfulfilledRequest) (*emissionstypes.IsWorkerNonceUnfulfilledResponse, error) {
	resp := &emissionstypes.IsWorkerNonceUnfulfilledResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_worker_nonce_unfulfilled/{topic_id}/{block_height}",
		[]string{"topic_id", "block_height"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.IsWorkerNonceUnfulfilled")
}

func (c *EmissionsRESTClient) IsReputerNonceUnfulfilled(ctx context.Context, req *emissionstypes.IsReputerNonceUnfulfilledRequest) (*emissionstypes.IsReputerNonceUnfulfilledResponse, error) {
	resp := &emissionstypes.IsReputerNonceUnfulfilledResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_reputer_nonce_unfulfilled/{topic_id}/{block_height}",
		[]string{"topic_id", "block_height"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.IsReputerNonceUnfulfilled")
}

func (c *EmissionsRESTClient) GetUnfulfilledWorkerNonces(ctx context.Context, req *emissionstypes.GetUnfulfilledWorkerNoncesRequest) (*emissionstypes.GetUnfulfilledWorkerNoncesResponse, error) {
	resp := &emissionstypes.GetUnfulfilledWorkerNoncesResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/unfulfilled_worker_nonces/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetUnfulfilledWorkerNonces")
}

func (c *EmissionsRESTClient) GetUnfulfilledReputerNonces(ctx context.Context, req *emissionstypes.GetUnfulfilledReputerNoncesRequest) (*emissionstypes.GetUnfulfilledReputerNoncesResponse, error) {
	resp := &emissionstypes.GetUnfulfilledReputerNoncesResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/unfulfilled_reputer_nonces/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetUnfulfilledReputerNonces")
}

func (c *EmissionsRESTClient) GetInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetInfererNetworkRegretRequest) (*emissionstypes.GetInfererNetworkRegretResponse, error) {
	resp := &emissionstypes.GetInfererNetworkRegretResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/inferer_network_regret/{topic_id}/{actor_id}",
		[]string{"topic_id", "actor_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetInfererNetworkRegret")
}

func (c *EmissionsRESTClient) GetForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetForecasterNetworkRegretRequest) (*emissionstypes.GetForecasterNetworkRegretResponse, error) {
	resp := &emissionstypes.GetForecasterNetworkRegretResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/forecaster_network_regret/{topic_id}/{worker}",
		[]string{"topic_id", "worker"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetForecasterNetworkRegret")
}

func (c *EmissionsRESTClient) GetOneInForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetOneInForecasterNetworkRegretRequest) (*emissionstypes.GetOneInForecasterNetworkRegretResponse, error) {
	resp := &emissionstypes.GetOneInForecasterNetworkRegretResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/one_in_forecaster_network_regret/{topic_id}/{forecaster}/{inferer}",
		[]string{"topic_id", "forecaster", "inferer"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetOneInForecasterNetworkRegret")
}

func (c *EmissionsRESTClient) IsWhitelistAdmin(ctx context.Context, req *emissionstypes.IsWhitelistAdminRequest) (*emissionstypes.IsWhitelistAdminResponse, error) {
	resp := &emissionstypes.IsWhitelistAdminResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/whitelist_admin/{address}",
		[]string{"address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.IsWhitelistAdmin")
}

func (c *EmissionsRESTClient) GetTopicLastWorkerCommitInfo(ctx context.Context, req *emissionstypes.GetTopicLastWorkerCommitInfoRequest) (*emissionstypes.GetTopicLastWorkerCommitInfoResponse, error) {
	resp := &emissionstypes.GetTopicLastWorkerCommitInfoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/topic_last_worker_commit_info/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetTopicLastWorkerCommitInfo")
}

func (c *EmissionsRESTClient) GetTopicLastReputerCommitInfo(ctx context.Context, req *emissionstypes.GetTopicLastReputerCommitInfoRequest) (*emissionstypes.GetTopicLastReputerCommitInfoResponse, error) {
	resp := &emissionstypes.GetTopicLastReputerCommitInfoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/topic_last_reputer_commit_info/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetTopicLastReputerCommitInfo")
}

func (c *EmissionsRESTClient) GetTopicRewardNonce(ctx context.Context, req *emissionstypes.GetTopicRewardNonceRequest) (*emissionstypes.GetTopicRewardNonceResponse, error) {
	resp := &emissionstypes.GetTopicRewardNonceResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/topic_reward_nonce/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetTopicRewardNonce")
}

func (c *EmissionsRESTClient) GetReputerLossBundlesAtBlock(ctx context.Context, req *emissionstypes.GetReputerLossBundlesAtBlockRequest) (*emissionstypes.GetReputerLossBundlesAtBlockResponse, error) {
	resp := &emissionstypes.GetReputerLossBundlesAtBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/reputer_loss_bundles/{topic_id}/{block_height}",
		[]string{"topic_id", "block_height"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetReputerLossBundlesAtBlock")
}

func (c *EmissionsRESTClient) GetStakeReputerAuthority(ctx context.Context, req *emissionstypes.GetStakeReputerAuthorityRequest) (*emissionstypes.GetStakeReputerAuthorityResponse, error) {
	resp := &emissionstypes.GetStakeReputerAuthorityResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/stake_reputer_authority/{topic_id}/{reputer}",
		[]string{"topic_id", "reputer"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetStakeReputerAuthority")
}

func (c *EmissionsRESTClient) GetDelegateStakePlacement(ctx context.Context, req *emissionstypes.GetDelegateStakePlacementRequest) (*emissionstypes.GetDelegateStakePlacementResponse, error) {
	resp := &emissionstypes.GetDelegateStakePlacementResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/delegate_stake_placement/{topic_id}/{delegator}/{target}",
		[]string{"topic_id", "delegator", "target"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetDelegateStakePlacement")
}

func (c *EmissionsRESTClient) GetDelegateStakeUponReputer(ctx context.Context, req *emissionstypes.GetDelegateStakeUponReputerRequest) (*emissionstypes.GetDelegateStakeUponReputerResponse, error) {
	resp := &emissionstypes.GetDelegateStakeUponReputerResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/delegate_stake_upon_reputer/{topic_id}/{target}",
		[]string{"topic_id", "target"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetDelegateStakeUponReputer")
}

func (c *EmissionsRESTClient) GetDelegateRewardPerShare(ctx context.Context, req *emissionstypes.GetDelegateRewardPerShareRequest) (*emissionstypes.GetDelegateRewardPerShareResponse, error) {
	resp := &emissionstypes.GetDelegateRewardPerShareResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/delegate_reward_per_share/{topic_id}/{reputer}",
		[]string{"topic_id", "reputer"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetDelegateRewardPerShare")
}

func (c *EmissionsRESTClient) GetStakeRemovalForReputerAndTopicId(ctx context.Context, req *emissionstypes.GetStakeRemovalForReputerAndTopicIdRequest) (*emissionstypes.GetStakeRemovalForReputerAndTopicIdResponse, error) {
	resp := &emissionstypes.GetStakeRemovalForReputerAndTopicIdResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/stake_removal/{reputer}/{topic_id}",
		[]string{"reputer", "topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetStakeRemovalForReputerAndTopicId")
}

func (c *EmissionsRESTClient) GetDelegateStakeRemoval(ctx context.Context, req *emissionstypes.GetDelegateStakeRemovalRequest) (*emissionstypes.GetDelegateStakeRemovalResponse, error) {
	resp := &emissionstypes.GetDelegateStakeRemovalResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/delegate_stake_removal/{block_height}/{topic_id}/{delegator}/{reputer}",
		[]string{"block_height", "topic_id", "delegator", "reputer"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetDelegateStakeRemoval")
}

func (c *EmissionsRESTClient) GetPreviousTopicWeight(ctx context.Context, req *emissionstypes.GetPreviousTopicWeightRequest) (*emissionstypes.GetPreviousTopicWeightResponse, error) {
	resp := &emissionstypes.GetPreviousTopicWeightResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/previous_topic_weight/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetPreviousTopicWeight")
}

func (c *EmissionsRESTClient) GetTotalSumPreviousTopicWeights(ctx context.Context, req *emissionstypes.GetTotalSumPreviousTopicWeightsRequest) (*emissionstypes.GetTotalSumPreviousTopicWeightsResponse, error) {
	resp := &emissionstypes.GetTotalSumPreviousTopicWeightsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/sum_previous_total_topic_weight",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetTotalSumPreviousTopicWeights")
}

func (c *EmissionsRESTClient) TopicExists(ctx context.Context, req *emissionstypes.TopicExistsRequest) (*emissionstypes.TopicExistsResponse, error) {
	resp := &emissionstypes.TopicExistsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/topic_exists/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.TopicExists")
}

func (c *EmissionsRESTClient) IsTopicActive(ctx context.Context, req *emissionstypes.IsTopicActiveRequest) (*emissionstypes.IsTopicActiveResponse, error) {
	resp := &emissionstypes.IsTopicActiveResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_topic_active/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.IsTopicActive")
}

func (c *EmissionsRESTClient) GetTopicFeeRevenue(ctx context.Context, req *emissionstypes.GetTopicFeeRevenueRequest) (*emissionstypes.GetTopicFeeRevenueResponse, error) {
	resp := &emissionstypes.GetTopicFeeRevenueResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/topic_fee_revenue/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetTopicFeeRevenue")
}

func (c *EmissionsRESTClient) GetInfererScoreEma(ctx context.Context, req *emissionstypes.GetInfererScoreEmaRequest) (*emissionstypes.GetInfererScoreEmaResponse, error) {
	resp := &emissionstypes.GetInfererScoreEmaResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/inferer_score_ema/{topic_id}/{inferer}",
		[]string{"topic_id", "inferer"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetInfererScoreEma")
}

func (c *EmissionsRESTClient) GetForecasterScoreEma(ctx context.Context, req *emissionstypes.GetForecasterScoreEmaRequest) (*emissionstypes.GetForecasterScoreEmaResponse, error) {
	resp := &emissionstypes.GetForecasterScoreEmaResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/forecaster_score_ema/{topic_id}/{forecaster}",
		[]string{"topic_id", "forecaster"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetForecasterScoreEma")
}

func (c *EmissionsRESTClient) GetReputerScoreEma(ctx context.Context, req *emissionstypes.GetReputerScoreEmaRequest) (*emissionstypes.GetReputerScoreEmaResponse, error) {
	resp := &emissionstypes.GetReputerScoreEmaResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/reputer_score_ema/{topic_id}/{reputer}",
		[]string{"topic_id", "reputer"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetReputerScoreEma")
}

func (c *EmissionsRESTClient) GetInferenceScoresUntilBlock(ctx context.Context, req *emissionstypes.GetInferenceScoresUntilBlockRequest) (*emissionstypes.GetInferenceScoresUntilBlockResponse, error) {
	resp := &emissionstypes.GetInferenceScoresUntilBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/inference_scores_until_block/{topic_id}/{block_height}",
		[]string{"topic_id", "block_height"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetInferenceScoresUntilBlock")
}

func (c *EmissionsRESTClient) GetPreviousTopicQuantileForecasterScoreEma(ctx context.Context, req *emissionstypes.GetPreviousTopicQuantileForecasterScoreEmaRequest) (*emissionstypes.GetPreviousTopicQuantileForecasterScoreEmaResponse, error) {
	resp := &emissionstypes.GetPreviousTopicQuantileForecasterScoreEmaResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/topic_quantile_forecaster_score_ema/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetPreviousTopicQuantileForecasterScoreEma")
}

func (c *EmissionsRESTClient) GetPreviousTopicQuantileInfererScoreEma(ctx context.Context, req *emissionstypes.GetPreviousTopicQuantileInfererScoreEmaRequest) (*emissionstypes.GetPreviousTopicQuantileInfererScoreEmaResponse, error) {
	resp := &emissionstypes.GetPreviousTopicQuantileInfererScoreEmaResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/topic_quantile_inferer_score_ema/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetPreviousTopicQuantileInfererScoreEma")
}

func (c *EmissionsRESTClient) GetPreviousTopicQuantileReputerScoreEma(ctx context.Context, req *emissionstypes.GetPreviousTopicQuantileReputerScoreEmaRequest) (*emissionstypes.GetPreviousTopicQuantileReputerScoreEmaResponse, error) {
	resp := &emissionstypes.GetPreviousTopicQuantileReputerScoreEmaResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/topic_quantile_reputer_score_ema/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetPreviousTopicQuantileReputerScoreEma")
}

func (c *EmissionsRESTClient) GetWorkerInferenceScoresAtBlock(ctx context.Context, req *emissionstypes.GetWorkerInferenceScoresAtBlockRequest) (*emissionstypes.GetWorkerInferenceScoresAtBlockResponse, error) {
	resp := &emissionstypes.GetWorkerInferenceScoresAtBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/worker_inference_scores_at_block/{topic_id}/{block_height}",
		[]string{"topic_id", "block_height"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetWorkerInferenceScoresAtBlock")
}

func (c *EmissionsRESTClient) GetCurrentLowestInfererScore(ctx context.Context, req *emissionstypes.GetCurrentLowestInfererScoreRequest) (*emissionstypes.GetCurrentLowestInfererScoreResponse, error) {
	resp := &emissionstypes.GetCurrentLowestInfererScoreResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/current_lowest_inferer_score/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetCurrentLowestInfererScore")
}

func (c *EmissionsRESTClient) GetForecastScoresUntilBlock(ctx context.Context, req *emissionstypes.GetForecastScoresUntilBlockRequest) (*emissionstypes.GetForecastScoresUntilBlockResponse, error) {
	resp := &emissionstypes.GetForecastScoresUntilBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/forecast_scores_until_block/{topic_id}/{block_height}",
		[]string{"topic_id", "block_height"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetForecastScoresUntilBlock")
}

func (c *EmissionsRESTClient) GetWorkerForecastScoresAtBlock(ctx context.Context, req *emissionstypes.GetWorkerForecastScoresAtBlockRequest) (*emissionstypes.GetWorkerForecastScoresAtBlockResponse, error) {
	resp := &emissionstypes.GetWorkerForecastScoresAtBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/worker_forecast_scores_at_block/{topic_id}/{block_height}",
		[]string{"topic_id", "block_height"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetWorkerForecastScoresAtBlock")
}

func (c *EmissionsRESTClient) GetCurrentLowestForecasterScore(ctx context.Context, req *emissionstypes.GetCurrentLowestForecasterScoreRequest) (*emissionstypes.GetCurrentLowestForecasterScoreResponse, error) {
	resp := &emissionstypes.GetCurrentLowestForecasterScoreResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/current_lowest_forecaster_score/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetCurrentLowestForecasterScore")
}

func (c *EmissionsRESTClient) GetReputersScoresAtBlock(ctx context.Context, req *emissionstypes.GetReputersScoresAtBlockRequest) (*emissionstypes.GetReputersScoresAtBlockResponse, error) {
	resp := &emissionstypes.GetReputersScoresAtBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/reputers_scores_at_block/{topic_id}/{block_height}",
		[]string{"topic_id", "block_height"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetReputersScoresAtBlock")
}

func (c *EmissionsRESTClient) GetCurrentLowestReputerScore(ctx context.Context, req *emissionstypes.GetCurrentLowestReputerScoreRequest) (*emissionstypes.GetCurrentLowestReputerScoreResponse, error) {
	resp := &emissionstypes.GetCurrentLowestReputerScoreResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/current_lowest_reputer_score/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetCurrentLowestReputerScore")
}

func (c *EmissionsRESTClient) GetListeningCoefficient(ctx context.Context, req *emissionstypes.GetListeningCoefficientRequest) (*emissionstypes.GetListeningCoefficientResponse, error) {
	resp := &emissionstypes.GetListeningCoefficientResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/listening_coefficient/{topic_id}/{reputer}",
		[]string{"topic_id", "reputer"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetListeningCoefficient")
}

func (c *EmissionsRESTClient) GetPreviousReputerRewardFraction(ctx context.Context, req *emissionstypes.GetPreviousReputerRewardFractionRequest) (*emissionstypes.GetPreviousReputerRewardFractionResponse, error) {
	resp := &emissionstypes.GetPreviousReputerRewardFractionResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/previous_reputer_reward_fraction/{topic_id}/{reputer}",
		[]string{"topic_id", "reputer"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetPreviousReputerRewardFraction")
}

func (c *EmissionsRESTClient) GetPreviousInferenceRewardFraction(ctx context.Context, req *emissionstypes.GetPreviousInferenceRewardFractionRequest) (*emissionstypes.GetPreviousInferenceRewardFractionResponse, error) {
	resp := &emissionstypes.GetPreviousInferenceRewardFractionResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/previous_inference_reward_fraction/{topic_id}/{worker}",
		[]string{"topic_id", "worker"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetPreviousInferenceRewardFraction")
}

func (c *EmissionsRESTClient) GetPreviousForecastRewardFraction(ctx context.Context, req *emissionstypes.GetPreviousForecastRewardFractionRequest) (*emissionstypes.GetPreviousForecastRewardFractionResponse, error) {
	resp := &emissionstypes.GetPreviousForecastRewardFractionResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/previous_forecast_reward_fraction/{topic_id}/{worker}",
		[]string{"topic_id", "worker"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetPreviousForecastRewardFraction")
}

func (c *EmissionsRESTClient) GetPreviousPercentageRewardToStakedReputers(ctx context.Context, req *emissionstypes.GetPreviousPercentageRewardToStakedReputersRequest) (*emissionstypes.GetPreviousPercentageRewardToStakedReputersResponse, error) {
	resp := &emissionstypes.GetPreviousPercentageRewardToStakedReputersResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/previous_percentage_reward_to_staked_reputers",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetPreviousPercentageRewardToStakedReputers")
}

func (c *EmissionsRESTClient) GetTotalRewardToDistribute(ctx context.Context, req *emissionstypes.GetTotalRewardToDistributeRequest) (*emissionstypes.GetTotalRewardToDistributeResponse, error) {
	resp := &emissionstypes.GetTotalRewardToDistributeResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/total_reward_to_distribute",
		nil, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetTotalRewardToDistribute")
}

func (c *EmissionsRESTClient) GetNaiveInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetNaiveInfererNetworkRegretRequest) (*emissionstypes.GetNaiveInfererNetworkRegretResponse, error) {
	resp := &emissionstypes.GetNaiveInfererNetworkRegretResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/native_inferer_network_regret",
		nil, []string{"topic_id", "inferer"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetNaiveInfererNetworkRegret")
}

func (c *EmissionsRESTClient) GetOneOutInfererInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutInfererInfererNetworkRegretRequest) (*emissionstypes.GetOneOutInfererInfererNetworkRegretResponse, error) {
	resp := &emissionstypes.GetOneOutInfererInfererNetworkRegretResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/one_out_inferer_inferer_network_regret",
		nil, []string{"topic_id", "one_out_inferer", "inferer"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetOneOutInfererInfererNetworkRegret")
}

func (c *EmissionsRESTClient) GetOneOutInfererForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutInfererForecasterNetworkRegretRequest) (*emissionstypes.GetOneOutInfererForecasterNetworkRegretResponse, error) {
	resp := &emissionstypes.GetOneOutInfererForecasterNetworkRegretResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/one_out_inferer_forecaster_network_regret",
		nil, []string{"topic_id", "one_out_inferer", "forecaster"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetOneOutInfererForecasterNetworkRegret")
}

func (c *EmissionsRESTClient) GetOneOutForecasterInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutForecasterInfererNetworkRegretRequest) (*emissionstypes.GetOneOutForecasterInfererNetworkRegretResponse, error) {
	resp := &emissionstypes.GetOneOutForecasterInfererNetworkRegretResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/one_out_forecaster_inferer_network_regret",
		nil, []string{"topic_id", "one_out_forecaster", "inferer"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetOneOutForecasterInfererNetworkRegret")
}

func (c *EmissionsRESTClient) GetOneOutForecasterForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutForecasterForecasterNetworkRegretRequest) (*emissionstypes.GetOneOutForecasterForecasterNetworkRegretResponse, error) {
	resp := &emissionstypes.GetOneOutForecasterForecasterNetworkRegretResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/one_out_forecaster_forecaster_network_regret",
		nil, []string{"topic_id", "one_out_forecaster", "forecaster"},
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetOneOutForecasterForecasterNetworkRegret")
}

func (c *EmissionsRESTClient) GetActiveTopicsAtBlock(ctx context.Context, req *emissionstypes.GetActiveTopicsAtBlockRequest) (*emissionstypes.GetActiveTopicsAtBlockResponse, error) {
	resp := &emissionstypes.GetActiveTopicsAtBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/active_topics_at_block/{block_height}",
		[]string{"block_height"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetActiveTopicsAtBlock")
}

func (c *EmissionsRESTClient) GetNextChurningBlockByTopicId(ctx context.Context, req *emissionstypes.GetNextChurningBlockByTopicIdRequest) (*emissionstypes.GetNextChurningBlockByTopicIdResponse, error) {
	resp := &emissionstypes.GetNextChurningBlockByTopicIdResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/next_churning_block_by_topic_id/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetNextChurningBlockByTopicId")
}

func (c *EmissionsRESTClient) GetCountInfererInclusionsInTopic(ctx context.Context, req *emissionstypes.GetCountInfererInclusionsInTopicRequest) (*emissionstypes.GetCountInfererInclusionsInTopicResponse, error) {
	resp := &emissionstypes.GetCountInfererInclusionsInTopicResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/count_inferer_inclusions_in_topic/{topic_id}/{inferer}",
		[]string{"topic_id", "inferer"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetCountInfererInclusionsInTopic")
}

func (c *EmissionsRESTClient) GetCountForecasterInclusionsInTopic(ctx context.Context, req *emissionstypes.GetCountForecasterInclusionsInTopicRequest) (*emissionstypes.GetCountForecasterInclusionsInTopicResponse, error) {
	resp := &emissionstypes.GetCountForecasterInclusionsInTopicResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/count_forecaster_inclusions_in_topic/{topic_id}/{forecaster}",
		[]string{"topic_id", "forecaster"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetCountForecasterInclusionsInTopic")
}

func (c *EmissionsRESTClient) GetActiveReputersForTopic(ctx context.Context, req *emissionstypes.GetActiveReputersForTopicRequest) (*emissionstypes.GetActiveReputersForTopicResponse, error) {
	resp := &emissionstypes.GetActiveReputersForTopicResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/active_reputers/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetActiveReputersForTopic")
}

func (c *EmissionsRESTClient) GetActiveForecastersForTopic(ctx context.Context, req *emissionstypes.GetActiveForecastersForTopicRequest) (*emissionstypes.GetActiveForecastersForTopicResponse, error) {
	resp := &emissionstypes.GetActiveForecastersForTopicResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/active_forecasters/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetActiveForecastersForTopic")
}

func (c *EmissionsRESTClient) GetActiveInferersForTopic(ctx context.Context, req *emissionstypes.GetActiveInferersForTopicRequest) (*emissionstypes.GetActiveInferersForTopicResponse, error) {
	resp := &emissionstypes.GetActiveInferersForTopicResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/active_inferers/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetActiveInferersForTopic")
}

func (c *EmissionsRESTClient) IsWhitelistedGlobalWorker(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalWorkerRequest) (*emissionstypes.IsWhitelistedGlobalWorkerResponse, error) {
	resp := &emissionstypes.IsWhitelistedGlobalWorkerResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_whitelisted_global_worker/{address}",
		[]string{"address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.IsWhitelistedGlobalWorker")
}

func (c *EmissionsRESTClient) IsWhitelistedGlobalReputer(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalReputerRequest) (*emissionstypes.IsWhitelistedGlobalReputerResponse, error) {
	resp := &emissionstypes.IsWhitelistedGlobalReputerResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_whitelisted_global_reputer/{address}",
		[]string{"address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.IsWhitelistedGlobalReputer")
}

func (c *EmissionsRESTClient) IsWhitelistedGlobalAdmin(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalAdminRequest) (*emissionstypes.IsWhitelistedGlobalAdminResponse, error) {
	resp := &emissionstypes.IsWhitelistedGlobalAdminResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_whitelisted_global_admin/{address}",
		[]string{"address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.IsWhitelistedGlobalAdmin")
}

func (c *EmissionsRESTClient) IsTopicWorkerWhitelistEnabled(ctx context.Context, req *emissionstypes.IsTopicWorkerWhitelistEnabledRequest) (*emissionstypes.IsTopicWorkerWhitelistEnabledResponse, error) {
	resp := &emissionstypes.IsTopicWorkerWhitelistEnabledResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_topic_worker_whitelist_enabled/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.IsTopicWorkerWhitelistEnabled")
}

func (c *EmissionsRESTClient) IsTopicReputerWhitelistEnabled(ctx context.Context, req *emissionstypes.IsTopicReputerWhitelistEnabledRequest) (*emissionstypes.IsTopicReputerWhitelistEnabledResponse, error) {
	resp := &emissionstypes.IsTopicReputerWhitelistEnabledResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_topic_reputer_whitelist_enabled/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.IsTopicReputerWhitelistEnabled")
}

func (c *EmissionsRESTClient) IsWhitelistedTopicCreator(ctx context.Context, req *emissionstypes.IsWhitelistedTopicCreatorRequest) (*emissionstypes.IsWhitelistedTopicCreatorResponse, error) {
	resp := &emissionstypes.IsWhitelistedTopicCreatorResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_whitelisted_topic_creator/{address}",
		[]string{"address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.IsWhitelistedTopicCreator")
}

func (c *EmissionsRESTClient) IsWhitelistedGlobalActor(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalActorRequest) (*emissionstypes.IsWhitelistedGlobalActorResponse, error) {
	resp := &emissionstypes.IsWhitelistedGlobalActorResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_whitelisted_global_actor/{address}",
		[]string{"address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.IsWhitelistedGlobalActor")
}

func (c *EmissionsRESTClient) IsWhitelistedTopicWorker(ctx context.Context, req *emissionstypes.IsWhitelistedTopicWorkerRequest) (*emissionstypes.IsWhitelistedTopicWorkerResponse, error) {
	resp := &emissionstypes.IsWhitelistedTopicWorkerResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_whitelisted_topic_worker/{topic_id}/{address}",
		[]string{"topic_id", "address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.IsWhitelistedTopicWorker")
}

func (c *EmissionsRESTClient) IsWhitelistedTopicReputer(ctx context.Context, req *emissionstypes.IsWhitelistedTopicReputerRequest) (*emissionstypes.IsWhitelistedTopicReputerResponse, error) {
	resp := &emissionstypes.IsWhitelistedTopicReputerResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_whitelisted_topic_reputer/{topic_id}/{address}",
		[]string{"topic_id", "address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.IsWhitelistedTopicReputer")
}

func (c *EmissionsRESTClient) CanUpdateAllGlobalWhitelists(ctx context.Context, req *emissionstypes.CanUpdateAllGlobalWhitelistsRequest) (*emissionstypes.CanUpdateAllGlobalWhitelistsResponse, error) {
	resp := &emissionstypes.CanUpdateAllGlobalWhitelistsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/can_update_all_global_whitelists/{address}",
		[]string{"address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.CanUpdateAllGlobalWhitelists")
}

func (c *EmissionsRESTClient) CanUpdateGlobalWorkerWhitelist(ctx context.Context, req *emissionstypes.CanUpdateGlobalWorkerWhitelistRequest) (*emissionstypes.CanUpdateGlobalWorkerWhitelistResponse, error) {
	resp := &emissionstypes.CanUpdateGlobalWorkerWhitelistResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/can_update_global_worker_whitelist/{address}",
		[]string{"address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.CanUpdateGlobalWorkerWhitelist")
}

func (c *EmissionsRESTClient) CanUpdateGlobalReputerWhitelist(ctx context.Context, req *emissionstypes.CanUpdateGlobalReputerWhitelistRequest) (*emissionstypes.CanUpdateGlobalReputerWhitelistResponse, error) {
	resp := &emissionstypes.CanUpdateGlobalReputerWhitelistResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/can_update_global_reputer_whitelist/{address}",
		[]string{"address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.CanUpdateGlobalReputerWhitelist")
}

func (c *EmissionsRESTClient) CanUpdateParams(ctx context.Context, req *emissionstypes.CanUpdateParamsRequest) (*emissionstypes.CanUpdateParamsResponse, error) {
	resp := &emissionstypes.CanUpdateParamsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/can_update_params/{address}",
		[]string{"address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.CanUpdateParams")
}

func (c *EmissionsRESTClient) CanUpdateTopicWhitelist(ctx context.Context, req *emissionstypes.CanUpdateTopicWhitelistRequest) (*emissionstypes.CanUpdateTopicWhitelistResponse, error) {
	resp := &emissionstypes.CanUpdateTopicWhitelistResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/can_update_topic_whitelist/{topic_id}/{address}",
		[]string{"topic_id", "address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.CanUpdateTopicWhitelist")
}

func (c *EmissionsRESTClient) CanCreateTopic(ctx context.Context, req *emissionstypes.CanCreateTopicRequest) (*emissionstypes.CanCreateTopicResponse, error) {
	resp := &emissionstypes.CanCreateTopicResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/can_create_topic/{address}",
		[]string{"address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.CanCreateTopic")
}

func (c *EmissionsRESTClient) CanSubmitWorkerPayload(ctx context.Context, req *emissionstypes.CanSubmitWorkerPayloadRequest) (*emissionstypes.CanSubmitWorkerPayloadResponse, error) {
	resp := &emissionstypes.CanSubmitWorkerPayloadResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/can_submit_worker_payload/{topic_id}/{address}",
		[]string{"topic_id", "address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.CanSubmitWorkerPayload")
}

func (c *EmissionsRESTClient) CanSubmitReputerPayload(ctx context.Context, req *emissionstypes.CanSubmitReputerPayloadRequest) (*emissionstypes.CanSubmitReputerPayloadResponse, error) {
	resp := &emissionstypes.CanSubmitReputerPayloadResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/can_submit_reputer_payload/{topic_id}/{address}",
		[]string{"topic_id", "address"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.CanSubmitReputerPayload")
}

func (c *EmissionsRESTClient) GetTopicInitialInfererEmaScore(ctx context.Context, req *emissionstypes.GetTopicInitialInfererEmaScoreRequest) (*emissionstypes.GetTopicInitialInfererEmaScoreResponse, error) {
	resp := &emissionstypes.GetTopicInitialInfererEmaScoreResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/initial_inferer_ema_score/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetTopicInitialInfererEmaScore")
}

func (c *EmissionsRESTClient) GetTopicInitialForecasterEmaScore(ctx context.Context, req *emissionstypes.GetTopicInitialForecasterEmaScoreRequest) (*emissionstypes.GetTopicInitialForecasterEmaScoreResponse, error) {
	resp := &emissionstypes.GetTopicInitialForecasterEmaScoreResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/initial_forecaster_ema_score/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetTopicInitialForecasterEmaScore")
}

func (c *EmissionsRESTClient) GetTopicInitialReputerEmaScore(ctx context.Context, req *emissionstypes.GetTopicInitialReputerEmaScoreRequest) (*emissionstypes.GetTopicInitialReputerEmaScoreResponse, error) {
	resp := &emissionstypes.GetTopicInitialReputerEmaScoreResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/initial_reputer_ema_score/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetTopicInitialReputerEmaScore")
}

func (c *EmissionsRESTClient) GetLatestRegretStdNorm(ctx context.Context, req *emissionstypes.GetLatestRegretStdNormRequest) (*emissionstypes.GetLatestRegretStdNormResponse, error) {
	resp := &emissionstypes.GetLatestRegretStdNormResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/latest_regret_stdnorm/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetLatestRegretStdNorm")
}

func (c *EmissionsRESTClient) GetLatestInfererWeight(ctx context.Context, req *emissionstypes.GetLatestInfererWeightRequest) (*emissionstypes.GetLatestInfererWeightResponse, error) {
	resp := &emissionstypes.GetLatestInfererWeightResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/latest_inferer_weight/{topic_id}/{actor_id}",
		[]string{"topic_id", "actor_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetLatestInfererWeight")
}

func (c *EmissionsRESTClient) GetLatestForecasterWeight(ctx context.Context, req *emissionstypes.GetLatestForecasterWeightRequest) (*emissionstypes.GetLatestForecasterWeightResponse, error) {
	resp := &emissionstypes.GetLatestForecasterWeightResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/latest_forecaster_weight/{topic_id}/{actor_id}",
		[]string{"topic_id", "actor_id"}, nil,
		req, resp,
	)
	return resp, errors.Wrap(err, "while calling EmissionsRESTClient.GetLatestForecasterWeight")
}
