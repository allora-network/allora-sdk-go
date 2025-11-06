package interfaces

import (
	"github.com/allora-network/allora-sdk-go/pool"
)

type CosmosClientPool interface {
	Close() error

	Distribution() DistributionClient
	Gov() GovClient
	Evidence() EvidenceClient
	Tendermint() TendermintClient
	Consensus() ConsensusClient
	Tx() TxClient
	Authz() AuthzClient
	Node() NodeClient
	Mint() MintClient
	Emissions() EmissionsClient
	Staking() StakingClient
	Feegrant() FeegrantClient
	Slashing() SlashingClient
	Params() ParamsClient
	Bank() BankClient
	Auth() AuthClient
}

type CosmosClient interface {
	CosmosClientPool
	pool.PoolParticipant
}
