package wasm

import (
	"context"
	"io"

	"github.com/tetratelabs/wazero"
)

type WasmRuntime struct {
	ctx       context.Context
	moduleCfg wazero.ModuleConfig
	runner    wazero.Runtime
}

func NewRuntime(w io.Writer) WasmRuntime {

	var r WasmRuntime

	r.moduleCfg = wazero.NewModuleConfig().
		WithName("plugin").
		WithStdout(w).
		WithStderr(w).
		WithStdin(nil)

	r.ctx = context.Background()
	r.runner = wazero.NewRuntime(r.ctx)

	return r
}
