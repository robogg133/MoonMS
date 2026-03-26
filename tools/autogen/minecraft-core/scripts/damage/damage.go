//go:generate go build -o exec .
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	damagetype "github.com/robogg133/MoonMS/tools/autogen/minecraft-core/scripts/damage/damage_type"
)

func main() {

	startingDir := os.Args[1]
	damageFolder := os.Args[2]
	releaseName := os.Args[3]

	os.MkdirAll(startingDir, 0777)

	dir, err := os.ReadDir(damageFolder)
	if err != nil {
		panic(err)
	}

	os.Chdir(startingDir)
	f, err := os.OpenFile(
		"types.go",
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0777,
	)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := damagetype.NewDamageTypeFileWriter(f, releaseName)

	for _, d := range dir {
		if d.IsDir() {
			continue
		}

		fmt.Printf("====== FILE: %s ======\n", d.Name())

		rawBytes, err := os.ReadFile(filepath.Join(damageFolder, d.Name()))
		if err != nil {
			panic(err)
		}

		var raw damagetype.RawDamageType
		if err := json.Unmarshal(rawBytes, &raw); err != nil {
			panic(err)
		}

		target := damagetype.DamageType{
			Exhaustion:       raw.Exhaustion,
			MessageID:        raw.MessageID,
			Scaling:          raw.Scaling,
			DeathMessageType: raw.DeathMessageType,
		}

		w.WriteObject(strings.TrimSuffix(d.Name(), ".json"), target)
	}

	w.Finish()
}
