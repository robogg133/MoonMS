package main

import (
	"bytes"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"runtime/debug"
	"time"

	"github.com/robogg133/MoonMS/app"
	"github.com/robogg133/MoonMS/internal/packets"
)

const DEFAULT_PROTOCOL_VERSION int32 = 774

func main() {

	defer func() {
		r := recover()

		if r == nil {
			os.Exit(0)
		}

		fmt.Println(r)
		fmt.Println(string(debug.Stack()))
		os.Exit(1)

	}()

	n := flag.Int("protocol", int(DEFAULT_PROTOCOL_VERSION), "protocol version for hello packet")

	flag.Parse()
	conn, err := net.Dial("tcp", flag.Arg(0))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	host, port, err := net.SplitHostPort(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	var portInt int
	fmt.Sscanf(port, "%d", &portInt)

	a := &packets.HelloPacket{
		ProtocolVersion: int32(*n),
		ServerAdress:    host,
		ServerPort:      uint16(portInt),
		Intent:          app.INTENT_STATUS,
	}
	b, err := packets.MarshalPacket(a, nil, -1)
	if err != nil {
		panic(err)
	}

	if _, err := conn.Write(b); err != nil {
		panic(err)
	}

	if _, err := conn.Write([]byte{1, 0}); err != nil {
		panic(err)
	}

	rawpkg, err := packets.ReadPackageFromConnecion(conn)
	if err != nil {
		panic(err)
	}

	knownPackages := make(packets.KnownPackets)

	knownPackages.RegisterPacket(packets.PACKET_STATUS, func() packets.Packet {
		return &packets.StatusPacket{}
	})

	knownPackages.RegisterPacket(packets.PACKET_PING_PONG, func() packets.Packet {
		return &packets.PingPong{}
	})

	pkg, err := packets.UnmarshalPacket(packets.NewReader(rawpkg), -1, knownPackages)
	if err != nil {
		panic(err)
	}

	pkt := pkg.(*packets.StatusPacket)

	pingData := make([]byte, 8)

	if _, err := io.ReadFull(rand.Reader, pingData); err != nil {
		panic(err)
	}

	ping := &packets.PingPong{
		Bytes: [8]byte(pingData),
	}

	b, err = packets.MarshalPacket(ping, nil, -1)
	if err != nil {
		panic(err)
	}

	pingTime := time.Now()
	if _, err := conn.Write(b); err != nil {
		panic(err)
	}

	rawpkg, err = packets.ReadPackageFromConnecion(conn)
	if err != nil {
		panic(err)
	}
	rcv := time.Since(pingTime)

	pingResp, err := packets.UnmarshalPacket(packets.NewReader(rawpkg), -1, knownPackages)
	if err != nil {
		panic(err)
	}

	if !bytes.Equal(pingResp.(*packets.PingPong).Bytes[:], pingData) {
		fmt.Printf("Players: %d/%d \nVersion: %s (%d)\nDescription: %s\n\nPing ??? ms\n", pkt.Players.OnlinePlayers, pkt.Players.MaxPlayers, pkt.Version.Name, pkt.Version.ProtocolVersion, pkt.Description.Text)
		os.Exit(0)
	}

	fmt.Printf("Players: %d/%d \nVersion: %s (%d)\nDescription: %s\nPing %d ms\n", pkt.Players.OnlinePlayers, pkt.Players.MaxPlayers, pkt.Version.Name, pkt.Version.ProtocolVersion, pkt.Description.Text, rcv.Milliseconds())

	os.Exit(0)
}
