package level

import (
	"os"
	"testing"

	"github.com/Tnze/go-mc/nbt"
)

type Level struct {
	Data Data `nbt:"Data"`
}

type Data struct {
	DifficultySettings struct {
		Difficulty string `nbt:"difficulty"`
		Hardcore   bool   `nbt:"hardcore"`
		Locked     bool   `nbt:"locked"`
	} `nbt:"difficulty_settings"`

	Time     int64  `nbt:"Time"`
	GameType string `nbt:"GameType"`

	ServerBrands []string `nbt:"ServerBrands"`

	Version_ int32 `nbt:"version"`

	LastPlayed int64 `nbt:"LastPlayed"`

	Spawn struct {
		Pos []int32 `nbt:"pos"`

		Pitch     int32  `nbt:"pitch"`
		Dimension string `nbt:"dimension"`
		Yaw       int32  `nbt:"yaw"`
	} `nbt:"spawn"`

	Version struct {
		Snapshot bool   `nbt:"Snapshot"`
		Series   string `nbt:"Series"`
		Id       int32  `nbt:"Id"`
		Name     int32  `nbt:"Name"`
	} `nbt:"Version"`

	LevelName string `nbt:"LevelName"`

	Initialized   bool  `nbt:"initialized"`
	WasModded     bool  `nbt:"WasModded"`
	DataVersion   int32 `nbt:"DataVersion"`
	AllowCommands bool  `nbt:"allowCommands"`

	DataPacks struct {
		Enabled  []string `nbt:"Enabled"`
		Disabled []string `nbt:"Disabled"`
	} `nbt:"DataPacks"`
}

func Test(t *testing.T) {

	a := Level{}

	b, err := nbt.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}

	os.WriteFile("test.nbt", b, 0777)
}
