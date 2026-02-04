package plugin_wasmi

import (
	"MoonMS/internal/server"
	"context"
	"fmt"
	"time"
)

const _EVENT_SERVER_STOP = "on_server_stop"

func RunServerStopping(r Runtime, pluginName string) error {

	mod := r.R.Module("plugin")

	ctx, cancel := context.WithDeadline(r.Ctx, time.Now().Add(30*time.Second))

	defer cancel()

	if mod == nil || mod.IsClosed() {
		server.LogInfo(fmt.Sprintf("Trying to bind event for plugin: \"%s\", but the plugin is \"offline\"", pluginName))
		return nil
	}

	fn := mod.ExportedFunction(_EVENT_SERVER_STOP)

	_, err := fn.Call(ctx)
	if err != nil {
		return err
	}
	return nil
}
