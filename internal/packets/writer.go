package packets

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"

	"github.com/robogg133/KernelCraft/internal/datatypes"
)

type BufferWriter interface {
	io.Writer

	WriteByte(byte) error
	WriteString(string) (int, error)
	WriteTo(w io.Writer) (n int64, err error)

	Available() int
	Reset()

	Bytes() []byte
	Len() int
}

type Writer struct {
	buf *bytes.Buffer
}

func NewWriter() *Writer {

	return &Writer{buf: new(bytes.Buffer)}
}

func (w *Writer) Bytes() []byte {
	return w.buf.Bytes()
}

func (w *Writer) Len() int {
	return w.buf.Len()
}

func (w *Writer) WriteVarInt(n int32) error {
	v := datatypes.NewVarInt(n)
	_, err := w.buf.Write(v)
	return err
}

func (w *Writer) WriteString(s string) error {
	if err := w.WriteVarInt(int32(len(s))); err != nil {
		return err
	}
	_, err := w.buf.WriteString(s)
	return err
}

const TRUE_VALUE byte = 0x01
const FALSE_VALUE byte = 0x00

func (w *Writer) WriteBoolean(b bool) error {
	if b {
		return w.buf.WriteByte(TRUE_VALUE)
	}
	return w.buf.WriteByte(FALSE_VALUE)
}

func (w *Writer) Write(b []byte) error {
	_, err := w.buf.Write(b)
	return err
}

func (w *Writer) WritePrefixed(b []byte) error {
	if err := w.WriteVarInt(int32(len(b))); err != nil {
		return err
	}
	_, err := w.buf.Write(b)
	return err
}

func (w *Writer) WriteByte(n byte) error { return w.buf.WriteByte(n) }

func (w *Writer) WriteUnsignedByte(n uint8) error { return w.buf.WriteByte(n) }

func (w *Writer) WriteShort(n int16) error {
	buff := make([]byte, 2)
	binary.BigEndian.PutUint16(buff, uint16(n))
	_, err := w.buf.Write(buff)
	return err
}

func (w *Writer) WriteUnsignedShort(n uint16) error {
	buff := make([]byte, 2)
	binary.BigEndian.PutUint16(buff, n)
	_, err := w.buf.Write(buff)
	return err
}

func (w *Writer) WriteInt(n int32) error {
	buff := make([]byte, 4)
	binary.BigEndian.PutUint32(buff, uint32(n))
	_, err := w.buf.Write(buff)
	return err
}

func (w *Writer) WriteLong(n int64) error {
	buff := make([]byte, 8)
	binary.BigEndian.PutUint64(buff, uint64(n))
	_, err := w.buf.Write(buff)
	return err
}

func (w *Writer) WriteFloat(n float32) error {
	_, err := w.buf.Write(float32ToByte(n))
	return err
}

func (w *Writer) WriteDouble(n float64) error {
	_, err := w.buf.Write(float64ToByte(n))
	return err
}

//

func float64ToByte(f float64) []byte {
	bits := math.Float64bits(f)

	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, bits)
	return b
}

func float32ToByte(f float32) []byte {
	bits := math.Float32bits(f)

	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, bits)
	return b
}
