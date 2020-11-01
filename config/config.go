package config

// Import built-in packages
import (
	"math"
	"time"       // used for pausing, measuring duration, etc
)

// You can set config here that can be passed around in a struct
// To add a value, add an entry in the struct below, and then
// set it in GetConfig()

// To use the config some place, first import this subpackage:
//		import "flood_go/config"
//
// Then, you can get the struct by using:
//		var cfg = config.GetConfig()
//
// And finally, use the data:
// 		fmt.Println(cfg.ScreenTitle)

type Config struct {
	ScreenTitle string
	ScreenWidth int32
	ScreenHeight int32
	CellSize int32
	Cols int32
	Rows int32
	SleepMiliSeconds time.Duration
	ShowDebugText bool
}
func (this *Config) Init() {  
	this.Rows = int32(math.Floor(float64(this.ScreenHeight)/float64(this.CellSize)))
	this.Cols = int32(math.Floor(float64(this.ScreenWidth)/float64(this.CellSize)))
}

func GetConfig() Config {
	var config = Config{
		ScreenTitle: "SDL Test Application",
		ScreenWidth: 800,
		ScreenHeight: 600,
		CellSize: 6,
		SleepMiliSeconds: 0,
		ShowDebugText: false,
	}

	config.Init()

	return config

}