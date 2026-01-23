package plugins_lua

import (
	"github.com/Shopify/go-lua"
)

const callbacks_server_stop = "__server_stop_callbacks"

func lOnServerStop(l *lua.State) int {

	if l.IsFunction(1) {
		l.PushString("expected function")
		l.Error()
		return 0
	}

	l.Global(callbacks_server_stop)

	l.Length(-1)
	n, _ := l.ToInteger(-1)
	l.Pop(1)
	l.PushValue(2)
	l.RawSetInt(-2, n+1)

	l.Pop(1)
	return 0
}

func RunEventServerStopping(l *lua.State) {

	l.Global(callbacks_server_stop)

	l.Length(-1)
	n, _ := l.ToInteger(-1)
	l.Pop(1)

	for i := 1; i <= n; i++ {
		l.RawGetInt(-1, i)
		l.Call(0, 0)
	}

	l.Pop(1)

}
