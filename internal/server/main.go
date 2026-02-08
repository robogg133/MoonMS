package server

import (
	"MoonMS/cmd/server/config"
	"os"
)

var PROTOCOL_VERSION int32 = 774
var CURRENT_VERSION string = "1.21.11"

const MOONMS_VERSION string = "1.0.0"

var LatestLogFile *os.File

var ServerDataPublic ServerData

type ServerData struct {
	PROTOCOL_VERSION  int32
	MINECRAFT_VERSION string

	Motd string

	VERSION string // Imutable

	Difficulty string

	Gamemode string

	LevelName string // Imutable

	Seed string // Imutable

	MaxPlayers uint

	OnlineMode bool

	ServerPort uint16 // Imutable

	AllowServerList bool

	ServerIcon string

	RSAKeyBits int

	Threshold int32
}

func InitServerData() (*ServerData, error) {
	LogInfo("Reading main configuration file")
	configs, err := config.ReadConfigurationFile()
	if err != nil {
		return nil, err
	}

	LatestLogFile, err = os.Create("logs/latest.log")
	if err != nil {
		return nil, err
	}

	if err := LatestLogFile.Chmod(0644); err != nil {
		return nil, err
	}

	ServerDataPublic = ServerData{
		PROTOCOL_VERSION:  PROTOCOL_VERSION,
		MINECRAFT_VERSION: CURRENT_VERSION,

		VERSION: MOONMS_VERSION,

		Difficulty: configs.Proprieties.Difficulty,

		Gamemode: configs.Proprieties.Gamemode,

		LevelName: configs.Proprieties.LevelName,

		Seed: configs.Proprieties.LevelSeed,

		MaxPlayers: configs.Proprieties.MaxPlayers,

		OnlineMode: configs.Proprieties.OnlineMode,

		ServerPort: configs.Proprieties.ServerPort,

		AllowServerList: configs.Proprieties.AllowServerList,

		ServerIcon: configs.Proprieties.ServerIcon,

		Motd: configs.Proprieties.Motd,

		RSAKeyBits: configs.Proprieties.RSAKeyBits,

		Threshold: configs.Proprieties.ServerThreshold,
	}

	return &ServerDataPublic, nil
}

func GetServerData() *ServerData { return &ServerDataPublic }
