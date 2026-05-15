package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func extract(releaseName *string) (dataDir string) {

	fmt.Println("[1/8] Downloading latest release server.jar")
	blob := downloadServerJar(releaseName)

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

	reader, err := zip.OpenReader(filepath.Join(tempDir, "versions", *releaseName, fmt.Sprintf("server-%s.jar", *releaseName)))
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	prefixDir := "minecraft_core_datapack"
	langPrefix := "minecraft_lang_codes"
	os.Mkdir(prefixDir, 0777)

	coreDatapackTrigger := filepath.Join("data", "minecraft")
	langCodesTrigger := filepath.Join("assets", "minecraft", "lang")
	for _, f := range reader.File {
		var currentPrefix, currentTrigger string

		if strings.HasPrefix(f.Name, coreDatapackTrigger) {
			currentTrigger = coreDatapackTrigger
			currentPrefix = prefixDir
		} else if strings.HasPrefix(f.Name, langCodesTrigger) {
			currentTrigger = langCodesTrigger
			currentPrefix = langPrefix
		} else {
			continue
		}

		f.Name = strings.TrimPrefix(f.Name, currentTrigger)
		if f.Mode().IsDir() {
			os.MkdirAll(filepath.Join(currentPrefix, f.Name), 0777)
			continue
		}

		os.MkdirAll(filepath.Join(currentPrefix, filepath.Dir(f.Name)), 0777)
		tf, err := os.OpenFile(filepath.Join(currentPrefix, f.Name), os.O_CREATE|os.O_WRONLY, 0777)
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

	return filepath.Join(tempDir, prefixDir)
}
