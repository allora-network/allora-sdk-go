package allora

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/stretchr/testify/require"
)

const testMnemonic = "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

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
