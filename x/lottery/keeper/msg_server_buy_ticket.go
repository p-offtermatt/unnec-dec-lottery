package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	"lottery/x/lottery/types"
)

func (k msgServer) SendBuyTicket(goCtx context.Context, msg *types.MsgSendBuyTicket) (*types.MsgSendBuyTicketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: logic before transmitting the packet

	// Construct the packet
	var packet types.BuyTicketPacketData

	packet.LotteryId = msg.LotteryId
	packet.Price = msg.Price

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
