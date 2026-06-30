package allora

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CreateUnsignedTx builds an unsigned transaction from the given messages and
// parameters. It is the generic entry point for any message type — a bank send,
// a feegrant grant, an emissions payload, or any combination — and replaces
// the single-message CreateUnsignedSendTx when callers need more flexibility.
//
// The returned bytes are an encoded, unsigned Cosmos transaction. Pass them to
// SignTransactionWith to produce a signed tx ready for broadcast, or to a
// txsend.TxBroadcaster for gas estimation before signing.
//
// At least one message is required. An empty slice is rejected rather than
// producing a tx the chain would reject anyway; the SDK's TxBuilder panics on a
// nil or empty message set, so this guard keeps the failure surface within the
// SDK's error handling.
func CreateUnsignedTx(msgs []sdk.Msg, params *TxParams) ([]byte, error) {
	if len(msgs) == 0 {
		return nil, fmt.Errorf("at least one message is required")
	}
	if params == nil {
		return nil, fmt.Errorf("tx params are required")
	}
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("invalid transaction parameters: %w", err)
	}

	builder := newTxBuilder()
	return builder.buildUnsignedTx(msgs, params)
}