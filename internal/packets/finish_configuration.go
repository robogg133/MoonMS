package packets

const PACKET_FINISH_CONFIGURATION int32 = 0x03

type FinishConfigurationPacket struct{}

func (*FinishConfigurationPacket) ID() int32 { return PACKET_FINISH_CONFIGURATION }

func (*FinishConfigurationPacket) Encode(w *Writer) error { return nil }

func (*FinishConfigurationPacket) Decode(r *Reader) error { return nil }
