package keeper

import (
	"context"

	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/peersyst/cbdc-node/x/cbdc/types"
)

func (k msgServer) Mint(goCtx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(gov.ErrInvalidSigner, "expected %s got %s", k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.ExecuteMint(ctx, msg.Address, msg.Amount); err != nil {
		return nil, err
	}

	return &types.MsgMintResponse{}, nil
}
