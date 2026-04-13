package js

import (
	"io"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/robogg133/MoonMS/pkg/plugin/js/methods"
)

type Runtime struct {
	r *goja.Runtime

	program *goja.Program

	logWriter io.Writer
}

func NewRuntime(srcCode string, name string, logWriter io.Writer) (*Runtime, error) {
	registry := new(require.Registry)
	r := &Runtime{}

	r.r = goja.New()
	req := registry.Enable(r.r)

	p, err := goja.Compile(name, srcCode, false)
	if err != nil {
		return nil, err
	}
	r.program = p

	r.logWriter = logWriter

	if err := r.r.Set("console", methods.NewConsole(r.logWriter)); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Runtime) Load() error {

	r.r.RunProgram(r.program)

	return nil
}
