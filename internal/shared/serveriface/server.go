package serveriface

import (
	"crypto/rsa"

	"github.com/robogg133/MoonMS/internal/packets"
)

type Server interface {
	GetMinecraftConfig() *MinecraftServerConfig

	SetMotd(motd string) error
	SetMaxPlayers(max uint) error
	SetDifficulty(diff string) error
	SetGamemode(gm string) error

	GetOnlinePlayers() uint32
	GetPlayerList() []packets.PlayerListInfo

	Log(args ...any)

	GetPrivateKey() *rsa.PrivateKey
}

type MinecraftServerConfig struct {
	Proprieties struct {
		Motd               string
		Difficulty         string
		Gamemode           string
		ForceGamemode      bool
		LevelName          string
		Seed               int64
		Hardcore           bool
		MaxPlayer          uint
		OnlineMode         bool
		AllowServerList    bool
		ServerIcon         string
		ServerPort         uint16
		ViewDistance       uint8
		SimluationDistance uint8
		AllowNether        bool
		AllowEnd           bool
		TPS                float32
		Whitelist          bool
	}

	Advanced struct {
		OfflineEncryption bool
		RSAKeyBitAmmount  uint
		Threshold         int32
	}

	ProtcolVersion   int32
	MinecraftVersion string
}
