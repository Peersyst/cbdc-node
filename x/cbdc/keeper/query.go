package keeper

import (
	"github.com/peersyst/cbdc-node/x/cbdc/types"
)

var _ types.QueryServer = Keeper{}
