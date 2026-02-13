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

	s.LogInfo(fmt.Sprintf("Starting minecraft %s server on port: %d", s.Config.StartName, s.MinecraftConfig.Proprieties.ServerPort))
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.MinecraftConfig.Proprieties.ServerPort))
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
		defer conn.Close()

	}

}

func (s *Server) Stop() error {
	return nil
}
