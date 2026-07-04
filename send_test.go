package allora

import (
	"context"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/allora-network/allora-sdk-go/txmsg"
	"github.com/allora-network/allora-sdk-go/txsend"
	txsendmocks "github.com/allora-network/allora-sdk-go/txsend/mocks"
)

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

// testWallet returns a wallet derived from the shared test mnemonic so keys
// are deterministic across test runs.
func testWallet(t *testing.T) *Wallet {
	t.Helper()
	w, err := NewWalletFromMnemonic(testMnemonic, DefaultHDPath)
	require.NoError(t, err)
	return w
}

// testMsgs returns a minimal set of valid Cosmos messages (a bank send to self)
// for exercising SendTx without importing txmsg in every subtest.
func testMsgs(t *testing.T, wallet *Wallet) []sdk.Msg {
	t.Helper()
	msg, err := txmsg.NewSend(
		wallet.GetAddress(),
		wallet.GetAddress(),
		sdk.NewCoins(sdk.NewInt64Coin("uallo", 1000)),
	)
	require.NoError(t, err)
	return []sdk.Msg{msg}
}

// newMockBroadcaster returns a fresh MockTxBroadcaster for a test.
func newMockBroadcaster(t *testing.T) *txsendmocks.MockTxBroadcaster {
	return txsendmocks.NewMockTxBroadcaster(t)
}

// stubSender creates a defaultSender backed by a mock broadcaster and a no-op
// logger for tests. Panics on nil broadcaster.
func stubSender(t *testing.T, b txsend.TxBroadcaster) *defaultSender {
	t.Helper()
	return &defaultSender{
		broadcaster: b,
		logger:      zerolog.Nop(),
	}
}

// ---------------------------------------------------------------------------
// Happy path
// ---------------------------------------------------------------------------

func TestSendTx_HappyPath(t *testing.T) {
	w := testWallet(t)
	msgs := testMsgs(t, w)
	mockBC := newMockBroadcaster(t)

	// AccountInfo → 1 call.
	mockBC.EXPECT().
		AccountInfo(mock.Anything, w.GetAddress()).
		Return(uint64(7), uint64(3), error(nil)).
		Once()

	// EstimateGas → 1 call.
	mockBC.EXPECT().
		EstimateGas(mock.Anything, mock.Anything).
		Return(uint64(150_000), error(nil)).
		Once()

	// Broadcast → succeeds with code 0.
	mockBC.EXPECT().
		Broadcast(mock.Anything, mock.Anything, txsend.BroadcastModeSync).
		Return(&txsend.BroadcastResult{TxHash: "HASH1", Code: 0}, error(nil)).
		Once()

	// WaitForTx → committed with code 0.
	mockBC.EXPECT().
		WaitForTx(mock.Anything, "HASH1").
		Return(&txsend.TxResult{TxHash: "HASH1", Code: 0, Height: 100}, error(nil)).
		Once()

	s := stubSender(t, mockBC)
	result, err := s.SendTx(
		context.Background(),
		w.PrivKey,
		msgs,
		SendOptions{ChainID: "allora-testnet-1"},
	)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "HASH1", result.TxHash)
	require.Equal(t, uint32(0), result.Code)
	require.Equal(t, int64(100), result.Height)
}

// ---------------------------------------------------------------------------
// Explicit GasLimit — EstimateGas NOT called
// ---------------------------------------------------------------------------

func TestSendTx_ExplicitGasLimitSkipsSimulate(t *testing.T) {
	w := testWallet(t)
	msgs := testMsgs(t, w)
	mockBC := newMockBroadcaster(t)

	mockBC.EXPECT().
		AccountInfo(mock.Anything, w.GetAddress()).
		Return(uint64(7), uint64(3), error(nil)).
		Once()

	// EstimateGas is NOT expected — we set GasLimit explicitly.

	mockBC.EXPECT().
		Broadcast(mock.Anything, mock.Anything, txsend.BroadcastModeSync).
		Return(&txsend.BroadcastResult{TxHash: "HASH2", Code: 0}, error(nil)).
		Once()

	mockBC.EXPECT().
		WaitForTx(mock.Anything, "HASH2").
		Return(&txsend.TxResult{TxHash: "HASH2", Code: 0, Height: 200}, error(nil)).
		Once()

	s := stubSender(t, mockBC)
	result, err := s.SendTx(
		context.Background(),
		w.PrivKey,
		msgs,
		SendOptions{ChainID: "allora-testnet-1", GasLimit: 300_000},
	)
	require.NoError(t, err)
	require.Equal(t, "HASH2", result.TxHash)

	// Assert EstimateGas was never called.
	mockBC.AssertNotCalled(t, "EstimateGas", mock.Anything, mock.Anything)
}

// ---------------------------------------------------------------------------
// SkipWait=true — WaitForTx NOT called
// ---------------------------------------------------------------------------

func TestSendTx_SkipWait(t *testing.T) {
	w := testWallet(t)
	msgs := testMsgs(t, w)
	mockBC := newMockBroadcaster(t)

	mockBC.EXPECT().
		AccountInfo(mock.Anything, w.GetAddress()).
		Return(uint64(7), uint64(3), error(nil)).
		Once()

	mockBC.EXPECT().
		EstimateGas(mock.Anything, mock.Anything).
		Return(uint64(150_000), error(nil)).
		Once()

	mockBC.EXPECT().
		Broadcast(mock.Anything, mock.Anything, txsend.BroadcastModeSync).
		Return(&txsend.BroadcastResult{TxHash: "HASH3", Code: 0}, error(nil)).
		Once()

	s := stubSender(t, mockBC)
	result, err := s.SendTx(
		context.Background(),
		w.PrivKey,
		msgs,
		SendOptions{ChainID: "allora-testnet-1", SkipWait: true},
	)
	require.NoError(t, err)
	require.Equal(t, "HASH3", result.TxHash)
	require.Equal(t, uint32(0), result.Code)

	// WaitForTx must NOT be called.
	mockBC.AssertNotCalled(t, "WaitForTx", mock.Anything, mock.Anything)
}

// ---------------------------------------------------------------------------
// Retry: sequence-mismatch (code 32) once then success
// ---------------------------------------------------------------------------

func TestSendTx_RetrySequenceMismatch(t *testing.T) {
	w := testWallet(t)
	msgs := testMsgs(t, w)
	mockBC := newMockBroadcaster(t)

	// First AccountInfo → seq=3.
	mockBC.EXPECT().
		AccountInfo(mock.Anything, w.GetAddress()).
		Return(uint64(7), uint64(3), error(nil)).
		Once()

	// EstimateGas → OK. Only the initial attempt simulates; retries reuse the
	// established gas limit (threaded via retryAdjustment.nextGasLimit), so a
	// sequence-mismatch retry does not re-estimate.
	mockBC.EXPECT().
		EstimateGas(mock.Anything, mock.Anything).
		Return(uint64(150_000), error(nil)).
		Once()

	// First Broadcast → sequence mismatch (code 32, codespace "sdk").
	mockBC.EXPECT().
		Broadcast(mock.Anything, mock.Anything, txsend.BroadcastModeSync).
		Return(&txsend.BroadcastResult{
			TxHash:    "HASH_SEQ1",
			Code:      32,
			Codespace: "sdk",
			RawLog:    "account sequence mismatch",
		}, error(nil)).
		Once()

	// Refetch AccountInfo after sequence-mismatch.
	mockBC.EXPECT().
		AccountInfo(mock.Anything, w.GetAddress()).
		Return(uint64(7), uint64(4), error(nil)).
		Once()

	// Second Broadcast → success.
	mockBC.EXPECT().
		Broadcast(mock.Anything, mock.Anything, txsend.BroadcastModeSync).
		Return(&txsend.BroadcastResult{TxHash: "HASH_SEQ2", Code: 0}, error(nil)).
		Once()

	// WaitForTx → committed.
	mockBC.EXPECT().
		WaitForTx(mock.Anything, "HASH_SEQ2").
		Return(&txsend.TxResult{TxHash: "HASH_SEQ2", Code: 0, Height: 500}, error(nil)).
		Once()

	s := stubSender(t, mockBC)
	result, err := s.SendTx(
		context.Background(),
		w.PrivKey,
		msgs,
		SendOptions{ChainID: "allora-testnet-1"},
	)
	require.NoError(t, err)
	require.Equal(t, "HASH_SEQ2", result.TxHash)
	require.Equal(t, uint32(0), result.Code)

	// AccountInfo should have been called twice (initial + refetch).
	mockBC.AssertNumberOfCalls(t, "AccountInfo", 2)
}

// ---------------------------------------------------------------------------
// Retry: out-of-gas (code 11) once then success
// ---------------------------------------------------------------------------

func TestSendTx_RetryOutOfGas(t *testing.T) {
	w := testWallet(t)
	msgs := testMsgs(t, w)
	mockBC := newMockBroadcaster(t)

	mockBC.EXPECT().
		AccountInfo(mock.Anything, w.GetAddress()).
		Return(uint64(7), uint64(3), error(nil)).
		Once()

	mockBC.EXPECT().
		EstimateGas(mock.Anything, mock.Anything).
		Return(uint64(150_000), error(nil)).
		Once() // only the initial attempt simulates; the out-of-gas retry reuses the bumped gas limit (no re-estimation)

	// Capture the gas limit of each broadcast tx so we can assert the retry
	// actually raised it (an out-of-gas retry that re-sends the same gas would
	// fail identically — see asg-pvd.8 fix).
	var broadcastGas []uint64
	captureGas := func(_ context.Context, signedTx []byte, _ txsend.BroadcastMode) {
		decoded, derr := ParseTxBytes(signedTx)
		require.NoError(t, derr)
		feeTx, ok := (*decoded).(sdk.FeeTx)
		require.True(t, ok, "decoded tx must be a FeeTx")
		broadcastGas = append(broadcastGas, feeTx.GetGas())
	}

	// First Broadcast → out-of-gas.
	mockBC.EXPECT().
		Broadcast(mock.Anything, mock.Anything, txsend.BroadcastModeSync).
		Run(captureGas).
		Return(&txsend.BroadcastResult{
			TxHash:    "HASH_OOG1",
			Code:      11,
			Codespace: "sdk",
			RawLog:    "out of gas",
		}, error(nil)).
		Once()

	// Second Broadcast → success (with bumped gas).
	mockBC.EXPECT().
		Broadcast(mock.Anything, mock.Anything, txsend.BroadcastModeSync).
		Run(captureGas).
		Return(&txsend.BroadcastResult{TxHash: "HASH_OOG2", Code: 0}, error(nil)).
		Once()

	mockBC.EXPECT().
		WaitForTx(mock.Anything, "HASH_OOG2").
		Return(&txsend.TxResult{TxHash: "HASH_OOG2", Code: 0, Height: 600}, error(nil)).
		Once()

	s := stubSender(t, mockBC)
	result, err := s.SendTx(
		context.Background(),
		w.PrivKey,
		msgs,
		SendOptions{ChainID: "allora-testnet-1"},
	)
	require.NoError(t, err)
	require.Equal(t, "HASH_OOG2", result.TxHash)
	require.Equal(t, uint32(0), result.Code)

	// Broadcast should be called 2 times (first failed, retry succeeded).
	mockBC.AssertNumberOfCalls(t, "Broadcast", 2)

	// The retry must carry STRICTLY MORE gas than the first attempt — otherwise
	// an out-of-gas retry is a no-op that fails identically. Guards the
	// asg-pvd.8 fix wiring retryBumpGasFactor through apply().
	require.Len(t, broadcastGas, 2)
	require.Greater(t, broadcastGas[1], broadcastGas[0],
		"out-of-gas retry must bump the gas limit (got %d then %d)", broadcastGas[0], broadcastGas[1])
}

// ---------------------------------------------------------------------------
// Retry: insufficient fee (code 13) — the fee must STRICTLY INCREASE on retry
// ---------------------------------------------------------------------------

func TestSendTx_RetryInsufficientFee(t *testing.T) {
	w := testWallet(t)
	msgs := testMsgs(t, w)
	mockBC := newMockBroadcaster(t)

	mockBC.EXPECT().
		AccountInfo(mock.Anything, w.GetAddress()).
		Return(uint64(7), uint64(3), error(nil)).
		Once()

	mockBC.EXPECT().
		EstimateGas(mock.Anything, mock.Anything).
		Return(uint64(150_000), error(nil)).
		Once() // only the initial attempt simulates; retries reuse the established gas limit

	// Capture the fee of each broadcast tx so we can assert the retry actually
	// raised it. An insufficient-fee retry that re-sends the same fee fails
	// identically — this test guards the fix that propagates apply()'s bumped
	// fee into the retry (previously discarded with `_ = feeAmount`).
	var broadcastFee []sdk.Coins
	captureFee := func(_ context.Context, signedTx []byte, _ txsend.BroadcastMode) {
		decoded, derr := ParseTxBytes(signedTx)
		require.NoError(t, derr)
		feeTx, ok := (*decoded).(sdk.FeeTx)
		require.True(t, ok, "decoded tx must be a FeeTx")
		broadcastFee = append(broadcastFee, feeTx.GetFee())
	}

	// First Broadcast → insufficient fee (code 13).
	mockBC.EXPECT().
		Broadcast(mock.Anything, mock.Anything, txsend.BroadcastModeSync).
		Run(captureFee).
		Return(&txsend.BroadcastResult{
			TxHash:    "HASH_FEE1",
			Code:      13,
			Codespace: "sdk",
			RawLog:    "insufficient fee",
		}, error(nil)).
		Once()

	// Second Broadcast → success (with bumped fee).
	mockBC.EXPECT().
		Broadcast(mock.Anything, mock.Anything, txsend.BroadcastModeSync).
		Run(captureFee).
		Return(&txsend.BroadcastResult{TxHash: "HASH_FEE2", Code: 0}, error(nil)).
		Once()

	mockBC.EXPECT().
		WaitForTx(mock.Anything, "HASH_FEE2").
		Return(&txsend.TxResult{TxHash: "HASH_FEE2", Code: 0, Height: 700}, error(nil)).
		Once()

	s := stubSender(t, mockBC)
	result, err := s.SendTx(
		context.Background(),
		w.PrivKey,
		msgs,
		SendOptions{ChainID: "allora-testnet-1"},
	)
	require.NoError(t, err)
	require.Equal(t, "HASH_FEE2", result.TxHash)
	require.Equal(t, uint32(0), result.Code)

	// Broadcast should be called 2 times (first failed, retry succeeded).
	mockBC.AssertNumberOfCalls(t, "Broadcast", 2)

	// The retry must carry STRICTLY MORE fee than the first attempt — otherwise
	// an insufficient-fee retry is a no-op that fails identically.
	require.Len(t, broadcastFee, 2)
	firstAmt := broadcastFee[0].AmountOf("uallo")
	secondAmt := broadcastFee[1].AmountOf("uallo")
	require.True(t, secondAmt.GT(firstAmt),
		"insufficient-fee retry must bump the fee (got %s then %s)", firstAmt, secondAmt)
}

// ---------------------------------------------------------------------------
// Retry: insufficient fee (code 13) exhausted
// ---------------------------------------------------------------------------

func TestSendTx_RetryExhaustedInsufficientFee(t *testing.T) {
	w := testWallet(t)
	msgs := testMsgs(t, w)
	mockBC := newMockBroadcaster(t)

	mockBC.EXPECT().
		AccountInfo(mock.Anything, w.GetAddress()).
		Return(uint64(7), uint64(3), error(nil)).
		Once()

	mockBC.EXPECT().
		EstimateGas(mock.Anything, mock.Anything).
		Return(uint64(150_000), error(nil)).
		Once() // only the initial attempt simulates; retries reuse the established gas limit

	// Three attempts all return code 13.
	for i := 0; i < 3; i++ {
		mockBC.EXPECT().
			Broadcast(mock.Anything, mock.Anything, txsend.BroadcastModeSync).
			Return(&txsend.BroadcastResult{
				TxHash:    "HASH_FEE",
				Code:      13,
				Codespace: "sdk",
				RawLog:    "insufficient fee",
			}, error(nil)).
			Once()
	}

	s := stubSender(t, mockBC)
	result, err := s.SendTx(
		context.Background(),
		w.PrivKey,
		msgs,
		SendOptions{ChainID: "allora-testnet-1", MaxRetries: 2}, // 0 + 2 = 3 attempts
	)
	require.Error(t, err)
	require.NotNil(t, result)
	require.Contains(t, err.Error(), "exhausted")
	require.Equal(t, uint32(13), result.Code)

	// Broadcast should be called 3 times.
	mockBC.AssertNumberOfCalls(t, "Broadcast", 3)
}

// ---------------------------------------------------------------------------
// Account not found — no broadcast attempted
// ---------------------------------------------------------------------------

func TestSendTx_AccountNotFound(t *testing.T) {
	w := testWallet(t)
	msgs := testMsgs(t, w)
	mockBC := newMockBroadcaster(t)

	mockBC.EXPECT().
		AccountInfo(mock.Anything, w.GetAddress()).
		Return(uint64(0), uint64(0), fmt.Errorf("account %s not found on-chain", w.GetAddress())).
		Once()

	s := stubSender(t, mockBC)
	_, err := s.SendTx(
		context.Background(),
		w.PrivKey,
		msgs,
		SendOptions{ChainID: "allora-testnet-1"},
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "fetching account info")

	// Broadcast must NOT be called.
	mockBC.AssertNotCalled(t, "Broadcast", mock.Anything, mock.Anything, mock.Anything)
}

// ---------------------------------------------------------------------------
// Sender interface compile-time check
// ---------------------------------------------------------------------------

func TestNewSenderPanicsOnNilBroadcaster(t *testing.T) {
	require.Panics(t, func() {
		NewSender(nil, zerolog.Nop())
	})
}

// ---------------------------------------------------------------------------
// classifyCheckTxError unit tests
// ---------------------------------------------------------------------------

func TestClassifyCheckTxError(t *testing.T) {
	tests := []struct {
		name        string
		code        uint32
		codespace   string
		wantRetry   bool
	}{
		{"out-of-gas sdk", 11, "sdk", true},
		{"insufficient-fee sdk", 13, "sdk", true},
		{"sequence-mismatch sdk", 32, "sdk", true},
		{"unknown code sdk", 99, "sdk", false},
		{"out-of-gas wrong codespace", 11, "bank", false},
		{"success code", 0, "sdk", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ra := classifyCheckTxError(tt.code, tt.codespace)
			require.Equal(t, tt.wantRetry, ra.retryable)
		})
	}
}
