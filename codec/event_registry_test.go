package codec_test

import (
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/reflect/protoreflect"

	alloramath "github.com/allora-network/allora-chain/math"
	emissionsv10api "github.com/allora-network/allora-chain/x/emissions/api/emissions/v10"
	emissionsv10 "github.com/allora-network/allora-chain/x/emissions/types"
	mintv3api "github.com/allora-network/allora-chain/x/mint/api/mint/v3"
	"github.com/allora-network/allora-sdk-go/codec"
)

func TestCodecRegistersAllEventMessages(t *testing.T) {
	cdc := codec.NewCodec()
	files := []protoreflect.FileDescriptor{
		emissionsv10api.File_emissions_v10_events_proto,
		mintv3api.File_mint_v3_events_proto,
	}
	require.NotEmpty(t, files)
	for _, file := range files {
		messages := file.Messages()
		for i := range messages.Len() {
			name := string(messages.Get(i).FullName())
			t.Run(name, func(t *testing.T) {
				event := abcitypes.Event{Type: name, Attributes: []abcitypes.EventAttribute{{Key: "unused", Value: "null"}}}
				require.True(t, cdc.IsTypedEvent(&event))
				_, err := cdc.ParseTypedEvent(&event)
				require.NoError(t, err)
			})
		}
	}
}

func TestEmittedEnumsUseNames(t *testing.T) {
	ctx := sdk.Context{}.WithEventManager(sdk.NewEventManager())
	dec := alloramath.MustNewDecFromString
	emissionsv10.EmitNewRewardPayoutFailedEvent(ctx, 44, 1301, "unpaid-reputer", emissionsv10.ActorType_ACTOR_TYPE_REPUTER, dec("22.000000000000000003"), emissionsv10.RewardPayoutFailureReason_REWARD_PAYOUT_FAILURE_REASON_ACCOUNT_SEND)
	attrs := ctx.EventManager().ABCIEvents()[0].Attributes
	values := make(map[string]string, len(attrs))
	for _, attr := range attrs {
		values[attr.Key] = attr.Value
	}
	require.Equal(t, `"ACTOR_TYPE_REPUTER"`, values["actor_type"])
	require.Equal(t, `"REWARD_PAYOUT_FAILURE_REASON_ACCOUNT_SEND"`, values["reason"])
}
