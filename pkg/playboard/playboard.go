package playboard

import (
	"fmt"
	"strconv"
	"strings"
)

const N = 19
const EmptySymbol = "."
const lenPositions = 2
const SymbolPlayer1 = "0"
const SymbolPlayer2 = "1"
const numOfCaptureStone = 2
const numOfCaptureStoneToWin = 10
const nextFromCapturedStone = numOfCaptureStone + 1
const numOfCheckFreeThree = 3

type Player struct {
	captures int
	symbol   string
}

var Player1 = Player{captures: 0, symbol: SymbolPlayer1}
var Player2 = Player{captures: 0, symbol: SymbolPlayer2}

type Pos struct {
	x int
	y int
}

type conditionFn func(int, int) bool

func PrintPlayBoard(playBoard string) {
	fmt.Println("current play board:")
	for i, val := range playBoard {
		if i%N == 0 {
			fmt.Println()
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
			pos.x = num
		} else if i == 1 {
			pos.y = num
		}
	}
	return &pos, nil
}

func conditionHorizontalCapture(j int, i int) bool {
	return j >= 0 && j/N == i/N //if the same string
}

func conditionRightDiagonalCapture(j int, _ int) bool { // diagonal is \
	return j >= 0 && j < N*N //till not out of index
}

func checkCapturedByCondition(step int, condition conditionFn, playBoard string, index int, currentPlayer string) (bool, *int, *int) {
	j := index + nextFromCapturedStone*step

	if condition(j, index) && j >= 0 && j < N*N && string(playBoard[j]) == currentPlayer {
		index1 := index + step
		index2 := index + step*2
		symbol1 := string(playBoard[index1])
		symbol2 := string(playBoard[index2])

		if symbol1 != currentPlayer && symbol1 != EmptySymbol && symbol2 != currentPlayer && symbol2 != EmptySymbol { //TO DO check another player
			return true, &index1, &index2
		}
		fmt.Println(index1, index2)
	}
	return false, nil, nil
}

func isCaptured(playBoard string, index int, currentPlayer string) (bool, *int, *int) {
	setRules := map[int]conditionFn{
		1:      conditionHorizontalCapture,
		N:      conditionVertical,
		N + 1:  conditionRightDiagonalCapture, //TO DO delete duplicate conditionRightDiagonal
		N - 1:  conditionLeftDiagonal,
		-1:     conditionHorizontalCapture,
		-N:     conditionVertical,
		-N - 1: conditionRightDiagonalCapture,
		-N + 1: conditionLeftUpperDiagonal,
	}

	for step, condition := range setRules {
		if isCapture, index1, index2 := checkCapturedByCondition(step, condition, playBoard, index, currentPlayer); isCapture {
			return isCapture, index1, index2
		}
	}

	return false, nil, nil
}

func conditionLeftDiagonalCheckFreeThree(j int, i int) bool {
	columnDiff := j%N - i%N
	return columnDiff >= 0 && columnDiff <= numOfCheckFreeThree+1 || columnDiff <= 0 && columnDiff >= -(numOfCheckFreeThree+1)
	// upper and new column diff less than/equal 3 or lower and new column diff more than/equal -3
	//+1 is for empty
}

func isFreeThree(step int, condition conditionFn, playBoard string, index int, currentPlayer string) bool {
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
	setRules := map[int]conditionFn{
		1:     conditionHorizontalCapture,
		N:     conditionVertical,
		N + 1: conditionRightDiagonalCapture, //TO DO delete duplicate conditionRightDiagonal
		N - 1: conditionLeftDiagonalCheckFreeThree,
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

func PutStone(playBoard string, pos *Pos, currentPlayer *Player) (string, error) {

	index := pos.y*N + pos.x
	fmt.Println(index)
	if string(playBoard[index]) != EmptySymbol {
		return "", fmt.Errorf("position is busy")
	}

	newPlayBoard := strings.Join([]string{playBoard[:index], currentPlayer.symbol, playBoard[index+1:]}, "")

	capture, index1, index2 := isCaptured(newPlayBoard, index, currentPlayer.symbol)
	if capture {
		fmt.Println("Capture: ", capture, *index1, *index2)
		if *index1 > *index2 {
			index1, index2 = index2, index1
		}
		newPlayBoard = strings.Join([]string{newPlayBoard[:*index1], EmptySymbol, newPlayBoard[*index1+1 : *index2], EmptySymbol, newPlayBoard[*index2+1:]}, "")
		currentPlayer.captures += 1
	} else if isForbidden(newPlayBoard, index, currentPlayer.symbol) {
		return "", fmt.Errorf("position is forbidden")
	}

	return newPlayBoard, nil
}

func FiveInRow(i int, step int, condition conditionFn, playBoard string, symbol string) bool {
	count := 1 //TO DO add if i % n + 5 >= n
	j := i + step
	for condition(j, i) && j < N*N {
		if string(playBoard[j]) == symbol {
			count += 1
		} else {
			break
		}
		j += step
	}
	if count >= 5 {
		fmt.Println("Game is over, CONGRATULATIONS TO PLAYER ", symbol)
		return true
	}
	return false
}

func conditionHorizontal(j int, _ int) bool {
	return j%N != 0 //till next row
}

func conditionVertical(j int, _ int) bool {
	return j >= 0 && j/N < N //till last + 1 raw
}

func conditionRightDiagonal(j int, _ int) bool {
	return j < N*N //till not out of index
}

func conditionLeftDiagonal(j int, i int) bool {
	return j%N < i%N //till column of left put stones less than start stone index
}

func conditionLeftUpperDiagonal(j int, i int) bool {
	return j%N > i%N //till column of right(left upper) put stones more than start stone index
}

func checkFive(playBoard string, i int, symbol string) bool {

	setRules := map[int]conditionFn{
		1:     conditionHorizontal,
		N:     conditionVertical,
		N + 1: conditionRightDiagonal,
		N - 1: conditionLeftDiagonal,
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

func IsOver(playBoard string, player1 *Player, player2 *Player) bool {
	for _, player := range []*Player{player1, player2} {
		if player != nil && player.captures >= numOfCaptureStoneToWin/numOfCaptureStone {
			fmt.Println("Game is over, CONGRATULATIONS TO PLAYER ", player.symbol)
			return true
		}
	}

	for i, val := range playBoard { // TO DO not all check but only 1 last put stone  && not range but while i < len
		value := string(val)
		if value == SymbolPlayer1 || value == SymbolPlayer2 {
			if checkFive(playBoard, i, value) {
				return true
			}
			i += 1
		}
		if string(val) == EmptySymbol {
			i += 1
		}
	}

	if containEmpty := strings.Contains(playBoard, EmptySymbol); !containEmpty {
		fmt.Println("Game is over, no space left, both players win")
		return true
	}
	return false
}
