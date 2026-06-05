package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// KeyOwner is the param store key for the mint/burn owner address.
var KeyOwner = []byte("Owner")

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(owner string) Params {
	return Params{Owner: owner}
}

// DefaultParams returns a default set of parameters. The owner is intentionally
// empty (and therefore invalid): each chain must set a mint/burn owner in
// genesis, otherwise InitChain fails rather than booting with no minter.
func DefaultParams() Params {
	return NewParams("")
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyOwner, &p.Owner, validateOwner),
	}
}

// Validate validates the set of params
func (p *Params) Validate() error {
	return validateOwner(p.Owner)
}

// validateOwner requires a non-empty, well-formed bech32 address. The owner is
// the mint/burn signer, so an unset owner is rejected: a genesis that forgets to
// set it fails InitChain instead of silently coming up with no minter.
func validateOwner(i interface{}) error {
	owner, ok := i.(string)
	if !ok {
		return ErrInvalidOwner.Wrapf("invalid parameter type: %T", i)
	}
	if owner == "" {
		return ErrInvalidOwner.Wrap("owner must not be empty")
	}
	if _, err := sdk.AccAddressFromBech32(owner); err != nil {
		return ErrInvalidOwner.Wrapf("invalid owner address (%s): %s", owner, err)
	}
	return nil
}
