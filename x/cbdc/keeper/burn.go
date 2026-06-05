package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peersyst/cbdc-node/x/cbdc/types"
)

// burningAllowed performs all pre-flight gating for a burn. The address param is
// kept for symmetry with mintingAllowed and for future address-based gating that
// would otherwise require a signature change.
func (k Keeper) burningAllowed(ctx sdk.Context, _ sdk.AccAddress, amount sdk.Coin) error {
	if k.GetParams(ctx).IssuancePaused {
		return types.ErrIssuancePaused
	}
	if err := amount.Validate(); err != nil {
		return err
	}
	if amount.Denom != k.cbdcDenom {
		return types.ErrInvalidDenom
	}
	if !amount.IsPositive() {
		return types.ErrInvalidAmount
	}
	return nil
}

func (k Keeper) BurnCoins(ctx sdk.Context, owner string, address sdk.AccAddress, amount sdk.Coin) error {
	if err := k.burningAllowed(ctx, address, amount); err != nil {
		return err
	}

	coins := sdk.NewCoins(amount)
	if err := k.bk.SendCoinsFromAccountToModule(ctx, address, types.ModuleName, coins); err != nil {
		return err
	}
	if err := k.bk.BurnCoins(ctx, types.ModuleName, coins); err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBurn,
			sdk.NewAttribute(types.AttributeOwner, owner),
			sdk.NewAttribute(types.AttributeAddress, address.String()),
			sdk.NewAttribute(types.AttributeAmount, amount.String()),
		),
	)

	return nil
}
