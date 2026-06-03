package types //nolint:dupl // mint/burn tests are structurally similar by design

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peersyst/cbdc-node/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgMint_ValidateBasic(t *testing.T) {
	tt := []struct {
		name      string
		msg       MsgMint
		expectErr bool
	}{
		{
			name: "valid",
			msg:  MsgMint{Owner: sample.AccAddress(), Address: sample.AccAddress(), Amount: sdk.NewCoin("acbdc", math.NewInt(100))},
		},
		{
			name:      "invalid owner",
			msg:       MsgMint{Owner: "invalid", Address: sample.AccAddress(), Amount: sdk.NewCoin("acbdc", math.NewInt(100))},
			expectErr: true,
		},
		{
			name:      "invalid recipient address",
			msg:       MsgMint{Owner: sample.AccAddress(), Address: "invalid", Amount: sdk.NewCoin("acbdc", math.NewInt(100))},
			expectErr: true,
		},
		{
			name:      "zero amount",
			msg:       MsgMint{Owner: sample.AccAddress(), Address: sample.AccAddress(), Amount: sdk.NewCoin("acbdc", math.NewInt(0))},
			expectErr: true,
		},
		{
			// struct literal bypasses sdk.NewCoin, which panics on negatives
			name:      "negative amount",
			msg:       MsgMint{Owner: sample.AccAddress(), Address: sample.AccAddress(), Amount: sdk.Coin{Denom: "acbdc", Amount: math.NewInt(-100)}},
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
