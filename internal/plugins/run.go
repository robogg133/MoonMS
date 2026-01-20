package plugins

import "github.com/Shopify/go-lua"

func (plugin *Plugin) LoadPlugin() error {

	if err := plugin.startVM(); err != nil {
		return err
	}

	if err := lua.DoString(plugin.LuaVM, plugin.MainLuaFile); err != nil {
		return err
	}

	AllPlugins = append(AllPlugins, *plugin)
	return nil
}
