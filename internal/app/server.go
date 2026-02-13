package app

import (
	"crypto/rsa"
	"fmt"
	"io"
	"net"
	"os"
	"runtime/debug"
)

type Server struct {
	MinecraftConfig MinecraftServerConfig

	Config Config

	ServerPrivateKey *rsa.PrivateKey
}

type MinecraftServerConfig struct {
	ProtcolVersion int32

	MinecraftVersion string

	Motd string

	Difficulty    string
	Gamemode      string
	ForceGamemode bool
	LevelName     string
	Seed          string
	Hardcore      bool

	MaxPlayer uint

	OnlineMode bool

	AllowServerList bool

	ServerIcon string

	ServerPort uint16

	RSAKeyBitAmmount uint

	Threshold uint

	ViewDistance uint16

	SimluationDistance uint16

	AllowNether bool

	AllowEnd bool

	TPS float32

	Encryption bool

	Whitelist bool
}

type Config struct {
	LatestLogFile io.Writer
	DebugEnabled  bool

	StartName string

	PluginsFolder string
}

func (s *Server) Start() {
	defer func() {
		if r := recover(); r != nil {
			s.LogPanic(fmt.Sprintf("SERVER CRASH: %v\n%s", r, debug.Stack()))
			s.Stop()
			os.Exit(1)
		}
	}()

	s.LogInfo(fmt.Sprintf("Starting minecraft %s server on port: %d", s.Config.StartName, s.MinecraftConfig.ServerPort))
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.MinecraftConfig.ServerPort))
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.LogPanic(err)
			continue
		}

	}

}

func (s *Server) Stop() error {
	return nil
}
