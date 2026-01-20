package server

import (
	"fmt"
	"os"
	"time"
)

func LogInfo(args ...any) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf("%v", args...)
	fmt.Printf("[%s] \033[34m[INFO]:\033[0m %v\n", timestamp, message)
	fmt.Fprintf(LatestLogFile, "[%s] [INFO]: %v\n", timestamp, message)
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
