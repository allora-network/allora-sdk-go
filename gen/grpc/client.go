package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// GRPCClient implements the Client interface using gRPC
type GRPCClient struct {
	endpointURL  string
	conn         *grpc.ClientConn
	tendermint   *TendermintGRPCClient
	mint         *MintGRPCClient
	slashing     *SlashingGRPCClient
	consensus    *ConsensusGRPCClient
	distribution *DistributionGRPCClient
	node         *NodeGRPCClient
	emissions    *EmissionsGRPCClient
	staking      *StakingGRPCClient
	evidence     *EvidenceGRPCClient
	authz        *AuthzGRPCClient
	auth         *AuthGRPCClient
	bank         *BankGRPCClient
	feegrant     *FeegrantGRPCClient
	params       *ParamsGRPCClient
	tx           *TxGRPCClient
	gov          *GovGRPCClient
}

var _ interfaces.Client = (*GRPCClient)(nil)

// NewGRPCClient creates a new gRPC aggregated client
func NewGRPCClient(cfg config.EndpointConfig, logger zerolog.Logger) (*GRPCClient, error) {
	if cfg.Protocol != "grpc" {
		return nil, fmt.Errorf("unsupported protocol: %s, only 'grpc' is supported", cfg.Protocol)
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
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", address, err)
	}

	client := &GRPCClient{
		endpointURL:  cfg.URL,
		conn:         conn,
		tendermint:   NewTendermintGRPCClient(conn, logger),
		mint:         NewMintGRPCClient(conn, logger),
		slashing:     NewSlashingGRPCClient(conn, logger),
		consensus:    NewConsensusGRPCClient(conn, logger),
		distribution: NewDistributionGRPCClient(conn, logger),
		node:         NewNodeGRPCClient(conn, logger),
		emissions:    NewEmissionsGRPCClient(conn, logger),
		staking:      NewStakingGRPCClient(conn, logger),
		evidence:     NewEvidenceGRPCClient(conn, logger),
		authz:        NewAuthzGRPCClient(conn, logger),
		auth:         NewAuthGRPCClient(conn, logger),
		bank:         NewBankGRPCClient(conn, logger),
		feegrant:     NewFeegrantGRPCClient(conn, logger),
		params:       NewParamsGRPCClient(conn, logger),
		tx:           NewTxGRPCClient(conn, logger),
		gov:          NewGovGRPCClient(conn, logger),
	}

	return client, nil
}

// Close closes the gRPC connection
func (c *GRPCClient) Close() error {
	return c.conn.Close()
}
func (c *GRPCClient) Tendermint() interfaces.TendermintClient {
	return c.tendermint
}
func (c *GRPCClient) Mint() interfaces.MintClient {
	return c.mint
}
func (c *GRPCClient) Slashing() interfaces.SlashingClient {
	return c.slashing
}
func (c *GRPCClient) Consensus() interfaces.ConsensusClient {
	return c.consensus
}
func (c *GRPCClient) Distribution() interfaces.DistributionClient {
	return c.distribution
}
func (c *GRPCClient) Node() interfaces.NodeClient {
	return c.node
}
func (c *GRPCClient) Emissions() interfaces.EmissionsClient {
	return c.emissions
}
func (c *GRPCClient) Staking() interfaces.StakingClient {
	return c.staking
}
func (c *GRPCClient) Evidence() interfaces.EvidenceClient {
	return c.evidence
}
func (c *GRPCClient) Authz() interfaces.AuthzClient {
	return c.authz
}
func (c *GRPCClient) Auth() interfaces.AuthClient {
	return c.auth
}
func (c *GRPCClient) Bank() interfaces.BankClient {
	return c.bank
}
func (c *GRPCClient) Feegrant() interfaces.FeegrantClient {
	return c.feegrant
}
func (c *GRPCClient) Params() interfaces.ParamsClient {
	return c.params
}
func (c *GRPCClient) Tx() interfaces.TxClient {
	return c.tx
}
func (c *GRPCClient) Gov() interfaces.GovClient {
	return c.gov
}

func (c *GRPCClient) GetEndpointURL() string {
	return c.endpointURL
}

func (c *GRPCClient) GetProtocol() string {
	return "grpc"
}

// Status implements the Status method required by ClientPoolManager
func (c *GRPCClient) Status(ctx context.Context) error {
	_, err := c.tendermint.GetSyncing(ctx, &cmtservice.GetSyncingRequest{})
	return err
}
