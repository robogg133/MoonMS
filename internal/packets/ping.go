package packets

import "MoonMS/internal/datatypes"

func SerializePong(num []byte) []byte {

	protocolID := datatypes.NewVarInt(int32(PACKET_PONG))
	packetLenght := datatypes.NewVarInt(int32(len(protocolID) + len(num)))

	response := make([]byte, 10)

	offset := 0
	offset += copy(response, packetLenght)
	offset += copy(response[offset:], protocolID)
	copy(response[offset:], num)

	return response
}
