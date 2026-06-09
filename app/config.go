package app

import "os"

type EVMOptionsFn func(uint64) error

// AccountAddressPrefix is the Bech32 prefix for account addresses. It reads
// CHAIN_PREFIX at process start (e.g. "hnl" for the HNL CBDC deployment) and
// falls back to the Ethermint-conventional "ethm" when unset (tests, local).
var AccountAddressPrefix = func() string {
	if v := os.Getenv("CHAIN_PREFIX"); v != "" {
		return v
	}
	return "ethm"
}()

const (
	Bip44CoinType = 60
	Name          = "cbdc"
	// BaseDenom defines to the default denomination used in EVM
	BaseDenom                  = "hnl"
	Denom                      = "CBDC"
	DenomDescription           = "CBDC is the digital central bank currency."
	DenomName                  = "CBDC"
	DenomSymbol                = "CBDC"
	NativeErc20ContractAddress = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
)
