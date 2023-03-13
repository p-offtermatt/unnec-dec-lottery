package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	"lottery/x/lottery/types"
)

func (k msgServer) SendSayhello(goCtx context.Context, msg *types.MsgSendSayhello) (*types.MsgSendSayhelloResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: logic before transmitting the packet

	// Construct the packet
	var packet types.SayhelloPacketData

	packet.Id = msg.Id

	// Transmit the packet
	_, err := k.TransmitSayhelloPacket(
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

	return &types.MsgSendSayhelloResponse{}, nil
}
