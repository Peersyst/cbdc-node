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

// DefaultParams returns a default set of parameters
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

// validateOwner accepts an empty owner (mint/burn stays disabled until gov sets
// one) or a well-formed bech32 address.
func validateOwner(i interface{}) error {
	owner, ok := i.(string)
	if !ok {
		return ErrInvalidOwner.Wrapf("invalid parameter type: %T", i)
	}
	if owner == "" {
		return nil
	}
	if _, err := sdk.AccAddressFromBech32(owner); err != nil {
		return ErrInvalidOwner.Wrapf("invalid owner address (%s): %s", owner, err)
	}
	return nil
}
