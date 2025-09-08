package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"

	emissionstypes "github.com/allora-network/allora-chain/x/emissions/types"
	minttypes "github.com/allora-network/allora-chain/x/mint/types"
	allorasdk "github.com/allora-network/allora-sdk-go"
	cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/allora-network/allora-sdk-go/config"
)

func main() {
	// Setup logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Create client configuration with mixed protocols
	config := &config.ClientConfig{
		Endpoints: []config.EndpointConfig{
			{
				URL:      "grpc://allora-grpc.testnet.allora.network:443",
				Protocol: config.ProtocolGRPC,
			},
			// {
			// 	URL:      "https://allora-api.testnet.allora.network",
			// 	Protocol: config.ProtocolREST,
			// },
			// {
			// 	URL:      "http://100.73.209.49:1317",
			// 	Protocol: config.ProtocolREST,
			// },
			// {
			// 	URL:      "100.73.209.49:9090",
			// 	Protocol: config.ProtocolGRPC,
			// },
			// {
			// 	URL:      "http://100.75.16.36:1317",
			// 	Protocol: config.ProtocolREST,
			// },
			// {
			// 	URL:      "100.75.16.36:9090",
			// 	Protocol: config.ProtocolGRPC,
			// },
		},
		RequestTimeout:    30 * time.Second,
		ConnectionTimeout: 10 * time.Second,
	}

	// Create the Allora client
	client, err := allorasdk.NewClient(config, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create Allora client")
	}
	defer client.Close()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Execute multiple requests to test all endpoints through round-robin
	fmt.Println("🔄 Testing gRPC and JSON-RPC endpoints with comprehensive module demonstrations...")

	// Request 1: Get staking params (demonstrating new module-based API)
	fmt.Println("\n🔍 Request 1: Getting staking parameters...")
	stakingParams, err := client.Staking().Params(ctx, &stakingtypes.QueryParamsRequest{})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get staking params")
	} else {
		fmt.Printf("✅ Unbonding time: %s\n", stakingParams.Params.UnbondingTime)
		fmt.Printf("✅ Max validators: %d\n", stakingParams.Params.MaxValidators)
	}

	// Request 2: Get staking pool (demonstrating new module-based API)
	fmt.Println("\n📦 Request 2: Getting staking pool...")
	stakingPool, err := client.Staking().Pool(ctx, &stakingtypes.QueryPoolRequest{})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get staking pool")
	} else {
		fmt.Printf("✅ Bonded tokens: %s\n", stakingPool.Pool.BondedTokens)
		fmt.Printf("✅ Not bonded tokens: %s\n", stakingPool.Pool.NotBondedTokens)
	}

	// Request 3: Get validators (demonstrating new module-based API)
	fmt.Println("\n👥 Request 3: Getting validators...")
	validators, err := client.Staking().Validators(ctx, &stakingtypes.QueryValidatorsRequest{
		Status: "",
		Pagination: &query.PageRequest{
			Limit: 5, // Get only 5 validators for this example
		},
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get validators")
	} else {
		fmt.Printf("✅ Found %d validators (showing first 5):\n", len(validators.Validators))
		for i, validator := range validators.Validators {
			if i >= 5 { // Safety check
				break
			}
			fmt.Printf("  %d. %s (status: %s)\n",
				i+1,
				validator.Description.Moniker,
				validator.Status.String())
		}
	}

	// Additional requests to demonstrate load balancing across both endpoints
	fmt.Println("\n🔄 Request 4-8: Making additional requests to test load balancing...")

	for i := 4; i <= 8; i++ {
		fmt.Printf("\n📊 Request %d: Getting staking pool...\n", i)
		pool, err := client.Staking().Pool(ctx, &stakingtypes.QueryPoolRequest{})
		if err != nil {
			logger.Error().Err(err).Int("request", i).Msg("failed to get staking pool")
		} else {
			fmt.Printf("✅ Request %d successful - Bonded tokens: %s\n", i, pool.Pool.BondedTokens)
		}
	}

	// Request 9: Demonstrate Bank module - Get total supply
	fmt.Println("\n💰 Request 9: Getting total token supply (Bank module)...")
	totalSupply, err := client.Bank().TotalSupply(ctx, &banktypes.QueryTotalSupplyRequest{
		Pagination: &query.PageRequest{Limit: 10},
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get total supply")
	} else {
		fmt.Printf("✅ Found %d token denominations\n", len(totalSupply.Supply))
		if len(totalSupply.Supply) > 0 {
			fmt.Printf("✅ Example: %s has supply %s\n", totalSupply.Supply[0].Denom, totalSupply.Supply[0].Amount)
		}
	}

	// Request 10: Demonstrate Tendermint module - Get node sync status
	fmt.Println("\n🌐 Request 10: Getting node sync status (Tendermint module)...")
	syncStatus, err := client.Tendermint().GetSyncing(ctx, &cmtservice.GetSyncingRequest{})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get sync status")
	} else {
		fmt.Printf("✅ Node syncing: %v\n", syncStatus.Syncing)
	}

	// Request 11: Demonstrate Emissions module - Get parameters
	fmt.Println("\n🔥 Request 11: Getting emissions parameters (Emissions module)...")
	emissionsParams, err := client.Emissions().GetParams(ctx, &emissionstypes.GetParamsRequest{})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get emissions params")
	} else {
		fmt.Printf("✅ Emissions params retrieved successfully\n")
		// Avoid potential type issues by just confirming retrieval
		_ = emissionsParams
	}

	// Request 12: Demonstrate Mint module - Get mint parameters
	fmt.Println("\n🌟 Request 12: Getting mint parameters (Mint module)...")
	_, err = client.Mint().Params(ctx, &minttypes.QueryServiceParamsRequest{})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get mint params")
	} else {
		fmt.Printf("✅ Mint params retrieved successfully\n")
	}

	denomOwners, err := client.Bank().DenomOwners(ctx, &banktypes.QueryDenomOwnersRequest{
		Denom:      "uallo",
		Pagination: &query.PageRequest{Offset: 0, Limit: 5},
	})
	bs, _ := json.MarshalIndent(denomOwners, "", "    ")
	fmt.Println("YOOOO", string(bs))

	denomOwners, err = client.Bank().DenomOwners(ctx, &banktypes.QueryDenomOwnersRequest{
		Denom:      "uallo",
		Pagination: &query.PageRequest{Offset: 1, Limit: 10},
	})
	bs, _ = json.MarshalIndent(denomOwners, "", "    ")
	fmt.Println("HEYYYYY", string(bs))

	// Show client pool health status
	fmt.Println("\n📊 Client pool health status:")
	healthStatus := client.GetHealthStatus()
	fmt.Printf("✅ Active clients: %v\n", healthStatus["active_clients"])
	fmt.Printf("✅ Cooling clients: %v\n", healthStatus["cooling_clients"])
	if currentClient, ok := healthStatus["current_client"]; ok {
		fmt.Printf("✅ Current client: %s\n", currentClient)
	}

	fmt.Println("\n🎉 Example completed successfully!")
	fmt.Println("📋 Demonstrated modules: Staking, Bank, Tendermint, Emissions, Mint")
	fmt.Printf("📊 Total requests made: 12 (demonstrating load balancing across %d endpoints)\n", len(config.Endpoints))
	fmt.Println("🚀 Mixed protocol support: gRPC and JSON-RPC clients in the same pool")
}
