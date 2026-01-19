package datatypes

import "bytes"

type String []byte

func NewString(value string) String {

	byteString := []byte(value)
	byteStringLen := int32(len(byteString))

	stringLenght := NewVarInt(byteStringLen)

	var payload bytes.Buffer

	payload.Write(stringLenght)
	payload.Write(byteString)

	return String(payload.Bytes())
}
