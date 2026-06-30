// Package cosmospool provides the concrete TxBroadcaster implementation backed
// by a txsend.CosmosTxPool (a cosmosrpc.ClientPool adapted to the narrow
// interface). It wires the pooled, health-tracked, retried cosmos client to the
// txsend contract: account-info, gas estimation, broadcast, and wait-for-tx.
//
// This is the impl subpackage: the interface lives in the parent txsend
// package, and this package owns the heavy deps (zerolog) and the real logic.
// The methods are currently stubs — each is implemented in a later bead
// (asg-pvd.3/.4/.5) — but the constructor, option set, and static interface
// check are real and final so the wiring shape is fixed.
package cosmospool

import (
	"context"
	"math"
	"reflect"

	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"

	"github.com/allora-network/allora-sdk-go/txsend"
)

// Broadcaster is the concrete txsend.TxBroadcaster backed by a
// txsend.CosmosTxPool. It holds the narrow pool, a logger, and a Clock (the
// time seam for WaitForTx polling). Methods are stubbed until their beads land.
type Broadcaster struct {
	pool          txsend.CosmosTxPool
	logger        zerolog.Logger
	clock         txsend.Clock
	gasAdjustment float64
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
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

// AccountInfo returns the on-chain account number and sequence for address.
//
// implemented in bead asg-pvd.3
func (b *Broadcaster) AccountInfo(ctx context.Context, address string) (uint64, uint64, error) {
	return 0, 0, errors.New("not implemented: bead asg-pvd.3")
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
// broadcast result.
//
// implemented in bead asg-pvd.5
func (b *Broadcaster) Broadcast(ctx context.Context, signedTx []byte, mode txsend.BroadcastMode) (*txsend.BroadcastResult, error) {
	return nil, errors.New("not implemented: bead asg-pvd.5")
}

// WaitForTx blocks until txHash is committed (or ctx is cancelled) and returns
// the final result.
//
// implemented in bead asg-pvd.5
func (b *Broadcaster) WaitForTx(ctx context.Context, txHash string) (*txsend.TxResult, error) {
	return nil, errors.New("not implemented: bead asg-pvd.5")
}
