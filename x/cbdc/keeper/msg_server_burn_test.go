package keeper

import (
	"errors"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/peersyst/cbdc-node/x/cbdc/testutil"
	"github.com/peersyst/cbdc-node/x/cbdc/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServer_Burn(t *testing.T) { //nolint:dupl
	cbdcKeeper, ctx := cbdcKeeperTestSetup(t)
	msgServer := NewMsgServerImpl(*cbdcKeeper)

	tt := []struct {
		name        string
		authority   string
		address     string
		amount      sdk.Coin
		expectedErr error
	}{
		{
			name:        "should fail - invalid authority",
			authority:   "invalidauthority",
			address:     "ethm1a0pd5cyew47pvgf7rd7axxy3humv9ev0nnkprp",
			amount:      sdk.NewCoin(testCBDCDenom, math.NewInt(100)),
			expectedErr: types.ErrUnauthorized,
		},
		{
			name:        "should fail - invalid address",
			authority:   cbdcKeeper.GetAuthority(),
			address:     "invalidaddress",
			amount:      sdk.NewCoin(testCBDCDenom, math.NewInt(100)),
			expectedErr: errors.New("decoding bech32 failed"),
		},
		{
			name:        "should fail - zero amount",
			authority:   cbdcKeeper.GetAuthority(),
			address:     "ethm1a0pd5cyew47pvgf7rd7axxy3humv9ev0nnkprp",
			amount:      sdk.NewCoin(testCBDCDenom, math.NewInt(0)),
			expectedErr: types.ErrInvalidAmount,
		},
		{
			// struct literal bypasses sdk.NewCoin, which panics on negatives
			name:        "should fail - negative amount",
			authority:   cbdcKeeper.GetAuthority(),
			address:     "ethm1a0pd5cyew47pvgf7rd7axxy3humv9ev0nnkprp",
			amount:      sdk.Coin{Denom: testCBDCDenom, Amount: math.NewInt(-100)},
			expectedErr: errors.New("negative coin amount"),
		},
		{
			name:        "should fail - wrong denom",
			authority:   cbdcKeeper.GetAuthority(),
			address:     "ethm1a0pd5cyew47pvgf7rd7axxy3humv9ev0nnkprp",
			amount:      sdk.NewCoin("axrp", math.NewInt(100)),
			expectedErr: types.ErrInvalidDenom,
		},
		{
			name:      "should pass",
			authority: cbdcKeeper.GetAuthority(),
			address:   "ethm1a0pd5cyew47pvgf7rd7axxy3humv9ev0nnkprp",
			amount:    sdk.NewCoin(testCBDCDenom, math.NewInt(100)),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			msg := &types.MsgBurn{
				Owner:   tc.authority,
				Address: tc.address,
				Amount:  tc.amount,
			}

			_, err := msgServer.Burn(ctx, msg)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestKeeper_BurnCoins(t *testing.T) {
	address := sdk.MustAccAddressFromBech32("ethm1a0pd5cyew47pvgf7rd7axxy3humv9ev0nnkprp")

	tt := []struct {
		name          string
		address       sdk.AccAddress
		amount        sdk.Coin
		bankMocks     func(ctx sdk.Context, bankKeeper *testutil.MockBankKeeper)
		expectedError error
	}{
		{
			name:          "should fail - zero amount",
			address:       address,
			amount:        sdk.NewCoin(testCBDCDenom, math.NewInt(0)),
			bankMocks:     func(_ sdk.Context, _ *testutil.MockBankKeeper) {},
			expectedError: types.ErrInvalidAmount,
		},
		{
			name:          "should fail - wrong denom",
			address:       address,
			amount:        sdk.NewCoin("axrp", math.NewInt(100)),
			bankMocks:     func(_ sdk.Context, _ *testutil.MockBankKeeper) {},
			expectedError: types.ErrInvalidDenom,
		},
		{
			name:    "should fail - SendCoinsFromAccountToModule returns error",
			address: address,
			amount:  sdk.NewCoin(testCBDCDenom, math.NewInt(100)),
			bankMocks: func(ctx sdk.Context, bankKeeper *testutil.MockBankKeeper) {
				bankKeeper.EXPECT().SendCoinsFromAccountToModule(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("bank send error"))
			},
			expectedError: errors.New("bank send error"),
		},
		{
			name:    "should fail - BurnCoins returns error",
			address: address,
			amount:  sdk.NewCoin(testCBDCDenom, math.NewInt(100)),
			bankMocks: func(ctx sdk.Context, bankKeeper *testutil.MockBankKeeper) {
				bankKeeper.EXPECT().SendCoinsFromAccountToModule(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bankKeeper.EXPECT().BurnCoins(ctx, gomock.Any(), gomock.Any()).Return(errors.New("bank burn error"))
			},
			expectedError: errors.New("bank burn error"),
		},
		{
			name:    "should pass",
			address: address,
			amount:  sdk.NewCoin(testCBDCDenom, math.NewInt(100)),
			bankMocks: func(ctx sdk.Context, bankKeeper *testutil.MockBankKeeper) {
				coins := sdk.NewCoins(sdk.NewCoin(testCBDCDenom, math.NewInt(100)))
				bankKeeper.EXPECT().SendCoinsFromAccountToModule(ctx, gomock.Any(), types.ModuleName, coins).Return(nil)
				bankKeeper.EXPECT().BurnCoins(ctx, types.ModuleName, coins).Return(nil)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			keeper, ctx := setupCbdcKeeper(t, tc.bankMocks)

			err := keeper.BurnCoins(ctx, keeper.GetAuthority(), tc.address, tc.amount)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
