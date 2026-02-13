package packets

import (
	"MoonMS/internal/datatypes"
	"bytes"
	"compress/zlib"
	"crypto/cipher"
	"fmt"
	"io"
)

const (
	SERVER_BOUND uint8 = 0
	CLIENT_BOUND uint8 = 1
)

type Packet interface {
	ID() int32

	Encode(w *Writer) error
	Decode(r *Reader) error
}

type connStatus struct {
	Threshold int
}

var packetRegistry = map[int32]func() Packet{}

func RegisterPacket(id int32, fn func() Packet) {
	packetRegistry[id] = fn
}

func Init() {
	RegisterPacket(PACKET_PING_PONG, func() Packet {
		return &PingPong{}
	})
	RegisterPacket(PACKET_HANDSHAKE, func() Packet {
		return &Handshake{}
	})

}

func MarshalPacket(p Packet, encryptionKey cipher.Stream, t int) ([]byte, error) {
	body := NewWriter()

	body.WriteVarInt(p.ID())
	if err := p.Encode(body); err != nil {
		return nil, err
	}

	out := NewWriter()
	if t > -1 {

		if body.Len() < t {
			if err := out.WriteVarInt(int32(body.Len()) + 1); err != nil { // +1 For data Length
				return nil, err
			}
			if err := out.WriteVarInt(0); err != nil {
				return nil, err
			}

			out.buf.Write(body.Bytes())
		} else {
			dataLen := datatypes.NewVarInt(int32(body.Len()))

			var compressedData bytes.Buffer

			w := zlib.NewWriter(&compressedData)
			if _, err := w.Write(body.Bytes()); err != nil {
				return nil, err
			}
			w.Close()
			w = nil
			body.buf.Reset()
			body = nil

			/* debug
			r, err := zlib.NewReader(&compressedData)
			if err != nil {
				return nil, err
			}

			b, err := io.ReadAll(r)
			if err != nil {
				return nil, err
			}

			fmt.Println(hex.Dump(b))
			*/

			if err := out.WriteVarInt(int32(len(dataLen) + compressedData.Len())); err != nil {
				return nil, err
			}
			if _, err := out.buf.Write(dataLen); err != nil {
				return nil, err
			}
			out.buf.Write(compressedData.Bytes())
		}
	} else {
		out.WriteVarInt(int32(body.Len()))
		out.buf.Write(body.Bytes())
	}

	if encryptionKey != nil {
		buff := make([]byte, out.Len())
		encryptionKey.XORKeyStream(buff, out.Bytes())
		return buff, nil
	}

	return out.Bytes(), nil
}

func UnmarshalPacket(r *Reader, t int) (Packet, error) {

	data, err := r.ReadPrefixed()
	if err != nil {
		return nil, err
	}
	afterLen := NewReader(data)

	if t > -1 {
		dataLen, err := afterLen.ReadVarInt()
		if err != nil {
			return nil, err
		}
		if dataLen != 0 {
			reader, err := zlib.NewReader(r.r)
			if err != nil {
				return nil, err
			}

			b, err := io.ReadAll(reader)
			if err != nil {
				return nil, err
			}
			afterLen = NewReader(b)
		}
	}

	packetID, err := afterLen.ReadVarInt()
	if err != nil {
		return nil, err
	}

	fn, ok := packetRegistry[packetID]
	if !ok {
		return nil, fmt.Errorf("unknown packet id: %x", packetID)
	}

	pkt := fn()
	if err := pkt.Decode(afterLen); err != nil {
		return nil, err
	}

	return pkt, nil
}
