package plugin

type State uint8

type Plugin struct {
	ID      string
	Meta    Manifest
	Runtime Runtime

	State State
	//Caps  []Capability
}

type Runtime interface {
	Load() error
	Enable() error
	Disable() error
	Close() error
}

const (
	StateLoaded State = iota
	StateEnabled
	StateDisabled
	StateCrashed
)
