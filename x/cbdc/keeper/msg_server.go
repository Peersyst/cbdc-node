package keeper

import (
	"context"

	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peersyst/cbdc-node/x/cbdc/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) Mint(goCtx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	if k.authority != msg.Owner {
		return nil, errors.Wrapf(types.ErrUnauthorized, "expected %s got %s", k.authority, msg.Owner)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.executeMint(ctx, msg.Owner, msg.Address, msg.Amount); err != nil {
		return nil, err
	}

	return &types.MsgMintResponse{}, nil
}

func (k msgServer) Burn(goCtx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	if k.authority != msg.Owner {
		return nil, errors.Wrapf(types.ErrUnauthorized, "expected %s got %s", k.authority, msg.Owner)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.executeBurn(ctx, msg.Owner, msg.Address, msg.Amount); err != nil {
		return nil, err
	}

	return &types.MsgBurnResponse{}, nil
}
