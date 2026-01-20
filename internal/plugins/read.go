package plugins

import (
	"archive/zip"
	"fmt"
	"io"
	"strings"

	"github.com/Shopify/go-lua"
	"go.yaml.in/yaml/v4"
)

func ReadPluginFile(r io.ReaderAt, size int64, filepath string) (Plugin, error) {
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
	pluginObj.Identifier = fmt.Sprintf("plugin:%s", strings.ToLower(manifest.Name))
	pluginObj.Filepath = filepath

	f, err := reader.Open(MAIN_LUA_FILE_NAME)
	if err != nil {
		return Plugin{}, err
	}
	defer f.Close()

	luaFile, err := io.ReadAll(f)
	if err != nil {
		return Plugin{}, err
	}

	pluginObj.MainLuaFile = string(luaFile)

	pluginObj.LuaVM = lua.NewState()

	return pluginObj, nil
}
