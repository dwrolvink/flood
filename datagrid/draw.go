package datagrid

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


func NormalizeSmellColor(smell int) int {
	
	if smell == 0 { return 0}
	if smell < 20 { return 20}
	if smell < 100 { return smell}
	if smell < 10000 { return int(smell / 100) + 99}
	return 255
}

func NormalizeAmountColor(amount int, max_out bool) int {
	s := amount
	if max_out {
		if s > 0 {
			return 255
		}
	}
	if s > 255 {
		return 255
	}	
	return s
}

/* un-implemented atm, in favour of pixelbased drawing */
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

func (this *DataGrid) Draw_PixelBased(done chan bool, row_start, row_end int, numbers_text *[256]*text.TextObject, print_val int, print_text bool){
	// Shorthand
	cell_size := int32(cfg.CELL_SIZE)

	// references
	var Cells = &(this.Cells)

	// Convert SDL color to pixel value
	var color = [4]byte{this.BaseColor.R, this.BaseColor.G, this.BaseColor.B, this.BaseColor.A}
	
	// Draw & Rearm Loop
	var x = int32(0)
	var y = int32(0) + (int32(row_start)*cell_size)
	var draw_value = 0
	for row := row_start; row < row_end; row++ {
		for col := 0; col < cfg.COLS; col++ {
			// get cell 
			//current_cell := (*Cells)[row][col]
			
			// get draw value
			draw_value = (*Cells)[row][col][uint8(print_val)]

			//index := row * col + col 
			//draw_value := int((*this.Amount)[3])
			
			if draw_value > 0 {
				/*
				if print_val == 0 {
					draw_value = NormalizeAmountColor(draw_value, false)
				}
				*/	
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
	
	done <- true
}

func (this *DataGrid) Draw_PixelBasedDotted(done chan bool, row_start, row_end int, numbers_text *[256]*text.TextObject, print_val int, print_text bool){
	// Shorthand
	cell_size := int32(cfg.CELL_SIZE)

	// references
	var Cells = &(this.Cells)

	// Vars
	print_value := uint8(print_val)
	row_skip := cfg.COLS * cfg.CELL_SIZE * 4 * (cfg.CELL_SIZE)
	index  := int32(3) + int32(row_skip * row_start)

	if cell_size > 1 {
		for row := row_start; row < row_end; row++ {
			for col := 0; col < cfg.COLS; col++ {
				(*this.Pixels)[index] = byte((*Cells)[row][col][print_value])
				index += cell_size * 4
			}
			index += cfg.COLS * cfg.CELL_SIZE * 4 * (cfg.CELL_SIZE - 1)
		}
	} else {
		for row := row_start; row < row_end; row++ {
			for col := 0; col < cfg.COLS; col++ {
				(*this.Pixels)[index] = byte((*Cells)[row][col][print_value])
				index += cell_size * 4
			}
		}
	}

	done <- true
}
