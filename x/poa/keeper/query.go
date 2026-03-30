package keeper

import (
	"github.com/peersyst/cbdc-node/x/poa/types"
)

var _ types.QueryServer = Keeper{}
