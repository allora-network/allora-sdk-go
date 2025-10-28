package allora

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/allora-network/allora-sdk-go/codec"
)

// AccountInfo contains the account number and sequence for a given address
type AccountInfo struct {
	AccountNumber uint64
	Sequence      uint64
}

// QueryAccountInfo queries the blockchain for account information needed for transaction signing
//
// This function retrieves the account number and sequence from the auth module.
// These values are required when building transactions.
//
// Parameters:
//   - ctx: Context for the query
//   - client: An Allora SDK client instance
//   - address: The address to query
//
// Returns:
//   - *AccountInfo: Account number and sequence
//   - error: Any error that occurred
//
// Example:
//
//	client, _ := allora.NewClient(config, logger)
//	addr, _ := sdk.AccAddressFromBech32("allo1...")
//
//	info, err := allora.QueryAccountInfo(ctx, client, addr)
//	if err != nil {
//	    panic(err)
//	}
//
//	fmt.Printf("Account Number: %d, Sequence: %d\n", info.AccountNumber, info.Sequence)
func QueryAccountInfo(ctx context.Context, client Client, address sdk.AccAddress) (*AccountInfo, error) {
	if address.Empty() {
		return nil, fmt.Errorf("address is required")
	}

	// Query account from the auth module
	req := &authtypes.QueryAccountRequest{
		Address: address.String(),
	}

	resp, err := client.Cosmos().Auth().Account(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to query account: %w", err)
	}

	// Unpack the account
	var account authtypes.AccountI
	alloraCodec := codec.CosmosCodec()
	err = alloraCodec.UnpackAny(resp.Account, &account)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack account: %w", err)
	}

	return &AccountInfo{
		AccountNumber: account.GetAccountNumber(),
		Sequence:      account.GetSequence(),
	}, nil
}

// GetChainID queries the blockchain for its chain ID
//
// The chain ID is required when signing transactions. This function retrieves it
// from the Tendermint node info.
//
// Parameters:
//   - ctx: Context for the query
//   - client: An Allora SDK client instance
//
// Returns:
//   - string: The chain ID (e.g., "allora-mainnet-1")
//   - error: Any error that occurred
//
// Example:
//
//	client, _ := allora.NewClient(config, logger)
//	chainID, err := allora.GetChainID(ctx, client)
//	if err != nil {
//	    panic(err)
//	}
//
//	fmt.Printf("Chain ID: %s\n", chainID)
func GetChainID(ctx context.Context, client Client) (string, error) {
	// Use Tendermint RPC to get node info
	status, err := client.Tendermint().Status(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to query status: %w", err)
	}

	return status.NodeInfo.Network, nil
}

// TxParamsBuilder provides a fluent API for building TxParams
//
// This builder helps construct TxParams with values either manually provided
// or automatically queried from the blockchain.
//
// Example:
//
//	client, _ := allora.NewClient(config, logger)
//	addr, _ := sdk.AccAddressFromBech32("allo1...")
//
//	params, err := allora.NewTxParamsBuilder(ctx, client).
//	    WithAddress(addr).
//	    WithGasLimit(250000).
//	    WithFee(sdk.NewCoins(sdk.NewInt64Coin("uallo", 6000))).
//	    WithMemo("my-transaction").
//	    QueryAndBuild()
type TxParamsBuilder struct {
	ctx    context.Context
	client Client

	// Manually set values
	chainID       *string
	accountNumber *uint64
	sequence      *uint64
	gasLimit      uint64
	feeAmount     sdk.Coins
	memo          string
	timeoutHeight uint64

	// Address to query if account info not manually set
	address sdk.AccAddress
}

// NewTxParamsBuilder creates a new TxParamsBuilder
func NewTxParamsBuilder(ctx context.Context, client Client) *TxParamsBuilder {
	return &TxParamsBuilder{
		ctx:      ctx,
		client:   client,
		gasLimit: 200000, // Default
		feeAmount: sdk.NewCoins(sdk.NewInt64Coin("uallo", 5000)), // Default
	}
}

// WithAddress sets the address to query for account information
func (b *TxParamsBuilder) WithAddress(address sdk.AccAddress) *TxParamsBuilder {
	b.address = address
	return b
}

// WithChainID manually sets the chain ID (skips automatic query)
func (b *TxParamsBuilder) WithChainID(chainID string) *TxParamsBuilder {
	b.chainID = &chainID
	return b
}

// WithAccountNumber manually sets the account number (skips automatic query)
func (b *TxParamsBuilder) WithAccountNumber(accountNumber uint64) *TxParamsBuilder {
	b.accountNumber = &accountNumber
	return b
}

// WithSequence manually sets the sequence (skips automatic query)
func (b *TxParamsBuilder) WithSequence(sequence uint64) *TxParamsBuilder {
	b.sequence = &sequence
	return b
}

// WithGasLimit sets the gas limit
func (b *TxParamsBuilder) WithGasLimit(gasLimit uint64) *TxParamsBuilder {
	b.gasLimit = gasLimit
	return b
}

// WithFee sets the transaction fee
func (b *TxParamsBuilder) WithFee(feeAmount sdk.Coins) *TxParamsBuilder {
	b.feeAmount = feeAmount
	return b
}

// WithMemo sets the transaction memo
func (b *TxParamsBuilder) WithMemo(memo string) *TxParamsBuilder {
	b.memo = memo
	return b
}

// WithTimeoutHeight sets the timeout height
func (b *TxParamsBuilder) WithTimeoutHeight(height uint64) *TxParamsBuilder {
	b.timeoutHeight = height
	return b
}

// QueryAndBuild queries any missing parameters from the blockchain and builds TxParams
//
// This method will query:
// - Chain ID (if not manually set)
// - Account number and sequence (if not manually set and address is provided)
//
// Returns:
//   - *TxParams: The constructed transaction parameters
//   - error: Any error that occurred during queries
func (b *TxParamsBuilder) QueryAndBuild() (*TxParams, error) {
	params := &TxParams{
		GasLimit:      b.gasLimit,
		FeeAmount:     b.feeAmount,
		Memo:          b.memo,
		TimeoutHeight: b.timeoutHeight,
	}

	// Query chain ID if not manually set
	if b.chainID == nil {
		chainID, err := GetChainID(b.ctx, b.client)
		if err != nil {
			return nil, fmt.Errorf("failed to query chain ID: %w", err)
		}
		params.ChainID = chainID
	} else {
		params.ChainID = *b.chainID
	}

	// Query account info if not manually set
	if b.accountNumber == nil || b.sequence == nil {
		if b.address.Empty() {
			return nil, fmt.Errorf("address is required to query account info")
		}

		accountInfo, err := QueryAccountInfo(b.ctx, b.client, b.address)
		if err != nil {
			return nil, fmt.Errorf("failed to query account info: %w", err)
		}

		if b.accountNumber == nil {
			params.AccountNumber = accountInfo.AccountNumber
		}
		if b.sequence == nil {
			params.Sequence = accountInfo.Sequence
		}
	}

	// Use manually set values if provided
	if b.accountNumber != nil {
		params.AccountNumber = *b.accountNumber
	}
	if b.sequence != nil {
		params.Sequence = *b.sequence
	}

	return params, nil
}

// Build builds TxParams without querying the blockchain
//
// This method is useful for offline signing where all parameters are manually provided.
// It will return an error if required parameters are missing.
//
// Returns:
//   - *TxParams: The constructed transaction parameters
//   - error: Error if required parameters are missing
func (b *TxParamsBuilder) Build() (*TxParams, error) {
	if b.chainID == nil {
		return nil, fmt.Errorf("chain ID is required")
	}
	if b.accountNumber == nil {
		return nil, fmt.Errorf("account number is required")
	}
	if b.sequence == nil {
		return nil, fmt.Errorf("sequence is required")
	}

	params := &TxParams{
		ChainID:       *b.chainID,
		AccountNumber: *b.accountNumber,
		Sequence:      *b.sequence,
		GasLimit:      b.gasLimit,
		FeeAmount:     b.feeAmount,
		Memo:          b.memo,
		TimeoutHeight: b.timeoutHeight,
	}

	return params, nil
}
