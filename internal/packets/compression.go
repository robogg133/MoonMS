package packets

import (
	"MoonMS/internal/datatypes"
)

type CompressionStart struct {
	Threshould int32
}

func (cs *CompressionStart) Serialize() []byte {
	packet_id := datatypes.NewVarInt(int32(PACKET_SET_COMPRESSION))

	totalLenght := len(packet_id)

	thresold := datatypes.NewVarInt(int32(cs.Threshould))
	totalLenght += len(thresold)

	lenght := datatypes.NewVarInt(int32(totalLenght))

	packet := make([]byte, totalLenght+len(lenght))

	offest := copy(packet, lenght)
	offest += copy(packet[offest:], packet_id)
	offest += copy(packet[offest:], thresold)

	return packet
}
