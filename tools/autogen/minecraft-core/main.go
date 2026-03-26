//go:generate go run .

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type LockFile struct {
	Version string `json:"mc_version"`
	Path    string `json:"path"`
}

func main() {

	var a string

	for {
		var err error
		a, err = filepath.Abs(".")
		if err != nil {
			panic(err)
		}

		if _, err := os.Stat("go.mod"); err != nil {
			os.Chdir("..")
			continue
		}
		break

	}

	startingDir, err := filepath.Abs(a)
	if err != nil {
		panic(err)
	}
	fmt.Println(startingDir)

	var releaseName string

	flag.StringVar(&releaseName, "release", "latest", "Minecraft release version to download the core datapack")
	flag.Parse()

	var folder string

	b, err := os.ReadFile("core_datapack.lock")
	if err != nil {
		if os.IsNotExist(err) {
			folder = extract(&releaseName)
			os.Chdir(startingDir)

			var lf LockFile

			lf.Path = folder
			lf.Version = releaseName

			b, err := json.Marshal(lf)
			if err != nil {
				panic(err)
			}

			if err := os.WriteFile("core_datapack.lock", b, 0777); err != nil {
				panic(err)
			}

		} else {
			panic(err)
		}
	} else {

		var lf LockFile
		if err := json.Unmarshal(b, &lf); err != nil {
			panic(err)
		}
		folder = lf.Path
		releaseName = lf.Version
	}

	os.Chdir(startingDir)

	scriptsFolder := "./tools/autogen/minecraft-core/scripts/"

	if err := exec.Command("go", "generate", scriptsFolder+"...").Run(); err != nil {
		panic(err)
	}

	switch flag.Arg(0) {
	case "biome":
		execFile(filepath.Join(scriptsFolder, "biome"), []string{filepath.Join(startingDir, "internal", "gen", "core", "worldgen"), filepath.Join(folder, "worldgen", "biome"), releaseName})
	case "damage_type":
		execFile(filepath.Join(scriptsFolder, "damage"), []string{filepath.Join(startingDir, "internal", "gen", "core", "damage"), filepath.Join(folder, "damage_type"), releaseName})
	default:
		execFile(filepath.Join(scriptsFolder, "biome"), []string{filepath.Join(startingDir, "internal", "gen", "core", "worldgen"), filepath.Join(folder, "worldgen", "biome"), releaseName})
		execFile(filepath.Join(scriptsFolder, "damage"), []string{filepath.Join(startingDir, "internal", "gen", "core", "damage"), filepath.Join(folder, "damage_type"), releaseName})
	}

	fmt.Println("end")
}

func execFile(s string, args []string) {

	fmt.Printf("Running %s\n", s)

	execPath := filepath.Join(s, "exec")

	cmd := exec.Command(execPath, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()

	if err != nil {
		panic(err)
	}
}
