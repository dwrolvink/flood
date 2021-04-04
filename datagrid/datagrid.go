package datagrid

import (
	"fmt"        // used for outputting to the terminal
	//"time"       // used for pausing, measuring duration, etc
	//"math/rand"  // random number generator
	//"math"
	//"sync"
)

// Import internal packages
import (
	cfg "flood_go/config"
)

// Import external packages
import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	COLS = cfg.COLS
	ROWS = cfg.ROWS
)

type DataGrid struct {
	Cells [ROWS][COLS][5] int		// amount, smell, _smell, _amount, unused
	//Smell [ROWS][COLS][2]float64    // smell, _smell [todo]

	UserId int			
	BaseColor *sdl.Color
	Enemy *DataGrid

	Pixels *[]byte

	NeighbourLUT [ROWS][COLS][4][3]int // row, col, exists

	// temp
	Amount *[]byte
}

// temp
func (this *DataGrid) SetAmount(index int, value byte) {
	this.Amount[index] = value
}
func (this *DataGrid) GetAmount(index int, value byte) byte {
	return this.Amount[index]
}
func (this *DataGrid) GetByteIndex(row, col int) {
	var cell_length = cfg.CELL_SIZE * 4 	// 4 bytes per pixel, N pixels per cell
	var start = 3 							// first pixel starts at position 3 (4th byte)
	var index = start + cell_length * col   // get position in row
	index += cell_length * col * row		// move N rows down to get position in grid

	return index
}

func (this *DataGrid) Init() {
	this.CalculateNeighbourLUT()
	fmt.Println("")
}

/* Sets entire Cell array to zeros */
func (this *DataGrid) Clear() {
	for row := range this.Cells {
		for col := range this.Cells[row] {
			this.Cells[row][col][0] = 0
			this.Cells[row][col][1] = 0
			this.Cells[row][col][2] = 0
			this.Cells[row][col][3] = 0
			this.Cells[row][col][4] = 0
		}
	}
}

/* Set cfg.KEY_AMOUNT and cfg.KEY_I_AMOUNT simultaneously */
func (this *DataGrid) SetCell(row, col int, amount int) {  
	this.Cells[row][col][cfg.KEY_AMOUNT] = amount 
	this.Cells[row][col][cfg.KEY_I_AMOUNT] = amount 

	// temp
	index := this.GetByteIndex(row, col)
	this.SetAmount(index, uint8(amount))
}

/* Alias for SetCell(r, c, 0) */
func (this *DataGrid) Kill(row, col int) {
	this.SetCell(row, col, 0)
}

/* 	Makes a lookup table that allows us to lookup if cell r,c has 
   	top/bottom/left/right neighbour  
*/
func (this *DataGrid) CalculateNeighbourLUT() {
	for row := 0; row < cfg.ROWS; row++ {
		for col := 0; col < cfg.COLS; col++ {
				
			if (row > 0) { 			// :top
				this.NeighbourLUT[row][col][0][0] = row - 1 
				this.NeighbourLUT[row][col][0][1] = col 
				this.NeighbourLUT[row][col][0][2] = 1 // exists
			}
			if (row < ROWS-1) { 	// :bottom
				this.NeighbourLUT[row][col][1][0] = row + 1
				this.NeighbourLUT[row][col][1][1] = col
				this.NeighbourLUT[row][col][1][2] = 1 // exists				
			}
			if (col > 0) {  		// :left
				this.NeighbourLUT[row][col][2][0] = row
				this.NeighbourLUT[row][col][2][1] = col - 1
				this.NeighbourLUT[row][col][2][2] = 1 // exists
			}
			if (col < COLS-1) {  	// :right
				this.NeighbourLUT[row][col][3][0] = row
				this.NeighbourLUT[row][col][3][1] = col + 1
				this.NeighbourLUT[row][col][3][2] = 1 // exists			
			}	
		}
	}
}

/* Not used atm */
func (this *DataGrid) ClearSmell() {  
	// reset smell of every cell to 0
	
	for row := 0; row < cfg.ROWS; row++ {
		for col := 0; col < cfg.COLS; col++ {
			this.Cells[row][col][1] = 0
		}
	}
}


func (this *DataGrid) GetAvgSmell(row, col, depth int) int {
	// References 
	var NeighbourLUT = &(this.NeighbourLUT)
	var Cells = &(this.Cells)

	// Ephemeral
	var nb_row int
	var nb_col int
	var nb_smell int
	var nb_amount int

	// Used to calc the avg
	number_of_nbs := 0
	total_sum := 0
	amount := (*Cells)[row][col][cfg.KEY_AMOUNT]

	// Collect data from all (existing) neighbours
	for i := 0; i < 4; i++ {
		if (*NeighbourLUT)[row][col][i][cfg.LUTKEY_EXISTS] == 1 {
			// Keep a tally of total existing neighbours to get a good avg
			number_of_nbs ++

			nb_row  = (*NeighbourLUT)[row][col][i][cfg.LUTKEY_ROW]
			nb_col  = (*NeighbourLUT)[row][col][i][cfg.LUTKEY_COL]

			// If depth > 0: just ask the cell to give its average smell & add to total
			if depth > 0 {
				total_sum += this.GetAvgSmell(nb_row, nb_col, depth-1)

			// Else: calc new smell for the cell & add to total
			} else {
				nb_smell = (*Cells)[nb_row][nb_col][cfg.KEY_SMELL]
				nb_amount = (*Cells)[nb_row][nb_col][cfg.KEY_AMOUNT]
				total_sum += nb_smell + nb_amount
			}
		}
	}

	// calc intermediate smell
	return (amount + (total_sum / number_of_nbs))
}


func (this *DataGrid) UpdateSmell(f float64) {  
	// References 
	var Cells = &(this.Cells)

	// calc intermediate smell
	/*
	for row := 0; row < cfg.ROWS; row++ {
		for col := 0; col < cfg.COLS; col++ {

			// Get Smell Average for current cell
			_smell = this.GetAvgSmell(row, col, 0)
			
			// Adjust smell to allow for dissipation
			if _smell > 200 {
				_smell = int(float64(_smell) * 0.999)
			}
			if _smell > 0 && f > float64(_smell)/200 {
				_smell -= 1
			}			
		
			// update intermediate smell
			(*Cells)[row][col][cfg.KEY_I_SMELL] = _smell
		}
	}
	*/

	// temp
	done_top := make(chan bool)
	done_bottom := make(chan bool)

	go this.UpdateIntermediateSmell(done_top, 0, cfg.ROWS/2, f)
	go this.UpdateIntermediateSmell(done_bottom, cfg.ROWS/2, cfg.ROWS, f)

	<- done_top
	<- done_bottom

	// update smell
	for row := 0; row < cfg.ROWS; row++ {
		for col := 0; col < cfg.COLS; col++ {
			(*Cells)[row][col][cfg.KEY_SMELL] = (*Cells)[row][col][2]
		}
	}
}

func (this *DataGrid) UpdateIntermediateSmell(done chan bool, row_start, row_end int, f float64) {  
	var _smell int

	// References 
	var Cells = &(this.Cells)

	// calc intermediate smell
	for row := row_start; row < row_end; row++ {
		for col := 0; col < cfg.COLS; col++ {

			// Get Smell Average for current cell
			_smell = this.GetAvgSmell(row, col, 0)
			
			// Adjust smell to allow for dissipation
			if _smell > 200 {
				_smell = int(float64(_smell) * 0.999)
			}
			if _smell > 0 && f > float64(_smell)/200 {
				_smell -= 1
			}			
		
			// update intermediate smell
			(*Cells)[row][col][cfg.KEY_I_SMELL] = _smell
		}
	}

	// signal that we're done
	done <- true
}


func (this *DataGrid) KillAll() {
	for row := 0; row < cfg.ROWS; row++ {
		for col := 0; col < cfg.COLS; col++ {
			this.SetCell(row, col, 0)
		}
	}
}