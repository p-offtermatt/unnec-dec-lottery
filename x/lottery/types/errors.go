package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/lottery module sentinel errors
var (
	ErrSample                = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrInvalidPacketTimeout  = sdkerrors.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion        = sdkerrors.Register(ModuleName, 1501, "invalid version")
	ErrLotteryDeadlinePassed = sdkerrors.Register(ModuleName, 2, "lottery deadline has passed, so lottery is not active")
)
