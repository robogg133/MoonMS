package packets

import (
	"fmt"
	"io"
)

const PACKET_COOKIE_REQUEST int32 = 0x05
const PACKET_COOKIE_RESPONSE int32 = 0x04

const COOKIE_MAX_LEN int = 5120

type CookieResponsePacket struct {
	Identifier string
	Data       []byte
}

func (*CookieResponsePacket) ID() int32 { return PACKET_COOKIE_RESPONSE }

func (s *CookieResponsePacket) Encode(w *Writer) error {
	if len(s.Data) > COOKIE_MAX_LEN {
		return fmt.Errorf("cookie max len overflow, max: %d, using: %d", COOKIE_MAX_LEN, len(s.Data))
	}

	if err := w.WriteString(s.Identifier); err != nil {
		return err
	}

	err := w.Write(s.Data)
	return err
}

func (s *CookieResponsePacket) Decode(r *Reader) error {

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

//

type CookieRequestPacket struct {
	Identifier string
}

func (*CookieRequestPacket) ID() int32 { return PACKET_COOKIE_REQUEST }

func (s *CookieRequestPacket) Encode(w *Writer) error {
	err := w.WriteString(s.Identifier)
	return err
}

func (s *CookieRequestPacket) Decode(r *Reader) error {

	var err error
	s.Identifier, err = r.ReadString()
	if err != nil {
		return err
	}
	return err
}
