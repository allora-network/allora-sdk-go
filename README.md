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