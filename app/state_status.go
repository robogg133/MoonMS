package app

import (
	"MoonMS/internal/packets"
	"MoonMS/internal/server"
	"encoding/base64"
	"fmt"
	"os"
)

const STATE_NAME_STATUS = "status"

type StatusState struct{}

func (s StatusState) Name() string { return "status" }

func (s *StatusState) Handle(sess *Session) error {

	statuspkg, err := packets.ReadPackageFromConnecion(sess.Conn)
	if err != nil {
		return err
	}

	if int32(statuspkg[1]) != packets.PACKET_HANDSHAKE {
		sess.Server.LogDebug("got pkg id: ", int32(statuspkg[1]))
		return err
	}
	statuspkg = nil

	var status packets.HandShakeResponseStatus

	status.Version.Name = server.CURRENT_VERSION
	status.Version.ProtocolVersion = sess.Server.MinecraftConfig.ProtcolVersion

	status.Players.MaxPlayers = sess.Server.MinecraftConfig.Proprieties.MaxPlayer
	status.Players.OnlinePlayers = sess.Server.OnlinePlayers
	status.Players.PlayerStatus = []packets.PlayerMinimunInfo{}

	if sess.Server.MinecraftConfig.Proprieties.ServerIcon == "" {
		status.Favicon = ""
	} else {
		status.Favicon, err = getBase64Image(sess.Server.MinecraftConfig.Proprieties.ServerIcon)
		if err != nil {
			server.LogError(err)
			status.Favicon = ""
		} else {
			status.Favicon = fmt.Sprintf("data:image/png;base64,%s", status.Favicon)
		}
	}

	status.Description.Text = sess.Server.MinecraftConfig.Proprieties.Motd

	if err := sess.WritePacket(&status); err != nil {
		return err
	}

	sess.KnownPkgs.RegisterPacket(packets.PACKET_PING_PONG, func() packets.Packet { return &packets.PingPong{} })
	ping, err := sess.ReadPacket()
	if err != nil {
		return err
	}

	if ping.ID() == packets.PACKET_PING_PONG {
		if err := sess.WritePacket(ping); err != nil {
			return err
		}
	}
	return nil
}

func getBase64Image(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("Can't find specified server image using none")
		}
		return "", err
	}

	return base64.StdEncoding.EncodeToString(content), nil
}
