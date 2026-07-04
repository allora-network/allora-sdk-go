package allora

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog/log"
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

	// Optional fields. Appended at the end so adding new optional fields does not
	// shift the position of existing ones, which would break any caller using
	// positional (unkeyed) struct literals.
	Memo          string
	TimeoutHeight uint64

	// FeeGranter, when set, is the address that pays the transaction fee via an
	// on-chain feegrant (e.g. a master subsidy wallet). Empty means the signer pays.
	FeeGranter sdk.AccAddress
}

// DefaultTxParams returns TxParams with sensible defaults for Allora Network
// Note: ChainID, AccountNumber, and Sequence must still be set by the caller
func DefaultTxParams() *TxParams {
	return &TxParams{
		GasLimit:  200000, // Standard gas limit for simple transfers
		FeeAmount: sdk.NewCoins(sdk.NewInt64Coin("uallo", 5000)), // Default fee
	}
}

const (
	// feeGranterEnvVar is the canonical fee-granter (master subsidy wallet) env var, shared
	// across the Allora SDKs (allora-sdk-py / allora-sdk-ts).
	feeGranterEnvVar = "FORGE_MASTER_GRANTER_ADDRESS"
	// legacyFeeGranterEnvVar is the pre-rename name. allora-sdk-py still accepts it for one
	// release with a deprecation warning, so FeeGranterFromEnv honors it too (with a warning)
	// to avoid silently dropping the granter on a Python->Go migration.
	// @@TODO: remove the FEE_GRANTER fallback in the next minor release.
	legacyFeeGranterEnvVar = "FEE_GRANTER"
)

// feeGranterDeprecationOnce limits the FEE_GRANTER deprecation notice to one emission per
// process, mirroring Python's once-per-location DeprecationWarning (FeeGranterFromEnv may be
// called once per transaction).
var feeGranterDeprecationOnce sync.Once

// FeeGranterFromEnv reads the canonical FORGE_MASTER_GRANTER_ADDRESS environment variable and
// parses it into an AccAddress suitable for TxParams.FeeGranter. This is the fee-granter (master
// subsidy wallet) env var shared across the Allora SDKs (allora-sdk-py / allora-sdk-ts), giving Go
// callers a single discovery path for the same value. It returns (nil, nil) when neither variable
// is set or blank, in which case the signing wallet pays its own fees.
//
// For parity with allora-sdk-py, the pre-rename FEE_GRANTER is still accepted as a fallback for one
// release (with a one-time deprecation warning): a deployment migrating from Python that still sets
// FEE_GRANTER keeps working instead of silently falling through to no-granter. Rename it to
// FORGE_MASTER_GRANTER_ADDRESS. @@TODO: drop the FEE_GRANTER fallback in the next minor release.
func FeeGranterFromEnv() (sdk.AccAddress, error) {
	envVar := feeGranterEnvVar
	addr := strings.TrimSpace(os.Getenv(envVar))
	if addr == "" {
		// Fall back to the deprecated name so a Python->Go migration does not silently lose the
		// granter (allora-sdk-py accepts FEE_GRANTER for one release with a deprecation warning).
		if legacy := strings.TrimSpace(os.Getenv(legacyFeeGranterEnvVar)); legacy != "" {
			envVar, addr = legacyFeeGranterEnvVar, legacy
			feeGranterDeprecationOnce.Do(func() {
				log.Warn().Msgf("%s is deprecated; rename it to %s", legacyFeeGranterEnvVar, feeGranterEnvVar)
			})
		}
	}
	if addr == "" {
		return nil, nil
	}
	granter, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid %s %q: %w", envVar, addr, err)
	}
	return granter, nil
}

// Validate checks that all required parameters are set. The receiver must be non-nil; the
// SDK's call sites guard for a nil *TxParams before invoking Validate, keeping the nil check
// at the call site rather than masking it inside a method on a nil pointer.
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
	// FeeGranter, when set, is encoded verbatim into the tx AuthInfo; a wrong-length []byte
	// serializes without error and only fails on-chain (feegrant not found) far from here,
	// after wasting a broadcast round-trip. Cosmos account addresses are 20 bytes, so reject
	// any other non-empty length up front.
	if len(p.FeeGranter) != 0 && len(p.FeeGranter) != 20 {
		return fmt.Errorf("fee granter address must be 20 bytes when set, got %d", len(p.FeeGranter))
	}
	return nil
}

// validateForSigning checks only the params fields consumed at sign time. Unlike Validate, it does
// not require GasLimit or FeeAmount: SignTransactionWith signs over the gas/fee/memo/granter values
// already encoded in the unsigned tx and reads only ChainID, AccountNumber, and Sequence from
// params (see SignTransactionWith). Requiring gas/fee here would break the documented two-phase
// flow — build + persist the unsigned tx, then sign later with minimal sign-time params — by
// rejecting params that legitimately omit the already-encoded fee/gas fields.
func (p *TxParams) validateForSigning() error {
	if p.ChainID == "" {
		return fmt.Errorf("chain ID is required")
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
	if params == nil {
		return nil, fmt.Errorf("tx params are required")
	}
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
	if wallet == nil {
		return nil, fmt.Errorf("wallet is required")
	}
	// wallet.PrivKey is a cryptotypes.PrivKey, which satisfies Signer. Local signing is
	// CPU-bound, so a background context is sufficient.
	return SignTransactionWith(context.Background(), unsignedTx, wallet.PrivKey, params)
}

// SignTransactionWith signs an unsigned transaction with any Signer. This is the
// general form of SignTransaction: pass a local wallet's key for self-managed signing,
// or a *RemoteSigner to delegate signing to the Forge backend (Privy-managed wallet).
//
// ctx is propagated to signers that perform I/O: a *RemoteSigner issues an HTTP call to
// the backend, so a deadline or cancellation on ctx bounds that call. Local-key signers
// ignore it.
//
// Only ChainID, AccountNumber, and Sequence from params are consumed here (they populate the
// SignerData). GasLimit, FeeAmount, FeeGranter, Memo, and TimeoutHeight are NOT re-read from
// params at sign time — the signature is computed over the values already encoded in
// unsignedTx — so passing params whose fee/gas/memo/timeout-height differ from those used to
// build unsignedTx does NOT change what is signed and is silently ignored. Consequently only
// ChainID is required at sign time (validated non-empty); GasLimit and FeeAmount need not be
// re-supplied, so a two-phase flow can sign later with a minimal params carrying just ChainID,
// AccountNumber, and Sequence. Passing the same params you built the unsigned tx with remains the
// simplest correct choice, since AccountNumber and Sequence must still match the encoded tx.
//
// Parameters:
//   - ctx: Context for cancellation/deadlines, honored by I/O-backed signers
//   - unsignedTx: The unsigned transaction bytes from CreateUnsignedSendTx
//   - signer: The signer (local key or remote signer)
//   - params: The same TxParams used to create the unsigned transaction. Only ChainID,
//     AccountNumber, and Sequence are read at sign time; the fee/gas/memo/timeout-height
//     fields come from the encoded tx, not from params (see the note above).
//
// Example:
//
//	signer, err := allora.NewRemoteSigner(ctx, allora.RemoteSignerConfig{
//	    BackendURL: "https://forge.allora.network",
//	    APIKey:     apiKey,
//	    WalletID:   walletID,
//	})
//	if err != nil {
//	    return err
//	}
//	signedTx, err := allora.SignTransactionWith(ctx, unsignedTx, signer, params)
func SignTransactionWith(
	ctx context.Context,
	unsignedTx []byte,
	signer Signer,
	params *TxParams,
) ([]byte, error) {
	// A nil context would panic in http.NewRequestWithContext when the signer is a
	// RemoteSigner. Default it to context.Background() (as SignTransaction already does) so a
	// nil ctx degrades to "no deadline/cancellation" instead of crashing during signing.
	if ctx == nil {
		ctx = context.Background()
	}
	if params == nil {
		return nil, fmt.Errorf("tx params are required")
	}
	// Only validate the fields actually read at sign time (ChainID). Gas/fee/granter/memo/timeout
	// are already encoded in unsignedTx, so requiring them again here would break the two-phase
	// flow documented above (build now, sign later with minimal params).
	if err := params.validateForSigning(); err != nil {
		return nil, fmt.Errorf("invalid transaction parameters: %w", err)
	}
	if len(unsignedTx) == 0 {
		return nil, fmt.Errorf("unsigned transaction is empty")
	}
	if isNilSigner(signer) {
		return nil, fmt.Errorf("signer is required")
	}

	builder := newTxBuilder()
	return builder.signTx(ctx, unsignedTx, signer, params)
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
