package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"

	emissionstypes "github.com/allora-network/allora-chain/x/emissions/types"
	minttypes "github.com/allora-network/allora-chain/x/mint/types"
	allorasdk "github.com/allora-network/allora-sdk-go"
	butils "github.com/brynbellomy/go-utils"
	ctypes "github.com/cometbft/cometbft/types"
	cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/tmrpc"
)

func main() {
	// Setup logger
	consoleOutput := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(consoleOutput).With().Caller().Timestamp().Logger()

	// Create client configuration with mixed protocols
	cfg := &config.ClientConfig{
		Endpoints: []config.EndpointConfig{
			{
				URL:      "allora-grpc.testnet.allora.network:443",
				Protocol: config.ProtocolGRPC,
			},
			{
				URL:      "https://allora-api.testnet.allora.network",
				Protocol: config.ProtocolREST,
			},
			{
				URL:          "https://allora-rpc.testnet.allora.network",
				WebsocketURL: "https://allora-rpc.testnet.allora.network/websocket",
				Protocol:     config.ProtocolTendermintRPC,
			},
		},
		RequestTimeout:    30 * time.Second,
		ConnectionTimeout: 10 * time.Second,
	}

	// Create the Allora client
	client, err := allorasdk.NewClient(cfg, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create Allora client")
	}
	defer client.Close()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.Cosmos().Mint().EmissionInfo(ctx, &minttypes.QueryServiceEmissionInfoRequest{}, config.Height(5601386))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create Allora client")
	}
	fmt.Println(butils.PrettyJSON(resp))

	// Execute multiple requests to test all endpoints through round-robin
	fmt.Println("🔄 Testing gRPC and JSON-RPC endpoints with comprehensive module demonstrations...")

	// Request 1: Get staking params (demonstrating new module-based API)
	fmt.Println("\n🔍 Request 1: Getting staking parameters...")
	stakingParams, err := client.Cosmos().Staking().Params(ctx, &stakingtypes.QueryParamsRequest{})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get staking params")
	} else {
		fmt.Printf("✅ Unbonding time: %s\n", stakingParams.Params.UnbondingTime)
		fmt.Printf("✅ Max validators: %d\n", stakingParams.Params.MaxValidators)
	}

	// Request 2: Get staking pool (demonstrating new module-based API)
	fmt.Println("\n📦 Request 2: Getting staking pool...")
	stakingPool, err := client.Cosmos().Staking().Pool(ctx, &stakingtypes.QueryPoolRequest{})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get staking pool")
	} else {
		fmt.Printf("✅ Bonded tokens: %s\n", stakingPool.Pool.BondedTokens)
		fmt.Printf("✅ Not bonded tokens: %s\n", stakingPool.Pool.NotBondedTokens)
	}

	// Request 3: Get validators (demonstrating new module-based API)
	fmt.Println("\n👥 Request 3: Getting validators...")
	validators, err := client.Cosmos().Staking().Validators(ctx, &stakingtypes.QueryValidatorsRequest{
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
		pool, err := client.Cosmos().Staking().Pool(ctx, &stakingtypes.QueryPoolRequest{})
		if err != nil {
			logger.Error().Err(err).Int("request", i).Msg("failed to get staking pool")
		} else {
			fmt.Printf("✅ Request %d successful - Bonded tokens: %s\n", i, pool.Pool.BondedTokens)
		}
	}

	// Request 9: Demonstrate Bank module - Get total supply
	fmt.Println("\n💰 Request 9: Getting total token supply (Bank module)...")
	totalSupply, err := client.Cosmos().Bank().TotalSupply(ctx, &banktypes.QueryTotalSupplyRequest{
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
	syncStatus, err := client.Cosmos().Tendermint().GetSyncing(ctx, &cmtservice.GetSyncingRequest{})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get sync status")
	} else {
		fmt.Printf("✅ Node syncing: %v\n", syncStatus.Syncing)
	}

	// Request 11: Demonstrate Emissions module - Get parameters
	fmt.Println("\n🔥 Request 11: Getting emissions parameters (Emissions module)...")
	emissionsParams, err := client.Cosmos().Emissions().GetParams(ctx, &emissionstypes.GetParamsRequest{})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get emissions params")
	} else {
		fmt.Printf("✅ Emissions params retrieved successfully\n")
		// Avoid potential type issues by just confirming retrieval
		_ = emissionsParams
	}

	// Request 12: Demonstrate Mint module - Get mint parameters
	fmt.Println("\n🌟 Request 12: Getting mint parameters (Mint module)...")
	_, err = client.Cosmos().Mint().Params(ctx, &minttypes.QueryServiceParamsRequest{})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get mint params")
	} else {
		fmt.Printf("✅ Mint params retrieved successfully\n")
	}

	// Show client pool health status
	fmt.Println("\n📊 Client pool health status:")
	healthStatus := client.GetHealthStatus()
	fmt.Printf("✅ Active clients: %v\n", healthStatus["active_clients"])
	fmt.Printf("✅ Cooling clients: %v\n", healthStatus["cooling_clients"])
	if currentClient, ok := healthStatus["current_client"]; ok {
		fmt.Printf("✅ Current client: %s\n", currentClient)
	}

	block, err := client.Tendermint().Block(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Msg("[tmrpc] failed to get latest block")
	} else if block == nil {
		logger.Error().Msg("[tmrpc] latest block is nil")
	} else {
		fmt.Printf("✅ [tmrpc] latest block retrieved successfully\n")
	}

	height := block.Block.Header.Height - 1
	block, err = client.Tendermint().Block(ctx, &height)
	if err != nil {
		logger.Error().Err(err).Msg("[tmrpc] failed to get recent block by height")
	} else if block == nil {
		logger.Error().Msg("[tmrpc] recent block by height is nil")
	} else if block.Block.Header.Height != height {
		logger.Error().Int64("expected_height", height).Int64("actual_height", block.Block.Header.Height).Msg("[tmrpc] recent block by height mismatch")
	} else {
		fmt.Printf("✅ [tmrpc] recent block by height retrieved successfully\n")
	}

	mb := tmrpc.NewMailbox(100)
	client.Subscribe(mb, "tm.event='NewBlock'")

	chBlock := make(chan struct{})

	go func() {
		defer close(chBlock)
		select {
		case <-mb.Notify():
			evt, ok := mb.Retrieve()
			if !ok {
				logger.Error().Msg("failed to cast event data to NewBlock")
				return
			}
			blockEvt, ok := evt.(ctypes.EventDataNewBlock)
			if !ok {
				logger.Error().Msg("failed to cast event data to NewBlock")
				return
			}
			fmt.Printf("✅ Received new block event: %v\n", blockEvt.Block.Height)
		case <-time.After(15 * time.Second):
			logger.Error().Msg("timeout waiting for new block event")
		}
	}()

	<-chBlock

	fmt.Println("\n🎉 Example completed successfully!")
	fmt.Println("📋 Demonstrated modules: Staking, Bank, Tendermint, Emissions, Mint")
	fmt.Printf("📊 Total requests made: 12 (demonstrating load balancing across %d endpoints)\n", len(cfg.Endpoints))
	fmt.Println("🚀 Mixed protocol support: gRPC and JSON-RPC clients in the same pool")
}
