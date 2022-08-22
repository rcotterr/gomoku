package visualize

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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
	GetTurns() int
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
	turns         int
	forbiddenMove bool
}

func (g HumanGame) GetPlayBoard() string { return g.playBoard }
func (g HumanGame) GetTurns() int        { return g.turns }

type AIGame struct {
	//screen *ebiten.Image
	playBoard     string
	humanPlayer   *playboard.Player
	machinePlayer *playboard.Player
	machineTurn   bool
	index         int
	isOver        bool
	turns         int
	forbiddenMove bool
}

func (g AIGame) GetPlayBoard() string { return g.playBoard }
func (g AIGame) GetTurns() int        { return g.turns }

const (
	screenWidth  = 1280
	screenHeight = 720
)

var _ = ebiten.NewImage(screenWidth, screenHeight)

var WhiteStone *ebiten.Image
var BlackStone *ebiten.Image

var normalFont font.Face

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

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	normalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    12,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

}

func HumanTurnVis(currentPlayer playboard.Player) (int, error) {
	mx, my := ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		mx, my = int(math.Round(float64(mx-start)/width)), int(math.Round(float64(my-start)/width))
		index := my*playboard.N + mx
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
		//g.currentPlayer.Captures += newPlayBoard.Captures // its in PutStone
		playboard.PrintPlayBoard(g.playBoard)
		g.currentPlayer, g.anotherPlayer = g.anotherPlayer, g.currentPlayer
		g.turns += 1
	} else {
		g.isOver = true
	}
	fmt.Println("test update", time.Now())

	return nil
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *AIGame) Update() error {

	if !g.isOver && !playboard.GameOver(g.playBoard, g.machinePlayer, g.humanPlayer, g.index) {
		if g.machineTurn {
			g.index = playboard.Algo(g.playBoard, *g.machinePlayer, *g.humanPlayer)
			newPlayBoard, err := playboard.PutStone(g.playBoard, g.index, g.machinePlayer)
			if err != nil {
				fmt.Println("Invalid machine algo!!!!!", err)
				log.Fatal()
			}
			g.playBoard = newPlayBoard.Node
			playboard.PrintPlayBoard(g.playBoard) //TO DO delete print
			fmt.Println(g.machinePlayer)
			if g.machinePlayer.IndexAlmostWin != nil {
				fmt.Println(g.machinePlayer, *g.machinePlayer.IndexAlmostWin)
			}
			fmt.Println(g.humanPlayer)
			if g.humanPlayer.IndexAlmostWin != nil {
				fmt.Println(g.humanPlayer, *g.humanPlayer.IndexAlmostWin)
			}
			g.machineTurn = false
			g.turns += 1
			g.forbiddenMove = false
		} else {
			newIndex, err := HumanTurnVis(*g.humanPlayer)
			if err != nil {
				return nil
			}
			newPlayBoard, err := playboard.PutStone(g.playBoard, newIndex, g.humanPlayer)
			if err != nil {
				switch e := err.(type) {
				case *playboard.PositionForbiddenError:
					log.Println(e)
					g.forbiddenMove = true
					//TO DO add position is busy?
				default:
					log.Println(e)
				}
				return nil
			}
			g.playBoard = newPlayBoard.Node
			g.index = newIndex
			playboard.PrintPlayBoard(g.playBoard)
			fmt.Println(g.machinePlayer)
			if g.machinePlayer.IndexAlmostWin != nil {
				fmt.Println(g.machinePlayer, *g.machinePlayer.IndexAlmostWin)
			}
			fmt.Println(g.humanPlayer)
			if g.humanPlayer.IndexAlmostWin != nil {
				fmt.Println(g.humanPlayer, *g.humanPlayer.IndexAlmostWin)
			}
			g.machineTurn = true
			g.turns += 1
			g.forbiddenMove = false
		}
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
			op.GeoM.Translate(float64(mx-widthStone/2), float64(my-widthStone/2))

			if string(stone) != playboard.SymbolPlayer1 {
				screen.DrawImage(WhiteStone, op)
			} else {
				screen.DrawImage(BlackStone, op)
			}
		}
	}

	text.Draw(screen, fmt.Sprintf("Turns: %d", g.GetTurns()), normalFont, 300, 80, color.Black)

}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *HumanGame) Draw(screen *ebiten.Image) {
	_draw(screen, g)
	text.Draw(screen, fmt.Sprintf("Turn to play for Player: %s", g.currentPlayer.Symbol), normalFont, 300, 60, color.Black)
	var y1, y2 int
	if g.currentPlayer.Symbol == playboard.SymbolPlayer1 {
		y1, y2 = 100, 120
	} else {
		y1, y2 = 120, 100
	}
	text.Draw(screen, fmt.Sprintf("Captures player %s: %d", g.currentPlayer.Symbol, g.currentPlayer.Captures), normalFont, 300, y1, color.Black)
	text.Draw(screen, fmt.Sprintf("Captures player %s: %d", g.anotherPlayer.Symbol, g.anotherPlayer.Captures), normalFont, 300, y2, color.Black)
	if g.forbiddenMove { //move to common _draw
		text.Draw(screen, fmt.Sprintf("Move is forbidden because of double free three"), normalFont, 30, 300, color.Black)
	}
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *AIGame) Draw(screen *ebiten.Image) {
	_draw(screen, g)
	startX := 300
	startY := 60
	text.Draw(screen, fmt.Sprintf("AI timer : %s", playboard.AITimer), normalFont, startX, startY, color.Black)
	startY += 20
	text.Draw(screen, fmt.Sprintf("Captures player %s: %d", g.machinePlayer.Symbol, g.machinePlayer.Captures), normalFont, 300, 100, color.Black)
	text.Draw(screen, fmt.Sprintf("Captures player %s: %d", g.humanPlayer.Symbol, g.humanPlayer.Captures), normalFont, 300, 120, color.Black)
	if g.forbiddenMove { //move to common _draw
		text.Draw(screen, fmt.Sprintf("Move is forbidden because of double free three"), normalFont, 30, 300, color.Black)
	}
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
		turns:         0,
	}
	return game
}

func NewAIGame() GameInterface {
	game := &AIGame{
		playBoard:     strings.Repeat(playboard.EmptySymbol, playboard.N*playboard.N),
		humanPlayer:   &playboard.Player1,
		machinePlayer: &playboard.MachinePlayer,
		index:         -1,
		isOver:        false,
		machineTurn:   true,
		turns:         0,
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
