package packets

import (
	"MoonMS/internal/datatypes"
	"bytes"
)

func RecongnizePacket(b []byte) (int32, error) {
	_, offset, err := datatypes.ParseVarInt(bytes.NewReader(b))
	if err != nil {
		return 0, err
	}
	id, _, err := datatypes.ParseVarInt(bytes.NewReader(b[offset:]))
	if err != nil {
		return 0, err
	}

	return id, nil
}
