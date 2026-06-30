// Package cosmospool provides the concrete TxBroadcaster implementation backed
// by a txsend.CosmosTxPool (a cosmosrpc.ClientPool adapted to the narrow
// interface). It wires the pooled, health-tracked, retried cosmos client to the
// txsend contract: account-info, gas estimation, broadcast, and wait-for-tx.
//
// This is the impl subpackage: the interface lives in the parent txsend
// package, and this package owns the heavy deps (zerolog) and the real logic.
// The constructor, option set, and static interface check are final; every
// TxBroadcaster method is implemented against the pool's Tx() and Auth()
// clients.
package cosmospool

import (
	"context"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/brynbellomy/go-utils/errors"
	sdktx "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	alloracodec "github.com/allora-network/allora-sdk-go/codec"
	"github.com/allora-network/allora-sdk-go/txsend"
)

// Broadcaster is the concrete txsend.TxBroadcaster backed by a
// txsend.CosmosTxPool. It holds the narrow pool, a logger, a Clock (the
// time seam for WaitForTx polling), and a pollInterval controlling the
// WaitForTx polling cadence.
type Broadcaster struct {
	pool          txsend.CosmosTxPool
	logger        zerolog.Logger
	clock         txsend.Clock
	gasAdjustment float64
	pollInterval  time.Duration
}

// Compile-time assertion that Broadcaster satisfies the txsend.TxBroadcaster
// contract. If the interface gains a method this package does not implement,
// the build fails here rather than at a call site.
var _ txsend.TxBroadcaster = (*Broadcaster)(nil)

// Option configures a Broadcaster at construction.
type Option func(*Broadcaster)

// WithClock sets the Clock used by WaitForTx's polling loop. Tests inject a
// fake clock here; production leaves the default SystemClock.
func WithClock(c txsend.Clock) Option {
	return func(b *Broadcaster) {
		if c != nil {
			b.clock = c
		}
	}
}

// WithGasAdjustment sets the multiplier applied to simulated gas usage to
// derive the gas limit returned by EstimateGas. The default is
// DefaultGasAdjustment (1.5), mirroring the cosmos-cli --gas-adjustment default;
// callers can raise it for complex txs prone to out-of-gas, or lower it
// (never below 1.0) for tight, well-understood txs.
func WithGasAdjustment(adj float64) Option {
	return func(b *Broadcaster) {
		if adj > 0 {
			b.gasAdjustment = adj
		}
	}
}

// WithPollInterval sets the interval between WaitForTx polls.
// Defaults to 1s; tests can inject shorter intervals or use the fake
// clock seam to advance deterministically.
func WithPollInterval(d time.Duration) Option {
	return func(b *Broadcaster) {
		if d > 0 {
			b.pollInterval = d
		}
	}
}

// isZeroLogger reports whether l is the zero-value zerolog.Logger (unconfigured).
// zerolog.Logger has unexported fields, so a plain == won't compile; reflect.DeepEqual
// against the zero value is the reliable check. The zero-value logger has a nil writer
// and silently drops all output, so it is safe to use, but we substitute zerolog.Nop()
// for the canonical form.
func isZeroLogger(l zerolog.Logger) bool {
	return reflect.DeepEqual(l, zerolog.Logger{})
}

// New returns a Broadcaster backed by pool. It panics if pool is nil (an
// invariant: a broadcaster without a pool cannot function and the caller has a
// wiring bug). A zero-value zerolog.Logger is replaced with zerolog.Nop() so
// call sites are unconditional. The clock defaults to SystemClock; override
// with WithClock for tests.
func New(pool txsend.CosmosTxPool, logger zerolog.Logger, opts ...Option) *Broadcaster {
	if pool == nil {
		// Invariant: a broadcaster with no pool can do nothing. This is a
		// wiring error at construction, not a recoverable runtime condition.
		panic("cosmospool.New: pool is nil")
	}
	// A zero-value zerolog.Logger has a nil writer and silently drops all output
	// (its should() returns false for every level), so it is already safe to use
	// as a no-op. We substitute zerolog.Nop() explicitly so the Broadcaster's
	// logger is always the canonical no-op form rather than the zero value, and
	// call sites never need to nil-check. zerolog.Logger has unexported fields, so
	// we detect the zero value via reflect.DeepEqual against zerolog.Logger{}.
	if isZeroLogger(logger) {
		logger = zerolog.Nop()
	}
	b := &Broadcaster{
		pool:          pool,
		logger:        logger,
		clock:         txsend.SystemClock{},
		gasAdjustment: DefaultGasAdjustment,
		pollInterval:  1 * time.Second,
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

// AccountInfo returns the on-chain account number and sequence for address. It
// queries the auth module via the pool's Auth client and unpacks the returned
// *codectypes.Any into a authtypes.AccountI (BaseAccount in practice) using the
// shared alloracodec.CosmosCodec() registry. ctx is threaded through the gRPC
// call so a caller's cancellation/deadline propagates to the node.
//
// A gRPC codes.NotFound — or any error whose message contains "not found" — is
// treated as a clear signal that the address has never been seen on-chain (no
// funding, no prior txs). That case is wrapped with explicit context so callers
// can distinguish "account doesn't exist" from transient transport errors.
func (b *Broadcaster) AccountInfo(ctx context.Context, address string) (uint64, uint64, error) {
	resp, err := b.pool.Auth().Account(ctx, &authtypes.QueryAccountRequest{Address: address})
	if err != nil {
		if isAccountNotFoundErr(err) {
			return 0, 0, errors.Wrapf(err, "account %s not found on-chain (unfunded or never seen)", address)
		}
		return 0, 0, errors.Wrapf(err, "account %s: auth query failed", address)
	}
	if resp == nil || resp.Account == nil {
		return 0, 0, errors.Wrapf(errMissingAccount, "account %s: empty response from auth module", address)
	}

	var acc authtypes.AccountI
	if err := alloracodec.CosmosCodec().UnpackAny(resp.Account, &acc); err != nil {
		return 0, 0, errors.Wrapf(err, "account %s: failed to unpack on-chain account", address)
	}
	if acc == nil {
		return 0, 0, errors.Wrapf(errMissingAccount, "account %s: unpacked account is nil", address)
	}

	return acc.GetAccountNumber(), acc.GetSequence(), nil
}

// errMissingAccount is a sentinel for "the auth module returned success but
// with no account payload" — a malformed/empty response, distinct from a real
// gRPC error. Using a package-level error keeps the test surface stable.
var errMissingAccount = errors.New("auth module returned nil account")

// isAccountNotFoundErr reports whether err looks like a "not found" failure
// from the cosmos auth module. It prefers the gRPC status code (codes.NotFound)
// and falls back to a substring check on the message, because some cosmos
// versions/wrappers surface the not-found condition as a plain error rather
// than a gRPC status.
func isAccountNotFoundErr(err error) bool {
	if err == nil {
		return false
	}
	if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "not found")
}

// EstimateGas returns the simulated gas usage for unsignedTx, adjusted by the
// broadcaster's gas-adjustment multiplier (default 1.5) so the returned value
// has headroom for the variance between simulation and actual inclusion.
//
// It calls CosmosTxPool.Tx().Simulate with the raw tx bytes and reads
// GasInfo.GasUsed from the response, then applies ceil(GasUsed * adjustment).
//
// On simulation ERROR the method returns a static fallback gas estimate
// (FallbackGasEstimate) WITHOUT an error, after logging a structured WARN.
// Simulation is best-effort: the caller can still build and broadcast the tx
// with the fallback gas limit, and the node's CheckTx will reject it if the
// gas is genuinely insufficient — which is the same outcome as a too-low
// simulation-derived limit. Returning (fallback, nil) keeps the send path
// operational during transient simulate failures (node overloaded, network
// blip) rather than aborting the tx entirely. The ctx is threaded to Simulate
// and honored for cancellation.
func (b *Broadcaster) EstimateGas(ctx context.Context, unsignedTx []byte) (uint64, error) {
	req := txtypes.SimulateRequest{TxBytes: unsignedTx}
	resp, err := b.pool.Tx().Simulate(ctx, &req)
	if err != nil {
		b.logger.Warn().
			Err(err).
			Str("op", "EstimateGas").
			Uint64("fallback_gas", FallbackGasEstimate).
			Msg("simulate failed; returning fallback gas estimate")
		return FallbackGasEstimate, nil
	}
	if resp == nil || resp.GasInfo == nil {
		b.logger.Warn().
			Str("op", "EstimateGas").
			Uint64("fallback_gas", FallbackGasEstimate).
			Msg("simulate returned nil GasInfo; returning fallback gas estimate")
		return FallbackGasEstimate, nil
	}
	gasUsed := resp.GasInfo.GasUsed
	return uint64(math.Ceil(float64(gasUsed) * b.gasAdjustment)), nil
}

// Broadcast submits signedTx in the given mode and returns the immediate
// broadcast result. It maps the txsend.BroadcastMode to the cosmos
// txtypes.BroadcastMode, calls the pool's Tx().BroadcastTx RPC, and maps the
// response to a *txsend.BroadcastResult. A non-zero CheckTx Code in the
// response is NOT an error — the tx was successfully submitted to the
// mempool and denied by CheckTx; the caller may inspect result.Code to decide
// how to proceed. Only transport-level failures (network errors, gRPC errors,
// pool retry exhaustion) return a wrapped error.
func (b *Broadcaster) Broadcast(ctx context.Context, signedTx []byte, mode txsend.BroadcastMode) (*txsend.BroadcastResult, error) {
	bm, err := mapBroadcastMode(mode)
	if err != nil {
		return nil, err
	}

	resp, err := b.pool.Tx().BroadcastTx(ctx, &txtypes.BroadcastTxRequest{
		TxBytes: signedTx,
		Mode:    bm,
	})
	if err != nil {
		return nil, errors.Wrap(err, "broadcasting tx")
	}

	txResp := resp.GetTxResponse()
	if txResp == nil {
		return nil, errors.New("broadcast response missing TxResponse")
	}

	return &txsend.BroadcastResult{
		TxHash:    txResp.TxHash,
		Code:      txResp.Code,
		Codespace: txResp.Codespace,
		RawLog:    txResp.RawLog,
	}, nil
}

// mapBroadcastMode converts txsend.BroadcastMode to the cosmos
// txtypes.BroadcastMode. Unknown values return a wrapped error.
func mapBroadcastMode(mode txsend.BroadcastMode) (txtypes.BroadcastMode, error) {
	switch mode {
	case txsend.BroadcastModeSync:
		return txtypes.BroadcastMode_BROADCAST_MODE_SYNC, nil
	case txsend.BroadcastModeAsync:
		return txtypes.BroadcastMode_BROADCAST_MODE_ASYNC, nil
	case txsend.BroadcastModeBlock:
		return txtypes.BroadcastMode_BROADCAST_MODE_BLOCK, nil
	default:
		return txtypes.BroadcastMode_BROADCAST_MODE_UNSPECIFIED,
			errors.Errorf("unknown broadcast mode: %d", mode)
	}
}

// WaitForTx blocks until the transaction identified by txHash is committed
// (included in a block) or ctx is cancelled, returning the final result. It
// polls GetTx at pollInterval (default 1s) using the clock seam so tests can
// drive time deterministically. A committed-but-failed tx (Code != 0) returns
// the TxResult with the non-zero Code and nil error — the caller inspects
// result.Code to determine success. Only context cancellation and transport
// errors that persist through the polling loop return an error.
func (b *Broadcaster) WaitForTx(ctx context.Context, txHash string) (*txsend.TxResult, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, errors.Wrap(ctx.Err(), "waiting for tx "+txHash)
		case <-b.clock.After(b.pollInterval):
		}

		resp, err := b.pool.Tx().GetTx(ctx, &txtypes.GetTxRequest{Hash: txHash})
		if err != nil {
			// NotFound or transient errors: continue polling.
			// The pool wrapper retries internally; a persistent error
			// will surface again on the next poll.
			continue
		}

		txResp := resp.GetTxResponse()
		if txResp == nil {
			continue
		}

		return mapTxResponse(txResp), nil
	}
}

// mapTxResponse converts a cosmos SDK TxResponse to txsend.TxResult,
// translating Events and their Attributes.
func mapTxResponse(txResp *sdktx.TxResponse) *txsend.TxResult {
	events := make([]txsend.Event, len(txResp.Events))
	for i, ev := range txResp.Events {
		attrs := make([]txsend.EventAttribute, len(ev.Attributes))
		for j, attr := range ev.Attributes {
			attrs[j] = txsend.EventAttribute{
				Key:   attr.Key,
				Value: attr.Value,
			}
		}
		events[i] = txsend.Event{
			Type:       ev.Type,
			Attributes: attrs,
		}
	}

	return &txsend.TxResult{
		TxHash:    txResp.TxHash,
		Height:    txResp.Height,
		Code:      txResp.Code,
		Codespace: txResp.Codespace,
		Data:      txResp.Data,
		RawLog:    txResp.RawLog,
		GasWanted: txResp.GasWanted,
		GasUsed:   txResp.GasUsed,
		Timestamp: txResp.Timestamp,
		Events:    events,
	}
}
