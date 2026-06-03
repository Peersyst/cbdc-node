package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgMint{}

func NewMsgMint(owner string, address string, amount sdk.Coin) *MsgMint {
	return &MsgMint{
		Owner:   owner,
		Address: address,
		Amount:  amount,
	}
}

// ValidateBasic performs stateless validation of the message. The amount denom
// is only checked for well-formedness here; the match against the configured
// CBDC denom is stateful and stays in the keeper.
func (msg *MsgMint) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return errorsmod.Wrapf(err, "invalid owner address (%s)", msg.Owner)
	}
	if _, err := sdk.AccAddressFromBech32(msg.Address); err != nil {
		return errorsmod.Wrapf(err, "invalid recipient address (%s)", msg.Address)
	}
	if err := msg.Amount.Validate(); err != nil {
		return err
	}
	if !msg.Amount.IsPositive() {
		return ErrInvalidAmount
	}
	return nil
}
