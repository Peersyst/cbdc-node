package integration

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/evm/testutil/integration/evm/factory"
	"github.com/cosmos/evm/testutil/integration/evm/grpc"
	"github.com/cosmos/evm/testutil/keyring"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/peersyst/cbdc-node/app"
	cbdccommon "github.com/peersyst/cbdc-node/testutil/integration/cbdc/common"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite

	network     *Network
	keyring     keyring.Keyring
	factory     factory.TxFactory
	grpcHandler grpc.Handler
}

func (s *TestSuite) Network() *Network {
	return s.network
}

func (s *TestSuite) SetupSuite() {
	s.network.SetupSdkConfig()
	s.Require().Equal(sdk.GetConfig().GetBech32AccountAddrPrefix(), "ethm")
}

func (s *TestSuite) SetupTest() {
	// Check that the network was created successfully
	kr := keyring.New(5)

	customGenesis := cbdccommon.CustomGenesisState{}

	evmGen := evmtypes.DefaultGenesisState()

	evmGen.Params.EvmDenom = app.BaseDenom

	customGenesis[evmtypes.ModuleName] = evmGen

	s.network = NewIntegrationNetwork(
		cbdccommon.WithPreFundedAccounts(kr.GetAllAccAddrs()...),
		cbdccommon.WithAmountOfValidators(5),
		cbdccommon.WithCustomGenesis(customGenesis),
		cbdccommon.WithBondDenom("apoa"),
		cbdccommon.WithMaxValidators(7),
		cbdccommon.WithMinDepositAmt(sdkmath.NewInt(1)),
		cbdccommon.WithValidatorOperators(kr.GetAllAccAddrs()),
	)
	s.Require().NotNil(s.network)

	// TODO: Update when migrating to v10
	grpcHandler := grpc.NewIntegrationHandler(s.network)

	// TODO: Update when migrating to v10
	s.factory = factory.New(s.network, grpcHandler)
	s.keyring = kr
	s.grpcHandler = grpcHandler
}
