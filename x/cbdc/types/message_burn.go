package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgBurn{}

func NewMsgBurn(owner string, address string, amount sdk.Coin) *MsgBurn {
	return &MsgBurn{
		Owner:   owner,
		Address: address,
		Amount:  amount,
	}
}

// ValidateBasic performs stateless validation of the message. The amount denom
// is only checked for well-formedness here; the match against the configured
// CBDC denom is stateful and stays in the keeper.
func (msg *MsgBurn) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return errorsmod.Wrapf(err, "invalid owner address (%s)", msg.Owner)
	}
	if _, err := sdk.AccAddressFromBech32(msg.Address); err != nil {
		return errorsmod.Wrapf(err, "invalid address (%s)", msg.Address)
	}
	if err := msg.Amount.Validate(); err != nil {
		return err
	}
	if !msg.Amount.IsPositive() {
		return ErrInvalidAmount
	}
	return nil
}
