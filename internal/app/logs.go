package app

import (
	"fmt"
	"time"
)

func (s *Server) LogInfo(args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf("%v", args...)
	fmt.Printf("[%s] \033[34m[INFO]:\033[0m %v\n", timestamp, message)
	fmt.Fprintf(s.Config.LatestLogFile, "[%s] [INFO]: %v\n", timestamp, message)
}

func (s *Server) LogDebug(args ...any) {
	if !s.Config.DebugEnabled {
		return
	}
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf("%v", args...)
	fmt.Printf("[%s] \033[33m[DEBUG]:\033[0m %v\n", timestamp, message)
	fmt.Fprintf(s.Config.LatestLogFile, "[%s] [DEBUG]: %v\n", timestamp, message)
}

func (s *Server) LogWarn(args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf("%v", args...)
	fmt.Printf("[%s] \033[33m[WARN]:\033[0m %v\n", timestamp, message)
	fmt.Fprintf(s.Config.LatestLogFile, "[%s] [WARN]: %v\n", timestamp, message)
}

func (s *Server) LogError(args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf("%v", args...)
	fmt.Printf("[%s] \033[31m[ERROR]:\033[0m %v\n", timestamp, message)
	fmt.Fprintf(s.Config.LatestLogFile, "[%s] [ERROR]: %v\n", timestamp, message)
}

func (s *Server) LogPanic(args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf("%v", args...)
	fmt.Printf("[%s] \033[31m[PANIC]:\033[0m %v\n", timestamp, message)
	fmt.Fprintf(s.Config.LatestLogFile, "[%s] [PANIC]: %v\n", timestamp, message)
}

func (s *Server) LogPlugin(pluginName string, args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf("%v", args...)
	fmt.Printf("[%s] \033[32m[%s]:\033[0m %v \n", timestamp, pluginName, message)
	fmt.Fprintf(s.Config.LatestLogFile, "[%s] [%s]: %v\n", timestamp, pluginName, message)
}
