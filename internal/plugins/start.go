package plugins

import (
	plugins_lua "MoonMS/internal/plugins/lua"
	plugins_wasmi "MoonMS/internal/plugins/wasmi"
	"MoonMS/internal/server"
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Shopify/go-lua"
	"go.yaml.in/yaml/v4"
)

const MANIFEST_FILE_NAME = "plugin-manifest.yml"

type PluginManifest struct {
	APIVersion string `yaml:"api-version"`

	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	Authors     string `yaml:"authors"`
	Website     string `yaml:"website"`

	Runtime string `yaml:"runtime"` // lua | wasm

	Entry string `yaml:"entry"` // main.lua | plugin.wasm

	Provides []PluginDependencie `yaml:"provides"`
	Requires []PluginDependencie `yaml:"requires"`
}

type PluginDependencie struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type Plugin struct {
	Identifier string
	Filepath   string
	Manifest   PluginManifest

	MainFileContent []byte
	Runtime         any
}

func ReadPluginFile(r io.ReaderAt, size int64, path string) (Plugin, error) {
	reader, err := zip.NewReader(r, size)
	if err != nil {
		return Plugin{}, err
	}

	file, err := reader.Open(MANIFEST_FILE_NAME)
	if err != nil {
		return Plugin{}, err
	}
	defer file.Close()

	yamlBuffer, err := io.ReadAll(file)
	if err != nil {
		return Plugin{}, err
	}
	var manifest PluginManifest

	if err := yaml.Unmarshal(yamlBuffer, &manifest); err != nil {
		return Plugin{}, err
	}

	var pluginObj Plugin

	pluginObj.Manifest = manifest
	pluginObj.Identifier = filepath.Base(path)
	pluginObj.Filepath = path

	if manifest.Runtime != "wasm" && manifest.Runtime != "lua" {
		return Plugin{}, fmt.Errorf("invalid runtime")
	}

	if manifest.APIVersion != server.GetServerData().MINECRAFT_VERSION {
		return Plugin{}, fmt.Errorf("mismatched api version")
	}

	f, err := reader.Open(manifest.Entry)
	if err != nil {
		return Plugin{}, err
	}

	mainFile, err := io.ReadAll(f)
	if err != nil {
		f.Close()
		return Plugin{}, err
	}
	f.Close()

	pluginObj.MainFileContent = mainFile

	thisPluginDir := filepath.Join("plugins", manifest.Name)
	_, err = os.Stat(thisPluginDir)
	if err == nil {
		return pluginObj, nil
	}
	if os.IsNotExist(err) {

		if err := os.Mkdir(thisPluginDir, 0755); err != nil {
			return Plugin{}, err
		}
		for _, f := range reader.File {
			base := filepath.Base(f.Name)
			switch {
			case strings.HasSuffix(base, ".wasm"):
				continue
			case strings.HasSuffix(base, ".lua"):
				continue
			case strings.HasSuffix(base, ".so"):
				continue
			case strings.HasSuffix(base, manifest.Entry):
				continue
			case base == MANIFEST_FILE_NAME:
				continue
			}

			fileName := filepath.Join(thisPluginDir, f.Name)

			if strings.Contains(filepath.Dir(f.Name), "/data/") {
				fileName = filepath.Join(thisPluginDir, strings.Replace(f.Name, "/data/", "/", 1))
			}

			if f.FileInfo().IsDir() {
				if err := os.Mkdir(fileName, f.Mode().Perm()); err != nil {
					return Plugin{}, err
				}
				continue
			}
			unzipedFile, err := f.Open()
			if err != nil {
				return Plugin{}, err
			}

			newFile, err := os.Create(fileName)
			if err != nil {
				unzipedFile.Close()
				return Plugin{}, err
			}

			io.Copy(newFile, unzipedFile)
			newFile.Close()
			unzipedFile.Close()
		}
		os.RemoveAll(filepath.Join(thisPluginDir, "data"))
	} else {
		return Plugin{}, err
	}

	return pluginObj, nil
}

func (p *Plugin) LoadPlugin() error {

	switch p.Manifest.Runtime {
	case "lua":
		p.Runtime = lua.NewState()

		if err := plugins_lua.LoadPlugin(p.Runtime.(*lua.State), string(p.MainFileContent), p.Manifest.Name); err != nil {
			return err
		}
	case "wasm":
		p.Runtime = plugins_wasmi.CreateRuntime(p.Manifest.Name)
		if err := plugins_wasmi.RunFile(p.MainFileContent, p.Runtime.(plugins_wasmi.Runtime)); err != nil {
			return err
		}
	}

	return nil
}
