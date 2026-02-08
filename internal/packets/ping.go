package packets

const PACKET_PING_PONG int32 = 0x01

type PingPong struct {
	Bytes []byte
}

func (p *PingPong) ID() int32 { return PACKET_PING_PONG }

func (p *PingPong) Encode(w *Writer) error {

	if err := w.Write(p.Bytes); err != nil {
		return err
	}

	return nil
}

func (p *PingPong) Decode(r *Reader) error {

	buff := make([]byte, 8)
	if _, err := r.Read(buff); err != nil {
		return err
	}

	p.Bytes = buff

	return nil
}
