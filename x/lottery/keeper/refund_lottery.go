package keeper

import (
	"errors"
	"time"

	"lottery/x/lottery/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

// TransmitRefundLotteryPacket transmits the packet over IBC with the specified source port and source channel
func (k Keeper) TransmitRefundLotteryPacket(
	ctx sdk.Context,
	packetData types.RefundLotteryPacketData,
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

// OnRecvRefundLotteryPacket processes packet reception
func (k Keeper) OnRecvRefundLotteryPacket(ctx sdk.Context, packet channeltypes.Packet, data types.RefundLotteryPacketData) (packetAck types.RefundLotteryPacketAck, err error) {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return packetAck, err
	}

	// TODO: packet reception logic

	return packetAck, nil
}

// OnAcknowledgementRefundLotteryPacket responds to the the success or failure of a packet
// acknowledgement written on the receiving chain.
func (k Keeper) OnAcknowledgementRefundLotteryPacket(ctx sdk.Context, packet channeltypes.Packet, data types.RefundLotteryPacketData, ack channeltypes.Acknowledgement) error {
	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:

		// resend the packet
		k.TransmitRefundLotteryPacket(ctx, data, packet.DestinationPort, packet.DestinationChannel, clienttypes.ZeroHeight(), DefaultRelativePacketTimeoutTimestamp)

		return nil
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck types.RefundLotteryPacketAck

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

// OnTimeoutRefundLotteryPacket responds to the case where a packet has not been transmitted because of a timeout
func (k Keeper) OnTimeoutRefundLotteryPacket(ctx sdk.Context, packet channeltypes.Packet, data types.RefundLotteryPacketData) error {

	// just resend it
	k.TransmitRefundLotteryPacket(ctx, data, packet.DestinationPort, packet.DestinationChannel, clienttypes.ZeroHeight(), DefaultRelativePacketTimeoutTimestamp)

	return nil
}
