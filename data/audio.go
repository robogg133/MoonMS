package data

type BackgroundMusic struct {
	MaxDelay int32  `json:"max_delay"`
	MinDelay int32  `json:"min_delay"`
	Sound    string `json:"sound"`
}

type Audio_BackgroundMusic struct {
	Default    BackgroundMusic `json:"default"`
	Creative   BackgroundMusic `json:"creative"`
	UnderWater BackgroundMusic `json:"under_water"`
}

func (*Audio_BackgroundMusic) ID() string { return "minecraft:audio/background_music" }

//

type Audio_AmbientSounds struct {
	Additions struct {
		Sound      string  `json:"sound"`
		TickChance float32 `json:"tick_chance"`
	} `json:"additions"`

	Loop string `json:"loop"`

	Mood struct {
		BlockSearchExtent uint8   `json:"block_search_extent"`
		Offset            float32 `json:"offset"`
		Sound             string  `json:"sound"`
		TickDelay         float32 `json:"tick_delay"`
	} `json:"mood"`
}

func (*Audio_AmbientSounds) ID() string { return "minecraft:audio/ambient_sounds" }

//

type Audio_MusicVolume float32

func (*Audio_MusicVolume) ID() string { return "minecraft:audio/music_volume" }
