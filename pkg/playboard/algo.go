package playboard

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

type void struct{}
type intSet map[int]void
type stringSet map[string]void

var member void
var t = 0

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

func Heuristic(state State, symbol string, captures int, capturedIndexes []int) float64 {
	defer TimeTrack(time.Now(), "Heuristic", RunTimesHeuristic, AllTimesHeuristic)
	vulnerable := false
	num := 0.0

	//count по каждой стороне
	//symbol := string(node[index])
	if symbol == "" {
		symbol = string(state.Node[state.index])
	}

	setRules := map[int]ConditionFn{
		1:     ConditionHorizontal,
		N:     ConditionVertical,
		N + 1: ConditionBackDiagonal,
		N - 1: ConditionForwardDiagonal,
	}

	for step, condition := range setRules {
		count, halfFree, free, _ := CountInRow(state.Node, state.index, step, condition, symbol)
		if count >= 5 || captures >= 5 {
			return 1000000000000
		} else if count == 4 && free {
			num += 10000000000
		} else if count == 4 && halfFree {
			num += 1000000000
		} else if count == 3 && free { //TO DO free more than one step
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

	//for _, capturedIndex := range capturedIndexes {
	//num += Heuristic(State{state.Node, capturedIndex, 0, int[]{}}, 0,)
	//}

	num += float64(1000000000 * state.Captures)
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
	h1 = Heuristic(state, player.Symbol, player.Captures, state.capturedIndexes)
	state.Node = strings.Join([]string{state.Node[:state.index], EmptySymbol, state.Node[state.index+1:]}, "")
	info, err := PutStone(state.Node, state.index, &opponent)
	newState := State{
		Node:            info.Node,
		index:           info.index,
		Captures:        info.Captures,
		capturedIndexes: info.capturedIndexes,
		machinePlayer:   player,
		humanPlayer:     opponent,
	}
	if err == nil {
		h2 = Heuristic(newState, opponent.Symbol, opponent.Captures, state.capturedIndexes)
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
	defer TimeTrack(time.Now(), "getChildren", RunTimesgetChildren, AllTimesgetChildren)
	var children []State

	//if index != -1 {
	//	UpdateSetChildren(index, node, childIndexesSet)
	//}

	//var updatePlayer *Player
	//
	//if currentPlayer.Symbol == SymbolPlayerMachine{
	//	updatePlayer = &state.machinePlayer
	//} else {
	//	updatePlayer = &state.humanPlayer
	//}

	for k := range childIndexesSet {
		//var updatePlayer *Player
		//
		//if currentPlayer.Symbol == SymbolPlayerMachine{
		//	updatePlayer = &state.machinePlayer
		//} else {
		//	updatePlayer = &state.humanPlayer
		//}
		var updatePlayer Player
		if currentPlayer.Symbol == SymbolPlayerMachine {
			updatePlayer = state.machinePlayer
		} else {
			updatePlayer = state.humanPlayer
		}

		infoChild, err := PutStone(state.Node, k, &updatePlayer)

		var machinePlayer Player
		var humanPlayer Player

		if currentPlayer.Symbol == SymbolPlayerMachine {
			machinePlayer = currentPlayer
			machinePlayer.Captures = updatePlayer.Captures
			humanPlayer = state.humanPlayer
		} else {
			machinePlayer = state.machinePlayer
			humanPlayer = currentPlayer
			machinePlayer.Captures = updatePlayer.Captures
		}
		newStateChild := State{
			Node:            infoChild.Node,
			index:           infoChild.index,
			Captures:        infoChild.Captures,
			capturedIndexes: infoChild.capturedIndexes,
			machinePlayer:   machinePlayer,
			humanPlayer:     humanPlayer,
		}
		if err == nil {
			children = append(children, newStateChild)
			//children[k] = newPlayBoard
		}
	}

	//TO DO cache

	return children
}

func getAllIndexChildren(playBoard string) intSet {
	defer TimeTrack(time.Now(), "getAllIndexChildren", nil, nil)
	set := make(intSet)

	for index, val := range playBoard {
		value := string(val)
		if value != EmptySymbol {
			UpdateSetChildren(index, playBoard, set)
		}
	}

	return set
}

func copySet(children []State) intSet {
	defer TimeTrack(time.Now(), "copySet", RunTimesCopySet, AllTimesCopySet)
	setNewChildIndexes := make(intSet)

	for _, stateChild := range children {
		setNewChildIndexes[stateChild.index] = member
	}

	return setNewChildIndexes
}

func copySet2(children intSet) intSet {
	defer TimeTrack(time.Now(), "copySet2", RunTimesCopySet, AllTimesCopySet)
	setNewChildIndexes := make(intSet)

	for index, val := range children {
		setNewChildIndexes[index] = val
	}

	return setNewChildIndexes
}

type Child struct {
	h1    float64
	h2    float64
	Value float64
	State State
}

func sortChildren(children []State, transpositions stringSet, player Player, opponent Player, multiplier int) []Child {
	var new_ []Child

	for _, childState := range children {
		//_, ok := transpositions[childState.Node]
		//if ok {
		//	t += 1
		//	continue
		//}
		transpositions[childState.Node] = member
		h1, h2 := getHeuristic(childState, player, opponent)
		new_ = append(new_, Child{h1, h2, math.Max(h1, h2), childState})
	}

	sort.Slice(new_, func(i, j int) bool {
		if new_[i].Value == new_[j].Value {
			return new_[i].h1 > new_[j].h1
		}

		return new_[i].Value > new_[j].Value
	})

	//for i := range children {
	//	if i+1 < len(new_) {
	//		//if new_[i+1].Value < 100000 {
	//		//	new_ = new_[:i+1]
	//		//}
	//		if new_[i].Value > new_[i+1].Value {
	//			return new_[:i+1]
	//		} else if new_[i].h1 > new_[i+1].h1 {
	//			return new_[:i+1]
	//		}
	//	} else {
	//		break
	//	}
	//}

	return new_
}

type State struct {
	Node            string
	index           int
	Captures        int
	capturedIndexes []int
	machinePlayer   Player
	humanPlayer     Player
}

func (a Algo) NegaScout(state State, depth int, alpha float64, beta float64, multiplier int, childIndexesSet intSet, transpositions stringSet) (float64, int) {
	//PrintPlayBoard(state.Node)
	if depth == 0 || GameOver(state.Node, &state.machinePlayer, &state.humanPlayer, state.index) {
		h1, h2 := getHeuristic(state, state.machinePlayer, state.humanPlayer)

		if h1 == 1000000000000 {
			return float64(multiplier) * (h1 + (float64(depth) * 0.1)), state.index
		} else if h2 == 1000000000000 {
			return float64(multiplier) * (-h2 - (float64(depth) * 0.1)), state.index
		}

		return float64(multiplier) * (math.Max(h1, h2)), state.index
	}

	maxEval := math.Inf(-1)
	maxIndex := -1
	var children []State
	var childrenSlice []Child

	if state.index != -1 {
		UpdateSetChildren(state.index, state.Node, childIndexesSet)
	}
	if multiplier == 1 {
		children = getChildren(state, state.machinePlayer, childIndexesSet)
		childrenSlice = sortChildren(children, transpositions, state.machinePlayer, state.humanPlayer, multiplier)
	} else {
		children = getChildren(state, state.humanPlayer, childIndexesSet)
		childrenSlice = sortChildren(children, transpositions, state.humanPlayer, state.machinePlayer, multiplier)
	}

	for i, child := range childrenSlice {
		//setNewChildIndexes := copySet(children)
		setNewChildIndexes := copySet2(childIndexesSet)
		eval, _ := a.NegaScout(child.State, depth-1, -beta, -alpha, -multiplier, setNewChildIndexes, transpositions)
		eval = -eval

		if eval > maxEval {
			maxEval = eval

			if depth == a.Depth {
				maxIndex = child.State.index
			}
		}

		//if depth == a.Depth {
		//	PrintPlayBoard(state.Node)
		//}

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
	defer TimeTrackPrint(time.Now(), fmt.Sprintf("Algo "))
	if RunTimesHeuristic != nil {
		*RunTimesHeuristic = 0
		*RunTimesIsOver = 0
		*RunTimesgetChildren = 0
		*RunTimesCopySet = 0

		*AllTimesHeuristic = 0
		*AllTimesIsOver = 0
		*AllTimesgetChildren = 0
		*AllTimesCopySet = 0
	} else {
		_RunTimesHeuristic := 0
		RunTimesHeuristic = &_RunTimesHeuristic
		_RunTimesIsOver := 0
		RunTimesIsOver = &_RunTimesIsOver
		_RunTimesgetChildren := 0
		RunTimesgetChildren = &_RunTimesgetChildren
		_RunTimesCopySet := 0
		RunTimesCopySet = &_RunTimesCopySet

		var _AllTimesHeuristic time.Duration = 0
		AllTimesHeuristic = &_AllTimesHeuristic
		var _AllTimesIsOver time.Duration = 0
		AllTimesIsOver = &_AllTimesIsOver
		var _AllTimesgetChildren time.Duration = 0
		AllTimesgetChildren = &_AllTimesgetChildren
		var _AllTimesCopySet time.Duration = 0
		AllTimesCopySet = &_AllTimesCopySet
	}

	setChildren := getAllIndexChildren(playBoard)

	if len(setChildren) == 0 {
		return 9*19 + 9
	}

	var transpositions = make(stringSet)

	_, index := a.NegaScout(State{playBoard, -1, 0, []int{}, machinePlayer, humanPlayer}, a.Depth, math.Inf(-1), math.Inf(1), 1, setChildren, transpositions)

	return index
}
