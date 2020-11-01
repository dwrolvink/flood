package cell

// Import built-in packages
import (
	"math/rand"  // random number generator
	//"fmt"
	"time"       // used for pausing, measuring duration, etc
)

// Import external packages
import (
	"github.com/veandco/go-sdl2/sdl"
)

type Expedition struct {
	Source *Cell
	Target *Cell
	Amount uint
	Register1 float64
	Register2 float64
	Register3 float64
}

type Cell struct {
	Rect sdl.Rect
	Col int
	Row int
	BaseColor *sdl.Color
	Amount uint				
	UserId uint8			// 0 for no user
	Register1 float64
	Register2 float64
	Register3 float64
	CanAct bool
	Nbc []*Cell
}

func (this *Cell) SetNeighbours(cell_grid *[][]Cell) {  
	// get neighbours
	var nbc []*Cell
	// top
	if this.Row  > 0 {
		nbc = append(nbc, &(*cell_grid)[this.Row-1][this.Col])
	}
	// left
	if this.Col > 0 {
		nbc = append(nbc, &(*cell_grid)[this.Row][this.Col-1])
	}		
	// bottom
	if this.Row + 1 < len(*cell_grid) {
		nbc = append(nbc, &(*cell_grid)[this.Row+1][this.Col])
	}	
	// right
	if this.Col + 1 < len((*cell_grid)[this.Row]) {
		nbc = append(nbc, &(*cell_grid)[this.Row][this.Col+1])
	}
	
	// shuffle neighbours
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(nbc), func(i, j int) { nbc[i], nbc[j] = nbc[j], nbc[i] })

	this.Nbc = nbc
}


func (this *Cell) CallExpeditionFunction(cell_grid *[][]Cell) Expedition{  

	// run function
	return this.ExpeditionFunction1()

}




func (this *Cell) ExpeditionFunction1() Expedition { 
	// shortcuts
	var nbc = this.Nbc

	// Output
	var exp = Expedition{}
	exp.Source = this

	// UPDATE REGISTERS
	// -----------------------------------------------------------------
	// "enemy smell"
	
		exp.Register1 = 0.0
		for i := range nbc {
			if nbc[i].UserId != 0 && nbc[i].UserId != this.UserId {
				exp.Register1 += float64(nbc[i].Amount)

			} else if nbc[i].UserId == this.UserId {
				exp.Register1 += 0.2 * nbc[i].Register1
			}
		}

	
	
	// "empty smell"
	exp.Register2 = 0.0
	var strongest_empty_smell_strength = 0.0
	var number_of_empty_cells = 0
	for i := range nbc {
		if nbc[i].UserId == 0 {
			number_of_empty_cells += 1

		} else if nbc[i].UserId == this.UserId {
			if nbc[i].Register2 > strongest_empty_smell_strength {
				strongest_empty_smell_strength = nbc[i].Register2
			}
		}
	}
	if number_of_empty_cells > 0 {
		exp.Register2 = 1.0 	
	} else {
		exp.Register2 = 0.9 * strongest_empty_smell_strength
	}

	// SET AMOUNT TO SHIP
	// -----------------------------------------------------------------
	// test if cell has enough resources
	if this.Amount < 2 {
		exp.Amount = 0
		return exp
	}
	
	// set amount to ship
	exp.Amount = this.Amount - 1


	var weakest_enemy_cell *Cell

	// ATTACK WEAKEST ENENMY CELL
	// -----------------------------------------------------------------
	for i := range nbc {
		if nbc[i].UserId != 0 && nbc[i].UserId != this.UserId {
			if weakest_enemy_cell == nil {
				weakest_enemy_cell = nbc[i]
			} else if weakest_enemy_cell.Amount > nbc[i].Amount {
				weakest_enemy_cell = nbc[i]
			}
		}
		if (weakest_enemy_cell != nil) {
			exp.Target = weakest_enemy_cell
			return exp
		}
	}

	// ATTACK EMPTY CELL WHEN FOUND
	// -----------------------------------------------------------------
	for i := range nbc {
		if nbc[i].UserId == 0 {
			// ship to empty cell
			exp.Target = nbc[i]
			return exp
		}
	}		

	// COLLECT SMELLS 
	// -----------------------------------------------------------------
	var strongest_enemy_smell *Cell
	var strongest_empty_smell *Cell
	for i := range nbc {
		// friendly cell
		if nbc[i].UserId == this.UserId {
			// set friendly neighbor with strongest enemy smell
			if strongest_enemy_smell == nil {
				strongest_enemy_smell = nbc[i]

			} else {
				if nbc[i].Register1 > strongest_enemy_smell.Register1 {
					strongest_enemy_smell = nbc[i]
				}
			}
			
			// set friendly neighbor with strongest empty smell
			if strongest_empty_smell == nil {
				strongest_empty_smell = nbc[i]
			} else {
				if nbc[i].Register2 > strongest_empty_smell.Register2 {
					strongest_empty_smell = nbc[i]
				}
			}			
		}
	}	

	// SEND TO STRONGEST ENEMY SMELLING NEIGHBOUR
	// -----------------------------------------------------------------
	if strongest_enemy_smell != nil && strongest_enemy_smell.Register1 > 0.0 {
		exp.Target = strongest_enemy_smell
		return exp
	}
	
	// SEND TO STRONGEST EMPTY SMELLING NEIGHBOUR
	// -----------------------------------------------------------------
	if strongest_empty_smell != nil && strongest_empty_smell.Register2 > 0.0 {
		exp.Target = strongest_empty_smell
		return exp
	}	

	exp.Amount = 0
	return exp
	
}
