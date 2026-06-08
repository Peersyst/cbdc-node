package cbdc

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/appmodule"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/peersyst/cbdc-node/x/cbdc/keeper"
	"github.com/peersyst/cbdc-node/x/cbdc/types"
)

var (
	_ module.AppModuleBasic      = (*AppModule)(nil)
	_ module.AppModuleSimulation = (*AppModule)(nil)
	_ module.HasGenesis          = (*AppModule)(nil)
	_ appmodule.AppModule        = (*AppModule)(nil)
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

type AppModuleBasic struct {
	cdc codec.BinaryCodec
}

func NewAppModuleBasic(cdc codec.BinaryCodec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

func (AppModuleBasic) Name() string {
	return types.ModuleName
}

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

func (a AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
}

func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var genState types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return genState.Validate()
}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

type AppModule struct {
	AppModuleBasic
	keeper   keeper.Keeper
	ak       types.AccountKeeper
	registry cdctypes.InterfaceRegistry
}

func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	ak types.AccountKeeper,
	registry cdctypes.InterfaceRegistry,
) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc),
		keeper:         keeper,
		ak:             ak,
		registry:       registry,
	}
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQuerier(am.keeper))
}

func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	ir.RegisterRoute(types.ModuleName, "module-account-exists", am.moduleAccountInvariant())
	ir.RegisterRoute(types.ModuleName, "owner-valid", am.ownerValidInvariant())
}

// moduleAccountInvariant checks that the cbdc module account still exists.
// If it were ever wiped, MintCoins/BurnCoins would fail at runtime, so this
// halts the chain instead. The account is created in InitGenesis.
func (am AppModule) moduleAccountInvariant() sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		addr := authtypes.NewModuleAddress(types.ModuleName)
		broken := am.ak.GetAccount(ctx, addr) == nil
		return sdk.FormatInvariant(
			types.ModuleName,
			"module-account-exists",
			fmt.Sprintf("cbdc module account %s does not exist\n", addr),
		), broken
	}
}

// ownerValidInvariant checks that the stored mint/burn owner is well-formed
// (empty, meaning disabled, or a valid bech32 address). It catches state
// corruption at a block boundary instead of at the next mint/burn attempt.
func (am AppModule) ownerValidInvariant() sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		params := am.keeper.GetParams(ctx)
		broken := params.Validate() != nil
		return sdk.FormatInvariant(
			types.ModuleName,
			"owner-valid",
			fmt.Sprintf("cbdc params owner %q is not a valid address\n", params.Owner),
		), broken
	}
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) {
	var genState types.GenesisState
	cdc.MustUnmarshalJSON(gs, &genState)

	am.keeper.InitGenesis(ctx, genState)

	// To create module account
	am.ak.GetModuleAccount(ctx, am.Name())
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := am.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(genState)
}

func (AppModule) ConsensusVersion() uint64 { return 1 }

func (am AppModule) IsOnePerModuleType() {}

func (am AppModule) IsAppModule() {}
