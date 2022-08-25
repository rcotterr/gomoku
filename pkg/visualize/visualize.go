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
	"os"
	"strings"
)

const (
	screenWidth  = 1280
	screenHeight = 720

	cellWidth  = 12
	boardLines = 19
	boardStart = 30
	boardEnd   = boardStart + cellWidth*(boardLines-1)

	xInfo      = 300
	yStartInfo = 60
	diff       = 20
	yInfo1     = yStartInfo + diff
	yInfo2     = yStartInfo + 2*diff
	yInfo3     = yStartInfo + 3*diff
	yInfo4     = yStartInfo + 4*diff

	widthStone     = 14
	widthStoneHalf = widthStone / 2
)

type GameInterface interface {
	Update() error
	Draw(screen *ebiten.Image)
	Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
	GetPlayBoard() string
	GetPly() int
	GetForbiddenMove() bool
	GetIsOver() bool
}

type HumanGame struct {
	//screen *ebiten.Image
	playBoard     string
	currentPlayer *playboard.Player
	anotherPlayer *playboard.Player
	index         int
	isOver        bool
	ply           int
	forbiddenMove bool
	algo          playboard.Algo
	file          *os.File
}

func (g HumanGame) GetPlayBoard() string   { return g.playBoard }
func (g HumanGame) GetPly() int            { return g.ply }
func (g HumanGame) GetForbiddenMove() bool { return g.forbiddenMove }
func (g HumanGame) GetIsOver() bool        { return g.isOver }

type AIGame struct {
	//screen *ebiten.Image
	playBoard      string
	humanPlayer    *playboard.Player
	machinePlayer  *playboard.Player
	machineTurn    bool
	index          int
	isOver         bool
	ply            int
	forbiddenMove  bool
	humanMoveFirst bool
	algo           playboard.Algo
	file           *os.File
}

func (g AIGame) GetPlayBoard() string   { return g.playBoard }
func (g AIGame) GetPly() int            { return g.ply }
func (g AIGame) GetForbiddenMove() bool { return g.forbiddenMove }
func (g AIGame) GetIsOver() bool        { return g.isOver }

var _ = ebiten.NewImage(screenWidth, screenHeight)

var WhiteStone *ebiten.Image
var BlackStone *ebiten.Image

var normalFont font.Face

var colorScreen color.NRGBA
var colorGrid *color.RGBA

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

	colorScreen = color.NRGBA{R: 222, G: 184, B: 135, A: 255}
	colorGrid = &color.RGBA{A: 50}
}

func checkValidCoordinate(coord int) bool {
	return coord >= boardStart-widthStoneHalf && coord <= boardEnd+widthStoneHalf
}

func normalizeCoordinate(coord int) int {
	return int(math.Round(float64(coord-boardStart) / cellWidth))
}

func HumanTurnVis() (int, error) {
	mx, my := ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		if checkValidCoordinate(mx) && checkValidCoordinate(my) {
			nx, ny := normalizeCoordinate(mx), normalizeCoordinate(my)
			index := ny*playboard.N + nx
			if index >= 0 && index < playboard.N*playboard.N {
				return index, nil
			}
		}
	}
	return -1, fmt.Errorf("no step")

}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *HumanGame) Update() error {

	if !g.isOver && !playboard.GameOver(g.playBoard, g.currentPlayer, g.anotherPlayer, g.index) {
		newIndex, err := HumanTurnVis()
		if err != nil {
			return nil
		}
		newPlayBoard, err := playboard.PutStone(g.playBoard, newIndex, g.currentPlayer)
		if err != nil {
			fmt.Println(err) //TODO position forbidden
			return nil
		}
		g.playBoard = newPlayBoard.Node
		g.index = newIndex
		playboard.PrintPlayBoard(g.playBoard)
		g.currentPlayer, g.anotherPlayer = g.anotherPlayer, g.currentPlayer
		g.ply += 1
	} else {
		g.isOver = true
	}

	return nil
}

func FPrintCurrentState(g *AIGame) {
	playboard.FPrintPlayBoard(g.playBoard, g.file)
	fmt.Fprintln(g.file, g.machinePlayer)
	if g.machinePlayer.IndexAlmostWin != nil {
		fmt.Fprintln(g.file, g.machinePlayer, *g.machinePlayer.IndexAlmostWin)
	}
	fmt.Fprintln(g.file, g.humanPlayer)
	if g.humanPlayer.IndexAlmostWin != nil {
		fmt.Fprintln(g.file, g.humanPlayer, *g.humanPlayer.IndexAlmostWin)
	}
	if !g.machineTurn {
		fmt.Fprintln(g.file, "Algo took ", playboard.AITimer)
	}
	//position forbidden
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *AIGame) Update() error {
	if g.isOver {
		return nil
	}

	if !playboard.GameOver(g.playBoard, g.machinePlayer, g.humanPlayer, g.index) {
		if g.machineTurn {
			g.index = g.algo.GetIndex(g.playBoard, *g.machinePlayer, *g.humanPlayer)
			newPlayBoard, err := playboard.PutStone(g.playBoard, g.index, g.machinePlayer)
			if err != nil {
				fmt.Println("Invalid machine algo!!!!!", err)
				log.Fatal()
			}
			g.playBoard = newPlayBoard.Node
			g.machineTurn = false
			g.ply += 1
			g.forbiddenMove = false
			FPrintCurrentState(g)
		} else {
			newIndex, err := HumanTurnVis()
			if err != nil {
				return nil
			}
			newPlayBoard, err := playboard.PutStone(g.playBoard, newIndex, g.humanPlayer)
			if err != nil {
				switch err.(type) {
				case *playboard.PositionForbiddenError:
					g.forbiddenMove = true
				}
				fmt.Fprintln(g.file, err)
				return nil
			}
			g.playBoard = newPlayBoard.Node
			g.index = newIndex
			g.machineTurn = true
			g.ply += 1
			g.forbiddenMove = false
			FPrintCurrentState(g)
		}
	} else {
		g.isOver = true
	}

	return nil
}

func _drawGrid(screen *ebiten.Image) {
	xStart, yStart := boardStart, boardStart
	xEnd, yEnd := boardEnd, boardEnd

	for i := 0; i < boardLines; i++ {
		ebitenutil.DrawLine(screen, float64(xStart), float64(yStart), float64(xStart), float64(yEnd), colorGrid)
		xStart += cellWidth
	}
	xStart = boardStart
	for i := 0; i < boardLines; i++ {
		ebitenutil.DrawLine(screen, float64(xStart), float64(yStart), float64(xEnd), float64(yStart), colorGrid)
		yStart += cellWidth
	}
}

func _drawBoard(screen *ebiten.Image, board string, humanMovesFirst bool) {
	for index, stone := range board {
		if string(stone) != playboard.EmptySymbol {
			op := &ebiten.DrawImageOptions{}
			mx, my := index%playboard.N, index/playboard.N
			mx, my = boardStart+cellWidth*(mx), boardStart+cellWidth*(my)
			op.GeoM.Translate(float64(mx-widthStoneHalf), float64(my-widthStoneHalf))

			if string(stone) == playboard.SymbolPlayer1 || string(stone) == playboard.SymbolPlayerMachine && !humanMovesFirst {
				screen.DrawImage(WhiteStone, op)
			} else {
				screen.DrawImage(BlackStone, op)
			}
		}
	}
}

func _drawAdditionalText(screen *ebiten.Image, ply int, forbiddenMove, isOver bool) {
	text.Draw(screen, fmt.Sprintf("Turns: %d", ply/2), normalFont, xInfo, yInfo1, color.Black)
	text.Draw(screen, fmt.Sprintf("Ply: %d", ply), normalFont, xInfo, yInfo2, color.Black)

	if forbiddenMove {
		text.Draw(screen, fmt.Sprintf("Move is forbidden because of double free three"), normalFont, boardStart, 300, color.Black)
	}
	if isOver {
		text.Draw(screen, fmt.Sprintf("Game over!"), normalFont, boardStart, 300, color.Black)
	}
}

func _draw(screen *ebiten.Image, g GameInterface) {
	screen.Fill(colorScreen)

	_drawGrid(screen)

	var humanMovesFirst = false
	switch game := g.(type) {
	case *AIGame:
		humanMovesFirst = game.humanMoveFirst
	}
	_drawBoard(screen, g.GetPlayBoard(), humanMovesFirst)
	_drawAdditionalText(screen, g.GetPly(), g.GetForbiddenMove(), g.GetIsOver())
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *HumanGame) Draw(screen *ebiten.Image) {
	_draw(screen, g)
	text.Draw(screen, fmt.Sprintf("Turn to play for Player: %s", g.currentPlayer.Symbol), normalFont, xInfo, yStartInfo, color.Black)
	var y1, y2 int
	if g.currentPlayer.Symbol == playboard.SymbolPlayer1 {
		y1, y2 = yInfo3, yInfo4
	} else {
		y1, y2 = yInfo4, yInfo3
	}
	text.Draw(screen, fmt.Sprintf("Captures player %s: %d", g.currentPlayer.Symbol, g.currentPlayer.Captures), normalFont, xInfo, y1, color.Black)
	text.Draw(screen, fmt.Sprintf("Captures player %s: %d", g.anotherPlayer.Symbol, g.anotherPlayer.Captures), normalFont, xInfo, y2, color.Black)
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *AIGame) Draw(screen *ebiten.Image) {
	_draw(screen, g)
	text.Draw(screen, fmt.Sprintf("AI timer : %s", playboard.AITimer), normalFont, xInfo, yStartInfo, color.Black)
	text.Draw(screen, fmt.Sprintf("Captures player %s: %d", g.machinePlayer.Symbol, g.machinePlayer.Captures), normalFont, xInfo, yInfo3, color.Black)
	text.Draw(screen, fmt.Sprintf("Captures player %s: %d", g.humanPlayer.Symbol, g.humanPlayer.Captures), normalFont, xInfo, yInfo4, color.Black)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *HumanGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth / 2, outsideHeight / 2
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *AIGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth / 2, outsideHeight / 2
}

func NewHumanGame(file *os.File) GameInterface {
	game := &HumanGame{
		playBoard:     strings.Repeat(playboard.EmptySymbol, playboard.N*playboard.N),
		currentPlayer: &playboard.Player1,
		anotherPlayer: &playboard.Player2,
		index:         -1,
		isOver:        false,
		ply:           0,
		file:          file,
	}
	return game
}

func NewAIGame(depth int, humanMoveFirst bool, file *os.File) GameInterface {
	var humanPlayer playboard.Player
	if humanMoveFirst {
		humanPlayer = playboard.Player1
	} else {
		humanPlayer = playboard.Player2
	}
	game := &AIGame{
		playBoard:      strings.Repeat(playboard.EmptySymbol, playboard.N*playboard.N),
		humanPlayer:    &humanPlayer,
		machinePlayer:  &playboard.MachinePlayer,
		index:          -1,
		isOver:         false,
		machineTurn:    !humanMoveFirst,
		ply:            0,
		humanMoveFirst: humanMoveFirst,
		algo:           playboard.Algo{Depth: depth},
		file:           file,
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
