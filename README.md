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

## Transaction sending

The SDK provides a single high-level entry point — `client.Tx().SendTx(ctx, signer, msgs, opts)` — that drives a message set through the full send lifecycle: account discovery, build, sign, gas estimation, broadcast, and on-chain confirmation. It is the path most callers want; the lower-level `CreateUnsignedTx` / `SignTransactionWith` primitives are also exposed for the two-phase build-now / sign-later flow (below).

### The one-call path

`SendTx` takes an `allora.Signer` (a local `Wallet`'s `PrivKey` or a `*RemoteSigner`), a `[]sdk.Msg`, and a `SendOptions` value. The returned `*txsend.TxResult` carries the committed tx's hash, block height, and DeliverTx code; a non-nil error means the chain rejected the tx and retries were exhausted.

```go
client, err := allora.NewClient(cfg, logger) // see Quick Start
if err != nil { /* ... */ }
defer client.Close()

// 1) build messages via the pure constructors in the txmsg package
//    (validation up front; returns sdk.Msg, no network calls).
sendMsg, err := txmsg.NewSend(fromAddr, toAddr, sdk.NewCoins(sdk.NewInt64Coin("uallo", 1_000_000)))
if err != nil { /* ... */ }
regMsg, err := txmsg.NewRegisterWorker(txmsg.RegisterWorkerParams{
    Sender:  fromAddr,
    TopicId: 42,
    Owner:   fromAddr,
})
if err != nil { /* ... */ }

msgs := []sdk.Msg{sendMsg, regMsg}

// 2) construct a Signer — *Wallet.PrivKey or *RemoteSigner both fit.
wallet, err := allora.NewWalletFromMnemonic(mnemonic, allora.DefaultHDPath)
if err != nil { /* ... */ }

// 3) optional: a master/subsidy wallet paying gas via an on-chain feegrant.
//    Empty (nil) means the signing wallet pays its own fees.
granter, err := allora.FeeGranterFromEnv() // reads FORGE_MASTER_GRANTER_ADDRESS
if err != nil { /* ... */ }

// 4) send. SendTx applies sensible defaults for any zero-valued
//    SendOptions field: GasPrice=0.01, GasAdjustment=1.5,
//    BroadcastMode=BroadcastModeSync, MaxRetries=2. Set GasLimit
//    explicitly to skip the simulate round-trip.
result, err := client.Tx().SendTx(ctx, wallet.PrivKey, msgs, allora.SendOptions{
    ChainID:    "allora-testnet-1",
    FeeGranter: granter,
    Memo:       "deposit-12345",
})
if err != nil { /* chain rejected the tx */ }
fmt.Printf("tx %s committed at height %d (code %d)\n", result.TxHash, result.Height, result.Code)
```

`SendOptions` fields (all optional; zero values produce defaults):

| Field | Default | Notes |
|---|---|---|
| `ChainID` | **required** | The chain id (e.g. `allora-testnet-1`). SendTx errors if empty. |
| `Memo` | `""` | Attached to the tx body. |
| `FeeGranter` | `nil` | `sdk.AccAddress` of a feegrant granter (master/subsidy wallet). See below. |
| `GasAdjustment` | `1.5` | Multiplier applied to simulated gas. Clamped to `>= 1.0`. |
| `GasLimit` | `0` (simulate) | When non-zero, skips `EstimateGas` and uses this value directly. |
| `GasPrice` | `0.01` | Per-gas-unit price in uallo. |
| `BroadcastMode` | `BroadcastModeSync` | `BroadcastModeSync` / `Async` / `Block`. |
| `SkipWait` | `false` | When true, returns after broadcast (no `WaitForTx` poll). |
| `MaxRetries` | `2` | Retries for retryable CheckTx rejections. See "Retry behavior" below. |

### The two-phase path (build now, sign later)

When signing must be decoupled from building — e.g. an unsigned tx is persisted or queued on one host and signed on another, or signing is delegated to a remote service — the SDK exposes the lower-level primitives directly:

```go
// The high-level client.Tx() returns a Sender (one-call SendTx). For the
// lower-level two-phase primitives you need a TxBroadcaster directly:
broadcaster := cosmospool.New(client.Cosmos(), logger)

// Phase 1: build the unsigned tx. Account info MUST be fresh here, since
// the sequence is encoded into the tx.
accNum, seq, err := broadcaster.AccountInfo(ctx, fromAddr.String())
if err != nil { /* ... */ }

params := &allora.TxParams{
    ChainID:       "allora-testnet-1",
    AccountNumber: accNum,
    Sequence:      seq,
    GasLimit:      200_000,
    FeeAmount:     sdk.NewCoins(sdk.NewInt64Coin("uallo", 5_000)),
    FeeGranter:    granter, // optional
    Memo:          "deposit-12345",
}
unsignedTx, err := allora.CreateUnsignedTx(msgs, params)
if err != nil { /* ... */ }
// → persist unsignedTx (db, queue, file, ...)

// Phase 2: sign on a different host / process / service. Only ChainID,
// AccountNumber, and Sequence are read from params at sign time — the fee,
// gas, memo, granter, and timeout-height are already encoded in the tx.
signedTx, err := allora.SignTransactionWith(ctx, unsignedTx, signer, params)
if err != nil { /* ... */ }

// Phase 3: broadcast the signed bytes through the same Tx() entry point
// (use a Sender backed by the broadcaster of your choice, or wire
// signedTx into your own broadcast path).
```

Why split build and sign? Two common cases:

- **Remote signing.** The build step can run on a host that has the account/sequence info but does not hold the key; the sign step runs wherever the key lives (locally or behind a remote signer service). See "Privy-Managed Signing" below for the `*RemoteSigner` integration.
- **Replay batching.** A worker builds and persists many unsigned txs up front (sharing one sequence-refetch), then signs and broadcasts them serially — minimizing the wall-clock between `AccountInfo` and broadcast so the sequence is still valid.

### RemoteSigner (Privy) — any Signer works

`SendTx` takes an `allora.Signer`, an interface satisfied by both a local `Wallet`'s `*secp256k1.PrivKey` (the default for self-managed keys) and a `*RemoteSigner` (Privy-managed, key never leaves the Forge backend). The remote path is documented in full in the "Privy-Managed Signing" section below; the only change for `SendTx` is that you pass the `*RemoteSigner` directly where the example above shows `wallet.PrivKey`.

### Fee granter (gas subsidy)

`SendOptions.FeeGranter` is an `sdk.AccAddress` of a feegrant granter — a master/subsidy wallet that pays the fee on behalf of the signer via an on-chain feegrant. When set, the signer's own ALLO balance is not debited. Discovery + resolution helpers, in precedence order:

1. `allora.FeeGranterFromEnv()` — reads `FORGE_MASTER_GRANTER_ADDRESS` (and the legacy `FEE_GRANTER` fallback with a one-time deprecation warning). Returns `(nil, nil)` when unset.
2. `signer.(*RemoteSigner).ResolveFeeGranter()` — returns the backend-discovered master granter for a `*RemoteSigner`. Local wallets have no discovery path; the call is a no-op for them.
3. `nil` — the signing wallet pays its own fees.

The combined result drops straight into `SendOptions.FeeGranter`.

### Message constructors (txmsg)

The `txmsg` package provides pure, validated constructors that return `sdk.Msg` values without touching signing, broadcasting, or the network. They are the recommended way to build messages for `SendTx`. Highlights:

- **Bank** — `txmsg.NewSend(from, to, amount)`, `txmsg.NewMultiSend(inputs, outputs)` (input/output sum equality is checked up front).
- **Emissions** — `txmsg.NewRegisterWorker(RegisterWorkerParams{Sender, TopicId, Owner, IsReputer})`, `txmsg.NewInsertWorkerPayload(InsertWorkerPayloadParams{Sender, WorkerDataBundle})` (validates bundle nonce and signature non-empty).
- **Feegrant** — `txmsg.NewGrantAllowance(granter, grantee, allowance)` (validates `BasicAllowance.SpendLimit` and `PeriodicAllowance.PeriodSpendLimit` / `PeriodCanSpend`), `txmsg.NewRevokeAllowance(granter, grantee)`.

Every constructor validates its inputs (non-empty addresses, positive amounts, valid bech32 against the Allora prefix) so chain rejections of bad-shaped messages are replaced by up-front errors at the call site.

### Retry behavior

`SendTx` classifies CheckTx rejections by `(code, codespace)` and retries the build → sign → broadcast cycle automatically, up to `MaxRetries` attempts. The retryable codes (all in the canonical `sdk` codespace) are:

| Code | Meaning | Action on retry |
|---|---|---|
| `11` | out-of-gas | gas limit × `1.3` (`retryBumpGasFactor`) |
| `13` | insufficient fee | fee × `2.0` (`retryBumpFeeFactor`) |
| `32` | sequence mismatch | refetch account info, re-sign with the new sequence |

A non-retryable CheckTx rejection (any other code, or any code in a non-`sdk` codespace) is returned as an error immediately. When retries are exhausted, `SendTx` returns the last `*TxResult` alongside an error so a caller can inspect the chain's raw `Codespace` / `RawLog` for triage.

### A runnable example

`example/txsend/main.go` is a runnable end-to-end example that wires it all together: it constructs a `Client`, picks a local `Wallet` or a `*RemoteSigner` based on env vars, builds a `NewSend` + `NewRegisterWorker` message set with a fee-granter subsidy, and prints the committed `TxResult`. By default it prints usage and exits — set `ALLORA_RUN_EXAMPLE=1` (and the other documented env vars) to actually broadcast.

```bash
go run ./example/txsend                 # prints usage, does not dial
ALLORA_RUN_EXAMPLE=1 \
  ALLORA_CHAIN_ID=allora-testnet-1 \
  ALLORA_GRPC=https://allora-grpc.testnet.allora.network:443 \
  ALLORA_MNEMONIC="<24-word mnemonic>" \
  ALLORA_TO=allo1... \
  ALLORA_AMOUNT_UALLO=1000 \
  go run ./example/txsend
```

A tag-gated end-to-end integration test (`example/txsend/integration_test.go`, `//go:build integration`) exercises the same path against a live node; it `t.Skip`s when its `ALLORA_TEST_*` env vars are unset, so the default `go test ./...` is never affected:

```bash
mise exec -- go test -tags integration -run TestIntegrationSendTx -count=1 -v ./example/txsend
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