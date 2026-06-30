// Package cosmospool tests for the AccountInfo path. The AuthClient and
// CosmosTxPool used here are hand-rolled fakes (not the generated mockery
// mocks): the only method the production code touches is pool.Auth().Account,
// so building a focused double is clearer than pulling in the full generated
// surfaces and lets this file double as a contract check for the slice of
// the cosmos client the broadcaster actually consumes for account discovery.
package cosmospool_test

import (
	"context"
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	tx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/txsend/cosmospool"
)

// -- fakes --

// fakeAuthClient implements interfaces.AuthClient with only Account() doing
// anything meaningful. All other methods panic — they're not part of this
// bead's contract, and a test that accidentally exercises one of them should
// fail loudly rather than silently do the wrong thing.
type fakeAuthClient struct {
	// accountResp / accountErr are what Account() returns; tests set these to
	// drive the success, not-found, and transport-error cases.
	accountResp *authtypes.QueryAccountResponse
	accountErr  error
	// lastReq captures the request the production code sent, for assertions.
	lastReq *authtypes.QueryAccountRequest
}

func (f *fakeAuthClient) Account(ctx context.Context, req *authtypes.QueryAccountRequest, _ ...config.CallOpt) (*authtypes.QueryAccountResponse, error) {
	f.lastReq = req
	return f.accountResp, f.accountErr
}

func (f *fakeAuthClient) AccountAddressByID(context.Context, *authtypes.QueryAccountAddressByIDRequest, ...config.CallOpt) (*authtypes.QueryAccountAddressByIDResponse, error) {
	panic("AccountAddressByID not expected in this test")
}
func (f *fakeAuthClient) AccountInfo(context.Context, *authtypes.QueryAccountInfoRequest, ...config.CallOpt) (*authtypes.QueryAccountInfoResponse, error) {
	panic("AccountInfo (auth) not expected in this test")
}
func (f *fakeAuthClient) Accounts(context.Context, *authtypes.QueryAccountsRequest, ...config.CallOpt) (*authtypes.QueryAccountsResponse, error) {
	panic("Accounts not expected in this test")
}
func (f *fakeAuthClient) AddressBytesToString(context.Context, *authtypes.AddressBytesToStringRequest, ...config.CallOpt) (*authtypes.AddressBytesToStringResponse, error) {
	panic("AddressBytesToString not expected in this test")
}
func (f *fakeAuthClient) AddressStringToBytes(context.Context, *authtypes.AddressStringToBytesRequest, ...config.CallOpt) (*authtypes.AddressStringToBytesResponse, error) {
	panic("AddressStringToBytes not expected in this test")
}
func (f *fakeAuthClient) Bech32Prefix(context.Context, *authtypes.Bech32PrefixRequest, ...config.CallOpt) (*authtypes.Bech32PrefixResponse, error) {
	panic("Bech32Prefix not expected in this test")
}
func (f *fakeAuthClient) ModuleAccountByName(context.Context, *authtypes.QueryModuleAccountByNameRequest, ...config.CallOpt) (*authtypes.QueryModuleAccountByNameResponse, error) {
	panic("ModuleAccountByName not expected in this test")
}
func (f *fakeAuthClient) ModuleAccounts(context.Context, *authtypes.QueryModuleAccountsRequest, ...config.CallOpt) (*authtypes.QueryModuleAccountsResponse, error) {
	panic("ModuleAccounts not expected in this test")
}
func (f *fakeAuthClient) Params(context.Context, *authtypes.QueryParamsRequest, ...config.CallOpt) (*authtypes.QueryParamsResponse, error) {
	panic("Params not expected in this test")
}

// fakeTxClient implements interfaces.TxClient with every method panicking.
// AccountInfo must not touch Tx(); a panic here would catch a regression where
// the broadcaster's call graph accidentally reached into the tx service.
type fakeTxClient struct{}

func (fakeTxClient) BroadcastTx(context.Context, *tx.BroadcastTxRequest, ...config.CallOpt) (*tx.BroadcastTxResponse, error) {
	panic("BroadcastTx not expected in this test")
}
func (fakeTxClient) GetBlockWithTxs(context.Context, *tx.GetBlockWithTxsRequest, ...config.CallOpt) (*tx.GetBlockWithTxsResponse, error) {
	panic("GetBlockWithTxs not expected in this test")
}
func (fakeTxClient) GetTx(context.Context, *tx.GetTxRequest, ...config.CallOpt) (*tx.GetTxResponse, error) {
	panic("GetTx not expected in this test")
}
func (fakeTxClient) GetTxsEvent(context.Context, *tx.GetTxsEventRequest, ...config.CallOpt) (*tx.GetTxsEventResponse, error) {
	panic("GetTxsEvent not expected in this test")
}
func (fakeTxClient) Simulate(context.Context, *tx.SimulateRequest, ...config.CallOpt) (*tx.SimulateResponse, error) {
	panic("Simulate not expected in this test")
}
func (fakeTxClient) TxDecode(context.Context, *tx.TxDecodeRequest, ...config.CallOpt) (*tx.TxDecodeResponse, error) {
	panic("TxDecode not expected in this test")
}
func (fakeTxClient) TxDecodeAmino(context.Context, *tx.TxDecodeAminoRequest, ...config.CallOpt) (*tx.TxDecodeAminoResponse, error) {
	panic("TxDecodeAmino not expected in this test")
}
func (fakeTxClient) TxEncode(context.Context, *tx.TxEncodeRequest, ...config.CallOpt) (*tx.TxEncodeResponse, error) {
	panic("TxEncode not expected in this test")
}
func (fakeTxClient) TxEncodeAmino(context.Context, *tx.TxEncodeAminoRequest, ...config.CallOpt) (*tx.TxEncodeAminoResponse, error) {
	panic("TxEncodeAmino not expected in this test")
}

// testPool is the minimum a test needs to drive AccountInfo: a configurable
// AuthClient and a TxClient that should never be called. It satisfies
// txsend.CosmosTxPool (defined in cosmospool.go's sibling package as
// "Tx() interfaces.TxClient; Auth() interfaces.AuthClient").
type testPool struct {
	auth *fakeAuthClient
	tx   *fakeTxClient
}

func (p *testPool) Auth() interfaces.AuthClient { return p.auth }
func (p *testPool) Tx() interfaces.TxClient     { return p.tx }

// packBaseAccount wraps a *authtypes.BaseAccount in a *codectypes.Any the way
// the cosmos auth module would on the wire: via codectypes.NewAnyWithValue,
// which produces a canonical TypeURL and a deterministic marshal. This
// exercises the same code path the production unpack will see at runtime.
func packBaseAccount(ba *authtypes.BaseAccount) (*codectypes.Any, error) {
	return codectypes.NewAnyWithValue(ba)
}

// errString is a tiny stdlib-error type for tests that need a plain (non-gRPC)
// error whose message can be substring-matched.
type errString string

func (e errString) Error() string { return string(e) }

// -- tests --

func TestAccountInfo(t *testing.T) {
	// Known number/sequence we expect to round-trip through the codec pack/unpack.
	const (
		wantNumber   uint64 = 42
		wantSequence uint64 = 7
		address             = "cosmos1abc"
	)

	tests := []struct {
		name       string
		setup      func(ac *fakeAuthClient)
		wantErr    bool
		wantSubstr string // substring expected in the returned error (empty = expect no error)
		wantNumber uint64
		wantSeq    uint64
	}{
		{
			name: "success: packed BaseAccount round-trips accountNumber and sequence",
			setup: func(ac *fakeAuthClient) {
				ba := &authtypes.BaseAccount{
					Address:       address,
					AccountNumber: wantNumber,
					Sequence:      wantSequence,
				}
				any, err := packBaseAccount(ba)
				require.NoError(t, err, "pack BaseAccount for test")
				ac.accountResp = &authtypes.QueryAccountResponse{Account: any}
			},
			wantErr:    false,
			wantNumber: wantNumber,
			wantSeq:    wantSequence,
		},
		{
			name: "not-found: gRPC codes.NotFound → wrapped error mentions not found",
			setup: func(ac *fakeAuthClient) {
				ac.accountErr = status.Error(codes.NotFound, "account "+address+" not found")
			},
			wantErr:    true,
			wantSubstr: "not found",
		},
		{
			name: "not-found substring: plain error without gRPC status still detected",
			setup: func(ac *fakeAuthClient) {
				ac.accountErr = errString("rpc error: account " + address + " not found in state")
			},
			wantErr:    true,
			wantSubstr: "not found",
		},
		{
			name: "transport error: generic error is wrapped and returned",
			setup: func(ac *fakeAuthClient) {
				ac.accountErr = errString("connection refused")
			},
			wantErr:    true,
			wantSubstr: "auth query failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := &fakeAuthClient{}
			tt.setup(auth)

			pool := &testPool{auth: auth, tx: &fakeTxClient{}}
			b := cosmospool.New(pool, zerolog.Nop())

			gotNum, gotSeq, err := b.AccountInfo(context.Background(), address)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantSubstr,
					"error should mention the expected condition")
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantNumber, gotNum, "accountNumber mismatch")
			require.Equal(t, tt.wantSeq, gotSeq, "sequence mismatch")
			require.NotNil(t, auth.lastReq, "Account() should have been called")
			require.Equal(t, address, auth.lastReq.Address, "address forwarded to auth client")
		})
	}
}
