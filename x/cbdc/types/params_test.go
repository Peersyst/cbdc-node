package types

import (
	"testing"

	"github.com/peersyst/cbdc-node/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestParams_Validate(t *testing.T) {
	tt := []struct {
		name      string
		owner     string
		expectErr bool
	}{
		{name: "empty owner is allowed", owner: "", expectErr: false},
		{name: "valid bech32 owner", owner: sample.AccAddress(), expectErr: false},
		{name: "malformed owner", owner: "not-an-address", expectErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			params := NewParams(tc.owner)
			err := params.Validate()
			if tc.expectErr {
				require.ErrorIs(t, err, ErrInvalidOwner)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
