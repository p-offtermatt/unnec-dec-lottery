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

	return packetAck, nil
}

// OnAcknowledgementRefundLotteryPacket responds to the the success or failure of a packet
// acknowledgement written on the receiving chain.
func (k Keeper) OnAcknowledgementRefundLotteryPacket(ctx sdk.Context, packet channeltypes.Packet, data types.RefundLotteryPacketData, ack channeltypes.Acknowledgement) error {
	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:

		return sdkerrors.ErrLogic.Wrapf("could not cancel lottery: IBC packet was acknowledged with error %s", dispatchedAck.Error)
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck types.RefundLotteryPacketAck

		if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
			// The counter-party module doesn't implement the correct acknowledgment format
			return errors.New("cannot unmarshal acknowledgment")
		}

		lottery, found := k.GetLottery(ctx, data.Id)

		if !found {
			return sdkerrors.ErrNotFound.Wrapf("lottery with id %d not found", data.Id)
		}

		if lottery.Deadline < uint64(ctx.BlockHeight()) {
			return types.ErrLotteryDeadlinePassed.Wrap("deadline has passed")
		}

		// lotteries are cancelled by setting their deadine to block 0, i.e. always expired
		lottery.Deadline = 0

		k.SetLottery(ctx, lottery)

		return nil
	default:
		// The counter-party module doesn't implement the correct acknowledgment format
		return errors.New("invalid acknowledgment format")
	}
}

// OnTimeoutRefundLotteryPacket responds to the case where a packet has not been transmitted because of a timeout
func (k Keeper) OnTimeoutRefundLotteryPacket(ctx sdk.Context, packet channeltypes.Packet, data types.RefundLotteryPacketData) error {
	return sdkerrors.ErrTxTimeoutHeight.Wrapf("could not cancel lottery: packet timed out")
}
