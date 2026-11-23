package codec_test

import (
	"reflect"
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"

	alloramath "github.com/allora-network/allora-chain/math"
	emissionsv10 "github.com/allora-network/allora-chain/x/emissions/types"
	"github.com/allora-network/allora-sdk-go/codec"
)

// TestCodecDecodesV10EventFieldMappings covers nonce aliases and the actor type that chain emitters duplicate or set to zero.
func TestCodecDecodesV10EventFieldMappings(t *testing.T) {
	cdc := codec.NewCodec()
	for _, expected := range fieldMappingEvents() {
		t.Run(proto.MessageName(expected), func(t *testing.T) {
			requireDistinctNonZeroScalars(t, expected)
			event, err := sdk.TypedEventToEvent(expected)
			require.NoError(t, err)
			abciEvent := abcitypes.Event(event)
			decoded, err := cdc.ParseTypedEvent(&abciEvent)
			require.NoError(t, err)
			requireMessageFields(t, expected, decoded)
		})
	}
}

func fieldMappingEvents() []proto.Message {
	dec := alloramath.MustNewDecFromString
	return []proto.Message{
		&emissionsv10.EventForecastTaskScoreSet{TopicId: 1, Score: dec("1.000000000000000002"), NonceBlockHeight: 3, Nonce: 4},
		&emissionsv10.EventActiveReputersSet{TopicId: 5, Addresses: []string{"active-reputer"}, NonceBlockHeight: 6, Nonce: 7},
		&emissionsv10.EventActiveInferersSet{TopicId: 8, Addresses: []string{"active-inferer"}, NonceBlockHeight: 9, Nonce: 10},
		&emissionsv10.EventActiveForecastersSet{TopicId: 11, Addresses: []string{"active-forecaster"}, NonceBlockHeight: 12, Nonce: 13, DiscardedAddresses: []string{"discarded-forecaster"}},
		&emissionsv10.EventNetworkInferenceInfererWeightsSet{TopicId: 14, NonceBlockHeight: 15, Addresses: []string{"network-inferer"}, Weights: []alloramath.Dec{dec("16.000000000000000017")}, Nonce: 18},
		&emissionsv10.EventNetworkInferenceForecasterWeightsSet{TopicId: 19, NonceBlockHeight: 20, Addresses: []string{"network-forecaster"}, Weights: []alloramath.Dec{dec("21.000000000000000022")}, Nonce: 23},
		&emissionsv10.EventNetworkInferenceInfererRegretsUsedSet{TopicId: 24, NonceBlockHeight: 25, Addresses: []string{"regret-inferer"}, Regrets: []alloramath.Dec{dec("26.000000000000000027")}, Nonce: 28},
		&emissionsv10.EventNetworkInferenceForecasterRegretsUsedSet{TopicId: 29, NonceBlockHeight: 30, Addresses: []string{"regret-forecaster"}, Regrets: []alloramath.Dec{dec("31.000000000000000032")}, Nonce: 33},
		&emissionsv10.EventWorkerSubmissionWindowOpened{TopicId: 34, NonceBlockHeight: 35, WindowEndBlock: 36, Nonce: 37},
		&emissionsv10.EventWorkerSubmissionWindowClosed{TopicId: 38, NonceBlockHeight: 39, Nonce: 40},
		&emissionsv10.EventReputerSubmissionWindowOpened{TopicId: 41, NonceBlockHeight: 42, WindowEndBlock: 43, Nonce: 44},
		&emissionsv10.EventReputerSubmissionWindowClosed{TopicId: 45, NonceBlockHeight: 46, Nonce: 47},
		&emissionsv10.EventEpochLabelRegistryFrozen{TopicId: 48, NonceBlockHeight: 49, RegistrySize: 50, Nonce: 51, Labels: []string{"label-alpha", "label-beta"}},
		&emissionsv10.EventListeningCoefficientsSet{ActorType: emissionsv10.ActorType_ACTOR_TYPE_FORECASTER, TopicId: 52, BlockHeight: 53, Addresses: []string{"listening-forecaster"}, Coefficients: []alloramath.Dec{dec("54.000000000000000055")}, Nonce: 56},
		&emissionsv10.EventInfererNetworkRegretSet{TopicId: 57, BlockHeight: 58, Addresses: []string{"inferer-regret"}, Regrets: []alloramath.Dec{dec("59.000000000000000060")}, Nonce: 61},
		&emissionsv10.EventForecasterNetworkRegretSet{TopicId: 62, BlockHeight: 63, Addresses: []string{"forecaster-regret"}, Regrets: []alloramath.Dec{dec("64.000000000000000065")}, Nonce: 66},
		&emissionsv10.EventNaiveInfererNetworkRegretSet{TopicId: 67, BlockHeight: 68, Addresses: []string{"naive-inferer-regret"}, Regrets: []alloramath.Dec{dec("69.000000000000000070")}, Nonce: 71},
		&emissionsv10.EventTopicInitialRegretSet{TopicId: 72, BlockHeight: 73, Regret: dec("74.000000000000000075"), Nonce: 76},
		&emissionsv10.EventTopicInitialEmaScoreSet{ActorType: emissionsv10.ActorType_ACTOR_TYPE_REPUTER, TopicId: 77, BlockHeight: 78, Score: dec("79.000000000000000080"), Nonce: 81},
		&emissionsv10.EventRegretStdNormSet{TopicId: 82, BlockHeight: 83, Stdnorm: dec("84.000000000000000085"), Nonce: 86},
		&emissionsv10.EventInfererWeightsSet{TopicId: 87, BlockHeight: 88, Addresses: []string{"inferer-weight"}, Weights: []alloramath.Dec{dec("89.000000000000000090")}, Nonce: 91},
		&emissionsv10.EventForecasterWeightsSet{TopicId: 92, BlockHeight: 93, Addresses: []string{"forecaster-weight"}, Weights: []alloramath.Dec{dec("94.000000000000000095")}, Nonce: 96},
		&emissionsv10.EventRewardsSettled{ActorType: emissionsv10.ActorType_ACTOR_TYPE_REPUTER, TopicId: 97, BlockHeight: 98, Addresses: []string{"reward-reputer"}, Rewards: []alloramath.Dec{dec("99.000000000000000100")}, BlockHeightTx: 101, Nonce: 102},
		&emissionsv10.EventScoresSet{ActorType: emissionsv10.ActorType_ACTOR_TYPE_FORECASTER, TopicId: 103, BlockHeight: 104, Addresses: []string{"score-forecaster"}, Scores: []alloramath.Dec{dec("105.000000000000000106")}, Nonce: 107},
	}
}

// requireDistinctNonZeroScalars ensures each scalar field can expose a dropped or swapped decode.
func requireDistinctNonZeroScalars(t *testing.T, expected proto.Message) {
	t.Helper()
	value := reflect.ValueOf(expected).Elem()
	typeOfValue := value.Type()
	values := make(map[any]string)
	for i := range value.NumField() {
		field := value.Field(i)
		fieldName := typeOfValue.Field(i).Name
		require.Falsef(t, field.IsZero(), "scalar field %s has its zero value", fieldName)
		if !field.Type().Comparable() {
			continue
		}
		value := field.Interface()
		previous, duplicate := values[value]
		require.Falsef(t, duplicate, "scalar fields %s and %s both have value %v", previous, fieldName, value)
		values[value] = fieldName
	}
}
