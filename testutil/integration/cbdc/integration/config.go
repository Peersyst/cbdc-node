// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package cbdcintegration

import (
	"fmt"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	testtx "github.com/cosmos/evm/testutil/tx"
	cbdccommon "github.com/peersyst/cbdc-node/testutil/integration/cbdc/common"
)

// DefaultIntegrationConfig returns the default configuration for a chain.
func DefaultIntegrationConfig() cbdccommon.Config {
	account, _ := testtx.NewAccAddressAndKey()
	config := cbdccommon.DefaultConfig()
	config.AmountOfValidators = 3
	config.PreFundedAccounts = []sdktypes.AccAddress{account}
	return config
}

// getGenAccountsAndBalances takes the network configuration and returns the used
// genesis accounts and balances.
//
// NOTE: If the balances are set, the pre-funded accounts are ignored.
func getGenAccountsAndBalances(cfg cbdccommon.Config, validators []stakingtypes.Validator) (genAccounts []authtypes.GenesisAccount, balances []banktypes.Balance) {
	if len(cfg.Balances) > 0 {
		balances = cfg.Balances
		accounts := getAccAddrsFromBalances(balances)
		genAccounts = createGenesisAccounts(accounts)
	} else {
		genAccounts = createGenesisAccounts(cfg.PreFundedAccounts)
		balances = createBalances(cfg.PreFundedAccounts, append(cfg.OtherCoinDenom, cfg.Denom))
	}

	// append validators to genesis accounts and balances
	valAccs := make([]sdktypes.AccAddress, len(validators))
	for i, v := range validators {
		valAddr, err := sdktypes.ValAddressFromBech32(v.OperatorAddress)
		if err != nil {
			panic(fmt.Sprintf("failed to derive validator address from %q: %s", v.OperatorAddress, err.Error()))
		}
		valAccs[i] = sdktypes.AccAddress(valAddr.Bytes())
	}
	genAccounts = append(genAccounts, createGenesisAccounts(valAccs)...)

	return
}
