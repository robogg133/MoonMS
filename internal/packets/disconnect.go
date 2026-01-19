package packets

import (
	"MoonMS/internal/datatypes"
	"encoding/json"
)

type Text struct {
	Text any `json:"text"`
}

func DisconnectLogin(reason any) ([]byte, error) {

	packetID := datatypes.NewVarInt(int32(PACKET_LOGIN_DISCONNECT))

	var txt Text

	txt.Text = reason

	reasonPayload, err := json.Marshal(&txt)
	if err != nil {
		return nil, err
	}

	reasonLenght := datatypes.NewVarInt(int32(len(reasonPayload)))
	var totalLenght int32 = int32(len(packetID) + len(reasonLenght) + len(reasonPayload))

	lenght := datatypes.NewVarInt(totalLenght)

	disconnectPacket := make([]byte, totalLenght+int32(len(lenght)))

	offset := copy(disconnectPacket, lenght)
	offset += copy(disconnectPacket[offset:], packetID)
	offset += copy(disconnectPacket[offset:], reasonLenght)
	copy(disconnectPacket[offset:], reasonPayload)

	return disconnectPacket, nil
}
