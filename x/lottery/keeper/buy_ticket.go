package keeper

import (
	"errors"

	"lottery/x/lottery/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
)

// TransmitBuyTicketPacket transmits the packet over IBC with the specified source port and source channel
func (k Keeper) TransmitBuyTicketPacket(
	ctx sdk.Context,
	packetData types.BuyTicketPacketData,
	sourcePort,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
) (uint64, error) {
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return 0, sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	packetBytes, err := packetData.GetBytes()
	if err != nil {
		return 0, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, "cannot marshal the packet: %w", err)
	}

	return k.channelKeeper.SendPacket(ctx, channelCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, packetBytes)
}

// OnRecvBuyTicketPacket processes packet reception
func (k Keeper) OnRecvBuyTicketPacket(ctx sdk.Context, packet channeltypes.Packet, data types.BuyTicketPacketData) (packetAck types.BuyTicketPacketAck, err error) {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return packetAck, err
	}

	lottery, found := k.GetLottery(ctx, data.LotteryId)
	if !found {
		return packetAck, sdkerrors.ErrNotFound.Wrapf("lottery with id %d not found", data.LotteryId)
	}

	if ctx.BlockHeight() > int64(lottery.Deadline) {
		return packetAck, types.ErrLotteryDeadlinePassed.Wrap("deadline to enter lottery has passed")
	}

	lottery.Users = append(lottery.Users, data.Creator)
	k.SetLottery(ctx, lottery)

	return packetAck, nil
}

func (k Keeper) RefundUser(ctx sdk.Context, data types.BuyTicketPacketData) error {
	// refund user
	creatorAddr := sdk.MustAccAddressFromBech32(data.Creator)
	coins := sdk.NewCoins(data.Price)

	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, coins)
	if err != nil {
		panic(err)
	}

	return nil
}

// OnAcknowledgementBuyTicketPacket responds to the the success or failure of a packet
// acknowledgement written on the receiving chain.
func (k Keeper) OnAcknowledgementBuyTicketPacket(ctx sdk.Context, packet channeltypes.Packet, data types.BuyTicketPacketData, ack channeltypes.Acknowledgement) error {
	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:

		// failed acknowledgement logic

		// refund user
		return k.RefundUser(ctx, data)
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck types.BuyTicketPacketAck

		if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
			// The counter-party module doesn't implement the correct acknowledgment format
			return errors.New("cannot unmarshal acknowledgment")
		}

		// TODO: successful acknowledgement logic

		return nil
	default:
		// The counter-party module doesn't implement the correct acknowledgment format
		return errors.New("invalid acknowledgment format")
	}
}

// OnTimeoutBuyTicketPacket responds to the case where a packet has not been transmitted because of a timeout
func (k Keeper) OnTimeoutBuyTicketPacket(ctx sdk.Context, packet channeltypes.Packet, data types.BuyTicketPacketData) error {

	return k.RefundUser(ctx, data)
}
