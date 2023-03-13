package types

// ValidateBasic is used for validating the packet
func (p BuyTicketPacketData) ValidateBasic() error {

	// TODO: Validate the packet data

	return nil
}

// GetBytes is a helper for serialising
func (p BuyTicketPacketData) GetBytes() ([]byte, error) {
	var modulePacket LotteryPacketData

	modulePacket.Packet = &LotteryPacketData_BuyTicketPacket{&p}

	return modulePacket.Marshal()
}
