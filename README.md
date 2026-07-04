# Allora SDK Go

A high-performance Go client SDK for the Allora Network blockchain, providing comprehensive access to query services across Allora custom modules and Cosmos SDK base modules.

## Features

- **Multi-endpoint support**: Connect to multiple RPC endpoints with automatic load balancing
- **Fault tolerance**: Automatic failover and retry mechanisms with exponential backoff
- **Health tracking**: Built-in client health monitoring and automatic recovery
- **Protocol abstraction**: Currently supports gRPC, extensible for other protocols
- **Comprehensive coverage**: Query support for:
  - Allora Emissions module (v9) - 100+ query methods
  - Allora Mint module (v5) - 3 query methods
  - Cosmos SDK Auth module - 5 query methods
  - Cosmos SDK Bank module - 9 query methods  
  - Cosmos SDK Staking module - 14 query methods
  - Tendermint/CometBFT service - 7 service methods
- **Production ready**: Comprehensive error handling, logging, and connection management

## Installation

```bash
go get github.com/allora-network/allora-sdk-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/rs/zerolog"
    allora "github.com/allora-network/allora-sdk-go"
)

func main() {
    // Setup logger
    logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

    // Create client configuration with multiple endpoints for redundancy
    config := &allora.ClientConfig{
        Endpoints: []allora.EndpointConfig{
            {
                URL:      "grpc://localhost:9090",
                Protocol: "grpc",
            },
            {
                URL:      "grpc://backup.allora.com:9090",
                Protocol: "grpc",
            },
        },
        RequestTimeout:    30 * time.Second,
        ConnectionTimeout: 10 * time.Second,
    }

    // Create the client
    client, err := allora.NewClient(config, logger)
    if err != nil {
        logger.Fatal().Err(err).Msg("failed to create client")
    }
    defer client.Close()

    ctx := context.Background()

    // Get node information
    nodeInfo, err := client.GetNodeInfo(ctx)
    if err != nil {
        logger.Error().Err(err).Msg("failed to get node info")
        return
    }
    fmt.Printf("Connected to: %s\\n", nodeInfo.DefaultNodeInfo.Moniker)

    // Get latest block
    block, err := client.GetLatestBlock(ctx)
    if err != nil {
        logger.Error().Err(err).Msg("failed to get latest block")
        return
    }
    fmt.Printf("Latest block height: %d\\n", block.SdkBlock.Header.Height)

    // Get account balance
    balance, err := client.BankBalance(ctx, "allo1...", "uallo")
    if err != nil {
        logger.Error().Err(err).Msg("failed to get balance")
        return
    }
    fmt.Printf("Balance: %s\\n", balance.Balance.Amount)

    // Get validators
    validators, err := client.StakingValidators(ctx, "", &allora.PageRequest{Limit: 10})
    if err != nil {
        logger.Error().Err(err).Msg("failed to get validators")
        return
    }
    fmt.Printf("Found %d validators\\n", len(validators.Validators))
}
```

## Privy-Managed Signing (delegated)

By default the SDK signs with a local key (`Wallet`). Alternatively, you can delegate
signing to the Forge backend, which signs with a Privy-managed server wallet — the worker
holds no private key. A local wallet's key (`wallet.PrivKey`) and the backend-backed
`RemoteSigner` both satisfy the `Signer` interface, so `SignTransactionWith` accepts
either. (For a local wallet you can also use the `SignTransaction(unsignedTx, wallet,
params)` convenience wrapper.)

```go
ctx := context.Background()

// The wallet ID and API key are minted in the Forge web app. The signer fetches its
// address + public key from the backend on construction.
signer, err := allora.NewRemoteSigner(ctx, allora.RemoteSignerConfig{
    BackendURL: "https://forge.allora.network",
    APIKey:     os.Getenv("FORGE_API_KEY"),
    WalletID:   os.Getenv("FORGE_SIGNING_WALLET_ID"),
})
if err != nil {
    panic(err)
}

params := &allora.TxParams{
    ChainID:       "allora-testnet-1",
    AccountNumber: accountNumber,
    Sequence:      sequence,
    GasLimit:      200000,
    FeeAmount:     sdk.NewCoins(sdk.NewInt64Coin("uallo", 5000)),
    // FeeGranter: masterAddr, // optional: subsidize gas from a master wallet via feegrant
}

// signer.PubKey().Address() works for any Signer (local key or RemoteSigner); a
// *RemoteSigner additionally offers the signer.AccAddress() convenience helper.
fromAddr := sdk.AccAddress(signer.PubKey().Address())
unsignedTx, _ := allora.CreateUnsignedSendTx(fromAddr, toAddr, amount, params)
signedTx, _ := allora.SignTransactionWith(ctx, unsignedTx, signer, params)
// broadcast signedTx via client.Cosmos().Tx().BroadcastTx(...)
```

The self-managed path — `SignTransaction(unsignedTx, wallet, params)` with a local
`Wallet` — is unchanged.

### Fee granter (optional gas subsidy)

Set `TxParams.FeeGranter` to have the transaction fee paid by another account via an
on-chain feegrant, so a Privy-managed worker can hold no ALLO of its own. The granter is
the Forge **master wallet address**: forge-v2 auto-creates a feegrant from it to each new
signing wallet.

A `RemoteSigner` **discovers that address at runtime**: the signing-wallet GET and the
provision response now carry a `master_granter` field, which the signer reads from the
backend response and exposes — as the raw string, verbatim — via `signer.MasterGranter()`
(empty when the backend has no master wallet configured). `MasterGranter()` does not parse
or validate the value; parsing into an `sdk.AccAddress` happens in `signer.ResolveFeeGranter()`
(below). Runtime discovery means a master-wallet rotation no longer forces every SDK consumer
to reconfigure.

For 12-factor deployments you can **override** discovery with the canonical
`FORGE_MASTER_GRANTER_ADDRESS` environment variable (the same name used by allora-sdk-py
and allora-sdk-ts), read and parsed by `allora.FeeGranterFromEnv()` (returns `(nil, nil)`
when unset). The former name `FEE_GRANTER` is still accepted as a fallback for one release
(with a one-time deprecation warning), for parity with allora-sdk-py; rename it to
`FORGE_MASTER_GRANTER_ADDRESS`.

`signer.ResolveFeeGranter()` applies the precedence for you and returns a parsed
`sdk.AccAddress` — env override first, then the discovered granter, then `nil` (the signing
wallet pays its own fees) — so it drops straight into `TxParams.FeeGranter`:

```go
granter, err := signer.ResolveFeeGranter()
if err != nil {
    return err
}
params.FeeGranter = granter // nil ⇒ the signing wallet pays its own fees
```

To wire the precedence yourself instead — preferring an explicit env value and falling back
to the discovered granter:

```go
granter, err := allora.FeeGranterFromEnv() // FORGE_MASTER_GRANTER_ADDRESS, nil when unset
if err != nil {
    return err
}
if granter == nil {
    if discovered := signer.MasterGranter(); discovered != "" {
        if granter, err = sdk.AccAddressFromBech32(discovered); err != nil {
            return err
        }
    }
}
params.FeeGranter = granter // nil ⇒ the signing wallet pays its own fees
```

### Topic-bound provisioning (one worker = one topic)

`NewRemoteSigner` binds to an existing wallet id you already hold. For a worker that should
get-or-create a wallet bound to a specific topic (ENGN-8572 "one worker = one topic"), use
`NewRemoteSignerForTopic` instead — it idempotently provisions a managed wallet for the
given topic and returns a `RemoteSigner` built directly from the provision response (no
second wallet-info round-trip):

```go
signer, err := allora.NewRemoteSignerForTopic(ctx, allora.RemoteSignerConfig{
    BackendURL: "https://forge.allora.network",
    APIKey:     os.Getenv("FORGE_API_KEY"),
    // WalletID must be empty — it is filled from the provision response.
}, topicID, "my-worker")
if err != nil {
    panic(err)
}
```

Each topic binding counts against a per-user wallet cap enforced by the backend. When a
worker is retired or rotated to a different topic, release its binding so the slot stops
counting against the cap.

### Unbinding and revoking wallets

A `*RemoteSigner` exposes two lifecycle methods for the wallet it wraps:

- **`signer.ClearAssociation(ctx)`** — unbinds the wallet from its topic (Forge-side
  bookkeeping only; it does NOT unregister the worker on-chain). Reversible by
  re-provisioning. Use it before re-provisioning a wallet against a new topic or before
  decommissioning it.
- **`signer.Revoke(ctx)`** — permanently decommissions the wallet (DELETE
  `/api/v1/signing-wallets/{id}`). Destructive counterpart to `ClearAssociation`: clearing
  only unbinds, revoking tears the wallet down for good.

When you only hold the wallet id (e.g. retiring a worker without constructing a signer —
no wallet-info fetch), use the standalone by-id variants:

```go
// Unbind a wallet from its topic by id:
err := allora.ClearWalletAssociation(ctx, allora.RemoteSignerConfig{
    BackendURL: "https://forge.allora.network",
    APIKey:     os.Getenv("FORGE_API_KEY"),
}, walletID)

// Permanently decommission a wallet by id:
err := allora.RevokeWallet(ctx, allora.RemoteSignerConfig{
    BackendURL: "https://forge.allora.network",
    APIKey:     os.Getenv("FORGE_API_KEY"),
}, walletID)
```

Both return a non-2xx response (e.g. 404 for an unknown/foreign/already-cleared wallet) as
an error so the caller decides whether the failure is fatal or best-effort.

## Configuration

### Client Configuration

```go
type ClientConfig struct {
    // List of RPC endpoints to connect to
    Endpoints []EndpointConfig
    
    // Timeout for individual requests (default: 30s)
    RequestTimeout time.Duration
    
    // Connection timeout (default: 10s)
    ConnectionTimeout time.Duration
}

type EndpointConfig struct {
    // URL of the endpoint (e.g., "grpc://localhost:9090")
    URL string
    
    // Protocol to use ("grpc" is required for now)
    Protocol string
}
```

### Default Configuration

```go
config := allora.DefaultClientConfig()
// Provides sensible defaults but you'll need to add endpoints
config.Endpoints = []allora.EndpointConfig{
    {URL: "grpc://localhost:9090", Protocol: "grpc"},
}
```

## API Reference

### Cosmos SDK Modules

#### Auth Module
- `AuthAccount(ctx, address)` - Get account information
- `AuthAccounts(ctx, pagination)` - List all accounts
- `AuthParams(ctx)` - Get auth module parameters
- `AuthModuleAccount(ctx, name)` - Get module account by name
- `AuthModuleAccounts(ctx)` - List all module accounts

#### Bank Module
- `BankBalance(ctx, address, denom)` - Get account balance for specific denom
- `BankAllBalances(ctx, address, pagination)` - Get all balances for account
- `BankSpendableBalances(ctx, address, pagination)` - Get spendable balances
- `BankTotalSupply(ctx, pagination)` - Get total supply of all tokens
- `BankSupplyOf(ctx, denom)` - Get supply of specific denom
- `BankParams(ctx)` - Get bank module parameters
- `BankDenomMetadata(ctx, denom)` - Get metadata for specific denom
- `BankDenomsMetadata(ctx, pagination)` - Get metadata for all denoms
- `BankDenomOwners(ctx, denom, pagination)` - Get owners of specific denom

#### Staking Module
- `StakingValidators(ctx, status, pagination)` - List validators
- `StakingValidator(ctx, validatorAddr)` - Get specific validator
- `StakingValidatorDelegations(ctx, validatorAddr, pagination)` - Get validator delegations
- `StakingDelegation(ctx, delegatorAddr, validatorAddr)` - Get specific delegation
- `StakingDelegatorDelegations(ctx, delegatorAddr, pagination)` - Get all delegator delegations
- `StakingPool(ctx)` - Get staking pool information
- `StakingParams(ctx)` - Get staking parameters
- And more...

#### Tendermint/CometBFT Service
- `GetNodeInfo(ctx)` - Get node information
- `GetSyncing(ctx)` - Get sync status
- `GetLatestBlock(ctx)` - Get latest block
- `GetBlockByHeight(ctx, height)` - Get block by height
- `GetLatestValidatorSet(ctx, pagination)` - Get latest validator set
- `GetValidatorSetByHeight(ctx, height, pagination)` - Get validator set by height
- `Status(ctx)` - Get node status

### Allora Modules

#### Emissions Module (v9)
100+ query methods including:
- `EmissionsGetParams(ctx)` - Get emissions parameters
- `EmissionsGetTopic(ctx, topicId)` - Get topic information
- `EmissionsGetTotalStake(ctx)` - Get total stake in system
- `EmissionsTopicExists(ctx, topicId)` - Check if topic exists
- `EmissionsIsTopicActive(ctx, topicId)` - Check if topic is active
- And many more...

#### Mint Module (v5)
- `MintParams(ctx)` - Get mint parameters
- `MintInflation(ctx)` - Get current inflation rate
- `MintAnnualProvisions(ctx)` - Get annual provisions

**Note**: Allora module methods currently return `interface{}` until protobuf types are generated. They will be updated to return strongly-typed responses in future versions.

## Load Balancing & Fault Tolerance

The client automatically handles:

- **Round-robin load balancing** among healthy endpoints
- **Automatic failover** when endpoints become unavailable  
- **Exponential backoff** for failed endpoints
- **Health monitoring** with automatic recovery
- **Request retries** with configurable timeouts

### Monitoring Client Health

```go
// Get health status of all clients in the pool
healthStatus := client.GetHealthStatus()
fmt.Printf("Active clients: %v\\n", healthStatus["active_clients"])
fmt.Printf("Cooling clients: %v\\n", healthStatus["cooling_clients"])
```

## Error Handling

The client provides comprehensive error handling:

```go
balance, err := client.BankBalance(ctx, address, denom)
if err != nil {
    // Handle different types of errors
    switch {
    case strings.Contains(err.Error(), "not found"):
        fmt.Println("Account not found")
    case strings.Contains(err.Error(), "connection refused"):
        fmt.Println("Cannot connect to any endpoints")
    default:
        fmt.Printf("Other error: %v\\n", err)
    }
    return
}
```

## Pagination

Many query methods support pagination:

```go
// Get first 50 validators
validators, err := client.StakingValidators(ctx, "", &allora.PageRequest{
    Limit:      50,
    CountTotal: true,
})

// Get next page using key from previous response
if validators.Pagination.NextKey != nil {
    nextValidators, err := client.StakingValidators(ctx, "", &allora.PageRequest{
        Key:   validators.Pagination.NextKey,
        Limit: 50,
    })
}
```

## Development

### Building

```bash
make build
```

### Running Tests

```bash
make test
```

### Linting

```bash
make lint
```

### Protocol Buffer Generation

```bash
make proto-gen
```

## Roadmap

- [ ] Complete protobuf type generation for Allora modules
- [ ] Add support for JSON-RPC endpoints
- [ ] Add support for Cosmos LCD/REST endpoints
- [ ] Add caching layer for frequently accessed data
- [ ] Add metrics and monitoring integration
- [ ] Add transaction broadcasting support

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.