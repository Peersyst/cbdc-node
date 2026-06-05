package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// Param store keys.
var (
	// KeyOwner is the param store key for the mint/burn owner address.
	KeyOwner = []byte("Owner")
	// KeyIssuancePaused is the param store key for the issuance pause switch.
	KeyIssuancePaused = []byte("IssuancePaused")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(owner string, issuancePaused bool) Params {
	return Params{Owner: owner, IssuancePaused: issuancePaused}
}

// DefaultParams returns a default set of parameters. The owner is intentionally
// empty (and therefore invalid): each chain must set a mint/burn owner in
// genesis, otherwise InitChain fails rather than booting with no minter.
func DefaultParams() Params {
	return NewParams("", false)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyOwner, &p.Owner, validateOwner),
		paramtypes.NewParamSetPair(KeyIssuancePaused, &p.IssuancePaused, validateIssuancePaused),
	}
}

// Validate validates the set of params
func (p *Params) Validate() error {
	if err := validateOwner(p.Owner); err != nil {
		return err
	}
	return validateIssuancePaused(p.IssuancePaused)
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

// validateIssuancePaused only checks the param is a bool; both values are valid.
func validateIssuancePaused(i interface{}) error {
	if _, ok := i.(bool); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
