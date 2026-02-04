package seed_test

import (
	"MoonMS/pkg/minecraft/world/seed"
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {

	fmt.Println(seed.GenerateSeed())
	fmt.Println(seed.GenerateByString("robogg133"))
}
