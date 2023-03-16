package keeper

import (
	"context"
	"fmt"

	"lottery/x/lottery/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateLottery(goCtx context.Context, msg *types.MsgCreateLottery) (*types.MsgCreateLotteryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	fmt.Printf("lottery creation received %v", msg)

	// other fields will be initialized by keeper (id)
	// or left empty (users)
	lottery := types.Lottery{
		Deadline: uint64(ctx.BlockHeight()) + msg.Deadline,
		Creator:  msg.Creator,
	}

	k.AppendLottery(ctx, lottery)

	return &types.MsgCreateLotteryResponse{}, nil
}
