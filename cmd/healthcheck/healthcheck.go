package main

import (
	"flag"
	"net"
	"net/netip"

	"github.com/robogg133/MoonMS/app"
	"github.com/robogg133/MoonMS/internal/packets"
)

const DEFAULT_PROTOCOL_VERSION int32 = 774

func main() {

	n := flag.Int("protocol", int(DEFAULT_PROTOCOL_VERSION), "protocol version for hello packet")

	flag.Parse()

	conn, err := net.Dial("tcp", flag.Arg(0))
	if err != nil {
		panic(err)
	}

	ip, err := netip.ParseAddrPort(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	a := &packets.HelloPacket{
		ProtocolVersion: int32(*n),
		ServerAdress:    ip.Addr().String(),
		ServerPort:      ip.Port(),
		Intent:          app.INTENT_STATUS,
	}

	b, err := packets.MarshalPacket(a, nil, 0)
	if err != nil {
		panic(err)
	}

	conn.Write(b)

}
