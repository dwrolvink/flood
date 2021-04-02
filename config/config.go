package config

// Import built-in packages
import (
	//"math"
	"time"       // used for pausing, measuring duration, etc
)

const (

	// Cell Array Keys (used for datagrid.Cells[r][c][KEY])
	KEY_AMOUNT = 0
	KEY_I_AMOUNT = 3	// intermediate, used to update amount when moving	
	KEY_SMELL = 1
	KEY_I_SMELL = 2		// intermediate, used to update amount when calculating smell

	// LUT Keys (used for datagrid.NeighbourLUT[r][c][KEY])
	LUTKEY_ROW = 0
	LUTKEY_COL = 1
	LUTKEY_EXISTS = 2

	// Screen Settings
	COLS = 100
	ROWS = 100
	CELL_SIZE = 4

	// Loop Settings
	INTERVAL = 0		// Amount of ms to sleep after each game loop

	// Draw Settings
	DRAW_BETWEEN_BATTLE = false
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