//go:build integration

// This is the end-to-end integration test for the SDK's transaction-sending
// pipeline. It exercises the full SendTx lifecycle against a live Allora
// node (account discovery → build → gas-simulate → sign → broadcast → wait
// for confirmation) and is gated behind the `integration` build tag so the
// default `go build ./...` and `go test -short ./...` never compile or run
// it — i.e. it is safe to leave in the repo and CI is not broken by the
// absent env vars below.
//
// Run it locally with:
//
//	mise exec -- go test -tags integration -run TestIntegrationSendTx -count=1 -v ./example/txsend
//
// Required env:
//
//	ALLORA_TEST_CHAIN_ID   e.g. allora-testnet-1
//	ALLORA_TEST_GRPC       e.g. allora-grpc.testnet.allora.network:443
//	ALLORA_TEST_MNEMONIC   24-word BIP-39 mnemonic for the signing key
//	ALLORA_TEST_TO         recipient bech32 address (allo1...)
//	ALLORA_TEST_AMOUNT_UALLO   amount to send, in uallo (positive int)
//
// Optional env:
//
//	ALLORA_TEST_TOPIC_ID   topic id to register a worker against; omit to skip
//	                          the NewRegisterWorker message
//	ALLORA_TEST_FEE_GRANTER  allo1... address of a master/subsidy wallet
//	                          with an active feegrant to the signer; omit to
//	                          have the signer pay its own gas
//	ALLORA_TEST_MEMO       memo to attach to the tx
//
// The test is `t.Skip(...)`'d when ANY required env var is missing or
// unparseable, so the file is safe to compile under the integration tag in
// any environment.
//
// On success the test prints the broadcast tx hash, block height, and final
// DeliverTx code; on failure it surfaces the SendTx error verbatim so a
// triage can read the chain's CheckTx/DeliverTx reason directly.
package main

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/rs/zerolog"

	allora "github.com/allora-network/allora-sdk-go"
	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/txmsg"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// TestIntegrationSendTx is the end-to-end SendTx smoke test. It is the
// integration counterpart to the unit tests in send_test.go, which exercise
// the retry/classification logic against a mock TxBroadcaster; this one
// drives a real node through the full pipeline.
func TestIntegrationSendTx(t *testing.T) {
	chainID, grpcURL, mnemonic, toAddr, amount, ok := integrationEnv(t)
	if !ok {
		return // t.Skip was already called inside integrationEnv
	}

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		With().Timestamp().Logger()

	cfg := &config.ClientConfig{
		Endpoints: []config.EndpointConfig{
			{URL: grpcURL, Protocol: config.ProtocolGRPC},
		},
		RequestTimeout:    30 * time.Second,
		ConnectionTimeout: 10 * time.Second,
	}
	client, err := allora.NewClient(cfg, logger)
	require.NoError(t, err, "allora.NewClient")
	defer client.Close()

	wallet, err := allora.NewWalletFromMnemonic(mnemonic, allora.DefaultHDPath)
	require.NoError(t, err, "NewWalletFromMnemonic")
	fromAddr := wallet.GetAddress()
	logger.Info().Str("from", fromAddr).Str("to", toAddr).Int64("amount", amount).Msg("integration tx params")

	// --- fee granter (optional) -------------------------------------------
	var granter sdk.AccAddress
	if fg := os.Getenv("ALLORA_TEST_FEE_GRANTER"); fg != "" {
		parsed, perr := sdk.AccAddressFromBech32(fg)
		require.NoError(t, perr, "ALLORA_TEST_FEE_GRANTER must be a bech32 allo1... address")
		granter = parsed
	}

	// --- message set -------------------------------------------------------
	msgs := []sdk.Msg{}

	sendMsg, err := txmsg.NewSend(
		fromAddr, toAddr,
		sdk.NewCoins(sdk.NewInt64Coin("uallo", amount)),
	)
	require.NoError(t, err, "txmsg.NewSend")
	msgs = append(msgs, sendMsg)

	if topicStr := os.Getenv("ALLORA_TEST_TOPIC_ID"); topicStr != "" {
		topicID, perr := strconv.ParseUint(topicStr, 10, 64)
		require.NoError(t, perr, "ALLORA_TEST_TOPIC_ID must be a positive uint64")
		require.Greater(t, topicID, uint64(0), "ALLORA_TEST_TOPIC_ID must be > 0")
		regMsg, rerr := txmsg.NewRegisterWorker(txmsg.RegisterWorkerParams{
			Sender:  fromAddr,
			TopicId: topicID,
			Owner:   fromAddr,
		})
		require.NoError(t, rerr, "txmsg.NewRegisterWorker")
		msgs = append(msgs, regMsg)
	}

	// --- SendTx end-to-end -------------------------------------------------
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	result, err := client.Tx().SendTx(ctx, wallet.PrivKey, msgs, allora.SendOptions{
		ChainID:    chainID,
		FeeGranter: granter,
		Memo:       os.Getenv("ALLORA_TEST_MEMO"),
	})
	require.NoError(t, err, "client.Tx().SendTx")
	require.NotNil(t, result, "SendTx returned a nil result on success")
	require.Equal(t, uint32(0), result.Code, "DeliverTx code must be 0 (success)")
	require.NotEmpty(t, result.TxHash, "tx hash must be populated")
	require.Greater(t, result.Height, int64(0), "tx must be committed in a non-zero-height block")

	t.Logf("✅ integration SendTx committed: hash=%s height=%d gas_wanted=%d gas_used=%d",
		result.TxHash, result.Height, result.GasWanted, result.GasUsed)
}

// integrationEnv validates the required env vars for the integration test.
// On missing/invalid input it calls t.Skip (so `go test -tags integration`
// against an unconfigured environment simply reports SKIP, not FAIL) and
// returns ok=false so the caller bails out cleanly.
func integrationEnv(t *testing.T) (chainID, grpcURL, mnemonic, toAddr string, amount int64, ok bool) {
	t.Helper()

	chainID = os.Getenv("ALLORA_TEST_CHAIN_ID")
	grpcURL = os.Getenv("ALLORA_TEST_GRPC")
	mnemonic = os.Getenv("ALLORA_TEST_MNEMONIC")
	toAddr = os.Getenv("ALLORA_TEST_TO")
	amountStr := os.Getenv("ALLORA_TEST_AMOUNT_UALLO")

	if chainID == "" || grpcURL == "" || mnemonic == "" || toAddr == "" || amountStr == "" {
		t.Skip("integration test skipped: set ALLORA_TEST_CHAIN_ID, ALLORA_TEST_GRPC, ALLORA_TEST_MNEMONIC, ALLORA_TEST_TO, ALLORA_TEST_AMOUNT_UALLO to enable")
		return "", "", "", "", 0, false
	}

	parsedAmount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil || parsedAmount <= 0 {
		t.Skipf("integration test skipped: ALLORA_TEST_AMOUNT_UALLO=%q must be a positive integer", amountStr)
		return "", "", "", "", 0, false
	}
	return chainID, grpcURL, mnemonic, toAddr, parsedAmount, true
}
