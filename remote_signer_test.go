package allora

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/stretchr/testify/require"
)

// newFakeForgeBackend mimics the Forge signing-wallet API, signing with the given local
// wallet so tests can assert that remote signing equals local signing.
func newFakeForgeBackend(t *testing.T, wallet *Wallet, walletID string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/signing-wallets/"+walletID, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.NotEmpty(t, r.Header.Get(apiKeyHeader))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"id":      walletID,
			"address": wallet.GetAddress(),
			"pubkey":  hex.EncodeToString(wallet.GetPublicKeyBytes()),
		})
	})

	mux.HandleFunc("/api/v1/signing-wallets/"+walletID+"/sign", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.NotEmpty(t, r.Header.Get(apiKeyHeader))

		var body signRequest
		raw, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.NoError(t, json.Unmarshal(raw, &body))
		require.False(t, body.Prehashed, "tx signing must request hash-then-sign")

		payload, err := hex.DecodeString(body.Payload)
		require.NoError(t, err)
		// cosmos secp256k1.PrivKey.Sign computes sha256(payload) then signs, exactly what
		// the real backend does via Privy's raw sign over sha256(SignDoc).
		sig, err := wallet.PrivKey.Sign(payload)
		require.NoError(t, err)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"signature": hex.EncodeToString(sig),
			"pubkey":    hex.EncodeToString(wallet.GetPublicKeyBytes()),
		})
	})
	return httptest.NewServer(mux)
}

func TestRemoteSigner_MatchesLocalSigning(t *testing.T) {
	wallet, err := NewWalletFromMnemonic(testMnemonic, DefaultHDPath)
	require.NoError(t, err)

	const walletID = "11111111-1111-1111-1111-111111111111"
	srv := newFakeForgeBackend(t, wallet, walletID)
	defer srv.Close()

	rs, err := NewRemoteSigner(context.Background(), RemoteSignerConfig{
		BackendURL: srv.URL,
		APIKey:     "forge_sk_test",
		WalletID:   walletID,
	})
	require.NoError(t, err)
	require.Equal(t, wallet.GetAddress(), rs.Address())
	require.Equal(t, wallet.GetPublicKeyBytes(), rs.PubKey().Bytes())

	amount := sdk.NewCoins(sdk.NewInt64Coin("uallo", 1000))
	params := &TxParams{
		ChainID:       "allora-testnet-1",
		AccountNumber: 7,
		Sequence:      3,
		GasLimit:      200000,
		FeeAmount:     sdk.NewCoins(sdk.NewInt64Coin("uallo", 5000)),
	}
	unsigned, err := CreateUnsignedSendTx(wallet.Address, wallet.Address, amount, params)
	require.NoError(t, err)

	localSigned, err := SignTransaction(unsigned, wallet, params)
	require.NoError(t, err)
	remoteSigned, err := SignTransactionWith(context.Background(), unsigned, rs, params)
	require.NoError(t, err)

	// The fake backend signs with the same local key and secp256k1 is deterministic
	// (RFC 6979), so byte-equality holds here. The real Privy backend may choose a
	// different (still-valid) nonce, so the contract check that survives a
	// non-deterministic backend is that the remote-signed tx carries a signature that
	// verifies against the signer's pubkey.
	require.Equal(t, localSigned, remoteSigned)
	assertTxSignatureVerifies(t, remoteSigned, params)
}

// assertTxSignatureVerifies decodes a signed tx, recomputes the SignDoc bytes from the
// embedded signer info, and verifies the signature against the signer's pubkey. Unlike a
// byte-equality check against local signing, this passes for any valid ECDSA signature
// regardless of the nonce the backend chose.
func assertTxSignatureVerifies(t *testing.T, signedTx []byte, params *TxParams) {
	t.Helper()
	b := newTxBuilder()
	decoded, err := b.txConfig.TxDecoder()(signedTx)
	require.NoError(t, err)

	sigTx, ok := decoded.(authsigning.Tx)
	require.True(t, ok, "decoded tx must expose signatures")
	sigs, err := sigTx.GetSignaturesV2()
	require.NoError(t, err)
	require.Len(t, sigs, 1)

	pubKey := sigs[0].PubKey
	require.NotNil(t, pubKey)
	single, ok := sigs[0].Data.(*signingtypes.SingleSignatureData)
	require.True(t, ok)

	signMode, err := authsigning.APISignModeToInternal(b.txConfig.SignModeHandler().DefaultMode())
	require.NoError(t, err)
	signerData := authsigning.SignerData{
		ChainID:       params.ChainID,
		AccountNumber: params.AccountNumber,
		Sequence:      params.Sequence,
		Address:       sdk.AccAddress(pubKey.Address()).String(),
		PubKey:        pubKey,
	}
	signBytes, err := authsigning.GetSignBytesAdapter(
		context.Background(),
		b.txConfig.SignModeHandler(),
		signMode,
		signerData,
		decoded,
	)
	require.NoError(t, err)
	require.True(t, pubKey.VerifySignature(signBytes, single.Signature),
		"remote-signed tx signature must verify against the signer pubkey")
}

func TestNewRemoteSigner_RequiresConfig(t *testing.T) {
	_, err := NewRemoteSigner(context.Background(), RemoteSignerConfig{})
	require.Error(t, err)
}

// TestNewRemoteSigner_CanonicalizesWalletID pins that a valid but non-canonical UUID
// (here upper-case) is canonicalized at construction, so the request path and the backend
// wallet-id cross-check both use the canonical lower-case form rather than producing a
// false "wallet id mismatch".
func TestNewRemoteSigner_CanonicalizesWalletID(t *testing.T) {
	wallet, err := NewWalletFromMnemonic(testMnemonic, DefaultHDPath)
	require.NoError(t, err)

	// The backend is keyed by (and reports) the canonical lower-case UUID...
	const canonicalID = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
	srv := newFakeForgeBackend(t, wallet, canonicalID)
	defer srv.Close()

	// ...but the caller supplies the upper-case equivalent.
	rs, err := NewRemoteSigner(context.Background(), RemoteSignerConfig{
		BackendURL: srv.URL,
		APIKey:     "forge_sk_test",
		WalletID:   "AAAAAAAA-BBBB-CCCC-DDDD-EEEEEEEEEEEE",
	})
	require.NoError(t, err)
	require.Equal(t, wallet.GetAddress(), rs.Address())
}

// TestNewRemoteSignerForTopic_BuildsFromProvisionResponse pins that the topic constructor
// builds the signer from the POST provision response (which carries id, address, and pubkey)
// without a second wallet-info GET — matching the allora-sdk-py / allora-sdk-ts provisioning
// path. A GET to /signing-wallets/{id} would be the redundant round-trip this removes, so the
// GET handler fails the test if it is ever hit.
func TestNewRemoteSignerForTopic_BuildsFromProvisionResponse(t *testing.T) {
	wallet, err := NewWalletFromMnemonic(testMnemonic, DefaultHDPath)
	require.NoError(t, err)

	const walletID = "44444444-4444-4444-4444-444444444444"
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/signing-wallets", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.NotEmpty(t, r.Header.Get(apiKeyHeader))
		var body provisionWalletRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Equal(t, int64(42), body.TopicID)
		require.Equal(t, "worker-label", body.Label)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"id":      walletID,
			"address": wallet.GetAddress(),
			"pubkey":  hex.EncodeToString(wallet.GetPublicKeyBytes()),
		})
	})
	mux.HandleFunc("/api/v1/signing-wallets/"+walletID, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected redundant wallet-info round-trip: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	rs, err := NewRemoteSignerForTopic(context.Background(), RemoteSignerConfig{
		BackendURL: srv.URL,
		APIKey:     "forge_sk_test",
		// WalletID intentionally omitted: it is filled from the provision response.
	}, 42, "worker-label")
	require.NoError(t, err)
	require.Equal(t, walletID, rs.cfg.WalletID)
	require.Equal(t, wallet.GetAddress(), rs.Address())
	require.Equal(t, wallet.GetPublicKeyBytes(), rs.PubKey().Bytes())
}

// TestNewRemoteSignerForTopic_RejectsNonEmptyWalletID pins that the topic constructor rejects a
// caller-supplied cfg.WalletID (synth-005): the field is filled from the provision response, so a
// stale value copy-pasted from a NewRemoteSigner config must fail loudly rather than be silently
// overwritten and bind the signer to an unexpected wallet.
func TestNewRemoteSignerForTopic_RejectsNonEmptyWalletID(t *testing.T) {
	_, err := NewRemoteSignerForTopic(context.Background(), RemoteSignerConfig{
		BackendURL: "https://forge.allora.network",
		APIKey:     "forge_sk_test",
		WalletID:   "44444444-4444-4444-4444-444444444444",
	}, 42, "")
	require.Error(t, err)
	require.ErrorContains(t, err, "cfg.WalletID must be empty")
}

// TestNewRemoteSignerForTopic_RejectsNonUUIDProvisionID pins that a backend returning a
// non-UUID wallet id from the provision POST is rejected at the provision step with a clear
// provision-response error, rather than surfacing later as a confusing wallet-ID validation
// failure attributed to the caller's config.
func TestNewRemoteSignerForTopic_RejectsNonUUIDProvisionID(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/signing-wallets", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"id":      "12345", // not a UUID
			"address": "allo1placeholder",
			"pubkey":  "abcd",
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	_, err := NewRemoteSignerForTopic(context.Background(), RemoteSignerConfig{
		BackendURL: srv.URL,
		APIKey:     "forge_sk_test",
	}, 42, "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "non-UUID wallet id")
}

// TestNewRemoteSignerForTopic_ProvisionErrorIdentifiesStep pins that a provision-step failure
// is wrapped with the topic-provisioning context, so it is distinguishable in logs from a later
// sign/fetch failure that shares the generic "forge backend returned status" message.
func TestNewRemoteSignerForTopic_ProvisionErrorIdentifiesStep(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/signing-wallets", func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	_, err := NewRemoteSignerForTopic(context.Background(), RemoteSignerConfig{
		BackendURL: srv.URL,
		APIKey:     "forge_sk_test",
	}, 42, "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "provisioning wallet for topic 42")
}

// TestSignTransactionWith_NilContextDoesNotPanic pins that a nil ctx is defaulted to
// context.Background() instead of panicking in http.NewRequestWithContext when the signer is
// a RemoteSigner that performs I/O during signing.
func TestSignTransactionWith_NilContextDoesNotPanic(t *testing.T) {
	wallet, err := NewWalletFromMnemonic(testMnemonic, DefaultHDPath)
	require.NoError(t, err)

	const walletID = "11111111-1111-1111-1111-111111111111"
	srv := newFakeForgeBackend(t, wallet, walletID)
	defer srv.Close()

	rs, err := NewRemoteSigner(context.Background(), RemoteSignerConfig{
		BackendURL: srv.URL,
		APIKey:     "forge_sk_test",
		WalletID:   walletID,
	})
	require.NoError(t, err)

	amount := sdk.NewCoins(sdk.NewInt64Coin("uallo", 1000))
	params := &TxParams{
		ChainID:       "allora-testnet-1",
		AccountNumber: 7,
		Sequence:      3,
		GasLimit:      200000,
		FeeAmount:     sdk.NewCoins(sdk.NewInt64Coin("uallo", 5000)),
	}
	unsigned, err := CreateUnsignedSendTx(wallet.Address, wallet.Address, amount, params)
	require.NoError(t, err)

	var nilCtx context.Context // explicitly nil to exercise the guard
	signed, err := SignTransactionWith(nilCtx, unsigned, rs, params)
	require.NoError(t, err)
	assertTxSignatureVerifies(t, signed, params)
}

// TestRemoteSigner_ConcurrentSign pins the documented goroutine-safety guarantee: a
// single RemoteSigner shared across goroutines must produce a valid signature on every
// call. Run with -race to catch any future field that turns the struct into a data race.
func TestRemoteSigner_ConcurrentSign(t *testing.T) {
	wallet, err := NewWalletFromMnemonic(testMnemonic, DefaultHDPath)
	require.NoError(t, err)

	const walletID = "11111111-1111-1111-1111-111111111111"
	srv := newFakeForgeBackend(t, wallet, walletID)
	defer srv.Close()

	rs, err := NewRemoteSigner(context.Background(), RemoteSignerConfig{
		BackendURL: srv.URL,
		APIKey:     "forge_sk_test",
		WalletID:   walletID,
	})
	require.NoError(t, err)

	msg := []byte("allora sign doc bytes")
	// The fake backend signs deterministically (RFC 6979) with the same key, so every
	// concurrent result must equal this reference signature.
	want, err := rs.Sign(msg)
	require.NoError(t, err)

	const n = 16
	var wg sync.WaitGroup
	errs := make(chan error, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := rs.Sign(msg)
			if err != nil {
				errs <- err
				return
			}
			if !bytes.Equal(got, want) {
				errs <- fmt.Errorf("concurrent signature mismatch")
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}
}

const negWalletID = "22222222-2222-2222-2222-222222222222"

// walletInfoJSON returns a GET handler that responds with the given JSON object.
func walletInfoJSON(m map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(m)
	}
}

// serveBackend starts a fake backend with the given wallet-info and (optional) sign
// handlers, registered under negWalletID, and closes it at test end.
func serveBackend(t *testing.T, walletInfo, sign http.HandlerFunc) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/signing-wallets/"+negWalletID, walletInfo)
	if sign != nil {
		mux.HandleFunc("/api/v1/signing-wallets/"+negWalletID+"/sign", sign)
	}
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func connectRemoteSigner(t *testing.T, baseURL string) (*RemoteSigner, error) {
	t.Helper()
	return NewRemoteSigner(context.Background(), RemoteSignerConfig{
		BackendURL: baseURL,
		APIKey:     "forge_sk_test",
		WalletID:   negWalletID,
	})
}

func TestNewRemoteSigner_RejectsNonUUIDWalletID(t *testing.T) {
	_, err := NewRemoteSigner(context.Background(), RemoteSignerConfig{
		BackendURL: "https://forge.allora.network",
		APIKey:     "forge_sk_test",
		WalletID:   "not-a-uuid",
	})
	require.Error(t, err)
}

func TestNewRemoteSigner_RejectsPlaintextNonLocalhost(t *testing.T) {
	_, err := NewRemoteSigner(context.Background(), RemoteSignerConfig{
		BackendURL: "http://forge.allora.network",
		APIKey:     "forge_sk_test",
		WalletID:   negWalletID,
	})
	require.Error(t, err)
}

// TestNewRemoteSigner_RejectsBackendURLWithPath pins that a backend URL carrying a path is
// rejected (synth-008): request paths are built by concatenation, so a base like
// "https://host/api" would silently produce "https://host/api/api/v1/..." 404s that look correct
// at a glance.
func TestNewRemoteSigner_RejectsBackendURLWithPath(t *testing.T) {
	_, err := NewRemoteSigner(context.Background(), RemoteSignerConfig{
		BackendURL: "https://forge.allora.network/api",
		APIKey:     "forge_sk_test",
		WalletID:   negWalletID,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "must not contain a path")
}

// TestNewRemoteSigner_RejectsBackendURLWithQuery pins that a backend URL carrying a query
// string is rejected: request paths are built by concatenation, so a query on the base URL
// would be glued onto every request path (e.g. "https://host/?token=x/api/v1/...") and corrupt it.
func TestNewRemoteSigner_RejectsBackendURLWithQuery(t *testing.T) {
	_, err := NewRemoteSigner(context.Background(), RemoteSignerConfig{
		BackendURL: "https://forge.allora.network/?token=x",
		APIKey:     "forge_sk_test",
		WalletID:   negWalletID,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "query string or fragment")
}

// TestNewRemoteSigner_RejectsBackendURLWithFragment pins that a fragment on the backend URL is
// rejected for the same path-corruption reason as a query string.
func TestNewRemoteSigner_RejectsBackendURLWithFragment(t *testing.T) {
	_, err := NewRemoteSigner(context.Background(), RemoteSignerConfig{
		BackendURL: "https://forge.allora.network/#frag",
		APIKey:     "forge_sk_test",
		WalletID:   negWalletID,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "query string or fragment")
}

// TestNewRemoteSigner_RejectsBackendURLWithUserinfo pins that embedded userinfo is rejected:
// net/http would otherwise emit Basic Auth on every request alongside X-Forge-API-Key, and the
// credentials would leak into *url.Error strings (and thus logs) on any network failure.
func TestNewRemoteSigner_RejectsBackendURLWithUserinfo(t *testing.T) {
	_, err := NewRemoteSigner(context.Background(), RemoteSignerConfig{
		BackendURL: "https://user:pass@forge.allora.network",
		APIKey:     "forge_sk_test",
		WalletID:   negWalletID,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "userinfo")
}

// TestIsLoopbackHost_NarrowedToCanonicalSet pins that the plaintext-http allowlist is the
// canonical {localhost, 127.0.0.1, ::1} set shared with allora-sdk-py, not the whole
// 127.0.0.0/8 block that net.IP.IsLoopback accepts. This keeps the cross-SDK policy aligned.
func TestIsLoopbackHost_NarrowedToCanonicalSet(t *testing.T) {
	// Mixed-case variants must also be accepted: Go's url.Hostname() preserves input case, so
	// isLoopbackHost lower-cases before lookup to stay aligned with the Python sibling.
	for _, h := range []string{"localhost", "127.0.0.1", "::1", "LOCALHOST", "Localhost"} {
		require.True(t, isLoopbackHost(h), "%q must be allowed as loopback", h)
	}
	// 127.0.0.2 is loopback per net.IP.IsLoopback but is rejected by the Python sibling, so it
	// must be rejected here too for parity.
	for _, h := range []string{"127.0.0.2", "127.1.2.3", "0.0.0.0", "example.com", ""} {
		require.False(t, isLoopbackHost(h), "%q must not be treated as loopback", h)
	}
}

func TestNewRemoteSigner_RejectsEmptyAddress(t *testing.T) {
	wallet, err := NewWalletFromMnemonic(testMnemonic, DefaultHDPath)
	require.NoError(t, err)
	srv := serveBackend(t, walletInfoJSON(map[string]string{
		"id":      negWalletID,
		"address": "",
		"pubkey":  hex.EncodeToString(wallet.GetPublicKeyBytes()),
	}), nil)
	_, err = connectRemoteSigner(t, srv.URL)
	require.Error(t, err)
}

func TestNewRemoteSigner_RejectsWrongWalletID(t *testing.T) {
	wallet, err := NewWalletFromMnemonic(testMnemonic, DefaultHDPath)
	require.NoError(t, err)
	srv := serveBackend(t, walletInfoJSON(map[string]string{
		"id":      "33333333-3333-3333-3333-333333333333",
		"address": wallet.GetAddress(),
		"pubkey":  hex.EncodeToString(wallet.GetPublicKeyBytes()),
	}), nil)
	_, err = connectRemoteSigner(t, srv.URL)
	require.Error(t, err)
}

func TestNewRemoteSigner_RejectsMalformedPubKey(t *testing.T) {
	srv := serveBackend(t, walletInfoJSON(map[string]string{
		"id":      negWalletID,
		"address": "allo1placeholder",
		"pubkey":  "abcd", // valid hex, but only 2 bytes (not 33) — must not panic
	}), nil)
	_, err := connectRemoteSigner(t, srv.URL)
	require.Error(t, err)
}

func TestNewRemoteSigner_RejectsNonJSONResponse(t *testing.T) {
	srv := serveBackend(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<html>captive portal</html>"))
	}, nil)
	_, err := connectRemoteSigner(t, srv.URL)
	require.Error(t, err)
}

// TestRemoteSigner_RejectsSameOriginRedirect pins reject-all parity with the sibling SDKs
// (allora-sdk-py allow_redirects=False; allora-sdk-ts redirect: "error"): even a same-origin
// 307 redirect must NOT be followed, because Go re-sends the method, body, and headers on
// 307/308 and would re-POST the SignDoc bytes and the X-Forge-API-Key. The redirect surfaces
// as a non-2xx error and the redirect target handler is never reached.
func TestRemoteSigner_RejectsSameOriginRedirect(t *testing.T) {
	wallet, err := NewWalletFromMnemonic(testMnemonic, DefaultHDPath)
	require.NoError(t, err)

	var followed bool
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/signing-wallets/"+negWalletID, func(w http.ResponseWriter, _ *http.Request) {
		// Same-origin 307; the SDK must refuse to follow it.
		w.Header().Set("Location", "/api/v1/signing-wallets/"+negWalletID+"/moved")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	mux.HandleFunc("/api/v1/signing-wallets/"+negWalletID+"/moved", func(w http.ResponseWriter, _ *http.Request) {
		followed = true
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"id": negWalletID, "address": wallet.GetAddress(),
			"pubkey": hex.EncodeToString(wallet.GetPublicKeyBytes()),
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	_, err = connectRemoteSigner(t, srv.URL)
	require.Error(t, err)
	require.False(t, followed, "same-origin redirect must not be followed")
}

func TestRemoteSigner_Sign_RejectsRotatedKey(t *testing.T) {
	wallet, err := NewWalletFromMnemonic(testMnemonic, DefaultHDPath)
	require.NoError(t, err)
	other, err := GenerateWallet()
	require.NoError(t, err)

	srv := serveBackend(t,
		walletInfoJSON(map[string]string{
			"id":      negWalletID,
			"address": wallet.GetAddress(),
			"pubkey":  hex.EncodeToString(wallet.GetPublicKeyBytes()),
		}),
		func(w http.ResponseWriter, r *http.Request) {
			var body signRequest
			raw, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			require.NoError(t, json.Unmarshal(raw, &body))
			payload, err := hex.DecodeString(body.Payload)
			require.NoError(t, err)
			sig, err := wallet.PrivKey.Sign(payload)
			require.NoError(t, err)
			// Sign with the real key but report a different pubkey (simulated rotation).
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{
				"signature": hex.EncodeToString(sig),
				"pubkey":    hex.EncodeToString(other.GetPublicKeyBytes()),
			})
		},
	)
	rs, err := connectRemoteSigner(t, srv.URL)
	require.NoError(t, err)
	_, err = rs.Sign([]byte("allora sign doc bytes"))
	require.Error(t, err)
}

func TestRemoteSigner_Sign_RejectsMissingPubKeyEcho(t *testing.T) {
	// Fail closed: a sign response that omits the pubkey echo must be rejected, matching
	// allora-sdk-ts. Otherwise a backend could simply drop the field to dodge the
	// rotation/mis-route check.
	wallet, err := NewWalletFromMnemonic(testMnemonic, DefaultHDPath)
	require.NoError(t, err)

	srv := serveBackend(t,
		walletInfoJSON(map[string]string{
			"id":      negWalletID,
			"address": wallet.GetAddress(),
			"pubkey":  hex.EncodeToString(wallet.GetPublicKeyBytes()),
		}),
		func(w http.ResponseWriter, r *http.Request) {
			var body signRequest
			raw, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			require.NoError(t, json.Unmarshal(raw, &body))
			payload, err := hex.DecodeString(body.Payload)
			require.NoError(t, err)
			sig, err := wallet.PrivKey.Sign(payload)
			require.NoError(t, err)
			// Valid signature over the payload, but no pubkey echo field at all.
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{
				"signature": hex.EncodeToString(sig),
			})
		},
	)
	rs, err := connectRemoteSigner(t, srv.URL)
	require.NoError(t, err)
	_, err = rs.Sign([]byte("allora sign doc bytes"))
	require.ErrorContains(t, err, "omitted the pubkey echo")
}

func TestRemoteSigner_Sign_RejectsUnverifiableSignature(t *testing.T) {
	wallet, err := NewWalletFromMnemonic(testMnemonic, DefaultHDPath)
	require.NoError(t, err)

	srv := serveBackend(t,
		walletInfoJSON(map[string]string{
			"id":      negWalletID,
			"address": wallet.GetAddress(),
			"pubkey":  hex.EncodeToString(wallet.GetPublicKeyBytes()),
		}),
		func(w http.ResponseWriter, _ *http.Request) {
			// Valid length (64 bytes) but not a real signature over the payload.
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{
				"signature": hex.EncodeToString(make([]byte, 64)),
				"pubkey":    hex.EncodeToString(wallet.GetPublicKeyBytes()),
			})
		},
	)
	rs, err := connectRemoteSigner(t, srv.URL)
	require.NoError(t, err)
	_, err = rs.Sign([]byte("allora sign doc bytes"))
	require.Error(t, err)
}

func TestRemoteSigner_ClearAssociation(t *testing.T) {
	wallet, err := NewWalletFromMnemonic(testMnemonic, DefaultHDPath)
	require.NoError(t, err)
	const walletID = "11111111-1111-1111-1111-111111111111"

	var cleared bool
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/signing-wallets/"+walletID, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"id": walletID, "address": wallet.GetAddress(),
			"pubkey": hex.EncodeToString(wallet.GetPublicKeyBytes()),
		})
	})
	mux.HandleFunc("/api/v1/signing-wallets/"+walletID+"/clear-association", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.NotEmpty(t, r.Header.Get(apiKeyHeader))
		cleared = true
		w.WriteHeader(http.StatusNoContent) // 204, empty body — must be accepted as success
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	rs, err := NewRemoteSigner(context.Background(), RemoteSignerConfig{
		BackendURL: srv.URL, APIKey: "forge_sk_test", WalletID: walletID,
	})
	require.NoError(t, err)

	require.NoError(t, rs.ClearAssociation(context.Background()))
	require.True(t, cleared, "backend clear-association endpoint must be called")
}

func TestRemoteSigner_ClearAssociation_NonOKIsError(t *testing.T) {
	wallet, err := NewWalletFromMnemonic(testMnemonic, DefaultHDPath)
	require.NoError(t, err)
	const walletID = "11111111-1111-1111-1111-111111111111"

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/signing-wallets/"+walletID, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"id": walletID, "address": wallet.GetAddress(),
			"pubkey": hex.EncodeToString(wallet.GetPublicKeyBytes()),
		})
	})
	mux.HandleFunc("/api/v1/signing-wallets/"+walletID+"/clear-association", func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	rs, err := NewRemoteSigner(context.Background(), RemoteSignerConfig{
		BackendURL: srv.URL, APIKey: "forge_sk_test", WalletID: walletID,
	})
	require.NoError(t, err)

	err = rs.ClearAssociation(context.Background())
	require.Error(t, err)
	require.Contains(t, err.Error(), "404")
}

// TestRemoteSigner_Revoke pins the signer-level decommission method (synth-002): a caller holding
// a *RemoteSigner can revoke its own wallet via DELETE /api/v1/signing-wallets/{id} without
// rebuilding config, mirroring allora-sdk-ts signer.revoke() and Go's own ClearAssociation symmetry.
func TestRemoteSigner_Revoke(t *testing.T) {
	wallet, err := NewWalletFromMnemonic(testMnemonic, DefaultHDPath)
	require.NoError(t, err)
	const walletID = "11111111-1111-1111-1111-111111111111"

	var revoked bool
	mux := http.NewServeMux()
	// The wallet-info GET (construction) and the DELETE (revoke) share the same path, so
	// multiplex on the method.
	mux.HandleFunc("/api/v1/signing-wallets/"+walletID, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{
				"id": walletID, "address": wallet.GetAddress(),
				"pubkey": hex.EncodeToString(wallet.GetPublicKeyBytes()),
			})
		case http.MethodDelete:
			require.NotEmpty(t, r.Header.Get(apiKeyHeader))
			revoked = true
			w.WriteHeader(http.StatusNoContent) // 204, empty body — must be accepted as success
		default:
			t.Fatalf("unexpected method %q", r.Method)
		}
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	rs, err := NewRemoteSigner(context.Background(), RemoteSignerConfig{
		BackendURL: srv.URL, APIKey: "forge_sk_test", WalletID: walletID,
	})
	require.NoError(t, err)

	require.NoError(t, rs.Revoke(context.Background()))
	require.True(t, revoked, "backend revoke endpoint must be called")
}

// TestClearWalletAssociation_Standalone pins the by-id unbind path: it releases the binding
// using only the wallet id, without constructing a RemoteSigner or issuing a wallet-info GET
// (no GET handler is registered, so any fetch would fail the test). Mirrors the Python sibling's
// client-level ForgeBackendClient.clear_association(wallet_id).
func TestClearWalletAssociation_Standalone(t *testing.T) {
	const walletID = "11111111-1111-1111-1111-111111111111"
	var cleared bool
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/signing-wallets/"+walletID+"/clear-association", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.NotEmpty(t, r.Header.Get(apiKeyHeader))
		cleared = true
		w.WriteHeader(http.StatusNoContent) // 204, empty body — must be accepted as success
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	err := ClearWalletAssociation(context.Background(), RemoteSignerConfig{
		BackendURL: srv.URL, APIKey: "forge_sk_test",
	}, walletID)
	require.NoError(t, err)
	require.True(t, cleared, "backend clear-association endpoint must be called")
}

// TestClearWalletAssociation_RejectsNonUUID pins that a non-UUID wallet id is rejected before
// any request is issued, so it cannot inject path segments into the request URL.
func TestClearWalletAssociation_RejectsNonUUID(t *testing.T) {
	err := ClearWalletAssociation(context.Background(), RemoteSignerConfig{
		BackendURL: "https://forge.allora.network", APIKey: "forge_sk_test",
	}, "not-a-uuid")
	require.Error(t, err)
}

// TestRemoteSigner_RedactsAPIKeyInError pins that the Forge API key is stripped from a backend
// error excerpt (synth-009): if a TLS-terminating intermediary echoes the X-Forge-API-Key request
// header into its error body, the key must not leak into the returned error string (and thus logs).
func TestRemoteSigner_RedactsAPIKeyInError(t *testing.T) {
	const apiKey = "forge_sk_supersecret"
	const walletID = "11111111-1111-1111-1111-111111111111"
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/signing-wallets/"+walletID, func(w http.ResponseWriter, r *http.Request) {
		// Simulate an intermediary that mirrors the request's API-key header into its error body.
		http.Error(w, "gateway error; saw header X-Forge-API-Key: "+r.Header.Get(apiKeyHeader), http.StatusBadGateway)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	_, err := NewRemoteSigner(context.Background(), RemoteSignerConfig{
		BackendURL: srv.URL, APIKey: apiKey, WalletID: walletID,
	})
	require.Error(t, err)
	require.NotContains(t, err.Error(), apiKey, "API key must not appear in the error string")
	require.Contains(t, err.Error(), "[REDACTED]")
}

// masterGranterBackend serves a wallet-info GET for walletID that reports the given
// master_granter, so the discovery + resolution tests can vary the advertised granter.
func masterGranterBackend(t *testing.T, wallet *Wallet, walletID, masterGranter string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/signing-wallets/"+walletID, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"id": walletID, "address": wallet.GetAddress(),
			"pubkey":         hex.EncodeToString(wallet.GetPublicKeyBytes()),
			"master_granter": masterGranter,
		})
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

// TestRemoteSigner_DiscoversMasterGranter pins runtime discovery (synth-003): the backend's
// master_granter field is captured at construction and surfaced via MasterGranter(), so a worker
// learns the gas granter from the API instead of configuring it out-of-band.
func TestRemoteSigner_DiscoversMasterGranter(t *testing.T) {
	wallet, err := NewWalletFromMnemonic(testMnemonic, DefaultHDPath)
	require.NoError(t, err)
	granter, err := GenerateWallet() // a real allo1... bech32 address to advertise as the granter
	require.NoError(t, err)

	const walletID = "11111111-1111-1111-1111-111111111111"
	srv := masterGranterBackend(t, wallet, walletID, granter.GetAddress())

	rs, err := NewRemoteSigner(context.Background(), RemoteSignerConfig{
		BackendURL: srv.URL, APIKey: "forge_sk_test", WalletID: walletID,
	})
	require.NoError(t, err)
	require.Equal(t, granter.GetAddress(), rs.MasterGranter())
}

// TestRemoteSigner_MasterGranterAbsent pins graceful degradation: a backend that omits
// master_granter leaves MasterGranter() empty and ResolveFeeGranter() nil (the wallet pays its
// own fees) rather than failing signer construction.
func TestRemoteSigner_MasterGranterAbsent(t *testing.T) {
	wallet, err := NewWalletFromMnemonic(testMnemonic, DefaultHDPath)
	require.NoError(t, err)

	const walletID = "11111111-1111-1111-1111-111111111111"
	srv := masterGranterBackend(t, wallet, walletID, "")

	rs, err := NewRemoteSigner(context.Background(), RemoteSignerConfig{
		BackendURL: srv.URL, APIKey: "forge_sk_test", WalletID: walletID,
	})
	require.NoError(t, err)
	require.Empty(t, rs.MasterGranter())

	t.Setenv("FORGE_MASTER_GRANTER_ADDRESS", "")
	got, err := rs.ResolveFeeGranter()
	require.NoError(t, err)
	require.Nil(t, got)
}

// TestRemoteSigner_MasterGranterVerbatim pins that a backend master_granter is stored verbatim,
// without trimming (synth-001): MasterGranter() returns exactly what the backend sent — as the
// README advertises ("raw string, verbatim") — keeping Go at parity with allora-sdk-py /
// allora-sdk-ts, which assign the field unmodified. A padded value therefore surfaces as a bech32
// parse error in ResolveFeeGranter, identically across the three SDKs, rather than being silently
// rescued by a Go-only trim.
func TestRemoteSigner_MasterGranterVerbatim(t *testing.T) {
	wallet, err := NewWalletFromMnemonic(testMnemonic, DefaultHDPath)
	require.NoError(t, err)
	granter, err := GenerateWallet()
	require.NoError(t, err)

	const walletID = "11111111-1111-1111-1111-111111111111"
	padded := "  " + granter.GetAddress() + "\n"
	srv := masterGranterBackend(t, wallet, walletID, padded)

	rs, err := NewRemoteSigner(context.Background(), RemoteSignerConfig{
		BackendURL: srv.URL, APIKey: "forge_sk_test", WalletID: walletID,
	})
	require.NoError(t, err)
	require.Equal(t, padded, rs.MasterGranter(), "master_granter must be returned verbatim, untrimmed")

	// A padded value is not silently rescued: it fails bech32 parsing in ResolveFeeGranter, the
	// same way it would on the Python/TypeScript siblings.
	t.Setenv("FORGE_MASTER_GRANTER_ADDRESS", "")
	_, err = rs.ResolveFeeGranter()
	require.Error(t, err)
}

// TestRevokeWallet_Standalone pins the by-id decommission path (synth-015): it deletes the wallet
// using only the wallet id, without constructing a RemoteSigner or issuing a wallet-info GET (no
// GET handler is registered, so any fetch would fail the test). Mirrors the server's
// RevokeSigningWallet (DELETE /api/v1/signing-wallets/{id}), which answers 200 + JSON.
func TestRevokeWallet_Standalone(t *testing.T) {
	const walletID = "11111111-1111-1111-1111-111111111111"
	var revoked bool
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/signing-wallets/"+walletID, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.NotEmpty(t, r.Header.Get(apiKeyHeader))
		revoked = true
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "wallet revoked"})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	err := RevokeWallet(context.Background(), RemoteSignerConfig{
		BackendURL: srv.URL, APIKey: "forge_sk_test",
	}, walletID)
	require.NoError(t, err)
	require.True(t, revoked, "backend revoke endpoint must be called")
}

// TestRevokeWallet_RejectsNonUUID pins that a non-UUID wallet id is rejected before any request
// is issued, so it cannot inject path segments into the request URL.
func TestRevokeWallet_RejectsNonUUID(t *testing.T) {
	err := RevokeWallet(context.Background(), RemoteSignerConfig{
		BackendURL: "https://forge.allora.network", APIKey: "forge_sk_test",
	}, "not-a-uuid")
	require.Error(t, err)
}

// TestRevokeWallet_NonOKIsError pins that a non-2xx response surfaces as an error so the caller
// can decide whether a failed decommission is fatal or best-effort.
func TestRevokeWallet_NonOKIsError(t *testing.T) {
	const walletID = "11111111-1111-1111-1111-111111111111"
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/signing-wallets/"+walletID, func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	err := RevokeWallet(context.Background(), RemoteSignerConfig{
		BackendURL: srv.URL, APIKey: "forge_sk_test",
	}, walletID)
	require.Error(t, err)
	require.Contains(t, err.Error(), "404")
}

// TestRemoteSigner_ResolveFeeGranter pins the discovery-with-override precedence: with no env
// override the discovered master granter is used, and FORGE_MASTER_GRANTER_ADDRESS overrides it.
func TestRemoteSigner_ResolveFeeGranter(t *testing.T) {
	wallet, err := NewWalletFromMnemonic(testMnemonic, DefaultHDPath)
	require.NoError(t, err)
	discovered, err := GenerateWallet()
	require.NoError(t, err)
	override, err := GenerateWallet()
	require.NoError(t, err)

	const walletID = "11111111-1111-1111-1111-111111111111"
	srv := masterGranterBackend(t, wallet, walletID, discovered.GetAddress())

	rs, err := NewRemoteSigner(context.Background(), RemoteSignerConfig{
		BackendURL: srv.URL, APIKey: "forge_sk_test", WalletID: walletID,
	})
	require.NoError(t, err)

	// No env override: fall back to the discovered master granter.
	t.Setenv("FORGE_MASTER_GRANTER_ADDRESS", "")
	got, err := rs.ResolveFeeGranter()
	require.NoError(t, err)
	require.Equal(t, discovered.GetAddress(), got.String())

	// Env override wins over the discovered value.
	t.Setenv("FORGE_MASTER_GRANTER_ADDRESS", override.GetAddress())
	got, err = rs.ResolveFeeGranter()
	require.NoError(t, err)
	require.Equal(t, override.GetAddress(), got.String())
}
