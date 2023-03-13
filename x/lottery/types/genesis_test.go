package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"lottery/x/lottery/types"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				PortId: types.PortID,
				LotteryPotsList: []types.LotteryPots{
					{
						Id: 0,
					},
					{
						Id: 1,
					},
				},
				LotteryPotsCount: 2,
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated lotteryPots",
			genState: &types.GenesisState{
				LotteryPotsList: []types.LotteryPots{
					{
						Id: 0,
					},
					{
						Id: 0,
					},
				},
			},
			valid: false,
		},
		{
			desc: "invalid lotteryPots count",
			genState: &types.GenesisState{
				LotteryPotsList: []types.LotteryPots{
					{
						Id: 1,
					},
				},
				LotteryPotsCount: 0,
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
