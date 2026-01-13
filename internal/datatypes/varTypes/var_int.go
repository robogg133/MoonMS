package datatypes

import (
	"fmt"
	"io"
)

type VarInt []byte // cannot be more than 5 bytes

func ReadVarInt(r io.ByteReader) (int32, int, error) {
	var value int32 = 0
	var position int = 0
	bytesRead := 0

	for {
		if bytesRead >= 5 {
			return 0, bytesRead, fmt.Errorf("varInt too big")
		}

		currentByte, err := r.ReadByte()
		if err != nil {
			if err == io.EOF && bytesRead > 0 {
				return 0, bytesRead, fmt.Errorf("malformed varInt")
			}
			return 0, bytesRead, err
		}

		bytesRead++

		value |= int32(currentByte&SEGMENT_BITS) << position

		if (currentByte & CONTINUE_BIT) == 0 {
			break
		}

		position += 7

		if position >= 32 {
			return 0, bytesRead, fmt.Errorf("to big varInt")
		}
	}

	return value, bytesRead, nil
}

func WriteVarInt(value int32) []byte {
	var result []byte

	for {
		temp := byte(value & int32(SEGMENT_BITS))

		value = int32(uint32(value) >> 7)

		if value != 0 {
			temp |= CONTINUE_BIT
		}

		result = append(result, temp)

		if value == 0 {
			break
		}
	}

	return result
}
