package wrapper

import (
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/pool"
)

type WrapperClient struct {
	tendermint   *TendermintClientWrapper
	mint         *MintClientWrapper
	slashing     *SlashingClientWrapper
	consensus    *ConsensusClientWrapper
	distribution *DistributionClientWrapper
	node         *NodeClientWrapper
	emissions    *EmissionsClientWrapper
	staking      *StakingClientWrapper
	evidence     *EvidenceClientWrapper
	authz        *AuthzClientWrapper
	auth         *AuthClientWrapper
	bank         *BankClientWrapper
	feegrant     *FeegrantClientWrapper
	params       *ParamsClientWrapper
	tx           *TxClientWrapper
	gov          *GovClientWrapper
}

func NewWrapperClient(poolManager *pool.ClientPoolManager, logger zerolog.Logger) *WrapperClient {
	return &WrapperClient{
		tendermint:   NewTendermintClientWrapper(poolManager, logger),
		mint:         NewMintClientWrapper(poolManager, logger),
		slashing:     NewSlashingClientWrapper(poolManager, logger),
		consensus:    NewConsensusClientWrapper(poolManager, logger),
		distribution: NewDistributionClientWrapper(poolManager, logger),
		node:         NewNodeClientWrapper(poolManager, logger),
		emissions:    NewEmissionsClientWrapper(poolManager, logger),
		staking:      NewStakingClientWrapper(poolManager, logger),
		evidence:     NewEvidenceClientWrapper(poolManager, logger),
		authz:        NewAuthzClientWrapper(poolManager, logger),
		auth:         NewAuthClientWrapper(poolManager, logger),
		bank:         NewBankClientWrapper(poolManager, logger),
		feegrant:     NewFeegrantClientWrapper(poolManager, logger),
		params:       NewParamsClientWrapper(poolManager, logger),
		tx:           NewTxClientWrapper(poolManager, logger),
		gov:          NewGovClientWrapper(poolManager, logger),
	}
}

func (c *WrapperClient) Close() error {
	return nil
}
func (c *WrapperClient) Tendermint() interfaces.TendermintClient {
	return c.tendermint
}

func (c *WrapperClient) Mint() interfaces.MintClient {
	return c.mint
}

func (c *WrapperClient) Slashing() interfaces.SlashingClient {
	return c.slashing
}

func (c *WrapperClient) Consensus() interfaces.ConsensusClient {
	return c.consensus
}

func (c *WrapperClient) Distribution() interfaces.DistributionClient {
	return c.distribution
}

func (c *WrapperClient) Node() interfaces.NodeClient {
	return c.node
}

func (c *WrapperClient) Emissions() interfaces.EmissionsClient {
	return c.emissions
}

func (c *WrapperClient) Staking() interfaces.StakingClient {
	return c.staking
}

func (c *WrapperClient) Evidence() interfaces.EvidenceClient {
	return c.evidence
}

func (c *WrapperClient) Authz() interfaces.AuthzClient {
	return c.authz
}

func (c *WrapperClient) Auth() interfaces.AuthClient {
	return c.auth
}

func (c *WrapperClient) Bank() interfaces.BankClient {
	return c.bank
}

func (c *WrapperClient) Feegrant() interfaces.FeegrantClient {
	return c.feegrant
}

func (c *WrapperClient) Params() interfaces.ParamsClient {
	return c.params
}

func (c *WrapperClient) Tx() interfaces.TxClient {
	return c.tx
}

func (c *WrapperClient) Gov() interfaces.GovClient {
	return c.gov
}
