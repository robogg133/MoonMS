package plugins_lua

import "github.com/Shopify/go-lua"

func createCallbacks(l *lua.State) {
	l.NewTable()
	l.SetGlobal(callbacks_server_stop)
}

func createEvents(l *lua.State, stackIndex *int) {

	l.PushGoFunction(lOnServerStop)
	l.SetField(*stackIndex, "OnServerStop")
}
