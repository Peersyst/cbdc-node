package keeper

import (
	"errors"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/golang/mock/gomock"
	"github.com/peersyst/cbdc-node/x/poa/testutil"
	"github.com/peersyst/cbdc-node/x/poa/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServer_Burn(t *testing.T) {
	poaKeeper, ctx := poaKeeperTestSetup(t)
	msgServer := NewMsgServerImpl(*poaKeeper)

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
			amount:      sdk.NewCoin("XRP", math.NewInt(100)),
			expectedErr: govtypes.ErrInvalidSigner,
		},
		{
			name:        "should fail - invalid address",
			authority:   poaKeeper.GetAuthority(),
			address:     "invalidaddress",
			amount:      sdk.NewCoin("XRP", math.NewInt(100)),
			expectedErr: errors.New("decoding bech32 failed"),
		},
		{
			name:        "should fail - zero amount",
			authority:   poaKeeper.GetAuthority(),
			address:     "ethm1a0pd5cyew47pvgf7rd7axxy3humv9ev0nnkprp",
			amount:      sdk.NewCoin("XRP", math.NewInt(0)),
			expectedErr: types.ErrInvalidAmount,
		},
		{
			name:      "should pass",
			authority: poaKeeper.GetAuthority(),
			address:   "ethm1a0pd5cyew47pvgf7rd7axxy3humv9ev0nnkprp",
			amount:    sdk.NewCoin("XRP", math.NewInt(100)),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			msg := &types.MsgBurn{
				Authority: tc.authority,
				Address:   tc.address,
				Amount:    tc.amount,
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

func TestKeeper_ExecuteBurn(t *testing.T) {
	tt := []struct {
		name          string
		address       string
		amount        sdk.Coin
		bankMocks     func(ctx sdk.Context, bankKeeper *testutil.MockBankKeeper)
		expectedError error
	}{
		{
			name:          "should fail - zero amount",
			address:       "ethm1a0pd5cyew47pvgf7rd7axxy3humv9ev0nnkprp",
			amount:        sdk.NewCoin("XRP", math.NewInt(0)),
			bankMocks:     func(_ sdk.Context, _ *testutil.MockBankKeeper) {},
			expectedError: types.ErrInvalidAmount,
		},
		{
			name:          "should fail - invalid address",
			address:       "invalidaddress",
			amount:        sdk.NewCoin("XRP", math.NewInt(100)),
			bankMocks:     func(_ sdk.Context, _ *testutil.MockBankKeeper) {},
			expectedError: errors.New("decoding bech32 failed"),
		},
		{
			name:    "should fail - SendCoinsFromAccountToModule returns error",
			address: "ethm1a0pd5cyew47pvgf7rd7axxy3humv9ev0nnkprp",
			amount:  sdk.NewCoin("XRP", math.NewInt(100)),
			bankMocks: func(ctx sdk.Context, bankKeeper *testutil.MockBankKeeper) {
				bankKeeper.EXPECT().SendCoinsFromAccountToModule(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("bank send error"))
			},
			expectedError: errors.New("bank send error"),
		},
		{
			name:    "should fail - BurnCoins returns error",
			address: "ethm1a0pd5cyew47pvgf7rd7axxy3humv9ev0nnkprp",
			amount:  sdk.NewCoin("XRP", math.NewInt(100)),
			bankMocks: func(ctx sdk.Context, bankKeeper *testutil.MockBankKeeper) {
				bankKeeper.EXPECT().SendCoinsFromAccountToModule(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bankKeeper.EXPECT().BurnCoins(ctx, gomock.Any(), gomock.Any()).Return(errors.New("bank burn error"))
			},
			expectedError: errors.New("bank burn error"),
		},
		{
			name:    "should pass",
			address: "ethm1a0pd5cyew47pvgf7rd7axxy3humv9ev0nnkprp",
			amount:  sdk.NewCoin("XRP", math.NewInt(100)),
			bankMocks: func(ctx sdk.Context, bankKeeper *testutil.MockBankKeeper) {
				coins := sdk.NewCoins(sdk.NewCoin("XRP", math.NewInt(100)))
				bankKeeper.EXPECT().SendCoinsFromAccountToModule(ctx, gomock.Any(), types.ModuleName, coins).Return(nil)
				bankKeeper.EXPECT().BurnCoins(ctx, types.ModuleName, coins).Return(nil)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			keeper, ctx := setupPoaKeeper(t, func(_ sdk.Context, _ *testutil.MockStakingKeeper) {}, tc.bankMocks)

			err := keeper.ExecuteBurn(ctx, tc.address, tc.amount)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
