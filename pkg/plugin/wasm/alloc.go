package wasm

import (
	"context"

	"github.com/tetratelabs/wazero/api"
)

type WasmAlloc interface {
	Alloc(uint32) uint32
	Write(uint32, []byte)
}

type WasmBasicAllocator struct {
	mem api.Memory
	fn  api.Function
	ctx *context.Context
}

func newBasicAlloc(mem api.Memory, fn api.Function, ctx *context.Context) WasmAlloc {
	return &WasmBasicAllocator{
		mem: mem,
		fn:  fn,
		ctx: ctx,
	}
}

func (w *WasmBasicAllocator) Alloc(size uint32) uint32 {
	ptrRes, _ := w.fn.Call(*w.ctx, uint64(size))
	return uint32(ptrRes[0])
}

func (w *WasmBasicAllocator) Write(ptr uint32, data []byte) {
	w.mem.Write(ptr, data)
}
