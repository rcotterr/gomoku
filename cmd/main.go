package main

import (
	"flag"
	"fmt"
	"gomoku/pkg/play"
)

func main() {
	Human := flag.Bool("human", false, "Play against another human player")
	flag.Parse()

	fmt.Println("it's gomoku, let's play")

	if *Human {
		play.HumanPlay()
	} else {
		play.AIPlay()
	}

}

//Possible optimizations from aromny-w:
//make array of stone structure
//bits operations

//Bonuses
//- players' names
//- tests
//- choose colour for player instead of machine
//- some words in win
