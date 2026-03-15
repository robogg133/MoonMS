package packets

import "encoding/json"

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

const PACKET_LOGIN_ACKNOWLEDGED int32 = 3

type LoginAcknowledgedPacket struct{}

func (l *LoginAcknowledgedPacket) ID() int32 { return PACKET_LOGIN_ACKNOWLEDGED }

func (l *LoginAcknowledgedPacket) Encode(w *Writer) error { return nil }

func (l *LoginAcknowledgedPacket) Decode(r *Reader) error { return nil }

const PACKET_LOGIN_DISCONNECT int32 = 0x00

type LoginDisconnectPacket struct {
	Reason json.RawMessage
}

func (l *LoginDisconnectPacket) ID() int32 { return PACKET_LOGIN_DISCONNECT }

func (l *LoginDisconnectPacket) Encode(w *Writer) error {

	b, err := l.Reason.MarshalJSON()
	if err != nil {
		return err
	}

	return w.WritePrefixed(b)
}

func (l *LoginDisconnectPacket) Decode(r *Reader) error {

	b, err := r.ReadPrefixed()
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &l.Reason)
}
