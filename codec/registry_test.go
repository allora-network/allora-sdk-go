package codec_test

import (
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"

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
	attrs := []abcitypes.EventAttribute{
		{Key: "topic_id", Value: "\"42\""},
		{Key: "block_height", Value: "\"12340\""},
		{Key: "addresses", Value: "[\"addr1\", \"addr2\"]"},
		{Key: "scores", Value: "[\"95.5\", \"87.2\"]"},
	}

	v9Event, err := cdc.ParseTypedEvent(&abcitypes.Event{
		Type:       "emissions.v9.EventScoresSet",
		Attributes: attrs,
	})
	require.NoError(t, err)
	require.IsType(t, &emissionsv9.EventScoresSet{}, v9Event)

	v10Event, err := cdc.ParseTypedEvent(&abcitypes.Event{
		Type:       "emissions.v10.EventScoresSet",
		Attributes: attrs,
	})
	require.NoError(t, err)
	require.IsType(t, &emissionsv10.EventScoresSet{}, v10Event)
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
