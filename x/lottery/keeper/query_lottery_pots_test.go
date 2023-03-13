package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "lottery/testutil/keeper"
	"lottery/testutil/nullify"
	"lottery/x/lottery/types"
)

func TestLotteryPotsQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.LotteryKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNLotteryPots(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetLotteryPotsRequest
		response *types.QueryGetLotteryPotsResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetLotteryPotsRequest{Id: msgs[0].Id},
			response: &types.QueryGetLotteryPotsResponse{LotteryPots: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetLotteryPotsRequest{Id: msgs[1].Id},
			response: &types.QueryGetLotteryPotsResponse{LotteryPots: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetLotteryPotsRequest{Id: uint64(len(msgs))},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.LotteryPots(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestLotteryPotsQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.LotteryKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNLotteryPots(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllLotteryPotsRequest {
		return &types.QueryAllLotteryPotsRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.LotteryPotsAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.LotteryPots), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.LotteryPots),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.LotteryPotsAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.LotteryPots), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.LotteryPots),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.LotteryPotsAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.LotteryPots),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.LotteryPotsAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
