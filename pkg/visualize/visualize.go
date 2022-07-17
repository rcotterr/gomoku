package visualize

import (
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	_ "image/jpeg"
	"log"
)
import (
	"github.com/hajimehoshi/ebiten/v2"
	//"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	//"github.com/hajimehoshi/ebiten/v2/inpututil"
	//"bytes"
)

type Game struct{}

const (
	screenWidth  = 640
	screenHeight = 480
	//screenWidth  = 1600
	//screenHeight = 1040
)

var image = ebiten.NewImage(screenWidth, screenHeight)

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Write your game's logical update.
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{222, 184, 135, 255})
	//screen.DrawImage(image, nil)
	// Write your game's rendering.
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
	//return screenWidth/2, screenHeight/2
}

func init() {
	_, _, err := ebitenutil.NewImageFromFile("/Users/marina.romashkova/tests/gomoku_git/whiteStone.jpg")
	_, _, err = ebitenutil.NewImageFromFile("/Users/marina.romashkova/tests/gomoku_git/blackStone.jpg")
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Draw_(screen *ebiten.Image) {
	//const (
	//	ox = 10
	//	oy = 10
	//	dx = 60
	//	dy = 50
	//)
	//screen.Fill(color.NRGBA{0x00, 0x40, 0x80, 0xff})
	//
	//op := &ebiten.DrawImageOptions{}
	//op.GeoM.Translate(ox, oy)
	//screen.DrawImage(ebitenImage, op)
	//
	//// Fill with solid colors
	//for i, c := range colors {
	//	op := &ebiten.DrawImageOptions{}
	//	x := i % 4
	//	y := i/4 + 1
	//	op.GeoM.Translate(ox+float64(dx*x), oy+float64(dy*y))
	//
	//	// Reset RGB (not Alpha) 0 forcibly
	//	op.ColorM.Scale(0, 0, 0, 1)
	//
	//	// Set color
	//	r := float64(c.R) / 0xff
	//	g := float64(c.G) / 0xff
	//	b := float64(c.B) / 0xff
	//	op.ColorM.Translate(r, g, b, 0)
	//	screen.DrawImage(ebitenImage, op)
	//}
}

func Vis() {
	game := &Game{}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Gomoku")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
