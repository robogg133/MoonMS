package app

import (
	"crypto/rsa"
	"fmt"
	"io"
	"net"
	"runtime/debug"
	"sync"

	"github.com/robogg133/MoonMS/internal/packets"
	"github.com/robogg133/MoonMS/plugin"
)

type Server struct {
	MinecraftConfig MinecraftServerConfig

	logFile io.Writer

	Config Config

	Plugins       map[string]*plugin.Plugin
	OnlinePlayers uint32

	PlayerList []packets.PlayerListInfo

	Sessions map[string]*Session

	ServerPrivateKey *rsa.PrivateKey

	op struct {
		lock  sync.RWMutex
		check map[string]*OPEntry
	}
	whitelist struct {
		lock  sync.RWMutex
		check map[string]bool
	}
	ban struct {
		lock  sync.RWMutex
		check map[string]*BanEntry
	}

	Bans struct {
		lock sync.RWMutex
		list []BanEntry
	}
	OPs struct {
		lock sync.RWMutex
		list []OPEntry
	}
	Whitelisteds struct {
		lock sync.RWMutex
		list []WhitelistEntry
	}
}

func New(m MinecraftServerConfig, cfg Config, sk *rsa.PrivateKey) *Server {

	server := &Server{
		MinecraftConfig:  m,
		Config:           cfg,
		Plugins:          make(map[string]*plugin.Plugin),
		ServerPrivateKey: sk,
	}

	server.op.check = make(map[string]*OPEntry)
	server.ban.check = make(map[string]*BanEntry)
	server.whitelist.check = make(map[string]bool)

	return server
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
			s.LogPanic("SERVER CRASH: %v\n%s", r, debug.Stack())
			s.Stop()
		}
	}()

	if err := s.basicFiles(); err != nil {
		panic(err)
	}
	if err := s.loadFiles(); err != nil {
		panic(err)
	}

	s.InitPlugins()

	s.LogInfo("Starting minecraft %s server on port: %d  (VERSION: %s)", s.Config.StartName, s.MinecraftConfig.Proprieties.ServerPort, s.MinecraftConfig.MinecraftVersion)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.MinecraftConfig.Proprieties.ServerPort))
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.LogPanic("%v", err)
			continue
		}
		go s.handleConn(conn)
	}

}

func (s *Server) Stop() error {
	s.LogInfo("Received stop signal, stopping the server")

	
	for id, plg := range s.Plugins {

		if plg.State != plugin.StateEnabled && plg.State != plugin.StateLoaded {
			continue
		}

		s.LogDebug("sending stop signal to plugin:%s ", id)

		if err := plg.Runtime.Call("server_stopping_event"); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) IsWhitelisted(plr string) bool {

	s.whitelist.lock.RLock()
	_, found := s.whitelist.check[plr]
	s.whitelist.lock.RUnlock()

	return found
}

// IsBanned checks from player uuid, or player ip if it is banned
func (s *Server) IsBanned(plr string) bool {

	s.ban.lock.RLock()
	_, found := s.ban.check[plr]
	s.ban.lock.RUnlock()

	return found
}

// IsOperator checks if the player uuid is operator in the server
func (s *Server) IsOperator(uuid string) bool {

	s.op.lock.RLock()
	_, found := s.op.check[uuid]
	s.op.lock.RUnlock()

	return found
}

func (s *Server) GetOpEntry(uuid string) *OPEntry {

	s.op.lock.RLock()
	entry := s.op.check[uuid]
	s.op.lock.RUnlock()

	return entry
}

func (s *Server) handleConn(conn net.Conn) {
	s.LogDebug("got connection from %s", conn.RemoteAddr().String())

	defer func() {
		if r := recover(); r != nil {
			s.LogPanic("CLOSING CONNECTION WITH %s: %v\n%s", conn.RemoteAddr().String(), r, debug.Stack())
		}
	}()

	sess := NewSession(conn, s)

	if err := sess.Run(); err != nil {
		s.LogError("%v", err)
		sess.Close()
	}
}
