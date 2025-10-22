package interfaces

import (
	"github.com/allora-network/allora-sdk-go/pool"
)

type CosmosClientPool interface {
	Close() error

	Bank() BankClient
	Tendermint() TendermintClient
	Params() ParamsClient
	Emissions() EmissionsClient
	Distribution() DistributionClient
	Consensus() ConsensusClient
	Staking() StakingClient
	Auth() AuthClient
	Mint() MintClient
	Tx() TxClient
	Evidence() EvidenceClient
	Feegrant() FeegrantClient
	Slashing() SlashingClient
	Gov() GovClient
	Authz() AuthzClient
	Node() NodeClient
}

type CosmosClient interface {
	CosmosClientPool
	pool.PoolParticipant
}
