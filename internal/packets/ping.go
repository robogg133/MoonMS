package packets

import (
	vartypes "MoonMS/internal/datatypes/varTypes"
)

func SerializePong(num []byte) []byte {

	protocolID := vartypes.WriteVarInt(PACKET_PONG)
	packetLenght := vartypes.WriteVarInt(int32(len(protocolID) + len(num)))

	response := make([]byte, 10)

	offset := 0
	offset += copy(response, packetLenght)
	offset += copy(response[offset:], protocolID)
	copy(response[offset:], num)

	return response
}
