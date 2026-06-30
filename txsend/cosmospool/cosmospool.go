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
	"strings"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/brynbellomy/go-utils/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	alloracodec "github.com/allora-network/allora-sdk-go/codec"
	"github.com/allora-network/allora-sdk-go/txsend"
)

// Broadcaster is the concrete txsend.TxBroadcaster backed by a
// txsend.CosmosTxPool. It holds the narrow pool, a logger, and a Clock (the
// time seam for WaitForTx polling). Methods are stubbed until their beads land.
type Broadcaster struct {
	pool   txsend.CosmosTxPool
	logger zerolog.Logger
	clock  txsend.Clock
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
		pool:   pool,
		logger: logger,
		clock:  txsend.SystemClock{},
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

// EstimateGas returns the simulated gas usage for unsignedTx.
//
// implemented in bead asg-pvd.4
func (b *Broadcaster) EstimateGas(ctx context.Context, unsignedTx []byte) (uint64, error) {
	return 0, errors.New("not implemented: bead asg-pvd.4")
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
