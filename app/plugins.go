package app

import (
	"MoonMS/internal/plugins"
	"fmt"
	"os"
	"path/filepath"
)

func (s *Server) InitPlugins() {

	_ = os.Mkdir("plugins", 0755)

	allDirFiles, err := os.ReadDir(s.Config.PluginsFolder)
	if err != nil {
		s.LogError(err)
	}

	for _, d := range allDirFiles {
		if d.IsDir() {
			continue
		}

		path := filepath.Join("plugins", d.Name())

		f, err := os.Open(path)
		if err != nil {
			s.LogError(fmt.Sprintf("Error opening file: %v", err))
			s.LogInfo(fmt.Sprintf("SKIPPING %v", d.Name()))
			continue
		}

		stat, err := f.Stat()
		if err != nil {
			s.LogError(fmt.Sprintf("Error getting file status: %v", err))
			s.LogInfo(fmt.Sprintf("SKIPPING %v", d.Name()))
			continue
		}

		plugin, err := plugins.ReadPluginFile(f, stat.Size(), path)
		if err != nil {
			s.LogError(fmt.Sprintf("Error parsing plugin: %v", err))
			s.LogInfo(fmt.Sprintf("SKIPPING %v", d.Name()))
			continue
		}

		s.LogInfo(fmt.Sprintf("Starting %s", plugin.Identifier))
		s.Plugins[plugin.Identifier] = plugin
		go plugin.LoadPlugin()
	}
}
