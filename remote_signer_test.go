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
		raw, _ := io.ReadAll(r.Body)
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

	// secp256k1 signing is deterministic (RFC 6979), so signing the same SignDoc with the
	// same key via the remote path must yield byte-identical signed tx bytes.
	require.Equal(t, localSigned, remoteSigned)
}

func TestTxParams_FeeGranterIsEncoded(t *testing.T) {
	wallet, err := NewWalletFromMnemonic(testMnemonic, DefaultHDPath)
	require.NoError(t, err)

	granter, err := sdk.AccAddressFromBech32(wallet.GetAddress())
	require.NoError(t, err)

	amount := sdk.NewCoins(sdk.NewInt64Coin("uallo", 1000))
	base := &TxParams{
		ChainID:   "allora-testnet-1",
		GasLimit:  200000,
		FeeAmount: sdk.NewCoins(sdk.NewInt64Coin("uallo", 5000)),
	}
	withGranter := *base
	withGranter.FeeGranter = granter

	unsignedNoGranter, err := CreateUnsignedSendTx(wallet.Address, wallet.Address, amount, base)
	require.NoError(t, err)
	unsignedWithGranter, err := CreateUnsignedSendTx(wallet.Address, wallet.Address, amount, &withGranter)
	require.NoError(t, err)

	// Setting the fee granter must change the encoded AuthInfo (the granter address is
	// written into the tx). An empty granter leaves the tx untouched.
	require.NotEqual(t, unsignedNoGranter, unsignedWithGranter)
	require.Greater(t, len(unsignedWithGranter), len(unsignedNoGranter))
}

func TestNewRemoteSigner_RequiresConfig(t *testing.T) {
	_, err := NewRemoteSigner(context.Background(), RemoteSignerConfig{})
	require.Error(t, err)
}
