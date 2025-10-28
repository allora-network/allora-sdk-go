package allora

import (
	"crypto/rand"
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
)

const (
	// AlloraBech32Prefix is the human-readable prefix for Allora addresses
	AlloraBech32Prefix = "allo"

	// DefaultBIP44CoinType is the coin type for Allora (118 is standard for Cosmos chains)
	DefaultBIP44CoinType = 118

	// DefaultHDPath is the default HD derivation path for Allora wallets
	DefaultHDPath = "m/44'/118'/0'/0/0"
)

// Wallet represents an Allora wallet with a private key and address
type Wallet struct {
	PrivKey cryptotypes.PrivKey
	PubKey  cryptotypes.PubKey
	Address sdk.AccAddress
	Mnemonic string
}

// NewWalletFromMnemonic creates a wallet from a BIP39 mnemonic phrase
func NewWalletFromMnemonic(mnemonic string, hdPath string) (*Wallet, error) {
	if hdPath == "" {
		hdPath = DefaultHDPath
	}

	// Derive the private key from mnemonic
	algo := hd.Secp256k1
	derivedPriv, err := algo.Derive()(mnemonic, "", hdPath)
	if err != nil {
		return nil, fmt.Errorf("failed to derive private key: %w", err)
	}

	privKey := algo.Generate()(derivedPriv)
	return newWalletFromPrivKey(privKey, mnemonic)
}

// GenerateWallet creates a new wallet with a random mnemonic
func GenerateWallet() (*Wallet, error) {
	return GenerateWalletWithMnemonicLength(256)
}

// GenerateWalletWithMnemonicLength creates a new wallet with a mnemonic of specified entropy bits (128, 160, 192, 224, or 256)
func GenerateWalletWithMnemonicLength(entropyBits int) (*Wallet, error) {
	// Generate entropy
	entropy, err := bip39.NewEntropy(entropyBits)
	if err != nil {
		return nil, fmt.Errorf("failed to generate entropy: %w", err)
	}

	// Generate mnemonic from entropy
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	// Create wallet from mnemonic
	return NewWalletFromMnemonic(mnemonic, DefaultHDPath)
}

// NewWalletFromPrivateKey creates a wallet from a raw secp256k1 private key (32 bytes)
func NewWalletFromPrivateKey(privKeyBytes []byte) (*Wallet, error) {
	if len(privKeyBytes) != 32 {
		return nil, fmt.Errorf("private key must be 32 bytes, got %d", len(privKeyBytes))
	}

	privKey := &secp256k1.PrivKey{Key: privKeyBytes}
	return newWalletFromPrivKey(privKey, "")
}

// GenerateRandomPrivateKey generates a new random secp256k1 private key
func GenerateRandomPrivateKey() ([]byte, error) {
	privKeyBytes := make([]byte, 32)
	_, err := rand.Read(privKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return privKeyBytes, nil
}

// newWalletFromPrivKey is an internal helper to create a wallet from a private key
func newWalletFromPrivKey(privKey cryptotypes.PrivKey, mnemonic string) (*Wallet, error) {
	pubKey := privKey.PubKey()

	// Convert public key to Allora address with bech32 prefix
	addr := sdk.AccAddress(pubKey.Address())

	return &Wallet{
		PrivKey:  privKey,
		PubKey:   pubKey,
		Address:  addr,
		Mnemonic: mnemonic,
	}, nil
}

// GetAddress returns the bech32-encoded Allora address
func (w *Wallet) GetAddress() string {
	return w.Address.String()
}

// GetPrivateKeyBytes returns the raw private key bytes
func (w *Wallet) GetPrivateKeyBytes() []byte {
	return w.PrivKey.Bytes()
}

// GetPublicKeyBytes returns the raw public key bytes
func (w *Wallet) GetPublicKeyBytes() []byte {
	return w.PubKey.Bytes()
}

// GetMnemonic returns the mnemonic phrase (if wallet was created from one)
func (w *Wallet) GetMnemonic() string {
	return w.Mnemonic
}

// Sign signs a message with the wallet's private key
func (w *Wallet) Sign(message []byte) ([]byte, error) {
	signature, err := w.PrivKey.Sign(message)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %w", err)
	}
	return signature, nil
}

// VerifySignature verifies a signature against a message using the wallet's public key
func (w *Wallet) VerifySignature(message, signature []byte) bool {
	return w.PubKey.VerifySignature(message, signature)
}

// init configures the SDK to use the Allora bech32 prefix
func init() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AlloraBech32Prefix, AlloraBech32Prefix+"pub")
	config.SetBech32PrefixForValidator(AlloraBech32Prefix+"valoper", AlloraBech32Prefix+"valoperpub")
	config.SetBech32PrefixForConsensusNode(AlloraBech32Prefix+"valcons", AlloraBech32Prefix+"valconspub")
	config.Seal()
}
