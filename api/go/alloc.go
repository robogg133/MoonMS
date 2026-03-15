package moonms

var heap [128 * 1024]byte
var heapOffset uint32

//go:wasmexport alloc
func alloc(size uint32) uint32 {
	ptr := heapOffset
	heapOffset += size

	if heapOffset >= uint32(len(heap)) {
		panic("abi heap overflow")
	}

	return ptr
}

//go:wasmexport abi_reset
func abi_reset() {
	heapOffset = 0
}
