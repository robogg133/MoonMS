package seed_test

import (
	"fmt"
	"testing"

	"github.com/robogg133/KernelCraft/pkg/minecraft/world/seed"
)

func TestGenerate(t *testing.T) {

	fmt.Println(seed.GenerateSeed())
	fmt.Println(seed.GenerateByString("robogg133"))
}
