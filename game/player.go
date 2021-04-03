package game

// Import built-in packages
import (
	"fmt"
	//"time"       // used for pausing, measuring duration, etc
	"math"
	//"math/rand"
	"bytes"
)

// Import external packages
import (
	"github.com/veandco/go-sdl2/sdl"
)

// Import internal packages
import (
	"flood_go/graphicsx"
	"flood_go/datagrid"	
	"flood_go/text"
	"flood_go/misc"
	cfg "flood_go/config"
)

type Player struct {
	UserId int	
	Name string
	BaseColor *sdl.Color
	DataGrid *datagrid.DataGrid
	Enemy *Player

	// When a start cell is selected (.SetStartCell()), 
	// these settings will be saved here
	StartCell [3]int	

	// Used for drawing the texture
	Pixels []byte
	Texture *sdl.Texture

	// temp
	Amount []byte
}

// =================================================================================================================================================
// INIT
// -------------------------------------------------------------------------------------------------------------------------------------------------

func (this *Player) Init(graphics *graphicsx.Graphics) { 
	// Set empty pixel array
	this.InitPixels()
	this.InitAmount()

	// Create Texture
	this.InitTexture(graphics)

	// Init sub structs
	this.DataGrid = &datagrid.DataGrid{BaseColor: this.BaseColor, Pixels: &this.Pixels, Amount: &this.Amount}	
	this.DataGrid.Init()

	// Lame print, so we don't have to toggle imports when debugging
	if false { fmt.Println("") }
}

func (this *Player) InitPixels() {
	// Make a byte slice, and prefill it with (basecolor + alpha=0)
	number_of_pixels := (cfg.COLS * cfg.CELL_SIZE) * (cfg.ROWS * cfg.CELL_SIZE)

	// Set color
	pixel := []byte{this.BaseColor.R, this.BaseColor.G, this.BaseColor.B, 0}

	// Repeat the color N number of times to form the slice
	this.Pixels = bytes.Repeat(pixel, number_of_pixels)
}

/* 	project: read and write from amount byte slice so that we can draw on the screen directly from this slice
			when the cell size = 1 pixel
	[note: not-implemented]
*/ 
func (this *Player) InitAmount() {
	// Make Pixel array, in which we'll store the amount data (so that we don't have to build the pixel array)
	number_of_cells := cfg.ROWS * cfg.COLS

	// Fill r,g,b with basecolor value (a will be the amount)
	color := []byte{this.BaseColor.R, this.BaseColor.G, this.BaseColor.B, 0}
	Amount := bytes.Repeat(color, number_of_cells)
	this.Amount = Amount
}

func (this *Player) InitTexture(graphics *graphicsx.Graphics) {
	var ScreenWidth = cfg.COLS * cfg.CELL_SIZE
	var ScreenHeight = cfg.ROWS * cfg.CELL_SIZE
	// Don't forget to remove the texture when you're done with it! (defer texture.Destroy())
	texture, err := graphics.Renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(ScreenWidth), int32(ScreenHeight))
	if err != nil {
	  panic(err)
	}
	
	// This allows us to merge two textures together
	texture.SetBlendMode(sdl.BLENDMODE_ADD)

	// Write it to the struct
	this.Texture = texture
}

func (this *Player) SetEnemy(Enemy *Player) {
	this.Enemy = Enemy

	// Set enemy in datagrid as shorthand (temp)
	this.DataGrid.Enemy = Enemy.DataGrid
}

func (this *Player) SetStartCell(row, col, amount int) {
	// save for when reset is called
	this.StartCell = [3]int{row, col, amount}

	// apply
	this.DataGrid.SetCell(row, col, amount)
}

// =================================================================================================================================================
// CONTROL
// -------------------------------------------------------------------------------------------------------------------------------------------------

func (this *Player) Reset() { 
	// Empty datagrid
	this.DataGrid.Clear()

	// Reapply startcell
	this.DataGrid.SetCell(this.StartCell[0], this.StartCell[1], this.StartCell[2])

	// [note: idea] opt-in to use random starting positions
}

func (this *Player) KillAll() { 
	this.DataGrid.KillAll()
}


// =================================================================================================================================================
// GAME LOOP
// -------------------------------------------------------------------------------------------------------------------------------------------------

/* cfg.KEY_AMOUNT +- growth --> cfg.KEY_I_AMOUNT */
func (this *Player) Grow(done chan bool, row_start, row_end int, Config *cfg.Config) {  
	// references
	var NeighbourLUT = &(this.DataGrid.NeighbourLUT)
	var Cells = &(this.DataGrid.Cells)

	// Internal vars
	var nbs = 0
	var nb_row int
	var nb_col int
	var nb_exists int
	var nb_nz_amount int
	var nb_nz_amount_sum int
	var growth = 1

	for row := row_start; row < row_end; row++ {
		for col := 0; col < cfg.COLS; col++ {
			// Init
			if Config.FlashyEel == false {
				growth = 1
			}
			nb_nz_amount_sum = 0

			// only grow cells with at least 1 amount
			if (*Cells)[row][col][cfg.KEY_AMOUNT] == 0 {
				continue
			}

			// grow if more than two friendly neighbours, else: shrink
			nbs = 0
			for i := 0; i < 4; i++ {
				nb_exists = (*NeighbourLUT)[row][col][i][cfg.LUTKEY_EXISTS]
				if nb_exists == 0 { continue }

				nb_row  = (*NeighbourLUT)[row][col][i][cfg.LUTKEY_ROW]
				nb_col  = (*NeighbourLUT)[row][col][i][cfg.LUTKEY_COL]
				nb_nz_amount = (*Cells)[nb_row][nb_col][cfg.KEY_AMOUNT]
				if nb_nz_amount > 0 {
					nbs ++
					nb_nz_amount_sum += nb_nz_amount
				}
			}

			// only grow if smaller than 200
			if (*Cells)[row][col][cfg.KEY_AMOUNT] > 200 {
				growth = 0
			}

			// if too crowded, don't grow
			if nb_nz_amount_sum > 800 {
				growth = 0
			}

			// set intm amount, max at 255, min at 0
			(*Cells)[row][col][cfg.KEY_I_AMOUNT] = misc.Normalize((*Cells)[row][col][cfg.KEY_AMOUNT] + growth, math.MaxInt64, 0) 			
		}
	}

	done <- true
}

/* 	Looks around at its neighbours, and determines where to send resources
	The change in resources is tracked in the intermediate amount register (KEY_I_AMOUNT)
	This amount will be written to KEY_AMOUNT by Battle() (after the battle of course)
*/
/* cfg.KEY_AMOUNT --> cfg.KEY_I_AMOUNT */
func (this *Player) Move(done chan bool, row_start, row_end int, f float64) {  

	if this.UserId == 1 {
		for row := row_start; row < row_end; row++ {
			for col := 0; col < cfg.COLS; col++ {
				
				this.UpdateIntermediateAmount(row, col, f)
			}
		}
	} else {
		for row := row_start; row < row_end; row++ {
			for col := 0; col < cfg.COLS; col++ {
				this.UpdateIntermediateAmount(row, col, f)
			}
		}		
	}

	done <- true

	// Don't forget to write I_AMOUNT back to AMOUNT after this function has been called!
	// (This used to be here, but it was moved to main function, so that red and green can share a loop)
}
func (this *Player) UpdateIntermediateAmount(row, col int, f float64) {
	// references
	var Cells = &(this.DataGrid.Cells)
	var EnemyCells = &(this.DataGrid.Enemy.Cells)
	var NeighbourLUT = &(this.DataGrid.NeighbourLUT)
		
	// skip if amount < 2
	if (*Cells)[row][col][cfg.KEY_AMOUNT] < 2 {
		return
	}		
	amount := (*Cells)[row][col][cfg.KEY_AMOUNT]

	// Ephemeral vars (used in loops)
	var nb_row int
	var nb_col int
	var nb_exists int

	// Keep track of how much to send over, can be changed based on logic
	var exp_amount uint8

	// Catalog of neighbours that are exceptional in one way or another
	least_friendly_smell := math.MaxInt64
	least_friendly_smell_nb := 0

	most_enemy_smell := 0
	most_enemy_smell_nb := -1

	least_nz_enemy_amount := math.MaxInt64
	least_nz_enemy_amount_nb := -1	

	least_friendly_amount := math.MaxInt64
	least_friendly_amount_nb := -1

	least_amount := math.MaxInt64
	least_amount_nb := -1
	sum_amount := 0

	// When we pick a neighbour, we can set it here
	target_neighbour := -1

	// determine amount to send
	exp_amount = uint8(amount / 2)		//exp_amount = uint8(amount - 1)

	// get list of neighbours, and their data
	nbs := [4][7]int{} 
	for i := 0; i < 4; i++ {
		nb_exists = (*NeighbourLUT)[row][col][i][cfg.LUTKEY_EXISTS]
		if nb_exists == 0 { continue }
		nb_row  = (*NeighbourLUT)[row][col][i][cfg.LUTKEY_ROW]
		nb_col  = (*NeighbourLUT)[row][col][i][cfg.LUTKEY_COL]

		// Copy neighbour over into shorthand list
		nbs[i][0] = nb_row 
		nbs[i][1] = nb_col 
		nbs[i][2] = nb_exists

		// Lookup extra information about our neighbour
		nbs[i][3] = (*Cells)[nb_row][nb_col][cfg.KEY_AMOUNT]			// intermediate amount friendly
		nbs[i][4] = (*Cells)[nb_row][nb_col][cfg.KEY_SMELL] 				// smell friendly
		nbs[i][5] = (*EnemyCells)[nb_row][nb_col][cfg.KEY_AMOUNT]		// intermediate amount enemy
		nbs[i][6] = (*EnemyCells)[nb_row][nb_col][cfg.KEY_SMELL]		// smell enemy

		
		if nbs[i][3] + int(exp_amount) > math.MaxInt64 {
			continue
		}

		if nbs[i][3] < least_friendly_amount || nbs[i][3] == least_friendly_amount && f > 0.5 {
			least_friendly_amount = nbs[i][3]
			least_friendly_amount_nb = i
		}

		if nbs[i][4] < least_friendly_smell || nbs[i][4] == least_friendly_smell && f > 0.5{
			least_friendly_smell = nbs[i][4]
			least_friendly_smell_nb = i
		}

		if nbs[i][6] > 0 && (nbs[i][6] > most_enemy_smell || nbs[i][6] == most_enemy_smell && f > 0.5 ) {
			most_enemy_smell = nbs[i][6]
			most_enemy_smell_nb = i
		}	
		
		if nbs[i][5] > 0 && (nbs[i][5] < least_nz_enemy_amount || nbs[i][5] == least_nz_enemy_amount && f > 0.5) {
			least_nz_enemy_amount = nbs[i][6]
			least_nz_enemy_amount_nb = i
		}

		sum_amount = nbs[i][3] + nbs[i][5]
		if sum_amount < least_amount || sum_amount == least_amount && f > 0.5 {
			least_amount = sum_amount
			least_amount_nb = i
		}
	}

	
	// pick neighbour with highest enemy smell
	if most_enemy_smell_nb != -1 {
		target_neighbour = most_enemy_smell_nb

	} else if least_nz_enemy_amount_nb != -1 {
		target_neighbour = least_nz_enemy_amount_nb

	} else if least_amount_nb != -1{
			target_neighbour = least_amount_nb
		
	// pick neighbour with least friendly smell
	} else if least_friendly_smell_nb != -1 {
		target_neighbour = least_friendly_smell_nb

	// pick neighbour with least friendly amount (and half exp amount)
	} else if least_friendly_amount_nb != -1 {
		exp_amount /= 2
		target_neighbour = least_friendly_amount_nb
	} 


	// update intermediate amount
	if target_neighbour != -1 {
		// expedition
		var target_row = nbs[target_neighbour][cfg.LUTKEY_ROW]
		var target_col = nbs[target_neighbour][cfg.LUTKEY_COL]

		(*Cells)[row][col][cfg.KEY_I_AMOUNT] -= int(exp_amount)
		(*Cells)[target_row][target_col][cfg.KEY_I_AMOUNT] = misc.Max255Int((*Cells)[target_row][target_col][cfg.KEY_I_AMOUNT], int(exp_amount))

	} 
}
func (this *Player) UpdateIntermediateAmountRandom(row, col int, f float64) {
	// references
	var Cells = &(this.DataGrid.Cells)
	var NeighbourLUT = &(this.DataGrid.NeighbourLUT)
		
	// skip if amount < 2
	if (*Cells)[row][col][cfg.KEY_AMOUNT] < 2 {
		return
	}		
	amount := (*Cells)[row][col][cfg.KEY_AMOUNT]


	// Keep track of how much to send over, can be changed based on logic
	var exp_amount uint8

	// When we pick a neighbour, we can set it here
	target_neighbour := -1

	// determine amount to send
	exp_amount = uint8(amount / 2)		//exp_amount = uint8(amount - 1)

	// pick random existing neighbour
	r := 0.1	
	for target_neighbour == -1 {
		for i := 0; i < 4; i++ {
			if f <= r && (*NeighbourLUT)[row][col][i][cfg.LUTKEY_EXISTS] == 1 {
				target_neighbour = i 
				break
			}
			r += 0.1
		}
	}

	// update intermediate amount
	// expedition
	var target_row = (*NeighbourLUT)[row][col][target_neighbour][cfg.LUTKEY_ROW]
	var target_col = (*NeighbourLUT)[row][col][target_neighbour][cfg.LUTKEY_COL]

	(*Cells)[row][col][cfg.KEY_I_AMOUNT] -= int(exp_amount)
	(*Cells)[target_row][target_col][cfg.KEY_I_AMOUNT] = misc.Max255Int((*Cells)[target_row][target_col][cfg.KEY_I_AMOUNT], int(exp_amount))
}


/*	Looks at the intermediate amount of a cell, and compares it to the same cell of the enemy
	When both have a non-zero amount, the result is calculated
	In either case: write (resultant) intermediate amount back to amount 
		(cfg.KEY_I_AMOUNT --> ?? --> cfg.KEY_AMOUNT )
*/
func (this *Player) Battle(f float64) {  
	done_top := make(chan bool)
	done_bottom := make(chan bool)

	go this.Battle_Goroutine(done_top, 0, cfg.ROWS/2, f)
	go this.Battle_Goroutine(done_bottom, cfg.ROWS/2, cfg.ROWS, f)

	<- done_top
	<- done_bottom
}
func (this *Player) Battle_Goroutine(done chan bool, row_start, row_end int, f float64) {  
	// references
	var Cells = &(this.DataGrid.Cells)
	var EnemyCells = &(this.DataGrid.Enemy.Cells)	
	
	// Ephemeral vars
	var enemy_amount int
	var own_amount int

	for row := row_start; row < row_end; row++ {
		for col := 0; col < cfg.COLS; col++ {
		
			// UPDATE KEY_AMOUNT
			// --------------------------------------------------------------
			// get intermediate amount
			own_amount = (*Cells)[row][col][cfg.KEY_I_AMOUNT]
			enemy_amount = (*EnemyCells)[row][col][cfg.KEY_I_AMOUNT]

			// write back to cfg.KEY_AMOUNT, so that we can exit out early if no battle takes place
			this.DataGrid.SetCell(row, col, own_amount)
			this.DataGrid.Enemy.SetCell(row, col, enemy_amount)

			// no battle will change the update above, continue to next cell
			if own_amount == 0 || enemy_amount == 0 {
				continue
			}

			// BATTLE
			// --------------------------------------------------------------
			// Tie --> make new amounts, with random win, and then continue to the battle
			if enemy_amount == own_amount {
				enemy_amount = int(float64(enemy_amount) * f)
				own_amount = int(float64(enemy_amount) * (1-f))
			} 
			
			// Enemy wins
			if enemy_amount > own_amount {
				// total defeat
				if enemy_amount >= own_amount * 2 {
					this.DataGrid.Kill(row, col)										// we died
					this.DataGrid.Enemy.SetCell(row, col, enemy_amount)					// enemy can keep all their amounts
				// good defeat
				} else {
					this.DataGrid.Kill(row, col)										// we died
					this.DataGrid.Enemy.SetCell(row, col, enemy_amount - own_amount)	// enemy loses what we had
				}

			// We win
			} else {
				// total victory
				if own_amount >= enemy_amount * 2 {
					this.DataGrid.Enemy.Kill(row, col)									// enemy loses all
					this.DataGrid.SetCell(row, col, own_amount)							// we can keep all our amounts
				// good defeat
				} else {
					this.DataGrid.Enemy.Kill(row, col)									// enemy is still dead
					this.DataGrid.SetCell(row, col, own_amount - enemy_amount)			// we lose the amount that the enemy had
				}				
			}
		}
	}

	done <- true
}


// =================================================================================================================================================
// DRAW
// -------------------------------------------------------------------------------------------------------------------------------------------------
/* goroutine wrapper for datagrid.Draw_PixelBased() */
func (this *Player) DrawPixels(done chan bool, numbers_text *[256]*text.TextObject, print_val int, print_text bool) {

	done_top := make(chan bool)
	done_bottom := make(chan bool)

	midpoint := cfg.ROWS / 2

	// Draw needs a lot of data from the DataGrid, so it is housed there
	if cfg.CELL_SIZE > 2 {
		// Erase pixels
		graphicsx.Clear(this.DataGrid.Pixels)

		// The method below draws only cells with a nonzero value, so the clear above is necessary
		go this.DataGrid.Draw_PixelBased(done_top, 0, midpoint, numbers_text, print_val, print_text)
		go this.DataGrid.Draw_PixelBased(done_bottom, midpoint, cfg.ROWS, numbers_text, print_val, print_text)

	} else {
		// The method below only draws the topleft corner of a cell, so perfect for when cell_size == 1
		// Draws all cells
		go this.DataGrid.Draw_PixelBasedDotted(done_top, 0, midpoint, numbers_text, print_val, print_text)
		go this.DataGrid.Draw_PixelBasedDotted(done_bottom, midpoint, cfg.ROWS, numbers_text, print_val, print_text)
	}

	<- done_top
	<- done_bottom
	
	// Let channel know that we are done with our calculation
	done <- true
}

/* Couldn't put this in with DrawPixels :( SDL doesn't cooperate well with goroutines */
func (this *Player) UpdateTexture() {
	var ScreenWidth = cfg.COLS * cfg.CELL_SIZE
	this.Texture.Update(nil, this.Pixels, int(ScreenWidth*4))
}

