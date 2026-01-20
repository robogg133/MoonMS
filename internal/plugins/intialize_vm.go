package plugins

import (
	"MoonMS/internal/server"
	"strings"

	"github.com/Shopify/go-lua"
)

type LuaServerObject struct {
	Proprieties *server.ServerData
}

func createServerTable(l *lua.State, serverData *server.ServerData) {
	l.NewTable()

	tableStackIndex := l.Top()

	//

	l.PushString(serverData.MINECRAFT_VERSION)
	l.SetField(tableStackIndex, "minecraft_version")
	l.PushGoFunction(func(state *lua.State) int {
		serverData.MINECRAFT_VERSION = lua.CheckString(state, 1)
		return 0
	})
	l.SetField(tableStackIndex, "set_minecraft_version")

	//

	l.PushInteger(int(serverData.PROTOCOL_VERSION))
	l.SetField(tableStackIndex, "protocol_version")
	l.PushGoFunction(func(state *lua.State) int {
		serverData.PROTOCOL_VERSION = uint16(lua.CheckInteger(state, 1))
		return 0
	})
	l.SetField(tableStackIndex, "set_protocol_version")

	//

	l.PushString(serverData.Motd)
	l.SetField(tableStackIndex, "motd")
	l.PushGoFunction(func(state *lua.State) int {
		serverData.Motd = lua.CheckString(state, 1)
		return 0
	})
	l.SetField(tableStackIndex, "set_motd")

	//

	l.PushInteger(int(serverData.MaxPlayers))
	l.SetField(tableStackIndex, "max_players")
	l.PushGoFunction(func(state *lua.State) int {
		serverData.MaxPlayers = uint(lua.CheckInteger(state, 1))
		return 0
	})
	l.SetField(tableStackIndex, "set_max_players")

	//

	l.PushBoolean(serverData.OnlineMode)
	l.SetField(tableStackIndex, "online_mode")
	l.PushGoFunction(func(state *lua.State) int {
		serverData.OnlineMode = checkBoolean(state, 1)
		return 0
	})
	l.SetField(tableStackIndex, "set_online_mode")

	//

	l.PushBoolean(serverData.AllowServerList)
	l.SetField(tableStackIndex, "allow_server_list")
	l.PushGoFunction(func(state *lua.State) int {
		serverData.AllowServerList = checkBoolean(state, 1)
		return 0
	})
	l.SetField(tableStackIndex, "set_allow_server_list")

	l.PushString(serverData.Difficulty)
	l.SetField(tableStackIndex, "difficulty")
	l.PushGoFunction(func(state *lua.State) int {
		serverData.Difficulty = lua.CheckString(state, 1)
		return 0
	})
	l.SetField(tableStackIndex, "set_difficulty")

	l.PushString(serverData.Gamemode)
	l.SetField(tableStackIndex, "default_gamemode")
	l.PushGoFunction(func(state *lua.State) int {
		serverData.Gamemode = lua.CheckString(state, 1)
		return 0
	})
	l.SetField(tableStackIndex, "set_default_gamemode")

	//

	createEvents(l, &tableStackIndex)

	l.SetGlobal("Server")
}

type pluginName string

func (name pluginName) lPrint(state *lua.State) int {
	top := state.Top()

	var s strings.Builder

	for i := 1; i <= top; i++ {
		str, _ := state.ToString(i)
		s.WriteString(str)
	}

	server.LogPlugin(string(name), s.String())

	return 0
}

func (plugin *Plugin) startVM() error {

	lua.OpenLibraries(plugin.LuaVM)
	plugin.LuaVM.Register("print", pluginName(plugin.Manifest.Name).lPrint)

	createCallbacks(plugin.LuaVM)

	createServerTable(plugin.LuaVM, server.GetServerData())

	return nil
}
