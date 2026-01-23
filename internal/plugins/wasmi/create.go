package plugin_wasmi

import (
	"MoonMS/internal/server"
	"context"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type ServerInfo server.ServerData

type Runtime struct {
	R             wazero.Runtime
	Ctx           context.Context
	ModuleConfigs wazero.ModuleConfig
}

func CreateRuntime(pluginName string) Runtime {
	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	build := r.NewHostModuleBuilder("env")

	serverData := ServerInfo(*server.GetServerData())

	build.NewFunctionBuilder().
		WithFunc(serverData.GetMaxPlayers).
		Export("get_server_max_players")

	build.NewFunctionBuilder().
		WithFunc(SetMaxPlayers).
		Export("set_server_max_players")

	build.NewFunctionBuilder().
		WithFunc(serverData.GetServerThreshold).
		Export("get_server_threshold")

	build.NewFunctionBuilder().
		WithFunc(serverData.GetServerMotd).
		Export("get_server_motd")

	build.NewFunctionBuilder().
		WithFunc(SetServerMotd).
		Export("set_server_motd")

	build.Instantiate(ctx)

	writer := server.GetLogPluginWriter(os.Stdout, pluginName)
	writer2 := server.GetLogPluginWriter(os.Stderr, pluginName)

	return Runtime{
		R:   r,
		Ctx: ctx,
		ModuleConfigs: wazero.NewModuleConfig().
			WithName("plugin").
			WithStdout(writer).
			WithStderr(writer2),
	}
}
