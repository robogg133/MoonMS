package plugin

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
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
	StatePrepared
	StateEnabled
	StateDisabled
	StateCrashed
)

func consultZipOkAndUpdate(jsonFile, identifier, srcPath string) (bool, error) {
	f, err := os.Open(srcPath)
	if err != nil {
		return false, err
	}

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		f.Close()
		return false, err
	}
	digest := hex.EncodeToString(h.Sum(nil))
	f.Close()

	jsonData := make(map[string]string)

	jf, err := os.OpenFile(jsonFile, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return false, err
	}
	defer func() {
		defer jf.Close()
		if _, err := jf.Seek(0, io.SeekStart); err != nil {
			panic(err)
		}

		jsonData[identifier] = digest
		if err := json.NewEncoder(jf).Encode(jsonData); err != nil {
			panic(err)
		}
	}()

	if err := json.NewDecoder(jf).Decode(&jsonData); err != nil {
		return false, err
	}

	if s, ok := jsonData[identifier]; ok && s == digest {
		return true, nil
	}
	return false, nil
}

func parseManifest(path string) (Manifest, error) {
	var m Manifest

	r, err := zip.OpenReader(path)
	if err != nil {
		return Manifest{}, err
	}
	defer r.Close()

	f, err := r.Open(MANIFEST_FILE_NAME)
	if err != nil {
		return Manifest{}, err
	}
	defer f.Close()

	if err := m.decode(f); err != nil {
		return Manifest{}, err
	}

	return m, nil
}

func PreparePlugin(path string, pluginsShaFile string) (*Plugin, error) {

	m, err := parseManifest(path)
	if err != nil {
		if err.Error() == zip.ErrFormat.Error() {
			return nil, fmt.Errorf("invalid file: %s on plugins folder NOT a plugin!", path)
		}
		return nil, err
	}

	ok, err := consultZipOkAndUpdate(pluginsShaFile, m.Identifier, path)
	if err != nil {
		return nil, err
	}
	destDir := filepath.Join(filepath.Dir(pluginsShaFile), "plg!!"+m.Identifier)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, err
	}
	if !ok {
		reader, err := zip.OpenReader(path)
		if err != nil {
			reader.Close()
			return nil, err
		}

		for _, fd := range reader.File {
			if fd.Name == MANIFEST_FILE_NAME {
				continue
			}
			if err := copyWithPrefix(fd, destDir); err != nil {
				reader.Close()
				return nil, err
			}
		}
		reader.Close()
	}
	pl := &Plugin{
		Meta:     m,
		State:    StatePrepared,
		MyFolder: destDir,
	}
	return pl, nil
}

func copyWithPrefix(v *zip.File, prefix string) error {
	if v.FileInfo().IsDir() {
		err := os.MkdirAll(filepath.Join(filepath.Join(prefix, v.Name)), 0755)
		if err != nil {
			return err
		}
	}

	sF, err := v.Open()
	if err != nil {
		return err
	}
	defer sF.Close()
	tF, err := os.Create(filepath.Join(prefix, v.Name))
	if err != nil {
		return err
	}

	if _, err := io.Copy(tF, sF); err != nil {
		return err
	}
	return nil
}

func (pl *Plugin) initRuntime(logWriter io.Writer) {

}
