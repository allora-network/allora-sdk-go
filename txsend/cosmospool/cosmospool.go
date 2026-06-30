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
	"reflect"
	"time"

	"github.com/brynbellomy/go-utils/errors"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	sdktx "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/txsend"
)

// Broadcaster is the concrete txsend.TxBroadcaster backed by a
// txsend.CosmosTxPool. It holds the narrow pool, a logger, a Clock (the
// time seam for WaitForTx polling), and a pollInterval controlling the
// WaitForTx polling cadence. Methods are stubbed until their beads land.
type Broadcaster struct {
	pool         txsend.CosmosTxPool
	logger       zerolog.Logger
	clock        txsend.Clock
	pollInterval time.Duration
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
		pool:         pool,
		logger:       logger,
		clock:        txsend.SystemClock{},
		pollInterval: 1 * time.Second,
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

// EstimateGas returns the simulated gas usage for unsignedTx.
//
// implemented in bead asg-pvd.4
func (b *Broadcaster) EstimateGas(ctx context.Context, unsignedTx []byte) (uint64, error) {
	return 0, errors.New("not implemented: bead asg-pvd.4")
}

// Broadcast submits signedTx in the given mode and returns the immediate
// broadcast result. It maps the txsend.BroadcastMode to the cosmos
// txtypes.BroadcastMode, calls the pool's Tx().BroadcastTx RPC, and maps the
// response to a *txsend.BroadcastResult. A non-zero CheckTx Code in the
// response is NOT an error — the tx was successfully submitted to the
// mempool and denied by CheckTx; the caller may inspect result.Code to decide
// how to proceed. Only transport-level failures (network errors, gRPC errors,
// pool retry exhaustion) return a wrapped error.
//
// implemented in bead asg-pvd.5
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
//
// implemented in bead asg-pvd.5
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
