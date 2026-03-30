package integration

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
	precisebanktypes "github.com/cosmos/evm/x/precisebank/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	cbdccommon "github.com/peersyst/cbdc-node/testutil/integration/cbdc/common"
	cbdcintegration "github.com/peersyst/cbdc-node/testutil/integration/cbdc/integration"
	poatypes "github.com/peersyst/cbdc-node/x/poa/types"
)

// TODO: Update when migrating to v10
var _ cbdcintegration.Network = (*Network)(nil)

type Network struct {
	cbdcintegration.IntegrationNetwork
}

func NewIntegrationNetwork(opts ...cbdccommon.ConfigOption) *Network {
	network := cbdcintegration.New(opts...)
	return &Network{
		IntegrationNetwork: *network,
	}
}

func (n *Network) SetupSdkConfig() {
	cbdccommon.SetupSdkConfig()
}

func (n *Network) GetERC20Client() erc20types.QueryClient {
	return cbdccommon.GetERC20Client(n)
}

func (n *Network) GetEvmClient() evmtypes.QueryClient {
	return cbdccommon.GetEvmClient(n)
}

func (n *Network) GetGovClient() govtypes.QueryClient {
	return cbdccommon.GetGovClient(n)
}

func (n *Network) GetBankClient() banktypes.QueryClient {
	return cbdccommon.GetBankClient(n)
}

func (n *Network) GetFeeMarketClient() feemarkettypes.QueryClient {
	return cbdccommon.GetFeeMarketClient(n)
}

func (n *Network) GetAuthClient() authtypes.QueryClient {
	return cbdccommon.GetAuthClient(n)
}

func (n *Network) GetAuthzClient() authz.QueryClient {
	return cbdccommon.GetAuthzClient(n)
}

func (n *Network) GetStakingClient() stakingtypes.QueryClient {
	return cbdccommon.GetStakingClient(n)
}

func (n *Network) GetSlashingClient() slashingtypes.QueryClient {
	return cbdccommon.GetSlashingClient(n)
}

func (n *Network) GetDistrClient() distrtypes.QueryClient {
	return cbdccommon.GetDistrClient(n)
}

func (n *Network) GetPoaClient() poatypes.QueryClient {
	return cbdccommon.GetPoaClient(n)
}

// Not needed
func (n *Network) GetPreciseBankClient() precisebanktypes.QueryClient {
	return nil
}
