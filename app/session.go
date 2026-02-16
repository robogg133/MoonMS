package app

import (
	"crypto/cipher"
	"errors"
	"net"

	"github.com/robogg133/KernelCraft/internal/packets"
)

type Session struct {
	Conn      net.Conn
	Server    *Server
	PkgReader *packets.Reader
	KnownPkgs packets.KnownPackets

	ClientProtocolVersion int32

	Threshold     int32
	EncryptCipher cipher.Stream
	DecryptCipher cipher.Stream

	stop bool

	PlayerInformation *packets.ClientInformationPacket

	State State
}

type State interface {
	Name() string
	Handle(*Session) error
}

var (
	ErrNoReason = errors.New("disconnect")
)

func (s *Session) Run() error {
	for {
		if err := s.State.Handle(s); err != nil {
			switch err {
			case ErrNoReason:
				return nil
			default:
				return err
			}
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

		EncryptCipher: nil,
		DecryptCipher: nil,

		Threshold: -1,
		State:     &HandshakeState{},
	}
}

func (s *Session) WritePacket(p packets.Packet) error {
	b, err := packets.MarshalPacket(p, s.EncryptCipher, s.Threshold)
	if err != nil {
		return err
	}
	_, err = s.Conn.Write(b)
	return err
}

func (s *Session) ReadPacket() (packets.Packet, error) {
	return packets.UnmarshalPacket(s.PkgReader, s.Threshold, s.KnownPkgs)
}
