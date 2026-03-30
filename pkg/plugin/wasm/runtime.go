package wasm

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type Runtime struct {
	ctx    context.Context
	cancel context.CancelFunc

	alloc WasmAlloc

	enabled bool

	moduleCfg wazero.ModuleConfig
	mod       api.Module
	state     wazero.Runtime

	logwriter io.Writer

	file []byte
}

func NewRuntime(w io.Writer, file []byte, fs fs.FS) *Runtime {

	r := &Runtime{
		logwriter: w,
	}

	r.ctx, r.cancel = context.WithCancel(context.Background())
	r.state = wazero.NewRuntime(r.ctx)

	r.state.NewHostModuleBuilder("env")

	wasi_snapshot_preview1.MustInstantiate(r.ctx, r.state)

	r.moduleCfg = wazero.NewModuleConfig().
		WithName("plugin").
		WithStdout(r.logwriter).
		WithStderr(r.logwriter).
		WithFS(fs)

	r.file = file

	return r
}

func (r *Runtime) Load() error {
	deadline, cancel := context.WithDeadline(r.ctx, time.Now().Add(30*time.Second))
	defer cancel()

	var err error
	r.mod, err = r.state.InstantiateWithConfig(deadline, r.file, r.moduleCfg)
	if err != nil {
		return err
	}

	r.file = nil
	r.enabled = true

	r.alloc = newBasicAlloc(r.mod.Memory(), r.mod.ExportedFunction("alloc"), &r.ctx)

	for fnName := range r.mod.ExportedFunctionDefinitions() {
		fmt.Println(fnName)
	}

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
	if fn == nil {
		return fmt.Errorf("unkown function %s", name)
	}

	mParams := marshalParams(r.alloc, params...)

	ctx, cancel := context.WithDeadline(r.ctx, time.Now().Add(1*time.Second))
	defer cancel()

	_, err := fn.Call(ctx, mParams...)

	return err
}

func (r *Runtime) Close() error {
	defer r.cancel()
	err := r.state.Close(r.ctx)
	return err
}
