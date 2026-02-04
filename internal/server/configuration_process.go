package server

import (
	"MoonMS/internal/packets"
	"net"
	"sync"
)

func ReadEntireConfigurationPrcs(conn net.Conn, lock *sync.WaitGroup) {
	for {

		lock.Wait()
		pkg, err := packets.ReadPackageFromConnecion(conn)
		if err != nil {
			LogError(err)
			return
		}

		pkgId, err := packets.RecongnizePacket(pkg)
		if err != nil {
			LogError(err)
			return
		}

		switch pkgId {
		case packets.PACKET_CUSTOM_PAYLOAD:
			// NEED AN IMPLEMENTATION

		}

	}
}
