package config

// Import built-in packages
import (
	//"math"
	"time"       // used for pausing, measuring duration, etc
)


type Config struct {
	ScreenTitle string		// Set by caller
	ScreenWidth int32		// Set in Init()
	ScreenHeight int32		// Set in Init()
	CellSize int32			// From CONST
	Cols int32				// From CONST
	Rows int32				// From CONST
	Interval time.Duration	// From CONST / Adjusted in Init()
	ShowDebugText bool		// Set by caller
	DrawBetweenBattle bool	// From CONST
}
func (this *Config) Init() {  
	this.Rows = ROWS
	this.Cols = COLS
	this.CellSize = CELL_SIZE
	this.ScreenWidth = COLS * CELL_SIZE
	this.ScreenHeight = ROWS * CELL_SIZE
	this.Interval = INTERVAL
	this.DrawBetweenBattle = DRAW_BETWEEN_BATTLE

	if this.DrawBetweenBattle {
		this.Interval /= 2
	}
}

// If you need calculated values you can use this function
// Currently only for calculating ScreenWidth/-Height
// Otherwise, you can use the CONSTs directly
func GetConfig() Config {
	var config = Config{
		ScreenTitle: "SDL Test Application",
		ShowDebugText: false,
	}
	config.Init()
	return config
}