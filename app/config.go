package app

import (
	"MoonMS/pkg/minecraft/world/seed"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

const MAIN_CONFIG_FILE_PATH string = "server-config.toml"

const (
	GAMEMODE_SURVIVAL  = "survival"
	GAMEMODE_CREATIVE  = "creative"
	GAMEMODE_ADVENTURE = "adventure"
	GAMEMODE_SPECTATOR = "spectator"
)

const (
	DIFFICULTY_EASY   = "easy"
	DIFFICULTY_NORMAL = "normal"
	DIFFICULTY_HARD   = "hard"
)

const DEFAULT_MOTD string = "A Moon Minecraft Server!"

type MinecraftServerConfig struct {
	Proprieties struct {
		Motd string `toml:"motd"`

		Difficulty    string `toml:"difficulty"`
		Gamemode      string `toml:"default-gamemode"`
		ForceGamemode bool   `toml:"force-gamemode"`
		LevelName     string `toml:"level-name" `
		Seed          int64  `toml:"level-seed"`
		Hardcore      bool   `toml:"hardcore"`

		MaxPlayer uint `toml:"max-players"`

		OnlineMode bool `toml:"online-mode"`

		AllowServerList bool `toml:"allow-server-list"`

		ServerIcon string `toml:"sever-icon-path"`

		ServerPort uint16 `toml:"server-port"`

		ViewDistance uint16 `toml:"view-distance"`

		SimluationDistance uint16 `toml:"simulation-distance"`

		AllowNether bool `toml:"allow-nether"`

		AllowEnd bool `toml:"allow-end"`

		TPS float32 `toml:"tps"`

		Whitelist bool `toml:"whitelist"`
	} `toml:"Proprieties"`

	Advanced struct {
		OfflineEncryption bool `toml:"offline-encryption"`

		RSAKeyBitAmmount uint  `toml:"rsa-key-bit-ammount"`
		Threshold        int32 `toml:"threshold"`
	} `toml:"Advanced"`
	ProtcolVersion int32 `toml:",omitempty"`

	MinecraftVersion string `toml:",omitempty"`
}

func (cfg *MinecraftServerConfig) ConfigFile() error {
readAgain:
	b, err := os.ReadFile(MAIN_CONFIG_FILE_PATH)
	if err != nil {
		if os.IsNotExist(err) {

			_ = os.MkdirAll(filepath.Dir(MAIN_CONFIG_FILE_PATH), 0755)
			f, err := os.Create(MAIN_CONFIG_FILE_PATH)
			if err != nil {
				return err
			}

			b, err := toml.Marshal(getDefaultCfgFile())
			if err != nil {
				return err
			}

			_, err = f.Write(b)
			if err != nil {
				return err
			}
			goto readAgain
		}
		return err
	}
	if err := toml.Unmarshal(b, &cfg); err != nil {
		return err
	}

	cfg.ProtcolVersion = 774

	return nil
}

func getDefaultCfgFile() MinecraftServerConfig {

	return MinecraftServerConfig{
		Proprieties: struct {
			Motd               string  "toml:\"motd\""
			Difficulty         string  "toml:\"difficulty\""
			Gamemode           string  "toml:\"default-gamemode\""
			ForceGamemode      bool    "toml:\"force-gamemode\""
			LevelName          string  "toml:\"level-name\" "
			Seed               int64   "toml:\"level-seed\""
			Hardcore           bool    "toml:\"hardcore\""
			MaxPlayer          uint    "toml:\"max-players\""
			OnlineMode         bool    "toml:\"online-mode\""
			AllowServerList    bool    "toml:\"allow-server-list\""
			ServerIcon         string  "toml:\"sever-icon-path\""
			ServerPort         uint16  "toml:\"server-port\""
			ViewDistance       uint16  "toml:\"view-distance\""
			SimluationDistance uint16  "toml:\"simulation-distance\""
			AllowNether        bool    "toml:\"allow-nether\""
			AllowEnd           bool    "toml:\"allow-end\""
			TPS                float32 "toml:\"tps\""
			Whitelist          bool    "toml:\"whitelist\""
		}{
			Motd:               DEFAULT_MOTD,
			Difficulty:         DIFFICULTY_EASY,
			Gamemode:           GAMEMODE_SURVIVAL,
			ForceGamemode:      false,
			LevelName:          "world",
			Seed:               seed.GenerateSeed(),
			Hardcore:           false,
			MaxPlayer:          20,
			OnlineMode:         true,
			AllowServerList:    true,
			ServerPort:         25565,
			ViewDistance:       10,
			SimluationDistance: 16,
			AllowNether:        true,
			AllowEnd:           true,
			TPS:                20.0,
			Whitelist:          false,
		},
		Advanced: struct {
			OfflineEncryption bool  "toml:\"offline-encryption\""
			RSAKeyBitAmmount  uint  "toml:\"rsa-key-bit-ammount\""
			Threshold         int32 "toml:\"threshold\""
		}{
			OfflineEncryption: true,
			RSAKeyBitAmmount:  2048,
			Threshold:         256,
		},
	}
}
