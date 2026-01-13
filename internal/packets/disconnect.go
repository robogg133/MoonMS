package packets

import (
	datatypes "MoonMS/internal/datatypes/varTypes"
	"encoding/json"
)

type Text struct {
	Text any `json:"text"`
}

func DisconnectLogin(reason any) ([]byte, error) {

	packetID := datatypes.WriteVarInt(PACKET_LOGIN_DISCONNECT)

	var txt Text

	txt.Text = reason

	reasonPayload, err := json.Marshal(&txt)
	if err != nil {
		return nil, err
	}

	reasonLenght := datatypes.WriteVarInt(int32(len(reasonPayload)))

	totalLenght := len(packetID) + len(reasonLenght) + len(reasonPayload)

	lenght := datatypes.WriteVarInt(int32(totalLenght))

	disconnectPacket := make([]byte, totalLenght+len(lenght))

	offset := copy(disconnectPacket, lenght)
	offset += copy(disconnectPacket[offset:], packetID)
	offset += copy(disconnectPacket[offset:], reasonLenght)
	copy(disconnectPacket[offset:], reasonPayload)

	return disconnectPacket, nil
}
