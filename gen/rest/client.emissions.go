package rest

import (
	"context"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"

	emissionstypes "github.com/allora-network/allora-chain/x/emissions/types"

	"github.com/allora-network/allora-sdk-go/config"
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

func (c *EmissionsRESTClient) GetParams(ctx context.Context, req *emissionstypes.GetParamsRequest, opts ...config.CallOpt) (*emissionstypes.GetParamsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetParamsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/params",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetParams")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetNextTopicId(ctx context.Context, req *emissionstypes.GetNextTopicIdRequest, opts ...config.CallOpt) (*emissionstypes.GetNextTopicIdResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetNextTopicIdResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/next_topic_id",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetNextTopicId")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetTopic(ctx context.Context, req *emissionstypes.GetTopicRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetTopicResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/topics/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetTopic")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetWorkerLatestInferenceByTopicId(ctx context.Context, req *emissionstypes.GetWorkerLatestInferenceByTopicIdRequest, opts ...config.CallOpt) (*emissionstypes.GetWorkerLatestInferenceByTopicIdResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetWorkerLatestInferenceByTopicIdResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/topics/{topic_id}/workers/{worker_address}/latest_inference",
		[]string{"topic_id", "worker_address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetWorkerLatestInferenceByTopicId")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetInferencesAtBlock(ctx context.Context, req *emissionstypes.GetInferencesAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetInferencesAtBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetInferencesAtBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/inferences/{topic_id}/{block_height}",
		[]string{"topic_id", "block_height"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetInferencesAtBlock")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetLatestTopicInferences(ctx context.Context, req *emissionstypes.GetLatestTopicInferencesRequest, opts ...config.CallOpt) (*emissionstypes.GetLatestTopicInferencesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetLatestTopicInferencesResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/latest_inferences/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetLatestTopicInferences")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetForecastsAtBlock(ctx context.Context, req *emissionstypes.GetForecastsAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetForecastsAtBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetForecastsAtBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/forecasts/{topic_id}/{block_height}",
		[]string{"topic_id", "block_height"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetForecastsAtBlock")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetNetworkLossBundleAtBlock(ctx context.Context, req *emissionstypes.GetNetworkLossBundleAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetNetworkLossBundleAtBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetNetworkLossBundleAtBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/network_loss/{topic_id}/{block_height}",
		[]string{"topic_id", "block_height"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetNetworkLossBundleAtBlock")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetTotalStake(ctx context.Context, req *emissionstypes.GetTotalStakeRequest, opts ...config.CallOpt) (*emissionstypes.GetTotalStakeResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetTotalStakeResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/total_stake",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetTotalStake")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetReputerStakeInTopic(ctx context.Context, req *emissionstypes.GetReputerStakeInTopicRequest, opts ...config.CallOpt) (*emissionstypes.GetReputerStakeInTopicResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetReputerStakeInTopicResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/reputer_stake/{address}/{topic_id}",
		[]string{"address", "topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetReputerStakeInTopic")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetMultiReputerStakeInTopic(ctx context.Context, req *emissionstypes.GetMultiReputerStakeInTopicRequest, opts ...config.CallOpt) (*emissionstypes.GetMultiReputerStakeInTopicResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetMultiReputerStakeInTopicResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/reputers_stakes/{topic_id}",
		[]string{"topic_id"}, []string{"addresses"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetMultiReputerStakeInTopic")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetStakeFromReputerInTopicInSelf(ctx context.Context, req *emissionstypes.GetStakeFromReputerInTopicInSelfRequest, opts ...config.CallOpt) (*emissionstypes.GetStakeFromReputerInTopicInSelfResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetStakeFromReputerInTopicInSelfResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/reputer_stake_self/{reputer_address}/{topic_id}",
		[]string{"reputer_address", "topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetStakeFromReputerInTopicInSelf")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetDelegateStakeInTopicInReputer(ctx context.Context, req *emissionstypes.GetDelegateStakeInTopicInReputerRequest, opts ...config.CallOpt) (*emissionstypes.GetDelegateStakeInTopicInReputerResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetDelegateStakeInTopicInReputerResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/reputer_delegate_stake/{reputer_address}/{topic_id}",
		[]string{"reputer_address", "topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetDelegateStakeInTopicInReputer")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetStakeFromDelegatorInTopicInReputer(ctx context.Context, req *emissionstypes.GetStakeFromDelegatorInTopicInReputerRequest, opts ...config.CallOpt) (*emissionstypes.GetStakeFromDelegatorInTopicInReputerResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetStakeFromDelegatorInTopicInReputerResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/delegate_stake/{delegator_address}/{reputer_address}/{topic_id}",
		[]string{"delegator_address", "reputer_address", "topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetStakeFromDelegatorInTopicInReputer")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetStakeFromDelegatorInTopic(ctx context.Context, req *emissionstypes.GetStakeFromDelegatorInTopicRequest, opts ...config.CallOpt) (*emissionstypes.GetStakeFromDelegatorInTopicResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetStakeFromDelegatorInTopicResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/delegate_stake/{delegator_address}/{topic_id}",
		[]string{"delegator_address", "topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetStakeFromDelegatorInTopic")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetTopicStake(ctx context.Context, req *emissionstypes.GetTopicStakeRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicStakeResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetTopicStakeResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/stake/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetTopicStake")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetStakeRemovalsUpUntilBlock(ctx context.Context, req *emissionstypes.GetStakeRemovalsUpUntilBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetStakeRemovalsUpUntilBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetStakeRemovalsUpUntilBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/stake_removals/{block_height}",
		[]string{"block_height"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetStakeRemovalsUpUntilBlock")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetDelegateStakeRemovalsUpUntilBlock(ctx context.Context, req *emissionstypes.GetDelegateStakeRemovalsUpUntilBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetDelegateStakeRemovalsUpUntilBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetDelegateStakeRemovalsUpUntilBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/delegate_stake_removals/{block_height}",
		[]string{"block_height"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetDelegateStakeRemovalsUpUntilBlock")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetStakeRemovalInfo(ctx context.Context, req *emissionstypes.GetStakeRemovalInfoRequest, opts ...config.CallOpt) (*emissionstypes.GetStakeRemovalInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetStakeRemovalInfoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/stake_removal/{topic_id}/{reputer}",
		[]string{"topic_id", "reputer"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetStakeRemovalInfo")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetDelegateStakeRemovalInfo(ctx context.Context, req *emissionstypes.GetDelegateStakeRemovalInfoRequest, opts ...config.CallOpt) (*emissionstypes.GetDelegateStakeRemovalInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetDelegateStakeRemovalInfoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/delegate_stake_removal/{topic_id}/{delegator}/{reputer}",
		[]string{"topic_id", "delegator", "reputer"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetDelegateStakeRemovalInfo")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetWorkerNodeInfo(ctx context.Context, req *emissionstypes.GetWorkerNodeInfoRequest, opts ...config.CallOpt) (*emissionstypes.GetWorkerNodeInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetWorkerNodeInfoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/worker/{address}",
		[]string{"address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetWorkerNodeInfo")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetReputerNodeInfo(ctx context.Context, req *emissionstypes.GetReputerNodeInfoRequest, opts ...config.CallOpt) (*emissionstypes.GetReputerNodeInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetReputerNodeInfoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/reputer/{address}",
		[]string{"address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetReputerNodeInfo")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) IsWorkerRegisteredInTopicId(ctx context.Context, req *emissionstypes.IsWorkerRegisteredInTopicIdRequest, opts ...config.CallOpt) (*emissionstypes.IsWorkerRegisteredInTopicIdResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.IsWorkerRegisteredInTopicIdResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/worker_registered/{topic_id}/{address}",
		[]string{"topic_id", "address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.IsWorkerRegisteredInTopicId")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) IsReputerRegisteredInTopicId(ctx context.Context, req *emissionstypes.IsReputerRegisteredInTopicIdRequest, opts ...config.CallOpt) (*emissionstypes.IsReputerRegisteredInTopicIdResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.IsReputerRegisteredInTopicIdResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/reputer_registered/{topic_id}/{address}",
		[]string{"topic_id", "address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.IsReputerRegisteredInTopicId")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetNetworkInferencesAtBlock(ctx context.Context, req *emissionstypes.GetNetworkInferencesAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetNetworkInferencesAtBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetNetworkInferencesAtBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/network_inferences/{topic_id}/last_inference/{block_height_last_inference}",
		[]string{"topic_id", "block_height_last_inference"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetNetworkInferencesAtBlock")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetNetworkInferencesAtBlockOutlierResistant(ctx context.Context, req *emissionstypes.GetNetworkInferencesAtBlockOutlierResistantRequest, opts ...config.CallOpt) (*emissionstypes.GetNetworkInferencesAtBlockOutlierResistantResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetNetworkInferencesAtBlockOutlierResistantResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/network_inferences_outlier_resistant/{topic_id}/last_inference/{block_height_last_inference}",
		[]string{"topic_id", "block_height_last_inference"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetNetworkInferencesAtBlockOutlierResistant")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetLatestNetworkInferences(ctx context.Context, req *emissionstypes.GetLatestNetworkInferencesRequest, opts ...config.CallOpt) (*emissionstypes.GetLatestNetworkInferencesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetLatestNetworkInferencesResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/latest_network_inferences/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetLatestNetworkInferences")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetLatestNetworkInferencesOutlierResistant(ctx context.Context, req *emissionstypes.GetLatestNetworkInferencesOutlierResistantRequest, opts ...config.CallOpt) (*emissionstypes.GetLatestNetworkInferencesOutlierResistantResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetLatestNetworkInferencesOutlierResistantResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/latest_network_inferences_outlier_resistant/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetLatestNetworkInferencesOutlierResistant")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) IsWorkerNonceUnfulfilled(ctx context.Context, req *emissionstypes.IsWorkerNonceUnfulfilledRequest, opts ...config.CallOpt) (*emissionstypes.IsWorkerNonceUnfulfilledResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.IsWorkerNonceUnfulfilledResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_worker_nonce_unfulfilled/{topic_id}/{block_height}",
		[]string{"topic_id", "block_height"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.IsWorkerNonceUnfulfilled")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) IsReputerNonceUnfulfilled(ctx context.Context, req *emissionstypes.IsReputerNonceUnfulfilledRequest, opts ...config.CallOpt) (*emissionstypes.IsReputerNonceUnfulfilledResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.IsReputerNonceUnfulfilledResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_reputer_nonce_unfulfilled/{topic_id}/{block_height}",
		[]string{"topic_id", "block_height"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.IsReputerNonceUnfulfilled")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetUnfulfilledWorkerNonces(ctx context.Context, req *emissionstypes.GetUnfulfilledWorkerNoncesRequest, opts ...config.CallOpt) (*emissionstypes.GetUnfulfilledWorkerNoncesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetUnfulfilledWorkerNoncesResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/unfulfilled_worker_nonces/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetUnfulfilledWorkerNonces")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetUnfulfilledReputerNonces(ctx context.Context, req *emissionstypes.GetUnfulfilledReputerNoncesRequest, opts ...config.CallOpt) (*emissionstypes.GetUnfulfilledReputerNoncesResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetUnfulfilledReputerNoncesResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/unfulfilled_reputer_nonces/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetUnfulfilledReputerNonces")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetInfererNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetInfererNetworkRegretResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetInfererNetworkRegretResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/inferer_network_regret/{topic_id}/{actor_id}",
		[]string{"topic_id", "actor_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetInfererNetworkRegret")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetForecasterNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetForecasterNetworkRegretResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetForecasterNetworkRegretResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/forecaster_network_regret/{topic_id}/{worker}",
		[]string{"topic_id", "worker"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetForecasterNetworkRegret")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetOneInForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetOneInForecasterNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetOneInForecasterNetworkRegretResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetOneInForecasterNetworkRegretResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/one_in_forecaster_network_regret/{topic_id}/{forecaster}/{inferer}",
		[]string{"topic_id", "forecaster", "inferer"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetOneInForecasterNetworkRegret")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) IsWhitelistAdmin(ctx context.Context, req *emissionstypes.IsWhitelistAdminRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistAdminResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.IsWhitelistAdminResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/whitelist_admin/{address}",
		[]string{"address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.IsWhitelistAdmin")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetTopicLastWorkerCommitInfo(ctx context.Context, req *emissionstypes.GetTopicLastWorkerCommitInfoRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicLastWorkerCommitInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetTopicLastWorkerCommitInfoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/topic_last_worker_commit_info/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetTopicLastWorkerCommitInfo")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetTopicLastReputerCommitInfo(ctx context.Context, req *emissionstypes.GetTopicLastReputerCommitInfoRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicLastReputerCommitInfoResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetTopicLastReputerCommitInfoResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/topic_last_reputer_commit_info/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetTopicLastReputerCommitInfo")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetTopicRewardNonce(ctx context.Context, req *emissionstypes.GetTopicRewardNonceRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicRewardNonceResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetTopicRewardNonceResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/topic_reward_nonce/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetTopicRewardNonce")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetReputerLossBundlesAtBlock(ctx context.Context, req *emissionstypes.GetReputerLossBundlesAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetReputerLossBundlesAtBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetReputerLossBundlesAtBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/reputer_loss_bundles/{topic_id}/{block_height}",
		[]string{"topic_id", "block_height"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetReputerLossBundlesAtBlock")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetStakeReputerAuthority(ctx context.Context, req *emissionstypes.GetStakeReputerAuthorityRequest, opts ...config.CallOpt) (*emissionstypes.GetStakeReputerAuthorityResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetStakeReputerAuthorityResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/stake_reputer_authority/{topic_id}/{reputer}",
		[]string{"topic_id", "reputer"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetStakeReputerAuthority")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetDelegateStakePlacement(ctx context.Context, req *emissionstypes.GetDelegateStakePlacementRequest, opts ...config.CallOpt) (*emissionstypes.GetDelegateStakePlacementResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetDelegateStakePlacementResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/delegate_stake_placement/{topic_id}/{delegator}/{target}",
		[]string{"topic_id", "delegator", "target"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetDelegateStakePlacement")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetDelegateStakeUponReputer(ctx context.Context, req *emissionstypes.GetDelegateStakeUponReputerRequest, opts ...config.CallOpt) (*emissionstypes.GetDelegateStakeUponReputerResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetDelegateStakeUponReputerResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/delegate_stake_upon_reputer/{topic_id}/{target}",
		[]string{"topic_id", "target"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetDelegateStakeUponReputer")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetDelegateRewardPerShare(ctx context.Context, req *emissionstypes.GetDelegateRewardPerShareRequest, opts ...config.CallOpt) (*emissionstypes.GetDelegateRewardPerShareResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetDelegateRewardPerShareResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/delegate_reward_per_share/{topic_id}/{reputer}",
		[]string{"topic_id", "reputer"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetDelegateRewardPerShare")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetStakeRemovalForReputerAndTopicId(ctx context.Context, req *emissionstypes.GetStakeRemovalForReputerAndTopicIdRequest, opts ...config.CallOpt) (*emissionstypes.GetStakeRemovalForReputerAndTopicIdResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetStakeRemovalForReputerAndTopicIdResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/stake_removal/{reputer}/{topic_id}",
		[]string{"reputer", "topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetStakeRemovalForReputerAndTopicId")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetDelegateStakeRemoval(ctx context.Context, req *emissionstypes.GetDelegateStakeRemovalRequest, opts ...config.CallOpt) (*emissionstypes.GetDelegateStakeRemovalResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetDelegateStakeRemovalResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/delegate_stake_removal/{block_height}/{topic_id}/{delegator}/{reputer}",
		[]string{"block_height", "topic_id", "delegator", "reputer"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetDelegateStakeRemoval")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetPreviousTopicWeight(ctx context.Context, req *emissionstypes.GetPreviousTopicWeightRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousTopicWeightResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetPreviousTopicWeightResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/previous_topic_weight/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetPreviousTopicWeight")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetTotalSumPreviousTopicWeights(ctx context.Context, req *emissionstypes.GetTotalSumPreviousTopicWeightsRequest, opts ...config.CallOpt) (*emissionstypes.GetTotalSumPreviousTopicWeightsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetTotalSumPreviousTopicWeightsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/sum_previous_total_topic_weight",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetTotalSumPreviousTopicWeights")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) TopicExists(ctx context.Context, req *emissionstypes.TopicExistsRequest, opts ...config.CallOpt) (*emissionstypes.TopicExistsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.TopicExistsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/topic_exists/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.TopicExists")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) IsTopicActive(ctx context.Context, req *emissionstypes.IsTopicActiveRequest, opts ...config.CallOpt) (*emissionstypes.IsTopicActiveResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.IsTopicActiveResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_topic_active/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.IsTopicActive")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetTopicFeeRevenue(ctx context.Context, req *emissionstypes.GetTopicFeeRevenueRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicFeeRevenueResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetTopicFeeRevenueResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/topic_fee_revenue/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetTopicFeeRevenue")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetInfererScoreEma(ctx context.Context, req *emissionstypes.GetInfererScoreEmaRequest, opts ...config.CallOpt) (*emissionstypes.GetInfererScoreEmaResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetInfererScoreEmaResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/inferer_score_ema/{topic_id}/{inferer}",
		[]string{"topic_id", "inferer"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetInfererScoreEma")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetForecasterScoreEma(ctx context.Context, req *emissionstypes.GetForecasterScoreEmaRequest, opts ...config.CallOpt) (*emissionstypes.GetForecasterScoreEmaResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetForecasterScoreEmaResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/forecaster_score_ema/{topic_id}/{forecaster}",
		[]string{"topic_id", "forecaster"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetForecasterScoreEma")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetReputerScoreEma(ctx context.Context, req *emissionstypes.GetReputerScoreEmaRequest, opts ...config.CallOpt) (*emissionstypes.GetReputerScoreEmaResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetReputerScoreEmaResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/reputer_score_ema/{topic_id}/{reputer}",
		[]string{"topic_id", "reputer"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetReputerScoreEma")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetInferenceScoresUntilBlock(ctx context.Context, req *emissionstypes.GetInferenceScoresUntilBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetInferenceScoresUntilBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetInferenceScoresUntilBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/inference_scores_until_block/{topic_id}/{block_height}",
		[]string{"topic_id", "block_height"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetInferenceScoresUntilBlock")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetPreviousTopicQuantileForecasterScoreEma(ctx context.Context, req *emissionstypes.GetPreviousTopicQuantileForecasterScoreEmaRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousTopicQuantileForecasterScoreEmaResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetPreviousTopicQuantileForecasterScoreEmaResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/topic_quantile_forecaster_score_ema/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetPreviousTopicQuantileForecasterScoreEma")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetPreviousTopicQuantileInfererScoreEma(ctx context.Context, req *emissionstypes.GetPreviousTopicQuantileInfererScoreEmaRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousTopicQuantileInfererScoreEmaResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetPreviousTopicQuantileInfererScoreEmaResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/topic_quantile_inferer_score_ema/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetPreviousTopicQuantileInfererScoreEma")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetPreviousTopicQuantileReputerScoreEma(ctx context.Context, req *emissionstypes.GetPreviousTopicQuantileReputerScoreEmaRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousTopicQuantileReputerScoreEmaResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetPreviousTopicQuantileReputerScoreEmaResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/topic_quantile_reputer_score_ema/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetPreviousTopicQuantileReputerScoreEma")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetWorkerInferenceScoresAtBlock(ctx context.Context, req *emissionstypes.GetWorkerInferenceScoresAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetWorkerInferenceScoresAtBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetWorkerInferenceScoresAtBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/worker_inference_scores_at_block/{topic_id}/{block_height}",
		[]string{"topic_id", "block_height"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetWorkerInferenceScoresAtBlock")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetCurrentLowestInfererScore(ctx context.Context, req *emissionstypes.GetCurrentLowestInfererScoreRequest, opts ...config.CallOpt) (*emissionstypes.GetCurrentLowestInfererScoreResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetCurrentLowestInfererScoreResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/current_lowest_inferer_score/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetCurrentLowestInfererScore")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetForecastScoresUntilBlock(ctx context.Context, req *emissionstypes.GetForecastScoresUntilBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetForecastScoresUntilBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetForecastScoresUntilBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/forecast_scores_until_block/{topic_id}/{block_height}",
		[]string{"topic_id", "block_height"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetForecastScoresUntilBlock")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetWorkerForecastScoresAtBlock(ctx context.Context, req *emissionstypes.GetWorkerForecastScoresAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetWorkerForecastScoresAtBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetWorkerForecastScoresAtBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/worker_forecast_scores_at_block/{topic_id}/{block_height}",
		[]string{"topic_id", "block_height"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetWorkerForecastScoresAtBlock")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetCurrentLowestForecasterScore(ctx context.Context, req *emissionstypes.GetCurrentLowestForecasterScoreRequest, opts ...config.CallOpt) (*emissionstypes.GetCurrentLowestForecasterScoreResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetCurrentLowestForecasterScoreResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/current_lowest_forecaster_score/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetCurrentLowestForecasterScore")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetReputersScoresAtBlock(ctx context.Context, req *emissionstypes.GetReputersScoresAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetReputersScoresAtBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetReputersScoresAtBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/reputers_scores_at_block/{topic_id}/{block_height}",
		[]string{"topic_id", "block_height"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetReputersScoresAtBlock")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetCurrentLowestReputerScore(ctx context.Context, req *emissionstypes.GetCurrentLowestReputerScoreRequest, opts ...config.CallOpt) (*emissionstypes.GetCurrentLowestReputerScoreResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetCurrentLowestReputerScoreResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/current_lowest_reputer_score/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetCurrentLowestReputerScore")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetListeningCoefficient(ctx context.Context, req *emissionstypes.GetListeningCoefficientRequest, opts ...config.CallOpt) (*emissionstypes.GetListeningCoefficientResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetListeningCoefficientResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/listening_coefficient/{topic_id}/{reputer}",
		[]string{"topic_id", "reputer"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetListeningCoefficient")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetPreviousReputerRewardFraction(ctx context.Context, req *emissionstypes.GetPreviousReputerRewardFractionRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousReputerRewardFractionResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetPreviousReputerRewardFractionResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/previous_reputer_reward_fraction/{topic_id}/{reputer}",
		[]string{"topic_id", "reputer"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetPreviousReputerRewardFraction")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetPreviousInferenceRewardFraction(ctx context.Context, req *emissionstypes.GetPreviousInferenceRewardFractionRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousInferenceRewardFractionResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetPreviousInferenceRewardFractionResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/previous_inference_reward_fraction/{topic_id}/{worker}",
		[]string{"topic_id", "worker"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetPreviousInferenceRewardFraction")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetPreviousForecastRewardFraction(ctx context.Context, req *emissionstypes.GetPreviousForecastRewardFractionRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousForecastRewardFractionResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetPreviousForecastRewardFractionResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/previous_forecast_reward_fraction/{topic_id}/{worker}",
		[]string{"topic_id", "worker"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetPreviousForecastRewardFraction")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetPreviousPercentageRewardToStakedReputers(ctx context.Context, req *emissionstypes.GetPreviousPercentageRewardToStakedReputersRequest, opts ...config.CallOpt) (*emissionstypes.GetPreviousPercentageRewardToStakedReputersResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetPreviousPercentageRewardToStakedReputersResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/previous_percentage_reward_to_staked_reputers",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetPreviousPercentageRewardToStakedReputers")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetTotalRewardToDistribute(ctx context.Context, req *emissionstypes.GetTotalRewardToDistributeRequest, opts ...config.CallOpt) (*emissionstypes.GetTotalRewardToDistributeResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetTotalRewardToDistributeResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/total_reward_to_distribute",
		nil, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetTotalRewardToDistribute")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetNaiveInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetNaiveInfererNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetNaiveInfererNetworkRegretResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetNaiveInfererNetworkRegretResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/native_inferer_network_regret",
		nil, []string{"topic_id", "inferer"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetNaiveInfererNetworkRegret")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetOneOutInfererInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutInfererInfererNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetOneOutInfererInfererNetworkRegretResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetOneOutInfererInfererNetworkRegretResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/one_out_inferer_inferer_network_regret",
		nil, []string{"topic_id", "one_out_inferer", "inferer"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetOneOutInfererInfererNetworkRegret")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetOneOutInfererForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutInfererForecasterNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetOneOutInfererForecasterNetworkRegretResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetOneOutInfererForecasterNetworkRegretResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/one_out_inferer_forecaster_network_regret",
		nil, []string{"topic_id", "one_out_inferer", "forecaster"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetOneOutInfererForecasterNetworkRegret")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetOneOutForecasterInfererNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutForecasterInfererNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetOneOutForecasterInfererNetworkRegretResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetOneOutForecasterInfererNetworkRegretResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/one_out_forecaster_inferer_network_regret",
		nil, []string{"topic_id", "one_out_forecaster", "inferer"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetOneOutForecasterInfererNetworkRegret")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetOneOutForecasterForecasterNetworkRegret(ctx context.Context, req *emissionstypes.GetOneOutForecasterForecasterNetworkRegretRequest, opts ...config.CallOpt) (*emissionstypes.GetOneOutForecasterForecasterNetworkRegretResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetOneOutForecasterForecasterNetworkRegretResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/one_out_forecaster_forecaster_network_regret",
		nil, []string{"topic_id", "one_out_forecaster", "forecaster"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetOneOutForecasterForecasterNetworkRegret")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetActiveTopicsAtBlock(ctx context.Context, req *emissionstypes.GetActiveTopicsAtBlockRequest, opts ...config.CallOpt) (*emissionstypes.GetActiveTopicsAtBlockResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetActiveTopicsAtBlockResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/active_topics_at_block/{block_height}",
		[]string{"block_height"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetActiveTopicsAtBlock")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetNextChurningBlockByTopicId(ctx context.Context, req *emissionstypes.GetNextChurningBlockByTopicIdRequest, opts ...config.CallOpt) (*emissionstypes.GetNextChurningBlockByTopicIdResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetNextChurningBlockByTopicIdResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/next_churning_block_by_topic_id/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetNextChurningBlockByTopicId")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetCountInfererInclusionsInTopic(ctx context.Context, req *emissionstypes.GetCountInfererInclusionsInTopicRequest, opts ...config.CallOpt) (*emissionstypes.GetCountInfererInclusionsInTopicResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetCountInfererInclusionsInTopicResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/count_inferer_inclusions_in_topic/{topic_id}/{inferer}",
		[]string{"topic_id", "inferer"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetCountInfererInclusionsInTopic")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetCountForecasterInclusionsInTopic(ctx context.Context, req *emissionstypes.GetCountForecasterInclusionsInTopicRequest, opts ...config.CallOpt) (*emissionstypes.GetCountForecasterInclusionsInTopicResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetCountForecasterInclusionsInTopicResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/count_forecaster_inclusions_in_topic/{topic_id}/{forecaster}",
		[]string{"topic_id", "forecaster"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetCountForecasterInclusionsInTopic")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) IsWhitelistedGlobalWorker(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalWorkerRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistedGlobalWorkerResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.IsWhitelistedGlobalWorkerResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_whitelisted_global_worker/{address}",
		[]string{"address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.IsWhitelistedGlobalWorker")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) IsWhitelistedGlobalReputer(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalReputerRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistedGlobalReputerResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.IsWhitelistedGlobalReputerResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_whitelisted_global_reputer/{address}",
		[]string{"address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.IsWhitelistedGlobalReputer")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) IsWhitelistedGlobalAdmin(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalAdminRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistedGlobalAdminResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.IsWhitelistedGlobalAdminResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_whitelisted_global_admin/{address}",
		[]string{"address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.IsWhitelistedGlobalAdmin")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) IsTopicWorkerWhitelistEnabled(ctx context.Context, req *emissionstypes.IsTopicWorkerWhitelistEnabledRequest, opts ...config.CallOpt) (*emissionstypes.IsTopicWorkerWhitelistEnabledResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.IsTopicWorkerWhitelistEnabledResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_topic_worker_whitelist_enabled/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.IsTopicWorkerWhitelistEnabled")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) IsTopicReputerWhitelistEnabled(ctx context.Context, req *emissionstypes.IsTopicReputerWhitelistEnabledRequest, opts ...config.CallOpt) (*emissionstypes.IsTopicReputerWhitelistEnabledResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.IsTopicReputerWhitelistEnabledResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_topic_reputer_whitelist_enabled/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.IsTopicReputerWhitelistEnabled")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) IsWhitelistedTopicCreator(ctx context.Context, req *emissionstypes.IsWhitelistedTopicCreatorRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistedTopicCreatorResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.IsWhitelistedTopicCreatorResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_whitelisted_topic_creator/{address}",
		[]string{"address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.IsWhitelistedTopicCreator")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) IsWhitelistedGlobalActor(ctx context.Context, req *emissionstypes.IsWhitelistedGlobalActorRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistedGlobalActorResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.IsWhitelistedGlobalActorResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_whitelisted_global_actor/{address}",
		[]string{"address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.IsWhitelistedGlobalActor")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) IsWhitelistedTopicWorker(ctx context.Context, req *emissionstypes.IsWhitelistedTopicWorkerRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistedTopicWorkerResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.IsWhitelistedTopicWorkerResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_whitelisted_topic_worker/{topic_id}/{address}",
		[]string{"topic_id", "address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.IsWhitelistedTopicWorker")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) IsWhitelistedTopicReputer(ctx context.Context, req *emissionstypes.IsWhitelistedTopicReputerRequest, opts ...config.CallOpt) (*emissionstypes.IsWhitelistedTopicReputerResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.IsWhitelistedTopicReputerResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/is_whitelisted_topic_reputer/{topic_id}/{address}",
		[]string{"topic_id", "address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.IsWhitelistedTopicReputer")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) CanUpdateAllGlobalWhitelists(ctx context.Context, req *emissionstypes.CanUpdateAllGlobalWhitelistsRequest, opts ...config.CallOpt) (*emissionstypes.CanUpdateAllGlobalWhitelistsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.CanUpdateAllGlobalWhitelistsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/can_update_all_global_whitelists/{address}",
		[]string{"address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.CanUpdateAllGlobalWhitelists")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) CanUpdateGlobalWorkerWhitelist(ctx context.Context, req *emissionstypes.CanUpdateGlobalWorkerWhitelistRequest, opts ...config.CallOpt) (*emissionstypes.CanUpdateGlobalWorkerWhitelistResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.CanUpdateGlobalWorkerWhitelistResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/can_update_global_worker_whitelist/{address}",
		[]string{"address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.CanUpdateGlobalWorkerWhitelist")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) CanUpdateGlobalReputerWhitelist(ctx context.Context, req *emissionstypes.CanUpdateGlobalReputerWhitelistRequest, opts ...config.CallOpt) (*emissionstypes.CanUpdateGlobalReputerWhitelistResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.CanUpdateGlobalReputerWhitelistResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/can_update_global_reputer_whitelist/{address}",
		[]string{"address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.CanUpdateGlobalReputerWhitelist")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) CanUpdateParams(ctx context.Context, req *emissionstypes.CanUpdateParamsRequest, opts ...config.CallOpt) (*emissionstypes.CanUpdateParamsResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.CanUpdateParamsResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/can_update_params/{address}",
		[]string{"address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.CanUpdateParams")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) CanUpdateTopicWhitelist(ctx context.Context, req *emissionstypes.CanUpdateTopicWhitelistRequest, opts ...config.CallOpt) (*emissionstypes.CanUpdateTopicWhitelistResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.CanUpdateTopicWhitelistResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/can_update_topic_whitelist/{topic_id}/{address}",
		[]string{"topic_id", "address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.CanUpdateTopicWhitelist")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) CanCreateTopic(ctx context.Context, req *emissionstypes.CanCreateTopicRequest, opts ...config.CallOpt) (*emissionstypes.CanCreateTopicResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.CanCreateTopicResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/can_create_topic/{address}",
		[]string{"address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.CanCreateTopic")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) CanSubmitWorkerPayload(ctx context.Context, req *emissionstypes.CanSubmitWorkerPayloadRequest, opts ...config.CallOpt) (*emissionstypes.CanSubmitWorkerPayloadResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.CanSubmitWorkerPayloadResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/can_submit_worker_payload/{topic_id}/{address}",
		[]string{"topic_id", "address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.CanSubmitWorkerPayload")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) CanSubmitReputerPayload(ctx context.Context, req *emissionstypes.CanSubmitReputerPayloadRequest, opts ...config.CallOpt) (*emissionstypes.CanSubmitReputerPayloadResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.CanSubmitReputerPayloadResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/can_submit_reputer_payload/{topic_id}/{address}",
		[]string{"topic_id", "address"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.CanSubmitReputerPayload")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetTopicInitialInfererEmaScore(ctx context.Context, req *emissionstypes.GetTopicInitialInfererEmaScoreRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicInitialInfererEmaScoreResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetTopicInitialInfererEmaScoreResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/initial_inferer_ema_score/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetTopicInitialInfererEmaScore")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetTopicInitialForecasterEmaScore(ctx context.Context, req *emissionstypes.GetTopicInitialForecasterEmaScoreRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicInitialForecasterEmaScoreResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetTopicInitialForecasterEmaScoreResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/initial_forecaster_ema_score/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetTopicInitialForecasterEmaScore")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetTopicInitialReputerEmaScore(ctx context.Context, req *emissionstypes.GetTopicInitialReputerEmaScoreRequest, opts ...config.CallOpt) (*emissionstypes.GetTopicInitialReputerEmaScoreResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetTopicInitialReputerEmaScoreResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/initial_reputer_ema_score/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetTopicInitialReputerEmaScore")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetLatestRegretStdNorm(ctx context.Context, req *emissionstypes.GetLatestRegretStdNormRequest, opts ...config.CallOpt) (*emissionstypes.GetLatestRegretStdNormResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetLatestRegretStdNormResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/latest_regret_stdnorm/{topic_id}",
		[]string{"topic_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetLatestRegretStdNorm")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetLatestInfererWeight(ctx context.Context, req *emissionstypes.GetLatestInfererWeightRequest, opts ...config.CallOpt) (*emissionstypes.GetLatestInfererWeightResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetLatestInfererWeightResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/latest_inferer_weight/{topic_id}/{actor_id}",
		[]string{"topic_id", "actor_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetLatestInfererWeight")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetLatestForecasterWeight(ctx context.Context, req *emissionstypes.GetLatestForecasterWeightRequest, opts ...config.CallOpt) (*emissionstypes.GetLatestForecasterWeightResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetLatestForecasterWeightResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/latest_forecaster_weight/{topic_id}/{actor_id}",
		[]string{"topic_id", "actor_id"}, nil,
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetLatestForecasterWeight")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetWorkerSubmissionWindowStatus(ctx context.Context, req *emissionstypes.GetWorkerSubmissionWindowStatusRequest, opts ...config.CallOpt) (*emissionstypes.GetWorkerSubmissionWindowStatusResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetWorkerSubmissionWindowStatusResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/worker_submission_window_status/{topic_id}",
		[]string{"topic_id"}, []string{"address"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetWorkerSubmissionWindowStatus")
	}
	return resp, nil
}

func (c *EmissionsRESTClient) GetReputerSubmissionWindowStatus(ctx context.Context, req *emissionstypes.GetReputerSubmissionWindowStatusRequest, opts ...config.CallOpt) (*emissionstypes.GetReputerSubmissionWindowStatusResponse, error) {
	callOpts := config.DefaultCallOpts()
	callOpts.Apply(opts...)

	resp := &emissionstypes.GetReputerSubmissionWindowStatusResponse{}
	err := c.RESTClientCore.executeRequest(ctx,
		"GET", "/emissions/v9/reputer_submission_window_status/{topic_id}",
		[]string{"topic_id"}, []string{"address"},
		req, resp, callOpts.Height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "while calling EmissionsRESTClient.GetReputerSubmissionWindowStatus")
	}
	return resp, nil
}
