package codec_test

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"

	cosmossdktypes "github.com/cosmos/cosmos-sdk/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"
	"github.com/stretchr/testify/require"

	"github.com/allora-network/allora-sdk-go/codec"
)

// loadTxFixture reads a testdata file holding one base64-encoded transaction,
// exactly as it appears in a CometBFT /block response.
func loadTxFixture(t *testing.T, name string) []byte {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)
	txBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(raw)))
	require.NoError(t, err)
	return txBytes
}

// Both fixtures are mainnet governance proposals whose single proposal message
// belongs to a module the chain runs outside the SDK's emissions/mint/bank set:
// x/consensus (block and evidence params) and skip-mev/feemarket (EIP-1559
// fee params). Decoding a transaction unpacks every nested Any, so a module
// missing from the registry fails the whole transaction, not just the message.
func TestCodecParsesGovernanceProposalsForChainModules(t *testing.T) {
	cdc := codec.NewCodec()

	cases := []struct {
		fixture string
		typeURL string
		check   func(t *testing.T, msg cosmossdktypes.Msg)
	}{
		{
			fixture: "mainnet_4487135_tx0_consensus_update_params.b64",
			typeURL: "/cosmos.consensus.v1.MsgUpdateParams",
			check: func(t *testing.T, msg cosmossdktypes.Msg) {
				m, ok := msg.(*consensustypes.MsgUpdateParams)
				require.True(t, ok, "got %T", msg)
				require.NotEmpty(t, m.Authority)
			},
		},
		{
			fixture: "mainnet_4487571_tx1_feemarket_params.b64",
			typeURL: "/feemarket.feemarket.v1.MsgParams",
			check: func(t *testing.T, msg cosmossdktypes.Msg) {
				m, ok := msg.(*feemarkettypes.MsgParams)
				require.True(t, ok, "got %T", msg)
				require.NotEmpty(t, m.Authority)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.typeURL, func(t *testing.T) {
			tx, err := cdc.ParseTx(loadTxFixture(t, tc.fixture))
			require.NoError(t, err)
			require.Len(t, tx.Body.Messages, 1)
			require.Equal(t, "/cosmos.gov.v1.MsgSubmitProposal", tx.Body.Messages[0].TypeUrl)

			proposal, ok := tx.Body.Messages[0].GetCachedValue().(*govv1types.MsgSubmitProposal)
			require.True(t, ok, "got %T", tx.Body.Messages[0].GetCachedValue())
			require.Len(t, proposal.Messages, 1)
			require.Equal(t, tc.typeURL, proposal.Messages[0].TypeUrl)

			inner, err := cdc.ParseTxMessage(proposal.Messages[0])
			require.NoError(t, err)
			tc.check(t, inner)

			// The transaction must also render to JSON, which resolves every
			// nested Any through the same registry.
			jsonBz, err := cdc.MarshalJSON(tx)
			require.NoError(t, err)
			require.Contains(t, string(jsonBz), tc.typeURL)
		})
	}
}
