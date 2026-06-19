package allora

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
)

// apiKeyHeader is the header carrying a Forge API key on requests to the backend.
const apiKeyHeader = "X-Forge-API-Key"

// RemoteSignerConfig configures a RemoteSigner.
type RemoteSignerConfig struct {
	// BackendURL is the Forge API base URL, e.g. "https://forge.allora.network".
	BackendURL string
	// APIKey is a Forge API key minted for the user (sent as X-Forge-API-Key).
	APIKey string
	// WalletID is the signing wallet's UUID, as returned by the create endpoint.
	WalletID string
	// HTTPClient is optional; a 30s-timeout client is used when nil.
	HTTPClient *http.Client
}

// RemoteSigner signs transactions by delegating to the Forge backend's signing-wallet
// API. The private key never leaves Privy/the backend. It implements Signer, so it can
// be passed to SignTransactionWith exactly like a local wallet.
type RemoteSigner struct {
	cfg        RemoteSignerConfig
	httpClient *http.Client
	pubKey     cryptotypes.PubKey
	address    sdk.AccAddress
}

var (
	_ Signer        = (*RemoteSigner)(nil)
	_ ContextSigner = (*RemoteSigner)(nil)
)

// NewRemoteSigner builds a RemoteSigner and fetches the wallet's public key and address
// from the backend (needed to assemble the SignDoc before the wallet has transacted).
func NewRemoteSigner(ctx context.Context, cfg RemoteSignerConfig) (*RemoteSigner, error) {
	if cfg.BackendURL == "" || cfg.APIKey == "" || cfg.WalletID == "" {
		return nil, fmt.Errorf("backend URL, API key, and wallet ID are all required")
	}
	// Normalize the base URL so a configured trailing slash (or several) does not
	// produce a malformed "...//api/v1/..." request path.
	cfg.BackendURL = strings.TrimRight(cfg.BackendURL, "/")
	// Validate the base URL and require TLS for non-localhost hosts so the Forge API
	// key is never sent over plaintext or to a non-absolute target.
	base, err := url.Parse(cfg.BackendURL)
	if err != nil {
		return nil, fmt.Errorf("invalid backend URL: %w", err)
	}
	if base.Host == "" || (base.Scheme != "http" && base.Scheme != "https") {
		return nil, fmt.Errorf("backend URL must be an absolute http(s) URL")
	}
	if base.Scheme != "https" && !isLoopbackHost(base.Hostname()) {
		return nil, fmt.Errorf("backend URL must use https for non-localhost host %q", base.Hostname())
	}
	// The wallet ID is interpolated into request paths; require it to be a UUID so a
	// malformed value cannot inject path segments or query strings.
	if _, err := uuid.Parse(cfg.WalletID); err != nil {
		return nil, fmt.Errorf("wallet ID must be a UUID: %w", err)
	}
	rs := &RemoteSigner{cfg: cfg, httpClient: newGuardedClient(cfg.HTTPClient)}
	if err := rs.fetchWallet(ctx); err != nil {
		return nil, err
	}
	return rs, nil
}

// newGuardedClient returns the HTTP client used for backend calls. A nil client gets a
// default 30s-timeout client; a caller-supplied client is shallow-copied so the redirect
// policy can be installed without mutating the caller's instance.
func newGuardedClient(c *http.Client) *http.Client {
	guarded := &http.Client{Timeout: 30 * time.Second}
	if c != nil {
		cp := *c
		guarded = &cp
	}
	guarded.CheckRedirect = stripCredentialOnCrossOrigin
	return guarded
}

// stripCredentialOnCrossOrigin removes the Forge API key header before following any
// redirect that changes origin (scheme or host) and bounds the redirect chain. Go's
// default policy only strips its built-in sensitive headers (Authorization, Cookie) on
// cross-host redirects, so a custom credential header would otherwise leak to a
// plaintext or attacker-controlled target via an open redirect.
func stripCredentialOnCrossOrigin(req *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		return fmt.Errorf("stopped after 10 redirects")
	}
	prev := via[len(via)-1]
	if req.URL.Scheme != prev.URL.Scheme || req.URL.Host != prev.URL.Host {
		req.Header.Del(apiKeyHeader)
	}
	return nil
}

// isLoopbackHost reports whether host is localhost or a loopback IP, the only hosts for
// which a plaintext http backend URL is permitted.
func isLoopbackHost(host string) bool {
	if host == "localhost" {
		return true
	}
	if ip := net.ParseIP(host); ip != nil {
		return ip.IsLoopback()
	}
	return false
}

// PubKey returns the wallet's secp256k1 public key.
func (rs *RemoteSigner) PubKey() cryptotypes.PubKey { return rs.pubKey }

// AccAddress returns the wallet's allo1... account address.
func (rs *RemoteSigner) AccAddress() sdk.AccAddress { return rs.address }

// Address returns the wallet's allo1... account address as a string.
func (rs *RemoteSigner) Address() string { return rs.address.String() }

type signingWalletInfoResponse struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	PubKey  string `json:"pubkey"`
}

func (rs *RemoteSigner) fetchWallet(ctx context.Context) error {
	reqURL := fmt.Sprintf("%s/api/v1/signing-wallets/%s", rs.cfg.BackendURL, url.PathEscape(rs.cfg.WalletID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return fmt.Errorf("creating wallet-info request: %w", err)
	}
	req.Header.Set(apiKeyHeader, rs.cfg.APIKey)
	req.Header.Set("Accept", "application/json")

	body, err := rs.do(req)
	if err != nil {
		return err
	}

	var info signingWalletInfoResponse
	if err := json.Unmarshal(body, &info); err != nil {
		return fmt.Errorf("decoding wallet-info response: %w", err)
	}
	// Reject a response describing a different wallet than the one requested, so a
	// misrouted or buggy backend cannot silently bind this signer to the wrong key.
	if info.ID != "" && info.ID != rs.cfg.WalletID {
		return fmt.Errorf("backend returned wallet id %q, expected %q", info.ID, rs.cfg.WalletID)
	}
	pubBytes, err := hex.DecodeString(info.PubKey)
	if err != nil {
		return fmt.Errorf("decoding wallet pubkey: %w", err)
	}
	// secp256k1.PubKey.Address panics on a wrong-length key, so reject a malformed
	// backend pubkey here and return a normal error instead of crashing the worker.
	if len(pubBytes) != secp256k1.PubKeySize {
		return fmt.Errorf("backend returned %d-byte pubkey, expected %d", len(pubBytes), secp256k1.PubKeySize)
	}
	rs.pubKey = &secp256k1.PubKey{Key: pubBytes}

	// Require the backend to report the wallet address and cross-check it against the
	// pubkey-derived address so a misconfigured wallet fails here rather than on
	// broadcast. An empty address would silently disable this integrity check.
	if info.Address == "" {
		return fmt.Errorf("backend returned empty wallet address for wallet %s", rs.cfg.WalletID)
	}
	derived := sdk.AccAddress(rs.pubKey.Address())
	if derived.String() != info.Address {
		return fmt.Errorf("backend address %s does not match pubkey-derived address %s", info.Address, derived.String())
	}
	rs.address = derived
	return nil
}

type signRequest struct {
	Payload   string `json:"payload"`
	Prehashed bool   `json:"prehashed"`
}

type signResponse struct {
	Signature string `json:"signature"`
	PubKey    string `json:"pubkey"`
}

// Sign signs the SignDoc bytes using a background context bounded by the HTTP client
// timeout. Callers with a deadline or cancellation should use SignWithContext; the tx
// builder does this automatically via the ContextSigner interface.
func (rs *RemoteSigner) Sign(msg []byte) ([]byte, error) {
	return rs.SignWithContext(context.Background(), msg)
}

// SignWithContext delegates signing of the SignDoc bytes to the backend, honoring ctx
// for cancellation and deadlines. The backend SHA-256 hashes the payload and signs it
// with the Privy wallet, returning the 64-byte signature.
func (rs *RemoteSigner) SignWithContext(ctx context.Context, msg []byte) ([]byte, error) {
	reqBody, err := json.Marshal(signRequest{Payload: hex.EncodeToString(msg), Prehashed: false})
	if err != nil {
		return nil, fmt.Errorf("marshaling sign request: %w", err)
	}
	reqURL := fmt.Sprintf("%s/api/v1/signing-wallets/%s/sign", rs.cfg.BackendURL, url.PathEscape(rs.cfg.WalletID))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("creating sign request: %w", err)
	}
	req.Header.Set(apiKeyHeader, rs.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	body, err := rs.do(req)
	if err != nil {
		return nil, err
	}
	var sr signResponse
	if err := json.Unmarshal(body, &sr); err != nil {
		return nil, fmt.Errorf("decoding sign response: %w", err)
	}
	// The cached pubkey is snapshotted at construction and encoded into the tx AuthInfo.
	// If the backend rotated the Privy key (or returned another wallet's key), the cached
	// key would no longer match the signing key and broadcast would fail with an opaque
	// "signature verification failed". Detect that here with an actionable error.
	if sr.PubKey != "" {
		respPub, err := hex.DecodeString(sr.PubKey)
		if err != nil {
			return nil, fmt.Errorf("decoding response pubkey: %w", err)
		}
		if !bytes.Equal(respPub, rs.pubKey.Bytes()) {
			return nil, fmt.Errorf("backend signing key rotated for wallet %s; reconstruct the RemoteSigner", rs.cfg.WalletID)
		}
	}
	sig, err := hex.DecodeString(sr.Signature)
	if err != nil {
		return nil, fmt.Errorf("decoding signature: %w", err)
	}
	if len(sig) != 64 {
		return nil, fmt.Errorf("backend returned %d-byte signature, expected 64", len(sig))
	}
	// Verify locally against the cached pubkey. cosmos secp256k1.VerifySignature rejects
	// non-canonical (high-S) signatures and confirms the signature is valid over
	// SHA-256(msg), so a backend low-S normalization regression fails here with a clear
	// error rather than producing a tx the chain rejects opaquely.
	if !rs.pubKey.VerifySignature(msg, sig) {
		return nil, fmt.Errorf("backend signature failed local verification (non-canonical/high-S or wrong key)")
	}
	return sig, nil
}

// maxResponseBytes caps how much of a backend response the signer buffers, so a broken
// or malicious endpoint cannot drive unbounded memory use during construction or
// signing. Wallet-info and signature responses are tiny; 1 MiB is generous headroom.
const maxResponseBytes = 1 << 20

// do executes an HTTP request and returns the body, mapping non-2xx to an error.
func (rs *RemoteSigner) do(req *http.Request) ([]byte, error) {
	resp, err := rs.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling forge backend: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	if err != nil {
		return nil, fmt.Errorf("reading forge response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("forge backend returned status %d: %s", resp.StatusCode, truncateForError(body))
	}
	// A 2xx with a non-JSON body usually means a captive portal, auth proxy, or
	// misconfigured CDN; surface that clearly instead of an opaque JSON-decode error.
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		return nil, fmt.Errorf("forge backend returned non-JSON response (content-type %q): %s", ct, truncateForError(body))
	}
	return body, nil
}

// truncateForError bounds an error-message body excerpt so a large response body does
// not bloat logs or error chains.
func truncateForError(body []byte) string {
	const max = 512
	if len(body) > max {
		return string(body[:max]) + "..."
	}
	return string(body)
}
