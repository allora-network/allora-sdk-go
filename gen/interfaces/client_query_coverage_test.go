package interfaces

import (
	"reflect"
	"slices"
	"testing"

	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	emissionstypes "github.com/allora-network/allora-chain/x/emissions/types"
	minttypes "github.com/allora-network/allora-chain/x/mint/types"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	"github.com/cosmos/cosmos-sdk/client/grpc/node"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	proposaltypes "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func TestGeneratedInterfacesCoverPinnedQueryClients(t *testing.T) {
	tests := []struct {
		generated any
		upstream  any
	}{
		{(*AuthClient)(nil), (*authtypes.QueryClient)(nil)},
		{(*AuthzClient)(nil), (*authz.QueryClient)(nil)},
		{(*BankClient)(nil), (*banktypes.QueryClient)(nil)},
		{(*ConsensusClient)(nil), (*consensustypes.QueryClient)(nil)},
		{(*DistributionClient)(nil), (*distributiontypes.QueryClient)(nil)},
		{(*EmissionsClient)(nil), (*emissionstypes.QueryServiceClient)(nil)},
		{(*EvidenceClient)(nil), (*evidencetypes.QueryClient)(nil)},
		{(*FeegrantClient)(nil), (*feegrant.QueryClient)(nil)},
		{(*GovClient)(nil), (*govv1.QueryClient)(nil)},
		{(*MintClient)(nil), (*minttypes.QueryServiceClient)(nil)},
		{(*NodeClient)(nil), (*node.ServiceClient)(nil)},
		{(*ParamsClient)(nil), (*proposaltypes.QueryClient)(nil)},
		{(*SlashingClient)(nil), (*slashingtypes.QueryClient)(nil)},
		{(*StakingClient)(nil), (*stakingtypes.QueryClient)(nil)},
		{(*TendermintClient)(nil), (*cmtservice.ServiceClient)(nil)},
		{(*TxClient)(nil), (*txtypes.ServiceClient)(nil)},
	}

	for _, tt := range tests {
		generatedType := reflect.TypeOf(tt.generated).Elem()
		upstreamType := reflect.TypeOf(tt.upstream).Elem()
		t.Run(generatedType.Name(), func(t *testing.T) {
			generated := methodNames(generatedType)
			upstream := methodNames(upstreamType)
			missing := difference(upstream, generated)
			extra := difference(generated, upstream)
			if len(missing) > 0 || len(extra) > 0 {
				t.Errorf("generated method set differs from %s\nmissing: %v\nextra: %v", upstreamType, missing, extra)
			}
		})
	}
}

// methodNames returns the sorted names of an interface's methods.
func methodNames(interfaceType reflect.Type) []string {
	if interfaceType.Kind() != reflect.Interface {
		panic("methodNames requires an interface")
	}
	names := make([]string, interfaceType.NumMethod())
	for i := range names {
		names[i] = interfaceType.Method(i).Name
	}
	slices.Sort(names)
	return names
}

// difference returns names in a that are not in b.
func difference(a, b []string) []string {
	var diff []string
	for _, name := range a {
		if !slices.Contains(b, name) {
			diff = append(diff, name)
		}
	}
	return diff
}
