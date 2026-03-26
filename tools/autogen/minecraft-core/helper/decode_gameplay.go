package parser

import (
	"encoding/json"

	"github.com/robogg133/MoonMS/data"
)

func decodeGameplay_SnowGolemMelts(a json.RawMessage) (Attribute, error) {
	var v data.Gameplay_SnowGolemMelts
	err := json.Unmarshal(a, &v)
	return &v, err
}

func decodeGameplay_IncreasedFireBurnout(a json.RawMessage) (Attribute, error) {
	var v data.Gameplay_IncreasedFireBurnout
	err := json.Unmarshal(a, &v)
	return &v, err
}

func decodeGameplay_CanPillagerPatrolSpawn(a json.RawMessage) (Attribute, error) {
	var v data.Gameplay_CanPillagerPatrolSpawn
	err := json.Unmarshal(a, &v)
	return &v, err
}
