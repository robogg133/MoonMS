package app

import (
	"encoding/hex"
	"fmt"

	"github.com/robogg133/MoonMS/internal/packets"
)

const STATE_NAME_CONFIG = "config"

type ConfigState struct{}

func (s *ConfigState) Name() string { return STATE_NAME_CONFIG }

func (s *ConfigState) Handle(sess *Session) error {
	sess.Server.LogDebug("START STATE CONFIG")

	new := make(packets.KnownPackets)
	new.RegisterPacket(packets.PACKET_CLIENT_INFORMATION, func() packets.Packet {
		return &packets.ClientInformationPacket{}
	})
	new.RegisterPacket(packets.PACKET_SERVERBOUND_PLUGIN_MESSAGE, func() packets.Packet {
		return &packets.ServerBoundPluginMessagePacket{}
	})
	sess.KnownPkgs = new

readPacketAgain:
	p, err := sess.ReadPacket()
	if err != nil {
		return err
	}

	switch p.ID() {
	case packets.PACKET_CLIENT_INFORMATION:
		cliInfo := p.(*packets.ClientInformationPacket)

		sess.PlayerInformation = cliInfo

		fmt.Println(sess.PlayerInformation)

	case 2:
		aff := p.(*packets.ServerBoundPluginMessagePacket)
		sess.Server.LogDebug(fmt.Sprintf("first server_bound_plugin_message, identifier: %s", aff.Identifier))
		sess.Server.LogDebug(fmt.Sprintf("first server_bound_plugin_message, data: %s ", hex.Dump(aff.Data)))
		goto readPacketAgain
	}

	/*
	 * NEED TO IMPLEMENT
	 */

	return nil
}
