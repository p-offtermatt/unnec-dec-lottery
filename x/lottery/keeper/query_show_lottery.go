package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"lottery/x/lottery/types"
)

func (k Keeper) ShowLottery(goCtx context.Context, req *types.QueryShowLotteryRequest) (*types.QueryShowLotteryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Process the query
	_ = ctx

	return &types.QueryShowLotteryResponse{}, nil
}
