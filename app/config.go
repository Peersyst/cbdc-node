package app

import (
	"cmp"
	"os"
)

type EVMOptionsFn func(uint64) error

// AccountAddressPrefix (the Bech32 account prefix) and BaseDenom (the EVM / CBDC
// base denomination) are deployment-specific, so the same binary can serve any
// deployment. They are read from the environment at process start and fall back
// to the Ethermint-conventional defaults when unset (tests, local runs):
//
//	CHAIN_PREFIX=hnl   CHAIN_DENOM=hnl   // HNL CBDC production
var (
	AccountAddressPrefix = cmp.Or(os.Getenv("CHAIN_PREFIX"), "ethm")
	BaseDenom            = cmp.Or(os.Getenv("CHAIN_DENOM"), "acbdc")
)

const (
	Bip44CoinType              = 60
	Name                       = "cbdc"
	Denom                      = "CBDC"
	DenomDescription           = "CBDC is the digital central bank currency."
	DenomName                  = "CBDC"
	DenomSymbol                = "CBDC"
	NativeErc20ContractAddress = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
)
