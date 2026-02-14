package app

import (
	"MoonMS/internal/plugins"
	"crypto/rsa"
	"fmt"
	"io"
	"net"
	"os"
	"runtime/debug"
	"sync"
)

type Server struct {
	MinecraftConfig MinecraftServerConfig

	logFile io.Writer

	Config Config

	Plugins       map[string]plugins.Plugin
	OnlinePlayers uint32

	ServerPrivateKey *rsa.PrivateKey
}

func New(m MinecraftServerConfig, cfg Config, sk *rsa.PrivateKey) *Server {
	return &Server{
		MinecraftConfig:  m,
		Config:           cfg,
		Plugins:          make(map[string]plugins.Plugin),
		ServerPrivateKey: sk,
	}
}

type Config struct {
	LatestLogFile string

	DebugEnabled bool

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

	if err := s.basicFiles(); err != nil {
		panic(err)
	}

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
		go s.handleConn(conn)
	}

}

func (s *Server) handleConn(conn net.Conn) {
	s.LogDebug(fmt.Sprintf("got connection from %s", conn.RemoteAddr().String()))
	sess := NewSession(conn, s)

	sess.Run()
}

func (s *Server) Stop() error {

	var wg sync.WaitGroup

	for _, plg := range s.Plugins {
		wg.Add(1)
		go plg.RunEventServerStopping(&wg)
	}

	wg.Wait()
	return nil
}
