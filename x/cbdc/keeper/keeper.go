package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/peersyst/cbdc-node/x/cbdc/types"
)

var _ types.QueryServer = Querier{}

type Querier struct {
	Keeper
}

func NewQuerier(keeper Keeper) Querier {
	return Querier{Keeper: keeper}
}

type (
	Keeper struct {
		cdc        codec.Codec
		paramstore paramtypes.Subspace
		authority  string // the address capable of executing a cbdc mint/burn. Usually the gov module account
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

func (k Keeper) ExecuteMint(ctx sdk.Context, address string, amount sdk.Coin) error {
	if amount.Denom != k.cbdcDenom {
		return types.ErrInvalidDenom
	}
	if !amount.IsPositive() {
		return types.ErrInvalidAmount
	}

	toAddr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}

	coins := sdk.NewCoins(amount)
	if err := k.bk.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return err
	}
	if err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAddr, coins); err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeAddress, toAddr.String()),
			sdk.NewAttribute(types.AttributeAmount, amount.String()),
		),
	)

	return nil
}

func (k Keeper) ExecuteBurn(ctx sdk.Context, address string, amount sdk.Coin) error {
	if amount.Denom != k.cbdcDenom {
		return types.ErrInvalidDenom
	}
	if !amount.IsPositive() {
		return types.ErrInvalidAmount
	}

	fromAddr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}

	coins := sdk.NewCoins(amount)
	if err := k.bk.SendCoinsFromAccountToModule(ctx, fromAddr, types.ModuleName, coins); err != nil {
		return err
	}
	if err := k.bk.BurnCoins(ctx, types.ModuleName, coins); err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBurn,
			sdk.NewAttribute(types.AttributeAddress, fromAddr.String()),
			sdk.NewAttribute(types.AttributeAmount, amount.String()),
		),
	)

	return nil
}
