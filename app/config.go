package app

import (
	"compress/gzip"
	"os"
	"path/filepath"
	"time"

	"github.com/dgraph-io/badger/v4/options"
	"github.com/robogg133/MoonMS/pkg/minecraft/world/seed"

	"github.com/pelletier/go-toml/v2"
)

const (
	MainConfigFilePath     string = "configs/server-config.toml"
	DatabaseConfigFilePath string = "configs/database-config.toml"
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

const (
	CompressionSnappy string = "snappy"
	CompressionZSTD   string = "zstd"
	CompressionNone   string = "none"
)

const DEFAULT_MOTD string = "A Minecraft server"

type DatabaseConfig struct {
	Active      bool          `toml:"active"`
	MemTable    int64         `toml:"mem_table"`
	BlockCache  int64         `toml:"block_cache"`
	ValueLog    int64         `toml:"value_log"`
	Compression string        `toml:"compression"`
	SyncWrites  bool          `toml:"sync_writes"`
	GCInterval  time.Duration `toml:"gc_interval"`
	GCDiscard   float64       `toml:"gc_discard"`
}

type MinecraftServerConfig struct {
	Proprieties struct {
		Motd string `toml:"motd"`

		Difficulty    string `toml:"difficulty"`
		Gamemode      string `toml:"default-gamemode"`
		ForceGamemode bool   `toml:"force-gamemode"`
		LevelName     string `toml:"level-name" `
		Seed          int64  `toml:"level-seed"`
		Hardcore      bool   `toml:"hardcore"`

		MaxPlayer uint32 `toml:"max-players"`

		OnlineMode bool `toml:"online-mode"`

		AllowServerList bool `toml:"allow-server-list"`

		ServerIcon string `toml:"sever-icon-path"`

		ServerPort uint16 `toml:"server-port"`

		ViewDistance uint8 `toml:"view-distance"`

		SimluationDistance uint8 `toml:"simulation-distance"`

		TPS float32 `toml:"tps"`

		Whitelist bool `toml:"whitelist"`

		LogCompressionLevel int `toml:"log-compression-level" comment:"best compression level: 9, default: -1, best speed: 1"`
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
	b, err := os.ReadFile(MainConfigFilePath)
	if err != nil {
		if os.IsNotExist(err) {

			_ = os.MkdirAll(filepath.Dir(MainConfigFilePath), 0755)
			f, err := os.Create(MainConfigFilePath)
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

	return nil
}

func (cfg *DatabaseConfig) ConfigFile() error {
readAgain:
	b, err := os.ReadFile(DatabaseConfigFilePath)
	if err != nil {
		if os.IsNotExist(err) {

			_ = os.MkdirAll(filepath.Dir(DatabaseConfigFilePath), 0755)
			f, err := os.Create(DatabaseConfigFilePath)
			if err != nil {
				return err
			}

			b, err := toml.Marshal(getDefaultCfgFileDB())
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
	return nil
}

func getDefaultCfgFileDB() DatabaseConfig {
	return DatabaseConfig{
		Active:      true,
		MemTable:    64 << 20,
		BlockCache:  32 << 20,
		ValueLog:    256 << 20,
		Compression: CompressionZSTD,
		SyncWrites:  false,
		GCInterval:  10 * time.Minute,
		GCDiscard:   0.5,
	}
}
func compressionType(s string) options.CompressionType {
	switch s {
	case CompressionSnappy:
		return options.Snappy
	case CompressionZSTD:
		return options.ZSTD
	default:
		return options.None
	}
}

func getDefaultCfgFile() MinecraftServerConfig {

	return MinecraftServerConfig{
		Proprieties: struct {
			Motd                string  "toml:\"motd\""
			Difficulty          string  "toml:\"difficulty\""
			Gamemode            string  "toml:\"default-gamemode\""
			ForceGamemode       bool    "toml:\"force-gamemode\""
			LevelName           string  "toml:\"level-name\" "
			Seed                int64   "toml:\"level-seed\""
			Hardcore            bool    "toml:\"hardcore\""
			MaxPlayer           uint32  "toml:\"max-players\""
			OnlineMode          bool    "toml:\"online-mode\""
			AllowServerList     bool    "toml:\"allow-server-list\""
			ServerIcon          string  "toml:\"sever-icon-path\""
			ServerPort          uint16  "toml:\"server-port\""
			ViewDistance        uint8   "toml:\"view-distance\""
			SimluationDistance  uint8   "toml:\"simulation-distance\""
			TPS                 float32 "toml:\"tps\""
			Whitelist           bool    "toml:\"whitelist\""
			LogCompressionLevel int     "toml:\"log-compression-level\" comment:\"best compression level: 9, default: -1, best speed: 1\""
		}{
			Motd:                DEFAULT_MOTD,
			Difficulty:          DIFFICULTY_EASY,
			Gamemode:            GAMEMODE_SURVIVAL,
			ForceGamemode:       false,
			LevelName:           "world",
			Seed:                seed.GenerateSeed(),
			Hardcore:            false,
			MaxPlayer:           20,
			OnlineMode:          true,
			AllowServerList:     true,
			ServerPort:          25565,
			ViewDistance:        10,
			SimluationDistance:  16,
			TPS:                 20.0,
			Whitelist:           false,
			LogCompressionLevel: gzip.DefaultCompression,
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
