package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgMint{}

func NewMsgMint(authority string, address string, amount sdk.Coin) *MsgMint {
	return &MsgMint{
		Authority: authority,
		Address:   address,
		Amount:    amount,
	}
}
