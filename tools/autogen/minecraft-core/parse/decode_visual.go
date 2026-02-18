package parser

import (
	"encoding/json"

	"github.com/robogg133/MoonMS/data"
)

func decodeVisual_SkyColor(a json.RawMessage) (Attribute, error) {
	var s string
	if err := json.Unmarshal(a, &s); err != nil {
		return nil, err
	}

	n, err := ParseHexColor(s)
	v := data.Visual_SkyColor(n)
	return &v, err
}

func decodeVisual_AmbientParticles(a json.RawMessage) (Attribute, error) {
	var v data.Visual_AmbientParticles
	err := json.Unmarshal(a, &v)
	return &v, err
}

func decodeVisual_FogColor(a json.RawMessage) (Attribute, error) {
	var s string
	if err := json.Unmarshal(a, &s); err != nil {
		return nil, err
	}

	n, err := ParseHexColor(s)
	v := data.Visual_FogColor(n)
	return &v, err
}

func decodeVisual_WaterFogColor(a json.RawMessage) (Attribute, error) {
	var s string
	if err := json.Unmarshal(a, &s); err != nil {
		return nil, err
	}

	n, err := ParseHexColor(s)
	v := data.Visual_WaterFogColor(n)
	return &v, err
}

func decodeVisual_WaterFogEndDistance(a json.RawMessage) (Attribute, error) {
	var v data.Visual_WaterFogEndDistance
	err := json.Unmarshal(a, &v)
	return &v, err
}
