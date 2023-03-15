package keeper

import (
	"context"

	"lottery/x/lottery/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
)

func (k Keeper) GetPotForLottery(ctx sdk.Context, id uint64) (*types.LotteryPots, bool) {
	for _, pot := range k.GetAllLotteryPots(ctx) {
		if pot.Id == id {
			return &pot, true
		}
	}
	return nil, false
}

func (k msgServer) SendBuyTicket(goCtx context.Context, msg *types.MsgSendBuyTicket) (*types.MsgSendBuyTicketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creatorAddr := sdk.MustAccAddressFromBech32(msg.Creator)

	coins := sdk.NewCoins(msg.Price)

	// hold funds until we know whether the ticket could be bought
	k.bankKeeper.SendCoinsFromAccountToModule(ctx, creatorAddr, types.ModuleName, coins)

	// add funds to the pool
	pot, found := k.GetPotForLottery(ctx, msg.LotteryId)
	if !found {
		newPot := types.LotteryPots{
			Id:     msg.LotteryId,
			Amount: msg.Price.Amount.Uint64(),
		}
		k.AppendLotteryPots(ctx, newPot)
	} else {
		pot.Amount += msg.Price.Amount.Uint64()
		k.SetLotteryPots(ctx, *pot)
	}

	// Construct the packet
	var packet types.BuyTicketPacketData

	packet.LotteryId = msg.LotteryId
	packet.Price = msg.Price
	packet.Creator = msg.Creator

	// Transmit the packet
	_, err := k.TransmitBuyTicketPacket(
		ctx,
		packet,
		msg.Port,
		msg.ChannelID,
		clienttypes.ZeroHeight(),
		msg.TimeoutTimestamp,
	)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendBuyTicketResponse{}, nil
}
