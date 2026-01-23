package plugins_lua

import "github.com/Shopify/go-lua"

func LoadPlugin(l *lua.State, luaFile, pluginName string) error {

	if err := startVM(l, pluginName); err != nil {
		return err
	}

	if err := lua.DoString(l, luaFile); err != nil {
		return err
	}

	return nil
}
