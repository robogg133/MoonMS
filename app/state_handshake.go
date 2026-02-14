package app

import (
	"MoonMS/internal/packets"
	"fmt"
)

const (
	INTENT_STATUS int32 = 1
)
const STATE_NAME_HANDSHAKE = "handshake"

type HandshakeState struct{}

func (s HandshakeState) Name() string { return STATE_NAME_HANDSHAKE }

func (s *HandshakeState) Handle(sess *Session) error {
	pkg, err := sess.ReadPacket()
	if err != nil {
		return err
	}
	sess.Server.LogDebug(fmt.Sprintf("current state = %s, reading packet", s.Name()))

	if pkg.ID() != packets.PACKET_HANDSHAKE {
		sess.Server.LogDebug("pkg id sent: ", pkg.ID())
		return err
	}

	sess.ClientProtocolVersion = pkg.(*packets.Handshake).ProtocolVersion

	switch pkg.(*packets.Handshake).Intent {
	case INTENT_STATUS:
		sess.State = &StatusState{}
	}
	return nil
}
