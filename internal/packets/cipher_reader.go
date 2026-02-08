package packets

import (
	"crypto/cipher"
	"io"
)

type CipherReader struct {
	r      io.Reader
	stream cipher.Stream
}

func NewCipherReader(r io.Reader, stream cipher.Stream) *CipherReader {
	return &CipherReader{
		r:      r,
		stream: stream,
	}
}

func (cr *CipherReader) Read(p []byte) (int, error) {
	n, err := cr.r.Read(p)
	if n > 0 {
		cr.stream.XORKeyStream(p[:n], p[:n])
	}
	return n, err
}
