package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gomoku/pkg/playboard"
)

func main() {
	fmt.Println("it's gomoku, let's play")
	var playBoard = strings.Repeat(playboard.EmptySymbol, playboard.N*playboard.N)
	playboard.PrintPlayBoard(playBoard)
	currentPlayer := playboard.Player1
	anotherPlayer := playboard.Player2
	reader := bufio.NewReader(os.Stdin)
	for !playboard.IsOver(playBoard) {
		fmt.Println("Player ", currentPlayer, ", enter positions (like 1 2):")
		text, _ := reader.ReadString('\n')
		pos, err := playboard.ParsePositions(text)
		if err != nil {
			fmt.Println(err)
			continue
		}
		newPlayBoard, err := playboard.PutStone(playBoard, pos, currentPlayer)
		if err != nil {
			fmt.Println(err)
			continue
		}
		playBoard = newPlayBoard
		playboard.PrintPlayBoard(playBoard)
		currentPlayer, anotherPlayer = anotherPlayer, currentPlayer
	}
}

//Possible optimizations from aromny-w:
//make array of stone structure
//bits operations

//Bonuses
//- players' names
//- tests
//- choose colour for player instead of machine
//
