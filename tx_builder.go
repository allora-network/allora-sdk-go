package allora

import (
	"bytes"
	"context"
	"fmt"

	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"

	alloracodec "github.com/allora-network/allora-sdk-go/codec"
)

// txBuilder wraps the Cosmos SDK transaction building functionality
type txBuilder struct {
	txConfig sdkclient.TxConfig
	codec    codec.Codec
}

// newTxBuilder creates a new transaction builder with the Allora codec
func newTxBuilder() *txBuilder {
	// Use the Allora codec which has all interfaces registered
	txConfig := authtx.NewTxConfig(alloracodec.CosmosCodec(), authtx.DefaultSignModes)

	return &txBuilder{
		txConfig: txConfig,
		codec:    alloracodec.CosmosCodec(),
	}
}

// buildUnsignedTx builds an unsigned transaction from arbitrary sdk.Msg values,
// encoding it with the Allora codec. Callers are responsible for the messages'
// internal validity; this method validates only that at least one message is
// present.
func (b *txBuilder) buildUnsignedTx(
	msgs []sdk.Msg,
	params *TxParams,
) ([]byte, error) {
	if len(msgs) == 0 {
		return nil, fmt.Errorf("at least one message is required")
	}

	// Create transaction builder
	txBuilder := b.txConfig.NewTxBuilder()

	// Set the messages
	if err := txBuilder.SetMsgs(msgs...); err != nil {
		return nil, fmt.Errorf("failed to set messages: %w", err)
	}

	// Set transaction parameters
	txBuilder.SetGasLimit(params.GasLimit)
	txBuilder.SetFeeAmount(params.FeeAmount)
	txBuilder.SetMemo(params.Memo)

	// Set the fee granter so a master/subsidy wallet pays the gas (feegrant). Empty
	// means the signer pays its own fees.
	if !params.FeeGranter.Empty() {
		txBuilder.SetFeeGranter(params.FeeGranter)
	}

	if params.TimeoutHeight > 0 {
		txBuilder.SetTimeoutHeight(params.TimeoutHeight)
	}

	// Encode the unsigned transaction
	txBytes, err := b.txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, fmt.Errorf("failed to encode transaction: %w", err)
	}

	return txBytes, nil
}

// verifySignerForTx checks that pubKey's account address is among the required
// signers of every message in tx. It resolves signers via the codec's
// GetMsgV1Signers, which returns each signer as raw address bytes, and compares
// them to pubKey.Address() without bech32 round-tripping. It returns an error if
// any message has no signers or does not include this signer, mirroring the
// on-chain "signature verification failed" failure but surfacing it before
// broadcast.
func (b *txBuilder) verifySignerForTx(tx sdk.Tx, pubKey cryptotypes.PubKey) error {
	signerBytes := pubKey.Address().Bytes()

	for i, msg := range tx.GetMsgs() {
		msgSigners, _, err := b.codec.GetMsgV1Signers(msg)
		if err != nil {
			return fmt.Errorf("failed to resolve signers for message %d: %w", i, err)
		}
		if len(msgSigners) == 0 {
			return fmt.Errorf("message %d has no signers", i)
		}

		found := false
		for _, s := range msgSigners {
			if bytes.Equal(s, signerBytes) {
				found = true
				break
			}
		}
		if !found {
			signerAddr := sdk.AccAddress(signerBytes).String()
			return fmt.Errorf("signer address %s is not a required signer of message %d", signerAddr, i)
		}
	}

	return nil
}

// signTx signs a transaction with the provided signer (a local key or a remote signer).
// ctx is honored by signers that perform I/O (e.g. RemoteSigner) and during sign-bytes
// assembly.
func (b *txBuilder) signTx(
	ctx context.Context,
	txBytes []byte,
	signer Signer,
	params *TxParams,
) ([]byte, error) {
	// Decode the unsigned transaction
	txDecoder := b.txConfig.TxDecoder()
	decodedTx, err := txDecoder(txBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decode transaction: %w", err)
	}

	// Convert to TxBuilder for signing
	txBuilder, err := b.txConfig.WrapTxBuilder(decodedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap transaction: %w", err)
	}

	// Get public key
	pubKey := signer.PubKey()
	// Signer is a public extension point (any type satisfying the interface can be passed to
	// SignTransactionWith), so a buggy custom signer whose PubKey() returns nil — including a
	// typed-nil like (*secp256k1.PubKey)(nil), non-nil at the interface level but nil at the
	// concrete level — would panic at pubKey.Address() and crash the worker. isNilPubKey catches
	// both the literal-nil and typed-nil cases (mirroring isNilSigner one layer up); fail with a
	// clear error instead.
	if isNilPubKey(pubKey) {
		return nil, fmt.Errorf("signer returned nil public key")
	}
	signerAddr := sdk.AccAddress(pubKey.Address()).String()

	// Guard against signing a transaction whose message signer is not this signer; such a
	// tx is rejected on-chain ("signature verification failed") far from the cause. The
	// check is message-type-agnostic: it resolves each message's required signers via the
	// codec (GetMsgV1Signers) and asserts this signer is among them for every message, so
	// it covers bank, emissions, feegrant, and any other registered message type.
	if err := b.verifySignerForTx(decodedTx, pubKey); err != nil {
		return nil, err
	}

	// Convert API sign mode to internal
	apiSignMode := b.txConfig.SignModeHandler().DefaultMode()
	signMode, err := authsigning.APISignModeToInternal(apiSignMode)
	if err != nil {
		return nil, fmt.Errorf("failed to convert sign mode: %w", err)
	}

	// Create signature data with empty signature (first pass)
	sigData := &signingtypes.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}

	sig := signingtypes.SignatureV2{
		PubKey:   pubKey,
		Data:     sigData,
		Sequence: params.Sequence,
	}

	// Set the signature (with nil signature data) to get proper sign bytes
	if err := txBuilder.SetSignatures(sig); err != nil {
		return nil, fmt.Errorf("failed to set signatures: %w", err)
	}

	// Create signer data
	signerData := authsigning.SignerData{
		ChainID:       params.ChainID,
		AccountNumber: params.AccountNumber,
		Sequence:      params.Sequence,
		Address:       signerAddr,
		PubKey:        pubKey,
	}

	// Get the bytes to sign
	bytesToSign, err := authsigning.GetSignBytesAdapter(
		ctx,
		b.txConfig.SignModeHandler(),
		signMode,
		signerData,
		txBuilder.GetTx(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get sign bytes: %w", err)
	}

	// Sign the bytes
	signature, err := signWithContext(ctx, signer, bytesToSign)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Update signature data with actual signature
	sigData = &signingtypes.SingleSignatureData{
		SignMode:  signMode,
		Signature: signature,
	}

	sig = signingtypes.SignatureV2{
		PubKey:   pubKey,
		Data:     sigData,
		Sequence: params.Sequence,
	}

	// Set the actual signature
	if err := txBuilder.SetSignatures(sig); err != nil {
		return nil, fmt.Errorf("failed to set final signatures: %w", err)
	}

	// Encode the signed transaction
	signedTxBytes, err := b.txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, fmt.Errorf("failed to encode signed transaction: %w", err)
	}

	return signedTxBytes, nil
}
