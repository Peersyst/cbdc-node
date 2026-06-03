package types

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peersyst/cbdc-node/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgBurn_ValidateBasic(t *testing.T) { //nolint:dupl
	tt := []struct {
		name      string
		msg       MsgBurn
		expectErr bool
	}{
		{
			name: "valid",
			msg:  MsgBurn{Owner: sample.AccAddress(), Address: sample.AccAddress(), Amount: sdk.NewCoin("acbdc", math.NewInt(100))},
		},
		{
			name:      "invalid owner",
			msg:       MsgBurn{Owner: "invalid", Address: sample.AccAddress(), Amount: sdk.NewCoin("acbdc", math.NewInt(100))},
			expectErr: true,
		},
		{
			name:      "invalid address",
			msg:       MsgBurn{Owner: sample.AccAddress(), Address: "invalid", Amount: sdk.NewCoin("acbdc", math.NewInt(100))},
			expectErr: true,
		},
		{
			name:      "zero amount",
			msg:       MsgBurn{Owner: sample.AccAddress(), Address: sample.AccAddress(), Amount: sdk.NewCoin("acbdc", math.NewInt(0))},
			expectErr: true,
		},
		{
			// struct literal bypasses sdk.NewCoin, which panics on negatives
			name:      "negative amount",
			msg:       MsgBurn{Owner: sample.AccAddress(), Address: sample.AccAddress(), Amount: sdk.Coin{Denom: "acbdc", Amount: math.NewInt(-100)}},
			expectErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
