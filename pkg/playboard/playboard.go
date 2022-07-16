package playboard

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var File *os.File

var RunTimesHeuristic *int
var RunTimesIsOver *int
var RunTimesgetChildren *int
var RunTimesCopySet *int

var AllTimesHeuristic *time.Duration
var AllTimesIsOver *time.Duration
var AllTimesgetChildren *time.Duration
var AllTimesCopySet *time.Duration

func TimeTrack(start time.Time, name string, runTimes *int, allTime *time.Duration) {
	elapsed := time.Since(start)
	if runTimes != nil && allTime != nil {
		*runTimes += 1
		*allTime += elapsed
	}
	_, _ = fmt.Fprintf(File, "%s took %s\n", name, elapsed)
	if name == "alphaBeta depth {5}" {
		_, _ = fmt.Fprintf(File, "All took RunTimesHeuristic:%d, %s;\n RunTimesIsOver:%d, %s\n RunTimesgetChildren:%d, %s\n, CopySet: %d, %s\n",
			*RunTimesHeuristic, AllTimesHeuristic, *RunTimesIsOver, AllTimesIsOver, *RunTimesgetChildren, AllTimesgetChildren, *RunTimesCopySet, AllTimesCopySet)
	}
}

const N = 19
const EmptySymbol = "."
const lenPositions = 2
const SymbolPlayer1 = "0"
const SymbolPlayer2 = "1"
const SymbolPlayerMachine = "M"
const numOfCaptureStone = 2
const numOfCaptureStoneToWin = 10
const nextFromCapturedStone = numOfCaptureStone + 1
const numOfCheckFreeThree = 3

type Player struct {
	Captures int
	Symbol   string
}

var Player1 = Player{Captures: 0, Symbol: SymbolPlayer1}
var Player2 = Player{Captures: 0, Symbol: SymbolPlayer2}
var MachinePlayer = Player{Captures: 0, Symbol: SymbolPlayerMachine}

type Pos struct {
	X int
	Y int
}

type ConditionFn func(int, int) bool

func PrintPlayBoard(playBoard string) {
	fmt.Println("current play board:")

	fmt.Print("   0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18")
	for i, val := range playBoard {
		if i%N == 0 {
			fmt.Println()
			if i/N > 9 {
				fmt.Print(i/N, " ")
			} else {
				fmt.Print(i/N, "  ")
			}
		}
		fmt.Print(string(val), "  ")
	}
	fmt.Println()
}

func ParsePositions(text string) (*Pos, error) {
	words := strings.Fields(text)
	if len(words) != lenPositions {
		return nil, fmt.Errorf("need %d positions", lenPositions)
	}
	var pos = Pos{}
	for i, word := range words {
		num, err := strconv.Atoi(word)
		if err != nil {
			return nil, fmt.Errorf("invalid positions %s", err)
		}
		if num >= N || num < 0 {
			return nil, fmt.Errorf("invalid positions, can be from 0 to 18")
		}
		if i == 0 {
			pos.X = num
		} else if i == 1 {
			pos.Y = num
		}
	}
	return &pos, nil
}

func ConditionHorizontal(j int, i int) bool {
	return j/N == i/N //if the same string
}

func ConditionBackDiagonal(j int, i int) bool { // diagonal is \
	if i > j {
		return j%N <= i%N
	} else {
		return j%N >= i%N
	}
}

func ConditionForwardDiagonal(j int, i int) bool { // diagonal is /
	if i > j {
		return j%N >= i%N
	} else {
		return j%N <= i%N
	}
}

func checkCapturedByCondition(step int, condition ConditionFn, playBoard string, index int, currentPlayer string) (bool, *int, *int) {
	j := index + nextFromCapturedStone*step

	if condition(j, index) && j >= 0 && j < N*N && string(playBoard[j]) == currentPlayer {
		index1 := index + step
		index2 := index + step*2
		symbol1 := string(playBoard[index1])
		symbol2 := string(playBoard[index2])

		if symbol1 != currentPlayer && symbol1 != EmptySymbol && symbol2 != currentPlayer && symbol2 != EmptySymbol { //TO DO check another player
			return true, &index1, &index2
		}
		//fmt.Println(index1, index2)
	}
	return false, nil, nil
}

func isCaptured(playBoard string, index int, currentPlayer string) (bool, *int, *int) {
	setRules := map[int]ConditionFn{
		1:      ConditionHorizontal,
		N:      ConditionVertical,
		N + 1:  ConditionBackDiagonal,
		N - 1:  ConditionForwardDiagonal,
		-1:     ConditionHorizontal,
		-N:     ConditionVertical,
		-N - 1: ConditionBackDiagonal,
		-N + 1: ConditionForwardDiagonal,
	}

	for step, condition := range setRules {
		if isCapture, index1, index2 := checkCapturedByCondition(step, condition, playBoard, index, currentPlayer); isCapture {
			return isCapture, index1, index2
		}
	}

	return false, nil, nil
}

func isFreeThree(step int, condition ConditionFn, playBoard string, index int, currentPlayer string) bool {
	canBeEmpty := true
	startIndex := index + -1*numOfCheckFreeThree*step
	endIndex := index + numOfCheckFreeThree*step
	countStones := 0
	for startIndex < 0 || !condition(startIndex, index) { // TO DO for diagonal
		startIndex += step
	}
	for endIndex >= N*N || !condition(endIndex, index) {
		endIndex -= step
	}

	for j := startIndex; condition(j, index) && j >= 0 && j < N*N; j += step {
		if string(playBoard[j]) == currentPlayer {
			countStones += 1
			if countStones == 1 {
				startIndex = j
			} else if countStones == numOfCheckFreeThree {
				endIndex = j
				break
			}
		} else if string(playBoard[j]) == EmptySymbol {
			if countStones != 0 {
				if canBeEmpty {
					canBeEmpty = false
				} else {
					countStones = 0
					canBeEmpty = true
				}
			}
		} else { // another player
			countStones = 0
			canBeEmpty = true
		}
	}
	if countStones >= numOfCheckFreeThree {
		for _, j := range []int{startIndex - step, endIndex + step} { // TO DO different for diagonal
			if !(condition(j, index) && j >= 0 && j < N*N && string(playBoard[j]) == EmptySymbol) {
				return false
			}
		}
		return true
	}

	return false
}

func isForbidden(playBoard string, index int, currentPlayer string) bool {
	setRules := map[int]ConditionFn{
		1:     ConditionHorizontal,
		N:     ConditionVertical,
		N + 1: ConditionBackDiagonal,
		N - 1: ConditionForwardDiagonal,
	}
	countFreeThree := 0

	for step, condition := range setRules {
		if isFreeThreeRow := isFreeThree(step, condition, playBoard, index, currentPlayer); isFreeThreeRow {
			countFreeThree += 1
		}
	}

	if countFreeThree >= 2 {
		return true
	}
	return false
}

func PutStone(playBoard string, index int, currentPlayer *Player) (string, error) {

	//index := pos.Y*N + pos.X
	//fmt.Println(index)
	if string(playBoard[index]) != EmptySymbol {
		return "", fmt.Errorf("position is busy")
	}

	newPlayBoard := strings.Join([]string{playBoard[:index], currentPlayer.Symbol, playBoard[index+1:]}, "")

	capture, index1, index2 := isCaptured(newPlayBoard, index, currentPlayer.Symbol)
	if capture {
		//fmt.Println("Capture: ", capture, *index1, *index2)
		if *index1 > *index2 {
			index1, index2 = index2, index1
		}
		newPlayBoard = strings.Join([]string{newPlayBoard[:*index1], EmptySymbol, newPlayBoard[*index1+1 : *index2], EmptySymbol, newPlayBoard[*index2+1:]}, "")
		currentPlayer.Captures += 1
	} else if isForbidden(newPlayBoard, index, currentPlayer.Symbol) {
		return "", fmt.Errorf("position is forbidden")
	}

	return newPlayBoard, nil
}

func FiveInRow(i int, step int, condition ConditionFn, playBoard string, symbol string) bool {
	count := 1 //TO DO add if i % n + 5 >= n
	j := i + step
	for condition(j, i) && j >= 0 && j < N*N {
		if string(playBoard[j]) == symbol {
			count += 1
		} else {
			break
		}
		j += step
	}
	if count >= 5 {
		//fmt.Println("Game is over, CONGRATULATIONS TO PLAYER ", symbol)
		return true
	}
	return false
}

func ConditionVertical(j int, _ int) bool {
	return j >= 0 && j/N < N //till last + 1 raw
}

func checkFive(playBoard string, i int, symbol string) bool {

	setRules := map[int]ConditionFn{
		1:     ConditionHorizontal,
		N:     ConditionVertical,
		N + 1: ConditionBackDiagonal,
		N - 1: ConditionForwardDiagonal,
	}

	for step, condition := range setRules {
		if FiveInRow(i, step, condition, playBoard, symbol) {
			return true
		}
	}

	return false
	//TO DO add possibleCapture than not win
	// 6 stones and capture only in 6
	//TO DO check only from new stone
}

func IsOver(playBoard string, player1 *Player, player2 *Player) bool { //TO DO change func without print
	defer TimeTrack(time.Now(), "IsOver", RunTimesIsOver, AllTimesIsOver)

	for _, player := range []*Player{player1, player2} {
		if player != nil && player.Captures >= numOfCaptureStoneToWin/numOfCaptureStone {
			//fmt.Println("Game is over, CONGRATULATIONS TO PLAYER ", player.Symbol)
			return true
		}
	}

	for j, val := range playBoard { // TO DO not all check but only 1 last put stone  && not range but while i < len
		value := string(val)
		if value != EmptySymbol {
			if checkFive(playBoard, j, value) {
				return true
			}
		}
	}

	if containEmpty := strings.Contains(playBoard, EmptySymbol); !containEmpty {
		//fmt.Println("Game is over, no space left, both players win")
		return true
	}
	return false
}
