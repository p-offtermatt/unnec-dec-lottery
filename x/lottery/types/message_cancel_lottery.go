package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCancelLottery = "cancel_lottery"

var _ sdk.Msg = &MsgCancelLottery{}

func NewMsgCancelLottery(creator string, id uint64) *MsgCancelLottery {
	return &MsgCancelLottery{
		Creator: creator,
		Id:      id,
	}
}

func (msg *MsgCancelLottery) Route() string {
	return RouterKey
}

func (msg *MsgCancelLottery) Type() string {
	return TypeMsgCancelLottery
}

func (msg *MsgCancelLottery) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCancelLottery) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCancelLottery) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
