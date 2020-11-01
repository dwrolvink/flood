package main
// =====================================================================
// 				Imports
// =====================================================================
// Import built-in packages
import (
	"fmt"        // used for outputting to the terminal
	"time"       // used for pausing, measuring duration, etc
	"math/rand"  // random number generator
)

// Import external packages
import (
	"github.com/veandco/go-sdl2/sdl"
)

// subpackages
import (
	"flood_go/graphicsx"
	"flood_go/world"
	"flood_go/text"
	"flood_go/config"
	"flood_go/cell"
)

// Define constants
const (
	MOUSE_LEFT_CLICK = 1
	MOUSE_MIDDLE_CLICK = 2
	MOUSE_RIGHT_CLICK = 3
	BUTTON_DOWN = 0
	BUTTON_UP = 1

	DRAW_REGISTER_1 = 1
	DRAW_REGISTER_2 = 2
)

var (
	// This variable allows us to exit the otherwise endless loop when we want
	running = true

	// Pause game execution
	paused = false

	color_mode = 0
	SCROLL = 1

	// list of actions to be undertaken
	expeditions = []cell.Expedition{}
)

func load_users(cell_grid *[][]cell.Cell, cfg config.Config) {
	// pause main loop
	paused = true
	time.Sleep( 5 * time.Millisecond * cfg.SleepMiliSeconds)		

	// Set start location user2
	(*cell_grid)[cfg.Rows-1][cfg.Cols-1].BaseColor = &sdl.Color{255,0,0,255}	
	(*cell_grid)[cfg.Rows-1][cfg.Cols-1].UserId = 2
	(*cell_grid)[cfg.Rows-1][cfg.Cols-1].Amount = 20
	(*cell_grid)[cfg.Rows-1][cfg.Cols-1].CanAct = true	

	// Set start location user1
	(*cell_grid)[0][0].BaseColor = &sdl.Color{0,255,0,255}	
	(*cell_grid)[0][0].UserId = 1
	(*cell_grid)[0][0].Amount = 20
	(*cell_grid)[0][0].CanAct = true

	paused = false

}



// This is the entry point for our app. Code execution starts here.

func main() {

	// Get game config
	var cfg = config.GetConfig()

	// ========= Init step =========

	// Load SDL2, and get window and renderer.
	// See the file graphicsx/graphicsx.go for more information on the
	// graphics struct, and the initialization steps.

	// Endpoint is that we have a window object that we write to (and can
	// close). And a renderer object, which does the writing.
	graphics := graphicsx.Initialize_graphics()

	var renderer = graphics.Renderer
	var window = graphics.Window

	// Load images into memory
	graphics.LoadImage("src/images/icon.png") // --> graphics.Images[0]
	graphics.LoadImage("src/images/cat.png")  // --> graphics.Images[1]


	// Get screen dimensions so that we can position images relative to the corners.
	// Note that the screen dimensions are dictated in config/config.go, this code
	// is here to show you can get it from the window object too.
	screenWidth, screenHeight := window.GetSize()

	// Create grid of rectangles. These will be drawn at random in black
	// in a later step.
	var cell_grid = world.CreateRectGrid(cfg.Rows, cfg.Cols, cfg.CellSize)


	load_users(&cell_grid, cfg)

	// Define variables outside of loop, so that we don't have to recreate
	// them every iteration.
	var event sdl.Event


	var debug_text = text.NewTextObject(text.TextObjectConfig{
		Graphics: &graphics, 
		Text: "Press a key to show keyevent",
		Font: "SourceCodePro-Regular.ttf", 
		FontSize: 12,
		Color: &sdl.Color{0, 0, 0, 255},
		BgColor: &sdl.Color{255, 255, 255, 255},
	})	

	// The hello text never changes, so we can just statically define a Rect.
	// With the debug text though, the length of the text will change, and thus
	// also the size of the resulting Rect. If we then draw it with the smaller
	// Rect, the image will be squished. Also, because we want to horizontally 
	// center the text, this will need to be recalculated too.
	// To accommodate for this, we add a function to the struct that defines
	// how to make a new Rect on the fly. 
	debug_text.UpdateRect = func(textobj *text.TextObject) {  	
		textobj.Rect = &sdl.Rect{
			(screenWidth - textobj.Image.Width) / 2, 
			screenHeight - 20, 
			textobj.Image.Width, textobj.Image.Height,
		}
	}
	// Now we can update the Rect whenever we change the text and generate a new
	// Image.
	debug_text.UpdateRect(debug_text)
	
	// Define variables outside of loop that we want to increment/decrement
	// every iteration
	var current_cell *cell.Cell
	var current_exp cell.Expedition

	// ========= Game loop =========

	// Endless loop unless is running set to false
	// One iteration of this loop is one draw cycle.
	for running	{

		// Sleep a little so that we go the speed that we want
		time.Sleep(time.Millisecond * cfg.SleepMiliSeconds)		

		// when paused, skip the game actions
		if paused == false {

			// set draw color to white
			renderer.SetDrawColor(255, 255, 255, 255)                        // red, green, blue, alpha (alpha = transparency)

			// clear the window with specified color - in this case white.
			renderer.Clear()

			// set alpha
			var alpha uint8

			// Draw & Rearm Loop
			for row := range cell_grid {
				for col := range cell_grid[row] {
					// get cell 
					current_cell = &cell_grid[row][col]

					// don't do anything for empty cells
					if current_cell.UserId == 0 {
						continue
					}

					// set Cell color
					// set alpha
					if current_cell.Amount < 200 {
						alpha = 55 + uint8(current_cell.Amount)
					} else {
						alpha = 255
					}
					// set color
					graphics.SetSDLDrawColor(current_cell.BaseColor, alpha)

					// Change color based on color mode										
					switch color_mode {
					case 0:
						break
					case DRAW_REGISTER_1:
						if (current_cell.Register1 > 0.0) {
							alpha = 10 + uint8(current_cell.Register1*float64(SCROLL))
							if alpha > 255 { alpha = 255 }
							graphics.SetSDLDrawColor(&sdl.Color{255,0,255,255}, alpha)
						}
					case DRAW_REGISTER_2:
							alpha = 10 + uint8(current_cell.Register2*float64(SCROLL*10))
							if alpha > 255 { alpha = 255 }							
							graphics.SetSDLDrawColor(&sdl.Color{255,0,255,255}, alpha)
											
					}
						
					// draw cell
					renderer.FillRect(&current_cell.Rect)

					// rearm Cell
					current_cell.CanAct = true
					current_cell.Amount += 1
					
				}
			}		

			// Collect expeditions
			for row := range cell_grid {
				for col := range cell_grid[row] {
					// do cell action
					if (&cell_grid[row][col]).UserId != 0 {
						current_exp = (&cell_grid[row][col]).CallExpeditionFunction(&cell_grid)
						expeditions = append(expeditions, current_exp)
					}
					
				}
			}	
			
			// shuffle expeditions
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(expeditions), func(i, j int) { expeditions[i], expeditions[j] = expeditions[j], expeditions[i] })
			
			// execute expeditions
			for e := range expeditions {
				world.ExecuteExpedition(expeditions[e])
			}

			// clear expeditions
			expeditions = []cell.Expedition{}
		}


		// Draw debug text
		if (cfg.ShowDebugText){
			renderer.Copy(debug_text.Image.Texture, nil, debug_text.Rect)	
		}
		
		// Draw Screen
		// The rects have been drawn, now it is time to tell the renderer to show
		// what has been draw to the screen. "Present them."
		renderer.Present()

		// Handle events, in this case keyevents and close window
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
				
				// event that is sent when the window is closed
				case *sdl.QuitEvent:
					// setting running to false will end the game loop
					running = false

				// keydown/keyup events
				case *sdl.KeyboardEvent:

					// compile debug msg
					msg := fmt.Sprintf("[%d ms] screen_width:%d Keyboard, type:%d, sym:%c, code:%d modifiers:%d, state:%d, repeat:%d",
						t.Timestamp, screenWidth, t.Type, t.Keysym.Sym, t.Keysym.Scancode, t.Keysym.Mod, t.State, t.Repeat)
				
					// show on screen
					debug_text.SetText(msg) 
					// The above command  will automatically generate a new Image.
					// Because the size might be different, generate a new Rect.
					debug_text.UpdateRect(debug_text)

					// print in terminal
					//fmt.Println(msg)

					// on space: restart
					if t.Keysym.Sym == ' ' && t.State == 0{
						// clear expeditions
						expeditions = []cell.Expedition{}

						// clear grid
						cell_grid = world.CreateRectGrid(cfg.Rows, cfg.Cols, cfg.CellSize)

						// load users
						load_users(&cell_grid, cfg)
					}

					if t.Keysym.Sym == '1' && t.State == 1{
						color_mode = DRAW_REGISTER_1
					}
					if t.Keysym.Sym == '1' && t.State == 0{
						color_mode = 0
					}	
					
					if t.Keysym.Sym == '2' && t.State == 1{
						color_mode = DRAW_REGISTER_2
					}
					if t.Keysym.Sym == '2' && t.State == 0{
						color_mode = 0
					}						
					
					if t.Keysym.Scancode == 82 && t.State == 0{
						SCROLL += 100
						fmt.Println(SCROLL)
					}
					if t.Keysym.Scancode == 81 && t.State == 0{
						SCROLL -= 100
						fmt.Println(SCROLL)
					}					



				case *sdl.MouseButtonEvent:
					if (t.State == BUTTON_UP) {
						// convert mouse pos to row,col pos
						var pos = world.GetPos(t.X, t.Y, cfg.CellSize)

						// USER 1
						var uid uint8 = 1
						var baseColor = &sdl.Color{0,255,0,255}

						if t.Button == MOUSE_RIGHT_CLICK {
							// USER 2
							baseColor = &sdl.Color{255,0,0,255}
							uid = 2
						}

						// set user
						cell_grid[pos[0]][pos[1]].BaseColor = baseColor
						cell_grid[pos[0]][pos[1]].UserId = uid
						cell_grid[pos[0]][pos[1]].Amount = 20
						cell_grid[pos[0]][pos[1]].CanAct = true
					}
					
			}
		}
		
		//running = false
	} 
	
	// ========= End of Game loop =========

	// program is over, time to start shutting down. Keep in mind that sdl is written in C and does not have convenient
	// garbage collection like Go does
	graphics.Destroy()

}