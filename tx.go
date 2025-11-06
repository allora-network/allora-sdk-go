package allora

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TxParams contains all parameters needed to build and sign a transaction
type TxParams struct {
	// Chain identification
	ChainID string

	// Account information
	AccountNumber uint64
	Sequence      uint64

	// Gas and fees
	GasLimit  uint64
	FeeAmount sdk.Coins

	// Optional fields
	Memo          string
	TimeoutHeight uint64
}

// DefaultTxParams returns TxParams with sensible defaults for Allora Network
// Note: ChainID, AccountNumber, and Sequence must still be set by the caller
func DefaultTxParams() *TxParams {
	return &TxParams{
		GasLimit:  200000, // Standard gas limit for simple transfers
		FeeAmount: sdk.NewCoins(sdk.NewInt64Coin("uallo", 5000)), // Default fee
	}
}

// Validate checks that all required parameters are set
func (p *TxParams) Validate() error {
	if p.ChainID == "" {
		return fmt.Errorf("chain ID is required")
	}
	if p.GasLimit == 0 {
		return fmt.Errorf("gas limit must be greater than 0")
	}
	if p.FeeAmount.Empty() {
		return fmt.Errorf("fee amount is required")
	}
	return nil
}

// CreateUnsignedSendTx creates an unsigned send transaction
//
// This function creates a bank MsgSend transaction in unsigned form, which can be
// stored in a database, Kafka, or other persistent storage before signing.
//
// Parameters:
//   - fromAddr: The sender's Allora address (allo...)
//   - toAddr: The recipient's Allora address (allo...)
//   - amount: The amount to send (e.g., sdk.NewCoins(sdk.NewInt64Coin("uallo", 1000000)))
//   - params: Transaction parameters including chain ID, account info, gas, and fees
//
// Returns:
//   - []byte: The unsigned transaction bytes (can be stored and signed later)
//   - error: Any error that occurred during transaction creation
//
// Example:
//
//	fromAddr, _ := sdk.AccAddressFromBech32("allo1...")
//	toAddr, _ := sdk.AccAddressFromBech32("allo1...")
//	amount := sdk.NewCoins(sdk.NewInt64Coin("uallo", 1000000))
//
//	params := &allora.TxParams{
//	    ChainID:       "allora-mainnet-1",
//	    AccountNumber: 123,
//	    Sequence:      5,
//	    GasLimit:      200000,
//	    FeeAmount:     sdk.NewCoins(sdk.NewInt64Coin("uallo", 5000)),
//	    Memo:          "deposit-12345",
//	}
//
//	unsignedTx, err := allora.CreateUnsignedSendTx(fromAddr, toAddr, amount, params)
func CreateUnsignedSendTx(
	fromAddr sdk.AccAddress,
	toAddr sdk.AccAddress,
	amount sdk.Coins,
	params *TxParams,
) ([]byte, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("invalid transaction parameters: %w", err)
	}

	if fromAddr.Empty() {
		return nil, fmt.Errorf("from address is required")
	}
	if toAddr.Empty() {
		return nil, fmt.Errorf("to address is required")
	}
	if amount.Empty() {
		return nil, fmt.Errorf("amount is required")
	}
	if err := amount.Validate(); err != nil {
		return nil, fmt.Errorf("invalid amount: %w", err)
	}

	builder := newTxBuilder()
	return builder.buildUnsignedSendTx(fromAddr, toAddr, amount, params)
}

// SignTransaction signs a previously created unsigned transaction
//
// This function takes an unsigned transaction (created with CreateUnsignedSendTx)
// and signs it with the provided wallet. This enables a two-phase workflow where
// transactions can be created, stored, and signed at a later time.
//
// Parameters:
//   - unsignedTx: The unsigned transaction bytes from CreateUnsignedSendTx
//   - wallet: The wallet to sign the transaction with
//   - params: The same TxParams used to create the unsigned transaction
//
// Returns:
//   - []byte: The signed transaction bytes (ready to broadcast)
//   - error: Any error that occurred during signing
//
// Example:
//
//	// Create wallet
//	wallet, _ := allora.GenerateWallet()
//
//	// Sign the transaction
//	signedTx, err := allora.SignTransaction(unsignedTx, wallet, params)
func SignTransaction(
	unsignedTx []byte,
	wallet *Wallet,
	params *TxParams,
) ([]byte, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("invalid transaction parameters: %w", err)
	}

	if len(unsignedTx) == 0 {
		return nil, fmt.Errorf("unsigned transaction is empty")
	}

	if wallet == nil {
		return nil, fmt.Errorf("wallet is required")
	}

	builder := newTxBuilder()
	return builder.signTx(unsignedTx, wallet.PrivKey, params)
}

// CreateSignedSendTx is a convenience function that creates and signs a send transaction in one step
//
// This function combines CreateUnsignedSendTx and SignTransaction for cases where
// you want to create and sign a transaction immediately without storing the unsigned version.
//
// Parameters:
//   - fromAddr: The sender's Allora address (allo...)
//   - toAddr: The recipient's Allora address (allo...)
//   - amount: The amount to send
//   - wallet: The wallet to sign with (must match fromAddr)
//   - params: Transaction parameters
//
// Returns:
//   - []byte: The signed transaction bytes (ready to broadcast)
//   - error: Any error that occurred
//
// Example:
//
//	wallet, _ := allora.GenerateWallet()
//	fromAddr := wallet.Address
//	toAddr, _ := sdk.AccAddressFromBech32("allo1...")
//	amount := sdk.NewCoins(sdk.NewInt64Coin("uallo", 1000000))
//
//	params := &allora.TxParams{
//	    ChainID:       "allora-mainnet-1",
//	    AccountNumber: 123,
//	    Sequence:      5,
//	    GasLimit:      200000,
//	    FeeAmount:     sdk.NewCoins(sdk.NewInt64Coin("uallo", 5000)),
//	}
//
//	signedTx, err := allora.CreateSignedSendTx(fromAddr, toAddr, amount, wallet, params)
func CreateSignedSendTx(
	fromAddr sdk.AccAddress,
	toAddr sdk.AccAddress,
	amount sdk.Coins,
	wallet *Wallet,
	params *TxParams,
) ([]byte, error) {
	// Verify wallet address matches from address
	if !wallet.Address.Equals(fromAddr) {
		return nil, fmt.Errorf("wallet address does not match from address")
	}

	// Create unsigned transaction
	unsignedTx, err := CreateUnsignedSendTx(fromAddr, toAddr, amount, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create unsigned transaction: %w", err)
	}

	// Sign the transaction
	signedTx, err := SignTransaction(unsignedTx, wallet, params)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	return signedTx, nil
}

// ParseTxBytes parses transaction bytes and returns the decoded transaction
// This is useful for inspecting transaction contents before broadcasting
func ParseTxBytes(txBytes []byte) (*sdk.Tx, error) {
	builder := newTxBuilder()
	decoder := builder.txConfig.TxDecoder()
	tx, err := decoder(txBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decode transaction: %w", err)
	}
	return &tx, nil
}
