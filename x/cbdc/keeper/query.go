package keeper

import (
	"github.com/peersyst/cbdc-node/x/cbdc/types"
)

var _ types.QueryServer = Querier{}

// Querier is the gRPC query server for the cbdc module. It holds the keeper as
// an unexported field (rather than embedding it) so the read API stays
// decoupled from the keeper's state-mutating methods.
type Querier struct {
	k Keeper
}

// NewQuerier returns a Querier for the provided Keeper.
func NewQuerier(k Keeper) Querier {
	return Querier{k: k}
}
