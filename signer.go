package allora

import (
	"context"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

// Signer abstracts transaction signing so callers can choose between a self-managed
// local key and a Privy-managed wallet whose key lives in the Forge backend.
//
// The interface is intentionally the minimal surface the tx builder needs: a public key
// (to assemble the SignDoc and SignatureV2) and a Sign method over the SignDoc bytes.
// It is satisfied for free by any cosmos cryptotypes.PrivKey — and therefore by the
// local *Wallet (via wallet.PrivKey) — and by *RemoteSigner, which delegates signing to
// the backend over HTTP.
//
// Implementations MUST return a 64-byte, low-S secp256k1 signature computed over
// SHA-256(msg), exactly matching cosmos secp256k1 semantics; otherwise on-chain
// signature verification will fail.
//
// Signer is intentionally transaction-only: callers sign SignDoc bytes with
// hash-then-sign (RemoteSigner sends prehashed=false). It deliberately omits the
// prehashed/digest variant the Python SDK exposes via sign_digest; a Go worker that
// needs to sign a raw 32-byte application digest should use a dedicated API rather than
// widening this tx-signing contract.
type Signer interface {
	// PubKey returns the signer's public key.
	PubKey() cryptotypes.PubKey
	// Sign returns a signature over the given SignDoc bytes.
	Sign(msg []byte) ([]byte, error)
}

// ContextSigner is an optional extension of Signer for implementations whose signing
// performs I/O (e.g. RemoteSigner's HTTP call to the Forge backend). When a Signer also
// implements ContextSigner, the tx builder calls SignWithContext so the caller's
// cancellation and deadline propagate to the in-flight signing operation. Local-key
// signers need not implement it; their Sign is purely CPU-bound.
type ContextSigner interface {
	Signer
	// SignWithContext returns a signature over the SignDoc bytes, honoring ctx for
	// cancellation and deadlines.
	SignWithContext(ctx context.Context, msg []byte) ([]byte, error)
}

// signWithContext signs msg using the signer's context-aware path when it implements
// ContextSigner, falling back to the context-free Sign for local-key signers.
func signWithContext(ctx context.Context, signer Signer, msg []byte) ([]byte, error) {
	if cs, ok := signer.(ContextSigner); ok {
		return cs.SignWithContext(ctx, msg)
	}
	return signer.Sign(msg)
}

// Compile-time assurance that the concrete secp256k1 private key — the type held inside
// *Wallet via wallet.PrivKey — satisfies Signer, so the self-managed path needs no
// adapter. Asserting the concrete type (rather than the cryptotypes.PrivKey interface)
// exercises the exact code path callers hit and would catch a future secp256k1 method
// signature change that an interface-vs-interface check silently allows.
var _ Signer = (*secp256k1.PrivKey)(nil)
