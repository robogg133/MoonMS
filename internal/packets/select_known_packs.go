package packets

const PACKET_SELECT_KNOWN_PACKS int32 = 0x0E

type SelectKnownPacksPacket struct {
	Packs []Pack
}

type Pack struct {
	Namespace string
	Pathname  string
	Version   string
}

func (*SelectKnownPacksPacket) ID() int32 { return PACKET_SELECT_KNOWN_PACKS }

func (s *SelectKnownPacksPacket) Encode(w *Writer) error {

	if err := w.WriteVarInt(int32(len(s.Packs))); err != nil {
		return err
	}

	for _, v := range s.Packs {
		if err := w.WriteString(v.Namespace); err != nil {
			return err
		}

		if err := w.WriteString(v.Pathname); err != nil {
			return err
		}

		if err := w.WriteString(v.Version); err != nil {
			return err
		}

	}
	return nil
}

func (s *SelectKnownPacksPacket) Decode(r *Reader) error {

	numPacks, err := r.ReadVarInt()
	if err != nil {
		return err
	}

	for range numPacks {
		var a Pack
		a.Namespace, err = r.ReadString()
		if err != nil {
			return err
		}

		a.Pathname, err = r.ReadString()
		if err != nil {
			return err
		}

		a.Version, err = r.ReadString()
		if err != nil {
			return err
		}
		s.Packs = append(s.Packs, a)
	}

	return nil
}
