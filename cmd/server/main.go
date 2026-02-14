package main

import (
	"MoonMS/app"
	"MoonMS/internal/packets"
	"crypto/rand"
	"crypto/rsa"
	"os"
	"os/signal"
	"time"

	_ "embed"
)

const DEADLINE = time.Second * 30

type MojangAnswer struct {
	Properties []map[string]string
}

var AnonymousPlayer = &packets.PlayerMinimunInfo{Username: "Anonymous Player", UUID: "00000000-0000-0000-0000-000000000000"}

func main() {

	cfg := app.MinecraftServerConfig{}
	if err := cfg.ConfigFile(); err != nil {
		panic(err)
	}

	scfg := app.Config{
		LatestLogFile: "logs/latest.log",
		StartName:     "java",
		DebugEnabled:  false,
		PluginsFolder: "plugins",
	}

	if os.Getenv("DEBUG") == "true" {
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

		server.Stop()

		os.Exit(0)
	}()

	server.Start()
}
