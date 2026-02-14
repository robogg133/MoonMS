package app

import (
	"MoonMS/internal/packets"
	"crypto/cipher"
	"net"
)

type Session struct {
	Conn      net.Conn
	Server    *Server
	PkgReader *packets.Reader
	KnownPkgs packets.KnownPackets

	ClientProtocolVersion int32

	Threshold int
	Stream    cipher.Stream

	State State
}

type State interface {
	Name() string
	Handle(*Session) error
}

func (s *Session) Run() error {
	for {
		if err := s.State.Handle(s); err != nil {
			return err
		}
	}
}

// Returns new session object, only knowing handshake package, and with handshake state
func NewSession(conn net.Conn, server *Server) *Session {

	kpkg := make(packets.KnownPackets)
	kpkg.RegisterPacket(packets.PACKET_HANDSHAKE, func() packets.Packet {
		return &packets.Handshake{}
	})

	return &Session{
		Conn:      conn,
		Server:    server,
		PkgReader: packets.NewReaderFromReader(conn),
		KnownPkgs: kpkg,

		Threshold: -1,
		Stream:    nil,
		State:     &HandshakeState{},
	}
}

func (s *Session) WritePacket(p packets.Packet) error {
	b, err := packets.MarshalPacket(p, s.Stream, s.Threshold)
	if err != nil {
		return err
	}
	_, err = s.Conn.Write(b)
	return err
}

func (s *Session) ReadPacket() (packets.Packet, error) {
	return packets.UnmarshalPacket(s.PkgReader, s.Threshold, s.KnownPkgs)
}
