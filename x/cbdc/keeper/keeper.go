package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/peersyst/cbdc-node/x/cbdc/types"
)

type (
	Keeper struct {
		cdc        codec.Codec
		paramstore paramtypes.Subspace
		authority  string // the address capable of updating module params (e.g. rotating the mint/burn owner). Usually the gov module account
		bk         types.BankKeeper
		cbdcDenom  string // the only denom that can be minted/burned via this module
	}
)

func NewKeeper(
	cdc codec.Codec,
	ps paramtypes.Subspace,
	bk types.BankKeeper,
	authority string,
	cbdcDenom string,
) *Keeper {
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(err)
	}

	if err := sdk.ValidateDenom(cbdcDenom); err != nil {
		panic(err)
	}

	return &Keeper{
		cdc:        cdc,
		paramstore: ps,
		authority:  authority,
		bk:         bk,
		cbdcDenom:  cbdcDenom,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetAuthority() string {
	return k.authority
}
