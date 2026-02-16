package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {

	startingDir, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	var releaseName string

	flag.StringVar(&releaseName, "release", "latest", "Minecraft release version to download the core datapack")
	flag.Parse()

	fmt.Println("[1/8] Downloading latest release server.jar")
	blob := downloadServerJar(&releaseName)

	tempDir, err := os.MkdirTemp("", "minecraft-server")
	if err != nil {
		panic(err)
	}
	fmt.Println("[2/8] Create temp dir")
	f, err := os.CreateTemp(tempDir, "server.jar")
	if err != nil {
		panic(err)
	}

	if _, err := f.Write(blob); err != nil {
		panic(err)
	}
	fmt.Println("[3/8] Create jar file")

	if err := os.WriteFile(filepath.Join(tempDir, "eula.txt"), []byte("eula=false"), 0777); err != nil {
		panic(err)
	}

	os.Chdir(tempDir)
	cmd := exec.Command("java", "-jar", f.Name())
	fmt.Println("[4/8] Starting minecraft server")
	if err := cmd.Start(); err != nil {
		panic(err)
	}

	_, err = cmd.Process.Wait()
	if err != nil {
		panic(err)
	}
	fmt.Println("[5/8] Minecraft server done")

	reader, err := zip.OpenReader(filepath.Join(tempDir, "versions", releaseName, fmt.Sprintf("server-%s.jar", releaseName)))
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	prefixDir := "minecraft_core_datapack"
	os.Mkdir(prefixDir, 0777)

	for _, f := range reader.File {
		if !strings.HasPrefix(f.Name, "data/minecraft") {
			continue
		}
		f.Name = strings.TrimPrefix(f.Name, "data/minecraft")
		if f.Mode().IsDir() {
			os.MkdirAll(filepath.Join(prefixDir, f.Name), 0777)
			continue
		}

		os.MkdirAll(filepath.Join(prefixDir, filepath.Dir(f.Name)), 0777)
		tf, err := os.OpenFile(filepath.Join(prefixDir, f.Name), os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			panic(err)
		}

		inputFile, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(tf, inputFile); err != nil {
			panic(err)
		}

	}

	fmt.Println(filepath.Join(tempDir, prefixDir))
	fmt.Println(startingDir)
}
