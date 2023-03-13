package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCreateLottery = "create_lottery"

var _ sdk.Msg = &MsgCreateLottery{}

func NewMsgCreateLottery(creator string, deadline uint64) *MsgCreateLottery {
	return &MsgCreateLottery{
		Creator:  creator,
		Deadline: deadline,
	}
}

func (msg *MsgCreateLottery) Route() string {
	return RouterKey
}

func (msg *MsgCreateLottery) Type() string {
	return TypeMsgCreateLottery
}

func (msg *MsgCreateLottery) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateLottery) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateLottery) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
