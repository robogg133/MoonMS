package wasm

import (
	"context"
	"io"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

type Runtime struct {
	ctx    context.Context
	cancel context.CancelFunc

	alloc WasmAlloc

	enabled bool

	moduleCfg wazero.ModuleConfig
	mod       api.Module
	state     wazero.Runtime

	file []byte
}

func NewRuntime(w io.Writer, file []byte) Runtime {

	var r Runtime

	r.moduleCfg = wazero.NewModuleConfig().
		WithName("plugin").
		WithStdout(w).
		WithStderr(w).
		WithStdin(nil)

	r.ctx, r.cancel = context.WithCancel(context.Background())
	r.state = wazero.NewRuntime(r.ctx)
	r.file = file
	return r
}

func (r *Runtime) Load() error {

	deadline, cancel := context.WithDeadline(r.ctx, time.Now().Add(30*time.Second))
	defer cancel()

	cmpModule, err := r.state.CompileModule(deadline, r.file)
	if err != nil {
		return err
	}

	r.file = nil

	r.mod, err = r.state.InstantiateModule(deadline, cmpModule, r.moduleCfg)
	if err != nil {
		return err
	}

	r.alloc = newBasicAlloc(r.mod.Memory(), r.mod.ExportedFunction("alloc"), &r.ctx)

	return nil
}

func (r *Runtime) Pause() { r.enabled = false }

func (r *Runtime) Resume() { r.enabled = true }

func (r *Runtime) Tick(deadline time.Time) error {
	if !r.enabled {
		return nil
	}

	fn := r.mod.ExportedFunction("tick")

	ctx, cancel := context.WithDeadline(r.ctx, deadline)
	defer cancel()

	_, err := fn.Call(ctx)
	return err
}

func (r *Runtime) Call(name string, params ...any) error {
	if !r.enabled {
		return nil
	}

	fn := r.mod.ExportedFunction(name)

	mParams := marshalParams(r.alloc, params...)

	ctx, cancel := context.WithDeadline(r.ctx, time.Now().Add(1*time.Second))
	defer cancel()

	_, err := fn.Call(ctx, mParams...)

	return err
}

func (r *Runtime) Stop() error {
	r.cancel()
	err := r.state.Close(r.ctx)
	return err
}
