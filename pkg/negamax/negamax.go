package NegaMax

import (
	"math"
)

type void struct{}
type children map[int]void

type player struct {
	captures int
	symbol   string
}

type state struct {
	board string
	transpositions map[string]void
	player1 player
	player2 player
}

type move struct {
	index int
	weight float64
}

func Heuristic() (float64) {
	return 0.0
}

func getMoves(state state, children children) []move {
	var moves []move

	for index, child := range children {
		if (state.transpositions[])
	}
}

func NegaMax(state state, depth int, multiplier int, children children, alpha float64, beta float64) (float64) {
	if depth == 0 {
		return float64(multiplier) * Heuristic();
	}

	max := math.Inf(-1);
	moves := getMoves(state, children)


	for index, child := range children {

		print(index, child);
	}
}
