package plugins

import (
	"github.com/Shopify/go-lua"
)

const MANIFEST_FILE_NAME = "plugin-manifest.yml"
const MAIN_LUA_FILE_NAME = "init.lua"

type PluginManifest struct {
	APIVersion string `yaml:"api-version"`

	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	Authors     string `yaml:"authors"`
	Website     string `yaml:"website"`

	Dependencies []PluginDependencie `yaml:"dependencies"`
}

type PluginDependencie struct {
	Name       string `yaml:"name"`
	LoadBefore bool   `yaml:"load-before"`
	Required   bool   `yaml:"required"`
}

type Plugin struct {
	Identifier string
	Filepath   string
	Manifest   PluginManifest

	MainLuaFile string
	LuaVM       *lua.State
}

var AllPlugins []Plugin

func GetAllPlugins() *[]Plugin { return &AllPlugins }
