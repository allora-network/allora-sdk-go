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

	"github.com/allora-network/allora-sdk-go/codec"
	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// GRPCClient implements the Client interface using gRPC
type GRPCClient struct {
	endpointURL  string
	conn         *grpc.ClientConn
	consensus    *ConsensusGRPCClient
	gov          *GovGRPCClient
	params       *ParamsGRPCClient
	mint         *MintGRPCClient
	authz        *AuthzGRPCClient
	distribution *DistributionGRPCClient
	tendermint   *TendermintGRPCClient
	auth         *AuthGRPCClient
	slashing     *SlashingGRPCClient
	bank         *BankGRPCClient
	emissions    *EmissionsGRPCClient
	tx           *TxGRPCClient
	evidence     *EvidenceGRPCClient
	node         *NodeGRPCClient
	feegrant     *FeegrantGRPCClient
	staking      *StakingGRPCClient
}

var _ interfaces.CosmosClient = (*GRPCClient)(nil)

var (
	grpcCodecOnce sync.Once
	grpcCodec     encoding.Codec
)

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
		grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.GRPCCodec())),
	)
	if err != nil {
		return nil, errors.Errorf("failed to connect to %s: %w", address, err)
	}

	client := &GRPCClient{
		endpointURL:  cfg.URL,
		conn:         conn,
		consensus:    NewConsensusGRPCClient(conn, logger),
		gov:          NewGovGRPCClient(conn, logger),
		params:       NewParamsGRPCClient(conn, logger),
		mint:         NewMintGRPCClient(conn, logger),
		authz:        NewAuthzGRPCClient(conn, logger),
		distribution: NewDistributionGRPCClient(conn, logger),
		tendermint:   NewTendermintGRPCClient(conn, logger),
		auth:         NewAuthGRPCClient(conn, logger),
		slashing:     NewSlashingGRPCClient(conn, logger),
		bank:         NewBankGRPCClient(conn, logger),
		emissions:    NewEmissionsGRPCClient(conn, logger),
		tx:           NewTxGRPCClient(conn, logger),
		evidence:     NewEvidenceGRPCClient(conn, logger),
		node:         NewNodeGRPCClient(conn, logger),
		feegrant:     NewFeegrantGRPCClient(conn, logger),
		staking:      NewStakingGRPCClient(conn, logger),
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
func (c *GRPCClient) Consensus() interfaces.ConsensusClient {
	return c.consensus
}
func (c *GRPCClient) Gov() interfaces.GovClient {
	return c.gov
}
func (c *GRPCClient) Params() interfaces.ParamsClient {
	return c.params
}
func (c *GRPCClient) Mint() interfaces.MintClient {
	return c.mint
}
func (c *GRPCClient) Authz() interfaces.AuthzClient {
	return c.authz
}
func (c *GRPCClient) Distribution() interfaces.DistributionClient {
	return c.distribution
}
func (c *GRPCClient) Tendermint() interfaces.TendermintClient {
	return c.tendermint
}
func (c *GRPCClient) Auth() interfaces.AuthClient {
	return c.auth
}
func (c *GRPCClient) Slashing() interfaces.SlashingClient {
	return c.slashing
}
func (c *GRPCClient) Bank() interfaces.BankClient {
	return c.bank
}
func (c *GRPCClient) Emissions() interfaces.EmissionsClient {
	return c.emissions
}
func (c *GRPCClient) Tx() interfaces.TxClient {
	return c.tx
}
func (c *GRPCClient) Evidence() interfaces.EvidenceClient {
	return c.evidence
}
func (c *GRPCClient) Node() interfaces.NodeClient {
	return c.node
}
func (c *GRPCClient) Feegrant() interfaces.FeegrantClient {
	return c.feegrant
}
func (c *GRPCClient) Staking() interfaces.StakingClient {
	return c.staking
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
