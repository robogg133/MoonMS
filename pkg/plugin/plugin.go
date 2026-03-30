package plugin

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/robogg133/MoonMS/internal/shared"
	"github.com/robogg133/MoonMS/pkg/plugin/wasm"
)

type State uint8

type Plugin struct {
	ID       string
	Meta     Manifest
	Runtime  Runtime
	MyFolder string

	State State
	Caps  []Capability
}

type Capability struct {
	Identifier string
	Version    string
}

type Runtime interface {
	Load() error

	Pause()
	Tick(deadline time.Time) error
	Resume()

	Call(string, ...any) error

	Close() error
}

const (
	StateLoaded State = iota
	StateEnabled
	StateDisabled
	StateCrashed
)

func NewPlugin(path string, logWriter io.Writer) Plugin {
	var pl Plugin

	reader, err := zip.OpenReader(path)
	if err != nil {
		panic(err)
	}

	f, err := reader.Open(MANIFEST_FILE_NAME)
	if err != nil {
		panic(err)
	}
	pl.Meta = ReadManifest(f)

	pl.MyFolder = filepath.Join(pluginsFolder, pl.Meta.Name)

	pl.ID = pl.Meta.Identifier

	if _, err = os.Stat(pl.MyFolder); err == nil {
		for _, v := range reader.File {
			if v.Name == pl.Meta.Entry.File {
				pl.copyWithPrefix(v, ".objects")
			}
			if slices.Contains(pl.Meta.Objects, v.Name) {
				pl.copyWithPrefix(v, ".objects")
			}
		}
		pl.initRuntime(logWriter)
		return pl
	}
	pl.mkdirPluginFolderStructure()

	for _, v := range reader.File {
		if v.Name == pl.Meta.Entry.File {
			pl.copyWithPrefix(v, ".objects")
		}
		if slices.Contains(pl.Meta.Objects, v.Name) {
			pl.copyWithPrefix(v, ".objects")
		}
		var found bool
		if v.Name, found = strings.CutPrefix(v.Name, "data/"); found {
			pl.copyWithPrefix(v, ".data")
		}
	}

	pl.initRuntime(logWriter)

	return pl
}

func (pl *Plugin) mkdirPluginFolderStructure() {
	if err := os.MkdirAll(pl.MyFolder, 0755); err != nil {
		panic(err)
	}

	if err := os.Mkdir(filepath.Join(pl.MyFolder, ".objects"), 0755); err != nil {
		panic(err)
	}
	shared.SetHidden(filepath.Join(pl.MyFolder, ".objects"))

	if err := os.Mkdir(filepath.Join(pl.MyFolder, ".data"), 0755); err != nil {
		panic(err)
	}
	shared.SetHidden(filepath.Join(pl.MyFolder, ".data"))

}

func (pl *Plugin) copyWithPrefix(v *zip.File, prefix string) {
	if v.FileInfo().IsDir() {
		err := os.MkdirAll(filepath.Join(filepath.Join(pl.MyFolder, prefix, v.Name)), 0755)
		if err != nil {
			panic(err)
		}
	}

	sF, err := v.Open()
	if err != nil {
		panic(err)
	}
	tF, err := os.Create(filepath.Join(pl.MyFolder, prefix, v.Name))
	if err != nil {
		panic(err)
	}

	if _, err := io.Copy(tF, sF); err != nil {
		panic(err)
	}

}

func (pl *Plugin) initRuntime(logWriter io.Writer) {
	if pl.Meta.Entry.Type == "wasm" {
		b, err := os.ReadFile(filepath.Join(pl.MyFolder, ".objects", pl.Meta.Entry.File))
		if err != nil {
			panic(err)
		}
		pl.Runtime = wasm.NewRuntime(logWriter, b, os.DirFS(pl.MyFolder))
	}
}
