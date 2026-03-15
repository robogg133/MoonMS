package app

import (
	"encoding/hex"

	"github.com/robogg133/MoonMS/internal/packets"
)

const (
	MINECRAFT_BRAND_IDENTIFIER = "minecraft:brand"
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
	new.RegisterPacket(packets.PACKET_PLUGIN_MESSAGE, func() packets.Packet {
		return &packets.PluginMessagePacket{}
	})
	sess.KnownPkgs = new

	p, err := sess.ReadPacket()
	if err != nil {
		return err
	}

	switch p.ID() {
	case packets.PACKET_CLIENT_INFORMATION:
		cliInfo := p.(*packets.ClientInformationPacket)

		sess.PlayerInformation = cliInfo

		sess.Server.LogDebug("received client information")

	case packets.PACKET_PLUGIN_MESSAGE:
		plmsg := p.(*packets.PluginMessagePacket)

		if plmsg.Identifier == MINECRAFT_BRAND_IDENTIFIER {
			r := packets.NewReader(plmsg.Data)

			sess.Brand, err = r.ReadString()
			if err != nil {
				return err
			}
		}

		sess.Server.LogDebug("session brand = %s", sess.Brand)
		sess.Server.LogDebug("server_bound_plugin_message, identifier: %s", plmsg.Identifier)
		sess.Server.LogDebug("server_bound_plugin_message, data: %s ", hex.Dump(plmsg.Data))
	}

	/*
	 * NEED TO IMPLEMENT
	 */

	return nil
}
