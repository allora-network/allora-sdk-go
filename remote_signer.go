package allora

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
)

const (
	// apiKeyHeader is the header carrying a Forge API key on requests to the backend.
	apiKeyHeader = "X-Forge-API-Key"

	// maxResponseBytes caps how much of a backend response the signer buffers, so a broken or
	// malicious endpoint cannot drive unbounded memory use during construction or signing.
	// Wallet-info and signature responses are tiny; 1 MiB is generous headroom. Referenced by
	// sendForgeRequest (and thus do/doWalletRequest) and provisionWalletForTopic.
	maxResponseBytes = 1 << 20
)

// RemoteSignerConfig configures a RemoteSigner.
type RemoteSignerConfig struct {
	// BackendURL is the Forge API base URL, e.g. "https://forge.allora.network".
	BackendURL string
	// APIKey is a Forge API key minted for the user (sent as X-Forge-API-Key).
	APIKey string
	// WalletID is the signing wallet's UUID, as returned by the create endpoint.
	WalletID string
	// HTTPClient is optional; a 30s-timeout client is used when nil. That default 30s timeout is
	// an absolute per-call upper bound (dial+TLS+headers+body): because http.Client.Timeout is
	// independent of the request context, it caps any longer ctx deadline a caller sets, so a
	// caller that needs a budget longer than 30s must supply a custom HTTPClient with Timeout: 0
	// (relying on ctx alone). The SDK installs its own CheckRedirect to guard the Forge API key
	// and SignDoc across redirects: it refuses every redirect, so any CheckRedirect set on a
	// supplied client is replaced — a redirect is never followed and the supplied policy is not
	// consulted.
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
	// masterGranter is the master feegrant granter (allo1...) the backend reported for this
	// wallet, discovered from the wallet-info/provision response. Empty when the backend does
	// not advertise one; resolved into a TxParams.FeeGranter via ResolveFeeGranter.
	masterGranter string
}

var (
	_ Signer        = (*RemoteSigner)(nil)
	_ ContextSigner = (*RemoteSigner)(nil)
)

// NewRemoteSigner builds a RemoteSigner and fetches the wallet's public key and address
// from the backend (needed to assemble the SignDoc before the wallet has transacted).
func NewRemoteSigner(ctx context.Context, cfg RemoteSignerConfig) (*RemoteSigner, error) {
	// http.NewRequestWithContext panics on a nil context; reject it with a normal error so an
	// accidentally-nil ctx cannot crash a worker during signer construction.
	if ctx == nil {
		return nil, fmt.Errorf("ctx must not be nil")
	}
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
	} else if dt, ok := http.DefaultTransport.(*http.Transport); ok {
		// http.DefaultTransport's MaxIdleConnsPerHost is 2, which throttles a worker that fans out
		// many concurrent sign calls to one backend host: each excess in-flight request dials a
		// fresh TCP+TLS connection and closes (rather than pools) it on completion, paying a TLS
		// handshake per call. Clone the default transport — preserving its proxy/TLS/HTTP2 and
		// IdleConnTimeout settings — and raise the per-host idle cap so concurrent sign calls reuse
		// pooled connections. A caller-supplied client keeps its own transport, so operators who
		// care can tune it themselves.
		t := dt.Clone()
		t.MaxIdleConnsPerHost = 64
		guarded.Transport = t
	}
	// Refuse every redirect (rejectRedirect returns http.ErrUseLastResponse) so the Forge API key
	// and SignDoc are never re-sent to a redirect target. This replaces any CheckRedirect on a
	// caller-supplied client: a redirect is never followed (it surfaces as a non-2xx error via
	// sendForgeRequest), so there is no path on which a caller's policy would run.
	guarded.CheckRedirect = rejectRedirect
	return guarded
}

// rejectRedirect refuses every redirect, not just cross-origin ones. Go re-sends the request
// method and body on 307/308 redirects, so following any redirect — including a same-origin
// one — would re-POST the SignDoc bytes (and the X-Forge-API-Key header) to the redirect
// target. Returning http.ErrUseLastResponse stops the client at the 3xx response (do() then
// maps it to a non-2xx error) instead of following it. This matches the reject-all stance of
// the sibling SDKs (allora-sdk-py: allow_redirects=False; allora-sdk-ts: redirect: "error").
func rejectRedirect(_ *http.Request, _ []*http.Request) error {
	return http.ErrUseLastResponse
}

// loopbackHosts is the exact set of hosts for which a plaintext http backend URL is
// permitted. It is intentionally narrower than net.IP.IsLoopback (which accepts the whole
// 127.0.0.0/8 block, e.g. 127.0.0.2) and matches the Python sibling's allowlist
// (allora-sdk-py _LOOPBACK_HOSTS). NOTE: allora-sdk-ts currently accepts a broader set — it also
// treats 0.0.0.0 and the entire 127.0.0.0/8 block as loopback — so the three SDKs are not yet
// byte-for-byte identical here. Narrowing TS down to this canonical set is tracked as a cross-repo
// follow-up; the divergence is local-dev only, since production always uses https.
var loopbackHosts = map[string]bool{
	"localhost": true,
	"127.0.0.1": true,
	"::1":       true,
}

// isLoopbackHost reports whether host is one of the canonical loopback hosts for which a
// plaintext http backend URL is permitted. The host is lower-cased before lookup to match the
// auto-lowercasing the Python sibling gets for free via urllib.parse.hostname — Go's
// url.URL.Hostname preserves input case — so a case-mixed input like "LOCALHOST" or "Localhost"
// is accepted identically across SDKs.
func isLoopbackHost(host string) bool {
	return loopbackHosts[strings.ToLower(host)]
}

// PubKey returns the wallet's secp256k1 public key.
func (rs *RemoteSigner) PubKey() cryptotypes.PubKey { return rs.pubKey }

// AccAddress returns the wallet's allo1... account address.
func (rs *RemoteSigner) AccAddress() sdk.AccAddress { return rs.address }

// Address returns the wallet's allo1... account address as a string.
func (rs *RemoteSigner) Address() string { return rs.address.String() }

// MasterGranter returns the master feegrant granter (allo1...) the backend reported for this
// wallet, or "" when the backend does not advertise one. forge-v2 auto-creates a feegrant from
// this granter to each new signing wallet, so a worker can subsidize its gas without holding
// any ALLO. It is discovered at runtime from the wallet-info/provision response (the
// master_granter field), so a master-wallet rotation does not force consumers to reconfigure.
// Use ResolveFeeGranter to apply the env-override precedence and obtain a TxParams.FeeGranter.
//
// The value is snapshotted once at construction (NewRemoteSigner / NewRemoteSignerForTopic) and is
// not refreshed during the signer's lifetime, so a master-wallet rotation mid-life is picked up
// only by reconstructing the signer. For a long-lived signer that must follow a rotation without a
// restart, set FORGE_MASTER_GRANTER_ADDRESS (honored first by ResolveFeeGranter) as the live override.
func (rs *RemoteSigner) MasterGranter() string { return rs.masterGranter }

// ResolveFeeGranter resolves the fee granter to use for this wallet's transactions, applying
// discovery-with-override precedence:
//
//  1. FORGE_MASTER_GRANTER_ADDRESS (read and parsed via FeeGranterFromEnv) — an explicit
//     operator override, shared with allora-sdk-py / allora-sdk-ts;
//  2. the master granter the backend discovered for this wallet (MasterGranter());
//  3. nil — neither is set, so the signing wallet pays its own fees.
//
// Assign the result straight to TxParams.FeeGranter. It returns an error only when an override
// is set but malformed, or when the backend advertised a master_granter that is not a valid
// bech32 address.
func (rs *RemoteSigner) ResolveFeeGranter() (sdk.AccAddress, error) {
	envGranter, err := FeeGranterFromEnv()
	if err != nil {
		return nil, err
	}
	if envGranter != nil {
		return envGranter, nil
	}
	if rs.masterGranter == "" {
		return nil, nil
	}
	granter, err := sdk.AccAddressFromBech32(rs.masterGranter)
	if err != nil {
		return nil, fmt.Errorf("backend returned invalid master_granter %q: %w", rs.masterGranter, err)
	}
	return granter, nil
}

// ClearAssociation releases this wallet's topic binding on the Forge backend (Forge-side
// bookkeeping only — it does NOT unregister the worker on-chain). Call it before
// re-provisioning the wallet against a new topic or before decommissioning it. This mirrors
// the sibling SDKs (allora-sdk-py ForgeBackendClient.clear_association, allora-sdk-ts
// clearAssociation): POST /api/v1/signing-wallets/{id}/clear-association with no body. A
// non-2xx response (e.g. 404 for an unknown/foreign/already-cleared wallet) is returned as
// an error so the caller decides whether an unbind failure is fatal or best-effort.
//
// Use the standalone ClearWalletAssociation when you only hold the wallet id (e.g. retiring a
// worker) and do not want to construct a RemoteSigner (and pay its wallet-info fetch) first.
func (rs *RemoteSigner) ClearAssociation(ctx context.Context) error {
	return clearAssociation(ctx, rs.httpClient, rs.cfg.BackendURL, rs.cfg.WalletID, rs.cfg.APIKey)
}

// Revoke permanently decommissions this signer's wallet on the Forge backend (DELETE
// /api/v1/signing-wallets/{id}). This is the destructive counterpart to ClearAssociation:
// clearing only unbinds the wallet from its topic (reversible by re-provisioning), whereas
// revoking tears the wallet down for good. It mirrors allora-sdk-ts ForgeRemoteSigner.revoke()
// and the server's RevokeSigningWallet handler. A non-2xx response (e.g. 404 for an
// unknown/foreign/already-revoked wallet) is returned as an error so the caller decides whether
// the failure is fatal or best-effort.
//
// Use the standalone RevokeWallet when you only hold the wallet id and do not want to construct
// a RemoteSigner (no wallet-info fetch) first.
func (rs *RemoteSigner) Revoke(ctx context.Context) error {
	return revokeWallet(ctx, rs.httpClient, rs.cfg.BackendURL, rs.cfg.WalletID, rs.cfg.APIKey)
}

// ClearWalletAssociation releases walletID's topic binding on the Forge backend without first
// constructing a RemoteSigner (no wallet-info fetch). Use it to unbind a wallet you only hold
// the id for — e.g. retiring or rotating a worker — so its topic slot stops counting against
// the per-user wallet cap. It mirrors the sibling SDKs' client-level unbind (allora-sdk-py
// ForgeBackendClient.clear_association, allora-sdk-ts clearAssociation): POST
// /api/v1/signing-wallets/{walletID}/clear-association with no body. A non-2xx response (e.g.
// 404 for an unknown/foreign/already-cleared wallet) is returned as an error so the caller
// decides whether an unbind failure is fatal or best-effort.
//
// Only cfg.BackendURL, cfg.APIKey, and cfg.HTTPClient are read; cfg.WalletID is ignored — the
// walletID parameter is the authoritative identifier for the wallet to unbind.
func ClearWalletAssociation(ctx context.Context, cfg RemoteSignerConfig, walletID string) error {
	client, backendURL, canonicalID, err := prepareWalletByIDCall(cfg, walletID)
	if err != nil {
		return err
	}
	return clearAssociation(ctx, client, backendURL, canonicalID, cfg.APIKey)
}

// RevokeWallet decommissions walletID on the Forge backend (DELETE
// /api/v1/signing-wallets/{walletID}), permanently retiring the managed signing wallet. This is
// the destructive counterpart to ClearWalletAssociation: clearing only unbinds the wallet from
// its topic (and is reversible by re-provisioning), whereas revoking tears the wallet down for
// good. Use it when you only hold the wallet id — e.g. decommissioning a worker — without
// constructing a RemoteSigner (no wallet-info fetch). It mirrors the server's RevokeSigningWallet
// handler. A non-2xx response (e.g. 404 for an unknown/foreign/already-revoked wallet) is
// returned as an error so the caller decides whether the failure is fatal or best-effort.
//
// Only cfg.BackendURL, cfg.APIKey, and cfg.HTTPClient are read; cfg.WalletID is ignored — the
// walletID parameter is the authoritative identifier for the wallet to revoke.
func RevokeWallet(ctx context.Context, cfg RemoteSignerConfig, walletID string) error {
	client, backendURL, canonicalID, err := prepareWalletByIDCall(cfg, walletID)
	if err != nil {
		return err
	}
	return revokeWallet(ctx, client, backendURL, canonicalID, cfg.APIKey)
}

// prepareWalletByIDCall validates the config shared by the standalone by-id wallet operations
// (ClearWalletAssociation, RevokeWallet), normalizes the backend URL, and validates +
// canonicalizes walletID so a malformed value cannot inject path segments or query strings. It
// returns a guarded HTTP client, the normalized backend URL, and the canonical wallet id.
func prepareWalletByIDCall(cfg RemoteSignerConfig, walletID string) (client *http.Client, backendURL, canonicalID string, err error) {
	if cfg.BackendURL == "" || cfg.APIKey == "" {
		return nil, "", "", fmt.Errorf("backend URL and API key are required")
	}
	backendURL = strings.TrimRight(cfg.BackendURL, "/")
	if err := requireSecureBackend(backendURL); err != nil {
		return nil, "", "", err
	}
	parsed, err := uuid.Parse(walletID)
	if err != nil {
		return nil, "", "", fmt.Errorf("wallet ID must be a UUID: %w", err)
	}
	return newGuardedClient(cfg.HTTPClient), backendURL, parsed.String(), nil
}

// clearAssociation issues the POST /clear-association call for walletID with the given guarded
// client. walletID must already be canonical. It backs the standalone ClearWalletAssociation
// function and the (*RemoteSigner).ClearAssociation method, delegating transport and 2xx/error
// handling to doWalletRequest (shared with revokeWallet).
func clearAssociation(ctx context.Context, httpClient *http.Client, backendURL, walletID, apiKey string) error {
	endpoint := fmt.Sprintf("%s/api/v1/signing-wallets/%s/clear-association", backendURL, url.PathEscape(walletID))
	return doWalletRequest(ctx, httpClient, http.MethodPost, endpoint, apiKey, "clear-association")
}

// revokeWallet issues the DELETE /signing-wallets/{walletID} call with the given guarded client.
// walletID must already be canonical. It backs the standalone RevokeWallet function, delegating
// transport and 2xx/error handling to doWalletRequest (shared with clearAssociation).
func revokeWallet(ctx context.Context, httpClient *http.Client, backendURL, walletID, apiKey string) error {
	endpoint := fmt.Sprintf("%s/api/v1/signing-wallets/%s", backendURL, url.PathEscape(walletID))
	return doWalletRequest(ctx, httpClient, http.MethodDelete, endpoint, apiKey, "revoke-wallet")
}

// doWalletRequest issues a no-body request to a signing-wallet endpoint and maps the result to an
// error. It is the shared transport behind the clear-association (POST) and revoke/decommission
// (DELETE) calls, which have the same shape: a method + URL carrying only the X-Forge-API-Key
// header, with success defined as any 2xx regardless of body (the endpoints answer 204 No Content
// or 200 + JSON). It shares the bounded body read, drain, and non-2xx handling with do() and the
// provision path via sendForgeRequest, but unlike do() it does not require a JSON body. op labels
// the operation in the request-construction error; ctx must be non-nil (http.NewRequestWithContext
// panics on a nil context).
func doWalletRequest(ctx context.Context, httpClient *http.Client, method, endpoint, apiKey, op string) error {
	if ctx == nil {
		return fmt.Errorf("ctx must not be nil")
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, nil)
	if err != nil {
		return fmt.Errorf("creating %s request: %w", op, err)
	}
	req.Header.Set(apiKeyHeader, apiKey)

	_, _, err = sendForgeRequest(httpClient, req, "reading forge response")
	return err
}

type signingWalletInfoResponse struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	PubKey  string `json:"pubkey"`
	// MasterGranter is the master feegrant granter (allo1...) that subsidizes this wallet's
	// gas. forge-v2 omits it when no master wallet is configured; encoding/json leaves the
	// field empty in that case (and silently ignores it on backends that never send it).
	MasterGranter string `json:"master_granter"`
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
	return rs.applyWalletInfo(info)
}

// applyWalletInfo validates a wallet-info payload against rs.cfg and records the verified
// pubkey and address on the signer. It is shared by NewRemoteSigner (after the GET
// wallet-info call) and NewRemoteSignerForTopic (after the POST provision call), so a
// topic-bound signer is built directly from the provision response without a redundant
// wallet-info GET. rs.cfg.WalletID must already be canonical.
func (rs *RemoteSigner) applyWalletInfo(info signingWalletInfoResponse) error {
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
	// already canonical (canonicalized by the caller).
	returnedID, err := uuid.Parse(info.ID)
	if err != nil {
		return fmt.Errorf("backend returned malformed wallet id %q: %w", info.ID, err)
	}
	if returnedID.String() != rs.cfg.WalletID {
		return fmt.Errorf("backend returned wallet id %q, expected %q", info.ID, rs.cfg.WalletID)
	}
	// hex.DecodeString("") returns (nil, nil), so an omitted/empty pubkey would otherwise
	// surface as the misleading "0-byte pubkey, expected 33". Reject it explicitly, matching
	// the empty-id guard above and the empty-address guard below.
	if info.PubKey == "" {
		return fmt.Errorf("backend returned empty wallet pubkey for wallet %s", rs.cfg.WalletID)
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
	// The master granter is an optional discovery hint, not part of the wallet's identity, so it
	// is validated only when ResolveFeeGranter parses it. Store it verbatim — no trimming — so
	// MasterGranter() returns exactly what the backend sent (as the README documents) and Go stays
	// at parity with allora-sdk-py / allora-sdk-ts, which assign master_granter unmodified. An
	// absent value degrades gracefully to "no granter"; a padded value surfaces as a bech32 parse
	// error in ResolveFeeGranter, identically across the three SDKs.
	rs.masterGranter = info.MasterGranter
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
	// http.NewRequestWithContext panics on a nil context; reject it with a normal error so an
	// accidentally-nil ctx cannot crash a worker mid-signing. (Sign passes context.Background.)
	if ctx == nil {
		return nil, fmt.Errorf("ctx must not be nil")
	}
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
	// Fail closed when the echo is absent: a `sr.PubKey != ""` guard would let a backend
	// that simply drops the field skip this rotation/mis-route check entirely. This is a
	// load-bearing security check, kept at parity with allora-sdk-ts (commit 5cfded6);
	// forge-v2 always returns the echo today.
	if sr.PubKey == "" {
		return nil, fmt.Errorf("backend sign response for wallet %s omitted the pubkey echo; cannot verify the signing key", rs.cfg.WalletID)
	}
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

// sendForgeRequest executes req with httpClient and returns the response together with the
// bounded response body, applying the transport hardening shared by every Forge backend call:
// it caps the buffered body at maxResponseBytes (so a broken or malicious endpoint cannot drive
// unbounded memory use), drains a bounded amount past the cap so the keep-alive connection stays
// reusable, and maps a non-2xx status to an error carrying a truncated body excerpt. On success
// the returned response has already had its Body fully consumed and closed — inspect only its
// status and headers (e.g. Content-Type); the body bytes are returned separately. readErrLabel
// names the body-read error so each call site keeps its specific wording. Callers layer their own
// concerns (Content-Type/JSON checks, body decoding) on top.
func sendForgeRequest(httpClient *http.Client, req *http.Request, readErrLabel string) (*http.Response, []byte, error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("calling forge backend: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", readErrLabel, err)
	}
	// A normal response is fully consumed above, so the keep-alive transport can reuse the
	// connection. If the body exceeded the cap, drain a bounded amount past it so the
	// connection stays reusable without letting a hostile oversized body stream forever;
	// anything beyond that is left for Close to discard (a new connection is acceptable for
	// an abnormal response).
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxResponseBytes))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, nil, fmt.Errorf("forge backend returned status %d: %s", resp.StatusCode, redactSecret(truncateForError(body), req.Header.Get(apiKeyHeader)))
	}
	return resp, body, nil
}

// do executes an HTTP request and returns the body, mapping non-2xx to an error. It layers the
// JSON content-type guard on top of the shared sendForgeRequest transport.
func (rs *RemoteSigner) do(req *http.Request) ([]byte, error) {
	resp, body, err := sendForgeRequest(rs.httpClient, req, "reading forge response")
	if err != nil {
		return nil, err
	}
	// A 2xx with a non-JSON body usually means a captive portal, auth proxy, or
	// misconfigured CDN; surface that clearly instead of an opaque JSON-decode error.
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		return nil, fmt.Errorf("forge backend returned non-JSON response (content-type %q): %s", ct, redactSecret(truncateForError(body), rs.cfg.APIKey))
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

// redactSecret removes every occurrence of secret from s, so a backend or intermediary error body
// echoed into an error string cannot leak the Forge API key. The backend never echoes request
// headers, but a TLS-terminating reverse proxy or captive portal can mirror them into its error
// page; this is defense-in-depth for that operator-deployed case. A blank secret is a no-op.
func redactSecret(s, secret string) string {
	if secret == "" {
		return s
	}
	return strings.ReplaceAll(s, secret, "[REDACTED]")
}

// NewRemoteSignerForTopic idempotently gets-or-creates the user's managed wallet bound to
// topicID (ENGN-8572 "one worker = one topic") and returns a RemoteSigner for it. Safe to call
// on every worker start: the backend enforces one wallet per (user, topic). cfg.WalletID must be
// empty; it is filled from the provision response (a non-empty value is rejected rather than
// silently overwritten).
//
// Each topic binding counts against a per-user wallet cap enforced by the backend. A worker
// that is retired or rotated to a different topic leaves its binding in place, so in
// autoscaling deployments abandoned bindings can accumulate until the cap is reached and new
// provisions fail. Release a binding you no longer need with (*RemoteSigner).ClearAssociation,
// or — when you only hold the wallet id — the standalone ClearWalletAssociation.
func NewRemoteSignerForTopic(ctx context.Context, cfg RemoteSignerConfig, topicID int64, label string) (*RemoteSigner, error) {
	// http.NewRequestWithContext panics on a nil context; reject it with a normal error so an
	// accidentally-nil ctx cannot crash a worker during provisioning.
	if ctx == nil {
		return nil, fmt.Errorf("ctx must not be nil")
	}
	if cfg.BackendURL == "" || cfg.APIKey == "" {
		return nil, fmt.Errorf("backend URL and API key are required")
	}
	// WalletID is filled from the provision response, not supplied by the caller. Reject a
	// non-empty value instead of silently overwriting it, so a config copy-pasted from a
	// NewRemoteSigner call (which requires WalletID) cannot leave the caller believing the signer
	// binds to that stale id when it actually binds to the backend-assigned one.
	if cfg.WalletID != "" {
		return nil, fmt.Errorf("cfg.WalletID must be empty for topic-bound provisioning; it is filled from the provision response")
	}
	if topicID <= 0 {
		return nil, fmt.Errorf("topic id must be a positive integer")
	}
	cfg.BackendURL = strings.TrimRight(cfg.BackendURL, "/")
	if err := requireSecureBackend(cfg.BackendURL); err != nil {
		return nil, err
	}
	// The provision POST already returns the wallet's id, address, and pubkey, so build the
	// signer directly from that response instead of issuing a second wallet-info GET. This
	// halves the round-trips on every topic-bound worker start and matches the provisioning
	// path in the sibling SDKs (allora-sdk-py provision_remote_wallet, allora-sdk-ts
	// provisionForTopic), both of which construct from the provision response.
	// Build the guarded HTTP client once and thread it through both the provision POST and the
	// returned signer, so the provision call's keep-alive connection is available for reuse by
	// later sign requests instead of being allocated, used once, and discarded (mirrors the
	// single requests.Session the Python sibling reuses across wallet operations).
	client := newGuardedClient(cfg.HTTPClient)
	info, err := provisionWalletForTopic(ctx, client, cfg, topicID, label)
	if err != nil {
		// Prefix with the provisioning context so a provision-step failure is distinguishable in
		// logs from a later sign/fetch failure: provisionWalletForTopic and do() otherwise share
		// generic "calling forge backend" / "forge backend returned status" messages, and an
		// operator triaging a failed worker start needs to know whether a backend wallet was
		// created (provision succeeded, something later failed) or not.
		return nil, fmt.Errorf("provisioning wallet for topic %d: %w", topicID, err)
	}
	// cfg.WalletID was ignored on entry; adopt the backend-assigned id, already validated and
	// canonicalized by provisionWalletForTopic, so the applyWalletInfo cross-check and later
	// sign request paths use the canonical form.
	cfg.WalletID = info.ID
	rs := &RemoteSigner{cfg: cfg, httpClient: client}
	if err := rs.applyWalletInfo(info); err != nil {
		return nil, err
	}
	return rs, nil
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
	// A non-root Path is similarly prepended to every request URL (e.g. "https://host/api" ->
	// "https://host/api/api/v1/signing-wallets/..."), silently producing 404s on a
	// correctly-configured backend; reject it for the same reason.
	if base.Path != "" && base.Path != "/" {
		return fmt.Errorf("backend URL must not contain a path (got %q); pass only scheme+host", base.Path)
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
// get-or-create) and returns the full wallet info (id, address, pubkey) the backend reports
// on creation. A static /provision sub-route would collide with /:id in the backend router,
// so provisioning rides on the create endpoint. Returning the full info lets the caller build
// the signer without a second wallet-info GET. It issues the POST with the supplied guarded
// client so the caller can reuse that client (and its keep-alive connection) for the signer.
func provisionWalletForTopic(ctx context.Context, client *http.Client, cfg RemoteSignerConfig, topicID int64, label string) (signingWalletInfoResponse, error) {
	reqBody, err := json.Marshal(provisionWalletRequest{TopicID: topicID, Label: label})
	if err != nil {
		return signingWalletInfoResponse{}, fmt.Errorf("marshaling provision request: %w", err)
	}
	reqURL := cfg.BackendURL + "/api/v1/signing-wallets"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(reqBody))
	if err != nil {
		return signingWalletInfoResponse{}, fmt.Errorf("creating provision request: %w", err)
	}
	req.Header.Set(apiKeyHeader, cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, body, err := sendForgeRequest(client, req, "reading provision response")
	if err != nil {
		return signingWalletInfoResponse{}, err
	}
	// Mirror do(): a 2xx with a non-JSON body usually means a captive portal, auth proxy, or
	// misconfigured CDN. Provisioning is the call most likely to hit an unauthenticated gateway
	// (worker first start), so surface that clearly instead of an opaque JSON-decode error.
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		return signingWalletInfoResponse{}, fmt.Errorf("forge backend returned non-JSON provision response (content-type %q): %s", ct, redactSecret(truncateForError(body), cfg.APIKey))
	}
	var info signingWalletInfoResponse
	if err := json.Unmarshal(body, &info); err != nil {
		return signingWalletInfoResponse{}, fmt.Errorf("decoding provision response: %w", err)
	}
	if info.ID == "" {
		return signingWalletInfoResponse{}, fmt.Errorf("forge backend provision response missing wallet id")
	}
	// Validate (and canonicalize) the backend-assigned id at its source. A buggy backend that
	// returns a non-UUID id (numeric, slug, stale beta-format) is reported as a provision-
	// response fault here, instead of surfacing later as a confusing wallet-ID validation error
	// that blames the caller's config.
	parsed, err := uuid.Parse(info.ID)
	if err != nil {
		return signingWalletInfoResponse{}, fmt.Errorf("forge backend provision response returned non-UUID wallet id %q: %w", info.ID, err)
	}
	info.ID = parsed.String()
	return info, nil
}
