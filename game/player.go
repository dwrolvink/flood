package game

// Import built-in packages
import (
	"math/rand"  // random number generator
	"fmt"
	"time"       // used for pausing, measuring duration, etc
	"math"
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

var (
	Config = cfg.GetConfig()
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
}

// =================================================================================================================================================
// INIT
// -------------------------------------------------------------------------------------------------------------------------------------------------

func (this *Player) Init(graphics *graphicsx.Graphics) { 
	// Set empty pixel array
	this.InitPixels()

	// Create Texture
	this.InitTexture(graphics)

	// Init sub structs
	this.DataGrid = &datagrid.DataGrid{BaseColor: this.BaseColor, Pixels: &this.Pixels}	
	this.DataGrid.Init()

	// Lame print, so we don't have to toggle imports when debugging
	fmt.Println("")
}

func (this *Player) InitPixels() {
	this.Pixels = make([]byte, Config.ScreenWidth*Config.ScreenHeight*4)
}

func (this *Player) InitTexture(graphics *graphicsx.Graphics) {
	texture, err := graphics.Renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(Config.ScreenWidth), int32(Config.ScreenHeight))
	if err != nil {
	  panic(err)
	}
	//defer texture.Destroy()

	texture.SetBlendMode(sdl.BLENDMODE_ADD)
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

	// [idea] opt-in to use random starting positions
}

func (this *Player) KillAll() { 
	this.DataGrid.KillAll()
}


// =================================================================================================================================================
// GAME LOOP
// -------------------------------------------------------------------------------------------------------------------------------------------------

/*	Full update cycle, minus the battle stage (can't do that concurrently)
	Calls:
		- Move()
		- Grow()
		- UpdateSmell()
*/
func (this *Player) UpdateInternalState(done chan bool) {
	// Grow existent cells
	this.Grow()

	// Update smell so the cells know what's around them
	// (Use different methods for different players)
	if this.UserId == 1 {
		this.DataGrid.UpdateSmell()
	} else {
		this.DataGrid.UpdateSmell()
	}
	
	// Say where the cell wants to move to 
	// (saved under intermediate Amount (KEY_I_AMOUNT))
	// The move will be finalized when .Battle() is called
	this.Move()

	// Let channel know that we are done with our calculation
	done <- true
}

/* cfg.KEY_AMOUNT --> cfg.KEY_I_AMOUNT */
func (this *Player) Move() {  

	rand.Seed(time.Now().UnixNano())
	f := rand.Float64()

	for row := 0; row < cfg.ROWS; row++ {
		for col := 0; col < cfg.COLS; col++ {
			if this.UserId == 1 {
				this.UpdateIntermediateAmount2(row, col, f)
			} else {
				this.UpdateIntermediateAmount2(row, col, f)
			}
		}
	}
}

/* 	Looks around at its neighbours, and determines where to send resources
	The change in resources is tracked in the intermediate amount register (KEY_I_AMOUNT)
	This amount will be written to KEY_AMOUNT by Battle() (after the battle of course)
*/
func (this *Player) UpdateIntermediateAmount(row, col int, f float64) {
	// skip if amount < 2
	if this.DataGrid.Cells[row][col][cfg.KEY_AMOUNT] < 2 {
		return
	}		
	amount := this.DataGrid.Cells[row][col][cfg.KEY_AMOUNT]

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
	most_enemy_smell_nb := 0

	least_nz_enemy_amount := math.MaxInt64
	least_nz_enemy_amount_nb := 0	

	least_friendly_amount := math.MaxInt64
	least_friendly_amount_nb := 0

	// When we pick a neighbour, we can set it here
	target_neighbour := -1

	// determine amount to send
	exp_amount = uint8(amount / 2)		//exp_amount = uint8(amount - 1)

	// get list of neighbours, and their data
	nbs := [4][7]int{} 

	for i := 0; i < 4; i++ {
		nb_exists = this.DataGrid.NeighbourLUT[row][col][i][cfg.LUTKEY_EXISTS]
		if nb_exists == 0 { continue }

		nb_row  = this.DataGrid.NeighbourLUT[row][col][i][cfg.LUTKEY_ROW]
		nb_col  = this.DataGrid.NeighbourLUT[row][col][i][cfg.LUTKEY_COL]

		// Copy neighbour over into shorthand list
		nbs[i][0] = nb_row 
		nbs[i][1] = nb_col 
		nbs[i][2] = nb_exists

		// Lookup extra information about our neighbour
		nbs[i][3] = this.DataGrid.Cells[nb_row][nb_col][cfg.KEY_AMOUNT]			// intermediate amount friendly
		nbs[i][4] = this.DataGrid.Cells[nb_row][nb_col][cfg.KEY_SMELL] 				// smell friendly
		nbs[i][5] = this.DataGrid.Enemy.Cells[nb_row][nb_col][cfg.KEY_AMOUNT]		// intermediate amount enemy
		nbs[i][6] = this.DataGrid.Enemy.Cells[nb_row][nb_col][cfg.KEY_SMELL]		// smell enemy

		
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

		if nbs[i][6] > most_enemy_smell || nbs[i][6] == most_enemy_smell && f > 0.5{
			most_enemy_smell = nbs[i][6]
			most_enemy_smell_nb = i
		}	
		
		if nbs[i][5] > 0 && nbs[i][5] < least_nz_enemy_amount || nbs[i][5] == least_nz_enemy_amount && f > 0.5{
			least_nz_enemy_amount = nbs[i][6]
			least_nz_enemy_amount_nb = i
		}
	}

	if least_nz_enemy_amount < math.MaxInt64 {
		target_neighbour = least_nz_enemy_amount_nb
	
	// pick neighbour with highest enemy smell
	} else if most_enemy_smell > 0 {
		target_neighbour = most_enemy_smell_nb

	// pick neighbour with least friendly smell
	} else if least_friendly_smell < math.MaxInt32 {
		target_neighbour = least_friendly_smell_nb

	// pick neighbour with least friendly amount (and half exp amount)
	} else if least_friendly_amount < math.MaxInt32 {
		exp_amount /= 2
		target_neighbour = least_friendly_amount_nb
	} 


	// update intermediate amount
	if target_neighbour != -1 {
		// expedition
		var target_row = nbs[target_neighbour][cfg.LUTKEY_ROW]
		var target_col = nbs[target_neighbour][cfg.LUTKEY_COL]

		this.DataGrid.Cells[row][col][cfg.KEY_I_AMOUNT] -= int(exp_amount)
		this.DataGrid.Cells[target_row][target_col][cfg.KEY_I_AMOUNT] = misc.Max255Int(this.DataGrid.Cells[target_row][target_col][cfg.KEY_I_AMOUNT], int(exp_amount))

	} 
}

func (this *Player) UpdateIntermediateAmount2(row, col int, f float64) {
	// skip if amount < 2
	if this.DataGrid.Cells[row][col][cfg.KEY_AMOUNT] < 2 {
		return
	}		
	amount := this.DataGrid.Cells[row][col][cfg.KEY_AMOUNT]

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
		nb_exists = this.DataGrid.NeighbourLUT[row][col][i][cfg.LUTKEY_EXISTS]
		if nb_exists == 0 { continue }
		nb_row  = this.DataGrid.NeighbourLUT[row][col][i][cfg.LUTKEY_ROW]
		nb_col  = this.DataGrid.NeighbourLUT[row][col][i][cfg.LUTKEY_COL]

		// Copy neighbour over into shorthand list
		nbs[i][0] = nb_row 
		nbs[i][1] = nb_col 
		nbs[i][2] = nb_exists

		// Lookup extra information about our neighbour
		nbs[i][3] = this.DataGrid.Cells[nb_row][nb_col][cfg.KEY_AMOUNT]			// intermediate amount friendly
		nbs[i][4] = this.DataGrid.Cells[nb_row][nb_col][cfg.KEY_SMELL] 				// smell friendly
		nbs[i][5] = this.DataGrid.Enemy.Cells[nb_row][nb_col][cfg.KEY_AMOUNT]		// intermediate amount enemy
		nbs[i][6] = this.DataGrid.Enemy.Cells[nb_row][nb_col][cfg.KEY_SMELL]		// smell enemy

		
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

		this.DataGrid.Cells[row][col][cfg.KEY_I_AMOUNT] -= int(exp_amount)
		this.DataGrid.Cells[target_row][target_col][cfg.KEY_I_AMOUNT] = misc.Max255Int(this.DataGrid.Cells[target_row][target_col][cfg.KEY_I_AMOUNT], int(exp_amount))

	} 
}
/*	Looks at the intermediate amount of a cell, and compares it to the same cell of the enemy
	When both have a non-zero amount, the result is calculated
	In either case: write (resultant) intermediate amount back to amount 
		(cfg.KEY_I_AMOUNT --> ?? --> cfg.KEY_AMOUNT )
*/
func (this *Player) Battle() {  

	// Random float (0-1) to decide who wins ties (and by how much)
	// We reuse the same value the entire round because the performance of this call is shite
	rand.Seed(time.Now().UnixNano())
	var f = rand.Float64()
	
	// Ephemeral vars
	var enemy_amount int
	var own_amount int
	var row int
	var col int

	// Construct array with N = number of cells
	// Then fill it with 0, 1, .., N
	// Then shuffle, so we go through the cells in random order
	var index = [cfg.ROWS * cfg.COLS]int{}
	for i := range index { index[i] = i}
	//rand.Shuffle(len(index), func(i, j int) { index[i], index[j] = index[j], index[i] })

	for _, i := range index {
		row = int(i / cfg.COLS)
		col = i - (cfg.COLS * row)
		
		// UPDATE KEY_AMOUNT
		// --------------------------------------------------------------
		// get intermediate amount
		own_amount = this.DataGrid.Cells[row][col][cfg.KEY_I_AMOUNT]
		enemy_amount = this.DataGrid.Enemy.Cells[row][col][cfg.KEY_I_AMOUNT]

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

/* cfg.KEY_AMOUNT */
func (this *Player) Grow() {  
	var nbs = 0
	var nb_row int
	var nb_col int
	var nb_exists int
	var nb_nz_amount int

	var growth = 0
	
	for row := 0; row < cfg.ROWS; row++ {
		for col := 0; col < cfg.COLS; col++ {

			// only grow cells with at least 1 amount
			if this.DataGrid.Cells[row][col][cfg.KEY_AMOUNT] == 0 {
				continue
			}

			// grow if more than two friendly nbs, else: shrink
			nbs = 0
			for i := 0; i < 4; i++ {
				nb_row  = this.DataGrid.NeighbourLUT[row][col][i][cfg.LUTKEY_ROW]
				nb_col  = this.DataGrid.NeighbourLUT[row][col][i][cfg.LUTKEY_COL]
				nb_exists = this.DataGrid.NeighbourLUT[row][col][i][cfg.LUTKEY_EXISTS]

				if nb_exists == 0 { continue }

				nb_nz_amount = this.DataGrid.Cells[nb_row][nb_col][cfg.KEY_AMOUNT]
				if nb_nz_amount > 0 {
					nbs ++
				}
			}

			growth = nbs - 2

			// only grow if smaller than 200
			if this.DataGrid.Cells[row][col][cfg.KEY_AMOUNT] > 200 {
				growth = 0
			}

			if nbs < 2 {
				growth = -4
			}

			this.DataGrid.SetCell(row, col, misc.Normalize(this.DataGrid.Cells[row][col][cfg.KEY_AMOUNT] + growth, math.MaxInt64, 0) )			
		}
	}
}


// =================================================================================================================================================
// DRAW
// -------------------------------------------------------------------------------------------------------------------------------------------------
/* goroutine wrapper for datagrid.Draw_PixelBased() */
func (this *Player) DrawPixels(done chan bool, numbers_text *[256]*text.TextObject, print_val int, print_text bool) {
	// Draw needs a lot of data from the DataGrid, so it is housed there
	this.DataGrid.Draw_PixelBased(numbers_text, print_val, print_text)

	// Let channel know that we are done with our calculation
	done <- true
}

/* Couldn't put this in with DrawPixels :( SDL doesn't cooperate well with goroutines */
func (this *Player) UpdateTexture() {
	this.Texture.Update(nil, this.Pixels, int(Config.ScreenWidth*4))
}

