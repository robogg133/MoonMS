package plugin

import (
	"fmt"
	"io"

	"go.yaml.in/yaml/v4"
)

const MANIFEST_FILE_NAME string = "plugin-manifest.yml"

type Dependencie struct {
	Capabilty string `yaml:"capability"`
	Version   string `yaml:"version"`
}

type Manifest struct {
	Identifier  string `yaml:"identifier"`
	PluginKind  string `yaml:"kind"`
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Homepage    string `yaml:"home-page"`
	Description string `yaml:"description"`
	Authors     string `yaml:"authors"`

	MCVersion string `yaml:"mc-version"`

	Require  []Dependencie `yaml:"require"`
	Provides []Dependencie `yaml:"provides"`
}

func (m *Manifest) decode(r io.Reader) error {
	if err := yaml.NewDecoder(r).Decode(&m); err != nil {
		return err
	}
	if !validIdentifier(m.Identifier) {
		return fmt.Errorf("Invalid identifier: can only be alphanumeric and have _")
	}
	return nil
}
