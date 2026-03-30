package app

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/robogg133/MoonMS/pkg/plugin"
)

var ErrSameIdentifier = errors.New("plugins with the same identifier, please check your plugins")

func (s *Server) InitPlugins() {

	allDirFiles, err := os.ReadDir(s.Config.PluginsFolder)
	if err != nil {
		s.LogError("%v", err)
	}
	plugin.SetPluginsFolder(s.Config.PluginsFolder)

	allPlgStart := time.Now()
	s.LogInfo("Loading all plugins asynchronously...")
	var wg sync.WaitGroup

	for _, d := range allDirFiles {
		if d.IsDir() {
			continue
		}

		path := filepath.Join(s.Config.PluginsFolder, d.Name())

		w := s.NewPluginWriter("")
		pl := plugin.NewPlugin(path, w)

		w.SetName(pl.Meta.Name)

		if pl.Meta.MCVersion != s.MinecraftConfig.MinecraftVersion {
			s.LogError("(plugin:%s) Mismatched version between server and plugin. Plugin version: %s, Server version: %s. Remove this plugin or try checking if it has received an update. Plugin homepage: \"%s\"", pl.ID, pl.Meta.MCVersion, s.MinecraftConfig.MinecraftVersion, pl.Meta.Homepage)
			continue
		}

		defer func() {
			r := recover()
			if r == nil {
				return
			}

			if r == ErrSameIdentifier {
				s.LogPanic("THERE ARE PLUGINS RUNNING WITH SAME IDENTIFIER CAN'T CONTINUE")
				runtime.StartTrace()
				panic(r.(error))
			}

			s.LogError("plugin (%s) failed to load: %v", pl.Meta.Name, r)
		}()

		_, exists := s.Plugins[pl.ID]
		if exists {
			panic(ErrSameIdentifier)
		}

		s.Plugins[pl.ID] = &pl

		s.loadWrapper(&pl, &wg)
	}

	wg.Wait()

	s.LogInfo("The server took %dms to load all plugins!", time.Since(allPlgStart).Milliseconds())

}

func (s *Server) loadWrapper(plg *plugin.Plugin, wg *sync.WaitGroup) {

	s.LogInfo("Loading plugin \"%s\"", plg.Meta.Name)
	start := time.Now()
	wg.Add(1)
	plg.Load()
	wg.Done()
	s.LogInfo("Plugin \"%s\" loaded in %dms", plg.Meta.Name, time.Since(start).Milliseconds())

}
