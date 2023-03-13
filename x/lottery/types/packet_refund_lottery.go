package types

// ValidateBasic is used for validating the packet
func (p RefundLotteryPacketData) ValidateBasic() error {

	// TODO: Validate the packet data

	return nil
}

// GetBytes is a helper for serialising
func (p RefundLotteryPacketData) GetBytes() ([]byte, error) {
	var modulePacket LotteryPacketData

	modulePacket.Packet = &LotteryPacketData_RefundLotteryPacket{&p}

	return modulePacket.Marshal()
}
