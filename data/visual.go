package data

type Particle struct {
	Particle struct {
		Type string `json:"type"`
	}
	Probability float32 `json:"probability"`
}

type Visual_SkyColor uint32

func (*Visual_SkyColor) ID() string { return "minecraft:visual/sky_color" }

type Visual_AmbientParticles []Particle

func (*Visual_AmbientParticles) ID() string { return "minecraft:visual/ambient_particles" }

type Visual_FogColor uint32

func (*Visual_FogColor) ID() string { return "minecraft:visual/fog_color" }

type Visual_WaterFogColor uint32

func (*Visual_WaterFogColor) ID() string { return "minecraft:visual/water_fog_color" }

type Visual_WaterFogEndDistance struct {
	Argument float32 `json:"argument"`
	Modifier string  `json:"modifier"`
}

func (*Visual_WaterFogEndDistance) ID() string { return "minecraft:visual/water_fog_end_distance" }
