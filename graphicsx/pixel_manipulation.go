package graphicsx

// =====================================================================
// 				Functions: Pixel based
// =====================================================================
/*
	Usage:
		graphics := graphicsx.Initialize_graphics()
		var renderer = graphics.Renderer
		tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(cfg.ScreenWidth), int32(cfg.ScreenHeight))
		if err != nil {
		panic(err)
		}
		defer tex.Destroy()

		pixels := make([]byte, cfg.ScreenWidth*cfg.ScreenHeight*4)		

		graphicsx.Clear(pixels)
		graphicsx.SetSquare(200, 200, int(cfg.CellSize), [3]byte{255,0,0}, pixels)
		graphicsx.SetPixel(200, 300, [3]byte{255,0,0}, pixels)
		
		tex.Update(nil, pixels, int(cfg.ScreenWidth*4))
		renderer.Copy(tex, nil, nil)
		renderer.Present()
*/

/* Set all pixels to 0 */
func Clear(pixels *[]byte) {
	for i := range (*pixels) {
	  (*pixels)[i] = 0
	}
}

/* Set all pixels to 0 */
func SetPixel(x, y int, c [4]byte, pixels *[]byte) {
	index := (y* int(cfg.ScreenWidth) + x) * 4

	if index < len(*pixels)-4 && index >= 0 {
		(*pixels)[index] = c[0]
		(*pixels)[index+1] = c[1]
		(*pixels)[index+2] = c[2]	
		(*pixels)[index+3] = c[3]	
	}
}
// setPixel(200, 200, [4]byte{255,0,0,255}, pixels)

// Same as SetPixel, but uses cell size to draw more than one (so to draw a square)
func SetSquare(x, y, size int, c [4]byte, pixels *[]byte) {
	if size == 1 {
		SetPixel(x, y, c, pixels)
		return
	}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			SetPixel(x + i, y + j, c, pixels)
		}
	}	
}
// setSquare(200, 200, 4, [4]byte{255,0,0,255}, pixels)