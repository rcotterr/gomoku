package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const n = 19

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

func parsePositions(text string) ([]int, error) {
	words := strings.Fields(text)
	if len(words) != 2 {
		return nil, fmt.Errorf("need 2 positions")
	}
	var pos []int
	for _, word := range words {
		num, err := strconv.Atoi(word)
		if err != nil {
			return nil, fmt.Errorf("invalid positions %s", err)
		}
		if num >= n || num < 0 {
			return nil, fmt.Errorf("invalid positions, can be from 0 to 18")
		}
		pos = append(pos, num)
	}
	return pos, nil
}

func putStone(playBoard string, pos []int) (string, error) {

	i := pos[0]
	j := pos[1]
	index := j * n + i
	fmt.Println(index)
	if playBoard[index] != '.' {
		return "", fmt.Errorf("position is busy")
	}

	newPlayBoard := strings.Join([]string{playBoard[:index], string("1"), playBoard[index+1:]}, "")

	return newPlayBoard, nil
}

func main() {
	fmt.Println("it's gomoku, let's play")
	var playBoard = strings.Repeat(".", n*n)
	printPlayBoard(playBoard)
	currentPlayer := false //TODO 0-1
	reader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Println("Player ", currentPlayer, ", enter positions (like 1 2):")
		text, _ := reader.ReadString('\n')
		pos, err := parsePositions(text)
		if err != nil {
			fmt.Println(err)
			continue
		}
		newPlayBoard, err := putStone(playBoard, pos)
		if err != nil {
			fmt.Println(err)
			continue
		}
		playBoard = newPlayBoard
		printPlayBoard(playBoard)
		currentPlayer = !currentPlayer
	}
}