package packets

const PACKET_SET_COMPRESSION int32 = 0x03

type CompressionStart struct {
	Threshould int32
}

func (c *CompressionStart) ID() int32 { return PACKET_SET_COMPRESSION }

func (c *CompressionStart) Encode(w *Writer) error {
	return w.WriteVarInt(c.Threshould)
}

func (c *CompressionStart) Decode(r *Reader) error { return nil }
