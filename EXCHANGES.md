# Allora Network Integration Guide for Exchanges

This guide provides instructions for cryptocurrency exchanges integrating with the Allora Network blockchain.

## Installation

```bash
go get github.com/allora-network/allora-sdk-go
```

## Wallet Management

The Allora SDK provides comprehensive wallet utilities for managing keypairs with the `allo` bech32 address prefix.

### Generating a New Wallet

```go
import allora "github.com/allora-network/allora-sdk-go"

// Generate a new wallet with a 24-word mnemonic
wallet, err := allora.GenerateWallet()
if err != nil {
    panic(err)
}

// Access wallet properties
address := wallet.GetAddress()           // e.g., "allo1..."
mnemonic := wallet.GetMnemonic()         // 24-word recovery phrase
privKey := wallet.GetPrivateKeyBytes()   // 32-byte private key
pubKey := wallet.GetPublicKeyBytes()     // Public key bytes
```

### Generating Wallets with Different Mnemonic Lengths

```go
// 12-word mnemonic (128 bits entropy)
wallet12, _ := allora.GenerateWalletWithMnemonicLength(128)

// 24-word mnemonic (256 bits entropy) - recommended for maximum security
wallet24, _ := allora.GenerateWalletWithMnemonicLength(256)
```

### Recovering a Wallet from Mnemonic

```go
mnemonic := "your twenty four word mnemonic phrase here..."
hdPath := allora.DefaultHDPath // "m/44'/118'/0'/0/0"

wallet, err := allora.NewWalletFromMnemonic(mnemonic, hdPath)
if err != nil {
    panic(err)
}

address := wallet.GetAddress()
```

### Importing from Private Key

```go
// privKeyBytes is a 32-byte secp256k1 private key
wallet, err := allora.NewWalletFromPrivateKey(privKeyBytes)
if err != nil {
    panic(err)
}
```

### Generating a Random Private Key

```go
privKeyBytes, err := allora.GenerateRandomPrivateKey()
if err != nil {
    panic(err)
}

// Create wallet from the generated key
wallet, err := allora.NewWalletFromPrivateKey(privKeyBytes)
```

## Address Format

All Allora Network addresses use the `allo` bech32 prefix:
- **Account addresses**: `allo1...` (42 characters)
- **Validator operator addresses**: `allovaloper1...`
- **Consensus addresses**: `allovalcons1...`

Exchanges should only handle account addresses (`allo1...`).

## Signing and Verification

### Signing a Message

```go
message := []byte("transaction data")

signature, err := wallet.Sign(message)
if err != nil {
    panic(err)
}
```

### Verifying a Signature

```go
message := []byte("transaction data")
signature := []byte{...} // signature bytes

isValid := wallet.VerifySignature(message, signature)
```

## Security Best Practices

### Key Storage
- **Never store mnemonics or private keys in plain text**
- Use hardware security modules (HSMs) for production hot wallets
- Implement key encryption at rest using industry-standard methods
- Consider multi-signature schemes for large holdings

### Cold Storage
- Generate wallets offline for cold storage
- Store mnemonic phrases in secure, geographically distributed locations
- Use 24-word mnemonics (256-bit entropy) for enhanced security

### Address Validation
```go
import sdk "github.com/cosmos/cosmos-sdk/types"

func ValidateAlloraAddress(address string) error {
    // Check prefix
    if !strings.HasPrefix(address, "allo") {
        return fmt.Errorf("invalid address prefix")
    }

    // Validate bech32 format
    _, err := sdk.AccAddressFromBech32(address)
    return err
}
```

## Constants

```go
const (
    AlloraBech32Prefix = "allo"      // Address prefix
    DefaultBIP44CoinType = 118       // Cosmos standard coin type
    DefaultHDPath = "m/44'/118'/0'/0/0"  // Default derivation path
)
```

## Wallet Compatibility

Allora uses standard Cosmos SDK cryptography:
- **Signature algorithm**: secp256k1 (same as Bitcoin and Ethereum)
- **Derivation path**: BIP44 with coin type 118
- **Address encoding**: bech32 with `allo` prefix

Wallets are compatible with other Cosmos-based chains when using the Allora-specific prefix.

## Transaction Broadcasting

**Note**: Transaction signing and broadcasting functionality is currently under development. This document will be updated when those features are available.

For read-only operations (querying balances, network state, etc.), see the main README.md.

## Support

For integration support:
- GitHub Issues: https://github.com/allora-network/allora-sdk-go/issues
- Documentation: https://docs.allora.network

## Example: Complete Wallet Setup

```go
package main

import (
    "fmt"
    allora "github.com/allora-network/allora-sdk-go"
)

func main() {
    // Generate new deposit wallet
    wallet, err := allora.GenerateWallet()
    if err != nil {
        panic(err)
    }

    // Display wallet information
    fmt.Printf("Deposit Address: %s\n", wallet.GetAddress())
    fmt.Printf("Store this mnemonic securely: %s\n", wallet.GetMnemonic())

    // Later, recover the wallet
    recovered, err := allora.NewWalletFromMnemonic(
        wallet.GetMnemonic(),
        allora.DefaultHDPath,
    )
    if err != nil {
        panic(err)
    }

    // Verify addresses match
    if wallet.GetAddress() != recovered.GetAddress() {
        panic("address mismatch")
    }

    fmt.Println("Wallet successfully recovered!")
}
```
