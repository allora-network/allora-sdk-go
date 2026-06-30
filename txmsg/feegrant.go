package txmsg

import (
	"github.com/brynbellomy/go-utils/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	feegrant "cosmossdk.io/x/feegrant"
)

// FeeAllowance is the minimal interface a feegrant allowance payload must
// satisfy to be packed into a MsgGrantAllowance.Any. It is satisfied by
// *feegrant.BasicAllowance, *feegrant.PeriodicAllowance, and
// *feegrant.AllowedMsgAllowance — all three are proto messages registered
// with the shared codec (feegrant.RegisterInterfaces is called in
// codec/registry.go's init).
//
// This interface is exported so callers can pass any of the
// feegrant-provided allowance types directly without an intermediate
// wrapper. The SDK's own feegrant.FeeAllowanceI is unexported; mirroring
// it here as a re-exported alias of proto.Message keeps the txmsg API
// usable from outside the package without depending on an unexported
// SDK interface.
type FeeAllowance = proto.Message

// NewGrantAllowance constructs a feegrant MsgGrantAllowance granting
// the given allowance from granter to grantee. The two addresses are
// bech32-validated. The allowance is packed into a *codectypes.Any
// using codectypes.NewAnyWithValue, the same path the SDK's own
// feegrant.NewMsgGrantAllowance takes — so a *feegrant.BasicAllowance
// passes through unchanged and a *feegrant.PeriodicAllowance likewise.
//
// If you need to grant a BasicAllowance, build one directly:
//
//	ba := &feegrant.BasicAllowance{
//	    SpendLimit: sdk.NewCoins(sdk.NewInt64Coin("uallo", 1_000_000)),
//	    Expiration: &expiry,
//	}
//	msg, err := txmsg.NewGrantAllowance(granter, grantee, ba)
//
// The basic allow SpendLimit is validated up front: a nil or zero
// spend limit is rejected so a caller cannot accidentally mint an
// unlimited grant by leaving the field empty. (The chain enforces a
// non-empty limit too, but only after the broadcast round-trip.)
//
// PeriodicAllowance's period spend limit and period-can-spend are
// similarly validated. AllowedMsgAllowance is accepted as-is; the chain
// validates the inner allowance Any more strictly than we can without
// duplicating the proto schema.
//
// The same feegrant.FeeAllowanceI interface unexported by the SDK
// (and the FeeAllowance type alias above) limits what types satisfy
// the constructor from outside this package: only registered proto
// messages, which is the exact set the chain will accept when it
// unpacks the Any.
func NewGrantAllowance(granter, grantee string, allowance FeeAllowance) (sdk.Msg, error) {
	if _, err := validateBech32(granter, "granter"); err != nil {
		return nil, err
	}
	if _, err := validateBech32(grantee, "grantee"); err != nil {
		return nil, err
	}

	if allowance == nil {
		return nil, errors.New("fee allowance is required")
	}

	any, err := codectypes.NewAnyWithValue(allowance)
	if err != nil {
		return nil, errors.Wrap(err, "failed to pack fee allowance into Any")
	}

	// Optional up-front validation of the most common allowance shapes:
	// a BasicAllowance with a non-zero, SDK-valid SpendLimit, and a
	// PeriodicAllowance with non-zero period spend / can-spend. These
	// are the only cases where the SDK provides a clean field-level
	// validation surface that we can mirror here without re-parsing
	// the proto schema. Other allowance types (AllowedMsgAllowance) are
	// accepted as-is and surface any field-shape errors on-chain.
	if ba, ok := allowance.(*feegrant.BasicAllowance); ok {
		if err := validateCoins(ba.SpendLimit, "basic allowance spend limit"); err != nil {
			return nil, err
		}
	}
	if pa, ok := allowance.(*feegrant.PeriodicAllowance); ok {
		if err := validateCoins(pa.PeriodSpendLimit, "periodic allowance period spend limit"); err != nil {
			return nil, err
		}
		if err := validateCoins(pa.PeriodCanSpend, "periodic allowance period can spend"); err != nil {
			return nil, err
		}
	}

	return &feegrant.MsgGrantAllowance{
		Granter:   granter,
		Grantee:   grantee,
		Allowance: any,
	}, nil
}

// NewRevokeAllowance constructs a feegrant MsgRevokeAllowance that
// revokes any active allowance from granter to grantee. Both addresses
// are bech32-validated.
//
// Note: the SDK's feegrant.NewMsgRevokeAllowance is a struct
// constructor (no error return) — but the resulting MsgRevokeAllowance
// only sets the address fields; any on-chain error (e.g. no
// allowance to revoke) surfaces at broadcast time, not at
// construction. This constructor deliberately follows the same shape
// so the resulting message is byte-identical to one built with
// feegrant.NewMsgRevokeAllowance on valid input.
func NewRevokeAllowance(granter, grantee string) (sdk.Msg, error) {
	if _, err := validateBech32(granter, "granter"); err != nil {
		return nil, err
	}
	if _, err := validateBech32(grantee, "grantee"); err != nil {
		return nil, err
	}
	return &feegrant.MsgRevokeAllowance{
		Granter: granter,
		Grantee: grantee,
	}, nil
}
