package allora_test

import (
	"fmt"

	allora "github.com/allora-network/allora-sdk-go"
)

// Example_generateWallet demonstrates how to generate a new random wallet
func Example_generateWallet() {
	// Generate a new wallet with a random mnemonic
	wallet, err := allora.GenerateWallet()
	if err != nil {
		panic(err)
	}

	fmt.Println("Wallet generated successfully!")
	fmt.Printf("Address: %s\n", wallet.GetAddress())
	fmt.Printf("Mnemonic: %s\n", wallet.GetMnemonic())
	fmt.Printf("Public Key (hex): %x\n", wallet.GetPublicKeyBytes())

	// IMPORTANT: Store the mnemonic securely!
	// The mnemonic can be used to recover the wallet later
}

// Example_recoverWalletFromMnemonic demonstrates how to recover a wallet from a mnemonic
func Example_recoverWalletFromMnemonic() {
	// Use an existing mnemonic to recover a wallet
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

	wallet, err := allora.NewWalletFromMnemonic(mnemonic, allora.DefaultHDPath)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Wallet recovered from mnemonic\n")
	fmt.Printf("Address: %s\n", wallet.GetAddress())
}

// Example_importPrivateKey demonstrates how to create a wallet from a private key
func Example_importPrivateKey() {
	// Generate a random private key (in practice, you'd import an existing one)
	privKeyBytes, err := allora.GenerateRandomPrivateKey()
	if err != nil {
		panic(err)
	}

	// Create wallet from private key
	wallet, err := allora.NewWalletFromPrivateKey(privKeyBytes)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Wallet created from private key\n")
	fmt.Printf("Address: %s\n", wallet.GetAddress())
}

// Example_signAndVerify demonstrates how to sign a message and verify the signature
func Example_signAndVerify() {
	// Generate a wallet
	wallet, err := allora.GenerateWallet()
	if err != nil {
		panic(err)
	}

	// Message to sign
	message := []byte("Hello Allora Network!")

	// Sign the message
	signature, err := wallet.Sign(message)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Message signed successfully\n")
	fmt.Printf("Signature length: %d bytes\n", len(signature))

	// Verify the signature
	isValid := wallet.VerifySignature(message, signature)
	fmt.Printf("Signature valid: %v\n", isValid)
}

// Example_differentMnemonicLengths demonstrates generating wallets with different mnemonic lengths
func Example_differentMnemonicLengths() {
	// Generate 12-word mnemonic (128 bits entropy)
	wallet12, _ := allora.GenerateWalletWithMnemonicLength(128)
	fmt.Printf("12-word mnemonic wallet: %s\n", wallet12.GetAddress())

	// Generate 24-word mnemonic (256 bits entropy) - most secure
	wallet24, _ := allora.GenerateWalletWithMnemonicLength(256)
	fmt.Printf("24-word mnemonic wallet: %s\n", wallet24.GetAddress())
}
