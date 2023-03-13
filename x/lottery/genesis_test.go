package lottery_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "lottery/testutil/keeper"
	"lottery/testutil/nullify"
	"lottery/x/lottery"
	"lottery/x/lottery/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
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
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.LotteryKeeper(t)
	lottery.InitGenesis(ctx, *k, genesisState)
	got := lottery.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.PortId, got.PortId)

	require.ElementsMatch(t, genesisState.LotteryPotsList, got.LotteryPotsList)
	require.Equal(t, genesisState.LotteryPotsCount, got.LotteryPotsCount)
	// this line is used by starport scaffolding # genesis/test/assert
}
