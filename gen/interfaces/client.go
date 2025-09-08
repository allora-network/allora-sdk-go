package interfaces

import (
	"context"
)

// Client defines the interface for protocol clients that aggregate multiple blockchain modules.
type Client interface {
	// Connection management methods required by ClientPoolManager
	GetEndpointURL() string
	GetProtocol() string
	Status(ctx context.Context) error
	Close() error

	Tendermint() TendermintClient
	Mint() MintClient
	Slashing() SlashingClient
	Consensus() ConsensusClient
	Distribution() DistributionClient
	Node() NodeClient
	Emissions() EmissionsClient
	Staking() StakingClient
	Evidence() EvidenceClient
	Authz() AuthzClient
	Auth() AuthClient
	Bank() BankClient
	Feegrant() FeegrantClient
	Params() ParamsClient
	Tx() TxClient
	Gov() GovClient
}
