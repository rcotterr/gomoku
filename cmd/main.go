package main

import (
	"flag"
	"gomoku/pkg/visualize"
)

func main() {
	Human := flag.Bool("human", false, "Play against another human player")
	moveFirst := flag.Bool("moveFirst", false, "Human player moves first")
	depth := flag.Int("depth", 10, "Depth for algorithm")
	flag.Parse()

	var game visualize.GameInterface
	if *Human {
		game = visualize.NewHumanGame()
	} else {
		game = visualize.NewAIGame(*depth, *moveFirst)
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
