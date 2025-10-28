package allora

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/cosmos/go-bip39"
)

func TestGenerateWallet(t *testing.T) {
	wallet, err := GenerateWallet()
	if err != nil {
		t.Fatalf("failed to generate wallet: %v", err)
	}

	// Check that address has correct prefix
	addr := wallet.GetAddress()
	if !strings.HasPrefix(addr, AlloraBech32Prefix) {
		t.Errorf("address does not have correct prefix, got: %s", addr)
	}

	// Check that private key is 32 bytes
	privKey := wallet.GetPrivateKeyBytes()
	if len(privKey) != 32 {
		t.Errorf("private key should be 32 bytes, got %d", len(privKey))
	}

	// Check that public key is not empty
	pubKey := wallet.GetPublicKeyBytes()
	if len(pubKey) == 0 {
		t.Error("public key should not be empty")
	}

	// Check that mnemonic is valid
	if wallet.GetMnemonic() == "" {
		t.Error("mnemonic should not be empty")
	}

	if !bip39.IsMnemonicValid(wallet.GetMnemonic()) {
		t.Error("mnemonic should be valid")
	}

	t.Logf("Generated wallet address: %s", addr)
	t.Logf("Mnemonic: %s", wallet.GetMnemonic())
}

func TestGenerateWalletWithDifferentLengths(t *testing.T) {
	tests := []struct {
		name        string
		entropyBits int
		wordCount   int
	}{
		{"12 words", 128, 12},
		{"15 words", 160, 15},
		{"18 words", 192, 18},
		{"21 words", 224, 21},
		{"24 words", 256, 24},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wallet, err := GenerateWalletWithMnemonicLength(tt.entropyBits)
			if err != nil {
				t.Fatalf("failed to generate wallet: %v", err)
			}

			words := strings.Split(wallet.GetMnemonic(), " ")
			if len(words) != tt.wordCount {
				t.Errorf("expected %d words, got %d", tt.wordCount, len(words))
			}

			addr := wallet.GetAddress()
			if !strings.HasPrefix(addr, AlloraBech32Prefix) {
				t.Errorf("address does not have correct prefix, got: %s", addr)
			}
		})
	}
}

func TestNewWalletFromMnemonic(t *testing.T) {
	// Use a known mnemonic for deterministic testing
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

	wallet, err := NewWalletFromMnemonic(mnemonic, DefaultHDPath)
	if err != nil {
		t.Fatalf("failed to create wallet from mnemonic: %v", err)
	}

	// The address should be deterministic
	addr := wallet.GetAddress()
	if !strings.HasPrefix(addr, AlloraBech32Prefix) {
		t.Errorf("address does not have correct prefix, got: %s", addr)
	}

	// Create another wallet with same mnemonic - should produce same address
	wallet2, err := NewWalletFromMnemonic(mnemonic, DefaultHDPath)
	if err != nil {
		t.Fatalf("failed to create second wallet from mnemonic: %v", err)
	}

	if wallet.GetAddress() != wallet2.GetAddress() {
		t.Errorf("same mnemonic should produce same address, got %s and %s",
			wallet.GetAddress(), wallet2.GetAddress())
	}

	t.Logf("Deterministic address from test mnemonic: %s", addr)
}

func TestNewWalletFromPrivateKey(t *testing.T) {
	// Generate random private key
	privKeyBytes, err := GenerateRandomPrivateKey()
	if err != nil {
		t.Fatalf("failed to generate random private key: %v", err)
	}

	wallet, err := NewWalletFromPrivateKey(privKeyBytes)
	if err != nil {
		t.Fatalf("failed to create wallet from private key: %v", err)
	}

	// Check address has correct prefix
	addr := wallet.GetAddress()
	if !strings.HasPrefix(addr, AlloraBech32Prefix) {
		t.Errorf("address does not have correct prefix, got: %s", addr)
	}

	// Check that we can recreate the same wallet from the same private key
	wallet2, err := NewWalletFromPrivateKey(privKeyBytes)
	if err != nil {
		t.Fatalf("failed to create second wallet from private key: %v", err)
	}

	if wallet.GetAddress() != wallet2.GetAddress() {
		t.Errorf("same private key should produce same address")
	}

	t.Logf("Address from random private key: %s", addr)
}

func TestSignAndVerify(t *testing.T) {
	wallet, err := GenerateWallet()
	if err != nil {
		t.Fatalf("failed to generate wallet: %v", err)
	}

	message := []byte("hello allora network")

	// Sign the message
	signature, err := wallet.Sign(message)
	if err != nil {
		t.Fatalf("failed to sign message: %v", err)
	}

	if len(signature) == 0 {
		t.Error("signature should not be empty")
	}

	// Verify the signature
	valid := wallet.VerifySignature(message, signature)
	if !valid {
		t.Error("signature verification failed")
	}

	// Try with wrong message - should fail
	wrongMessage := []byte("wrong message")
	valid = wallet.VerifySignature(wrongMessage, signature)
	if valid {
		t.Error("signature verification should fail with wrong message")
	}

	t.Logf("Signature: %s", hex.EncodeToString(signature))
}

func TestInvalidPrivateKeyLength(t *testing.T) {
	// Try with wrong length private key
	invalidKey := make([]byte, 16) // Wrong length
	_, err := NewWalletFromPrivateKey(invalidKey)
	if err == nil {
		t.Error("should fail with invalid private key length")
	}
}

func TestInvalidMnemonic(t *testing.T) {
	invalidMnemonic := "invalid mnemonic phrase that is not real"
	_, err := NewWalletFromMnemonic(invalidMnemonic, DefaultHDPath)
	if err == nil {
		t.Error("should fail with invalid mnemonic")
	}
}
