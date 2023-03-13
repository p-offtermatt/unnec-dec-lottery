package keeper

import (
	"context"

	"lottery/x/lottery/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
)

func (k msgServer) SendBuyTicket(goCtx context.Context, msg *types.MsgSendBuyTicket) (*types.MsgSendBuyTicketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creatorAddr := sdk.MustAccAddressFromBech32(msg.Creator)

	coins := sdk.NewCoins(msg.Price)

	// hold funds until we know whether the ticket could be bought
	k.bankKeeper.SendCoinsFromAccountToModule(ctx, creatorAddr, types.ModuleName, coins)

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
