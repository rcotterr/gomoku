package main

import (
	"fmt"
	"strings"
)

const n = 19

func printPlayBoard(playBoard string) {
	fmt.Println("current play board:")
	for i, val := range playBoard {
		if  i % n == 0 {
			fmt.Println()
		} else {
			fmt.Print(string(val), "  ")
		}
	}
	fmt.Println()
}

func main() {
	fmt.Println("it's gomoku, let's play")
	var startArray = strings.Repeat(".", n*n)
	printPlayBoard(startArray)
}