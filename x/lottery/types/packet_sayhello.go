package types

// ValidateBasic is used for validating the packet
func (p SayhelloPacketData) ValidateBasic() error {

	// TODO: Validate the packet data

	return nil
}

// GetBytes is a helper for serialising
func (p SayhelloPacketData) GetBytes() ([]byte, error) {
	var modulePacket LotteryPacketData

	modulePacket.Packet = &LotteryPacketData_SayhelloPacket{&p}

	return modulePacket.Marshal()
}
