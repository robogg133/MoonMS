package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/robogg133/MoonMS/data"
	parser "github.com/robogg133/MoonMS/tools/autogen/minecraft-core/parse"
	"github.com/robogg133/MoonMS/tools/autogen/minecraft-core/parse/worldgen"
)

func doBiome(startingDir, biomeFolder, releaseName string) {
	dir, err := os.ReadDir(biomeFolder)
	if err != nil {
		panic(err)
	}

	os.Chdir(startingDir)
	if err := os.MkdirAll("internal/gen/core/damage/", 0777); err != nil {
		panic(err)
	}
	f, err := os.OpenFile("internal/gen/core/worldgen/biomes.go", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := worldgen.NewBiomeFileWriter(f, releaseName)

	for _, d := range dir {
		if d.IsDir() {
			fmt.Println("what, i found a directory")
			continue
		}

		fmt.Printf("====== FILE: %s ======\n", d.Name())

		b, err := os.ReadFile(filepath.Join(biomeFolder, d.Name()))
		if err != nil {
			panic(err)
		}

		var target worldgen.Biome
		var a worldgen.RawBiome
		a.Attributes = make(map[string]json.RawMessage)
		a.Effects = make(map[string]json.RawMessage)
		a.SpawnCosts = make(map[string]worldgen.RawCost)

		if err := json.Unmarshal(b, &a); err != nil {
			panic(err)
		}

		target.Downfall = a.Downfall
		target.HasPreciptation = a.HasPrecip
		target.Temperature = a.Temperature
		target.Spawners = worldgen.ConvertRawSpawners(a.Spawners)
		target.Features = a.Features
		target.MusicVolume = 100

		if err := json.Unmarshal(a.Carvers, &target.Carvers); err != nil {
			var s string

			if err := json.Unmarshal(a.Carvers, &s); err != nil {
				panic(err)
			}

			target.Carvers = []string{s}
		}

		for i, v := range a.Attributes {
			fmt.Printf("[DEBUG ATTRIBUTES]: %s : %s\n", i, string(v))
			fn := parser.AttributeDecoders[i]
			a, err := fn(v)
			if err != nil {
				panic(err)
			}
			switch i {
			case "minecraft:audio/background_music":
				target.BackgorundMusic = a.(*data.Audio_BackgroundMusic)
			case "minecraft:audio/ambient_sounds":
				target.AmbientSounds = a.(*data.Audio_AmbientSounds)
			case "minecraft:audio/music_volume":
				target.MusicVolume = float32(*a.(*data.Audio_MusicVolume))
			case "minecraft:gameplay/snow_golem_melts":
				target.SnowGolemMelts = bool(*a.(*data.Gameplay_SnowGolemMelts))
			case "minecraft:gameplay/increased_fire_burnout":
				target.IncreasedFireBurnout = bool(*a.(*data.Gameplay_IncreasedFireBurnout))

			case "minecraft:visual/ambient_particles":
				target.AmbientParticles = a.(*data.Visual_AmbientParticles)
			case "minecraft:visual/water_fog_color":
				target.WaterFogColor = uint32(*a.(*data.Visual_WaterFogColor))
			case "minecraft:visual/water_fog_end_distance":
				target.WaterFogEndDistance = a.(*data.Visual_WaterFogEndDistance)
			case "minecraft:visual/fog_color":
				target.FogColor = uint32(*a.(*data.Visual_FogColor))
			case "minecraft:visual/sky_color":
				target.SkyColor = uint32(*a.(*data.Visual_SkyColor))
			}

			fmt.Println(a.ID())
		}

		for i, v := range a.Effects {
			fmt.Printf("[DEBUG EFFECTS]: %s : %s\n", i, string(v))
			fn := parser.AttributeDecoders[i]
			a, err := fn(v)
			if err != nil {
				panic(err)
			}

			switch i {

			case "water_color":
				target.WaterColor = uint32(*a.(*data.Effect_WaterColor))

			case "grass_color":
				target.GrassColor = uint32(*a.(*data.Effect_GrassColor))

			case "foliage_color":
				target.FoliageColor = uint32(*a.(*data.Effect_FoliageColor))

			case "dry_foliage_color":
				target.DryFoliageColor = uint32(*a.(*data.Effect_DryFoliageColor))

			case "grass_color_modifier":
				target.GrassColorModifier = string(*a.(*data.Effect_GrassColorModifier))
			}

			fmt.Println(a.ID())
		}

		w.WriteObject(target, strings.TrimSuffix(d.Name(), ".json"))
	}

	w.Finish()
}
