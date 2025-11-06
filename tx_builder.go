package allora

import (
	"context"
	"fmt"

	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

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

// buildUnsignedSendTx creates an unsigned send transaction
func (b *txBuilder) buildUnsignedSendTx(
	fromAddr sdk.AccAddress,
	toAddr sdk.AccAddress,
	amount sdk.Coins,
	params *TxParams,
) ([]byte, error) {
	// Create the MsgSend
	msg := banktypes.NewMsgSend(fromAddr, toAddr, amount)

	// Create transaction builder
	txBuilder := b.txConfig.NewTxBuilder()

	// Set the message
	if err := txBuilder.SetMsgs(msg); err != nil {
		return nil, fmt.Errorf("failed to set messages: %w", err)
	}

	// Set transaction parameters
	txBuilder.SetGasLimit(params.GasLimit)
	txBuilder.SetFeeAmount(params.FeeAmount)
	txBuilder.SetMemo(params.Memo)

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

// signTx signs a transaction with the provided private key
func (b *txBuilder) signTx(
	txBytes []byte,
	privKey cryptotypes.PrivKey,
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
	pubKey := privKey.PubKey()

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
		Address:       sdk.AccAddress(pubKey.Address()).String(),
		PubKey:        pubKey,
	}

	// Get the bytes to sign
	bytesToSign, err := authsigning.GetSignBytesAdapter(
		context.Background(),
		b.txConfig.SignModeHandler(),
		signMode,
		signerData,
		txBuilder.GetTx(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get sign bytes: %w", err)
	}

	// Sign the bytes
	signature, err := privKey.Sign(bytesToSign)
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
