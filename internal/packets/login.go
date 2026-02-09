package packets

import (
	"MoonMS/internal/datatypes"
	"bytes"
)

const PACKET_LOGIN_SUCCESS int32 = 0x02

type GameProfile struct {
	UUID     []byte
	Username string

	Name          string
	Value         string
	HaveSignature bool
	Signature     string
}

type LoginSuccessPacket struct {
	Profile GameProfile
}

func (l *LoginSuccessPacket) ID() int32 { return PACKET_LOGIN_SUCCESS }

func (l *LoginSuccessPacket) Encode(w *Writer) error {

	if err := w.Write(l.Profile.UUID); err != nil {
		return err
	}
	if err := w.WriteString(l.Profile.Username); err != nil {
		return err
	}

	if l.Profile.Name != "" && l.Profile.Value != "" {
		propreties := NewWriter()

		propreties.WriteString(l.Profile.Name)
		propreties.WriteString(l.Profile.Value)

		propreties.WriteBoolean(l.Profile.HaveSignature)
		if l.Profile.HaveSignature {
			propreties.WriteString(l.Profile.Signature)
		}

		w.WriteVarInt(1)
		w.Write(propreties.Bytes())

		return nil
	}

	w.WriteVarInt(0)

	return nil
}

func (l *LoginSuccessPacket) Decode(r *Reader) error { return nil }

func (pkg *LoginSuccessPacket) _Serialize() []byte {

	var response bytes.Buffer
	var proprieties bytes.Buffer

	if pkg.Profile.Name == "" && pkg.Profile.Value == "" && !pkg.Profile.HaveSignature {

		response.Write(datatypes.NewVarInt(int32(len(pkg.Profile.Username) + len(pkg.Profile.UUID) + 2))) // 2 is for packet id and Proprieties prefix
		response.Write(datatypes.NewVarInt(int32(PACKET_LOGIN_SUCCESS)))
		response.Write(pkg.Profile.UUID)
		//	response.Write(pkg.Profile.Username)
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
	//	response.Write(pkg.Profile.Username)

	response.Write(proprietiesPrefix)
	response.Write(proprieties.Bytes())
	proprieties.Reset()

	return response.Bytes()
}
