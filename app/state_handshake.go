package app

import (
	"fmt"

	"github.com/robogg133/MoonMS/internal/packets"
)

const (
	INTENT_STATUS   int32 = 1
	INTENT_LOGIN    int32 = 2
	INTENT_TRANSFER int32 = 3
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
	case INTENT_LOGIN:
		sess.State = &LoginState{}
	}
	return nil
}
