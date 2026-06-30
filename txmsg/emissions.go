package txmsg

import (
	"github.com/brynbellomy/go-utils/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	emissionstypes "github.com/allora-network/allora-chain/x/emissions/types"
)

// RegisterWorkerParams carries the inputs to NewRegisterWorker.
//
// The brief's old draft (the deleted untracked tx/{bank,emissions,feegrant}_tx.go)
// used an analogous struct for the same purpose. The same shape is kept here
// so callers that were already constructing that draft struct can switch
// to txmsg without a rename.
//
// Sender is the address paying the tx fee; it is required. Owner is the
// actor being registered for the topic (a worker or a reputer, depending on
// IsReputer). On Allora's emissions module the chain checks that Sender
// and Owner are equal (only the actor itself can register) — that
// constraint is the chain's, not the constructor's, and is intentionally
// NOT mirrored here so the constructor stays a pure message builder.
type RegisterWorkerParams struct {
	Sender    string
	TopicId   uint64
	Owner     string
	IsReputer bool
}

// InsertWorkerPayloadParams carries the inputs to NewInsertWorkerPayload.
//
// The Sender is the worker paying the tx fee. WorkerDataBundle is the
// pre-built emissions worker bundle (which already carries the worker
// address, topic id, nonce, inference/forecast payload, and the
// application-level signature that authenticates the bundle's contents).
// The constructor does NOT re-validate the signature — the chain
// verifies it; this constructor only enforces that the bundle is
// non-nil, that its nonce is non-nil, and that its signature is
// non-empty (mirroring the old untracked draft's pre-flight guard so a
// missing signature is caught here, not on chain).
type InsertWorkerPayloadParams struct {
	Sender           string
	WorkerDataBundle *emissionstypes.InputWorkerDataBundle
}

// NewRegisterWorker constructs an emissions module RegisterRequest
// (the chain's "register as worker or reputer for a topic" message).
//
// The returned sdk.Msg is the concrete *emissionstypes.RegisterRequest —
// note the Go name: the v9 types package dropped the historical "Msg"
// prefix on these types (the older API versions, e.g. emissionsv2 and
// emissionsv3, named it MsgRegister; the live/current emissionsv9
// package, which the shared codec in codec/registry.go registers
// alongside the older API versions, names it RegisterRequest). The
// constructor name "NewRegisterWorker" follows the SDK convention of
// action-oriented names; the underlying Go type follows the v9
// package's naming. If a future version bump renames the type again,
// update this return type and the codec registration together.
func NewRegisterWorker(p RegisterWorkerParams) (sdk.Msg, error) {
	if _, err := validateBech32(p.Sender, "sender"); err != nil {
		return nil, err
	}
	if p.TopicId == 0 {
		return nil, errors.New("topic id must be greater than 0")
	}
	if _, err := validateBech32(p.Owner, "owner"); err != nil {
		return nil, err
	}
	return &emissionstypes.RegisterRequest{
		Sender:    p.Sender,
		TopicId:   p.TopicId,
		Owner:     p.Owner,
		IsReputer: p.IsReputer,
	}, nil
}

// NewInsertWorkerPayload constructs an emissions module
// InsertWorkerPayloadRequest (the chain's "submit a worker inference/
// forecast payload" message).
//
// As with NewRegisterWorker, the underlying Go type is from the v9
// emissions types package; historically this was MsgInsertWorkerPayload
// in the v2/v3 API versions, but the live/current v9 types package
// uses the request/response naming convention.
//
// The validation here mirrors the old untracked tx/emissions_tx.go
// draft's pre-flight guards: bundle non-nil, bundle.Nonce non-nil,
// and InferencesForecastsBundleSignature non-empty. The chain would
// reject any of these on its own (an empty signature is a hard
// "signature required" error on the emissions module), but catching
// them up front saves a broadcast round-trip and keeps the error
// message next to the call site.
func NewInsertWorkerPayload(p InsertWorkerPayloadParams) (sdk.Msg, error) {
	if _, err := validateBech32(p.Sender, "sender"); err != nil {
		return nil, err
	}
	if p.WorkerDataBundle == nil {
		return nil, errors.New("worker data bundle is required")
	}
	if p.WorkerDataBundle.Nonce == nil {
		return nil, errors.New("worker data bundle nonce is required")
	}
	if len(p.WorkerDataBundle.InferencesForecastsBundleSignature) == 0 {
		return nil, errors.New("worker data bundle signature is required")
	}

	return &emissionstypes.InsertWorkerPayloadRequest{
		Sender:           p.Sender,
		WorkerDataBundle: p.WorkerDataBundle,
	}, nil
}
