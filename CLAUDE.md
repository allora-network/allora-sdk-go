# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

The Allora SDK Go is a high-performance Go client SDK for the Allora Network blockchain, providing comprehensive access to query services across Allora custom modules and Cosmos SDK base modules. The SDK supports multiple RPC endpoints with automatic load balancing, fault tolerance, and health tracking.

## Development Commands

### Building and Testing
```bash
# Build the SDK
make build

# Install the SDK
make install

# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Clean build artifacts
make clean
```

### Linting and Code Quality
```bash
# Run golangci-lint
make lint

# Run golangci-lint with auto-fix
make lint-fix
```

### Protocol Buffer Generation
```bash
# Generate protobuf code (requires Docker)
make proto-gen

# Format protobuf files
make proto-format

# Lint protobuf files  
make proto-lint

# Both format and generate
make proto-all
```

### Code Generation
```bash
# Generate all module client wrappers
go run cmd/codegen/main.go generate-all

# Generate specific module wrapper
go run cmd/codegen/main.go generate <module>

# List available modules
go run cmd/codegen/main.go list
```

## Architecture Overview

### Core Components

1. **Client (`client_facade.go`)** - Main client interface that aggregates all module clients with health management
2. **ClientPoolManager (`client_pool_manager.go`)** - Manages pools of clients with load balancing, health tracking, and automatic failover
3. **AggregatedClient Interface (`client.go`)** - Protocol abstraction layer supporting both gRPC and JSON-RPC
4. **Protocol Implementations**:
   - **GRPCAggregatedClient (`grpc_client.go`)** - gRPC protocol implementation
   - **JSONRPCAggregatedClient (`jsonrpc_aggregated_client.go` + `jsonrpc_client.go`)** - JSON-RPC protocol implementation

### Module Structure

Each blockchain module has dedicated client files:
- `client.<module>.grpc.go` - gRPC implementation for the module
- `client.<module>.jsonrpc.go` - JSON-RPC implementation for the module (if HTTP annotations exist)

**Supported Modules:**
- **Cosmos SDK Modules**: auth, bank, staking, tendermint
- **Allora Modules**: emissions (100+ methods), mint

### Code Generation System

The `cmd/codegen/main.go` tool automatically generates module client wrappers by:
1. Parsing protobuf service definitions (`.proto` files) to extract HTTP annotations
2. Falling back to Go protobuf files (`.pb.go`) if proto sources unavailable
3. Generating both gRPC and JSON-RPC clients based on available annotations

### Dependencies and Proto Management

- **Proto Dependencies**: Located in `proto-deps/` directory containing Cosmos SDK and Allora Chain protobuf definitions
- **Code Generation**: Uses Docker containers for consistent protobuf compilation
- **Buf Configuration**: `buf.yaml` and `buf.gen.yaml` configure protobuf build pipeline

### Health Management and Fault Tolerance

The ClientPoolManager implements sophisticated health tracking:
- **Active Pool**: Healthy clients available for requests
- **Cooling Pool**: Temporarily unhealthy clients with exponential backoff
- **Round-robin Load Balancing**: Distributes requests across healthy clients
- **Automatic Recovery**: Failed clients are periodically retried and reactivated
- **Request Retry Logic**: `execute_with_retry.go` handles automatic retries with exponential backoff

### Configuration

Clients are configured via:
```go
type ClientConfig struct {
    Endpoints []EndpointConfig     // Multiple RPC endpoints
    RequestTimeout time.Duration   // Individual request timeout
    ConnectionTimeout time.Duration // Connection timeout
}

type EndpointConfig struct {
    URL string      // e.g., "grpc://localhost:9090"
    Protocol string // "grpc" or "json-rpc"  
}
```

## Development Patterns

### Adding New Modules

1. Add module configuration to `cmd/codegen/main.go` in `getModuleConfigs()`
2. Run code generator: `go run cmd/codegen/main.go generate <module>`
3. Update main client interface if needed
4. Add wrapper instantiation to `client_facade.go`

### Testing Patterns

- Use `client_test.go` as reference for integration tests
- Mock the AggregatedClient interface for unit tests
- Test both happy path and error scenarios including client failures

### Protocol Support

- **gRPC**: Full support for all module methods
- **JSON-RPC**: Only methods with HTTP annotations from protobuf definitions
- Protocol selection is per-endpoint, allowing mixed configurations

### Error Handling

- All client methods use the `executeWithRetry` pattern
- Errors from individual clients trigger health status updates
- Failed clients are moved to cooling pool with backoff delays
- Comprehensive logging via zerolog with structured context