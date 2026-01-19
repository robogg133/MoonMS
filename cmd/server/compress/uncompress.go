package compress

import (
	"MoonMS/internal/datatypes"
	"bytes"
	"compress/zlib"

	"io"
)

func Uncompress(buffer []byte) ([]byte, error) {

	_, packetLenghtLenght, err := datatypes.ParseVarInt(bytes.NewBuffer(buffer))
	if err != nil {
		return nil, err
	}

	dataLenght, dataLenghtLenght, err := datatypes.ParseVarInt(bytes.NewBuffer(buffer[packetLenghtLenght:]))
	if err != nil {
		return nil, err
	}

	if dataLenght == 0 {
		var responseBuffer bytes.Buffer

		responseBuffer.Write(datatypes.NewVarInt(int32(len(buffer[dataLenghtLenght+packetLenghtLenght:]))))
		responseBuffer.Write(buffer[dataLenghtLenght+packetLenghtLenght:])

		return responseBuffer.Bytes(), nil
	}

	reader, err := zlib.NewReader(bytes.NewBuffer(buffer[dataLenghtLenght+packetLenghtLenght:]))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var responseBuffer bytes.Buffer

	responseBuffer.Write(datatypes.NewVarInt(dataLenght))

	decompressed, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	responseBuffer.Write(decompressed)

	return responseBuffer.Bytes(), nil
}
