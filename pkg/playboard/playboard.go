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

func ConditionHorizontal(j int, i int) bool {
	return j/N == i/N //if the same string
}

func ConditionVertical(_ int, _ int) bool {
	return true
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

func isCaptured(playBoard string, index int, currentPlayer string) (int, []int) {
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

	var arrIndexes []int
	numCaptures := 0

	for step, condition := range setRules {
		if isCapture, index1, index2 := checkCapturedByCondition(step, condition, playBoard, index, currentPlayer); isCapture {
			numCaptures += 1
			arrIndexes = append(arrIndexes, *index1, *index2)
		}
	}

	return numCaptures, arrIndexes
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
			if countFreeThree == 2 {
				return true
			}
		}
	}

	return false
}

func PutStone(playBoard string, index int, currentPlayer *Player) (State, error) {

	//index := pos.Y*N + pos.X
	//fmt.Println(index)
	if string(playBoard[index]) != EmptySymbol {
		return State{}, fmt.Errorf("position is busy")
	}

	newPlayBoard := strings.Join([]string{playBoard[:index], currentPlayer.Symbol, playBoard[index+1:]}, "")

	numCaptures, arrIndexes := isCaptured(newPlayBoard, index, currentPlayer.Symbol) //TO DO more than one capture
	if numCaptures > 0 {
		for _, capturedIndex := range arrIndexes {
			newPlayBoard = strings.Join([]string{newPlayBoard[:capturedIndex], EmptySymbol, newPlayBoard[capturedIndex+1:]}, "")
		}
		currentPlayer.Captures += numCaptures
	} else if isForbidden(newPlayBoard, index, currentPlayer.Symbol) {
		return State{}, fmt.Errorf("position is forbidden")
	}

	return State{newPlayBoard, index, numCaptures}, nil //TO DO return not always 0 captures
}

func CountInRow(node string, index int, step int, condition ConditionFn, symbol string) (int, bool, bool) {

	empty := 0

	hreny := node[index:]
	hreny = string(hreny)
	startIndex := index
	for tmpIndex := index - step; tmpIndex >= 0 && tmpIndex > index-(step*5); tmpIndex -= step {
		if condition(tmpIndex, index) && string(node[tmpIndex]) == symbol {
			startIndex = tmpIndex
		} else {
			if string(node[tmpIndex]) == EmptySymbol {
				empty += 1
			}
			break
		}
	}

	endIndex := index
	for tmpIndex := index + step; tmpIndex < N*N && tmpIndex < index+(step*5); tmpIndex += step {
		if condition(tmpIndex, index) && string(node[tmpIndex]) == symbol {
			endIndex = tmpIndex
		} else {
			if string(node[tmpIndex]) == EmptySymbol {
				empty += 1
			}
			break
		}
	}

	count := ((endIndex - startIndex) / step) + 1

	return count, empty == 1, empty == 2

}

func checkFive(playBoard string, index int, symbol string) bool {

	setRules := map[int]ConditionFn{
		1:     ConditionHorizontal,
		N:     ConditionVertical,
		N + 1: ConditionBackDiagonal,
		N - 1: ConditionForwardDiagonal,
	}

	for step, condition := range setRules {
		count, _, _ := CountInRow(playBoard, index, step, condition, symbol)
		if count >= 5 { // TO DO and not capture
			return true
		}
	}

	return false
	//TO DO add possibleCapture than not win
	// 6 stones and capture only in 6
}

func GameOver(playBoard string, player1 *Player, player2 *Player, index int) bool { //TO DO change func without print
	defer TimeTrack(time.Now(), "GameOver", RunTimesIsOver, AllTimesIsOver)

	if index == -1 {
		return false //first launch
	}
	for _, player := range []*Player{player1, player2} {
		if player != nil && player.Captures >= numOfCaptureStoneToWin/numOfCaptureStone {
			//fmt.Println("Game is over, CONGRATULATIONS TO PLAYER ", player.Symbol)
			return true
		}
	}
	//if string(playBoard[index]) == EmptySymbol {
	//	return false
	//}

	if checkFive(playBoard, index, string(playBoard[index])) {
		return true
	}

	if containEmpty := strings.Contains(playBoard, EmptySymbol); !containEmpty {
		//fmt.Println("Game is over, no space left, both players win")
		return true
	}
	return false
}
