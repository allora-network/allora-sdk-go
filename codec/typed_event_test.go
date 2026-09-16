package codec_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"

	alloramath "github.com/allora-network/allora-chain/math"
	emissionsv9 "github.com/allora-network/allora-chain/x/emissions/api/emissions/v9"
	emissionsv10 "github.com/allora-network/allora-chain/x/emissions/types"
	"github.com/allora-network/allora-sdk-go/codec"
)

// loadEventFixture reads a testdata file holding a JSON array of ABCI events
// (trimmed from real `block_results` responses of the testnet archive node)
// and returns the first event.
func loadEventFixture(t *testing.T, name string) abcitypes.Event {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)
	var events []abcitypes.Event
	require.NoError(t, json.Unmarshal(raw, &events))
	require.NotEmpty(t, events, "fixture %s has no events", name)
	return events[0]
}

func attrKeys(ev *abcitypes.Event) []string {
	keys := make([]string, 0, len(ev.Attributes))
	for _, a := range ev.Attributes {
		keys = append(keys, a.Key)
	}
	return keys
}

func cloneEvent(ev abcitypes.Event) abcitypes.Event {
	out := ev
	out.Attributes = append([]abcitypes.EventAttribute(nil), ev.Attributes...)
	return out
}

// Legacy (pulsar/protov2) tx events carry baseapp's `msg_index` attribute.
// Mutation check: revert ParseTypedEvent to cosmossdktypes.ParseTypedEvent
// and this fails with `unknown field "msg_index"`.
func TestParseTypedEvent_LegacyTxEventWithMsgIndex(t *testing.T) {
	cdc := codec.NewCodec()
	ev := loadEventFixture(t, "v9_tx_insert_inferer_payload_9977429.json")
	require.Contains(t, attrKeys(&ev), "msg_index")
	before := cloneEvent(ev)

	require.True(t, cdc.IsTypedEvent(&ev))
	msg, err := cdc.ParseTypedEvent(&ev)
	require.NoError(t, err)

	payload, ok := msg.(*emissionsv9.EventInsertInfererPayload)
	require.True(t, ok, "got %T", msg)
	require.Equal(t, uint64(71), payload.TopicId)
	require.Equal(t, int64(9977427), payload.Nonce)
	require.Equal(t, "allo1xh7kcr6vrkp9lwsu2qnf0wvt93u23pf63aqrpz", payload.Inferer)
	require.Equal(t, "-0.012403974191526725", payload.Value)

	require.Equal(t, before, ev, "input event must not be mutated")
}

// EndBlock events carry an unquoted `mode` attribute. The chain emits it last,
// but the producer grafts `total_stake` onto Add/RemoveStake events after it,
// so `mode` must be excluded wherever it sits.
func TestParseTypedEvent_ModeAttributeAnyPosition(t *testing.T) {
	cdc := codec.NewCodec()
	base := loadEventFixture(t, "v9_endblock_add_stake_9977429.json")
	require.Equal(t, "mode", base.Attributes[len(base.Attributes)-1].Key)

	grafted := cloneEvent(base)
	grafted.Attributes = append(grafted.Attributes, abcitypes.EventAttribute{Key: "total_stake", Value: `"99"`})

	modeFirst := cloneEvent(base)
	modeFirst.Attributes = append(
		[]abcitypes.EventAttribute{modeFirst.Attributes[len(modeFirst.Attributes)-1]},
		modeFirst.Attributes[:len(modeFirst.Attributes)-1]...,
	)

	for name, ev := range map[string]abcitypes.Event{
		"mode last":                    base,
		"unknown attribute after mode": grafted,
		"mode first":                   modeFirst,
	} {
		t.Run(name, func(t *testing.T) {
			before := cloneEvent(ev)
			msg, err := cdc.ParseTypedEvent(&ev)
			require.NoError(t, err)
			stake, ok := msg.(*emissionsv9.EventAddStake)
			require.True(t, ok, "got %T", msg)
			require.Equal(t, uint64(41), stake.TopicId)
			require.Equal(t, "32660430", stake.Amount)
			require.Equal(t, "allo1aanuhulxtl2me3ygkwu8jghkjjrslhhgrluad5", stake.Reputer)
			require.Equal(t, before, ev, "input event must not be mutated")
		})
	}
}

// The v9 network inference/loss events put a 2-D array where the pulsar type
// declares `repeated string`; they are registered but undecodable and must be
// reported untyped so callers keep the raw attribute JSON.
func TestIsTypedEvent_LegacyNestedArrayEventsAreUntyped(t *testing.T) {
	cdc := codec.NewCodec()
	ev := loadEventFixture(t, "v9_endblock_network_loss_set_9977446.json")

	require.False(t, cdc.IsTypedEvent(&ev))
	_, err := cdc.ParseTypedEvent(&ev)
	require.Error(t, err)

	raw, err := cdc.ParseUntypedEvent(&ev)
	require.NoError(t, err)
	var attrs map[string]string
	require.NoError(t, json.Unmarshal(raw, &attrs))
	require.Contains(t, attrs, "bundle")

	for _, typ := range []string{
		"emissions.v9.EventNetworkInferences",
		"emissions.v9.EventOutlierResistantNetworkInferences",
		"emissions.v9.EventInsertReputerPayload",
		"emissions.v9.EventValueBundle",
	} {
		require.False(t, cdc.IsTypedEvent(&abcitypes.Event{Type: typ}), typ)
	}
}

// The v9 reputer payload embeds the same EventValueBundle as the loss/inference
// events, so its `bundle` attribute carries the 2-D array too. It is a tx event
// (carries `msg_index`) and must be reported untyped with the raw JSON intact.
func TestIsTypedEvent_LegacyReputerPayloadIsUntyped(t *testing.T) {
	cdc := codec.NewCodec()
	ev := loadEventFixture(t, "v9_tx_insert_reputer_payload_9977185.json")
	require.Equal(t, "emissions.v9.EventInsertReputerPayload", ev.Type)

	require.False(t, cdc.IsTypedEvent(&ev))
	_, err := cdc.ParseTypedEvent(&ev)
	require.Error(t, err)

	raw, err := cdc.ParseUntypedEvent(&ev)
	require.NoError(t, err)
	var attrs map[string]string
	require.NoError(t, json.Unmarshal(raw, &attrs))
	require.Contains(t, attrs, "bundle")
	require.Contains(t, attrs, "reputer")
	require.Contains(t, attrs, "topic_id")

	var bundle struct {
		OneOutInfererForecasterValues [][]string `json:"one_out_inferer_forecaster_values"`
	}
	require.NoError(t, json.Unmarshal([]byte(attrs["bundle"]), &bundle))
	require.NotEmpty(t, bundle.OneOutInfererForecasterValues, "fixture must carry the nested array that breaks protojson")
	require.NotEmpty(t, bundle.OneOutInfererForecasterValues[0])
}

// The current (gogo) EventNetworkLossSet decodes its nested DecArray through
// the customtype's UnmarshalJSON, so it stays on the typed path.
func TestParseTypedEvent_CurrentNetworkLossSetDecodes(t *testing.T) {
	cdc := codec.NewCodec()
	ev := loadEventFixture(t, "v10_endblock_network_loss_set_10090000.json")

	require.True(t, cdc.IsTypedEvent(&ev))
	msg, err := cdc.ParseTypedEvent(&ev)
	require.NoError(t, err)
	loss, ok := msg.(*emissionsv10.EventNetworkLossSet)
	require.True(t, ok, "got %T", msg)
	require.Equal(t, uint64(76), loss.TopicId)
	require.Equal(t, int64(10089760), loss.Nonce)
	require.NotNil(t, loss.Bundle)
	require.NotEmpty(t, loss.Bundle.OneOutInfererForecasterValues)
	require.IsType(t, alloramath.DecArray{}, loss.Bundle.OneOutInfererForecasterValues[0])
}

// Current (gogo) tx events also carry `msg_index`; the gogo jsonpb path must
// keep ignoring it.
func TestParseTypedEvent_CurrentTxEventWithMsgIndex(t *testing.T) {
	cdc := codec.NewCodec()
	ev := loadEventFixture(t, "v10_tx_insert_inferer_payload_10500000.json")
	require.Contains(t, attrKeys(&ev), "msg_index")

	msg, err := cdc.ParseTypedEvent(&ev)
	require.NoError(t, err)
	payload, ok := msg.(*emissionsv10.EventInsertInfererPayload)
	require.True(t, ok, "got %T", msg)
	require.Equal(t, uint64(83), payload.TopicId)
	require.Equal(t, int64(10499999), payload.Nonce)
	require.Equal(t, "allo19f8ljwal8hrsfv7udvv0pzcgyd2hukd6whh550", payload.Inferer)
	// Customtype (Dec), bytes-from-null, and empty repeated fields all go
	// through gogo jsonpb's own handling; pin them so a regression on the gogo
	// path is caught, not just the protov2 one.
	require.Equal(t, "0.126859700106903839", payload.Value.String())
	require.Empty(t, payload.ExtraData)
	require.Empty(t, payload.Values)
}

// Any unknown attribute — not just msg_index — is ignored on both decode paths.
func TestParseTypedEvent_UnknownAttributesIgnored(t *testing.T) {
	cdc := codec.NewCodec()

	for name, tc := range map[string]struct {
		fixture string
		topicID func(proto.Message) uint64
	}{
		"legacy protov2": {
			fixture: "v9_tx_insert_inferer_payload_9977429.json",
			topicID: func(m proto.Message) uint64 { return m.(*emissionsv9.EventInsertInfererPayload).TopicId },
		},
		"current gogo": {
			fixture: "v10_tx_insert_inferer_payload_10500000.json",
			topicID: func(m proto.Message) uint64 { return m.(*emissionsv10.EventInsertInfererPayload).TopicId },
		},
	} {
		t.Run(name, func(t *testing.T) {
			ev := loadEventFixture(t, tc.fixture)
			ev.Attributes = append(ev.Attributes, abcitypes.EventAttribute{Key: "some_future_field", Value: `{"nested":[1,2]}`})
			msg, err := cdc.ParseTypedEvent(&ev)
			require.NoError(t, err)
			// Ignoring the unknown key must not cost the known ones.
			require.NotZero(t, tc.topicID(msg))
		})
	}
}

// A typed event never carries the same key twice (EmitTypedEvent builds the
// attributes from a JSON object); a duplicate is upstream corruption and must
// be refused rather than decoded last-wins.
func TestParseTypedEvent_DuplicateAttributeKeyRejected(t *testing.T) {
	cdc := codec.NewCodec()
	ev := loadEventFixture(t, "v9_endblock_add_stake_9977429.json")
	ev.Attributes = append(ev.Attributes, abcitypes.EventAttribute{Key: "amount", Value: `"1"`})
	before := cloneEvent(ev)

	_, err := cdc.ParseTypedEvent(&ev)
	require.Error(t, err)
	require.Contains(t, err.Error(), `duplicate attribute key "amount"`)
	require.Equal(t, before, ev, "input event must not be mutated")
}

func TestParseTypedEvent_NoAttributes(t *testing.T) {
	cdc := codec.NewCodec()
	_, err := cdc.ParseTypedEvent(&abcitypes.Event{Type: "emissions.v9.EventAddStake"})
	require.Error(t, err)
}
