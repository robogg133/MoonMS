package packets

const PACKET_REGISTRY_DATA int32 = 0x07

type RegistryDataPacket struct {
	Identifier string
	Tags       []RegistryEntry
}

type RegistryEntry struct {
	Identifier string
	NBTData    []byte // need to implement NBT type
}

func (*RegistryDataPacket) ID() int32 { return PACKET_REGISTRY_DATA }

func (s *RegistryDataPacket) Encode(w *Writer) error {

	if err := w.WriteString(s.Identifier); err != nil {
		return err
	}

	if err := w.WriteVarInt(int32(len(s.Tags))); err != nil {
		return err
	}

	for _, v := range s.Tags {
		if err := w.WriteString(v.Identifier); err != nil {
			return err
		}

		if err := w.Write(v.NBTData); err != nil {
			return err
		}
	}

	return nil
}

func (s *RegistryDataPacket) Decode(r *Reader) error {

	var err error
	s.Identifier, err = r.ReadString()
	if err != nil {
		return err
	}

	length, err := r.ReadVarInt()
	if err != nil {
		return err
	}

	for range length {
		var a RegistryEntry

		a.Identifier, err = r.ReadString()
		if err != nil {
			return err
		}

		a.NBTData, err = r.ReadPrefixed()
		if err != nil {
			return err
		}

		s.Tags = append(s.Tags, a)
	}

	return nil
}
