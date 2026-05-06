package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgBurn{}

func NewMsgBurn(authority string, address string, amount sdk.Coin) *MsgBurn {
	return &MsgBurn{
		Authority: authority,
		Address:   address,
		Amount:    amount,
	}
}
