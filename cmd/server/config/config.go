package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

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

type Configs struct {
	Proprieties struct {
		Difficulty string `toml:"difficulty"`
		Gamemode   string `toml:"gamemode"`

		LevelName string `toml:"level-name"`
		LevelSeed string `toml:"level-seed"`

		MaxPlayers uint `toml:"max-players"`

		OnlineMode bool   `toml:"online-mode"`
		Motd       string `toml:"motd"`

		ServerPort uint16 `toml:"server-port"`

		AllowServerList bool `toml:"allow-server-list"`

		ServerIcon string `toml:"server-icon"`

		RSAKeyBits int `toml:"rsa-key-bits"`

		ServerThreshold int32 `toml:"server-threshold"`
	} `toml:"Proprieties"`
}

func ReadConfigurationFile() (Configs, error) {
	var cfg Configs

	content, err := os.Open("server-config.toml")
	if err != nil {
		if os.IsNotExist(err) {
			cfg = DefaultValuesForServerConfig()
			file, err := os.Create("server-config.toml")
			if err != nil {
				return Configs{}, err
			}

			encoder := toml.NewEncoder(file)
			if err := encoder.Encode(cfg); err != nil {
				return Configs{}, err
			}

			return cfg, nil
		}
		return Configs{}, err
	}
	defer content.Close()

	decoder := toml.NewDecoder(content)

	if err := decoder.Decode(&cfg); err != nil {
		return Configs{}, err
	}

	return cfg, nil
}

func DefaultValuesForServerConfig() Configs {
	var cfg Configs

	cfg.Proprieties.Gamemode = GAMEMODE_SURVIVAL
	cfg.Proprieties.Difficulty = DIFFICULTY_EASY
	cfg.Proprieties.LevelName = "world"
	cfg.Proprieties.OnlineMode = true
	cfg.Proprieties.ServerPort = 25565
	cfg.Proprieties.Motd = DEFAULT_MOTD
	cfg.Proprieties.MaxPlayers = 20
	cfg.Proprieties.LevelSeed = ""
	cfg.Proprieties.AllowServerList = true
	cfg.Proprieties.RSAKeyBits = 2048
	cfg.Proprieties.ServerThreshold = 256
	return cfg
}
