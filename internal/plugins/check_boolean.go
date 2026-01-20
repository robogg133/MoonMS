package plugins

import "github.com/Shopify/go-lua"

func checkBoolean(l *lua.State, arg int) bool {
	if !l.IsBoolean(arg) {
		lua.ArgumentError(l, arg, "booleano esperado")
	}
	return l.ToBoolean(arg)
}
