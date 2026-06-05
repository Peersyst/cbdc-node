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
		paused    bool
		expectErr bool
	}{
		{name: "empty owner is rejected", owner: "", expectErr: true},
		{name: "valid bech32 owner", owner: sample.AccAddress(), expectErr: false},
		{name: "malformed owner", owner: "not-an-address", expectErr: true},
		{name: "paused is valid", owner: sample.AccAddress(), paused: true, expectErr: false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			params := NewParams(tc.owner, tc.paused)
			err := params.Validate()
			if tc.expectErr {
				require.ErrorIs(t, err, ErrInvalidOwner)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
