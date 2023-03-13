package types

// IBC events
const (
	EventTypeTimeout             = "timeout"
	EventTypeRefundLotteryPacket = "refundLottery_packet"
	EventTypeSayhelloPacket      = "sayhello_packet"
	EventTypeBuyTicketPacket     = "buyTicket_packet"
	// this line is used by starport scaffolding # ibc/packet/event

	AttributeKeyAckSuccess = "success"
	AttributeKeyAck        = "acknowledgement"
	AttributeKeyAckError   = "error"
)
