package allora

import (
	"context"
	"testing"

	"cosmossdk.io/x/feegrant"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
)

// defaultGeneralTxParams returns TxParams suitable for the generalized-builder tests.
func defaultGeneralTxParams() *TxParams {
	return &TxParams{
		ChainID:       "allora-testnet-1",
		AccountNumber: 11,
		Sequence:      4,
		GasLimit:      200000,
		FeeAmount:     sdk.NewCoins(sdk.NewInt64Coin("uallo", 5000)),
	}
}

// TestCreateUnsignedTx_RejectsEmptyMsgs pins that the general entrypoint rejects a
// nil or empty message slice rather than producing a meaningless empty transaction.
func TestCreateUnsignedTx_RejectsEmptyMsgs(t *testing.T) {
	params := defaultGeneralTxParams()

	_, err := CreateUnsignedTx(nil, params)
	require.Error(t, err)
	require.ErrorContains(t, err, "at least one message is required")

	_, err = CreateUnsignedTx([]sdk.Msg{}, params)
	require.Error(t, err)
	require.ErrorContains(t, err, "at least one message is required")
}

// TestCreateUnsignedTx_MultiMessageRoundTrip pins that the general builder encodes
// multiple messages and that they decode back in order.
func TestCreateUnsignedTx_MultiMessageRoundTrip(t *testing.T) {
	wallet, err := GenerateWallet()
	require.NoError(t, err)
	recipient, err := GenerateWallet()
	require.NoError(t, err)

	amount := sdk.NewCoins(sdk.NewInt64Coin("uallo", 1000))
	msgs := []sdk.Msg{
		banktypes.NewMsgSend(wallet.Address, recipient.Address, amount),
		banktypes.NewMsgSend(wallet.Address, recipient.Address, amount),
	}

	unsigned, err := CreateUnsignedTx(msgs, defaultGeneralTxParams())
	require.NoError(t, err)

	decoded, err := ParseTxBytes(unsigned)
	require.NoError(t, err)
	require.Len(t, (*decoded).GetMsgs(), 2)
}

// TestSignTransactionWith_NonBankMsg_GuardMatches pins that the msg-agnostic sender
// guard accepts a non-bank message (feegrant MsgRevokeAllowance) whose required
// signer is the signing wallet, exercising the codec signer-resolution path for a
// message type the old MsgSend-only guard could not check.
func TestSignTransactionWith_NonBankMsg_GuardMatches(t *testing.T) {
	granter, err := GenerateWallet()
	require.NoError(t, err)
	grantee, err := GenerateWallet()
	require.NoError(t, err)

	msg := &feegrant.MsgRevokeAllowance{
		Granter: granter.Address.String(),
		Grantee: grantee.Address.String(),
	}

	params := defaultGeneralTxParams()
	unsigned, err := CreateUnsignedTx([]sdk.Msg{msg}, params)
	require.NoError(t, err)

	// Signing with the granter's key (the message's required signer) must succeed.
	signed, err := SignTransactionWith(context.Background(), unsigned, granter.PrivKey, params)
	require.NoError(t, err)
	require.NotEmpty(t, signed)
}

// TestSignTransactionWith_NonBankMsg_GuardRejectsMismatch pins that the guard rejects
// signing a non-bank message with a key that is not the message's required signer.
func TestSignTransactionWith_NonBankMsg_GuardRejectsMismatch(t *testing.T) {
	granter, err := GenerateWallet()
	require.NoError(t, err)
	grantee, err := GenerateWallet()
	require.NoError(t, err)
	stranger, err := GenerateWallet()
	require.NoError(t, err)

	msg := &feegrant.MsgRevokeAllowance{
		Granter: granter.Address.String(),
		Grantee: grantee.Address.String(),
	}

	params := defaultGeneralTxParams()
	unsigned, err := CreateUnsignedTx([]sdk.Msg{msg}, params)
	require.NoError(t, err)

	// stranger is not the granter, so the guard must reject before broadcast.
	_, err = SignTransactionWith(context.Background(), unsigned, stranger.PrivKey, params)
	require.Error(t, err)
	require.ErrorContains(t, err, "is not a required signer")
}
