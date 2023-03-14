package types

// ValidateBasic is used for validating the packet
func (p WinnerPickedPacketData) ValidateBasic() error {

	// TODO: Validate the packet data

	return nil
}

// GetBytes is a helper for serialising
func (p WinnerPickedPacketData) GetBytes() ([]byte, error) {
	var modulePacket LotteryPacketData

	modulePacket.Packet = &LotteryPacketData_WinnerPickedPacket{&p}

	return modulePacket.Marshal()
}
