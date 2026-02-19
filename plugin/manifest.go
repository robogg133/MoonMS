package plugin

import (
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
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Homepage    string `yaml:"home-page"`
	Description string `yaml:"description"`
	Authors     string `yaml:"authors"`

	ApiVersion string `yaml:"api-version"`

	Entry struct {
		Type string `yaml:"type"`
		File string `yaml:"file"`
	} `yaml:"entry"`

	Objects []string `yaml:"objects"`

	Require  []Dependencie `yaml:"require"`
	Provides []Dependencie `yaml:"provides"`
}

func ReadManifest(r io.Reader) Manifest {
	decoder := yaml.NewDecoder(r)

	var m Manifest

	if err := decoder.Decode(&m); err != nil {
		panic(err)
	}

	return m
}
