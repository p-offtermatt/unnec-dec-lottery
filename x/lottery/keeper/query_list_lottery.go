package keeper

import (
	"context"

	"lottery/x/lottery/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ListLottery(goCtx context.Context, req *types.QueryListLotteryRequest) (*types.QueryListLotteryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var lotteries []types.Lottery
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	lotteryStore := prefix.NewStore(store, types.KeyPrefix(types.LotteryKey))

	pageRes, err := query.Paginate(lotteryStore, req.Pagination, func(key []byte, value []byte) error {
		var lottery types.Lottery
		if err := k.cdc.Unmarshal(value, &lottery); err != nil {
			return err
		}

		lotteries = append(lotteries, lottery)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryListLotteryResponse{Lottery: lotteries, Pagination: pageRes}, nil
}
