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
