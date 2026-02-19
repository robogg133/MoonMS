package plugin

var pluginsFolder string = "plugins"

func SetPluginsFolder(folder string) {
	pluginsFolder = folder
}

func (p *Plugin) Load() {

	err := p.Runtime.Load()
	if err != nil {
		p.State = StateCrashed
		panic(err)
	}
	p.State = StateLoaded

}
