package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/cbdc module sentinel errors
var (
	ErrInvalidAmount = sdkerrors.Register(ModuleName, 1, "amount must be positive")
	ErrInvalidDenom  = sdkerrors.Register(ModuleName, 2, "amount denom must match the configured CBDC denom")
	ErrUnauthorized  = sdkerrors.Register(ModuleName, 3, "unauthorized signer")
	ErrInvalidOwner  = sdkerrors.Register(ModuleName, 4, "invalid owner address")
)
