package keeper

import (
	"context"

	"lottery/x/lottery/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateLottery(goCtx context.Context, msg *types.MsgCreateLottery) (*types.MsgCreateLotteryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// other fields will be initialized by keeper (id)
	// or left empty (users)
	lottery := types.Lottery{
		Deadline: msg.Deadline,
	}

	k.AppendLottery(ctx, lottery)

	return &types.MsgCreateLotteryResponse{}, nil
}
