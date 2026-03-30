package app

type EVMOptionsFn func(uint64) error

const (
	AccountAddressPrefix = "ethm"
	Bip44CoinType        = 60
	Name                 = "cbdc"
	// BaseDenom defines to the default denomination used in EVM
	BaseDenom                  = "acbdc"
	Denom                      = "CBDC"
	DenomDescription           = "CBDC is the digital central bank currency."
	DenomName                  = "CBDC"
	DenomSymbol                = "CBDC"
	NativeErc20ContractAddress = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
)
