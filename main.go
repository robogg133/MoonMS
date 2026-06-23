package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"slices"

	"github.com/robogg133/MoonMS/app"

	_ "embed"
)

var AdittionalVersionData = "(Developer Build)"

func main() {

	if slices.Contains(os.Args, "--version") {
		fmt.Printf("MoonMS r%d-mc%s %s/%s %s\n", 1, "26.1.2", runtime.GOOS, runtime.GOARCH, AdittionalVersionData)
		return
	}

	cfg := app.MinecraftServerConfig{}
	if err := cfg.ConfigFile(); err != nil {
		panic(err)
	}

	cfg.MinecraftVersion = "26.2"
	cfg.ProtcolVersion = 776

	scfg := app.Config{
		LatestLogFile:  "logs/latest.log",
		StartName:      "java",
		DebugEnabled:   false,
		PluginsFolder:  "plugins",
		AcessFolder:    "data/access",
		DatabaseFolder: "data/database",
	}

	if os.Getenv("DEBUG") == "1" {
		scfg.DebugEnabled = true
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, int(cfg.Advanced.RSAKeyBitAmmount))
	if err != nil {
		panic(err)
	}

	server := app.New(cfg, scfg, privateKey)
	if err := server.StartLogger(); err != nil {
		panic(err)
	}

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig

		err := server.Stop()
		if err != nil {
			server.LogError("%v", err)
		}
		os.Exit(0)
	}()

	server.Start()
}
