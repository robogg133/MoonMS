package app

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"runtime"

	"github.com/robogg133/MoonMS/internal/shared/plataforms"
)

type PluginWriter struct {
	plName string
	s      *Server
}

func (s *Server) NewPluginWriter(plName string) *PluginWriter {
	return &PluginWriter{
		plName: plName,
		s:      s,
	}
}

func (pl *PluginWriter) SetName(s string) {
	pl.plName = s
}
func (pl *PluginWriter) Write(b []byte) (int, error) {

	pl.s.LogPlugin(pl.plName, strings.TrimSuffix(string(b), "\n"))
	return 0, nil
}

func (s *Server) LogInfo(format string, args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf(format, args...)
	debug := s.debugStatus()+":"
	fmt.Printf("[%s] \033[34m[INFO]\033[0m%s %v\n", timestamp, addDebugColor(debug), message)
	fmt.Fprintf(s.logFile, "[%s] [INFO]%s %v\n", timestamp, debug, message)
}

func (s *Server) LogDebug(format string, args ...any) {
	if !s.Config.DebugEnabled {
		return
	}
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf(format, args...)
	debug := s.debugStatus()+":"
	fmt.Printf("[%s] \033[33m[DEBUG]%s\033[0m %v\n", timestamp, debug, message)
	fmt.Fprintf(s.logFile, "[%s] [DEBUG]%s %v\n", timestamp, debug, message)
}

func (s *Server) LogWarn(format string, args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf(format, args...)
	debug := s.debugStatus()+":"
	fmt.Printf("[%s] \033[33m[WARN]\033[0m%s %v\n", timestamp,addDebugColor(debug), message)
	fmt.Fprintf(s.logFile, "[%s] [WARN]%s %v\n", timestamp, debug, message)
}

func (s *Server) LogError(format string, args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf(format, args...)
	debug := s.debugStatus()+":"
	fmt.Printf("[%s] \033[31m[ERROR]\033[0m%s %v\n", timestamp, addDebugColor(debug), message)
	fmt.Fprintf(s.logFile, "[%s] [ERROR]%s %v\n", timestamp, debug, message)
}

func (s *Server) LogPanic(format string, args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf(format, args...)
	debug := s.debugStatus()+":"
	fmt.Printf("[%s] \033[31m[PANIC]:\033[0m%s %v\n", timestamp, addDebugColor(debug), message)
	fmt.Fprintf(s.logFile, "[%s] [PANIC]%s %v\n", timestamp, debug, message)
}

func (s *Server) LogPlugin(pluginName string, args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf("%v", args...)
	fmt.Printf("[%s] \033[32m[%s]\033[0m %v \n", timestamp, pluginName, message)
	fmt.Fprintf(s.logFile, "[%s] [%s] %v\n", timestamp, pluginName, message)
}

func (s *Server) StartLogger() error {

	plataforms.EnableANSII()

	_ = os.MkdirAll(filepath.Dir(s.Config.LatestLogFile), 0755)

	if err := s.compressLog(); err != nil {
		return err
	}

	f, err := os.Create(s.Config.LatestLogFile)
	if err != nil {
		return err
	}

	s.logFile = f
	return nil
}

func (s *Server) compressLog() error {
	oldLog, err := os.Open(s.Config.LatestLogFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer oldLog.Close()

	dir := filepath.Dir(s.Config.LatestLogFile)
	f, err := os.Create(fmt.Sprintf("%s/%s.log.gz", dir, time.Now().Format("2006-01-02 15-04-05")))
	if err != nil {
		return err
	}
	defer f.Close()

	writer, err := gzip.NewWriterLevel(f, s.MinecraftConfig.Proprieties.LogCompressionLevel)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, oldLog)
	if err != nil {
		writer.Close()
		return err
	}

	defer os.Remove(s.Config.LatestLogFile)

	return writer.Close()
}

//

type logWrapper struct {
	server *Server
}

func NewBadgerLogWrapper(s *Server) *logWrapper {

	return &logWrapper{
		server: s,
	}

}
func (l *logWrapper) Debugf(s string, args ...interface{}) {
	l.server.LogDebug(s, args...)
}

func (l *logWrapper) Errorf(s string, args ...interface{}) {
	l.server.LogError(s, args...)
}

func (l *logWrapper) Warningf(s string, args ...interface{}) {
	l.server.LogWarn(s, args...)
}

func (l *logWrapper) Infof(s string, args ...interface{}) {
	l.server.LogDebug(s, args...)
}

func addDebugColor(s string) string {
	if s == "" || s == ":" {
		return s
	}

	return "\033[33m"+s+"\033[0m"
}

func (s *Server) debugStatus() string {
	if !s.Config.DebugEnabled {
		return ""
	}
	pc, filename, line, _ := runtime.Caller(2)
	
	return fmt.Sprintf(" [%s] [%s:%d]", runtime.FuncForPC(pc).Name(), filename, line)	
}
