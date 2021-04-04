package config

const (
	// Screen Settings
	COLS = 200
	ROWS = 200
	CELL_SIZE = 3

	// Loop Settings
	INTERVAL_NS = 0 * 1000000		// Amount of nanoseconds to make each game loop last for (ns = 1000 us = 1000 000 ms)
									// If the loop takes longer, than this will be allowed

	// GAME SETTINGS
	FLASHY_EEL = false 				// Mystery setting! (toggle with 'E')

	// Cell Array Keys (used for datagrid.Cells[r][c][KEY])
	KEY_AMOUNT = 0
	KEY_I_AMOUNT = 3	// intermediate, used to update amount when moving	
	KEY_SMELL = 1
	KEY_I_SMELL = 2		// intermediate, used to update amount when calculating smell

	// LUT Keys (used for datagrid.NeighbourLUT[r][c][KEY])
	LUTKEY_ROW = 0
	LUTKEY_COL = 1
	LUTKEY_EXISTS = 2	

	// INPUT
	MOUSE_LEFT_CLICK = 1
	MOUSE_MIDDLE_CLICK = 2
	MOUSE_RIGHT_CLICK = 3
	BUTTON_DOWN = 0
	BUTTON_UP = 1
)