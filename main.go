package main
// =====================================================================
// 				Imports
// =====================================================================

// Import built-in packages
import (
	//"os"
	"fmt"        // used for outputting to the terminal
	"time"       // used for pausing, measuring duration, etc
	//"math"
	"math/rand"  // random number generator
	"strconv"	 // int to string
)

// Import external packages
import (
	"github.com/veandco/go-sdl2/sdl"
)

// subpackages
import (
	"flood_go/graphicsx"
	"flood_go/text"
	"flood_go/game"
	//"flood_go/misc"
	cfg "flood_go/config"
)

var (
	// IMPORTS / INITS
	Config = cfg.GetConfig()

	// CONTROLS
	Running = true						// Setting this to false will exit main loop (exit program)
	Paused = false						// Pause game execution (main loop continues on)
	SkipDraw = false

	ValueOfInterest = cfg.KEY_AMOUNT	// Which cell value to color the cells with
	PrintText = false					// Whether to print text to show value (broken atm bc pixel printing is implemented)

	// GRAPHICS / WINDOW / OS
	Graphics = graphicsx.Initialize_graphics()	// See the graphics package for explanations
	Renderer = Graphics.Renderer
	Event sdl.Event

	// ASSETS
	NumberImages = CreateNumberImages()	// holds images for the numbers 0-255
	
	// GAME GLOBALS
	player_red game.Player
	player_green game.Player
)



// This is the entry point for our app. Code execution starts here.
func main() {

	InitPlayers()

	// Set the color that the screen will be cleared at
	// (This will do little now that we use pixel based drawing)
	Renderer.SetDrawColor(0, 0, 0, 0)           // red, green, blue, alpha (alpha = transparency)

	// Slow down game so it doesn't get jumpy
	var TickStart time.Time
	var TickDurationNs time.Duration

	// Create channels to track when each player is done
	// when using goroutines
	done_red_top := make(chan bool)
	done_green_top := make(chan bool)
	done_red_bottom := make(chan bool)
	done_green_bottom := make(chan bool)

	// Seed so we have 
	rand.Seed(time.Now().UnixNano())
	
		
	for Running	{	
		
		TickStart = time.Now()

		// when paused, skip the game/draw actions, but stay in the main loop and receive events
		if Paused == false {

			// Randomness
			f := rand.Float64()

			// Grow existent cells
			// ------------------------
			go player_red.Grow(done_red_top, 0, cfg.ROWS/2, &Config)
			go player_red.Grow(done_red_bottom, cfg.ROWS/2, cfg.ROWS, &Config)
			go player_green.Grow(done_green_top, 0, cfg.ROWS/2, &Config)
			go player_green.Grow(done_green_bottom, cfg.ROWS/2, cfg.ROWS, &Config)			

			<- done_red_top
			<- done_red_bottom
			<- done_green_top
			<- done_green_bottom	
			
			// Write intermediate amount to final amount
			/*
			for row := 0; row < cfg.ROWS; row++ {
				for col := 0; col < cfg.COLS; col++ {
					player_red.DataGrid.Cells[row][col][cfg.KEY_AMOUNT] = player_red.DataGrid.Cells[row][col][cfg.KEY_I_AMOUNT]
					player_green.DataGrid.Cells[row][col][cfg.KEY_AMOUNT] = player_green.DataGrid.Cells[row][col][cfg.KEY_I_AMOUNT]
				}
			}
			*/				

			// Update smell so the cells know what's around them
			// ------------------------
			// (Use different methods for different players)
			go player_red.DataGrid.UpdateIntermediateSmell(done_red_top, 0, cfg.ROWS/2, f)
			go player_red.DataGrid.UpdateIntermediateSmell(done_red_bottom, cfg.ROWS/2, cfg.ROWS, f)
			go player_green.DataGrid.UpdateIntermediateSmell(done_green_top, 0, cfg.ROWS/2, f)
			go player_green.DataGrid.UpdateIntermediateSmell(done_green_bottom, cfg.ROWS/2, cfg.ROWS, f)
		
			<- done_red_top
			<- done_red_bottom
			<- done_green_top
			<- done_green_bottom

			// Write intermediate smell to final smell
			for row := 0; row < cfg.ROWS; row++ {
				for col := 0; col < cfg.COLS; col++ {
					player_red.DataGrid.Cells[row][col][cfg.KEY_SMELL] = player_red.DataGrid.Cells[row][col][cfg.KEY_I_SMELL]
					player_green.DataGrid.Cells[row][col][cfg.KEY_SMELL] = player_green.DataGrid.Cells[row][col][cfg.KEY_I_SMELL]
				}
			}			
			
			// Say where the cell wants to move to 
			// ------------------------
			// (saved under intermediate Amount (cfg.KEY_I_AMOUNT))
			// The move to cfg.KEY_AMOUNT will be finalized when .Battle() is called
			go player_red.Move(done_red_top, 0, cfg.ROWS/2, f)
			go player_red.Move(done_red_bottom, cfg.ROWS/2, cfg.ROWS, f)
			go player_green.Move(done_green_top, 0, cfg.ROWS/2, f)
			go player_green.Move(done_green_bottom, cfg.ROWS/2, cfg.ROWS, f)			

			<- done_red_top
			<- done_red_bottom
			<- done_green_top
			<- done_green_bottom	
	
			// Battle it out on cells where both players have an intermediate amount
			// ------------------------
			// Else, just move cfg.KEY_I_AMOUNT into cfg.KEY_AMOUNT
			player_red.Battle(f)
					
			// Draw result
			// ------------------------
			if SkipDraw == false {
				DrawFrame()	
			}
		}

		// Handle events, in this case keyevents and close window
		for Event = sdl.PollEvent(); Event != nil; Event = sdl.PollEvent() {
			switch t := Event.(type) {
				
				// event that is sent when the window is closed
				case *sdl.QuitEvent:
					// setting running to false will end the game loop
					Running = false

				// keydown/keyup events
				case *sdl.KeyboardEvent:

					// print in terminal
					fmt.Println(string(t.Keysym.Sym), t.State)

					// on space: restart
					if t.Keysym.Sym == ' ' && t.State == 0{
						ResetGame()

					} else if t.Keysym.Sym == 'f' && t.State == 1{
						SkipDraw = true
						
					} else if t.Keysym.Sym == 'f' && t.State == 0{
						SkipDraw = false
						
					} else if t.Keysym.Sym == 'e' && t.State == 0{
						Config.FlashyEel = !Config.FlashyEel
						fmt.Println(Config.FlashyEel)
						
					}else if t.Keysym.Sym == '1' && t.State == 1{
						ValueOfInterest = cfg.KEY_I_AMOUNT

					} else if t.Keysym.Sym == '1' && t.State == 0{
						ValueOfInterest = cfg.KEY_AMOUNT

					} else if t.Keysym.Sym == '2' && t.State == 1{
						ValueOfInterest = cfg.KEY_SMELL

					} else if t.Keysym.Sym == '2' && t.State == 0{
						ValueOfInterest = cfg.KEY_SMELL
					} else if t.Keysym.Sym == '0' && t.State == 1{
						//print_text = ! print_text
						player_green.KillAll()

					} else if t.Keysym.Scancode == 82 && t.State == 0{
						//SCROLL += 100
						fmt.Println("SCROLL UP")

					} else if t.Keysym.Scancode == 81 && t.State == 0{
						//SCROLL -= 100
						fmt.Println("SCROLL DOWN")
					}					


				case *sdl.MouseButtonEvent:
					if (t.State == cfg.BUTTON_UP) {
						continue
					}
					
			}
		}
		
		// Track framerate		
		TickDurationNs = time.Since(TickStart)
		if TickDurationNs > cfg.INTERVAL_NS {
			//fmt.Println("Frame missed by ", TickDurationNs - cfg.INTERVAL_NS)
		} else {
			time.Sleep(time.Nanosecond * time.Duration(cfg.INTERVAL_NS - TickDurationNs))
		}	
	} 

	// ========= End of Game loop =========

	// Program is over, time to start shutting down. Keep in mind that sdl is written in C and does not have convenient
	// garbage collection like Go does
	player_red.Texture.Destroy()
	player_green.Texture.Destroy()	
	Graphics.Destroy()
}

func InitPlayers() {
	player_red = game.Player{
		UserId: 1,
		Name: "Red",
		BaseColor: &sdl.Color{255,0,0,255},
	}
	player_green = game.Player{
		UserId: 2,
		Name: "Green",
		BaseColor: &sdl.Color{0,255,0,255},
	}

	player_red.Init(&Graphics)
	player_green.Init(&Graphics)

	player_red.SetEnemy(&player_green)
	player_green.SetEnemy(&player_red)	
	
	player_green.SetStartCell(0, 0, 64)
	player_red.SetStartCell(int(cfg.ROWS - 1), int(cfg.COLS - 1), 64)
}

func ResetGame() {
	player_red.Reset()
	player_green.Reset()
}

func DrawFrame() {
	// Clear screen
	Renderer.Clear()

	// Build pixel arrays based on game data
	done_red := make(chan bool)
	done_green := make(chan bool)

	go player_red.DrawPixels(done_red, &NumberImages, ValueOfInterest, PrintText)
	go player_green.DrawPixels(done_green, &NumberImages, ValueOfInterest, PrintText)

	<- done_red
	<- done_green

	// Build textures from the pixel arrays, for each player
	// 		This would be together in one function with drawpixels, 
	// 		were it not that SDL functions cannot be called from goroutines
	player_red.UpdateTexture()	
	player_green.UpdateTexture()
	
	// Merge the separate textures into the screen texture
	// 		Notice the "texture.SetBlendMode(sdl.BLENDMODE_ADD)" in game.Player.InitTexture()
	//		This makes it so that the two textures are blended, rather than one overwriting the other
	Renderer.Copy(player_red.Texture, nil, nil)	
	Renderer.Copy(player_green.Texture, nil, nil)	

	// Update screen
	Renderer.Present()
}

// Should be housed under flood_go/text, but I can't be bothered atm
func CreateNumberImages() [256]*text.TextObject{
	var NumberImages [256]*text.TextObject
	for i := range NumberImages {

		NumberImages[i] = text.NewTextObject(text.TextObjectConfig{
			Graphics: &Graphics, 
			Text: strconv.Itoa(i),
			Font: "SourceCodePro-Regular.ttf", 
			FontSize: 10,
			Color: &sdl.Color{0, 0, 0, 255},
			BgColor: &sdl.Color{255, 255, 255, 0},
		})
	} //NumberImages[0].Image.Texture.SetBlendMode(sdl.BLENDMODE_BLEND) 
	return NumberImages	
}



func Fib(n int) int {
	if n < 2 {
			return n
	}
	return Fib(n-1) + Fib(n-2)
}