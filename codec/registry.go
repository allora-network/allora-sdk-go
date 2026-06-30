package codec

import (
	"encoding/json"
	"fmt"

	txsigning "cosmossdk.io/x/tx/signing"

	coreaddress "cosmossdk.io/core/address"
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
	emissionsv9 "github.com/allora-network/allora-chain/x/emissions/types"
)

var (
	grpcCodec   encoding.Codec
	cosmosCodec *cosmoscodec.ProtoCodec
	registry    codectypes.InterfaceRegistry
)

func init() {
	// Build the interface registry with an address codec that delegates to the
	// SDK's global config at call time. This is required for the registry's
	// GetMsgV1Signers to resolve message signers (the tx builder's
	// message-agnostic sender guard depends on it); a bare
	// NewInterfaceRegistry() has no address codec and fails at signer-resolution
	// time. Using a lazy codec avoids init-order coupling with the allora wallet
	// package, which sets the bech32 prefix and seals the SDK config.
	var err error
	registry, err = codectypes.NewInterfaceRegistryWithOptions(codectypes.InterfaceRegistryOptions{
		SigningOptions: txsigning.Options{
			AddressCodec:          &lazyAddressCodec{},
			ValidatorAddressCodec: &lazyValidatorAddressCodec{},
		},
		// Use gogoproto's hybrid resolver (the same default NewInterfaceRegistry uses)
		// rather than protoregistry.GlobalFiles: the cosmos modules register their
		// descriptors through gogoproto, so GlobalFiles alone cannot resolve signers for
		// types like feegrant.MsgRevokeAllowance, while the hybrid resolver can.
		ProtoFiles: proto.HybridResolver,
	})
	if err != nil {
		panic(fmt.Errorf("failed to create interface registry: %w", err))
	}
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

// lazyAddressCodec delegates StringToBytes / BytesToString to the SDK global
// config's bech32 account prefix at call time. This avoids init-order coupling
// between the codec package and the allora wallet package (which sets the
// bech32 prefix and seals the SDK config in its own init).
type lazyAddressCodec struct{}

var _ coreaddress.Codec = (*lazyAddressCodec)(nil)

func (lazyAddressCodec) StringToBytes(text string) ([]byte, error) {
	return cosmossdktypes.AccAddressFromBech32(text)
}

func (lazyAddressCodec) BytesToString(bz []byte) (string, error) {
	return cosmossdktypes.AccAddress(bz).String(), nil
}

// lazyValidatorAddressCodec is the validator-address counterpart of
// lazyAddressCodec: it resolves validator bech32 (alloravaloper...) signers via
// the SDK global config's validator prefix. Using the account codec here would
// mis-resolve staking messages whose signer is a validator address.
type lazyValidatorAddressCodec struct{}

var _ coreaddress.Codec = (*lazyValidatorAddressCodec)(nil)

func (lazyValidatorAddressCodec) StringToBytes(text string) ([]byte, error) {
	return cosmossdktypes.ValAddressFromBech32(text)
}

func (lazyValidatorAddressCodec) BytesToString(bz []byte) (string, error) {
	return cosmossdktypes.ValAddress(bz).String(), nil
}
