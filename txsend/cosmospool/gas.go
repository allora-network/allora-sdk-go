package cosmospool

import (
	sdkmath "cosmossdk.io/math"
)

// DefaultGasAdjustment is the default multiplier applied to simulated gas usage
// to derive the gas limit set on a transaction. It mirrors the cosmos-cli
// --gas-adjustment default and leaves headroom for execution-time variance
// between simulation and inclusion so a tx is not rejected for out-of-gas.
const DefaultGasAdjustment = 1.5

// FallbackGasEstimate is the static gas returned when simulation fails. It is
// intentionally conservative: large enough for a typical multi-message tx
// while not overpaying. See EstimateGas for the rationale.
const FallbackGasEstimate uint64 = 200_000

// FeeTier selects a gas-price band for fee payment. It maps to a concrete
// legacy-decimal gas price via GetGasPrice, letting callers request "cheap,
// slow", "balanced", or "fast, expensive" inclusion without hardcoding prices.
type FeeTier string

const (
	// FeeTierLow is the cheapest tier: lowest gas price, slowest inclusion.
	// Use for non-urgent txs willing to wait in the mempool.
	FeeTierLow FeeTier = "low"
	// FeeTierMedium is the balanced tier: moderate gas price, typical
	// inclusion. This is the default for most txs.
	FeeTierMedium FeeTier = "medium"
	// FeeTierHigh is the priority tier: highest gas price, fastest inclusion.
	// Use for time-sensitive txs.
	FeeTierHigh FeeTier = "high"
)

// gasPrices holds the legacy-decimal gas price for each FeeTier. Values are in
// the smallest token unit per unit of gas and are intentionally conservative
// defaults; callers needing chain-specific prices can override by ignoring
// GetGasPrice and passing their own price when building the tx.
var gasPrices = map[FeeTier]sdkmath.LegacyDec{
	FeeTierLow:    sdkmath.LegacyMustNewDecFromStr("0.001"),
	FeeTierMedium: sdkmath.LegacyMustNewDecFromStr("0.01"),
	FeeTierHigh:   sdkmath.LegacyMustNewDecFromStr("0.025"),
}

// GetGasPrice returns the gas price (smallest token unit per gas unit) for the
// given tier. An unknown tier falls back to FeeTierMedium so a caller passing a
// typo'd or unsupported tier still gets a sane, moderate price rather than a
// zero (which would make the tx free and likely deprioritized).
func GetGasPrice(tier FeeTier) sdkmath.LegacyDec {
	if p, ok := gasPrices[tier]; ok {
		return p
	}
	return gasPrices[FeeTierMedium]
}
