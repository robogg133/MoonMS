package methods

import (
	"fmt"
	"io"

	"github.com/dop251/goja"
)

func NewConsole(w io.Writer) map[string]interface{} {

	console := map[string]interface{}{
		"log": func(call goja.FunctionCall) goja.Value {
			for _, arg := range call.Arguments {
				w.Write([]byte(fmt.Sprint(arg.Export())))
			}
			return goja.Undefined()
		},
	}

	return console
}
