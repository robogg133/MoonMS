package packets

import "io"

const PACKET_PLUGIN_MESSAGE int32 = 2

type PluginMessagePacket struct {
	Identifier string
	Data       []byte
}

func (*PluginMessagePacket) ID() int32 { return PACKET_PLUGIN_MESSAGE }

func (s *PluginMessagePacket) Encode(w *Writer) error {

	if err := w.WriteString(s.Identifier); err != nil {
		return err
	}

	err := w.Write(s.Data)
	return err
}

func (s *PluginMessagePacket) Decode(r *Reader) error {

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
