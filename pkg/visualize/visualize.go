package visualize

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/openlyinc/pointy"
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
	screenWidth  = 800
	screenHeight = 472

	cellWidth  = 12
	boardLines = 19
	boardStart = 10
	boardEnd   = boardStart + cellWidth*(boardLines-1)

	xInfo      = 250
	yStartInfo = 30
	diff       = 20
	yInfo1     = yStartInfo + diff
	yInfo2     = yStartInfo + 2*diff
	yInfo3     = yStartInfo + 3*diff
	yInfo4     = yStartInfo + 4*diff
	yInfo6     = yStartInfo + 6*diff

	xInfoOver  = 15
	yInfoOver  = 100
	yInfoOver1 = yInfoOver + 30

	widthStone     = 14
	widthStoneHalf = widthStone / 2
)

type GameInterface interface {
	Update() error
	Draw(screen *ebiten.Image)
	Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
	GetPlayBoard() string
	GetIndex() int
	GetPly() int
	GetForbiddenMove() bool
	GetIsOver() bool
	GetWinnerPhrase() *string
}

type HumanGame struct {
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
func (g HumanGame) GetIndex() int          { return g.index }
func (g HumanGame) GetPly() int            { return g.ply }
func (g HumanGame) GetForbiddenMove() bool { return g.forbiddenMove }
func (g HumanGame) GetIsOver() bool        { return g.isOver }
func (g HumanGame) GetWinnerPhrase() *string {
	if g.isOver {
		if g.currentPlayer.Winner {
			return pointy.String(getPhraseByPlayer(g.currentPlayer))
		}
		if g.anotherPlayer.Winner {
			return pointy.String(getPhraseByPlayer(g.anotherPlayer))
		}
	}
	return nil
}

type AIGame struct {
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
func (g AIGame) GetIndex() int          { return g.index }
func (g AIGame) GetPly() int            { return g.ply }
func (g AIGame) GetForbiddenMove() bool { return g.forbiddenMove }
func (g AIGame) GetIsOver() bool        { return g.isOver }
func (g AIGame) GetWinnerPhrase() *string {
	if g.machinePlayer.Winner {
		return pointy.String("AI won")
	}
	if g.humanPlayer.Winner {
		return pointy.String("Human won")
	}
	return nil
}

func getPhraseByPlayer(player *playboard.Player) string {
	if player.Symbol == playboard.SymbolPlayer1 {
		return "First player won"
	}
	if player.Symbol == playboard.SymbolPlayer2 {
		return "Second player won"
	}
	return "Some player won"
}

var _ = ebiten.NewImage(screenWidth, screenHeight)

var WhiteStone *ebiten.Image
var BlackStone *ebiten.Image

var normalFont font.Face
var bigFont font.Face
var middleFont font.Face

var colorScreen color.NRGBA
var colorGrid *color.RGBA
var colorHalfBlack *color.NRGBA
var colorRed *color.RGBA
var colorGrey *color.RGBA

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

	bigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    40,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	middleFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    20,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	colorScreen = color.NRGBA{R: 222, G: 184, B: 135, A: 255}
	colorGrid = &color.RGBA{A: 50}
	colorHalfBlack = &color.NRGBA{R: 0, G: 0, B: 0, A: 128}
	colorRed = &color.RGBA{R: 255, G: 0, B: 0, A: 128}
	colorGrey = &color.RGBA{169, 169, 169, 255}
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

func FPrintPlayer(file *os.File, player *playboard.Player) {
	fmt.Fprintf(file, "* Player symbol %s\n\tcaptures %d\n", player.Symbol, player.Captures)
	if player.IndexAlmostWin != nil {
		fmt.Fprintf(file, "\tindex almost win %d\n", player.IndexAlmostWin)
	}
}

func FPrintCurrentState(g GameInterface) {
	var file *os.File
	var board string
	var player1 *playboard.Player
	var player2 *playboard.Player
	var machineTurn bool

	switch game := g.(type) {
	case *AIGame:
		file = game.file
		board = game.playBoard
		player1 = game.machinePlayer
		player2 = game.humanPlayer
	case *HumanGame:
		file = game.file
		board = game.playBoard
		player1 = game.currentPlayer
		player2 = game.anotherPlayer
	}
	playboard.FPrintPlayBoard(board, file)

	fmt.Fprintln(file, "Players info: ")
	FPrintPlayer(file, player1)
	FPrintPlayer(file, player2)

	if machineTurn {
		fmt.Fprintln(file, "Algo took ", playboard.AITimer)
	}
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
			switch err.(type) {
			case *playboard.PositionForbiddenError:
				g.forbiddenMove = true
			}
			fmt.Fprintln(g.file, err)
			return nil
		}
		g.playBoard = newPlayBoard.Node
		g.index = newIndex
		FPrintCurrentState(g)
		g.currentPlayer, g.anotherPlayer = g.anotherPlayer, g.currentPlayer
		g.ply += 1
	} else {
		g.isOver = true
	}

	return nil
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
			newPlayBoard, _ := playboard.PutStone(g.playBoard, g.index, g.machinePlayer)
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

func _drawCircle(screen *ebiten.Image, x, y int, clr color.Color) {
	radius64 := float64(widthStoneHalf / 4)
	minAngle := math.Acos(1 - 1/radius64)

	for angle := float64(0); angle <= 360; angle += minAngle {
		xDelta := radius64 * math.Cos(angle)
		yDelta := radius64 * math.Sin(angle)

		x1 := int(math.Round(float64(x) + xDelta))
		y1 := int(math.Round(float64(y) + yDelta))

		screen.Set(x1, y1, clr)
	}
}

func _drawBoard(screen *ebiten.Image, board string, humanMovesFirst bool, lastIndex int) {
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
	mx, my := lastIndex%playboard.N, lastIndex/playboard.N
	mx, my = boardStart+cellWidth*(mx), boardStart+cellWidth*(my)
	_drawCircle(screen, mx, my, colorGrey)
}

func _drawAdditionalText(screen *ebiten.Image, ply int, forbiddenMove, isOver bool, winnerPhrase *string) {
	text.Draw(screen, fmt.Sprintf("Turns: %d", ply/2), normalFont, xInfo, yInfo1, color.Black)
	text.Draw(screen, fmt.Sprintf("Ply: %d", ply), normalFont, xInfo, yInfo2, color.Black)

	if forbiddenMove {
		text.Draw(screen, fmt.Sprintf("Move is forbidden\n(double free three)"), normalFont, xInfo, yInfo6, colorRed)
	}
	if isOver {
		text.Draw(screen, fmt.Sprintf("Game over"), bigFont, xInfoOver, yInfoOver, colorHalfBlack)
		if winnerPhrase != nil {
			text.Draw(screen, fmt.Sprintf(*winnerPhrase), middleFont, xInfoOver, yInfoOver1, colorHalfBlack)
		}
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
	_drawBoard(screen, g.GetPlayBoard(), humanMovesFirst, g.GetIndex())
	_drawAdditionalText(screen, g.GetPly(), g.GetForbiddenMove(), g.GetIsOver(), g.GetWinnerPhrase())
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
