package parser

import (
	"encoding/json"

	"github.com/robogg133/MoonMS/data"
)

func decodeEffect_WaterColor(a json.RawMessage) (Attribute, error) {
	var s string
	if err := json.Unmarshal(a, &s); err != nil {
		return nil, err
	}

	n, err := ParseHexColor(s)
	v := data.Effect_WaterColor(n)
	return &v, err
}

func decodeEffect_FoliageColor(a json.RawMessage) (Attribute, error) {
	var s string
	if err := json.Unmarshal(a, &s); err != nil {
		return nil, err
	}

	n, err := ParseHexColor(s)
	v := data.Effect_FoliageColor(n)
	return &v, err
}

func decodeEffect_DryFoliageColor(a json.RawMessage) (Attribute, error) {
	var s string
	if err := json.Unmarshal(a, &s); err != nil {
		return nil, err
	}

	n, err := ParseHexColor(s)
	v := data.Effect_DryFoliageColor(n)
	return &v, err
}

func decodeEffect_GrassColor(a json.RawMessage) (Attribute, error) {
	var s string
	if err := json.Unmarshal(a, &s); err != nil {
		return nil, err
	}

	n, err := ParseHexColor(s)
	v := data.Effect_GrassColor(n)
	return &v, err
}

func decodeEffect_GrassColorModifier(a json.RawMessage) (Attribute, error) {
	var s string
	err := json.Unmarshal(a, &s)

	v := data.Effect_GrassColorModifier(s)
	return &v, err
}
