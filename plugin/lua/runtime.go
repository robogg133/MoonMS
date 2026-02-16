package lua

import (
	"io"

	"github.com/Shopify/go-lua"
)

func NewRuntime(w io.Writer) {
	l := lua.NewState()

	lua.OpenLibraries(l)
}
