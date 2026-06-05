package keeper

import (
	"testing"
	"time"

	storetypes "cosmossdk.io/store/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/golang/mock/gomock"
	"github.com/peersyst/cbdc-node/x/cbdc/testutil"
	"github.com/peersyst/cbdc-node/x/cbdc/types"
)

const (
	accountAddressPrefix = "ethm"
	bip44CoinType        = 60
	testCBDCDenom        = "acbdc"
	// testGovAuthority is the params updater (gov module account stand-in).
	testGovAuthority = "ethm1wunfhl05vc8r8xxnnp8gt62wa54r6y52pg03zq"
	// testOwner is the mint/burn owner stored in params.
	testOwner = "ethm1j2arnn0sajut2w8gnaxumrlkkem2c9vz5sfj6k"
)

func setupSdkConfig() {
	accountPubKeyPrefix := accountAddressPrefix + "pub"
	validatorAddressPrefix := accountAddressPrefix + "valoper"
	validatorPubKeyPrefix := accountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := accountAddressPrefix + "valcons"
	consNodePubKeyPrefix := accountAddressPrefix + "valconspub"

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(accountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	config.SetCoinType(bip44CoinType)
	config.SetPurpose(sdk.Purpose)
}

func getBankKeeperMock(t *testing.T, ctx sdk.Context, setExpectations func(ctx sdk.Context, bankKeeper *testutil.MockBankKeeper)) *testutil.MockBankKeeper {
	ctrl := gomock.NewController(t)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	setExpectations(ctx, bankKeeper)
	return bankKeeper
}

func getCtxMock(t *testing.T, key *storetypes.KVStoreKey, tsKey *storetypes.TransientStoreKey) sdk.Context {
	setupSdkConfig()

	testCtx := sdktestutil.DefaultContextWithDB(t, key, tsKey)
	ctx := testCtx.Ctx.WithBlockHeader(tmproto.Header{Time: time.Now()})
	return ctx
}

func getMockedCbdcKeeper(_ *testing.T, key *storetypes.KVStoreKey, tsKey *storetypes.TransientStoreKey, ctx sdk.Context, bankKeeper *testutil.MockBankKeeper) *Keeper {
	encCfg := moduletestutil.MakeTestEncodingConfig()

	types.RegisterInterfaces(encCfg.InterfaceRegistry)

	cbdcKeeper := NewKeeper(
		encCfg.Codec,
		paramtypes.NewSubspace(encCfg.Codec, encCfg.Amino, key, tsKey, "cbdc"),
		bankKeeper,
		testGovAuthority,
		testCBDCDenom,
	)
	cbdcKeeper.SetParams(ctx, types.NewParams(testOwner))

	return cbdcKeeper
}

func setupCbdcKeeper(t *testing.T, setBankExpectations func(ctx sdk.Context, bankKeeper *testutil.MockBankKeeper)) (*Keeper, sdk.Context) {
	key := storetypes.NewKVStoreKey(types.StoreKey)
	tsKey := storetypes.NewTransientStoreKey("test")

	ctx := getCtxMock(t, key, tsKey)
	bankKeeper := getBankKeeperMock(t, ctx, setBankExpectations)

	return getMockedCbdcKeeper(t, key, tsKey, ctx, bankKeeper), ctx
}

func cbdcKeeperTestSetup(t *testing.T) (*Keeper, sdk.Context) {
	bankExpectations := func(ctx sdk.Context, bankKeeper *testutil.MockBankKeeper) {
		bankKeeper.EXPECT().BlockedAddr(gomock.Any()).Return(false).AnyTimes()
		bankKeeper.EXPECT().IsSendEnabledCoin(ctx, gomock.Any()).Return(true).AnyTimes()
		bankKeeper.EXPECT().MintCoins(ctx, gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		bankKeeper.EXPECT().SendCoinsFromModuleToAccount(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		bankKeeper.EXPECT().BurnCoins(ctx, gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		bankKeeper.EXPECT().SendCoinsFromAccountToModule(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	}

	return setupCbdcKeeper(t, bankExpectations)
}
