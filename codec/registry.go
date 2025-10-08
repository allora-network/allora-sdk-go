package codec

import (
	"encoding/json"
	"sync"

	"github.com/brynbellomy/go-utils/errors"
	"google.golang.org/grpc/encoding"

	feegrant "cosmossdk.io/x/feegrant"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	cosmoscodec "github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	stdtypes "github.com/cosmos/cosmos-sdk/std"
	cosmossdktypes "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gogo/protobuf/proto"

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
	emissionsv9 "github.com/allora-network/allora-chain/x/emissions/types"
)

var (
	once        sync.Once
	grpcCodec   encoding.Codec
	cosmosCodec *cosmoscodec.ProtoCodec
	registry    = codectypes.NewInterfaceRegistry()
)

func init() {
	registerFuncs := []func(codectypes.InterfaceRegistry){
		banktypes.RegisterInterfaces,
		distributiontypes.RegisterInterfaces,
		slashingtypes.RegisterInterfaces,
		stakingtypes.RegisterInterfaces,
		feegrant.RegisterInterfaces,
		stdtypes.RegisterInterfaces,
		cosmossdktypes.RegisterInterfaces,
		txtypes.RegisterInterfaces,
		bank.RegisterInterfaces,
		distribution.RegisterInterfaces,
		slashing.RegisterInterfaces,
		staking.RegisterInterfaces,
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
	}
	for _, register := range registerFuncs {
		register(registry)
	}
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

func (c *Codec) IsTypedEvent(event *abcitypes.Event) bool {
	concreteGoType := proto.MessageType(event.Type)
	return concreteGoType != nil
}

func (c *Codec) ParseTypedEvent(event *abcitypes.Event) (proto.Message, error) {
	if len(event.Attributes) == 0 {
		return nil, errors.New("event has no attributes")
	}

	eventCopy := *event
	// Remove the last attribute if the key is "mode"
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
