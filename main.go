package main

import (
	"crypto/rand"
	"crypto/rsa"
	"os"
	"os/signal"

	"github.com/robogg133/MoonMS/app"

	_ "embed"
)

func main() {

	cfg := app.MinecraftServerConfig{}
	if err := cfg.ConfigFile(); err != nil {
		panic(err)
	}

	cfg.MinecraftVersion = "26.1.2"
	cfg.ProtcolVersion = 775

	scfg := app.Config{
		LatestLogFile:  "logs/latest.log",
		StartName:      "java",
		DebugEnabled:   false,
		PluginsFolder:  "plugins",
		AcessFolder:    "data/access",
		DatabaseFolder: "data/database",
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

		err := server.Stop()
		if err != nil {
			server.LogError("%v", err)
		}
		os.Exit(0)
	}()

	server.Start()
}
