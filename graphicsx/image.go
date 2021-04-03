package graphicsx

// Import external packages
import (
	"github.com/veandco/go-sdl2/sdl"
	//"github.com/veandco/go-sdl2/img"
)

// =====================================================================
// 				Struct: Image
// =====================================================================

// Tidy little package to contain one loaded image
type Image struct {
	Texture *sdl.Texture
	Width int32
	Height int32
}