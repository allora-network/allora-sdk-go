package codec_test

import (
	"context"

	cosmosmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	alloramath "github.com/allora-network/allora-chain/math"
	emissionsv10 "github.com/allora-network/allora-chain/x/emissions/types"
)

func emittedEventFixtures() []emittedEventFixture {
	dec := alloramath.MustNewDecFromString
	infNonce := &emissionsv10.Nonce{BlockHeight: 103}
	repNonce := &emissionsv10.ReputerRequestNonce{ReputerNonce: &emissionsv10.Nonce{BlockHeight: 105}}
	inputInference := &emissionsv10.InputWorkerDataBundle{
		Worker: "inferer-payload", TopicId: 19, Nonce: infNonce,
		InferenceForecastsBundle: &emissionsv10.InputInferenceForecastBundle{Inference: &emissionsv10.InputInference{
			Value: alloramath.MustNewBoundedExp40DecFromString("31.1234"), Values: []*emissionsv10.InputLabeledValue{{Label: "label-a", Value: alloramath.MustNewBoundedExp40DecFromString("41.2345")}}, ExtraData: []byte{1, 2, 127, 128},
		}},
	}
	inputForecast := &emissionsv10.WorkerDataBundle{
		Worker: "forecaster-payload", TopicId: 20, Nonce: &emissionsv10.Nonce{BlockHeight: 104},
		InferenceForecastsBundle: &emissionsv10.InferenceForecastBundle{Forecast: &emissionsv10.Forecast{
			Forecaster: "forecaster-payload", ForecastElements: []*emissionsv10.ForecastElement{
				{Inferer: "forecast-inferer-a", Value: dec("51.000000000000000001")},
				{Inferer: "forecast-inferer-b", Value: dec("52.000000000000000002")},
			}, ExtraData: []byte{3, 4, 129, 130},
		}},
	}
	lossBundle := &emissionsv10.LossBundle{
		TopicId: 21, Reputer: "reputer-payload", ReputerRequestNonce: repNonce, ExtraData: []byte{5, 6, 131, 132},
		CombinedValue: dec("61.123456789012345678"), NaiveValue: dec("62.123456789012345678"),
		InfererValues:                 []*emissionsv10.WorkerAttributedValue{{Worker: "loss-inferer-a", Value: dec("63.123456789012345678")}},
		ForecasterValues:              []*emissionsv10.WorkerAttributedValue{{Worker: "loss-forecaster-a", Value: dec("64.123456789012345678")}},
		OneOutInfererValues:           []*emissionsv10.WithheldWorkerAttributedValue{{Worker: "loss-inferer-a", Value: dec("65.123456789012345678")}},
		OneOutForecasterValues:        []*emissionsv10.WithheldWorkerAttributedValue{{Worker: "loss-forecaster-a", Value: dec("66.123456789012345678")}},
		OneInForecasterValues:         []*emissionsv10.WorkerAttributedValue{{Worker: "loss-forecaster-a", Value: dec("67.123456789012345678")}},
		OneOutInfererForecasterValues: []*emissionsv10.OneOutInfererForecasterValues{{Forecaster: "loss-forecaster-a", OneOutInfererValues: []*emissionsv10.WithheldWorkerAttributedValue{{Worker: "loss-inferer-a", Value: dec("68.123456789012345678")}}}},
	}
	eventValueBundle := emissionsv10.ValueBundleToEventValueBundleBase(lossBundle)

	return []emittedEventFixture{
		{
			name: "EventInsertInfererPayload",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewInsertInfererPayloadEvent(ctx, inputInference, emissionsv10.PayloadAdmissionOutcome_PAYLOAD_ADMISSION_OUTCOME_NOT_ADMITTED)
			},
			expected: &emissionsv10.EventInsertInfererPayload{
				Inferer: "inferer-payload", Nonce: 103, TopicId: 19, Value: alloramath.MustNewBoundedExp40DecFromString("31.1234"), ExtraData: []byte{1, 2, 127, 128},
				Values: []*emissionsv10.InputLabeledValue{{Label: "label-a", Value: alloramath.MustNewBoundedExp40DecFromString("41.2345")}}, AdmissionOutcome: emissionsv10.PayloadAdmissionOutcome_PAYLOAD_ADMISSION_OUTCOME_NOT_ADMITTED,
			},
		},
		{
			name: "EventInsertForecasterPayload",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewInsertForecasterPayloadEvent(ctx, inputForecast, emissionsv10.PayloadAdmissionOutcome_PAYLOAD_ADMISSION_OUTCOME_NO_ELEMENTS_ACCEPTED)
			},
			expected: &emissionsv10.EventInsertForecasterPayload{
				Forecaster: "forecaster-payload", Nonce: 104, TopicId: 20, InfererAddresses: []string{"forecast-inferer-a", "forecast-inferer-b"},
				InfererValues: []alloramath.Dec{dec("51.000000000000000001"), dec("52.000000000000000002")}, ExtraData: []byte{3, 4, 129, 130},
				AdmissionOutcome: emissionsv10.PayloadAdmissionOutcome_PAYLOAD_ADMISSION_OUTCOME_NO_ELEMENTS_ACCEPTED,
			},
		},
		{
			name: "EventInsertReputerPayload",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewInsertReputerPayloadEvent(ctx, lossBundle, emissionsv10.PayloadAdmissionOutcome_PAYLOAD_ADMISSION_OUTCOME_ADMITTED)
			},
			expected: &emissionsv10.EventInsertReputerPayload{TopicId: 21, Nonce: 105, Reputer: "reputer-payload", Bundle: eventValueBundle, AdmissionOutcome: emissionsv10.PayloadAdmissionOutcome_PAYLOAD_ADMISSION_OUTCOME_ADMITTED},
		},
		{
			name: "EventInferenceEvictedFromActiveSet",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewInferenceEvictedFromActiveSetEvent(ctx, 22, 106, "evicted-inferer", "evicting-inferer")
			},
			expected: &emissionsv10.EventInferenceEvictedFromActiveSet{TopicId: 22, Nonce: 106, EvictedInferer: "evicted-inferer", EvictingInferer: "evicting-inferer"},
		},
		{
			name: "EventForecastEvictedFromActiveSet",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewForecastEvictedFromActiveSetEvent(ctx, 23, 107, "evicted-forecaster", "evicting-forecaster")
			},
			expected: &emissionsv10.EventForecastEvictedFromActiveSet{TopicId: 23, Nonce: 107, EvictedForecaster: "evicted-forecaster", EvictingForecaster: "evicting-forecaster"},
		},
		{
			name: "EventReputerLossEvictedFromActiveSet",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewReputerLossEvictedFromActiveSetEvent(ctx, 24, 108, "evicted-reputer", "evicting-reputer")
			},
			expected: &emissionsv10.EventReputerLossEvictedFromActiveSet{TopicId: 24, Nonce: 108, EvictedReputer: "evicted-reputer", EvictingReputer: "evicting-reputer"},
		},
		{
			name: "EventReputerSubmissionRetired",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewReputerSubmissionRetiredEvent(ctx, 25, 109, "retired-reputer")
			},
			expected: &emissionsv10.EventReputerSubmissionRetired{TopicId: 25, Nonce: 109, Reputer: "retired-reputer"},
		},
		{
			name: "EventGlobalForecasterWhitelistAdded",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewGlobalForecasterWhitelistAddedEvent(ctx, "global-forecaster-added")
			},
			expected: &emissionsv10.EventGlobalForecasterWhitelistAdded{Address: "global-forecaster-added"},
		},
		{
			name: "EventGlobalForecasterWhitelistRemoved",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewGlobalForecasterWhitelistRemovedEvent(ctx, "global-forecaster-removed")
			},
			expected: &emissionsv10.EventGlobalForecasterWhitelistRemoved{Address: "global-forecaster-removed"},
		},
		{
			name:     "EventTopicForecasterWhitelistEnabled",
			emit:     func(ctx context.Context) { emissionsv10.EmitNewTopicForecasterWhitelistEnabledEvent(ctx, 26) },
			expected: &emissionsv10.EventTopicForecasterWhitelistEnabled{TopicId: 26},
		},
		{
			name:     "EventTopicForecasterWhitelistDisabled",
			emit:     func(ctx context.Context) { emissionsv10.EmitNewTopicForecasterWhitelistDisabledEvent(ctx, 27) },
			expected: &emissionsv10.EventTopicForecasterWhitelistDisabled{TopicId: 27},
		},
		{
			name: "EventTopicForecasterWhitelistAdded",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewTopicForecasterWhitelistAddedEvent(ctx, 28, "topic-forecaster-added")
			},
			expected: &emissionsv10.EventTopicForecasterWhitelistAdded{TopicId: 28, Address: "topic-forecaster-added"},
		},
		{
			name: "EventTopicForecasterWhitelistRemoved",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewTopicForecasterWhitelistRemovedEvent(ctx, 29, "topic-forecaster-removed")
			},
			expected: &emissionsv10.EventTopicForecasterWhitelistRemoved{TopicId: 29, Address: "topic-forecaster-removed"},
		},
		{
			name: "EventRewardPayoutFailed",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewRewardPayoutFailedEvent(ctx, 30, 110, "unpaid-reputer", emissionsv10.ActorType_ACTOR_TYPE_REPUTER, dec("22.000000000000000003"), emissionsv10.RewardPayoutFailureReason_REWARD_PAYOUT_FAILURE_REASON_ACCOUNT_SEND)
			},
			expected: &emissionsv10.EventRewardPayoutFailed{TopicId: 30, Nonce: 110, Address: "unpaid-reputer", ActorType: emissionsv10.ActorType_ACTOR_TYPE_REPUTER, Amount: dec("22.000000000000000003"), Reason: emissionsv10.RewardPayoutFailureReason_REWARD_PAYOUT_FAILURE_REASON_ACCOUNT_SEND},
		},
		{
			name: "EventTopicRewardForfeitedNoRewardNonce",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewTopicRewardForfeitedNoRewardNonceEvent(ctx, 31, 111, cosmosmath.NewInt(31001))
			},
			expected: &emissionsv10.EventTopicRewardForfeitedNoRewardNonce{TopicId: 31, Epoch: 111, Amount: cosmosmath.NewInt(31001)},
		},
		{
			name: "EventTopicRewardForfeitedDistributionFailed",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewTopicRewardForfeitedDistributionFailedEvent(ctx, 32, 112, 113, cosmosmath.NewInt(32002))
			},
			expected: &emissionsv10.EventTopicRewardForfeitedDistributionFailed{TopicId: 32, Epoch: 112, RewardNonce: 113, Amount: cosmosmath.NewInt(32002)},
		},
		{
			name: "EventNonceExpired",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewNonceExpiredEvent(ctx, 33, 114, emissionsv10.ActorType_ACTOR_TYPE_REPUTER, emissionsv10.NonceExpiredReason_NONCE_EXPIRED_REASON_CLOSE_FAILED)
			},
			expected: &emissionsv10.EventNonceExpired{TopicId: 33, Nonce: 114, ActorType: emissionsv10.ActorType_ACTOR_TYPE_REPUTER, Reason: emissionsv10.NonceExpiredReason_NONCE_EXPIRED_REASON_CLOSE_FAILED},
		},
		{
			name: "EventTopicWeightUpdatedAfterStakeChange",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewTopicWeightUpdatedAfterStakeChangeEvent(ctx, 34, dec("3.000000000000000004"), cosmosmath.NewInt(34003), cosmosmath.NewInt(34004))
			},
			expected: &emissionsv10.EventTopicWeightUpdatedAfterStakeChange{TopicId: 34, NewWeight: dec("3.000000000000000004"), TopicStake: cosmosmath.NewInt(34003), TopicFeeRevenue: cosmosmath.NewInt(34004)},
		},
		{
			name: "EventAddStake",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewAddStakeEvent(ctx, 35, "stake-reputer", "stake-delegator", cosmosmath.NewInt(35004), cosmosmath.NewInt(35009), cosmosmath.NewInt(35005))
			},
			expected:     &emissionsv10.EventAddStake{TopicId: 35, Reputer: "stake-reputer", Delegator: "stake-delegator", Amount: cosmosmath.NewInt(35004), TopicStake: cosmosmath.NewInt(35005)},
			ignoreEvents: []string{"emissions.v10.EventReputerStakeUpdated"},
		},
		{
			name: "EventRemoveStake",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewRemoveStakeEvent(ctx, 36, "unstake-reputer", "unstake-delegator", cosmosmath.NewInt(36004), cosmosmath.NewInt(36009), cosmosmath.NewInt(36005))
			},
			expected:     &emissionsv10.EventRemoveStake{TopicId: 36, Reputer: "unstake-reputer", Delegator: "unstake-delegator", Amount: cosmosmath.NewInt(36004), TopicStake: cosmosmath.NewInt(36005)},
			ignoreEvents: []string{"emissions.v10.EventReputerStakeUpdated"},
		},
		{
			name: "EventForecastTaskScoreSet",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewForecastTaskUtilityScoreSetEvent(sdk.UnwrapSDKContext(ctx), 37, dec("3.000000000000000007"), 115)
			},
			expected: &emissionsv10.EventForecastTaskScoreSet{TopicId: 37, Score: dec("3.000000000000000007"), NonceBlockHeight: 115, Nonce: 115},
		},
		{
			name: "EventDelegateRewardShareUpdated",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewDelegateRewardShareUpdatedEvent(ctx, 38, "delegate-reputer", dec("3.000000000000000008"), 116)
			},
			expected: &emissionsv10.EventDelegateRewardShareUpdated{TopicId: 38, Reputer: "delegate-reputer", RewardPerShare: dec("3.000000000000000008"), Nonce: 116},
		},
		{
			name: "EventDelegateRewardDistributed",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewDelegateRewardDistributedEvent(ctx, 39, "delegate-reward-reputer", cosmosmath.NewInt(39003), 117)
			},
			expected: &emissionsv10.EventDelegateRewardDistributed{TopicId: 39, Reputer: "delegate-reward-reputer", Amount: cosmosmath.NewInt(39003), Nonce: 117},
		},
		{
			name: "EventActiveReputersSet",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewActiveReputersSetEvent(ctx, 40, 118, []string{"active-reputer-a", "active-reputer-b"})
			},
			expected: &emissionsv10.EventActiveReputersSet{TopicId: 40, Addresses: []string{"active-reputer-a", "active-reputer-b"}, NonceBlockHeight: 118, Nonce: 118},
		},
		{
			name: "EventActiveInferersSet",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewActiveInferersSetEvent(ctx, 41, 119, []string{"active-inferer-a", "active-inferer-b"})
			},
			expected: &emissionsv10.EventActiveInferersSet{TopicId: 41, Addresses: []string{"active-inferer-a", "active-inferer-b"}, NonceBlockHeight: 119, Nonce: 119},
		},
		{
			name: "EventActiveForecastersSet",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewActiveForecastersSetEvent(ctx, 42, 120, []string{"active-forecaster-a"}, []string{"discarded-forecaster-b", "discarded-forecaster-c"})
			},
			expected: &emissionsv10.EventActiveForecastersSet{TopicId: 42, Addresses: []string{"active-forecaster-a"}, NonceBlockHeight: 120, Nonce: 120, DiscardedAddresses: []string{"discarded-forecaster-b", "discarded-forecaster-c"}},
		},
		{
			name: "EventTopicStatusChanged",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewTopicStatusChangedEvent(ctx, 43, true, emissionsv10.TopicStatusChangeReason_TOPIC_STATUS_CHANGE_REASON_ACTIVATED_BY_FUNDING, dec("4.000000000000000001"), dec("4.000000000000000002"), dec("4.000000000000000003"))
			},
			expected: &emissionsv10.EventTopicStatusChanged{TopicId: 43, IsActive: true, Reason: emissionsv10.TopicStatusChangeReason_TOPIC_STATUS_CHANGE_REASON_ACTIVATED_BY_FUNDING, Weight: dec("4.000000000000000001"), MinWeight: dec("4.000000000000000002"), ThresholdWeight: dec("4.000000000000000003")},
		},
		{
			name: "EventNetworkInferenceInfererWeightsSet",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewNetworkInferenceInfererWeightsSetEvent(ctx, 45, 121, []string{"weight-inferer-a"}, []alloramath.Dec{dec("4.000000000000000005")})
			},
			expected: &emissionsv10.EventNetworkInferenceInfererWeightsSet{TopicId: 45, NonceBlockHeight: 121, Addresses: []string{"weight-inferer-a"}, Weights: []alloramath.Dec{dec("4.000000000000000005")}, Nonce: 121},
		},
		{
			name: "EventNetworkInferenceForecasterWeightsSet",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewNetworkInferenceForecasterWeightsSetEvent(ctx, 46, 122, []string{"weight-forecaster-a"}, []alloramath.Dec{dec("4.000000000000000006")})
			},
			expected: &emissionsv10.EventNetworkInferenceForecasterWeightsSet{TopicId: 46, NonceBlockHeight: 122, Addresses: []string{"weight-forecaster-a"}, Weights: []alloramath.Dec{dec("4.000000000000000006")}, Nonce: 122},
		},
		{
			name: "EventNetworkInferenceInfererRegretsUsedSet",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewNetworkInferenceInfererRegretsUsedSetEvent(ctx, 47, 123, []string{"regret-inferer-a"}, []alloramath.Dec{dec("4.000000000000000007")})
			},
			expected: &emissionsv10.EventNetworkInferenceInfererRegretsUsedSet{TopicId: 47, NonceBlockHeight: 123, Addresses: []string{"regret-inferer-a"}, Regrets: []alloramath.Dec{dec("4.000000000000000007")}, Nonce: 123},
		},
		{
			name: "EventNetworkInferenceForecasterRegretsUsedSet",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewNetworkInferenceForecasterRegretsUsedSetEvent(ctx, 48, 124, []string{"regret-forecaster-a"}, []alloramath.Dec{dec("4.000000000000000008")})
			},
			expected: &emissionsv10.EventNetworkInferenceForecasterRegretsUsedSet{TopicId: 48, NonceBlockHeight: 124, Addresses: []string{"regret-forecaster-a"}, Regrets: []alloramath.Dec{dec("4.000000000000000008")}, Nonce: 124},
		},
		{
			name:     "EventWorkerSubmissionWindowOpened",
			emit:     func(ctx context.Context) { emissionsv10.EmitNewWorkerSubmissionWindowOpenedEvent(ctx, 49, 125, 126) },
			expected: &emissionsv10.EventWorkerSubmissionWindowOpened{TopicId: 49, NonceBlockHeight: 125, WindowEndBlock: 126, Nonce: 125},
		},
		{
			name:     "EventWorkerSubmissionWindowClosed",
			emit:     func(ctx context.Context) { emissionsv10.EmitNewWorkerSubmissionWindowClosedEvent(ctx, 50, 127) },
			expected: &emissionsv10.EventWorkerSubmissionWindowClosed{TopicId: 50, NonceBlockHeight: 127, Nonce: 127},
		},
		{
			name:     "EventReputerSubmissionWindowOpened",
			emit:     func(ctx context.Context) { emissionsv10.EmitNewReputerSubmissionWindowOpenedEvent(ctx, 51, 128, 129) },
			expected: &emissionsv10.EventReputerSubmissionWindowOpened{TopicId: 51, NonceBlockHeight: 128, WindowEndBlock: 129, Nonce: 128},
		},
		{
			name:     "EventReputerSubmissionWindowClosed",
			emit:     func(ctx context.Context) { emissionsv10.EmitNewReputerSubmissionWindowClosedEvent(ctx, 52, 130) },
			expected: &emissionsv10.EventReputerSubmissionWindowClosed{TopicId: 52, NonceBlockHeight: 130, Nonce: 130},
		},
		{
			name: "EventEpochLabelRegistryFrozen",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewEpochLabelRegistryFrozenEvent(ctx, 53, 131, []string{"label-alpha", "label-beta", "label-gamma"})
			},
			expected: &emissionsv10.EventEpochLabelRegistryFrozen{TopicId: 53, NonceBlockHeight: 131, RegistrySize: 3, Nonce: 131, Labels: []string{"label-alpha", "label-beta", "label-gamma"}},
		},
		{
			name: "EventListeningCoefficientsSet",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewListeningCoefficientsSetEvent(ctx, emissionsv10.ActorType_ACTOR_TYPE_FORECASTER, 54, 132, []string{"listening-forecaster"}, []alloramath.Dec{dec("5.000000000000000001")})
			},
			expected: &emissionsv10.EventListeningCoefficientsSet{ActorType: emissionsv10.ActorType_ACTOR_TYPE_FORECASTER, TopicId: 54, BlockHeight: 132, Addresses: []string{"listening-forecaster"}, Coefficients: []alloramath.Dec{dec("5.000000000000000001")}, Nonce: 132},
		},
		{
			name: "EventInfererNetworkRegretSet",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewInfererNetworkRegretSetEvent(ctx, 55, 133, []string{"inferer-regret"}, []alloramath.Dec{dec("5.000000000000000002")})
			},
			expected: &emissionsv10.EventInfererNetworkRegretSet{TopicId: 55, BlockHeight: 133, Addresses: []string{"inferer-regret"}, Regrets: []alloramath.Dec{dec("5.000000000000000002")}, Nonce: 133},
		},
		{
			name: "EventForecasterNetworkRegretSet",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewForecasterNetworkRegretSetEvent(ctx, 56, 134, []string{"forecaster-regret"}, []alloramath.Dec{dec("5.000000000000000003")})
			},
			expected: &emissionsv10.EventForecasterNetworkRegretSet{TopicId: 56, BlockHeight: 134, Addresses: []string{"forecaster-regret"}, Regrets: []alloramath.Dec{dec("5.000000000000000003")}, Nonce: 134},
		},
		{
			name: "EventNaiveInfererNetworkRegretSet",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewNaiveInfererNetworkRegretSetEvent(ctx, 57, 135, []string{"naive-inferer-regret"}, []alloramath.Dec{dec("5.000000000000000004")})
			},
			expected: &emissionsv10.EventNaiveInfererNetworkRegretSet{TopicId: 57, BlockHeight: 135, Addresses: []string{"naive-inferer-regret"}, Regrets: []alloramath.Dec{dec("5.000000000000000004")}, Nonce: 135},
		},
		{
			name: "EventTopicInitialRegretSet",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewTopicInitialRegretSetEvent(ctx, 58, 136, dec("5.000000000000000005"))
			},
			expected: &emissionsv10.EventTopicInitialRegretSet{TopicId: 58, BlockHeight: 136, Regret: dec("5.000000000000000005"), Nonce: 136},
		},
		{
			name: "EventTopicInitialEmaScoreSet",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewTopicInitialEmaScoreSetEvent(ctx, emissionsv10.ActorType_ACTOR_TYPE_FORECASTER, 59, 137, dec("5.000000000000000006"))
			},
			expected: &emissionsv10.EventTopicInitialEmaScoreSet{ActorType: emissionsv10.ActorType_ACTOR_TYPE_FORECASTER, TopicId: 59, BlockHeight: 137, Score: dec("5.000000000000000006"), Nonce: 137},
		},
		{
			name: "EventRegretStdNormSet",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewRegretScaleSetEvent(ctx, 60, 138, dec("5.000000000000000007"))
			},
			expected: &emissionsv10.EventRegretStdNormSet{TopicId: 60, BlockHeight: 138, Stdnorm: dec("5.000000000000000007"), Nonce: 138},
		},
		{
			name: "EventInfererWeightsSet",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewInfererWeightsSetEvent(ctx, 61, 139, []string{"inferer-weight"}, []alloramath.Dec{dec("5.000000000000000008")})
			},
			expected: &emissionsv10.EventInfererWeightsSet{TopicId: 61, BlockHeight: 139, Addresses: []string{"inferer-weight"}, Weights: []alloramath.Dec{dec("5.000000000000000008")}, Nonce: 139},
		},
		{
			name: "EventForecasterWeightsSet",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewForecasterWeightsSetEvent(ctx, 62, 140, []string{"forecaster-weight"}, []alloramath.Dec{dec("5.000000000000000009")})
			},
			expected: &emissionsv10.EventForecasterWeightsSet{TopicId: 62, BlockHeight: 140, Addresses: []string{"forecaster-weight"}, Weights: []alloramath.Dec{dec("5.000000000000000009")}, Nonce: 140},
		},
		{
			name: "EventRewardsSettled",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewInfererRewardsSettledEvent(ctx, 141, 142, []emissionsv10.TaskReward{{TopicId: 63, Address: "reward-inferer", Reward: dec("10.000000000000000125")}})
			},
			expected: &emissionsv10.EventRewardsSettled{ActorType: emissionsv10.ActorType_ACTOR_TYPE_INFERER_UNSPECIFIED, TopicId: 63, BlockHeight: 141, Addresses: []string{"reward-inferer"}, Rewards: []alloramath.Dec{dec("10.000000000000000125")}, BlockHeightTx: 142, Nonce: 141},
		},
		{
			name: "EventScoresSet",
			emit: func(ctx context.Context) {
				emissionsv10.EmitNewActorScoresSetEvent(ctx, emissionsv10.ActorType_ACTOR_TYPE_FORECASTER, []emissionsv10.Score{{TopicId: 64, BlockHeight: 143, Address: "score-forecaster", Score: dec("123.456789012345678901")}})
			},
			expected: &emissionsv10.EventScoresSet{ActorType: emissionsv10.ActorType_ACTOR_TYPE_FORECASTER, TopicId: 64, BlockHeight: 143, Addresses: []string{"score-forecaster"}, Scores: []alloramath.Dec{dec("123.456789012345678901")}, Nonce: 143},
		},
	}
}
