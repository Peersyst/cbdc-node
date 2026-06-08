package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peersyst/cbdc-node/x/cbdc/types"
)

// mintingAllowed performs all pre-flight gating for a mint, including the owner
// authorization check so the keeper cannot be minted from by any caller that
// bypasses the msg server. Keeping it separate from MintCoins makes new gating
// easy to add and audit.
func (k Keeper) mintingAllowed(ctx sdk.Context, owner string, address sdk.AccAddress, amount sdk.Coin) error {
	params := k.GetParams(ctx)
	if owner != params.Owner {
		return types.ErrUnauthorized
	}
	if params.IssuancePaused {
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
	if k.bk.BlockedAddr(address) {
		return types.ErrBlockedAddr.Wrap(address.String())
	}
	if !k.bk.IsSendEnabledCoin(ctx, amount) {
		return types.ErrSendDisabled.Wrap(amount.Denom)
	}
	return nil
}

func (k Keeper) MintCoins(ctx sdk.Context, owner string, address sdk.AccAddress, amount sdk.Coin) error {
	if err := k.mintingAllowed(ctx, owner, address, amount); err != nil {
		return err
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
