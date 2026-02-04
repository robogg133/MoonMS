package plugins

import (
	plugins_lua "MoonMS/internal/plugins/lua"
	plugins_wasm "MoonMS/internal/plugins/wasmi"
	"MoonMS/internal/server"
	"sync"

	"github.com/Shopify/go-lua"
)

func (p *Plugin) RunEventServerStopping(wg *sync.WaitGroup) {
	defer wg.Done()

	switch p.Manifest.Runtime {
	case "wasm":
		if err := plugins_wasm.RunServerStopping(p.Runtime.(plugins_wasm.Runtime), p.Identifier); err != nil {
			server.LogError(err)
			return
		}
	case "lua":
		plugins_lua.RunEventServerStopping(p.Runtime.(*lua.State))
	}
}
