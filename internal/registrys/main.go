package registrys

import (
	"bytes"

	"github.com/robogg133/KernelCraft/internal/datatypes"
	"github.com/robogg133/KernelCraft/internal/packets"
)

type RegistryData struct {
	DimensionIdentifier string
	NBTData             []byte
}

type Registry struct {
	AllData []RegistryData
}

func (reg *Registry) Serialize() []byte {
	var response bytes.Buffer

	var allDataBuffer bytes.Buffer

	var prefixedLength int32 = 0

	for _, data := range reg.AllData {
		allDataBuffer.Write(datatypes.NewString(data.DimensionIdentifier))
		allDataBuffer.Write(data.NBTData)
		prefixedLength++
	}

	packetId := datatypes.NewVarInt(packets.PACKET_REGISTRY_DATA)

	prefixedLenghtSerialized := datatypes.NewVarInt(prefixedLength)

	response.Write(datatypes.NewVarInt(int32(len(packetId) + len(prefixedLenghtSerialized) + allDataBuffer.Len())))

	response.Write(packetId)
	response.Write(prefixedLenghtSerialized)
	response.Write(allDataBuffer.Bytes())

	return response.Bytes()
}
