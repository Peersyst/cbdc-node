package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peersyst/cbdc-node/x/cbdc/types"
)

func (k Keeper) executeMint(ctx sdk.Context, owner string, address string, amount sdk.Coin) error {
	if err := amount.Validate(); err != nil {
		return err
	}
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
			sdk.NewAttribute(types.AttributeOwner, owner),
			sdk.NewAttribute(types.AttributeAddress, toAddr.String()),
			sdk.NewAttribute(types.AttributeAmount, amount.String()),
		),
	)

	return nil
}
