package app

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func (s *Server) LogInfo(args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf("%v", args...)
	fmt.Printf("[%s] \033[34m[INFO]:\033[0m %v\n", timestamp, message)
	fmt.Fprintf(s.logFile, "[%s] [INFO]: %v\n", timestamp, message)
}

func (s *Server) LogDebug(args ...any) {
	if !s.Config.DebugEnabled {
		return
	}
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf("%v", args...)
	fmt.Printf("[%s] \033[33m[DEBUG]:\033[0m %v\n", timestamp, message)
	fmt.Fprintf(s.logFile, "[%s] [DEBUG]: %v\n", timestamp, message)
}

func (s *Server) LogWarn(args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf("%v", args...)
	fmt.Printf("[%s] \033[33m[WARN]:\033[0m %v\n", timestamp, message)
	fmt.Fprintf(s.logFile, "[%s] [WARN]: %v\n", timestamp, message)
}

func (s *Server) LogError(args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf("%v", args...)
	fmt.Printf("[%s] \033[31m[ERROR]:\033[0m %v\n", timestamp, message)
	fmt.Fprintf(s.logFile, "[%s] [ERROR]: %v\n", timestamp, message)
}

func (s *Server) LogPanic(args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf("%v", args...)
	fmt.Printf("[%s] \033[31m[PANIC]:\033[0m %v\n", timestamp, message)
	fmt.Fprintf(s.logFile, "[%s] [PANIC]: %v\n", timestamp, message)
}

func (s *Server) LogPlugin(pluginName string, args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf("%v", args...)
	fmt.Printf("[%s] \033[32m[%s]:\033[0m %v \n", timestamp, pluginName, message)
	fmt.Fprintf(s.logFile, "[%s] [%s]: %v\n", timestamp, pluginName, message)
}

func (s *Server) StartLogger() error {

	if err := enable(); err != nil {
		return err
	}

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

	writer := gzip.NewWriter(f)

	_, err = io.Copy(writer, oldLog)
	if err != nil {
		return err
	}

	return os.Remove(s.Config.LatestLogFile)
}
