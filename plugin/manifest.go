package plugin

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

	Require  []Dependencie `yaml:"require"`
	Provides []Dependencie `yaml:"provides"`
}
