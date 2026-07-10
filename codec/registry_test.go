package codec_test

import (
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cosmossdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"

	alloramath "github.com/allora-network/allora-chain/math"
	emissionsv9 "github.com/allora-network/allora-chain/x/emissions/api/emissions/v9"
	emissionsv10 "github.com/allora-network/allora-chain/x/emissions/types"
	"github.com/allora-network/allora-sdk-go/codec"
)

func TestCodecRecognizesLegacyAndCurrentEmissionsEvents(t *testing.T) {
	cdc := codec.NewCodec()

	require.True(t, cdc.IsTypedEvent(&abcitypes.Event{Type: "emissions.v9.EventScoresSet"}))
	require.True(t, cdc.IsTypedEvent(&abcitypes.Event{Type: "emissions.v10.EventScoresSet"}))
	require.True(t, cdc.IsTypedEvent(&abcitypes.Event{Type: "emissions.v10.EventNetworkInferenceBundle"}))
	require.True(t, cdc.IsTypedEvent(&abcitypes.Event{Type: "emissions.v10.EventEpochLabelRegistryFrozen"}))
}

func TestCodecParsesLegacyAndCurrentEmissionsEvents(t *testing.T) {
	cdc := codec.NewCodec()

	v9ABCIEvent, err := typedEventToABCIEvent(&emissionsv9.EventScoresSet{
		ActorType:   emissionsv9.ActorType_ACTOR_TYPE_INFERER_UNSPECIFIED,
		TopicId:     42,
		BlockHeight: 12340,
		Addresses:   []string{"addr1", "addr2"},
		Scores:      []string{"95.5", "87.2"},
	})
	require.NoError(t, err)

	v9Event, err := cdc.ParseTypedEvent(&v9ABCIEvent)
	require.NoError(t, err)
	parsedV9Event, ok := v9Event.(*emissionsv9.EventScoresSet)
	require.True(t, ok)
	require.Equal(t, emissionsv9.ActorType_ACTOR_TYPE_INFERER_UNSPECIFIED, parsedV9Event.ActorType)
	require.Equal(t, uint64(42), parsedV9Event.TopicId)
	require.Equal(t, int64(12340), parsedV9Event.BlockHeight)
	require.Equal(t, []string{"addr1", "addr2"}, parsedV9Event.Addresses)
	require.Equal(t, []string{"95.5", "87.2"}, parsedV9Event.Scores)

	v10ABCIEvent, err := typedEventToABCIEvent(&emissionsv10.EventScoresSet{
		ActorType:   emissionsv10.ActorType_ACTOR_TYPE_INFERER_UNSPECIFIED,
		TopicId:     42,
		BlockHeight: 12340,
		Addresses:   []string{"addr1", "addr2"},
		Scores: []alloramath.Dec{
			alloramath.MustNewDecFromString("95.5"),
			alloramath.MustNewDecFromString("87.2"),
		},
	})
	require.NoError(t, err)

	v10Event, err := cdc.ParseTypedEvent(&v10ABCIEvent)
	require.NoError(t, err)
	parsedV10Event, ok := v10Event.(*emissionsv10.EventScoresSet)
	require.True(t, ok)
	require.Equal(t, emissionsv10.ActorType_ACTOR_TYPE_INFERER_UNSPECIFIED, parsedV10Event.ActorType)
	require.Equal(t, uint64(42), parsedV10Event.TopicId)
	require.Equal(t, int64(12340), parsedV10Event.BlockHeight)
	require.Equal(t, []string{"addr1", "addr2"}, parsedV10Event.Addresses)
	require.Equal(t, []alloramath.Dec{
		alloramath.MustNewDecFromString("95.5"),
		alloramath.MustNewDecFromString("87.2"),
	}, parsedV10Event.Scores)
}

func TestCodecRegistersAllLegacyEmissionsV9Events(t *testing.T) {
	for _, msg := range []proto.Message{
		&emissionsv9.EventScoresSet{},
		&emissionsv9.EventRewardsSettled{},
		// EventNetworkLossSet, EventNetworkInferences, EventOutlierResistantNetworkInferences and
		// EventValueBundle are intentionally omitted: their typed decoding is currently unsupported
		// because they use multi-dimensional arrays, so they are not registered by registerV9TypedEvents.
		&emissionsv9.EventInsertInfererPayload{},
		&emissionsv9.EventInsertForecasterPayload{},
		&emissionsv9.EventCreateNewTopic{},
		&emissionsv9.EventTopicUpdated{},
		&emissionsv9.EventAddStake{},
		&emissionsv9.EventRemoveStake{},
		&emissionsv9.EventRequestStakeRemoval{},
		&emissionsv9.EventCancelStakeRemoval{},
		&emissionsv9.EventReputerStakeUpdated{},
		&emissionsv9.EventRewardDelegateStake{},
		&emissionsv9.EventInsertReputerPayload{},
		&emissionsv9.EventReputerRegistered{},
		&emissionsv9.EventWorkerRegistered{},
		&emissionsv9.EventNodeOwnerUpdated{},
		&emissionsv9.EventReputerUnregistered{},
		&emissionsv9.EventWorkerUnregistered{},
		&emissionsv9.EventFundTopic{},
		&emissionsv9.EventParamsSet{},
		&emissionsv9.EventWhitelistAdminAdded{},
		&emissionsv9.EventWhitelistAdminRemoved{},
		&emissionsv9.EventGlobalWhitelistAdded{},
		&emissionsv9.EventGlobalWhitelistRemoved{},
		&emissionsv9.EventGlobalWorkerWhitelistAdded{},
		&emissionsv9.EventGlobalWorkerWhitelistRemoved{},
		&emissionsv9.EventGlobalReputerWhitelistAdded{},
		&emissionsv9.EventGlobalReputerWhitelistRemoved{},
		&emissionsv9.EventGlobalAdminWhitelistAdded{},
		&emissionsv9.EventGlobalAdminWhitelistRemoved{},
		&emissionsv9.EventTopicWorkerWhitelistEnabled{},
		&emissionsv9.EventTopicWorkerWhitelistDisabled{},
		&emissionsv9.EventTopicReputerWhitelistEnabled{},
		&emissionsv9.EventTopicReputerWhitelistDisabled{},
		&emissionsv9.EventTopicCreatorWhitelistAdded{},
		&emissionsv9.EventTopicCreatorWhitelistRemoved{},
		&emissionsv9.EventTopicWorkerWhitelistAdded{},
		&emissionsv9.EventTopicWorkerWhitelistRemoved{},
		&emissionsv9.EventTopicReputerWhitelistAdded{},
		&emissionsv9.EventTopicReputerWhitelistRemoved{},
		&emissionsv9.EventForecastTaskScoreSet{},
		&emissionsv9.EventWorkerLastCommitSet{},
		&emissionsv9.EventReputerLastCommitSet{},
		&emissionsv9.EventTopicRewardsSet{},
		&emissionsv9.EventEMAScoresSet{},
		&emissionsv9.EventListeningCoefficientsSet{},
		&emissionsv9.EventInfererNetworkRegretSet{},
		&emissionsv9.EventForecasterNetworkRegretSet{},
		&emissionsv9.EventNaiveInfererNetworkRegretSet{},
		&emissionsv9.EventTopicInitialRegretSet{},
		&emissionsv9.EventTopicInitialEmaScoreSet{},
		&emissionsv9.EventRegretStdNormSet{},
		&emissionsv9.EventInfererWeightsSet{},
		&emissionsv9.EventForecasterWeightsSet{},
		&emissionsv9.EventPreviousPercentageRewardToStakedReputersSet{},
		&emissionsv9.EventPruneRecords{},
		&emissionsv9.EventDelegateRewardShareUpdated{},
		&emissionsv9.EventDelegateRewardDistributed{},
		&emissionsv9.EventActiveReputersSet{},
		&emissionsv9.EventActiveInferersSet{},
		&emissionsv9.EventActiveForecastersSet{},
		&emissionsv9.EventTopicStatusChanged{},
		&emissionsv9.EventNetworkInferenceInfererWeightsSet{},
		&emissionsv9.EventNetworkInferenceForecasterWeightsSet{},
		&emissionsv9.EventNetworkInferenceInfererRegretsUsedSet{},
		&emissionsv9.EventNetworkInferenceForecasterRegretsUsedSet{},
		&emissionsv9.EventTopicWeightUpdated{},
		&emissionsv9.EventTopicFeeRevenueDripped{},
		&emissionsv9.EventWorkerSubmissionWindowOpened{},
		&emissionsv9.EventWorkerSubmissionWindowClosed{},
		&emissionsv9.EventReputerSubmissionWindowOpened{},
		&emissionsv9.EventReputerSubmissionWindowClosed{},
	} {
		name := proto.MessageName(msg)
		require.NotNilf(t, proto.MessageType(name), "expected %s to be registered", name)
	}
}

func TestCodecUnpacksLegacyAndCurrentEmissionsMessages(t *testing.T) {
	cdc := codec.NewCodec()

	v9Msg, err := codectypes.NewAnyWithValue(&emissionsv9.AddStakeRequest{})
	require.NoError(t, err)
	parsedV9Msg, err := cdc.ParseTxMessage(v9Msg)
	require.NoError(t, err)
	require.IsType(t, &emissionsv9.AddStakeRequest{}, parsedV9Msg)

	v10Msg, err := codectypes.NewAnyWithValue(&emissionsv10.AddStakeRequest{})
	require.NoError(t, err)
	parsedV10Msg, err := cdc.ParseTxMessage(v10Msg)
	require.NoError(t, err)
	require.IsType(t, &emissionsv10.AddStakeRequest{}, parsedV10Msg)
}

func typedEventToABCIEvent(msg proto.Message) (abcitypes.Event, error) {
	event, err := cosmossdktypes.TypedEventToEvent(msg)
	if err != nil {
		return abcitypes.Event{}, err
	}

	return abcitypes.Event(event), nil
}
