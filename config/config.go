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
	IntervalNs time.Duration	// From CONST 
	ShowDebugText bool		// Set by caller
	DrawBetweenBattle bool	// From CONST
	FlashyEel bool			// From CONST / struct / toggle with 'e'
}
func (this *Config) Init() {  
	this.Rows = ROWS
	this.Cols = COLS
	this.CellSize = CELL_SIZE
	this.ScreenWidth = COLS * CELL_SIZE
	this.ScreenHeight = ROWS * CELL_SIZE
	this.IntervalNs = INTERVAL_NS
	this.FlashyEel = FLASHY_EEL
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