package visualize

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	_ "image/jpeg"
	"log"
	"time"
)
import (
	"github.com/hajimehoshi/ebiten/v2"
	//"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	//"github.com/hajimehoshi/ebiten/v2/inpututil"
	//"bytes"
)
import _ "image/png"

type Game struct {
	screen *ebiten.Image
}

const (
	//screenWidth  = 640
	//screenHeight = 480
	screenWidth  = 1280
	screenHeight = 720
)

//func init() {
//WhiteStone, _, err := ebitenutil.NewImageFromFile("img/whiteStone.jpg")
//_, _, err = ebitenutil.NewImageFromFile("img/blackStone.jpg")
//if err != nil {
//	log.Fatal(err)
//}
//}

var _ = ebiten.NewImage(screenWidth, screenHeight)

//var WhiteStone, _, err_ = ebitenutil.NewImageFromFile("img/whiteStone_copy.png")

var WhiteStone *ebiten.Image

func init() {
	var err error
	//WhiteStone, _, err = ebitenutil.NewImageFromFile("img/whiteStone.jpg")
	//WhiteStone, _, err = ebitenutil.NewImageFromFile("img/whiteStone_.png")
	//WhiteStone, _, err = ebitenutil.NewImageFromFile("img/whiteStone_copy.png")
	//WhiteStone, _, err = ebitenutil.NewImageFromFile("img/whiteStone_copy_8.png")
	WhiteStone, _, err = ebitenutil.NewImageFromFile("img/WhiteStone.png")
	if err != nil {
		log.Fatal(err)
	}
}

//var WhiteStone, _, _ = ebitenutil.NewImageFromFile("img/whiteStone.jpg")

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Write your game's logical update.
	//mx, my := ebiten.CursorPosition()
	////fmt.Println(mx, my)
	//if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
	//	op := &ebiten.DrawImageOptions{}
	//	op.GeoM.Translate(float64(mx), float64(my))
	//	screen.DrawImage(WhiteStone, op)
	//}
	fmt.Println("test update", time.Now())

	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{R: 222, G: 184, B: 135, A: 255})
	gridColor64 := &color.RGBA{A: 50}
	//if WhiteStone == nil {
	//	fmt.Println(err_)
	//}
	//gridColor32 := &color.RGBA{A: 20}
	//const w = screenWidth
	//const h = screenHeight
	//for y := 0.0; y < h; y += 32 {
	//	ebitenutil.DrawLine(screen, 0, y, w, y, gridColor32)
	//}
	//for y := 0.0; y < h; y += 64 {
	//	ebitenutil.DrawLine(screen, 0, y, w, y, gridColor64)
	//}
	//for x := 0.0; x < w; x += 32 {
	//	ebitenutil.DrawLine(screen, x, 0, x, h, gridColor32)
	//}
	//for x := 0.0; x < w; x += 64 {
	//	ebitenutil.DrawLine(screen, x, 0, x, h, gridColor64)
	//}

	width := 12
	lines := 19
	start := 10
	//plus := 32
	xStart, yStart := start, start
	xEnd, yEnd := xStart+width*(lines-1), yStart+width*(lines-1)
	lines = 19

	for i := 0; i < lines; i++ {
		ebitenutil.DrawLine(screen, float64(xStart), float64(yStart), float64(xStart), float64(yEnd), gridColor64)
		//fmt.Println(i)
		xStart += width
	}
	xStart = start
	for i := 0; i < lines; i++ {
		ebitenutil.DrawLine(screen, float64(xStart), float64(yStart), float64(xEnd), float64(yStart), gridColor64)
		//fmt.Println(i)
		yStart += width
	}
	//for i in range(lines):  // for x
	//pygame.draw.lines(display_screen, (0, 0, 0), True, ((x_start, y_start), (x_start, y_end)))
	//x_start += width
	//
	//x_start = start
	//for i in range(lines):  // for y
	//pygame.draw.lines(display_screen, (0, 0, 0), True, ((x_start, y_start), (x_end, y_start)))
	//y_start += width
	//
	//pygame.display.update()

	//screen.DrawImage(WhiteStone, nil)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(10), float64(10))
	//op.ColorM.Scale(1.0, 0.25, 0.25, 1.0)
	//theta := 2.0 * math.Pi * float64(count%60) / 60.0
	//op.ColorM.Concat(ebiten.RotateHue(theta))
	screen.DrawImage(WhiteStone, op)

	//mx, my := ebiten.CursorPosition()
	////fmt.Println(mx, my)
	//if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
	//	op = &ebiten.DrawImageOptions{}
	//	op.GeoM.Translate(float64(mx), float64(my))
	//	screen.DrawImage(WhiteStone, op)
	//}
	// Write your game's rendering.

	fmt.Println("test draw", time.Now())
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 360
	//return screenWidth/2, screenHeight/2
}

func (g *Game) Draw_(screen *ebiten.Image) {
	// Рисуем сетку (32x32 и 64x64)
	//gridColor64 := &color.RGBA{A: 50}
	//gridColor32 := &color.RGBA{A: 20}
	//for y := 0.0; y < h; y += 32 {
	//	ebitenutil.DrawLine(screen, 0, y, w, y, gridColor32)
	//}
	//for y := 0.0; y < h; y += 64 {
	//	ebitenutil.DrawLine(screen, 0, y, w, y, gridColor64)
	//}
	//for x := 0.0; x < w; x += 32 {
	//	ebitenutil.DrawLine(screen, x, 0, x, h, gridColor32)
	//}
	//for x := 0.0; x < w; x += 64 {
	//	ebitenutil.DrawLine(screen, x, 0, x, h, gridColor64)
	//}
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
	//mx, my := ebiten.CursorPosition()
	//fmt.Println(mx, my)
	//if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
	//
	//}

}
