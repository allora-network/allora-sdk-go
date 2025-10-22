package interfaces

import (
	"github.com/allora-network/allora-sdk-go/pool"
)

type CosmosClientPool interface {
	Close() error

	Gov() GovClient
	Evidence() EvidenceClient
	Auth() AuthClient
	Tx() TxClient
	Slashing() SlashingClient
	Emissions() EmissionsClient
	Feegrant() FeegrantClient
	Mint() MintClient
	Staking() StakingClient
	Authz() AuthzClient
	Node() NodeClient
	Consensus() ConsensusClient
	Distribution() DistributionClient
	Bank() BankClient
	Params() ParamsClient
	Tendermint() TendermintClient
}

type CosmosClient interface {
	CosmosClientPool
	pool.PoolParticipant
}
