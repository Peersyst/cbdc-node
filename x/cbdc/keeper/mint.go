package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peersyst/cbdc-node/x/cbdc/types"
)

func (k Keeper) executeMint(ctx sdk.Context, owner string, address sdk.AccAddress, amount sdk.Coin) error {
	if err := amount.Validate(); err != nil {
		return err
	}
	if amount.Denom != k.cbdcDenom {
		return types.ErrInvalidDenom
	}
	if !amount.IsPositive() {
		return types.ErrInvalidAmount
	}

	coins := sdk.NewCoins(amount)
	if err := k.bk.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return err
	}
	if err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, address, coins); err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeOwner, owner),
			sdk.NewAttribute(types.AttributeAddress, address.String()),
			sdk.NewAttribute(types.AttributeAmount, amount.String()),
		),
	)

	return nil
}
