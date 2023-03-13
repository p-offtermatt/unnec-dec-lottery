package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "lottery/testutil/keeper"
	"lottery/testutil/nullify"
	"lottery/x/lottery/keeper"
	"lottery/x/lottery/types"
)

func createNLotteryPots(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.LotteryPots {
	items := make([]types.LotteryPots, n)
	for i := range items {
		items[i].Id = keeper.AppendLotteryPots(ctx, items[i])
	}
	return items
}

func TestLotteryPotsGet(t *testing.T) {
	keeper, ctx := keepertest.LotteryKeeper(t)
	items := createNLotteryPots(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetLotteryPots(ctx, item.Id)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestLotteryPotsRemove(t *testing.T) {
	keeper, ctx := keepertest.LotteryKeeper(t)
	items := createNLotteryPots(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveLotteryPots(ctx, item.Id)
		_, found := keeper.GetLotteryPots(ctx, item.Id)
		require.False(t, found)
	}
}

func TestLotteryPotsGetAll(t *testing.T) {
	keeper, ctx := keepertest.LotteryKeeper(t)
	items := createNLotteryPots(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllLotteryPots(ctx)),
	)
}

func TestLotteryPotsCount(t *testing.T) {
	keeper, ctx := keepertest.LotteryKeeper(t)
	items := createNLotteryPots(keeper, ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, keeper.GetLotteryPotsCount(ctx))
}
