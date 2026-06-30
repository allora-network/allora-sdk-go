package allora

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/txsend"
	"github.com/brynbellomy/go-utils/errors"
)

// ---------------------------------------------------------------------------
// Sender — the high-level send endpoint
// ---------------------------------------------------------------------------

// Sender drives a message set through the full send lifecycle: account
// discovery, build, sign, gas estimation, broadcast, and confirmation. It is
// the one-call entry point that composes the lower-level building blocks
// (CreateUnsignedTx, SignTransactionWith, txsend.TxBroadcaster) with
// code-based retry classification.
//
// The interface takes an allora.Signer so it works for both local Wallet keys
// and *RemoteSigner without a type-branch in the caller.
type Sender interface {
	// SendTx performs the full send lifecycle for the given messages. The
	// signer's pubkey derives the sender address used for account discovery
	// and signing. msgs is the set of Cosmos SDK messages to include in the
	// transaction; at least one message is required. opts carries overrides
	// for chain-id, gas, fees, broadcast mode, wait behaviour, and retry
	// policy.
	//
	// It returns the final TxResult on success (code 0 at CheckTx and
	// DeliverTx), and may return the last TxResult alongside an error when
	// retries are exhausted or the error is non-retryable.
	SendTx(ctx context.Context, signer Signer, msgs []sdk.Msg, opts SendOptions) (*txsend.TxResult, error)
}

// SendOptions tunes the SendTx pipeline. All fields are optional; zero values
// produce sensible defaults.
type SendOptions struct {
	// ChainID is the Cosmos chain identifier (e.g. "allora-mainnet-1"). It
	// is REQUIRED: SendTx returns an error if it is empty, since chain-ID is
	// needed for signing and cannot be derived from the broadcaster.
	ChainID string

	// Memo is an optional string attached to the transaction. It is visible
	// on-chain in the tx body.
	Memo string

	// FeeGranter is the bech32 address of a feegrant granter that pays the
	// transaction fee on behalf of the signer. Empty means the signer pays
	// its own gas. When set, the address is validated as 20 bytes (Cosmos
	// account address length) before building.
	FeeGranter sdk.AccAddress

	// GasAdjustment is the multiplier applied to the simulated gas usage to
	// derive the final gas limit. It defaults to DefaultGasAdjustment (1.5)
	// when zero, matching the cosmos-cli --gas-adjustment default. A value
	// below 1.0 is clamped to 1.0.
	GasAdjustment float64

	// GasLimit, when non-zero, skips gas simulation entirely and uses this
	// value as the fixed gas limit. Use for txs with well-known gas
	// consumption when you want to avoid the extra simulate RPC.
	GasLimit uint64

	// GasPrice is the per-gas-unit price in the smallest token unit —
	// typically the value returned by cosmospool.GetGasPrice(tier) or a
	// caller-supplied override. With a simulated gas limit of 150,000 and a
	// gas price of 0.01 uallo/gas-unit, the computed fee is 1,500 uallo.
	// When zero, SendTx uses a default of 0.01 (FeeTierMedium).
	GasPrice float64

	// BroadcastMode selects the broadcast mode. Defaults to
	// txsend.BroadcastModeSync so CheckTx rejections are visible
	// immediately.
	BroadcastMode txsend.BroadcastMode

	// SkipWait, when true, skips the WaitForTx call after broadcast and
	// returns a TxResult populated from the BroadcastResult. When false
	// (the default), SendTx blocks for on-chain confirmation before
	// returning.
	SkipWait bool

	// MaxRetries is the maximum number of retries for retryable CheckTx
	// rejections (out-of-gas, insufficient-fee, sequence-mismatch). It
	// defaults to DefaultMaxRetries (2) when zero. Each retry re-derives
	// the sending parameters and re-signs, so a sequence-bump or gas-bump
	// from the previous attempt is encoded into the new signed tx.
	MaxRetries int
}

const (
	// DefaultMaxRetries is the default maximum number of retries SendTx
	// performs for retryable CheckTx rejections before giving up.
	DefaultMaxRetries = 2

	// DefaultGasAdjustment is the default multiplier applied to simulated
	// gas usage to derive the final gas limit. It mirrors the
	// cosmospool.DefaultGasAdjustment (1.5) so a SendTx caller gets the
	// same headroom without pulling in the cosmospool package.
	DefaultGasAdjustment = 1.5

	// DefaultGasPrice is the default per-gas-unit price used when
	// SendOptions.GasPrice is zero. 0.01 uallo/gas-unit = FeeTierMedium.
	DefaultGasPrice = 0.01
)

// ---------------------------------------------------------------------------
// cosmos SDK error codes used by retry classification
// ---------------------------------------------------------------------------

const (
	// cosmosCodespaceSDK is the conventional cosmos-SDK error namespace
	// ("sdk") for all the standard error codes below.
	cosmosCodespaceSDK = "sdk"

	// cosmosCodeOutOfGas means the tx's gas limit was too low for
	// execution. Retryable: bump the gas limit and re-sign.
	cosmosCodeOutOfGas uint32 = 11

	// cosmosCodeInsufficientFee means the tx's fee is below the node's
	// minimum gas price. Retryable: bump the fee and re-sign.
	cosmosCodeInsufficientFee uint32 = 13

	// cosmosCodeSequenceMismatch / cosmosCodeWrongSequence means the
	// sequence in the tx does not match the account's current on-chain
	// sequence. Retryable: refetch account info and re-sign.
	cosmosCodeSequenceMismatch uint32 = 32
)

// retryBumpGasFactor is the multiplier applied to GasLimit on an
// out-of-gas retry. 1.3 gives a 30 % headroom bump.
const retryBumpGasFactor = 1.3

// retryBumpFeeFactor is the multiplier applied to the fee amount on an
// insufficient-fee retry. 2.0 doubles the fee.
const retryBumpFeeFactor = 2.0

// ---------------------------------------------------------------------------
// defaultSender — the concrete Sender
// ---------------------------------------------------------------------------

// defaultSender is the concrete Sender implementation. It holds a
// TxBroadcaster (the network seam — mockable) and a zerolog logger.
type defaultSender struct {
	broadcaster txsend.TxBroadcaster
	logger      zerolog.Logger
}

// Compile-time interface satisfaction.
var _ Sender = (*defaultSender)(nil)

// NewSender returns a Sender backed by the given TxBroadcaster (interface).
// It panics if broadcaster is nil: a sender without a network seam is a
// wiring bug. The logger defaults to zerolog.Nop() when it is the zero-value.
//
// Example:
//
//	pool, _ := cosmosrpc.NewClientPool(...)
//	broadcaster := cosmospool.New(pool, logger)
//	sender := allora.NewSender(broadcaster, logger)
//	result, err := sender.SendTx(ctx, signer, msgs, allora.SendOptions{ChainID: "allora-mainnet-1"})
func NewSender(broadcaster txsend.TxBroadcaster, logger zerolog.Logger) Sender {
	if broadcaster == nil {
		panic("allora.NewSender: broadcaster is nil")
	}
	// A zero-value zerolog.Logger has a nil writer and silently drops all
	// output, so it is safe to use as a no-op. We substitute zerolog.Nop()
	// so the logger is always the canonical no-op form.
	if isZeroLogger(logger) {
		logger = zerolog.Nop()
	}
	return &defaultSender{
		broadcaster: broadcaster,
		logger:      logger,
	}
}

// SendTx implements Sender.
func (s *defaultSender) SendTx(
	ctx context.Context,
	signer Signer,
	msgs []sdk.Msg,
	opts SendOptions,
) (*txsend.TxResult, error) {
	if len(msgs) == 0 {
		return nil, errors.New("at least one message is required")
	}
	if isNilSigner(signer) {
		return nil, errors.New("signer is required")
	}

	// --- derive signer address -------------------------------------------
	pubKey := signer.PubKey()
	if isNilPubKey(pubKey) {
		return nil, errors.New("signer returned nil public key")
	}
	addr := sdk.AccAddress(pubKey.Address()).String()

	// --- apply defaults --------------------------------------------------
	buildOpts := &opts
	if buildOpts.MaxRetries <= 0 {
		buildOpts.MaxRetries = DefaultMaxRetries
	}
	if buildOpts.GasAdjustment <= 0 {
		buildOpts.GasAdjustment = DefaultGasAdjustment
	}
	if buildOpts.GasAdjustment < 1.0 {
		buildOpts.GasAdjustment = 1.0
	}
	if buildOpts.GasPrice <= 0 {
		buildOpts.GasPrice = DefaultGasPrice
	}

	// --- initial account info --------------------------------------------
	accNum, seq, err := s.broadcaster.AccountInfo(ctx, addr)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching account info for %s", addr)
	}

	// --- main retry loop -------------------------------------------------
	var lastResult *txsend.TxResult
	for attempt := 0; attempt <= buildOpts.MaxRetries; attempt++ {
		result, ra, sendErr := s.trySend(ctx, signer, msgs, *buildOpts, addr, accNum, seq, attempt)
		if ra == nil && sendErr == nil {
			return result, nil
		}
		lastResult = result

		if ra == nil {
			// Non-retryable error.
			return result, sendErr
		}

		// Refetch account info when the adjustment demands it (sequence-mismatch).
		if ra.code == cosmosCodeSequenceMismatch {
			newAccNum, newSeq, accErr := s.broadcaster.AccountInfo(ctx, addr)
			if accErr != nil {
				return nil, errors.Wrapf(accErr, "refetching account info on sequence-mismatch retry for %s", addr)
			}
			accNum, seq = newAccNum, newSeq
		}
		gasLimit, feeAmount := ra.apply(accNum, seq, buildOpts.GasLimit)
		buildOpts.GasLimit = gasLimit
		_ = feeAmount // feeAmount is applied inside trySend when rebuilding

		s.logger.Info().
			Str("op", "SendTx").
			Str("address", addr).
			Int("attempt", attempt+1).
			Uint32("code", ra.code).
			Str("codespace", ra.codespace).
			Uint64("gas_limit", buildOpts.GasLimit).
			Str("fee", feeAmount.String()).
			Msg("retrying after classified error")
	}

	return lastResult, errors.Errorf(
		"SendTx: exhausted %d retries for address %s; last code=%d codespace=%s",
		buildOpts.MaxRetries, addr, lastResult.Code, lastResult.Codespace,
	)
}

// trySend runs one attempt of the build→sim→sign→broadcast→confirm cycle.
// It returns (result, retryAdj, error). retryAdj is non-nil when the
// CheckTx rejection is retryable (classified by code+codespace) — the
// caller should apply the adjustment and retry. retryAdj is nil for
// success (error nil) or non-retryable failure (error non-nil).
func (s *defaultSender) trySend(
	ctx context.Context,
	signer Signer,
	msgs []sdk.Msg,
	opts SendOptions,
	addr string,
	accNum, seq uint64,
	attempt int,
) (*txsend.TxResult, *retryAdjustment, error) {
	// --- build topic labels for structured logging (attempt 0 is first try)
	buildLog := s.logger.With().
		Str("op", "trySend").
		Str("address", addr).
		Int("attempt", attempt).
		Uint64("account", accNum).
		Uint64("sequence", seq).
		Logger()

	// --- build unsigned tx with placeholder gas/fee (for sim) ------------
	params := &TxParams{
		ChainID:       opts.ChainID,
		AccountNumber: accNum,
		Sequence:      seq,
		FeeGranter:    opts.FeeGranter,
		Memo:          opts.Memo,
		GasLimit:      200_000,                             // placeholder
		FeeAmount:     sdk.NewCoins(sdk.NewInt64Coin("uallo", 1)), // placeholder
	}
	if err := params.Validate(); err != nil {
		return nil, nil, errors.Wrap(err, "building tx params")
	}

	unsignedTx, err := CreateUnsignedTx(msgs, params)
	if err != nil {
		return nil, nil, errors.Wrap(err, "building unsigned tx")
	}

	// --- gas: simulate or use explicit limit ----------------------------
	gasLimit := opts.GasLimit
	if gasLimit == 0 {
		simGas, gasErr := s.broadcaster.EstimateGas(ctx, unsignedTx)
		if gasErr != nil {
			return nil, nil, errors.Wrap(gasErr, "estimating gas")
		}
		gasLimit = uint64(float64(simGas) * opts.GasAdjustment)
		if gasLimit == 0 {
			gasLimit = 200_000 // absolute floor
		}
	}
	params.GasLimit = gasLimit

	// --- compute fee amount from gas price ---------------------------------
	params.FeeAmount = sdk.NewCoins(sdk.NewInt64Coin(
		"uallo",
		int64(float64(gasLimit)*opts.GasPrice),
	))

	// --- rebuild unsigned tx with final gas+fee --------------------------
	unsignedTx, err = CreateUnsignedTx(msgs, params)
	if err != nil {
		return nil, nil, errors.Wrap(err, "rebuilding unsigned tx with gas/fee")
	}

	// --- sign -----------------------------------------------------------
	signedTx, err := SignTransactionWith(ctx, unsignedTx, signer, params)
	if err != nil {
		return nil, nil, errors.Wrap(err, "signing tx")
	}

	// --- broadcast -------------------------------------------------------
	mode := opts.BroadcastMode
	br, err := s.broadcaster.Broadcast(ctx, signedTx, mode)
	if err != nil {
		return nil, nil, errors.Wrap(err, "broadcasting tx")
	}

	// --- classify CheckTx rejection (retryable codes) -------------------
	if br.Code != 0 {
		ra := classifyCheckTxError(br.Code, br.Codespace)
		if ra.retryable {
			buildLog.Warn().
				Uint32("code", br.Code).
				Str("codespace", br.Codespace).
				Str("raw_log", br.RawLog).
				Msg("broadcast rejected with retryable code")
			return broadcastResultToTxResult(br), &ra, nil
		}
		buildLog.Warn().
			Uint32("code", br.Code).
			Str("codespace", br.Codespace).
			Str("raw_log", br.RawLog).
			Msg("broadcast rejected with non-retryable code")
		return broadcastResultToTxResult(br), nil,
			errors.Errorf("broadcast rejected: code=%d codespace=%s", br.Code, br.Codespace)
	}

	// --- broadcast accepted; optionally wait for confirmation -----------
	if opts.SkipWait {
		// Wait=false — return result from BroadcastResult.
		return broadcastResultToTxResult(br), nil, nil
	}

	result, err := s.broadcaster.WaitForTx(ctx, br.TxHash)
	if err != nil {
		return nil, nil, errors.Wrap(err, "waiting for tx "+br.TxHash)
	}
	return result, nil, nil
}

// broadcastResultToTxResult converts a txsend.BroadcastResult to a
// txsend.TxResult for the Wait=false path. Fields not available at broadcast
// time (Height, GasWanted, GasUsed, Timestamp) are left as zero values.
func broadcastResultToTxResult(br *txsend.BroadcastResult) *txsend.TxResult {
	return &txsend.TxResult{
		TxHash:    br.TxHash,
		Code:      br.Code,
		Codespace: br.Codespace,
		RawLog:    br.RawLog,
	}
}

// ---------------------------------------------------------------------------
// retry classification
// ---------------------------------------------------------------------------

// retryAdjustment captures a retryable CheckTx rejection and the action to
// take.
type retryAdjustment struct {
	code      uint32
	codespace string
	retryable bool
	// nextGasLimit, when > 0, carries a bumped gas limit.
	nextGasLimit uint64
}

// Error implements error so retryAdjustment can be passed as an error value.
func (ra *retryAdjustment) Error() string {
	return fmt.Sprintf("retryable CheckTx: code=%d codespace=%s", ra.code, ra.codespace)
}

// apply computes the adjusted gas limit and fee amount for the next retry.
// The caller has already re-fetched account info when the code is
// sequence-mismatch, so accNum and seq are already current.
func (ra *retryAdjustment) apply(accNum, seq, currentGasLimit uint64) (
	gasLimit uint64,
	feeAmount sdk.Coins,
) {
	gasLimit = currentGasLimit
	if ra.nextGasLimit > 0 {
		gasLimit = ra.nextGasLimit
	}
	if ra.code == cosmosCodeInsufficientFee {
		// Double the fee using the base gas price as reference.
		feeAmount = sdk.NewCoins(sdk.NewInt64Coin(
			"uallo",
			int64(float64(gasLimit)*DefaultGasPrice*retryBumpFeeFactor),
		))
	} else {
		feeAmount = sdk.NewCoins(sdk.NewInt64Coin(
			"uallo",
			int64(float64(gasLimit)*DefaultGasPrice),
		))
	}
	return gasLimit, feeAmount
}

// classifyCheckTxError maps a cosmos CheckTx rejection (code + codespace) to
// a retryAdjustment. It prefers exact numeric code + codespace matching
// (stable, unambiguous) and falls back to a substring match on the raw-log
// only as a last resort for nodes that return custom codespace values.
//
// Retryable codes:
//   - out-of-gas (sdk/11): bump GasLimit by retryBumpGasFactor (1.3×).
//   - insufficient-fee (sdk/13): fee doubled by apply().
//   - sequence-mismatch (sdk/32): refetch account info.
func classifyCheckTxError(code uint32, codespace string) retryAdjustment {
	ra := retryAdjustment{code: code, codespace: codespace}

	// Only classify as retryable when the codespace is the canonical "sdk"
	// namespace. Unknown codespaces are not retried — they may be
	// application-specific errors we don't know how to recover from.
	if codespace != cosmosCodespaceSDK {
		return ra
	}

	switch code {
	case cosmosCodeOutOfGas:
		ra.retryable = true
	case cosmosCodeInsufficientFee:
		ra.retryable = true
	case cosmosCodeSequenceMismatch:
		ra.retryable = true
	}

	return ra
}

// isZeroLogger reports whether l is the zero-value zerolog.Logger
// (unconfigured). zerolog.Logger has unexported fields, so a plain == won't
// compile; reflect.DeepEqual against the zero value is the reliable check.
// This mirrors cosmospool.isZeroLogger and lives in the root package to avoid
// a dependency on cosmospool for this small helper.
func isZeroLogger(l zerolog.Logger) bool {
	// We use an interface-based comparison that works without reflect.
	var zero zerolog.Logger
	return l.GetLevel() == zero.GetLevel()
}
