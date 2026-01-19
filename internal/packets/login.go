package packets

import (
	"MoonMS/internal/datatypes"
	"bytes"
)

type GameProfile struct {
	UUID     []byte
	Username datatypes.String

	Name          string
	Value         string
	HaveSignature bool
	Signature     string
}

type LoginSuccessPacket struct {
	Profile GameProfile
}

func (pkg *LoginSuccessPacket) Serialize() []byte {

	var response bytes.Buffer
	var proprieties bytes.Buffer

	if pkg.Profile.Name == "" && pkg.Profile.Value == "" && !pkg.Profile.HaveSignature {

		response.Write(datatypes.NewVarInt(int32(len(pkg.Profile.Username) + len(pkg.Profile.UUID) + 2))) // 2 is for packet id and Proprieties prefix
		response.Write(datatypes.NewVarInt(int32(PACKET_LOGIN_SUCCESS)))
		response.Write(pkg.Profile.UUID)
		response.Write(pkg.Profile.Username)
		response.Write(datatypes.NewVarInt(0))

		return response.Bytes()
	}

	proprieties.Write(datatypes.NewString(pkg.Profile.Name))
	proprieties.Write(datatypes.NewString(pkg.Profile.Value))

	proprieties.WriteByte(byte(datatypes.NewBoolean(pkg.Profile.HaveSignature)))
	if pkg.Profile.HaveSignature {
		proprieties.Write(datatypes.NewString(pkg.Profile.Signature))
	}

	proprietiesPrefix := datatypes.NewVarInt(int32(proprieties.Len()))

	totalLenght := len(datatypes.NewVarInt(int32(PACKET_LOGIN_SUCCESS))) + len(pkg.Profile.UUID) + len(pkg.Profile.Username) + len(proprietiesPrefix) + proprieties.Len()

	response.Write(datatypes.NewVarInt(int32(totalLenght)))
	response.Write(datatypes.NewVarInt(int32(PACKET_LOGIN_SUCCESS)))

	response.Write(pkg.Profile.UUID)
	response.Write(pkg.Profile.Username)

	response.Write(proprietiesPrefix)
	response.Write(proprieties.Bytes())
	proprieties.Reset()

	return response.Bytes()
}
