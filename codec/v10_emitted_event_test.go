package codec_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"

	"github.com/allora-network/allora-sdk-go/codec"
)

// updateEmittedEventFixtures selects fixture regeneration from the pinned chain emitters.
var updateEmittedEventFixtures = os.Getenv("UPDATE_EMITTED_EVENT_FIXTURES") == "1"

type emittedEventFixture struct {
	name         string
	emit         func(context.Context)
	expected     proto.Message
	ignoreEvents []string
}

func TestCodecDecodesEmittedV10Events(t *testing.T) {
	cdc := codec.NewCodec()
	for _, tc := range emittedEventFixtures() {
		t.Run(tc.name, func(t *testing.T) {
			ctx := sdk.Context{}.WithEventManager(sdk.NewEventManager())
			tc.emit(ctx)
			events := ctx.EventManager().ABCIEvents()
			require.Len(t, events, len(tc.ignoreEvents)+1)
			event := events[0]
			for i, eventType := range tc.ignoreEvents {
				require.Equal(t, eventType, events[i+1].Type)
			}
			require.Equal(t, "emissions.v10."+tc.name, event.Type)
			require.True(t, cdc.IsTypedEvent(&event))

			path := filepath.Join("testdata", "generated_v10", proto.MessageName(tc.expected)+".json")
			if updateEmittedEventFixtures {
				writeEmittedEventFixture(t, path, event)
			}
			raw, err := os.ReadFile(path)
			require.NoError(t, err)
			var stored []abcitypes.Event
			require.NoError(t, json.Unmarshal(raw, &stored))
			require.Equal(t, event, stored[0])

			decoded, err := cdc.ParseTypedEvent(&event)
			require.NoError(t, err)
			requireMessageFields(t, tc.expected, decoded)
		})
	}
}

func writeEmittedEventFixture(t *testing.T, path string, event abcitypes.Event) {
	t.Helper()
	raw, err := json.MarshalIndent([]abcitypes.Event{event}, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, append(raw, '\n'), 0o644))
}

func requireMessageFields(t *testing.T, expected proto.Message, decoded proto.Message) {
	t.Helper()
	require.IsType(t, expected, decoded)
	expectedValue := reflect.ValueOf(expected).Elem()
	decodedValue := reflect.ValueOf(decoded).Elem()
	for i := range expectedValue.NumField() {
		fieldName := expectedValue.Type().Field(i).Name
		require.Equalf(t, expectedValue.Field(i).Interface(), decodedValue.Field(i).Interface(), "field %s", fieldName)
	}
}
