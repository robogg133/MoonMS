package packets

const PACKET_UPDATE_ENABLED_FEATURES int32 = 0x0C

type UpdateEnabledFeaturesPacket struct {
	Identifiers []string
}

func (*UpdateEnabledFeaturesPacket) ID() int32 { return PACKET_UPDATE_ENABLED_FEATURES }

func (s *UpdateEnabledFeaturesPacket) Encode(w *Writer) error {

	if err := w.WriteVarInt(int32(len(s.Identifiers))); err != nil {
		return err
	}

	for _, v := range s.Identifiers {
		if err := w.WriteString(v); err != nil {
			return err
		}
	}

	return nil
}

func (s *UpdateEnabledFeaturesPacket) Decode(r *Reader) error {

	length, err := r.ReadVarInt()
	if err != nil {
		return err
	}

	for range length {
		id, err := r.ReadString()
		if err != nil {
			return err
		}
		s.Identifiers = append(s.Identifiers, id)
	}

	return nil
}
