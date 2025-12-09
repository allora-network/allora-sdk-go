package interfaces

import (
	"github.com/allora-network/allora-sdk-go/pool"
)

type CosmosClientPool interface {
	Close() error

	Consensus() ConsensusClient
	Gov() GovClient
	Params() ParamsClient
	Mint() MintClient
	Authz() AuthzClient
	Distribution() DistributionClient
	Tendermint() TendermintClient
	Auth() AuthClient
	Slashing() SlashingClient
	Bank() BankClient
	Emissions() EmissionsClient
	Tx() TxClient
	Evidence() EvidenceClient
	Node() NodeClient
	Feegrant() FeegrantClient
	Staking() StakingClient
}

type CosmosClient interface {
	CosmosClientPool
	pool.PoolParticipant
}
