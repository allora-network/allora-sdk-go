package codec

import (
	"encoding/json"

	"cosmossdk.io/x/feegrant"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/brynbellomy/go-utils/errors"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	cosmoscodec "github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	stdtypes "github.com/cosmos/cosmos-sdk/std"
	cosmossdktypes "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/grpc/encoding"

	// IBC modules
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	ibcconnection "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	ibcchannel "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"

	ibclightclient "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"

	// OLD VERSIONS of MINT TRANSACTIONS
	mintv1beta1 "github.com/allora-network/allora-chain/x/mint/api/mint/v1beta1"
	mintv2 "github.com/allora-network/allora-chain/x/mint/api/mint/v2"
	mintv5 "github.com/allora-network/allora-chain/x/mint/types"

	// mintv3 "github.com/allora-network/allora-chain/x/mint/api/mint/v3"
	// mintv4 "github.com/allora-network/allora-chain/x/mint/api/mint/v4"
	// mintv5 "github.com/allora-network/allora-chain/x/mint/api/mint/v5"

	// emissions "github.com/allora-network/allora-chain/x/emissions/types"
	// OLD VERSIONS of EMISSIONS TRANSACTIONS
	emissionsv2 "github.com/allora-network/allora-chain/x/emissions/api/emissions/v2"
	emissionsv3 "github.com/allora-network/allora-chain/x/emissions/api/emissions/v3"
	emissionsv4 "github.com/allora-network/allora-chain/x/emissions/api/emissions/v4"
	emissionsv5 "github.com/allora-network/allora-chain/x/emissions/api/emissions/v5"
	emissionsv6 "github.com/allora-network/allora-chain/x/emissions/api/emissions/v6"
	emissionsv7 "github.com/allora-network/allora-chain/x/emissions/api/emissions/v7"
	emissionsv8 "github.com/allora-network/allora-chain/x/emissions/api/emissions/v8"
	emissionsv9 "github.com/allora-network/allora-chain/x/emissions/api/emissions/v9"
	emissionsv10 "github.com/allora-network/allora-chain/x/emissions/types"
)

var (
	grpcCodec   encoding.Codec
	cosmosCodec *cosmoscodec.ProtoCodec
	registry    = codectypes.NewInterfaceRegistry()
)

func init() {
	registerFuncs := []func(codectypes.InterfaceRegistry){
		upgradetypes.RegisterInterfaces,
		banktypes.RegisterInterfaces,
		distributiontypes.RegisterInterfaces,
		slashingtypes.RegisterInterfaces,
		stakingtypes.RegisterInterfaces,
		authz.RegisterInterfaces,
		feegrant.RegisterInterfaces,
		govv1types.RegisterInterfaces,
		govv1beta1types.RegisterInterfaces,
		stdtypes.RegisterInterfaces,
		cosmossdktypes.RegisterInterfaces,
		txtypes.RegisterInterfaces,
		ibctransfertypes.RegisterInterfaces,
		ibcclient.RegisterInterfaces,
		ibcconnection.RegisterInterfaces,
		ibcchannel.RegisterInterfaces,
		ibclightclient.RegisterInterfaces,
		mintv1beta1.RegisterInterfaces,
		mintv2.RegisterInterfaces,
		mintv5.RegisterInterfaces,
		emissionsv2.RegisterInterfaces,
		emissionsv3.RegisterInterfaces,
		emissionsv4.RegisterInterfaces,
		emissionsv5.RegisterInterfaces,
		emissionsv6.RegisterInterfaces,
		emissionsv7.RegisterInterfaces,
		emissionsv8.RegisterInterfaces,
		emissionsv9.RegisterInterfaces,
		emissionsv10.RegisterInterfaces,
	}
	for _, register := range registerFuncs {
		register(registry)
	}
	registerV9TypedEvents()
	cosmosCodec = cosmoscodec.NewProtoCodec(registry)
	grpcCodec = cosmosCodec.GRPCCodec()
}

func GRPCCodec() encoding.Codec {
	return grpcCodec
}

func CosmosCodec() *cosmoscodec.ProtoCodec {
	return cosmosCodec
}

type Codec struct {
	*cosmoscodec.ProtoCodec
}

func NewCodec() *Codec {
	return &Codec{cosmosCodec}
}

func (c *Codec) IsTypedEvent(event *abcitypes.Event) bool {
	concreteGoType := proto.MessageType(event.Type)
	return concreteGoType != nil
}

func (c *Codec) ParseTypedEvent(event *abcitypes.Event) (proto.Message, error) {
	if len(event.Attributes) == 0 {
		return nil, errors.New("event has no attributes")
	}

	eventCopy := *event
	if eventCopy.Attributes[len(eventCopy.Attributes)-1].Key == "mode" {
		eventCopy.Attributes = eventCopy.Attributes[:len(eventCopy.Attributes)-1]
	}

	protoEvent, err := cosmossdktypes.ParseTypedEvent(eventCopy)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse typed event")
	}

	return protoEvent, nil
}

func (c *Codec) ParseUntypedEvent(event *abcitypes.Event) (json.RawMessage, error) {
	attrMap := make(map[string]string)
	for _, attr := range event.Attributes {
		attrMap[attr.Key] = attr.Value
	}

	attrBytes, err := json.Marshal(attrMap)
	if err != nil {
		return nil, err
	}
	return attrBytes, nil
}

func (c *Codec) ParseTx(txBytes []byte) (*txtypes.Tx, error) {
	var txMsg txtypes.Tx
	err := c.Unmarshal(txBytes, &txMsg)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to unmarshal txBytes")
	}
	return &txMsg, nil
}

func (c *Codec) ParseTxMessages(messages []*codectypes.Any) ([]proto.Message, error) {
	var parsedMessages []proto.Message
	var parseErrors []error
	for _, msgAny := range messages {
		msg, err := c.ParseTxMessage(msgAny)
		if err != nil {
			parseErrors = append(parseErrors, err)
			continue
		}
		parsedMessages = append(parsedMessages, msg)
	}
	if len(parseErrors) > 0 {
		return parsedMessages, errors.WithMessage(errors.Join(parseErrors...), "failed to parse some messages")
	}
	return parsedMessages, nil
}

func (c *Codec) ParseTxMessage(message *codectypes.Any) (proto.Message, error) {
	var msg cosmossdktypes.Msg
	err := c.UnpackAny(message, &msg)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to unpack any")
	}
	return msg, nil
}

func registerV9TypedEvents() {
	typedEvents := []struct {
		msg  proto.Message
		name string
	}{
		{&emissionsv9.EventScoresSet{}, "emissions.v9.EventScoresSet"},
		{&emissionsv9.EventRewardsSettled{}, "emissions.v9.EventRewardsSettled"},
		{&emissionsv9.EventNetworkLossSet{}, "emissions.v9.EventNetworkLossSet"},
		{&emissionsv9.EventNetworkInferences{}, "emissions.v9.EventNetworkInferences"},
		{&emissionsv9.EventOutlierResistantNetworkInferences{}, "emissions.v9.EventOutlierResistantNetworkInferences"},
		{&emissionsv9.EventValueBundle{}, "emissions.v9.EventValueBundle"},
		{&emissionsv9.EventInsertInfererPayload{}, "emissions.v9.EventInsertInfererPayload"},
		{&emissionsv9.EventInsertForecasterPayload{}, "emissions.v9.EventInsertForecasterPayload"},
		{&emissionsv9.EventCreateNewTopic{}, "emissions.v9.EventCreateNewTopic"},
		{&emissionsv9.EventTopicUpdated{}, "emissions.v9.EventTopicUpdated"},
		{&emissionsv9.EventAddStake{}, "emissions.v9.EventAddStake"},
		{&emissionsv9.EventRemoveStake{}, "emissions.v9.EventRemoveStake"},
		{&emissionsv9.EventRequestStakeRemoval{}, "emissions.v9.EventRequestStakeRemoval"},
		{&emissionsv9.EventCancelStakeRemoval{}, "emissions.v9.EventCancelStakeRemoval"},
		{&emissionsv9.EventReputerStakeUpdated{}, "emissions.v9.EventReputerStakeUpdated"},
		{&emissionsv9.EventRewardDelegateStake{}, "emissions.v9.EventRewardDelegateStake"},
		{&emissionsv9.EventInsertReputerPayload{}, "emissions.v9.EventInsertReputerPayload"},
		{&emissionsv9.EventReputerRegistered{}, "emissions.v9.EventReputerRegistered"},
		{&emissionsv9.EventWorkerRegistered{}, "emissions.v9.EventWorkerRegistered"},
		{&emissionsv9.EventNodeOwnerUpdated{}, "emissions.v9.EventNodeOwnerUpdated"},
		{&emissionsv9.EventReputerUnregistered{}, "emissions.v9.EventReputerUnregistered"},
		{&emissionsv9.EventWorkerUnregistered{}, "emissions.v9.EventWorkerUnregistered"},
		{&emissionsv9.EventFundTopic{}, "emissions.v9.EventFundTopic"},
		{&emissionsv9.EventParamsSet{}, "emissions.v9.EventParamsSet"},
		{&emissionsv9.EventWhitelistAdminAdded{}, "emissions.v9.EventWhitelistAdminAdded"},
		{&emissionsv9.EventWhitelistAdminRemoved{}, "emissions.v9.EventWhitelistAdminRemoved"},
		{&emissionsv9.EventGlobalWhitelistAdded{}, "emissions.v9.EventGlobalWhitelistAdded"},
		{&emissionsv9.EventGlobalWhitelistRemoved{}, "emissions.v9.EventGlobalWhitelistRemoved"},
		{&emissionsv9.EventGlobalWorkerWhitelistAdded{}, "emissions.v9.EventGlobalWorkerWhitelistAdded"},
		{&emissionsv9.EventGlobalWorkerWhitelistRemoved{}, "emissions.v9.EventGlobalWorkerWhitelistRemoved"},
		{&emissionsv9.EventGlobalReputerWhitelistAdded{}, "emissions.v9.EventGlobalReputerWhitelistAdded"},
		{&emissionsv9.EventGlobalReputerWhitelistRemoved{}, "emissions.v9.EventGlobalReputerWhitelistRemoved"},
		{&emissionsv9.EventGlobalAdminWhitelistAdded{}, "emissions.v9.EventGlobalAdminWhitelistAdded"},
		{&emissionsv9.EventGlobalAdminWhitelistRemoved{}, "emissions.v9.EventGlobalAdminWhitelistRemoved"},
		{&emissionsv9.EventTopicWorkerWhitelistEnabled{}, "emissions.v9.EventTopicWorkerWhitelistEnabled"},
		{&emissionsv9.EventTopicWorkerWhitelistDisabled{}, "emissions.v9.EventTopicWorkerWhitelistDisabled"},
		{&emissionsv9.EventTopicReputerWhitelistEnabled{}, "emissions.v9.EventTopicReputerWhitelistEnabled"},
		{&emissionsv9.EventTopicReputerWhitelistDisabled{}, "emissions.v9.EventTopicReputerWhitelistDisabled"},
		{&emissionsv9.EventTopicCreatorWhitelistAdded{}, "emissions.v9.EventTopicCreatorWhitelistAdded"},
		{&emissionsv9.EventTopicCreatorWhitelistRemoved{}, "emissions.v9.EventTopicCreatorWhitelistRemoved"},
		{&emissionsv9.EventTopicWorkerWhitelistAdded{}, "emissions.v9.EventTopicWorkerWhitelistAdded"},
		{&emissionsv9.EventTopicWorkerWhitelistRemoved{}, "emissions.v9.EventTopicWorkerWhitelistRemoved"},
		{&emissionsv9.EventTopicReputerWhitelistAdded{}, "emissions.v9.EventTopicReputerWhitelistAdded"},
		{&emissionsv9.EventTopicReputerWhitelistRemoved{}, "emissions.v9.EventTopicReputerWhitelistRemoved"},
		{&emissionsv9.EventForecastTaskScoreSet{}, "emissions.v9.EventForecastTaskScoreSet"},
		{&emissionsv9.EventWorkerLastCommitSet{}, "emissions.v9.EventWorkerLastCommitSet"},
		{&emissionsv9.EventReputerLastCommitSet{}, "emissions.v9.EventReputerLastCommitSet"},
		{&emissionsv9.EventTopicRewardsSet{}, "emissions.v9.EventTopicRewardsSet"},
		{&emissionsv9.EventEMAScoresSet{}, "emissions.v9.EventEMAScoresSet"},
		{&emissionsv9.EventListeningCoefficientsSet{}, "emissions.v9.EventListeningCoefficientsSet"},
		{&emissionsv9.EventInfererNetworkRegretSet{}, "emissions.v9.EventInfererNetworkRegretSet"},
		{&emissionsv9.EventForecasterNetworkRegretSet{}, "emissions.v9.EventForecasterNetworkRegretSet"},
		{&emissionsv9.EventNaiveInfererNetworkRegretSet{}, "emissions.v9.EventNaiveInfererNetworkRegretSet"},
		{&emissionsv9.EventTopicInitialRegretSet{}, "emissions.v9.EventTopicInitialRegretSet"},
		{&emissionsv9.EventTopicInitialEmaScoreSet{}, "emissions.v9.EventTopicInitialEmaScoreSet"},
		{&emissionsv9.EventRegretStdNormSet{}, "emissions.v9.EventRegretStdNormSet"},
		{&emissionsv9.EventInfererWeightsSet{}, "emissions.v9.EventInfererWeightsSet"},
		{&emissionsv9.EventForecasterWeightsSet{}, "emissions.v9.EventForecasterWeightsSet"},
		{&emissionsv9.EventPreviousPercentageRewardToStakedReputersSet{}, "emissions.v9.EventPreviousPercentageRewardToStakedReputersSet"},
		{&emissionsv9.EventPruneRecords{}, "emissions.v9.EventPruneRecords"},
		{&emissionsv9.EventDelegateRewardShareUpdated{}, "emissions.v9.EventDelegateRewardShareUpdated"},
		{&emissionsv9.EventDelegateRewardDistributed{}, "emissions.v9.EventDelegateRewardDistributed"},
		{&emissionsv9.EventActiveReputersSet{}, "emissions.v9.EventActiveReputersSet"},
		{&emissionsv9.EventActiveInferersSet{}, "emissions.v9.EventActiveInferersSet"},
		{&emissionsv9.EventActiveForecastersSet{}, "emissions.v9.EventActiveForecastersSet"},
		{&emissionsv9.EventTopicStatusChanged{}, "emissions.v9.EventTopicStatusChanged"},
		{&emissionsv9.EventNetworkInferenceInfererWeightsSet{}, "emissions.v9.EventNetworkInferenceInfererWeightsSet"},
		{&emissionsv9.EventNetworkInferenceForecasterWeightsSet{}, "emissions.v9.EventNetworkInferenceForecasterWeightsSet"},
		{&emissionsv9.EventNetworkInferenceInfererRegretsUsedSet{}, "emissions.v9.EventNetworkInferenceInfererRegretsUsedSet"},
		{&emissionsv9.EventNetworkInferenceForecasterRegretsUsedSet{}, "emissions.v9.EventNetworkInferenceForecasterRegretsUsedSet"},
		{&emissionsv9.EventTopicWeightUpdated{}, "emissions.v9.EventTopicWeightUpdated"},
		{&emissionsv9.EventTopicFeeRevenueDripped{}, "emissions.v9.EventTopicFeeRevenueDripped"},
		{&emissionsv9.EventWorkerSubmissionWindowOpened{}, "emissions.v9.EventWorkerSubmissionWindowOpened"},
		{&emissionsv9.EventWorkerSubmissionWindowClosed{}, "emissions.v9.EventWorkerSubmissionWindowClosed"},
		{&emissionsv9.EventReputerSubmissionWindowOpened{}, "emissions.v9.EventReputerSubmissionWindowOpened"},
		{&emissionsv9.EventReputerSubmissionWindowClosed{}, "emissions.v9.EventReputerSubmissionWindowClosed"},
	}

	for _, event := range typedEvents {
		proto.RegisterType(event.msg, event.name)
	}
}
