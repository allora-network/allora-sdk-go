package txmsg

import (
	"github.com/brynbellomy/go-utils/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// NewSend constructs a bank MsgSend from the given bech32 addresses and
// amount. Both addresses are validated against the Allora bech32 prefix
// (configured in the parent allora package's init) and the amount is
// checked for non-zero + passes the SDK's own sdk.Coins validation.
//
// The returned sdk.Msg is the concrete *banktypes.MsgSend; the broad
// interface return type matches the signature expected by
// allora.CreateUnsignedTx so callers do not need an extra type assertion
// in the common case.
func NewSend(from, to string, amount sdk.Coins) (sdk.Msg, error) {
	if _, err := validateBech32(from, "from"); err != nil {
		return nil, err
	}
	if _, err := validateBech32(to, "to"); err != nil {
		return nil, err
	}
	if err := validateCoins(amount, "send"); err != nil {
		return nil, err
	}

	// Re-derive AccAddress from the (already validated) bech32 strings so
	// the FromAddress/ToAddress fields carry the canonical 20-byte form
	// rather than the bech32 string. The SDK's NewMsgSend takes
	// sdk.AccAddress and converts back to its bech32 string, so the
	// round-trip is a no-op for valid input but defends against any
	// future AccAddressFromBech32 returning a non-canonical form.
	fromAddr, _ := sdk.AccAddressFromBech32(from)
	toAddr, _ := sdk.AccAddressFromBech32(to)
	return banktypes.NewMsgSend(fromAddr, toAddr, amount), nil
}

// NewMultiSend constructs a bank MsgMultiSend from the given inputs and
// outputs. Each input and output's address is validated as a bech32
// Allora account, each entry's coins are checked with validateCoins, and
// the sum of input coins must equal the sum of output coins — a Cosmos
// invariant that the chain enforces on its own, but failing fast here
// saves a wasted broadcast round-trip and an on-chain error that
// surfaces far from the call site.
//
// Passing nil or empty slices is rejected; a multi-send with zero entries
// has no economic effect and the chain would reject it anyway.
//
// As with NewSend, the returned sdk.Msg is the concrete
// *banktypes.MsgMultiSend.
func NewMultiSend(inputs []banktypes.Input, outputs []banktypes.Output) (sdk.Msg, error) {
	if len(inputs) == 0 {
		return nil, errors.New("multi-send requires at least one input")
	}
	if len(outputs) == 0 {
		return nil, errors.New("multi-send requires at least one output")
	}

	var totalIn sdk.Coins
	for i, in := range inputs {
		if _, err := validateBech32(in.Address, fmtField("input", i, "address")); err != nil {
			return nil, err
		}
		if err := validateCoins(in.Coins, fmtField("input", i, "coins")); err != nil {
			return nil, err
		}
		totalIn = totalIn.Add(in.Coins...)
	}

	var totalOut sdk.Coins
	for i, out := range outputs {
		if _, err := validateBech32(out.Address, fmtField("output", i, "address")); err != nil {
			return nil, err
		}
		if err := validateCoins(out.Coins, fmtField("output", i, "coins")); err != nil {
			return nil, err
		}
		totalOut = totalOut.Add(out.Coins...)
	}

	// Cosmos requires input == output totals. sdk.Coins.Equal handles
	// sorting and denomination matching, so a caller-supplied
	// {100uatom, 50uallo} input vs {50uallo, 100uatom} output is treated
	// as equal and accepted.
	if !totalIn.Equal(totalOut) {
		return nil, errors.Errorf(
			"multi-send inputs (%s) do not equal outputs (%s)",
			totalIn, totalOut,
		)
	}

	return &banktypes.MsgMultiSend{Inputs: inputs, Outputs: outputs}, nil
}

// fmtField formats a field path like "input[2].address" for use in
// validation errors. It is package-private to keep the API surface
// small — callers should not need to construct these strings.
func fmtField(group string, index int, field string) string {
	return fmtGroupIndex(group, index) + "." + field
}

// fmtGroupIndex formats a group path like "input[2]". It avoids fmt
// so the call site stays free of an extra import and a runtime
// reflection cost per validation error.
func fmtGroupIndex(group string, index int) string {
	if index == 0 {
		return group + "[0]"
	}
	// Build the decimal digits backwards into a small stack buffer.
	const decimal = 10
	var buf [20]byte
	pos := len(buf)
	n := index
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%decimal)
		n /= decimal
	}
	return group + "[" + string(buf[pos:]) + "]"
}
