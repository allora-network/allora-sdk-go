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
	// HTTPClient is optional; a 30s-timeout client is used when nil. The SDK installs its
	// own CheckRedirect to guard the Forge API key and SignDoc across redirects; any
	// CheckRedirect set on a supplied client is composed after it (run once the SDK guard
	// permits the redirect), not discarded.
	HTTPClient *http.Client
}

// RemoteSigner signs transactions by delegating to the Forge backend's signing-wallet
// API. The private key never leaves Privy/the backend. It implements Signer, so it can
// be passed to SignTransactionWith exactly like a local wallet.
//
// A RemoteSigner is safe for concurrent use by multiple goroutines: every field is set
// once in NewRemoteSigner and only read thereafter, so a single signer can be shared
// across a worker's transactions for its lifetime.
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
	// Validate the base URL with the same guard NewRemoteSignerForTopic applies before its
	// provisioning POST, so both entrypoints reject unsafe/malformed URLs identically (no
	// plaintext key leak, no userinfo Basic-Auth leak, no query/fragment path corruption).
	if err := requireSecureBackend(cfg.BackendURL); err != nil {
		return nil, err
	}
	// The wallet ID is interpolated into request paths; require it to be a UUID so a
	// malformed value cannot inject path segments or query strings. Canonicalize it
	// (lowercase, hyphenated) so a valid but non-canonical input — uppercase, a "urn:uuid:"
	// prefix, or braces — is used consistently in request paths and in the backend wallet-id
	// cross-check, instead of surfacing as a false wallet-id mismatch in fetchWallet.
	parsedWalletID, err := uuid.Parse(cfg.WalletID)
	if err != nil {
		return nil, fmt.Errorf("wallet ID must be a UUID: %w", err)
	}
	cfg.WalletID = parsedWalletID.String()
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
	// Compose with (rather than silently discard) any CheckRedirect the caller installed:
	// run the SDK's credential/cross-origin guard first, then defer to the caller's policy,
	// so a caller with a stricter posture (e.g. http.ErrUseLastResponse) keeps it.
	prevCheck := guarded.CheckRedirect
	guarded.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if err := rejectCrossOriginRedirect(req, via); err != nil {
			return err
		}
		if prevCheck != nil {
			return prevCheck(req, via)
		}
		return nil
	}
	return guarded
}

// rejectCrossOriginRedirect refuses any redirect that changes origin (scheme or host)
// and bounds the redirect chain. Go preserves the method and body on 307/308 redirects,
// so following a cross-origin redirect would re-POST the SignDoc bytes to the redirect
// target; refusing outright (rather than only stripping the API key and continuing)
// prevents that body leak. This matches the reject-all stance of the sibling SDKs
// (allora-sdk-py: allow_redirects=False; allora-sdk-ts: redirect: "error").
func rejectCrossOriginRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		return fmt.Errorf("stopped after 10 redirects")
	}
	prev := via[len(via)-1]
	if req.URL.Scheme != prev.URL.Scheme || req.URL.Host != prev.URL.Host {
		return fmt.Errorf("refusing cross-origin redirect from %s://%s to %s://%s",
			prev.URL.Scheme, prev.URL.Host, req.URL.Scheme, req.URL.Host)
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
	// misrouted or buggy backend cannot silently bind this signer to the wrong key. An
	// empty id must fail rather than skip the check (a backend that omits it is broken by
	// definition) — the same non-empty posture applied to the address field below.
	if info.ID == "" {
		return fmt.Errorf("backend returned empty wallet id for wallet %s", rs.cfg.WalletID)
	}
	// Compare canonicalized UUIDs, not raw text: the backend may return an equivalent but
	// differently-formatted UUID (uppercase, braces, ...) that names the same wallet, and a
	// raw-string compare would surface that as a false wallet-id mismatch. rs.cfg.WalletID is
	// already canonical (canonicalized in NewRemoteSigner).
	returnedID, err := uuid.Parse(info.ID)
	if err != nil {
		return fmt.Errorf("backend returned malformed wallet id %q: %w", info.ID, err)
	}
	if returnedID.String() != rs.cfg.WalletID {
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
		// Length-check before comparing: an alternate but valid encoding of the *same* key
		// (uncompressed 65-byte SEC1, amino-prefixed, ...) would fail bytes.Equal and be
		// misreported as a key rotation, sending users down a wrong debugging path.
		if len(respPub) != secp256k1.PubKeySize {
			return nil, fmt.Errorf("backend returned %d-byte pubkey, expected %d-byte compressed secp256k1", len(respPub), secp256k1.PubKeySize)
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
	// A normal response is fully consumed above, so the keep-alive transport can reuse the
	// connection. If the body exceeded the cap, drain a bounded amount past it so the
	// connection stays reusable without letting a hostile oversized body stream forever;
	// anything beyond that is left for Close to discard (a new connection is acceptable for
	// an abnormal response).
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxResponseBytes))
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

// NewRemoteSignerForTopic idempotently gets-or-creates the user's managed wallet bound to
// topicID (ENGN-8572 "one worker = one topic") and returns a RemoteSigner for it. Safe to call
// on every worker start: the backend enforces one wallet per (user, topic). cfg.WalletID is
// ignored; it is filled from the provision response.
func NewRemoteSignerForTopic(ctx context.Context, cfg RemoteSignerConfig, topicID int64, label string) (*RemoteSigner, error) {
	if cfg.BackendURL == "" || cfg.APIKey == "" {
		return nil, fmt.Errorf("backend URL and API key are required")
	}
	if topicID <= 0 {
		return nil, fmt.Errorf("topic id must be a positive integer")
	}
	cfg.BackendURL = strings.TrimRight(cfg.BackendURL, "/")
	if err := requireSecureBackend(cfg.BackendURL); err != nil {
		return nil, err
	}
	walletID, err := provisionWalletForTopic(ctx, cfg, topicID, label)
	if err != nil {
		return nil, err
	}
	cfg.WalletID = walletID
	// NewRemoteSigner re-validates the config and fetches/cross-checks the wallet pubkey+address.
	return NewRemoteSigner(ctx, cfg)
}

// requireSecureBackend validates that backendURL is a safe Forge API base URL: an absolute
// http(s) URL with no query string, fragment, or embedded userinfo, requiring TLS for
// non-loopback hosts. It is the single validator shared by NewRemoteSigner and
// NewRemoteSignerForTopic so the first provisioning POST is held to exactly the same
// standard as signer construction (the Forge API key is never sent in cleartext, emitted
// as Basic Auth, or glued onto a request path).
func requireSecureBackend(backendURL string) error {
	base, err := url.Parse(backendURL)
	if err != nil {
		return fmt.Errorf("invalid backend URL: %w", err)
	}
	if base.Host == "" || (base.Scheme != "http" && base.Scheme != "https") {
		return fmt.Errorf("backend URL must be an absolute http(s) URL")
	}
	// Request paths are built by string concatenation (fmt.Sprintf), so a query string or
	// fragment on the base URL would be glued onto the path and corrupt every request
	// (e.g. "https://host/?token=x" -> "https://host/?token=x/api/v1/...").
	if base.RawQuery != "" || base.Fragment != "" {
		return fmt.Errorf("backend URL must not contain a query string or fragment")
	}
	// Reject embedded userinfo (e.g. "https://user:pass@host"): net/http would emit Basic
	// Auth on every request alongside the X-Forge-API-Key, and the credentials would leak
	// into *url.Error strings (and thus logs) on any network failure.
	if base.User != nil {
		return fmt.Errorf("backend URL must not contain userinfo")
	}
	if base.Scheme != "https" && !isLoopbackHost(base.Hostname()) {
		return fmt.Errorf("backend URL must use https for non-localhost host %q", base.Hostname())
	}
	return nil
}

// provisionWalletRequest is the POST body for topic wallet provisioning. A typed struct
// (matching signRequest/signResponse/signingWalletInfoResponse) keeps the wire contract
// visible and gives the field names compile-time safety; the omitempty tag reproduces the
// previous "include label only when non-empty" behavior.
type provisionWalletRequest struct {
	TopicID int64  `json:"topic_id"`
	Label   string `json:"label,omitempty"`
}

// provisionWalletForTopic POSTs to /api/v1/signing-wallets with a topic_id (idempotent
// get-or-create) and returns the wallet's UUID. A static /provision sub-route would collide
// with /:id in the backend router, so provisioning rides on the create endpoint.
func provisionWalletForTopic(ctx context.Context, cfg RemoteSignerConfig, topicID int64, label string) (string, error) {
	reqBody, err := json.Marshal(provisionWalletRequest{TopicID: topicID, Label: label})
	if err != nil {
		return "", fmt.Errorf("marshaling provision request: %w", err)
	}
	reqURL := cfg.BackendURL + "/api/v1/signing-wallets"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("creating provision request: %w", err)
	}
	req.Header.Set(apiKeyHeader, cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := newGuardedClient(cfg.HTTPClient).Do(req)
	if err != nil {
		return "", fmt.Errorf("calling forge backend: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	if err != nil {
		return "", fmt.Errorf("reading provision response: %w", err)
	}
	// Mirror do(): drain a bounded amount past the read cap so a slightly-oversized 2xx body
	// still leaves the keep-alive connection reusable instead of being torn down on Close.
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxResponseBytes))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("forge backend returned status %d: %s", resp.StatusCode, truncateForError(body))
	}
	// Mirror do(): a 2xx with a non-JSON body usually means a captive portal, auth proxy, or
	// misconfigured CDN. Provisioning is the call most likely to hit an unauthenticated gateway
	// (worker first start), so surface that clearly instead of an opaque JSON-decode error.
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		return "", fmt.Errorf("forge backend returned non-JSON provision response (content-type %q): %s", ct, truncateForError(body))
	}
	var info signingWalletInfoResponse
	if err := json.Unmarshal(body, &info); err != nil {
		return "", fmt.Errorf("decoding provision response: %w", err)
	}
	if info.ID == "" {
		return "", fmt.Errorf("forge backend provision response missing wallet id")
	}
	return info.ID, nil
}
