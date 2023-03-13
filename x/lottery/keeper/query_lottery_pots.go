package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"lottery/x/lottery/types"
)

func (k Keeper) LotteryPotsAll(goCtx context.Context, req *types.QueryAllLotteryPotsRequest) (*types.QueryAllLotteryPotsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var lotteryPotss []types.LotteryPots
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	lotteryPotsStore := prefix.NewStore(store, types.KeyPrefix(types.LotteryPotsKey))

	pageRes, err := query.Paginate(lotteryPotsStore, req.Pagination, func(key []byte, value []byte) error {
		var lotteryPots types.LotteryPots
		if err := k.cdc.Unmarshal(value, &lotteryPots); err != nil {
			return err
		}

		lotteryPotss = append(lotteryPotss, lotteryPots)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllLotteryPotsResponse{LotteryPots: lotteryPotss, Pagination: pageRes}, nil
}

func (k Keeper) LotteryPots(goCtx context.Context, req *types.QueryGetLotteryPotsRequest) (*types.QueryGetLotteryPotsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	lotteryPots, found := k.GetLotteryPots(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetLotteryPotsResponse{LotteryPots: lotteryPots}, nil
}
