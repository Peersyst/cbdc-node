package app

type EVMOptionsFn func(uint64) error

const (
	AccountAddressPrefix = "hn"
	Bip44CoinType        = 60
	Name                 = "hnl"
	// BaseDenom defines to the default denomination used in EVM
	BaseDenom                  = "mhnl"
	Denom                      = "hnl"
	DenomDescription           = "Lempira digital. La divisa nativa de la red CBDC de Honduras."
	DenomName                  = "HNL"
	DenomSymbol                = "HNL"
	NativeErc20ContractAddress = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
)
