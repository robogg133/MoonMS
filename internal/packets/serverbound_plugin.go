package packets

import "io"

const PACKET_SERVERBOUND_PLUGIN_MESSAGE int32 = 2

type ServerBoundPluginMessagePacket struct {
	Identifier string
	Data       []byte
}

func (*ServerBoundPluginMessagePacket) ID() int32 { return PACKET_SERVERBOUND_PLUGIN_MESSAGE }

func (s *ServerBoundPluginMessagePacket) Encode(w *Writer) error { return nil }

func (s *ServerBoundPluginMessagePacket) Decode(r *Reader) error {

	var err error
	s.Identifier, err = r.ReadString()
	if err != nil {
		return err
	}

	s.Data, err = io.ReadAll(r)
	if io.ErrUnexpectedEOF == err {
		return nil
	}
	return err
}
