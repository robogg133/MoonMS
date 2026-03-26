package parser

import (
	"encoding/json"

	"github.com/robogg133/MoonMS/data"
)

func decodeAudio_BackgroundMusic(a json.RawMessage) (Attribute, error) {
	var v data.Audio_BackgroundMusic
	err := json.Unmarshal(a, &v)
	return &v, err
}

func decodeAudio_AmbientSounds(a json.RawMessage) (Attribute, error) {
	var v data.Audio_AmbientSounds
	err := json.Unmarshal(a, &v)
	return &v, err
}

func decodeAudio_MusicVolume(a json.RawMessage) (Attribute, error) {
	var v data.Audio_MusicVolume
	err := json.Unmarshal(a, &v)
	return &v, err
}
