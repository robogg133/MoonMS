package plugins

import (
	plugins_lua "MoonMS/internal/plugins/lua"
	plugins_wasmi "MoonMS/internal/plugins/wasmi"
	"MoonMS/internal/server"
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"

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
	defer f.Close()

	luaFile, err := io.ReadAll(f)
	if err != nil {
		return Plugin{}, err
	}

	pluginObj.MainFileContent = luaFile

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

	AllPlugins = append(AllPlugins, *p)
	return nil
}
