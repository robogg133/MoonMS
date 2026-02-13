package server

import (
	"fmt"
	"io"
	"os"
	"time"
)

var DebugEnabled bool

func LogInfo(args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf("%v", args...)
	fmt.Printf("[%s] \033[34m[INFO]:\033[0m %v\n", timestamp, message)
	fmt.Fprintf(LatestLogFile, "[%s] [INFO]: %v\n", timestamp, message)
}

func Debug(args ...any) {
	if !DebugEnabled {
		return
	}
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf("%v", args...)
	fmt.Printf("[%s] \033[33m[DEBUG]:\033[0m %v\n", timestamp, message)
	fmt.Fprintf(LatestLogFile, "[%s] [DEBUG]: %v\n", timestamp, message)
}

func LogWarn(args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf("%v", args...)
	fmt.Printf("[%s] \033[33m[WARN]:\033[0m %v\n", timestamp, message)
	fmt.Fprintf(LatestLogFile, "[%s] [WARN]: %v\n", timestamp, message)
}

func LogError(args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf("%v", args...)
	fmt.Printf("[%s] \033[31m[ERROR]:\033[0m %v\n", timestamp, message)
	fmt.Fprintf(LatestLogFile, "[%s] [ERROR]: %v\n", timestamp, message)
}

func LogFatal(args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf("%v", args...)
	fmt.Printf("[%s] \033[31m[FATAL]:\033[0m %v\n", timestamp, message)
	fmt.Fprintf(LatestLogFile, "[%s] [FATAL]: %v\n", timestamp, message)
	os.Exit(1)
}

func LogPlugin(pluginName string, args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf("%v", args...)
	fmt.Printf("[%s] \033[32m[%s]:\033[0m %v \n", timestamp, pluginName, message)
	fmt.Fprintf(LatestLogFile, "[%s] [%s]: %v\n", timestamp, pluginName, message)
}

type PluginWriter struct {
	w    io.Writer
	name string
}

func GetLogPluginWriter(w io.Writer, pluginName string) *PluginWriter {
	return &PluginWriter{
		w:    w,
		name: pluginName,
	}
}

func (p *PluginWriter) Write(b []byte) (int, error) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")

	fmt.Fprintf(LatestLogFile, "[%s] [%s]: %v", timestamp, p.name, string(b))
	return fmt.Fprintf(p.w, "[%s] \033[32m[%s]\033[0m: %v", timestamp, p.name, string(b))
}

type PluginWriterErr struct {
	w    io.Writer
	name string
}

func GetLogPluginWriterErr(w io.Writer, pluginName string) *PluginWriterErr {
	return &PluginWriterErr{
		w:    w,
		name: pluginName,
	}
}

func (p *PluginWriterErr) Write(b []byte) (int, error) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")

	fmt.Fprintf(LatestLogFile, "[%s] [%s] (ERROR): %v", timestamp, p.name, string(b))
	return fmt.Fprintf(p.w, "[%s] \033[31m[%s] (ERROR)\033[0m: %v", timestamp, p.name, string(b))
}
