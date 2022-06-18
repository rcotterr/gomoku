package playboard

import (
	"fmt"
	"strconv"
	"strings"
)

const N = 19
const EmptySymbol = "."
const lenPositions = 2
const Player1 = "0"
const Player2 = "1"

type Pos struct {
	x int
	y int
}

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

func isCaptured(playBoard string, index int, currentPlayer string) (bool, *int, *int) {
	//thirdStoneCurrentPlayer := false

	// parse —
	j := index + 3
	if j/N == index/N && string(playBoard[j]) == currentPlayer { //TO DO check j is exactly not out of n*n
		//thirdStoneCurrentPlayer = true
		step := (j - index) / 2
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

func PutStone(playBoard string, pos *Pos, currentPlayer string) (string, error) {

	index := pos.y*N + pos.x
	fmt.Println(index)
	if string(playBoard[index]) != EmptySymbol {
		return "", fmt.Errorf("position is busy")
	}

	newPlayBoard := strings.Join([]string{playBoard[:index], currentPlayer, playBoard[index+1:]}, "")

	capture, index1, index2 := isCaptured(newPlayBoard, index, currentPlayer)
	if capture {
		fmt.Println("Capture: ", capture, *index1, *index2)
	}

	return newPlayBoard, nil
}

func checkFive(playBoard string, i int, symbol string) bool {
	count := 1 //TO DO add if i % n + 5 >= n
	j := i + 1
	for j%N != 0 && j < N*N {
		if string(playBoard[j]) == symbol {
			count += 1
		} else {
			break
		}
		j += 1
	}
	if count >= 5 {
		fmt.Println("Game is over, CONGRATULATIONS TO PLAYER ", symbol)
		return true
	}

	count = 1
	j = i + N
	for j/N != N && j < N*N {
		if string(playBoard[j]) == symbol {
			count += 1
		} else {
			break
		}
		j += N
	}
	if count >= 5 {
		fmt.Println("Game is over, CONGRATULATIONS TO PLAYER ", symbol)
		return true
	}

	count = 1
	j = i + N + 1
	for j < N*N {
		if string(playBoard[j]) == symbol {
			count += 1
		} else {
			break
		}
		j += N + 1
	}
	if count >= 5 {
		fmt.Println("Game is over, CONGRATULATIONS TO PLAYER ", symbol)
		return true
	}

	count = 1
	j = i + N - 1
	for j%N < i%N && j < N*N { //while column j < column i
		if string(playBoard[j]) == symbol {
			count += 1
		} else {
			break
		}
		j += N - 1
	}
	if count >= 5 {
		fmt.Println("Game is over, CONGRATULATIONS TO PLAYER ", symbol)
		return true
	}

	return false
	//TO DO add possibleCapture
	// 6 stones and capture only in 6
	// test not sequence
	// check panic: runtime error: index out of range
}

func IsOver(playBoard string) bool {
	for i, val := range playBoard {
		if string(val) == Player1 {
			if checkFive(playBoard, i, Player1) {
				return true
			}
			i += 1
		}
		if string(val) == Player2 {
			if checkFive(playBoard, i, Player2) {
				return true
			}
			i += 1
		}
		if string(val) == EmptySymbol {
			i += 1
		}
	}

	//TO DO add 10 captured
	//add no space left
	return false
}
