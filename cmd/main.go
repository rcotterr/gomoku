package main

import (
	"flag"
	"gomoku/pkg/visualize"
	"log"
	"os"
)

func main() {
	Human := flag.Bool("human", false, "Play against another human player")
	moveFirst := flag.Bool("moveFirst", false, "Human player moves first")
	depth := flag.Int("depth", 10, "Depth for algorithm")
	fileName := flag.String("game-history", "game-history", "File to write game history")
	flag.Parse()

	file, fileErr := os.Create(*fileName)
	if fileErr != nil {
		log.Fatal(fileErr)
	}

	var game visualize.GameInterface
	if *Human {
		game = visualize.NewHumanGame(file)
	} else {
		game = visualize.NewAIGame(*depth, *moveFirst, file)
	}
	visualize.Vis(game)
}

//Possible optimizations from aromny-w:
//make array of stone structure
//bits operations
//make map like in sapper
//return enum from free halfFree

//Bonuses
//- players' names
//- tests
//- choose colour for player instead of machine
//- some words in win
