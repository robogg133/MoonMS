package packets

import (
	"MoonMS/internal/datatypes"
	"bytes"
)

const PACKET_CUSTOM_PAYLOAD int32 = 0x01

type CustomPacket struct {
	Identifier string
	Data       []byte
}

func (p *CustomPacket) Serialize() []byte {

	var buffer bytes.Buffer

	packetID := datatypes.NewVarInt(PACKET_CUSTOM_PAYLOAD)

	serializedIdentifier := datatypes.NewString(p.Identifier)

	totalLenght := len(packetID) + len(serializedIdentifier) + len(p.Data)

	buffer.Write(datatypes.NewVarInt(int32(totalLenght)))

	buffer.Write(packetID)

	buffer.Write(serializedIdentifier)
	buffer.Write(p.Data)

	return buffer.Bytes()
}
