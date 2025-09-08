package interfaces

// Client defines the interface for protocol clients that aggregate multiple blockchain modules.
type Client interface {
	Close() error

	Auth() AuthClient
	Mint() MintClient
	Evidence() EvidenceClient
	Staking() StakingClient
	Distribution() DistributionClient
	Emissions() EmissionsClient
	Params() ParamsClient
	Feegrant() FeegrantClient
	Tx() TxClient
	Bank() BankClient
	Slashing() SlashingClient
	Node() NodeClient
	Authz() AuthzClient
	Consensus() ConsensusClient
	Tendermint() TendermintClient
	Gov() GovClient
}
