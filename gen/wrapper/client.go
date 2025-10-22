package wrapper

import (
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/pool"
)

type WrapperClient struct {
	bank         *BankClientWrapper
	tendermint   *TendermintClientWrapper
	params       *ParamsClientWrapper
	emissions    *EmissionsClientWrapper
	distribution *DistributionClientWrapper
	consensus    *ConsensusClientWrapper
	staking      *StakingClientWrapper
	auth         *AuthClientWrapper
	mint         *MintClientWrapper
	tx           *TxClientWrapper
	evidence     *EvidenceClientWrapper
	feegrant     *FeegrantClientWrapper
	slashing     *SlashingClientWrapper
	gov          *GovClientWrapper
	authz        *AuthzClientWrapper
	node         *NodeClientWrapper
}

func NewWrapperClient(poolManager *pool.ClientPoolManager[interfaces.CosmosClient], logger zerolog.Logger) *WrapperClient {
	return &WrapperClient{
		bank:         NewBankClientWrapper(poolManager, logger),
		tendermint:   NewTendermintClientWrapper(poolManager, logger),
		params:       NewParamsClientWrapper(poolManager, logger),
		emissions:    NewEmissionsClientWrapper(poolManager, logger),
		distribution: NewDistributionClientWrapper(poolManager, logger),
		consensus:    NewConsensusClientWrapper(poolManager, logger),
		staking:      NewStakingClientWrapper(poolManager, logger),
		auth:         NewAuthClientWrapper(poolManager, logger),
		mint:         NewMintClientWrapper(poolManager, logger),
		tx:           NewTxClientWrapper(poolManager, logger),
		evidence:     NewEvidenceClientWrapper(poolManager, logger),
		feegrant:     NewFeegrantClientWrapper(poolManager, logger),
		slashing:     NewSlashingClientWrapper(poolManager, logger),
		gov:          NewGovClientWrapper(poolManager, logger),
		authz:        NewAuthzClientWrapper(poolManager, logger),
		node:         NewNodeClientWrapper(poolManager, logger),
	}
}

func (c *WrapperClient) Close() error {
	return nil
}
func (c *WrapperClient) Bank() interfaces.BankClient {
	return c.bank
}

func (c *WrapperClient) Tendermint() interfaces.TendermintClient {
	return c.tendermint
}

func (c *WrapperClient) Params() interfaces.ParamsClient {
	return c.params
}

func (c *WrapperClient) Emissions() interfaces.EmissionsClient {
	return c.emissions
}

func (c *WrapperClient) Distribution() interfaces.DistributionClient {
	return c.distribution
}

func (c *WrapperClient) Consensus() interfaces.ConsensusClient {
	return c.consensus
}

func (c *WrapperClient) Staking() interfaces.StakingClient {
	return c.staking
}

func (c *WrapperClient) Auth() interfaces.AuthClient {
	return c.auth
}

func (c *WrapperClient) Mint() interfaces.MintClient {
	return c.mint
}

func (c *WrapperClient) Tx() interfaces.TxClient {
	return c.tx
}

func (c *WrapperClient) Evidence() interfaces.EvidenceClient {
	return c.evidence
}

func (c *WrapperClient) Feegrant() interfaces.FeegrantClient {
	return c.feegrant
}

func (c *WrapperClient) Slashing() interfaces.SlashingClient {
	return c.slashing
}

func (c *WrapperClient) Gov() interfaces.GovClient {
	return c.gov
}

func (c *WrapperClient) Authz() interfaces.AuthzClient {
	return c.authz
}

func (c *WrapperClient) Node() interfaces.NodeClient {
	return c.node
}
