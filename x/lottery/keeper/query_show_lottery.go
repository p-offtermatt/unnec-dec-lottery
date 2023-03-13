package keeper

import (
	"context"

	"lottery/x/lottery/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ShowLottery(goCtx context.Context, req *types.QueryShowLotteryRequest) (*types.QueryShowLotteryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	lottery, found := k.GetLottery(ctx, req.Id)

	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "lottery with id %d not found", req.Id)
	}

	return &types.QueryShowLotteryResponse{Lottery: lottery}, nil
}
