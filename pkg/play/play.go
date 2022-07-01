package play

import (
	"bufio"
	"fmt"
	"gomoku/pkg/algo"
	"gomoku/pkg/playboard"
	"os"
	"strings"
)

func HumanTurn(reader *bufio.Reader, currentPlayer playboard.Player, playBoard string) (string, error) {
	fmt.Println("Player ", currentPlayer, ", enter positions (like 1 2):")
	text, _ := reader.ReadString('\n')
	pos, err := playboard.ParsePositions(text)
	if err != nil {
		return "", err
	}

	index := pos.Y*playboard.N + pos.X
	newPlayBoard, err := playboard.PutStone(playBoard, index, &currentPlayer)
	if err != nil {
		return "", err
	}
	return newPlayBoard, nil

}

func HumanPlay() {
	var playBoard = strings.Repeat(playboard.EmptySymbol, playboard.N*playboard.N)
	reader := bufio.NewReader(os.Stdin)
	playboard.PrintPlayBoard(playBoard)
	currentPlayer := playboard.Player1
	anotherPlayer := playboard.Player2
	for !playboard.IsOver(playBoard, &currentPlayer, &anotherPlayer) {
		newPlayBoard, err := HumanTurn(reader, currentPlayer, playBoard)
		if err != nil {
			fmt.Println(err)
			continue
		}
		playBoard = newPlayBoard
		playboard.PrintPlayBoard(playBoard)
		currentPlayer, anotherPlayer = anotherPlayer, currentPlayer
	}
}

func AIPlay() {
	var playBoard = strings.Repeat(playboard.EmptySymbol, playboard.N*playboard.N)
	reader := bufio.NewReader(os.Stdin)
	humanPlayer := playboard.Player1
	machinePlayer := playboard.MachinePlayer
	machineTurn := true
	var err error
	var newPlayBoard string
	for !playboard.IsOver(playBoard, &machinePlayer, &humanPlayer) {
		if machineTurn {
			machineIndex := algo.Algo(playBoard, machinePlayer, humanPlayer)
			playBoard, err = playboard.PutStone(playBoard, machineIndex, &machinePlayer)
			if err != nil {
				fmt.Println("Invalid machine algo!!!!!", err)
				return
			}
			playboard.PrintPlayBoard(playBoard)
			machineTurn = false
		} else {
			newPlayBoard, err = HumanTurn(reader, humanPlayer, playBoard)
			if err != nil {
				fmt.Println(err)
				continue
			}
			playBoard = newPlayBoard
			playboard.PrintPlayBoard(playBoard)
			machineTurn = true
		}
	}
}