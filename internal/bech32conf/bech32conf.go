// Package bech32conf configures the Cosmos SDK's global bech32 address prefixes
// for the Allora network. It exists as a dedicated importable package (rather
// than living in the top-level allora package's init) so that subpackages such
// as txmsg can ensure the prefix is configured without depending on the parent
// allora package's initialization side effects.
//
// Importing this package (directly or transitively) triggers an init() that
// sets the account, validator, and consensus-node prefixes. The SetBech32*
// calls are idempotent field assignments on the SDK global config, so multiple
// imports are safe.
package bech32conf

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AlloraBech32Prefix is the human-readable prefix for Allora account addresses.
const AlloraBech32Prefix = "allo"

func init() {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(AlloraBech32Prefix, AlloraBech32Prefix+"pub")
	cfg.SetBech32PrefixForValidator(AlloraBech32Prefix+"valoper", AlloraBech32Prefix+"valoperpub")
	cfg.SetBech32PrefixForConsensusNode(AlloraBech32Prefix+"valcons", AlloraBech32Prefix+"valconspub")
}
