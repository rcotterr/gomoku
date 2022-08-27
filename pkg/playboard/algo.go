package playboard

import (
	"math"
	"sort"
	"strings"
	"time"
)

type void struct{}
type intSet map[int]void

var member void

var setRulesChildren = map[int]ConditionFn{
	1:      ConditionHorizontal,
	N:      ConditionVertical,
	N + 1:  ConditionBackDiagonal,
	N - 1:  ConditionForwardDiagonal,
	-1:     ConditionHorizontal,
	-N:     ConditionVertical,
	-N - 1: ConditionBackDiagonal,
	-N + 1: ConditionForwardDiagonal,
}

func Heuristic(move Move, symbol string, captures int) float64 {
	vulnerable := false
	num := 0.0

	if symbol == "" {
		symbol = string(move.Node[move.index])
	}

	setRules := map[int]ConditionFn{
		1:     ConditionHorizontal,
		N:     ConditionVertical,
		N + 1: ConditionBackDiagonal,
		N - 1: ConditionForwardDiagonal,
	}

	for step, condition := range setRules {
		count, halfFree, free, _ := CountInRow(move.Node, move.index, step, condition, symbol)
		if count >= 5 || captures >= 5 {
			return 1000000000000
		} else if count == 4 && free {
			num += 10000000000
		} else if count == 4 && halfFree {
			num += 1000000000
		} else if count == 3 && free {
			num += 10000000
		} else if count == 3 && halfFree {
			num += 1000000
		} else if count == 2 && free {
			num += 10000
		} else if count == 2 && halfFree {
			num += 1000
			vulnerable = true
		} else if count == 1 && free {
			num += 10
		} else if count == 1 && halfFree {
			num += 1
		}
	}

	num += float64(1000000000 * move.Captures)
	if vulnerable == true {
		num -= 1000000000
	}

	if num < 0 {
		return 0
	}

	return num
}

func getHeuristic(state State, player Player, opponent Player) (float64, float64) {
	var h1, h2 float64

	if player.Winner {
		return 1000000000000, 0
	}
	if opponent.Winner {
		return 0, 1000000000000
	}
	h1 = Heuristic(state.move, player.Symbol, player.Captures)
	state.move.Node = strings.Join([]string{state.move.Node[:state.move.index], EmptySymbol, state.move.Node[state.move.index+1:]}, "")
	move, err := PutStone(state.move.Node, state.move.index, &opponent)
	if err == nil {
		h2 = Heuristic(move, opponent.Symbol, opponent.Captures)
	} else {
		h2 = 0
	}

	if h1 != 1000000000000 {
		h1 *= 10
	}

	return h1, h2
}

func UpdateSetChildren(index int, playBoard string, set intSet) {
	for step, condition := range setRulesChildren {
		j := index + step
		if j >= 0 && j < N*N && condition(j, index) && string(playBoard[j]) == EmptySymbol {
			set[j] = member
		}
	}
}

func getChildren(state State, currentPlayer Player, childIndexesSet intSet) []State {
	var children []State

	for k := range childIndexesSet {
		updatePlayer := currentPlayer

		infoChild, err := PutStone(state.move.Node, k, &updatePlayer)

		if err == nil {
			var newStateChild State
			if currentPlayer.Symbol == SymbolPlayerMachine {
				newStateChild = State{
					move:          infoChild,
					machinePlayer: updatePlayer,
					humanPlayer:   state.humanPlayer,
				}
			} else {
				newStateChild = State{
					move:          infoChild,
					machinePlayer: state.machinePlayer,
					humanPlayer:   updatePlayer,
				}
			}

			children = append(children, newStateChild)
		}
	}

	return children
}

func getAllIndexChildren(playBoard string) intSet {
	set := make(intSet)

	for index, val := range playBoard {
		value := string(val)
		if value != EmptySymbol {
			UpdateSetChildren(index, playBoard, set)
		}
	}

	return set
}

func copySet(children intSet) intSet {
	setNewChildIndexes := make(intSet)

	for index, val := range children {
		setNewChildIndexes[index] = val
	}

	return setNewChildIndexes
}

type Child struct {
	State State
	Value float64
	h1    float64
	h2    float64
}

func sortChildren(children []State, multiplier int) []Child {
	var new_ []Child

	for _, childState := range children {
		var h1, h2 float64

		if multiplier == 1 {
			h1, h2 = getHeuristic(childState, childState.machinePlayer, childState.humanPlayer)
		} else {
			h1, h2 = getHeuristic(childState, childState.humanPlayer, childState.machinePlayer)
		}

		new_ = append(new_, Child{childState, math.Max(h1, h2), h1, h2})
	}

	sort.Slice(new_, func(i, j int) bool {
		if new_[i].Value == new_[j].Value {
			return new_[i].h1 > new_[j].h1
		}

		return new_[i].Value > new_[j].Value
	})

	return new_
}

type State struct {
	move          Move
	machinePlayer Player
	humanPlayer   Player
}

func (a Algo) NegaScout(state State, depth int, alpha float64, beta float64, multiplier int, childIndexesSet intSet) (float64, int) {
	if depth == 0 || GameOver(state.move.Node, &state.machinePlayer, &state.humanPlayer, state.move.index) {
		h1, h2 := getHeuristic(state, state.machinePlayer, state.humanPlayer)

		if h1 == 1000000000000 {
			return float64(multiplier) * (h1 + (float64(depth) * 0.1)), state.move.index
		} else if h2 == 1000000000000 {
			return float64(multiplier) * (-h2 - (float64(depth) * 0.1)), state.move.index
		}

		return float64(multiplier) * (math.Max(h1, h2)), state.move.index
	}

	maxEval := math.Inf(-1)
	maxIndex := -1

	var children []State
	var childrenSlice []Child

	if state.move.index != -1 {
		UpdateSetChildren(state.move.index, state.move.Node, childIndexesSet)
	}

	if multiplier == 1 {
		children = getChildren(state, state.machinePlayer, childIndexesSet)
	} else {
		children = getChildren(state, state.humanPlayer, childIndexesSet)
	}

	childrenSlice = sortChildren(children, multiplier)

	for i, child := range childrenSlice {
		setNewChildIndexes := copySet(childIndexesSet)
		eval, _ := a.NegaScout(child.State, depth-1, -beta, -alpha, -multiplier, setNewChildIndexes)
		eval = -eval

		if eval > maxEval {
			maxEval = eval

			if depth == a.Depth {
				maxIndex = child.State.move.index
			}
		}

		alpha = math.Max(alpha, eval)

		if alpha >= beta {
			break
		}

		if depth != a.Depth {
			if i+1 < len(childrenSlice) {
				if child.Value > childrenSlice[i+1].Value || child.h1 > childrenSlice[i+1].h1 {
					break
				}
			}
		} else if maxEval > 0 {
			if i+1 < len(childrenSlice) {
				if child.Value > childrenSlice[i+1].Value || child.h1 > childrenSlice[i+1].h1 {
					break
				}
			}
		}
	}

	return maxEval, maxIndex
}

type Algo struct {
	Depth int
}

func (a Algo) GetIndex(playBoard string, machinePlayer Player, humanPlayer Player) int {
	defer TimeTrackPrint(time.Now())

	setChildren := getAllIndexChildren(playBoard)

	if len(setChildren) == 0 {
		return 9*19 + 9
	}

	_, index := a.NegaScout(State{Move{playBoard, -1, 0, []int{}}, machinePlayer, humanPlayer}, a.Depth, math.Inf(-1), math.Inf(1), 1, setChildren)

	return index
}
