package txmsg_test

import (
	"testing"
	"time"

	stderrors "errors"

	"cosmossdk.io/x/feegrant"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"

	emissionstypes "github.com/allora-network/allora-chain/x/emissions/types"

	"github.com/allora-network/allora-sdk-go/txmsg"
)

// generateTestAddr returns a fresh bech32-valid Allora address. The
// parent allora package's init() configures the SDK to use the "allo"
// account prefix, so addresses derived from random secp256k1 keys are
// guaranteed to pass the txmsg bech32 validation.
func generateTestAddr(t *testing.T) string {
	t.Helper()
	pk := secp256k1.GenPrivKey()
	return sdk.AccAddress(pk.PubKey().Address()).String()
}

// The test bodies in this file are table-driven and assert on the
// returned sdk.Msg concrete type plus its key fields. The valid-path
// cases use a freshly-generated bech32 address and require the
// returned message to be non-nil and to carry the expected fields; the
// error-path cases use a deliberately bad address or zero-valued field
// and require a wrapped error mentioning the offending field name.

// -------- bank.NewSend --------

func TestNewSend_Success(t *testing.T) {
	from := generateTestAddr(t)
	to := generateTestAddr(t)
	amount := sdk.NewCoins(sdk.NewInt64Coin("uallo", 1_000_000))

	msg, err := txmsg.NewSend(from, to, amount)
	require.NoError(t, err)
	require.NotNil(t, msg)

	// Concrete type must be *banktypes.MsgSend so callers can do an
	// unchecked type assertion when they need to.
	send, ok := msg.(*banktypes.MsgSend)
	require.True(t, ok, "expected *banktypes.MsgSend, got %T", msg)
	require.Equal(t, from, send.FromAddress)
	require.Equal(t, to, send.ToAddress)
	require.True(t, send.Amount.Equal(amount), "amount %s != expected %s", send.Amount, amount)
}

func TestNewSend_ValidationErrors(t *testing.T) {
	validAddr := generateTestAddr(t)
	otherAddr := generateTestAddr(t)
	validAmount := sdk.NewCoins(sdk.NewInt64Coin("uallo", 1))

	cases := []struct {
		name    string
		from    string
		to      string
		amount  sdk.Coins
		wantSub string // substring that must appear in the error
	}{
		{
			name:    "empty from",
			from:    "",
			to:      otherAddr,
			amount:  validAmount,
			wantSub: "from address is required",
		},
		{
			name:    "empty to",
			from:    validAddr,
			to:      "",
			amount:  validAmount,
			wantSub: "to address is required",
		},
		{
			name:    "malformed from",
			from:    "not-a-bech32-address",
			to:      otherAddr,
			amount:  validAmount,
			wantSub: "invalid from bech32 address",
		},
		{
			name:    "nil amount",
			from:    validAddr,
			to:      otherAddr,
			amount:  nil,
			wantSub: "send amount is required",
		},
		{
			name:    "zero amount",
			from:    validAddr,
			to:      otherAddr,
			amount:  sdk.NewCoins(sdk.NewInt64Coin("uallo", 0)),
			wantSub: "send amount is required",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg, err := txmsg.NewSend(tc.from, tc.to, tc.amount)
			require.Error(t, err)
			require.Nil(t, msg)
			require.Contains(t, err.Error(), tc.wantSub)
			// Errors must be wrap-able so callers can errors.Is/As
			// against the underlying SDK error. go-utils/errors
			// wraps pkg/errors, which preserves the chain, so
			// standard-library errors.Unwrap either returns a
			// non-nil cause or a nil value (for a leaf error). The
			// important guarantee is that the call does not panic
			// and that the returned error is comparable.
			_ = stderrors.Unwrap(err)
		})
	}
}

// -------- bank.NewMultiSend --------

func TestNewMultiSend_Success(t *testing.T) {
	a := generateTestAddr(t)
	b := generateTestAddr(t)
	c := generateTestAddr(t)

	inputs := []banktypes.Input{
		{Address: a, Coins: sdk.NewCoins(sdk.NewInt64Coin("uallo", 1_000))},
	}
	outputs := []banktypes.Output{
		{Address: b, Coins: sdk.NewCoins(sdk.NewInt64Coin("uallo", 600))},
		{Address: c, Coins: sdk.NewCoins(sdk.NewInt64Coin("uallo", 400))},
	}

	msg, err := txmsg.NewMultiSend(inputs, outputs)
	require.NoError(t, err)
	require.NotNil(t, msg)

	ms, ok := msg.(*banktypes.MsgMultiSend)
	require.True(t, ok, "expected *banktypes.MsgMultiSend, got %T", msg)
	require.Len(t, ms.Inputs, 1)
	require.Len(t, ms.Outputs, 2)
	require.Equal(t, a, ms.Inputs[0].Address)
}

func TestNewMultiSend_ValidationErrors(t *testing.T) {
	a := generateTestAddr(t)
	b := generateTestAddr(t)

	cases := []struct {
		name    string
		inputs  []banktypes.Input
		outputs []banktypes.Output
		wantSub string
	}{
		{
			name:    "no inputs",
			inputs:  nil,
			outputs: []banktypes.Output{{Address: a, Coins: sdk.NewCoins(sdk.NewInt64Coin("uallo", 1))}},
			wantSub: "at least one input",
		},
		{
			name:    "no outputs",
			inputs:  []banktypes.Input{{Address: a, Coins: sdk.NewCoins(sdk.NewInt64Coin("uallo", 1))}},
			outputs: nil,
			wantSub: "at least one output",
		},
		{
			name: "input bech32 malformed",
			inputs: []banktypes.Input{
				{Address: "not-bech32", Coins: sdk.NewCoins(sdk.NewInt64Coin("uallo", 100))},
			},
			outputs: []banktypes.Output{
				{Address: a, Coins: sdk.NewCoins(sdk.NewInt64Coin("uallo", 100))},
			},
			wantSub: "input[0].address",
		},
		{
			name: "output coins zero",
			inputs: []banktypes.Input{
				{Address: a, Coins: sdk.NewCoins(sdk.NewInt64Coin("uallo", 100))},
			},
			outputs: []banktypes.Output{
				{Address: b, Coins: sdk.NewCoins(sdk.NewInt64Coin("uallo", 0))},
			},
			wantSub: "output[0].coins",
		},
		{
			name: "totals do not balance",
			inputs: []banktypes.Input{
				{Address: a, Coins: sdk.NewCoins(sdk.NewInt64Coin("uallo", 100))},
			},
			outputs: []banktypes.Output{
				{Address: b, Coins: sdk.NewCoins(sdk.NewInt64Coin("uallo", 99))},
			},
			wantSub: "do not equal outputs",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg, err := txmsg.NewMultiSend(tc.inputs, tc.outputs)
			require.Error(t, err)
			require.Nil(t, msg)
			require.Contains(t, err.Error(), tc.wantSub)
		})
	}
}

// -------- emissions.NewRegisterWorker --------

func TestNewRegisterWorker_Success(t *testing.T) {
	sender := generateTestAddr(t)
	owner := generateTestAddr(t)

	msg, err := txmsg.NewRegisterWorker(txmsg.RegisterWorkerParams{
		Sender:    sender,
		TopicId:   42,
		Owner:     owner,
		IsReputer: false,
	})
	require.NoError(t, err)
	require.NotNil(t, msg)

	rr, ok := msg.(*emissionstypes.RegisterRequest)
	require.True(t, ok, "expected *emissionstypes.RegisterRequest, got %T", msg)
	require.Equal(t, sender, rr.Sender)
	require.Equal(t, uint64(42), rr.TopicId)
	require.Equal(t, owner, rr.Owner)
	require.False(t, rr.IsReputer)

	// IsReputer=true path: same constructor, different flag.
	msg2, err := txmsg.NewRegisterWorker(txmsg.RegisterWorkerParams{
		Sender:    sender,
		TopicId:   7,
		Owner:     owner,
		IsReputer: true,
	})
	require.NoError(t, err)
	rr2, ok := msg2.(*emissionstypes.RegisterRequest)
	require.True(t, ok)
	require.True(t, rr2.IsReputer)
}

func TestNewRegisterWorker_ValidationErrors(t *testing.T) {
	valid := generateTestAddr(t)

	cases := []struct {
		name    string
		params  txmsg.RegisterWorkerParams
		wantSub string
	}{
		{
			name: "empty sender",
			params: txmsg.RegisterWorkerParams{
				Sender:  "",
				TopicId: 1,
				Owner:   valid,
			},
			wantSub: "sender address is required",
		},
		{
			name: "malformed sender",
			params: txmsg.RegisterWorkerParams{
				Sender:  "garbage",
				TopicId: 1,
				Owner:   valid,
			},
			wantSub: "invalid sender bech32 address",
		},
		{
			name: "zero topic id",
			params: txmsg.RegisterWorkerParams{
				Sender:  valid,
				TopicId: 0,
				Owner:   valid,
			},
			wantSub: "topic id must be greater than 0",
		},
		{
			name: "empty owner",
			params: txmsg.RegisterWorkerParams{
				Sender:  valid,
				TopicId: 1,
				Owner:   "",
			},
			wantSub: "owner address is required",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg, err := txmsg.NewRegisterWorker(tc.params)
			require.Error(t, err)
			require.Nil(t, msg)
			require.Contains(t, err.Error(), tc.wantSub)
		})
	}
}

// -------- emissions.NewInsertWorkerPayload --------

func TestNewInsertWorkerPayload_Success(t *testing.T) {
	sender := generateTestAddr(t)
	worker := generateTestAddr(t)

	// The chain requires a non-empty signature and a non-nil nonce; the
	// constructor mirrors both. We don't have to provide a real
	// signature value — any non-empty byte slice is enough to pass the
	// constructor's pre-flight guards (the chain will reject the
	// signature cryptographically, but that's a broadcast concern).
	bundle := &emissionstypes.InputWorkerDataBundle{
		Worker:                             worker,
		Nonce:                              &emissionstypes.Nonce{BlockHeight: 1234},
		TopicId:                            1,
		InferencesForecastsBundleSignature: []byte{0xde, 0xad, 0xbe, 0xef},
	}

	msg, err := txmsg.NewInsertWorkerPayload(txmsg.InsertWorkerPayloadParams{
		Sender:           sender,
		WorkerDataBundle: bundle,
	})
	require.NoError(t, err)
	require.NotNil(t, msg)

	ip, ok := msg.(*emissionstypes.InsertWorkerPayloadRequest)
	require.True(t, ok, "expected *emissionstypes.InsertWorkerPayloadRequest, got %T", msg)
	require.Equal(t, sender, ip.Sender)
	require.NotNil(t, ip.WorkerDataBundle)
	require.Equal(t, worker, ip.WorkerDataBundle.Worker)
	require.Equal(t, int64(1234), ip.WorkerDataBundle.Nonce.BlockHeight)
}

func TestNewInsertWorkerPayload_ValidationErrors(t *testing.T) {
	sender := generateTestAddr(t)
	worker := generateTestAddr(t)
	goodBundle := &emissionstypes.InputWorkerDataBundle{
		Worker:                             worker,
		Nonce:                              &emissionstypes.Nonce{BlockHeight: 1},
		InferencesForecastsBundleSignature: []byte{0x01},
	}

	cases := []struct {
		name    string
		params  txmsg.InsertWorkerPayloadParams
		wantSub string
	}{
		{
			name: "empty sender",
			params: txmsg.InsertWorkerPayloadParams{
				Sender:           "",
				WorkerDataBundle: goodBundle,
			},
			wantSub: "sender address is required",
		},
		{
			name: "nil bundle",
			params: txmsg.InsertWorkerPayloadParams{
				Sender:           sender,
				WorkerDataBundle: nil,
			},
			wantSub: "worker data bundle is required",
		},
		{
			name: "nil nonce",
			params: txmsg.InsertWorkerPayloadParams{
				Sender: sender,
				WorkerDataBundle: &emissionstypes.InputWorkerDataBundle{
					Worker:                             worker,
					InferencesForecastsBundleSignature: []byte{0x01},
				},
			},
			wantSub: "worker data bundle nonce is required",
		},
		{
			name: "empty signature",
			params: txmsg.InsertWorkerPayloadParams{
				Sender: sender,
				WorkerDataBundle: &emissionstypes.InputWorkerDataBundle{
					Worker: worker,
					Nonce:  &emissionstypes.Nonce{BlockHeight: 1},
				},
			},
			wantSub: "worker data bundle signature is required",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg, err := txmsg.NewInsertWorkerPayload(tc.params)
			require.Error(t, err)
			require.Nil(t, msg)
			require.Contains(t, err.Error(), tc.wantSub)
		})
	}
}

// -------- feegrant.NewGrantAllowance --------

func TestNewGrantAllowance_BasicAllowance_Success(t *testing.T) {
	granter := generateTestAddr(t)
	grantee := generateTestAddr(t)
	expiry := time.Now().Add(24 * time.Hour)
	allowance := &feegrant.BasicAllowance{
		SpendLimit: sdk.NewCoins(sdk.NewInt64Coin("uallo", 1_000_000)),
		Expiration: &expiry,
	}

	msg, err := txmsg.NewGrantAllowance(granter, grantee, allowance)
	require.NoError(t, err)
	require.NotNil(t, msg)

	ga, ok := msg.(*feegrant.MsgGrantAllowance)
	require.True(t, ok, "expected *feegrant.MsgGrantAllowance, got %T", msg)
	require.Equal(t, granter, ga.Granter)
	require.Equal(t, grantee, ga.Grantee)
	require.NotNil(t, ga.Allowance)
	require.Equal(t, "/cosmos.feegrant.v1beta1.BasicAllowance", ga.Allowance.TypeUrl)

	// Round-trip: assert on the concrete fields so a future codec
	// swap that breaks the packed TypeUrl shows up here.
	ba, ok := ga.Allowance.GetCachedValue().(*feegrant.BasicAllowance)
	require.True(t, ok, "Allowance.GetCachedValue() returned %T", ga.Allowance.GetCachedValue())
	require.True(t, ba.SpendLimit.Equal(allowance.SpendLimit))
	require.NotNil(t, ba.Expiration)
}

func TestNewGrantAllowance_PeriodicAllowance_Success(t *testing.T) {
	granter := generateTestAddr(t)
	grantee := generateTestAddr(t)
	allowance := &feegrant.PeriodicAllowance{
		Basic: feegrant.BasicAllowance{
			SpendLimit: sdk.NewCoins(sdk.NewInt64Coin("uallo", 5_000_000)),
		},
		Period:           time.Hour,
		PeriodSpendLimit: sdk.NewCoins(sdk.NewInt64Coin("uallo", 500_000)),
		PeriodCanSpend:   sdk.NewCoins(sdk.NewInt64Coin("uallo", 500_000)),
	}

	msg, err := txmsg.NewGrantAllowance(granter, grantee, allowance)
	require.NoError(t, err)
	require.NotNil(t, msg)

	ga, ok := msg.(*feegrant.MsgGrantAllowance)
	require.True(t, ok)
	require.Equal(t, "/cosmos.feegrant.v1beta1.PeriodicAllowance", ga.Allowance.TypeUrl)
}

func TestNewGrantAllowance_ValidationErrors(t *testing.T) {
	granter := generateTestAddr(t)
	grantee := generateTestAddr(t)
	goodBasic := &feegrant.BasicAllowance{
		SpendLimit: sdk.NewCoins(sdk.NewInt64Coin("uallo", 100)),
	}
	zeroSpendBasic := &feegrant.BasicAllowance{
		SpendLimit: sdk.NewCoins(sdk.NewInt64Coin("uallo", 0)),
	}
	zeroPeriod := &feegrant.PeriodicAllowance{
		Basic:            feegrant.BasicAllowance{SpendLimit: sdk.NewCoins(sdk.NewInt64Coin("uallo", 100))},
		Period:           time.Hour,
		PeriodSpendLimit: sdk.NewCoins(sdk.NewInt64Coin("uallo", 0)),
		PeriodCanSpend:   sdk.NewCoins(sdk.NewInt64Coin("uallo", 0)),
	}

	cases := []struct {
		name      string
		granter   string
		grantee   string
		allowance proto.Message
		wantSub   string
	}{
		{
			name:      "empty granter",
			granter:   "",
			grantee:   grantee,
			allowance: goodBasic,
			wantSub:   "granter address is required",
		},
		{
			name:      "malformed grantee",
			granter:   granter,
			grantee:   "not-bech32",
			allowance: goodBasic,
			wantSub:   "invalid grantee bech32 address",
		},
		{
			name:      "nil allowance",
			granter:   granter,
			grantee:   grantee,
			allowance: nil,
			wantSub:   "fee allowance is required",
		},
		{
			name:      "zero basic spend limit",
			granter:   granter,
			grantee:   grantee,
			allowance: zeroSpendBasic,
			wantSub:   "basic allowance spend limit",
		},
		{
			name:      "zero periodic spend limit",
			granter:   granter,
			grantee:   grantee,
			allowance: zeroPeriod,
			wantSub:   "periodic allowance period spend limit",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg, err := txmsg.NewGrantAllowance(tc.granter, tc.grantee, tc.allowance)
			require.Error(t, err)
			require.Nil(t, msg)
			require.Contains(t, err.Error(), tc.wantSub)
		})
	}
}

// -------- feegrant.NewRevokeAllowance --------

func TestNewRevokeAllowance_Success(t *testing.T) {
	granter := generateTestAddr(t)
	grantee := generateTestAddr(t)

	msg, err := txmsg.NewRevokeAllowance(granter, grantee)
	require.NoError(t, err)
	require.NotNil(t, msg)

	ra, ok := msg.(*feegrant.MsgRevokeAllowance)
	require.True(t, ok, "expected *feegrant.MsgRevokeAllowance, got %T", msg)
	require.Equal(t, granter, ra.Granter)
	require.Equal(t, grantee, ra.Grantee)
}

func TestNewRevokeAllowance_ValidationErrors(t *testing.T) {
	valid := generateTestAddr(t)
	cases := []struct {
		name    string
		granter string
		grantee string
		wantSub string
	}{
		{
			name:    "empty granter",
			granter: "",
			grantee: valid,
			wantSub: "granter address is required",
		},
		{
			name:    "empty grantee",
			granter: valid,
			grantee: "",
			wantSub: "grantee address is required",
		},
		{
			name:    "malformed granter",
			granter: "garbage",
			grantee: valid,
			wantSub: "invalid granter bech32 address",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg, err := txmsg.NewRevokeAllowance(tc.granter, tc.grantee)
			require.Error(t, err)
			require.Nil(t, msg)
			require.Contains(t, err.Error(), tc.wantSub)
		})
	}
}
