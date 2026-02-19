package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/robogg133/MoonMS/plugin"
)

func (s *Server) InitPlugins() {

	allDirFiles, err := os.ReadDir(s.Config.PluginsFolder)
	if err != nil {
		s.LogError(err)
	}
	plugin.SetPluginsFolder(s.Config.PluginsFolder)

	for _, d := range allDirFiles {
		if d.IsDir() {
			continue
		}

		path := filepath.Join(s.Config.PluginsFolder, d.Name())

		pl := plugin.NewPlugin(path)

		s.LogInfo(fmt.Sprintf("Starting %s", pl.Meta.Identifier))

		defer func() {
			r := recover()
			s.LogError(fmt.Sprintf("plugin %s failed to load: %v", pl.Meta.Name, r))
		}()

		go pl.Load()
	}
}
