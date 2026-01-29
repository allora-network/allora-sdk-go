package interfaces

import (
	"github.com/allora-network/allora-sdk-go/pool"
)

type CosmosClientPool interface {
	Close() error

	Evidence() EvidenceClient
	Feegrant() FeegrantClient
	Emissions() EmissionsClient
	Mint() MintClient
	Tendermint() TendermintClient
	Node() NodeClient
	Tx() TxClient
	Auth() AuthClient
	Authz() AuthzClient
	Bank() BankClient
	Consensus() ConsensusClient
	Distribution() DistributionClient
	Gov() GovClient
	Params() ParamsClient
	Slashing() SlashingClient
	Staking() StakingClient
}

type CosmosClient interface {
	CosmosClientPool
	pool.PoolParticipant
}
