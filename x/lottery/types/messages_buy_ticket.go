package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSendBuyTicket = "send_buy_ticket"

var _ sdk.Msg = &MsgSendBuyTicket{}

func NewMsgSendBuyTicket(
	creator string,
	port string,
	channelID string,
	timeoutTimestamp uint64,
	lotteryId uint64,
	price sdk.Coin,
) *MsgSendBuyTicket {
	return &MsgSendBuyTicket{
		Creator:          creator,
		Port:             port,
		ChannelID:        channelID,
		TimeoutTimestamp: timeoutTimestamp,
		LotteryId:        lotteryId,
		Price:            price,
	}
}

func (msg *MsgSendBuyTicket) Route() string {
	return RouterKey
}

func (msg *MsgSendBuyTicket) Type() string {
	return TypeMsgSendBuyTicket
}

func (msg *MsgSendBuyTicket) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSendBuyTicket) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSendBuyTicket) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.Port == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid packet port")
	}
	if msg.ChannelID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid packet channel")
	}
	if msg.TimeoutTimestamp == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid packet timeout")
	}
	return nil
}
