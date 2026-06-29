// Package txsend defines the contract for sending, simulating, and confirming
// Cosmos transactions through a pooled, health-tracked client. It is the seam
// the later beads (asg-pvd.3 account-info, asg-pvd.4 gas estimation,
// asg-pvd.5 broadcast+wait) fill in: the interface and value types live here,
// the concrete pooled implementation lives in txsend/cosmospool.
//
// txsend deliberately depends only on gen/interfaces and the standard library.
// It must NOT import cosmosrpc: cosmosrpc.ClientPool embeds
// gen/wrapper.WrapperClient (which transitively reaches back toward the
// top-level allora package via its constructor callers), and importing it here
// would close a cycle (allora -> cosmosrpc -> txsend -> cosmosrpc). Instead the
// narrow CosmosTxPool interface below captures exactly the two sub-clients
// (Tx and Auth) the broadcaster needs; the wiring layer adapts a
// cosmosrpc.ClientPool to it. Interface-segregation keeps the dependency
// surface minimal and the cycle impossible.
//
// Signing is intentionally NOT re-implemented here. The existing allora.Signer
// (signer.go) and the two-phase build/sign flow (tx.go, tx_builder.go) produce
// signed tx bytes; txsend takes those bytes as opaque input and is responsible
// only for account discovery, gas estimation, broadcast, and confirmation.
// The shared codec (alloracodec.CosmosCodec, codec/registry.go) is what the
// later beads use to (un)marshal SimulateRequest/BroadcastTxRequest payloads
// and to decode TxResponse events into the value types defined below.
package txsend

import (
	"context"
	"time"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// TxBroadcaster is the contract for the send-side of the transaction lifecycle:
// fetch the sender's account info, estimate gas, broadcast a signed tx, and wait
// for it to be committed. Each method is ctx-first so a caller's cancellation
// and deadline flow through every I/O boundary (pool retry, gRPC, polling).
//
// Methods return txsend's own value types (BroadcastResult, TxResult), not the
// raw cosmos SDK response types, so callers depend only on this contract. The
// concrete cosmospool.Broadcaster translates between the cosmos gRPC responses
// and these types.
type TxBroadcaster interface {
	// AccountInfo returns the on-chain account number and sequence for address,
	// used to populate TxParams before building/signing. A later bead (asg-pvd.3)
	// implements this via CosmosTxPool.Auth().Account; sequence must reflect the
	// latest committed state so a signed tx is not immediately stale.
	AccountInfo(ctx context.Context, address string) (accountNumber, sequence uint64, err error)

	// EstimateGas returns the simulated gas usage for unsignedTx (raw tx bytes
	// produced by allora.CreateUnsignedTx / CreateUnsignedSendTx). A later bead
	// (asg-pvd.4) implements this via CosmosTxPool.Tx().Simulate. The returned
	// gas is the simulation's GasUsed; callers typically apply a multiplier
	// before setting TxParams.GasLimit.
	EstimateGas(ctx context.Context, unsignedTx []byte) (gas uint64, err error)

	// Broadcast submits signedTx (raw bytes from allora.SignTransactionWith) in
	// the given mode and returns the immediate broadcast result. A later bead
	// (asg-pvd.5) implements this via CosmosTxPool.Tx().BroadcastTx. The result
	// carries the tx hash and the CheckTx code; a non-zero Code means the tx was
	// rejected before mempool acceptance.
	Broadcast(ctx context.Context, signedTx []byte, mode BroadcastMode) (*BroadcastResult, error)

	// WaitForTx blocks until the transaction identified by txHash is committed
	// (or ctx is cancelled), returning the final result. A later bead
	// (asg-pvd.5) implements this by polling GetTx / GetTxsEvent with backoff
	// gated by the Clock seam; it returns an error only if ctx is cancelled
	// before inclusion, not if the tx committed with a non-zero code (that is
	// reported in TxResult.Code).
	WaitForTx(ctx context.Context, txHash string) (*TxResult, error)
}

// CosmosTxPool is the narrow dependency the broadcaster needs from the cosmos
// client pool: the Tx service (broadcast, simulate, get-tx) and the Auth
// service (account queries). It is an interface-segregated slice of the full
// cosmosrpc.ClientPool / interfaces.CosmosClientPool so that txsend depends on
// the minimum surface and never imports cosmosrpc (which would create a cycle).
// The wiring layer satisfies this by passing a cosmosrpc.ClientPool, which
// embeds gen/wrapper.WrapperClient and therefore exposes Tx() and Auth().
type CosmosTxPool interface {
	// Tx returns the cosmos tx service client used for broadcast/simulate/get-tx.
	Tx() interfaces.TxClient
	// Auth returns the cosmos auth service client used for account queries.
	Auth() interfaces.AuthClient
}

// BroadcastMode selects how Broadcast submits a transaction. It is txsend's
// own type (not the cosmos txtypes.BroadcastMode enum) so the interface package
// stays decoupled from the cosmos SDK types; cosmospool.Broadcaster maps each
// value to the corresponding txtypes.BroadcastMode when it builds the
// BroadcastTxRequest.
type BroadcastMode int

const (
	// BroadcastModeSync submits the tx and waits for the CheckTx (mempool
	// validation) response before returning. Maps to
	// txtypes.BroadcastMode_BROADCAST_MODE_SYNC. This is the recommended mode:
	// it surfaces CheckTx rejections immediately while still returning fast.
	BroadcastModeSync BroadcastMode = iota
	// BroadcastModeAsync submits the tx and returns immediately without waiting
	// for CheckTx. Maps to txtypes.BroadcastMode_BROADCAST_MODE_ASYNC. Use only
	// when the caller will WaitForTx anyway and does not need the CheckTx code.
	BroadcastModeAsync
	// BroadcastModeBlock submits the tx and waits until it is included in a
	// block. Maps to txtypes.BroadcastMode_BROADCAST_MODE_BLOCK, which is
	// deprecated by the cosmos SDK since v0.47; prefer BroadcastModeSync plus
	// WaitForTx. It is retained for parity with the cosmos service enum.
	BroadcastModeBlock
)

// BroadcastResult is the immediate outcome of Broadcast: the tx hash and the
// CheckTx/DeliverTx code reported by the node at broadcast time. It is a
// snapshot of the node's response; it does not imply the tx was committed.
type BroadcastResult struct {
	// TxHash is the hex transaction hash the node assigned (or computed) for
	// the submitted tx; use it with WaitForTx.
	TxHash string
	// Code is the response code; 0 means the tx passed CheckTx acceptance.
	// Non-zero means the tx was rejected (see Codespace and RawLog).
	Code uint32
	// Codespace is the namespace for Code (e.g. "sdk", "bank").
	Codespace string
	// RawLog is the raw human-readable log from the node; non-deterministic.
	RawLog string
}

// TxResult is the final outcome of a committed transaction, returned by
// WaitForTx once the tx is included in a block. Code == 0 means success; any
// non-zero code means the tx was committed but failed during DeliverTx.
type TxResult struct {
	// TxHash is the hex transaction hash.
	TxHash string
	// Height is the block height at which the tx was committed.
	Height int64
	// Code is the DeliverTx response code; 0 means the tx executed successfully.
	Code uint32
	// Codespace is the namespace for Code.
	Codespace string
	// Data is the base64-encoded result data, if any.
	Data string
	// RawLog is the raw human-readable log; non-deterministic.
	RawLog string
	// GasWanted is the gas limit requested by the tx.
	GasWanted int64
	// GasUsed is the gas actually consumed during execution.
	GasUsed int64
	// Timestamp is the RFC3339 timestamp of the block that included the tx.
	Timestamp string
	// Events are the events emitted during tx execution (ante + messages).
	Events []Event
}

// Event is a single typed event emitted during transaction execution, grouped
// by Type with a set of key/value Attributes. Later beads decode these via
// alloracodec to typed proto messages when needed.
type Event struct {
	// Type is the event type string (e.g. "transfer", "message").
	Type string
	// Attributes are the key/value pairs attached to the event.
	Attributes []EventAttribute
}

// EventAttribute is a single key/value pair on an Event.
type EventAttribute struct {
	// Key is the attribute key.
	Key string
	// Value is the attribute value.
	Value string
}

// Clock is the time seam used by WaitForTx's polling loop. Injecting it (via
// cosmospool.WithClock) lets tests drive the poll/backoff deterministically
// with a fake clock instead of real time.Sleep; production wires SystemClock.
type Clock interface {
	// Now returns the current time.
	Now() time.Time
	// After returns a channel that delivers the current time after duration d,
	// mirroring time.After so a select on it can be cancelled by ctx.Done().
	After(d time.Duration) <-chan time.Time
}

// SystemClock is the Clock implementation backed by the time package; it is the
// constructor default for cosmospool.Broadcaster.
type SystemClock struct{}

// Now returns time.Now.
func (SystemClock) Now() time.Time { return time.Now() }

// After returns time.After(d).
func (SystemClock) After(d time.Duration) <-chan time.Time { return time.After(d) }
