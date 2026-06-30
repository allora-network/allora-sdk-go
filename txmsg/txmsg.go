// Package txmsg provides pure, validated constructors for Cosmos SDK message
// types used by the Allora SDK. The constructors return sdk.Msg values (or
// concrete typed pointers) without touching signing, broadcasting, network, or
// the allora package's TxManager — so they compose cleanly with
// allora.CreateUnsignedTx([]sdk.Msg, params) and any other tx-pipeline that
// only needs the message payload.
//
// The constructors in this package perform input validation up front
// (non-empty required fields, valid bech32 addresses, positive amounts and
// topic IDs, etc.) and wrap any failure with brynbellomy/go-utils/errors so
// callers can use errors.Is/As against the underlying validation errors
// returned by the SDK.
//
// All emitted messages are registered with the shared codec in
// codec/registry.go: bank (x/bank/types), feegrant
// (cosmossdk.io/x/feegrant), and the live emissions types package aliased as
// emissionsv9 in codec/registry.go. If the registered emissions version
// changes, the concrete types returned from the New*Workers and New*Payload
// constructors change with it — keep the type assertions in callers in sync.
package txmsg

import (
	"github.com/brynbellomy/go-utils/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// validateBech32 decodes addr with the SDK's bech32 decoder (which enforces
// the account bech32 prefix configured by allora.init) and returns a wrapped
// error identifying the field on failure. The decoded address is returned for
// callers that want to use the raw AccAddress form.
//
// Unlike sdk.AccAddressFromBech32, this returns a *non-nil* error for the
// empty string so callers can rely on a single error path for "missing
// address" and "malformed address" — the SDK treats an empty string as a
// valid AccAddress (20 zero bytes), which would slip past the constructor's
// guard and surface as a confusing on-chain failure later.
func validateBech32(addr, field string) (sdk.AccAddress, error) {
	if addr == "" {
		return nil, errors.Errorf("%s address is required", field)
	}
	decoded, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid %s bech32 address %q", field, addr)
	}
	return decoded, nil
}

// validateCoins returns a wrapped error if amount is empty or fails the SDK's
// own validation (negative amounts, duplicate denominations, etc.). Callers
// that need positive amounts in a specific denomination should chain an
// additional denom-specific check after this returns nil.
func validateCoins(amount sdk.Coins, field string) error {
	if amount == nil || amount.IsZero() {
		return errors.Errorf("%s amount is required and must be non-zero", field)
	}
	if err := amount.Validate(); err != nil {
		return errors.Wrapf(err, "invalid %s amount", field)
	}
	return nil
}
