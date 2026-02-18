package app

import (
	"crypto/rsa"
	"fmt"
	"io"
	"net"
	"runtime/debug"

	"github.com/robogg133/MoonMS/internal/packets"
)

type Server struct {
	MinecraftConfig MinecraftServerConfig

	logFile io.Writer

	Config Config

	//Plugins       map[string]plugins.Plugin
	OnlinePlayers uint32

	PlayerList []packets.PlayerListInfo

	ServerPrivateKey *rsa.PrivateKey
}

func New(m MinecraftServerConfig, cfg Config, sk *rsa.PrivateKey) *Server {
	return &Server{
		MinecraftConfig: m,
		Config:          cfg,
		//Plugins:          make(map[string]plugins.Plugin),
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

	defer func() {
		if r := recover(); r != nil {
			s.LogPanic(fmt.Sprintf("CLOSING CONNECTION WITH %s: %v\n%s", conn.RemoteAddr().String(), r, debug.Stack()))
		}
	}()

	sess := NewSession(conn, s)

	if err := sess.Run(); err != nil {
		s.LogError(err)
	}
}

func (s *Server) Stop() error {

	return nil
}
