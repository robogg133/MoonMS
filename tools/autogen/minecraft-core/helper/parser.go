package parser

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Attribute interface {
	ID() string
}

type AttributeDecoder func(json.RawMessage) (Attribute, error)

var AttributeDecoders = map[string]AttributeDecoder{
	"minecraft:audio/background_music": decodeAudio_BackgroundMusic,
	"minecraft:audio/ambient_sounds":   decodeAudio_AmbientSounds,
	"minecraft:audio/music_volume":     decodeAudio_MusicVolume,

	"minecraft:gameplay/snow_golem_melts":          decodeGameplay_SnowGolemMelts,
	"minecraft:gameplay/increased_fire_burnout":    decodeGameplay_IncreasedFireBurnout,
	"minecraft:gameplay/can_pillager_patrol_spawn": decodeGameplay_CanPillagerPatrolSpawn,

	"minecraft:visual/sky_color":              decodeVisual_SkyColor,
	"minecraft:visual/fog_color":              decodeVisual_FogColor,
	"minecraft:visual/water_fog_color":        decodeVisual_WaterFogColor,
	"minecraft:visual/ambient_particles":      decodeVisual_AmbientParticles,
	"minecraft:visual/water_fog_end_distance": decodeVisual_WaterFogEndDistance,

	"water_color":          decodeEffect_WaterColor,
	"foliage_color":        decodeEffect_FoliageColor,
	"dry_foliage_color":    decodeEffect_DryFoliageColor,
	"grass_color":          decodeEffect_GrassColor,
	"grass_color_modifier": decodeEffect_GrassColorModifier,
}

func ParseHexColor(s string) (uint32, error) {
	if len(s) != 7 || s[0] != '#' {
		return 0, fmt.Errorf("invalid color format: %q", s)
	}

	v, err := strconv.ParseUint(s[1:], 16, 32)
	if err != nil {
		return 0, err
	}

	return uint32(v), nil
}
