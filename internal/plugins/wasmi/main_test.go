package plugin_wasmi_test

import (
	plugin_wasmi "MoonMS/internal/plugins/wasmi"
	"MoonMS/internal/server"
	"fmt"
	"os"
	"testing"
)

func TestXxx(t *testing.T) {
	t.Run("asdadsd", func(t *testing.T) {
		da, err := server.InitServerData()
		if err != nil {
			t.Error(err)
		}

		fmt.Println(da.Motd)

		runtime := plugin_wasmi.CreateRuntime("test")

		f, err := os.ReadFile("/home/robo/MoonMS/plugin.wasm")
		if err != nil {
			t.Error(err)
		}

		if err := plugin_wasmi.RunFile(f, runtime); err != nil {
			t.Error(err)
		}
	})
}
