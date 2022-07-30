package visualize

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"gomoku/pkg/playboard"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"strings"
	"time"
)

const (
	width = 12
	lines = 19
	start = 30
)

type GameInterface interface {
	Update() error
	Draw(screen *ebiten.Image)
	Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
	GetPlayBoard() string
}

type MockGame struct{}

func (g *MockGame) Update() error {
	return nil
}

func (g *MockGame) Draw(_ *ebiten.Image) {}

func (g *MockGame) Layout(_, _ int) (_, _ int) {
	return 0, 0
}

type HumanGame struct {
	//screen *ebiten.Image
	playBoard     string
	currentPlayer *playboard.Player
	anotherPlayer *playboard.Player
	index         int
	isOver        bool
}

func (g HumanGame) GetPlayBoard() string { return g.playBoard }

type AIGame struct {
	//screen *ebiten.Image
	playBoard     string
	currentPlayer *playboard.Player
	anotherPlayer *playboard.Player
	index         int
	isOver        bool
}

func (g AIGame) GetPlayBoard() string { return g.playBoard }

const (
	screenWidth  = 1280
	screenHeight = 720
)

var _ = ebiten.NewImage(screenWidth, screenHeight)

var WhiteStone *ebiten.Image
var BlackStone *ebiten.Image

func init() {
	var err error

	BlackStone, _, err = ebitenutil.NewImageFromFile("img/BlackStone.png")
	if err != nil {
		log.Fatal(err)
	}
	WhiteStone, _, err = ebitenutil.NewImageFromFile("img/WhiteStone.png")
	if err != nil {
		log.Fatal(err)
	}
}

func HumanTurnVis(currentPlayer playboard.Player) (int, error) {
	mx, my := ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		mx, my = int(math.Round(float64(mx-start)/width)), int(math.Round(float64(my-start)/width))
		index := my*playboard.N + mx
		//fmt.Println("HERE X Y !", mx, my, index)
		if index >= 0 && index < playboard.N*playboard.N {
			return index, nil
		}
	}
	return -1, fmt.Errorf("no step")

}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *HumanGame) Update() error {

	if !g.isOver && !playboard.GameOver(g.playBoard, g.currentPlayer, g.anotherPlayer, g.index) {
		newIndex, err := HumanTurnVis(*g.currentPlayer)
		if err != nil {
			return nil
		}
		newPlayBoard, err := playboard.PutStone(g.playBoard, newIndex, g.currentPlayer)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		g.playBoard = newPlayBoard.Node
		g.index = newIndex
		g.currentPlayer.Captures += newPlayBoard.Captures
		playboard.PrintPlayBoard(g.playBoard)
		g.currentPlayer, g.anotherPlayer = g.anotherPlayer, g.currentPlayer
	} else {
		g.isOver = true
	}
	fmt.Println("test update", time.Now())

	return nil
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *AIGame) Update() error {

	if !g.isOver && !playboard.GameOver(g.playBoard, g.currentPlayer, g.anotherPlayer, g.index) {
		newIndex, err := HumanTurnVis(*g.currentPlayer)
		if err != nil {
			return nil
		}
		newPlayBoard, err := playboard.PutStone(g.playBoard, newIndex, g.currentPlayer)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		g.playBoard = newPlayBoard.Node
		g.index = newIndex
		g.currentPlayer.Captures += newPlayBoard.Captures
		playboard.PrintPlayBoard(g.playBoard)
		g.currentPlayer, g.anotherPlayer = g.anotherPlayer, g.currentPlayer
	} else {
		g.isOver = true
	}
	fmt.Println("test update", time.Now())

	return nil
}

func _draw(screen *ebiten.Image, g GameInterface) {
	screen.Fill(color.NRGBA{R: 222, G: 184, B: 135, A: 255})
	gridColor64 := &color.RGBA{A: 50}

	widthStone := 14
	xStart, yStart := start, start
	xEnd, yEnd := xStart+width*(lines-1), yStart+width*(lines-1)

	for i := 0; i < lines; i++ {
		ebitenutil.DrawLine(screen, float64(xStart), float64(yStart), float64(xStart), float64(yEnd), gridColor64)
		xStart += width
	}
	xStart = start
	for i := 0; i < lines; i++ {
		ebitenutil.DrawLine(screen, float64(xStart), float64(yStart), float64(xEnd), float64(yStart), gridColor64)
		yStart += width
	}
	yStart = start

	for index, stone := range g.GetPlayBoard() {
		if string(stone) != playboard.EmptySymbol {
			op := &ebiten.DrawImageOptions{}
			mx, my := index%playboard.N, index/playboard.N
			mx, my = xStart+width*(mx), yStart+width*(my)
			//fmt.Println(mx, my)
			op.GeoM.Translate(float64(mx-widthStone/2), float64(my-widthStone/2))

			if string(stone) != playboard.SymbolPlayer2 {
				screen.DrawImage(WhiteStone, op)
			} else {
				screen.DrawImage(BlackStone, op)
			}
		}
	}
	//if g.currentPlayer.Symbol == playboard.SymbolPlayer1 || g.currentPlayer.Symbol == playboard.SymbolPlayer2 {
	//	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Turn player (%s)", g.currentPlayer.Symbol), 300, 100)
	//	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Capture player (%s) is (%d)", g.currentPlayer.Symbol, g.currentPlayer.Captures), 300, 150)
	//	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Capture player (%s) is (%d)", g.anotherPlayer.Symbol, g.anotherPlayer.Captures), 300, 200)
	//}
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *HumanGame) Draw(screen *ebiten.Image) {
	_draw(screen, g)
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *AIGame) Draw(screen *ebiten.Image) {
	_draw(screen, g)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *HumanGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 360
	//return screenWidth/2, screenHeight/2
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *AIGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 360
	//return screenWidth/2, screenHeight/2
}

func NewHumanGame() GameInterface {
	game := &HumanGame{
		playBoard:     strings.Repeat(playboard.EmptySymbol, playboard.N*playboard.N),
		currentPlayer: &playboard.Player1,
		anotherPlayer: &playboard.Player2,
		index:         -1,
		isOver:        false,
	}
	return game
}

func NewAIGame() GameInterface {
	game := &AIGame{
		playBoard:     strings.Repeat(playboard.EmptySymbol, playboard.N*playboard.N),
		currentPlayer: &playboard.Player1,
		anotherPlayer: &playboard.Player2,
		index:         -1,
		isOver:        false,
	}
	return game
}

func Vis(game GameInterface) {

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Gomoku")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal()
	}
}
