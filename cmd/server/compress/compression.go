package compress

import (
	"MoonMS/internal/datatypes"
	"bytes"
	"compress/zlib"
	"fmt"
)

type CompressedPackage struct {
	PacketLenght datatypes.VarInt
}

func Compress(packet []byte, Threshold int32) ([]byte, error) {

	if Threshold == -1 {
		return packet, nil
	}

	_, readed, err := datatypes.ParseVarInt(bytes.NewReader(packet))
	if err != nil {
		return nil, fmt.Errorf("error parsing varInt: %v", err)
	}

	if int32(len(packet[readed:])) < Threshold {
		return tooShort(packet, readed), nil
	}

	return packPackage(packet, readed)

}

func tooShort(packet []byte, packetLenghtOffset int) []byte {
	var result bytes.Buffer

	result.Write(datatypes.NewVarInt(int32(len(packet[packetLenghtOffset:]) + 1))) // +1 for DataLenght byte

	result.WriteByte(0x00) // Data Lenght

	result.Write(packet[packetLenghtOffset:])

	return result.Bytes()
}

func packPackage(packet []byte, packetLenghtOffset int) ([]byte, error) {

	var compressedBuffer bytes.Buffer

	writerZlib := zlib.NewWriter(&compressedBuffer)

	dataLenght := datatypes.NewVarInt(int32(len(packet[packetLenghtOffset:])))

	_, err := writerZlib.Write(packet[packetLenghtOffset:])
	if err != nil {
		return nil, fmt.Errorf("error writing package to compressed buffer package: %v", err)
	}
	writerZlib.Close()

	compressedSize := len(compressedBuffer.Bytes())

	packetLenght := datatypes.NewVarInt(int32(compressedSize + len(dataLenght)))

	var resultBuffer bytes.Buffer

	resultBuffer.Write(packetLenght)
	resultBuffer.Write(dataLenght)
	resultBuffer.Write(compressedBuffer.Bytes())

	return resultBuffer.Bytes(), nil
}
