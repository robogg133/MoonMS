package packets

import (
	"MoonMS/internal/datatypes"
	"bytes"
	"encoding/binary"
	"io"
	"math"
)

type Reader struct {
	r io.Reader
}

func NewReader(b []byte) *Reader {
	return &Reader{r: bytes.NewReader(b)}
}
func NewReaderFromReader(r io.Reader) *Reader {
	return &Reader{r: r}
}

func (r *Reader) ReadVarInt() (int32, error) {
	n, _, err := datatypes.ParseVarInt(r)
	return n, err
}

func (r *Reader) ReadString() (string, error) {
	l, err := r.ReadVarInt()
	if err != nil {
		return "", err
	}

	b := make([]byte, l)
	io.ReadFull(r.r, b)
	return string(b), nil
}

func (r *Reader) ReadBoolean() (bool, error) {

	b, err := r.ReadByte()
	var v bool

	if b == FALSE_VALUE {
		v = false
	} else {
		v = true
	}

	return v, err
}

func (r *Reader) Read(b []byte) (int, error) { return io.ReadFull(r.r, b) }
func (r *Reader) ReadPrefixed() ([]byte, error) {
	n, err := r.ReadVarInt()
	if err != nil {
		return nil, err
	}

	buff := make([]byte, n)
	_, err = io.ReadFull(r.r, buff)
	return buff, err
}
func (r *Reader) ReadByte() (byte, error) {
	buff := make([]byte, 1)
	_, err := io.ReadFull(r.r, buff)
	return buff[0], err
}
func (r *Reader) ReadUnsignedByte() (uint8, error) {
	buff := make([]byte, 1)
	_, err := io.ReadFull(r.r, buff)
	return buff[0], err
}

func (r *Reader) ReadShort() int16 {
	var v uint16
	binary.Read(r.r, binary.BigEndian, &v)
	return int16(v)
}

func (r *Reader) ReadUnsignedShort() uint16 {
	var v uint16
	binary.Read(r.r, binary.BigEndian, &v)
	return v
}

func (r *Reader) ReadInt() (int32, error) {
	var v uint32
	err := binary.Read(r.r, binary.BigEndian, &v)
	return int32(v), err
}

func (r *Reader) ReadLong() (int64, error) {
	var v uint64
	err := binary.Read(r.r, binary.BigEndian, &v)
	return int64(v), err
}

func (r *Reader) ReadFloat() (float32, error) {
	buff := make([]byte, 4)
	_, err := io.ReadFull(r.r, buff)
	return bit32ToFloat32(buff), err
}

func (r *Reader) ReadDouble() (float64, error) {
	buff := make([]byte, 8)
	_, err := io.ReadFull(r.r, buff)
	return bit64ToFloat64(buff), err
}

//

func bit32ToFloat32(b []byte) float32 {
	bits := binary.BigEndian.Uint32(b)

	return math.Float32frombits(bits)
}

func bit64ToFloat64(b []byte) float64 {
	bits := binary.BigEndian.Uint64(b)

	return math.Float64frombits(bits)
}
