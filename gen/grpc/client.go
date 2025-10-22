package grpc

import (
	"context"
	"crypto/tls"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/brynbellomy/go-utils/errors"
	sdkgrpc "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/metadata"

	cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	cosmoscodec "github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	emissionstypes "github.com/allora-network/allora-chain/x/emissions/types"
	minttypes "github.com/allora-network/allora-chain/x/mint/types"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// GRPCClient implements the Client interface using gRPC
type GRPCClient struct {
	endpointURL  string
	conn         *grpc.ClientConn
	gov          *GovGRPCClient
	evidence     *EvidenceGRPCClient
	auth         *AuthGRPCClient
	tx           *TxGRPCClient
	slashing     *SlashingGRPCClient
	emissions    *EmissionsGRPCClient
	feegrant     *FeegrantGRPCClient
	mint         *MintGRPCClient
	staking      *StakingGRPCClient
	authz        *AuthzGRPCClient
	node         *NodeGRPCClient
	consensus    *ConsensusGRPCClient
	distribution *DistributionGRPCClient
	bank         *BankGRPCClient
	params       *ParamsGRPCClient
	tendermint   *TendermintGRPCClient
}

var _ interfaces.CosmosClient = (*GRPCClient)(nil)

var (
	grpcCodecOnce sync.Once
	grpcCodec     encoding.Codec
)

func buildGRPCCodec() encoding.Codec {
	grpcCodecOnce.Do(func() {
		registry := codectypes.NewInterfaceRegistry()
		registerFuncs := []func(codectypes.InterfaceRegistry){
			std.RegisterInterfaces,
			banktypes.RegisterInterfaces,
			stakingtypes.RegisterInterfaces,
			slashingtypes.RegisterInterfaces,
			distributiontypes.RegisterInterfaces,
			minttypes.RegisterInterfaces,
			emissionstypes.RegisterInterfaces,
		}
		for _, register := range registerFuncs {
			register(registry)
		}
		grpcCodec = cosmoscodec.NewProtoCodec(registry).GRPCCodec()
	})
	return grpcCodec
}

// NewGRPCClient creates a new gRPC aggregated client
func NewGRPCClient(cfg config.EndpointConfig, logger zerolog.Logger) (*GRPCClient, error) {
	if cfg.Protocol != "grpc" {
		return nil, errors.Errorf("unsupported protocol: %s, only 'grpc' is supported", cfg.Protocol)
	}

	// Parse the URL to determine if TLS is needed and get the address
	var creds credentials.TransportCredentials
	var address string

	if strings.HasPrefix(cfg.URL, "grpcs://") {
		// Secure gRPC with TLS
		address = strings.TrimPrefix(cfg.URL, "grpcs://")
		creds = credentials.NewTLS(&tls.Config{})
	} else if strings.HasPrefix(cfg.URL, "grpc://") {
		// Check if port 443 (typically TLS)
		address = strings.TrimPrefix(cfg.URL, "grpc://")
		if strings.Contains(address, ":443") {
			// Port 443 or .network domains typically require TLS
			creds = credentials.NewTLS(&tls.Config{})
		} else {
			creds = insecure.NewCredentials()
		}
	} else {
		// Default to the URL as-is, assume TLS for port 443
		address = cfg.URL
		if strings.Contains(address, ":443") {
			creds = credentials.NewTLS(&tls.Config{})
		} else {
			creds = insecure.NewCredentials()
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(), //grpc.ForceCodec(buildGRPCCodec())),
	)
	if err != nil {
		return nil, errors.Errorf("failed to connect to %s: %w", address, err)
	}

	client := &GRPCClient{
		endpointURL:  cfg.URL,
		conn:         conn,
		gov:          NewGovGRPCClient(conn, logger),
		evidence:     NewEvidenceGRPCClient(conn, logger),
		auth:         NewAuthGRPCClient(conn, logger),
		tx:           NewTxGRPCClient(conn, logger),
		slashing:     NewSlashingGRPCClient(conn, logger),
		emissions:    NewEmissionsGRPCClient(conn, logger),
		feegrant:     NewFeegrantGRPCClient(conn, logger),
		mint:         NewMintGRPCClient(conn, logger),
		staking:      NewStakingGRPCClient(conn, logger),
		authz:        NewAuthzGRPCClient(conn, logger),
		node:         NewNodeGRPCClient(conn, logger),
		consensus:    NewConsensusGRPCClient(conn, logger),
		distribution: NewDistributionGRPCClient(conn, logger),
		bank:         NewBankGRPCClient(conn, logger),
		params:       NewParamsGRPCClient(conn, logger),
		tendermint:   NewTendermintGRPCClient(conn, logger),
	}

	return client, nil
}

// Close closes the gRPC connection
func (c *GRPCClient) Close() error {
	return c.conn.Close()
}

func (c *GRPCClient) GetProtocol() config.Protocol {
	return config.ProtocolGRPC
}

func (c *GRPCClient) GetEndpointURL() string {
	return c.endpointURL
}
func (c *GRPCClient) Gov() interfaces.GovClient {
	return c.gov
}
func (c *GRPCClient) Evidence() interfaces.EvidenceClient {
	return c.evidence
}
func (c *GRPCClient) Auth() interfaces.AuthClient {
	return c.auth
}
func (c *GRPCClient) Tx() interfaces.TxClient {
	return c.tx
}
func (c *GRPCClient) Slashing() interfaces.SlashingClient {
	return c.slashing
}
func (c *GRPCClient) Emissions() interfaces.EmissionsClient {
	return c.emissions
}
func (c *GRPCClient) Feegrant() interfaces.FeegrantClient {
	return c.feegrant
}
func (c *GRPCClient) Mint() interfaces.MintClient {
	return c.mint
}
func (c *GRPCClient) Staking() interfaces.StakingClient {
	return c.staking
}
func (c *GRPCClient) Authz() interfaces.AuthzClient {
	return c.authz
}
func (c *GRPCClient) Node() interfaces.NodeClient {
	return c.node
}
func (c *GRPCClient) Consensus() interfaces.ConsensusClient {
	return c.consensus
}
func (c *GRPCClient) Distribution() interfaces.DistributionClient {
	return c.distribution
}
func (c *GRPCClient) Bank() interfaces.BankClient {
	return c.bank
}
func (c *GRPCClient) Params() interfaces.ParamsClient {
	return c.params
}
func (c *GRPCClient) Tendermint() interfaces.TendermintClient {
	return c.tendermint
}

type callFn[In, Out any] func(ctx context.Context, i In, opts ...grpc.CallOption) (Out, error)

func queryWithHeight[In any, Out any](ctx context.Context, height int64, queryFn callFn[In, Out], in In) (Out, error) {
	var zero Out

	if height > 0 {
		ctx = metadata.AppendToOutgoingContext(ctx, sdkgrpc.GRPCBlockHeightHeader, strconv.FormatInt(height, 10))
	}

	resp, err := queryFn(ctx, in)
	if err != nil {
		return zero, err
	}
	return resp, nil
}

// Status implements the Status method required by ClientPoolManager
func (c *GRPCClient) Status(ctx context.Context) error {
	_, err := c.tendermint.GetSyncing(ctx, &cmtservice.GetSyncingRequest{})
	return err
}

// HealthCheck wraps Status to satisfy pool requirements
func (c *GRPCClient) HealthCheck(ctx context.Context) error {
	return c.Status(ctx)
}
