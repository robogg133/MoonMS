//Needs tinygo compiler
//go:build wasip1

package moonms_api

import (
	"unsafe"
)

//export alloc
func alloc(size uint32) uint32 {
	buf := make([]byte, size)
	return uint32(uintptr(unsafe.Pointer(&buf[0])))
}

func allocString(s string) (uint32, uint32) {
	old := []byte(s)

	buf := make([]byte, len(old))

	copy(buf, old)

	ptr := uint32(uintptr(unsafe.Pointer(&buf[0])))

	return ptr, uint32(len(buf))
}

//go:wasm-module env

/* GetMaxPlayers returns the max player number in the server */
//go:export get_server_max_players
func GetMaxPlayers() uint

/* SetMaxPlayers set the max player number in the server */
//go:export set_server_max_players
func SetMaxPlayers(n uint32)

/* GetServerThreshold returns server threshold */
//go:export set_server_threshold
func GetServerThreshold() int32

//go:export get_server_motd
func get_server_motd() uint64

/* GetServerMOTD returns the server motd (description) */
func GetServerMOTD() string {
	pack := get_server_motd()

	ptr := uint32(pack >> 32)
	length := uint32(pack & 0xffffffff)

	b := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), length)

	return string(b)
}

//go:export set_server_motd
func set_server_motd(ptr uint32, lenght uint32)

/* SetServerMOTD changes the server motd (description) */
func SetServerMOTD(s string) { set_server_motd(allocString(s)) }
