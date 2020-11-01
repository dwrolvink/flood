package world

// Import built-in packages
import (
	//"fmt"        // used for outputting to the terminal
	"math"
)

// Import external packages
import (
	"github.com/veandco/go-sdl2/sdl"
)

// Import internal packages
import (
	"flood_go/cell"
)

func CreateRectGrid(rows int32, cols int32, cell_size int32) [][]cell.Cell {
	
	// Make grid
	cell_grid := make([][]cell.Cell, rows)
	for row := range cell_grid {
		cell_grid[row] = make([]cell.Cell, cols)
	}	

	var x int32 = 0
	var y int32 = 0
	var rect sdl.Rect
	var color sdl.Color

	// initialize all cells
	for row := 0; row < int(rows); row++ {
		for col := 0; col < int(cols); col++ {
			// define rect
			rect = sdl.Rect{x, y, cell_size, cell_size}
			x += cell_size

			// define color
			color = sdl.Color{255,255,255,255}

			cell_grid[row][col] = cell.Cell{rect, col, row, &color, 0, 0, 0, 0, 0, false, []*cell.Cell{}}
		}
		// Go back to the left, and scroll down one row
		x = 0
		y += cell_size
	}

	// load neighbours
	for row := 0; row < int(rows); row++ {
		for col := 0; col < int(cols); col++ {
			// Load neighbours
			cell_grid[row][col].SetNeighbours(&cell_grid)
		}
	}

	return cell_grid
}

func LoopGrid(grid *[][]cell.Cell, def func(*cell.Cell)) {
	for row := range *grid {
		for col := range (*grid)[row] {
			var cell = &(*grid)[row][col]
			def(cell)
		}
	}
}

func GetPos(mouseX int32, mouseY int32, cellSize int32) [2]int32 {
	var row = int32(math.Floor(float64(mouseY)/float64(cellSize)))
	var col = int32(math.Floor(float64(mouseX)/float64(cellSize)))

	return [2]int32{row, col}
}


func ExecuteExpedition(exp cell.Expedition) {
	// update registers
	exp.Source.Register1 = exp.Register1
	exp.Source.Register2 = exp.Register2
	exp.Source.Register3 = exp.Register3

	// if amount == 0, do nothing
	if exp.Amount == 0 {
		return
	}

	// check if cell has already been beaten
	if exp.Source.CanAct == false {
		return
	}

	// update source
	exp.Source.Amount -= exp.Amount


	if exp.Target.UserId == 0 || exp.Target.UserId == exp.Source.UserId {
		// move amount to empty/friendly cell
		exp.Target.UserId = exp.Source.UserId
		exp.Target.BaseColor = exp.Source.BaseColor
		exp.Target.Amount += exp.Amount

	} else {
		// Attack enemy cell
		var remainder = int(exp.Amount) - int(exp.Target.Amount)

		// Target cell is attacked, so it can't act this round
		exp.Target.CanAct = false

		if remainder == 0 {
			// tie, make empty cell
			exp.Target.UserId = 0
			exp.Target.BaseColor = &sdl.Color{255,255,255,255}
			exp.Target.Amount = 0
			exp.Target.Register1 = 0
			exp.Target.Register2 = 0
			exp.Target.Register3 = 0
		} else if remainder > 0 {
			// won, take over cell
			exp.Target.UserId = exp.Source.UserId
			exp.Target.BaseColor = exp.Source.BaseColor
			exp.Target.Amount = uint(remainder)
			exp.Target.Register1 = 0
			exp.Target.Register2 = 0
			exp.Target.Register3 = 0			
		} else if remainder < 0 {
			// lost, just set negative of remainder
			exp.Target.Amount = uint(-1 * remainder)
		}
	}

}
