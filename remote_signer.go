package allora

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

var _ Signer = (*RemoteSigner)(nil)

// NewRemoteSigner builds a RemoteSigner and fetches the wallet's public key and address
// from the backend (needed to assemble the SignDoc before the wallet has transacted).
func NewRemoteSigner(ctx context.Context, cfg RemoteSignerConfig) (*RemoteSigner, error) {
	if cfg.BackendURL == "" || cfg.APIKey == "" || cfg.WalletID == "" {
		return nil, fmt.Errorf("backend URL, API key, and wallet ID are all required")
	}
	rs := &RemoteSigner{cfg: cfg, httpClient: cfg.HTTPClient}
	if rs.httpClient == nil {
		rs.httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	if err := rs.fetchWallet(ctx); err != nil {
		return nil, err
	}
	return rs, nil
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
	url := fmt.Sprintf("%s/api/v1/signing-wallets/%s", rs.cfg.BackendURL, rs.cfg.WalletID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("creating wallet-info request: %w", err)
	}
	req.Header.Set(apiKeyHeader, rs.cfg.APIKey)

	body, err := rs.do(req)
	if err != nil {
		return err
	}

	var info signingWalletInfoResponse
	if err := json.Unmarshal(body, &info); err != nil {
		return fmt.Errorf("decoding wallet-info response: %w", err)
	}
	pubBytes, err := hex.DecodeString(info.PubKey)
	if err != nil {
		return fmt.Errorf("decoding wallet pubkey: %w", err)
	}
	rs.pubKey = &secp256k1.PubKey{Key: pubBytes}

	// Derive the address from the pubkey and cross-check it against the backend's
	// reported address so a misconfigured wallet fails here rather than on broadcast.
	derived := sdk.AccAddress(rs.pubKey.Address())
	if info.Address != "" && derived.String() != info.Address {
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

// Sign delegates signing of the SignDoc bytes to the backend. The backend SHA-256
// hashes the payload and signs it with the Privy wallet, returning the 64-byte
// signature. Signer.Sign carries no context, so a background context bounded by the
// HTTP client timeout is used.
func (rs *RemoteSigner) Sign(msg []byte) ([]byte, error) {
	reqBody, err := json.Marshal(signRequest{Payload: hex.EncodeToString(msg), Prehashed: false})
	if err != nil {
		return nil, fmt.Errorf("marshaling sign request: %w", err)
	}
	url := fmt.Sprintf("%s/api/v1/signing-wallets/%s/sign", rs.cfg.BackendURL, rs.cfg.WalletID)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("creating sign request: %w", err)
	}
	req.Header.Set(apiKeyHeader, rs.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	body, err := rs.do(req)
	if err != nil {
		return nil, err
	}
	var sr signResponse
	if err := json.Unmarshal(body, &sr); err != nil {
		return nil, fmt.Errorf("decoding sign response: %w", err)
	}
	sig, err := hex.DecodeString(sr.Signature)
	if err != nil {
		return nil, fmt.Errorf("decoding signature: %w", err)
	}
	return sig, nil
}

// do executes an HTTP request and returns the body, mapping non-2xx to an error.
func (rs *RemoteSigner) do(req *http.Request) ([]byte, error) {
	resp, err := rs.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling forge backend: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading forge response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("forge backend returned status %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}
