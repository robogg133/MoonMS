package plugin_wasmi

import (
	"MoonMS/internal/server"
	"context"

	"github.com/tetratelabs/wazero/api"
)

func (data *ServerInfo) GetMaxPlayers(ctx context.Context, m api.Module) uint32 {
	return uint32(data.MaxPlayers)
}

func SetMaxPlayers(ctx context.Context, m api.Module, v uint32) {
	server.GetServerData().MaxPlayers = uint(v)
}

func (data *ServerInfo) GetServerThreshold(ctx context.Context, m api.Module) int32 {
	return data.Threshold
}

func (data *ServerInfo) GetServerMotd(ctx context.Context, m api.Module) uint64 {
	content := []byte(data.Motd)

	mem := m.Memory()

	alloc := m.ExportedFunction("alloc")

	r, err := alloc.Call(ctx, uint64(len(content)))
	if err != nil {
		server.LogError("Can not call alloc function")
		return 0
	}

	ptr := uint32(r[0])

	mem.Write(ptr, content)

	return (uint64(ptr) << 32) | (uint64(len(content)) & 0xffffffff)
}

func SetServerMotd(ctx context.Context, m api.Module, ptr, lenght uint32) {
	mem := m.Memory()

	str, _ := mem.Read(ptr, lenght)

	server.GetServerData().Motd = string(str)
}
