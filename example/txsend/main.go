// Command txsend demonstrates how to construct, sign, and broadcast a
// transaction against the Allora Network using the SDK's high-level
// `client.Tx().SendTx(...)` entry point, with a fee-granter (master/subsidy
// wallet) subsidising the gas.
//
// The example does not dial a node unless the ALLORA_RUN_EXAMPLE environment
// variable is set, so `go build ./...`, `go vet ./...`, and the default
// `go test ./...` stay clean. To actually broadcast a tx, configure the
// environment and re-run:
//
//	export ALLORA_RUN_EXAMPLE=1
//	export ALLORA_CHAIN_ID=allora-testnet-1          # or allora-mainnet-1
//	export ALLORA_GRPC=grpcs://allora-grpc.testnet.allora.network:443
//	# (optional) export ALLORA_REST=https://allora-api.testnet.allora.network
//	export ALLORA_MNEMONIC="word1 word2 ... word24"  # the signer key
//	# (optional) export FORGE_MASTER_GRANTER_ADDRESS=allo1... # master wallet paying gas
//	# (optional) export FORGE_API_KEY=...                       # to use a RemoteSigner (Privy)
//	# (optional) export FORGE_SIGNING_WALLET_ID=...            # to use a RemoteSigner (Privy)
//	export ALLORA_TO=allo1...                                 # recipient
//	export ALLORA_AMOUNT_UALLO=1000                           # send amount
//	export ALLORA_TOPIC_ID=42                                 # topic to register for
//	go run ./example/txsend
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"

	allora "github.com/allora-network/allora-sdk-go"
	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/txmsg"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func main() {
	// --- guard rail: refuse to dial by default -----------------------------
	// `go build` / `go vet` / `go test` never set ALLORA_RUN_EXAMPLE, so this
	// example only ever fires a tx when the operator explicitly opts in. This
	// keeps the example safe to leave in the repo.
	if os.Getenv("ALLORA_RUN_EXAMPLE") == "" {
		printUsage()
		return
	}

	// --- config / logger ---------------------------------------------------
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		With().Timestamp().Logger()

	chainID := requireEnv("ALLORA_CHAIN_ID")
	grpcURL := requireEnv("ALLORA_GRPC")
	toAddr := requireEnv("ALLORA_TO")
	amountStr := requireEnv("ALLORA_AMOUNT_UALLO")

	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		fatal(logger, "invalid ALLORA_AMOUNT_UALLO", err)
	}

	cfg := &config.ClientConfig{
		Endpoints: []config.EndpointConfig{
			{URL: grpcURL, Protocol: config.ProtocolGRPC},
		},
		RequestTimeout:    30 * time.Second,
		ConnectionTimeout: 10 * time.Second,
	}
	// Allow a co-deployed REST endpoint for the same client (purely additive
	// — the Tx sender uses gRPC; the REST endpoint is exposed for the
	// operator's monitoring convenience).
	if rest := os.Getenv("ALLORA_REST"); rest != "" {
		cfg.Endpoints = append(cfg.Endpoints, config.EndpointConfig{
			URL: rest, Protocol: config.ProtocolREST,
		})
	}

	client, err := allora.NewClient(cfg, logger)
	if err != nil {
		fatal(logger, "failed to create allora client", err)
	}
	defer client.Close()

	// --- signer: pick local (default) or RemoteSigner (Privy) --------------
	// The Tx() entry point takes an allora.Signer; both *Wallet.PrivKey and
	// *RemoteSigner satisfy it. Switch on FORGE_API_KEY to choose the path
	// without changing the rest of the example.
	signer, signerLabel, err := buildSigner(logger)
	if err != nil {
		fatal(logger, "building signer", err)
	}
	logger.Info().Str("signer", signerLabel).Msg("signer ready")

	// --- fee granter (optional) -------------------------------------------
	// Resolve the master/subsidy wallet address. Discovery order matches the
	// sibling SDKs:
	//   1. FORGE_MASTER_GRANTER_ADDRESS env var (operator override)
	//   2. (RemoteSigner-only) signer.ResolveFeeGranter() — backend-advertised
	//   3. nil — the signing wallet pays its own fees
	granter, err := resolveGranter(signer)
	if err != nil {
		fatal(logger, "resolving fee granter", err)
	}
	if granter != nil {
		logger.Info().Str("fee_granter", granter.String()).Msg("fee-granter subsidy will pay gas")
	} else {
		logger.Info().Msg("no fee granter configured; signer will pay its own gas")
	}

	// --- build the message set --------------------------------------------
	// The signer is the From address for the bank send (we resolve it from
	// the signer's pubkey so this example works for both local + remote).
	fromAddr := sdk.AccAddress(signer.PubKey().Address())
	sendMsg, err := txmsg.NewSend(
		fromAddr.String(), toAddr,
		sdk.NewCoins(sdk.NewInt64Coin("uallo", amount)),
	)
	if err != nil {
		fatal(logger, "building NewSend msg", err)
	}

	// Register a worker on a topic (ALLORA_TOPIC_ID required, default 1 if
	// unset). Demonstrates an emissions-module message in the same tx set.
	topicID := uint64(1)
	if v := os.Getenv("ALLORA_TOPIC_ID"); v != "" {
		n, perr := strconv.ParseUint(v, 10, 64)
		if perr != nil || n == 0 {
			fatal(logger, "invalid ALLORA_TOPIC_ID", perr)
		}
		topicID = n
	}
	regMsg, err := txmsg.NewRegisterWorker(txmsg.RegisterWorkerParams{
		Sender:  fromAddr.String(),
		TopicId: topicID,
		Owner:   fromAddr.String(),
	})
	if err != nil {
		fatal(logger, "building NewRegisterWorker msg", err)
	}

	msgs := []sdk.Msg{sendMsg, regMsg}

	// --- broadcast via the high-level SendTx entrypoint --------------------
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	result, err := client.Tx().SendTx(ctx, signer, msgs, allora.SendOptions{
		ChainID:    chainID,
		FeeGranter: granter,
		Memo:       "allora-sdk-go example tx",
		// GasAdjustment, GasLimit, GasPrice, BroadcastMode, SkipWait, MaxRetries
		// left at their zero values; SendTx applies the documented defaults
		// (simulate × 1.5, GasPrice 0.01, BroadcastModeSync, MaxRetries 2).
	})
	if err != nil {
		fatal(logger, "SendTx failed", err)
	}

	// --- print the result --------------------------------------------------
	fmt.Printf("✅ tx broadcast succeeded\n")
	fmt.Printf("   tx_hash: %s\n", result.TxHash)
	fmt.Printf("   height:  %d\n", result.Height)
	fmt.Printf("   code:    %d (%s)\n", result.Code, result.Codespace)
	fmt.Printf("   gas:     wanted=%d used=%d\n", result.GasWanted, result.GasUsed)
}

// buildSigner constructs the Signer the Tx() entry point will use. Returns a
// human-readable label too so the example's log line is informative.
//
// Two paths, side by side. Flip the conditional on FORGE_API_KEY (and
// supply FORGE_SIGNING_WALLET_ID) to take the RemoteSigner (Privy) path.
func buildSigner(logger zerolog.Logger) (allora.Signer, string, error) {
	// --- (a) RemoteSigner (Privy) — also shown commented below ------------
	// This is the path to use in production: the worker's private key lives
	// in Privy/Forge and never enters the worker process. Uncomment the block
	// and supply FORGE_API_KEY + FORGE_SIGNING_WALLET_ID to take it.
	//
	//     ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//     defer cancel()
	//     rs, err := allora.NewRemoteSigner(ctx, allora.RemoteSignerConfig{
	//         BackendURL: "https://forge.allora.network",
	//         APIKey:     os.Getenv("FORGE_API_KEY"),
	//         WalletID:   os.Getenv("FORGE_SIGNING_WALLET_ID"),
	//     })
	//     if err != nil {
	//         return nil, "", fmt.Errorf("new remote signer: %w", err)
	//     }
	//     return rs, "RemoteSigner (Privy)", nil
	//
	// For the *topic-bound* variant (one worker = one topic), use
	// NewRemoteSignerForTopic instead — it idempotently provisions a wallet
	// for a given topic ID. See remote_signer.go for the full contract.

	if os.Getenv("FORGE_API_KEY") != "" && os.Getenv("FORGE_SIGNING_WALLET_ID") != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		rs, err := allora.NewRemoteSigner(ctx, allora.RemoteSignerConfig{
			BackendURL: "https://forge.allora.network",
			APIKey:     os.Getenv("FORGE_API_KEY"),
			WalletID:   os.Getenv("FORGE_SIGNING_WALLET_ID"),
		})
		if err != nil {
			return nil, "", fmt.Errorf("new remote signer: %w", err)
		}
		return rs, "RemoteSigner (Privy)", nil
	}

	// --- (b) local Wallet — the default path for this example -------------
	// NewWalletFromMnemonic derives the secp256k1 key from a BIP-39 mnemonic
	// at the Allora HD path (m/44'/118'/0'/0/0). The returned *Wallet exposes
	// .PrivKey — a *secp256k1.PrivKey — which implements allora.Signer
	// directly, so we pass it straight into SendTx.
	mnemonic := requireEnv("ALLORA_MNEMONIC")
	w, err := allora.NewWalletFromMnemonic(mnemonic, allora.DefaultHDPath)
	if err != nil {
		return nil, "", fmt.Errorf("new wallet from mnemonic: %w", err)
	}
	return w.PrivKey, "local Wallet (mnemonic)", nil
}

// resolveGranter picks the fee granter in the documented precedence: env
// override, then (RemoteSigner-only) backend-discovered master granter, then
// nil. The result drops straight into SendOptions.FeeGranter; nil means the
// signing wallet pays its own gas.
func resolveGranter(signer allora.Signer) (sdk.AccAddress, error) {
	// 1. Env override (FORGE_MASTER_GRANTER_ADDRESS or legacy FEE_GRANTER) —
	// shared across allora-sdk-go / allora-sdk-py / allora-sdk-ts.
	if envGranter, err := allora.FeeGranterFromEnv(); err != nil {
		return nil, fmt.Errorf("FORGE_MASTER_GRANTER_ADDRESS: %w", err)
	} else if envGranter != nil {
		return envGranter, nil
	}
	// 2. RemoteSigner discovery: ask the backend which master wallet (if any)
	// has a feegrant set up to this wallet. Local signers have no discovery
	// path; skip silently.
	if rs, ok := signer.(*allora.RemoteSigner); ok {
		return rs.ResolveFeeGranter()
	}
	// 3. No granter — signer pays its own fees.
	return nil, nil
}

// --- small stdlib helpers (kept local so the example stays single-file) ---

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		fmt.Fprintf(os.Stderr, "missing required env var %q\n", key)
		os.Exit(2)
	}
	return v
}

func fatal(logger zerolog.Logger, msg string, err error) {
	logger.Error().Err(err).Msg(msg)
	os.Exit(1)
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `txsend: Allora SDK transaction-sending example.

This example does not broadcast a transaction unless ALLORA_RUN_EXAMPLE=1 is
set. To actually send a tx, set the variables below and re-run.

Required:
  ALLORA_RUN_EXAMPLE=1
  ALLORA_CHAIN_ID=allora-testnet-1
  ALLORA_GRPC=grpcs://allora-grpc.testnet.allora.network:443
  ALLORA_MNEMONIC="<24-word bip39 mnemonic>"
  ALLORA_TO=allo1...
  ALLORA_AMOUNT_UALLO=1000

Optional:
  ALLORA_REST=https://allora-api.testnet.allora.network      (co-deployed REST endpoint)
  ALLORA_TOPIC_ID=42                                        (topic for the register-worker msg)
  FORGE_MASTER_GRANTER_ADDRESS=allo1...                     (master wallet paying gas)
  FORGE_API_KEY=...            FORGE_SIGNING_WALLET_ID=...  (take the RemoteSigner/Privy path)`)
}
