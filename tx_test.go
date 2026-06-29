package allora

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestCreateUnsignedSendTx(t *testing.T) {
	// Generate test wallet
	wallet, err := GenerateWallet()
	if err != nil {
		t.Fatalf("failed to generate wallet: %v", err)
	}

	// Create recipient address
	recipientWallet, err := GenerateWallet()
	if err != nil {
		t.Fatalf("failed to generate recipient wallet: %v", err)
	}

	// Create transaction parameters
	params := &TxParams{
		ChainID:       "allora-testnet-1",
		AccountNumber: 123,
		Sequence:      5,
		GasLimit:      200000,
		FeeAmount:     sdk.NewCoins(sdk.NewInt64Coin("uallo", 5000)),
		Memo:          "test transaction",
	}

	// Create amount to send
	amount := sdk.NewCoins(sdk.NewInt64Coin("uallo", 1000000))

	// Create unsigned transaction
	unsignedTx, err := CreateUnsignedSendTx(wallet.Address, recipientWallet.Address, amount, params)
	if err != nil {
		t.Fatalf("failed to create unsigned transaction: %v", err)
	}

	if len(unsignedTx) == 0 {
		t.Error("unsigned transaction should not be empty")
	}

	t.Logf("Unsigned transaction created: %d bytes", len(unsignedTx))
}

func TestSignTransaction(t *testing.T) {
	// Generate test wallet
	wallet, err := GenerateWallet()
	if err != nil {
		t.Fatalf("failed to generate wallet: %v", err)
	}

	// Create recipient address
	recipientWallet, err := GenerateWallet()
	if err != nil {
		t.Fatalf("failed to generate recipient wallet: %v", err)
	}

	// Create transaction parameters
	params := &TxParams{
		ChainID:       "allora-testnet-1",
		AccountNumber: 123,
		Sequence:      5,
		GasLimit:      200000,
		FeeAmount:     sdk.NewCoins(sdk.NewInt64Coin("uallo", 5000)),
		Memo:          "test transaction",
	}

	// Create amount to send
	amount := sdk.NewCoins(sdk.NewInt64Coin("uallo", 1000000))

	// Create unsigned transaction
	unsignedTx, err := CreateUnsignedSendTx(wallet.Address, recipientWallet.Address, amount, params)
	if err != nil {
		t.Fatalf("failed to create unsigned transaction: %v", err)
	}

	// Sign the transaction
	signedTx, err := SignTransaction(unsignedTx, wallet, params)
	if err != nil {
		t.Fatalf("failed to sign transaction: %v", err)
	}

	if len(signedTx) == 0 {
		t.Error("signed transaction should not be empty")
	}

	// Signed tx should be larger than unsigned (contains signature)
	if len(signedTx) <= len(unsignedTx) {
		t.Error("signed transaction should be larger than unsigned transaction")
	}

	t.Logf("Unsigned: %d bytes, Signed: %d bytes", len(unsignedTx), len(signedTx))
}

func TestCreateSignedSendTx(t *testing.T) {
	// Generate test wallet
	wallet, err := GenerateWallet()
	if err != nil {
		t.Fatalf("failed to generate wallet: %v", err)
	}

	// Create recipient address
	recipientWallet, err := GenerateWallet()
	if err != nil {
		t.Fatalf("failed to generate recipient wallet: %v", err)
	}

	// Create transaction parameters
	params := &TxParams{
		ChainID:       "allora-testnet-1",
		AccountNumber: 123,
		Sequence:      5,
		GasLimit:      200000,
		FeeAmount:     sdk.NewCoins(sdk.NewInt64Coin("uallo", 5000)),
		Memo:          "test transaction",
	}

	// Create amount to send
	amount := sdk.NewCoins(sdk.NewInt64Coin("uallo", 1000000))

	// Create and sign in one step
	signedTx, err := CreateSignedSendTx(wallet.Address, recipientWallet.Address, amount, wallet, params)
	if err != nil {
		t.Fatalf("failed to create signed transaction: %v", err)
	}

	if len(signedTx) == 0 {
		t.Error("signed transaction should not be empty")
	}

	t.Logf("Signed transaction created: %d bytes", len(signedTx))
}

func TestTxParamsValidation(t *testing.T) {
	tests := []struct {
		name        string
		params      *TxParams
		expectError bool
	}{
		{
			name: "valid params",
			params: &TxParams{
				ChainID:       "allora-testnet-1",
				AccountNumber: 123,
				Sequence:      5,
				GasLimit:      200000,
				FeeAmount:     sdk.NewCoins(sdk.NewInt64Coin("uallo", 5000)),
			},
			expectError: false,
		},
		{
			name: "missing chain ID",
			params: &TxParams{
				AccountNumber: 123,
				Sequence:      5,
				GasLimit:      200000,
				FeeAmount:     sdk.NewCoins(sdk.NewInt64Coin("uallo", 5000)),
			},
			expectError: true,
		},
		{
			name: "zero gas limit",
			params: &TxParams{
				ChainID:       "allora-testnet-1",
				AccountNumber: 123,
				Sequence:      5,
				GasLimit:      0,
				FeeAmount:     sdk.NewCoins(sdk.NewInt64Coin("uallo", 5000)),
			},
			expectError: true,
		},
		{
			name: "empty fee",
			params: &TxParams{
				ChainID:       "allora-testnet-1",
				AccountNumber: 123,
				Sequence:      5,
				GasLimit:      200000,
				FeeAmount:     sdk.Coins{},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestDefaultTxParams(t *testing.T) {
	params := DefaultTxParams()

	if params == nil {
		t.Fatal("default params should not be nil")
	}

	if params.GasLimit == 0 {
		t.Error("default gas limit should not be zero")
	}

	if params.FeeAmount.Empty() {
		t.Error("default fee should not be empty")
	}

	t.Logf("Default gas limit: %d", params.GasLimit)
	t.Logf("Default fee: %s", params.FeeAmount)
}

func TestCreateUnsignedSendTxWithMemo(t *testing.T) {
	wallet, _ := GenerateWallet()
	recipientWallet, _ := GenerateWallet()

	params := &TxParams{
		ChainID:       "allora-testnet-1",
		AccountNumber: 123,
		Sequence:      5,
		GasLimit:      200000,
		FeeAmount:     sdk.NewCoins(sdk.NewInt64Coin("uallo", 5000)),
		Memo:          "deposit-id-12345",
	}

	amount := sdk.NewCoins(sdk.NewInt64Coin("uallo", 1000000))

	unsignedTx, err := CreateUnsignedSendTx(wallet.Address, recipientWallet.Address, amount, params)
	if err != nil {
		t.Fatalf("failed to create unsigned transaction: %v", err)
	}

	if len(unsignedTx) == 0 {
		t.Error("unsigned transaction should not be empty")
	}

	t.Logf("Transaction with memo created: %d bytes", len(unsignedTx))
}

func TestCreateUnsignedSendTxWithTimeoutHeight(t *testing.T) {
	wallet, _ := GenerateWallet()
	recipientWallet, _ := GenerateWallet()

	params := &TxParams{
		ChainID:       "allora-testnet-1",
		AccountNumber: 123,
		Sequence:      5,
		GasLimit:      200000,
		FeeAmount:     sdk.NewCoins(sdk.NewInt64Coin("uallo", 5000)),
		TimeoutHeight: 1000000,
	}

	amount := sdk.NewCoins(sdk.NewInt64Coin("uallo", 1000000))

	unsignedTx, err := CreateUnsignedSendTx(wallet.Address, recipientWallet.Address, amount, params)
	if err != nil {
		t.Fatalf("failed to create unsigned transaction: %v", err)
	}

	if len(unsignedTx) == 0 {
		t.Error("unsigned transaction should not be empty")
	}

	t.Logf("Transaction with timeout created: %d bytes", len(unsignedTx))
}

func TestInvalidInputs(t *testing.T) {
	wallet, _ := GenerateWallet()
	recipientWallet, _ := GenerateWallet()

	params := &TxParams{
		ChainID:       "allora-testnet-1",
		AccountNumber: 123,
		Sequence:      5,
		GasLimit:      200000,
		FeeAmount:     sdk.NewCoins(sdk.NewInt64Coin("uallo", 5000)),
	}

	amount := sdk.NewCoins(sdk.NewInt64Coin("uallo", 1000000))

	// Test with empty from address
	_, err := CreateUnsignedSendTx(sdk.AccAddress{}, recipientWallet.Address, amount, params)
	if err == nil {
		t.Error("expected error with empty from address")
	}

	// Test with empty to address
	_, err = CreateUnsignedSendTx(wallet.Address, sdk.AccAddress{}, amount, params)
	if err == nil {
		t.Error("expected error with empty to address")
	}

	// Test with empty amount
	_, err = CreateUnsignedSendTx(wallet.Address, recipientWallet.Address, sdk.Coins{}, params)
	if err == nil {
		t.Error("expected error with empty amount")
	}

	// Test signing with nil wallet
	unsignedTx, _ := CreateUnsignedSendTx(wallet.Address, recipientWallet.Address, amount, params)
	_, err = SignTransaction(unsignedTx, nil, params)
	if err == nil {
		t.Error("expected error with nil wallet")
	}

	// Test signing with empty tx bytes
	_, err = SignTransaction([]byte{}, wallet, params)
	if err == nil {
		t.Error("expected error with empty transaction bytes")
	}
}

func TestWalletAddressMismatch(t *testing.T) {
	wallet1, _ := GenerateWallet()
	wallet2, _ := GenerateWallet()
	recipientWallet, _ := GenerateWallet()

	params := &TxParams{
		ChainID:       "allora-testnet-1",
		AccountNumber: 123,
		Sequence:      5,
		GasLimit:      200000,
		FeeAmount:     sdk.NewCoins(sdk.NewInt64Coin("uallo", 5000)),
	}

	amount := sdk.NewCoins(sdk.NewInt64Coin("uallo", 1000000))

	// Try to sign with wallet2 when transaction is from wallet1
	_, err := CreateSignedSendTx(wallet1.Address, recipientWallet.Address, amount, wallet2, params)
	if err == nil {
		t.Error("expected error when wallet address doesn't match from address")
	}
}

func TestTxRoundTrip(t *testing.T) {
	// Create and sign a transaction
	wallet, _ := GenerateWallet()
	recipientWallet, _ := GenerateWallet()

	params := &TxParams{
		ChainID:       "allora-testnet-1",
		AccountNumber: 123,
		Sequence:      5,
		GasLimit:      200000,
		FeeAmount:     sdk.NewCoins(sdk.NewInt64Coin("uallo", 5000)),
		Memo:          "round-trip-test",
	}

	amount := sdk.NewCoins(sdk.NewInt64Coin("uallo", 1000000))

	// Create signed transaction
	signedTx, err := CreateSignedSendTx(wallet.Address, recipientWallet.Address, amount, wallet, params)
	if err != nil {
		t.Fatalf("failed to create signed transaction: %v", err)
	}

	// Parse it back
	parsedTx, err := ParseTxBytes(signedTx)
	if err != nil {
		t.Fatalf("failed to parse transaction: %v", err)
	}

	if parsedTx == nil {
		t.Fatal("parsed transaction should not be nil")
	}

	// Check that it has messages
	msgs := (*parsedTx).GetMsgs()
	if len(msgs) == 0 {
		t.Error("transaction should have at least one message")
	}

	t.Logf("Parsed transaction with %d message(s)", len(msgs))
}

// TestSignTransactionWith_RejectsSenderMismatch pins the signer-vs-sender guard in
// signTx: signing a MsgSend whose FromAddress is walletA with walletB's key must fail
// before producing a tx the chain would reject opaquely.
func TestSignTransactionWith_RejectsSenderMismatch(t *testing.T) {
	walletA, err := NewWalletFromMnemonic(testMnemonic, DefaultHDPath)
	require.NoError(t, err)
	walletB, err := GenerateWallet()
	require.NoError(t, err)

	amount := sdk.NewCoins(sdk.NewInt64Coin("uallo", 1000))
	params := &TxParams{
		ChainID:       "allora-testnet-1",
		AccountNumber: 7,
		Sequence:      3,
		GasLimit:      200000,
		FeeAmount:     sdk.NewCoins(sdk.NewInt64Coin("uallo", 5000)),
	}
	// Build the tx with walletA as sender, then try to sign it with walletB's key.
	unsigned, err := CreateUnsignedSendTx(walletA.Address, walletB.Address, amount, params)
	require.NoError(t, err)

	_, err = SignTransactionWith(context.Background(), unsigned, walletB.PrivKey, params)
	require.ErrorContains(t, err, "does not match transaction sender")
}

func TestTxParams_FeeGranterIsEncoded(t *testing.T) {
	wallet, err := NewWalletFromMnemonic(testMnemonic, DefaultHDPath)
	require.NoError(t, err)

	granter, err := sdk.AccAddressFromBech32(wallet.GetAddress())
	require.NoError(t, err)

	amount := sdk.NewCoins(sdk.NewInt64Coin("uallo", 1000))
	base := &TxParams{
		ChainID:   "allora-testnet-1",
		GasLimit:  200000,
		FeeAmount: sdk.NewCoins(sdk.NewInt64Coin("uallo", 5000)),
	}
	withGranter := *base
	withGranter.FeeGranter = granter

	unsignedNoGranter, err := CreateUnsignedSendTx(wallet.Address, wallet.Address, amount, base)
	require.NoError(t, err)
	unsignedWithGranter, err := CreateUnsignedSendTx(wallet.Address, wallet.Address, amount, &withGranter)
	require.NoError(t, err)

	// Setting the fee granter must change the encoded AuthInfo (the granter address is
	// written into the tx). An empty granter leaves the tx untouched.
	require.NotEqual(t, unsignedNoGranter, unsignedWithGranter)
	require.Greater(t, len(unsignedWithGranter), len(unsignedNoGranter))
}

func TestFeeGranterFromEnv(t *testing.T) {
	wallet, err := NewWalletFromMnemonic(testMnemonic, DefaultHDPath)
	require.NoError(t, err)
	valid := wallet.GetAddress()

	// Unset/blank → nil granter (the signing wallet pays its own fees).
	t.Setenv("FORGE_MASTER_GRANTER_ADDRESS", "")
	got, err := FeeGranterFromEnv()
	require.NoError(t, err)
	require.Nil(t, got)

	// A valid allo address is parsed into the matching AccAddress.
	t.Setenv("FORGE_MASTER_GRANTER_ADDRESS", valid)
	got, err = FeeGranterFromEnv()
	require.NoError(t, err)
	require.Equal(t, valid, got.String())

	// Whitespace is trimmed.
	t.Setenv("FORGE_MASTER_GRANTER_ADDRESS", "  "+valid+"  ")
	got, err = FeeGranterFromEnv()
	require.NoError(t, err)
	require.Equal(t, valid, got.String())

	// A malformed value is rejected at config time rather than reaching broadcast.
	t.Setenv("FORGE_MASTER_GRANTER_ADDRESS", "not-a-bech32-address")
	_, err = FeeGranterFromEnv()
	require.Error(t, err)
}
