package keeper

import (
	"testing"

	"github.com/peersyst/cbdc-node/x/cbdc/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServer_UpdateParams(t *testing.T) {
	newOwner := "ethm1ef8ep2et20ja5s6r99tafld67kph72h7e0577u"

	tt := []struct {
		name        string
		authority   string
		params      types.Params
		expectedErr error
	}{
		{
			name:        "should fail - unauthorized authority",
			authority:   testOwner,
			params:      types.NewParams(newOwner, false),
			expectedErr: types.ErrUnauthorized,
		},
		{
			name:        "should fail - invalid owner",
			authority:   testGovAuthority,
			params:      types.NewParams("invalidowner", false),
			expectedErr: types.ErrInvalidOwner,
		},
		{
			name:        "should fail - empty owner",
			authority:   testGovAuthority,
			params:      types.NewParams("", false),
			expectedErr: types.ErrInvalidOwner,
		},
		{
			name:      "should pass - rotate owner",
			authority: testGovAuthority,
			params:    types.NewParams(newOwner, false),
		},
		{
			name:      "should pass - pause issuance",
			authority: testGovAuthority,
			params:    types.NewParams(newOwner, true),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cbdcKeeper, ctx := cbdcKeeperTestSetup(t)
			msgServer := NewMsgServerImpl(*cbdcKeeper)

			msg := &types.MsgUpdateParams{
				Authority: tc.authority,
				Params:    tc.params,
			}

			_, err := msgServer.UpdateParams(ctx, msg)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr.Error())
				return
			}

			require.NoError(t, err)
			got := cbdcKeeper.GetParams(ctx)
			require.Equal(t, tc.params.Owner, got.Owner)
			require.Equal(t, tc.params.IssuancePaused, got.IssuancePaused)
		})
	}
}
