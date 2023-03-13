package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"lottery/testutil/sample"
)

func TestMsgSendBuyTicket_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgSendBuyTicket
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgSendBuyTicket{
				Creator:          "invalid_address",
				Port:             "port",
				ChannelID:        "channel-0",
				TimeoutTimestamp: 100,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "invalid port",
			msg: MsgSendBuyTicket{
				Creator:          sample.AccAddress(),
				Port:             "",
				ChannelID:        "channel-0",
				TimeoutTimestamp: 100,
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "invalid channel",
			msg: MsgSendBuyTicket{
				Creator:          sample.AccAddress(),
				Port:             "port",
				ChannelID:        "",
				TimeoutTimestamp: 100,
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "invalid timeout",
			msg: MsgSendBuyTicket{
				Creator:          sample.AccAddress(),
				Port:             "port",
				ChannelID:        "channel-0",
				TimeoutTimestamp: 0,
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "valid message",
			msg: MsgSendBuyTicket{
				Creator:          sample.AccAddress(),
				Port:             "port",
				ChannelID:        "channel-0",
				TimeoutTimestamp: 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
