package datagrid

import (
	"fmt"        // used for outputting to the terminal
	//"time"       // used for pausing, measuring duration, etc
	"math/rand"  // random number generator
	"math"
	//"sync"
)

// Import internal packages
import (
	"flood_go/graphicsx"	
	"flood_go/text"
	"flood_go/misc"
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

	// temp
	Q [5]int 
	Pixels *[]byte

	NeighbourLUT [ROWS][COLS][4][3]int // row, col, exists
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


func (this *DataGrid) SetCell(row, col int, amount int) {  
	this.Cells[row][col][cfg.KEY_AMOUNT] = amount 
	this.Cells[row][col][cfg.KEY_I_AMOUNT] = amount 
}

func (this *DataGrid) Kill(row, col int) {
	this.SetCell(row, col, 0)
	//this.Cells[row][col][1] = 0
}

/* 
	Makes a lookup table that allows us to lookup if cell r,c has 
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



func (this *DataGrid) ClearSmell() {  
	// reset smell of every cell to 0
	
	for row := range this.Cells {
		for col := range this.Cells[row] {
			this.Cells[row][col][1] = 0
		}
	}
}


func (this *DataGrid) UpdateSmell() {  

	//var amount int

	var nb_row int
	var nb_col int
	var nb_smell int
	var nb_amount int

	var highest_amount int
	var highest_smell int
	var highest_val int	


	// calc intermediate smell
	for row := range this.Cells {
		for col := range this.Cells[row] {

			// reset vars
			highest_amount = 0
			highest_smell = 0
			//amount = this.Cells[row][col][0]
			highest_val = 0

			// calc sum of nbs
			for i := 0; i < 4; i++ {
				if this.NeighbourLUT[row][col][i][2] == 1 {

					nb_row  = this.NeighbourLUT[row][col][i][0]
					nb_col  = this.NeighbourLUT[row][col][i][1]

					nb_smell = this.Cells[nb_row][nb_col][1]
					nb_amount = this.Cells[nb_row][nb_col][0]

					if nb_smell > highest_smell {
						highest_smell = nb_smell
					}
					if nb_amount > highest_amount {
						highest_amount = nb_amount
					}					
				}
			}

			highest_val = highest_amount
			if highest_smell > highest_val {
				highest_val = highest_smell
			}


			

			// update intermediate smell
			this.Cells[row][col][2] = misc.Normalize(highest_val - 1, math.MaxInt64, 0)
		}
	}

	// update smell
	for row := range this.Cells {
		for col := range this.Cells[row] {	
			this.Cells[row][col][1] = this.Cells[row][col][2]
		}
	}
}

func (this *DataGrid) GetAvgSmell(row, col, depth int) int {
	// Ephemeral
	var nb_row int
	var nb_col int
	var nb_smell int
	var nb_amount int

	// Used to calc the avg
	number_of_nbs := 0
	total_sum := 0
	amount := this.Cells[row][col][cfg.KEY_AMOUNT]

	// Collect data from all (existing) neighbours
	for i := 0; i < 4; i++ {
		if this.NeighbourLUT[row][col][i][cfg.LUTKEY_EXISTS] == 1 {
			// Keep a tally of total existing neighbours to get a good avg
			number_of_nbs ++

			nb_row  = this.NeighbourLUT[row][col][i][cfg.LUTKEY_ROW]
			nb_col  = this.NeighbourLUT[row][col][i][cfg.LUTKEY_COL]

			// If depth > 0: just ask the cell to give its average smell & add to total
			if depth > 0 {
				total_sum += this.GetAvgSmell(nb_row, nb_col, depth-1)

			// Else: calc new smell for the cell & add to total
			} else {
				nb_smell = this.Cells[nb_row][nb_col][cfg.KEY_SMELL]
				nb_amount = this.Cells[nb_row][nb_col][cfg.KEY_AMOUNT]
				total_sum += nb_smell + nb_amount
			}
		}
	}

	// calc intermediate smell
	return (amount + (total_sum / number_of_nbs))
}

func (this *DataGrid) UpdateSmell4() {  


	var _smell int

	f := rand.Float64()

	// calc intermediate smell
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
			this.Cells[row][col][cfg.KEY_I_SMELL] = _smell
		}
	}

	// update smell
	for row := range this.Cells {
		for col := range this.Cells[row] {	
			this.Cells[row][col][cfg.KEY_SMELL] = this.Cells[row][col][2]
		}
	}
}

func NormalizeSmellColor(smell int) int {
	
	if smell == 0 { return 0}
	if smell < 20 { return 20}
	if smell < 100 { return smell}
	if smell < 10000 { return int(smell / 100) + 99}
	return 255
}

func NormalizeAmountColor(amount int) int {
	s := amount
	if s > 255 {
		return 255
	}	
	return s
}


func (this *DataGrid) Draw(graphics *graphicsx.Graphics, numbers_text *[256]*text.TextObject, cell_size int32, print_val int, print_text bool){
	// set alpha
	var alpha int

	// Draw & Rearm Loop
	var x = int32(0)
	var y = int32(0)
	for row := range this.Cells {
		for col := range this.Cells[row] {
			// get cell 
			current_cell := this.Cells[row][col]
			
			// get draw value
			draw_value := current_cell[uint8(print_val)]

			if draw_value > 0 {
				// set Cell color
				// set alpha
				alpha = misc.Max255Int(20, draw_value)

				// set color
				(*graphics).SetSDLDrawColor(this.BaseColor, uint8(alpha))

				// create rect
				rect := sdl.Rect{x, y, cell_size, cell_size}
					
				// draw cell
				(*graphics).Renderer.FillRect(&rect)	
			
				// draw text
				if print_text {
					(*graphics).Renderer.Copy((*numbers_text)[draw_value].Image.Texture, nil, &sdl.Rect{
						x, 
						y, 
						(*numbers_text)[draw_value].Image.Width, (*numbers_text)[draw_value].Image.Height,
					})	
				}	
			}			

			x += cell_size
		}
		x = 0
		y += cell_size
	}		
}

func (this *DataGrid) Draw_PixelBased(numbers_text *[256]*text.TextObject, print_val int, print_text bool){
	// Shorthand
	cell_size := int32(cfg.CELL_SIZE)

	// Erase pixels
	graphicsx.Clear(this.Pixels)


	// Convert SDL color to pixel value
	var color = [4]byte{this.BaseColor.R, this.BaseColor.G, this.BaseColor.B, this.BaseColor.A}
	
	// Draw & Rearm Loop
	var x = int32(0)
	var y = int32(0)
	for row := range this.Cells {
		for col := range this.Cells[row] {
			// get cell 
			current_cell := this.Cells[row][col]
			
			// get draw value
			draw_value := current_cell[uint8(print_val)]
			
			//if draw_value > 0 {
			if true {

				if print_val == 0 {
					draw_value = NormalizeAmountColor(draw_value)
				}	
				if print_val == 1 {
					draw_value = NormalizeSmellColor(draw_value)
				}	

				color[3] = uint8(draw_value)
				
				graphicsx.SetSquare(int(x), int(y), int(cell_size), color, this.Pixels)
			}			

			x += cell_size
		}
		x = 0
		y += cell_size
	}		
}

func (this *DataGrid) KillAll() {
	for row := range this.Cells {
		for col := range this.Cells[row] {
			this.SetCell(row, col, 0)
		}
	}
}