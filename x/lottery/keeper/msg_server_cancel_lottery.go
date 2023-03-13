package keeper

import (
	"context"

	"lottery/x/lottery/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
)

func (k msgServer) CancelLottery(goCtx context.Context, msg *types.MsgCancelLottery) (*types.MsgCancelLotteryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	lottery, found := k.GetLottery(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "lottery with id %d not found", msg.Id)
	}

	if lottery.Creator != msg.Creator {
		return nil, sdkerrors.ErrUnauthorized.Wrap("cannot cancel lottery: not creator")
	}

	// hack with surely no adverse consequences:
	// lotteries are cancelled by setting their deadine to block 0, i.e. always expired
	lottery.Deadline = 0

	k.SetLottery(ctx, lottery)

	k.TransmitRefundLotteryPacket(ctx,
		types.RefundLotteryPacketData{Id: lottery.Id},
		msg.SourcePort,
		msg.SourceChannel,
		clienttypes.ZeroHeight(),
		msg.TimeoutTimestamp,
	)

	return &types.MsgCancelLotteryResponse{}, nil
}
