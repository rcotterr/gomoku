package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const n = 19
const emptySymbol = "."
const lenPositions = 2
const player1 = "0"
const player2 = "1"

type Pos struct {
	x int
	y int
}

func printPlayBoard(playBoard string) {
	fmt.Println("current play board:")
	for i, val := range playBoard {
		if  i % n == 0 {
			fmt.Println()
		}
		fmt.Print(string(val), "  ")
	}
	fmt.Println()
}

func parsePositions(text string) (*Pos, error) {
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
		if num >= n || num < 0 {
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

func putStone(playBoard string, pos *Pos, currentPlayer string) (string, error) {

	index := pos.y * n + pos.x
	fmt.Println(index)
	if string(playBoard[index]) != emptySymbol {
		return "", fmt.Errorf("position is busy")
	}

	newPlayBoard := strings.Join([]string{playBoard[:index], currentPlayer, playBoard[index+1:]}, "")

	return newPlayBoard, nil
}

func checkFive(playBoard string, i int, symbol string) bool {
	count := 1 //TO DO add if i % n + 5 >= n
	j := i + 1
	for j % n != 0 {
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
	j = i + n
	for j / n != n {
		if string(playBoard[j]) == symbol {
			count += 1
		} else {
			break
		}
		j += n
	}
	if count >= 5 {
		fmt.Println("Game is over, CONGRATULATIONS TO PLAYER ", symbol)
		return true
	}

	count = 1
	j = i + n + 1
	for j < n*n {
		if string(playBoard[j]) == symbol {
			count += 1
		} else {
			break
		}
		j += n + 1
	}
	if count >= 5 {
		fmt.Println("Game is over, CONGRATULATIONS TO PLAYER ", symbol)
		return true
	}

	count = 1
	j = i + n - 1
	for j % n  < i % n { //while column j < column i
		if string(playBoard[j]) == symbol {
			count += 1
		} else {
			break
		}
		j += n - 1
	}
	if count >= 5 {
		fmt.Println("Game is over, CONGRATULATIONS TO PLAYER ", symbol)
		return true
	}

	return false
	//TO DO add possibleCapture
	// 6 stones and capture only in 6
	// test not sequence
}

func isOver(playBoard string) bool {
	for i, val := range playBoard {
		if string(val) == player1 {
			if checkFive(playBoard, i, player1) {
				return true
			}
			i += 1
		}
		if string(val) == player2 {
			if checkFive(playBoard, i, player2) {
				return true
			}
			i += 1
		}
		if string(val) == emptySymbol {
			i += 1
		}
	}

	//TO DO add 10 captured
	//add no space left
	return false
}

func main() {
	fmt.Println("it's gomoku, let's play")
	var playBoard = strings.Repeat(emptySymbol, n*n)
	printPlayBoard(playBoard)
	currentPlayer := player1
	anotherPlayer := player2
	reader := bufio.NewReader(os.Stdin)
	for !isOver(playBoard) {
		fmt.Println("Player ", currentPlayer, ", enter positions (like 1 2):")
		text, _ := reader.ReadString('\n')
		pos, err := parsePositions(text)
		if err != nil {
			fmt.Println(err)
			continue
		}
		newPlayBoard, err := putStone(playBoard, pos, currentPlayer)
		if err != nil {
			fmt.Println(err)
			continue
		}
		playBoard = newPlayBoard
		printPlayBoard(playBoard)
		currentPlayer, anotherPlayer = anotherPlayer, currentPlayer
	}
}


//Bonuses
//- players' names
//- tests
//
//