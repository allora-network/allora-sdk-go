package wrapper

import (
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/pool"
)

type WrapperClient struct {
	auth         *AuthClientWrapper
	mint         *MintClientWrapper
	evidence     *EvidenceClientWrapper
	staking      *StakingClientWrapper
	distribution *DistributionClientWrapper
	emissions    *EmissionsClientWrapper
	params       *ParamsClientWrapper
	feegrant     *FeegrantClientWrapper
	tx           *TxClientWrapper
	bank         *BankClientWrapper
	slashing     *SlashingClientWrapper
	node         *NodeClientWrapper
	authz        *AuthzClientWrapper
	consensus    *ConsensusClientWrapper
	tendermint   *TendermintClientWrapper
	gov          *GovClientWrapper
}

func NewWrapperClient(poolManager *pool.ClientPoolManager, logger zerolog.Logger) *WrapperClient {
	return &WrapperClient{
		auth:         NewAuthClientWrapper(poolManager, logger),
		mint:         NewMintClientWrapper(poolManager, logger),
		evidence:     NewEvidenceClientWrapper(poolManager, logger),
		staking:      NewStakingClientWrapper(poolManager, logger),
		distribution: NewDistributionClientWrapper(poolManager, logger),
		emissions:    NewEmissionsClientWrapper(poolManager, logger),
		params:       NewParamsClientWrapper(poolManager, logger),
		feegrant:     NewFeegrantClientWrapper(poolManager, logger),
		tx:           NewTxClientWrapper(poolManager, logger),
		bank:         NewBankClientWrapper(poolManager, logger),
		slashing:     NewSlashingClientWrapper(poolManager, logger),
		node:         NewNodeClientWrapper(poolManager, logger),
		authz:        NewAuthzClientWrapper(poolManager, logger),
		consensus:    NewConsensusClientWrapper(poolManager, logger),
		tendermint:   NewTendermintClientWrapper(poolManager, logger),
		gov:          NewGovClientWrapper(poolManager, logger),
	}
}

func (c *WrapperClient) Close() error {
	return nil
}
func (c *WrapperClient) Auth() interfaces.AuthClient {
	return c.auth
}

func (c *WrapperClient) Mint() interfaces.MintClient {
	return c.mint
}

func (c *WrapperClient) Evidence() interfaces.EvidenceClient {
	return c.evidence
}

func (c *WrapperClient) Staking() interfaces.StakingClient {
	return c.staking
}

func (c *WrapperClient) Distribution() interfaces.DistributionClient {
	return c.distribution
}

func (c *WrapperClient) Emissions() interfaces.EmissionsClient {
	return c.emissions
}

func (c *WrapperClient) Params() interfaces.ParamsClient {
	return c.params
}

func (c *WrapperClient) Feegrant() interfaces.FeegrantClient {
	return c.feegrant
}

func (c *WrapperClient) Tx() interfaces.TxClient {
	return c.tx
}

func (c *WrapperClient) Bank() interfaces.BankClient {
	return c.bank
}

func (c *WrapperClient) Slashing() interfaces.SlashingClient {
	return c.slashing
}

func (c *WrapperClient) Node() interfaces.NodeClient {
	return c.node
}

func (c *WrapperClient) Authz() interfaces.AuthzClient {
	return c.authz
}

func (c *WrapperClient) Consensus() interfaces.ConsensusClient {
	return c.consensus
}

func (c *WrapperClient) Tendermint() interfaces.TendermintClient {
	return c.tendermint
}

func (c *WrapperClient) Gov() interfaces.GovClient {
	return c.gov
}
