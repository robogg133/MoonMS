package moonms

var heap [128 * 1024]byte
var heapOffset uint32

//export alloc
func alloc(size uint32) uint32 {
	ptr := heapOffset
	heapOffset += size

	if heapOffset >= uint32(len(heap)) {
		panic("abi heap overflow")
	}

	return ptr
}

//export abi_reset
func abiReset() {
	heapOffset = 0
}
