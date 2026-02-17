package worldgen

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/robogg133/KernelCraft/data"
)

type RawBiome struct {
	Attributes  map[string]json.RawMessage `json:"attributes"`
	Carvers     json.RawMessage            `json:"carvers"`
	Downfall    float32                    `json:"downfall"`
	Effects     map[string]json.RawMessage `json:"effects"`
	Features    [][]string                 `json:"features"`
	HasPrecip   bool                       `json:"has_precipitation"`
	SpawnCosts  map[string]RawCost         `json:"spawn_costs"`
	Spawners    RawSpawners                `json:"spawners"`
	Temperature float32                    `json:"temperature"`
}

type RawCost struct {
	Charge       float32 `json:"charge"`
	EnergyBudget float32 `json:"energy_budget"`
}

type RawMob struct {
	Type     string `json:"type"`
	MaxCount uint8  `json:"maxCount"`
	MinCount uint8  `json:"minCount"`

	Weight uint8 `json:"weight"`
}

type RawSpawners struct {
	Ambient                  []RawMob `json:"ambient"`
	Axolotls                 []RawMob `json:"axolotls"`
	Creature                 []RawMob `json:"creature"`
	Misc                     []RawMob `json:"misc"`
	Monster                  []RawMob `json:"monster"`
	UndergroundWaterCreature []RawMob `json:"underground_water_creature"`
	WaterAmbient             []RawMob `json:"water_ambient"`
	WaterCreature            []RawMob `json:"water_creature"`
}

type BiomeFileWriter struct {
	w           io.Writer
	knownBiomes []string
}

type Biome struct {
	Carvers  []string
	Downfall float32

	Features [][]string

	HasPreciptation bool

	Spawners Spawners

	Temperature float32

	SpawnCosts map[data.EntityID]Costs

	AmbientSounds   *data.Audio_AmbientSounds
	BackgorundMusic *data.Audio_BackgroundMusic
	MusicVolume     float32

	AmbientParticles    *data.Visual_AmbientParticles
	WaterFogEndDistance *data.Visual_WaterFogEndDistance
	FogColor            uint32
	SkyColor            uint32
	WaterFogColor       uint32

	IncreasedFireBurnout   bool
	SnowGolemMelts         bool
	CanPillagerPatrolSpawn bool

	GrassColor         uint32
	GrassColorModifier string
	WaterColor         uint32
	FoliageColor       uint32
	DryFoliageColor    uint32
}

type Mob struct {
	Type     string
	MaxCount uint8
	MinCount uint8

	Weight uint8
}

type Spawners struct {
	Ambient                  []Mob
	Axolotls                 []Mob
	Creature                 []Mob
	Misc                     []Mob
	Monster                  []Mob
	UndergroundWaterCreature []Mob
	WaterAmbient             []Mob
	WaterCreature            []Mob
}

type Costs struct {
	Charge       float32
	EnergyBudget float32
}

func NewBiomeFileWriter(w io.Writer, mcVer, mcProtocol string) BiomeFileWriter {

	w.Write(([]byte(fmt.Sprintf(`// Complete auto-generated file do not edit mannualy!
// Generated for Minecraft %s (Protocol %s)
package worldgen

import "github.com/robogg133/KernelCraft/data"

type Biome struct {
		Carvers  []string
		Downfall float32

		Features [][]string

		HasPreciptation bool

		Spawners Spawners

		Temperature float32

		SpawnCosts map[data.EntityID]Costs

		AmbientSounds   *data.Audio_AmbientSounds
		BackgorundMusic *data.Audio_BackgroundMusic
		MusicVolume		float32

		AmbientParticles    *data.Visual_AmbientParticles
		WaterFogEndDistance *data.Visual_WaterFogEndDistance
		FogColor            uint32
		SkyColor            uint32
		WaterFogColor       uint32

		IncreasedFireBurnout   bool
		SnowGolemMelts         bool
		CanPillagerPatrolSpawn bool

		GrassColor      uint32
		GrassColorModifier string
		WaterColor      uint32
		FoliageColor    uint32
		DryFoliageColor uint32
}

type Mob struct {
		Type     string
		MaxCount uint8
		MinCount uint8

		Weight uint8
}

type Spawners struct {
		Ambient                  []Mob
		Axolotls                 []Mob
		Creature                 []Mob
		Misc                     []Mob
		Monster                  []Mob
		UndergroundWaterCreature []Mob
		WaterAmbient             []Mob
		WaterCreature            []Mob
}

type Costs struct {
		Charge       float32
		EnergyBudget float32
}
`, mcVer, mcProtocol))))

	return BiomeFileWriter{
		w: w,
	}
}

func (w *BiomeFileWriter) WriteObject(b Biome, name string) {
	w.writeStringf("var Biome_%s = Biome{\n", name)
	w.knownBiomes = append(w.knownBiomes, fmt.Sprintf("minecraft:%s", name))

	// --- Climate ---
	w.writeStringf("\tTemperature: %f,\n", b.Temperature)
	w.writeStringf("\tDownfall: %f,\n", b.Downfall)
	w.writeStringf("\tHasPreciptation: %t,\n", b.HasPreciptation)

	// --- Colors ---
	w.writeStringf("\tFogColor: %d,\n", b.FogColor)
	w.writeStringf("\tSkyColor: %d,\n", b.SkyColor)
	w.writeStringf("\tWaterFogColor: %d,\n", b.WaterFogColor)
	w.writeStringf("\tWaterColor: %d,\n", b.WaterColor)
	w.writeStringf("\tGrassColor: %d,\n", b.GrassColor)
	w.writeStringf("\tFoliageColor: %d,\n", b.FoliageColor)
	w.writeStringf("\tDryFoliageColor: %d,\n", b.DryFoliageColor)
	w.writeStringf("\tGrassColorModifier: \"%s\",\n", b.GrassColorModifier)

	// --- Flags ---
	w.writeStringf("\tIncreasedFireBurnout: %t,\n", b.IncreasedFireBurnout)
	w.writeStringf("\tSnowGolemMelts: %t,\n", b.SnowGolemMelts)
	w.writeStringf("\tCanPillagerPatrolSpawn: %t,\n", b.CanPillagerPatrolSpawn)

	// --- Carvers ---
	w.writeStringf("\tCarvers: ")
	w.writeStringSlice(b.Carvers)
	w.writeStringf(",\n")

	// --- Features ---
	w.writeStringf("\tFeatures: ")
	w.write2DStringSlice(b.Features)
	w.writeStringf(",\n")

	// --- Spawners ---
	w.writeStringf("\tSpawners: Spawners{\n")
	w.writeMobList("Ambient", b.Spawners.Ambient)
	w.writeMobList("Axolotls", b.Spawners.Axolotls)
	w.writeMobList("Creature", b.Spawners.Creature)
	w.writeMobList("Misc", b.Spawners.Misc)
	w.writeMobList("Monster", b.Spawners.Monster)
	w.writeMobList("UndergroundWaterCreature", b.Spawners.UndergroundWaterCreature)
	w.writeMobList("WaterAmbient", b.Spawners.WaterAmbient)
	w.writeMobList("WaterCreature", b.Spawners.WaterCreature)
	w.writeStringf("\t},\n")

	// --- Spawn costs ---
	if len(b.SpawnCosts) > 0 {
		w.writeStringf("\tSpawnCosts: map[data.EntityID]Costs{\n")
		for id, cost := range b.SpawnCosts {
			w.writeStringf(
				"\t\t%d: {Charge: %f, EnergyBudget: %f},\n",
				id, cost.Charge, cost.EnergyBudget,
			)
		}
		w.writeStringf("\t},\n")
	}

	// --- Audio ---
	if b.AmbientSounds != nil {
		w.writeStringf("\tAmbientSounds: &%#v,\n", *b.AmbientSounds)
	}
	if b.BackgorundMusic != nil {
		w.writeStringf("\tBackgorundMusic: &%#v,\n", *b.BackgorundMusic)
	}
	w.writeStringf("\tMusicVolume: %f,\n", b.MusicVolume)

	// --- Visual ---
	if b.AmbientParticles != nil {
		w.writeStringf("\tAmbientParticles: &%#v,\n", *b.AmbientParticles)
	}
	if b.WaterFogEndDistance != nil {
		w.writeStringf("\tWaterFogEndDistance: &%#v,\n", *b.WaterFogEndDistance)
	}

	w.writeStringf("}\n\n")
}

func (w *BiomeFileWriter) Finish() {

	w.writeStringf("var KnownBiomes = ")
	w.writeStringSlice(w.knownBiomes)
	w.writeStringf("\n")
}

func (w *BiomeFileWriter) write2DStringSlice(v [][]string) {
	w.writeStringf("[][]string{\n")

	for _, row := range v {
		w.writeStringf("\t\t{")
		for i, s := range row {
			if i > 0 {
				w.writeStringf(", ")
			}
			w.writeStringf("\"%s\"", s)
		}
		w.writeStringf("},\n")
	}

	w.writeStringf("\t}")
}

func (w *BiomeFileWriter) writeStringf(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	_, err := w.w.Write([]byte(s))
	if err != nil {
		panic(err)
	}
}

func (w *BiomeFileWriter) writeMobList(name string, mobs []Mob) {
	w.writeStringf("\t\t%s: []Mob{\n", name)
	for _, m := range mobs {
		w.writeStringf(
			"\t\t\t{Type: \"%s\", Weight: %d, MinCount: %d, MaxCount: %d},\n",
			m.Type, m.Weight, m.MinCount, m.MaxCount,
		)
	}
	w.writeStringf("\t\t},\n")
}

func (w *BiomeFileWriter) writeStringSlice(a []string) {
	w.writeStringf("[]string{")
	for _, s := range a {
		w.writeStringf("\"%s\",", s)
	}
	w.writeStringf("}")
}

//

func ConvertRawSpawners(r RawSpawners) Spawners {
	return Spawners{
		Ambient:                  convertRawMobSlice(r.Ambient),
		Axolotls:                 convertRawMobSlice(r.Axolotls),
		Creature:                 convertRawMobSlice(r.Creature),
		Misc:                     convertRawMobSlice(r.Misc),
		Monster:                  convertRawMobSlice(r.Monster),
		UndergroundWaterCreature: convertRawMobSlice(r.UndergroundWaterCreature),
		WaterAmbient:             convertRawMobSlice(r.WaterAmbient),
		WaterCreature:            convertRawMobSlice(r.WaterCreature),
	}
}

func convertRawMobSlice(v []RawMob) []Mob {
	if len(v) == 0 {
		return nil
	}

	out := make([]Mob, 0, len(v))
	for _, r := range v {
		out = append(out, convertRawMob(r))
	}
	return out
}

func convertRawMob(r RawMob) Mob {
	return Mob{
		Type:     r.Type,
		Weight:   uint8(r.Weight),
		MinCount: uint8(r.MinCount),
		MaxCount: uint8(r.MaxCount),
	}
}
