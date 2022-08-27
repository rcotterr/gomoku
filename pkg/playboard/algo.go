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

func Heuristic(state State, symbol string, captures int, capturedIndexes []int) float64 {
	defer TimeTrack(time.Now(), "Heuristic", RunTimesHeuristic, AllTimesHeuristic)
	vulnerable := false
	num := 0.0

	//count по каждой стороне
	//symbol := string(node[index])
	if symbol == "" {
		symbol = string(state.node[state.index])
	}

	setRules := map[int]ConditionFn{
		1:     ConditionHorizontal,
		N:     ConditionVertical,
		N + 1: ConditionBackDiagonal,
		N - 1: ConditionForwardDiagonal,
	}

	for step, condition := range setRules {
		count, halfFree, free, _ := CountInRow(state.node, state.index, step, condition, symbol)
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

	num += float64(1000000000 * state.captures)
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
	state.node = strings.Join([]string{state.node[:state.index], EmptySymbol, state.node[state.index+1:]}, "")
	state, err := PutStone(state, state.index, &opponent)
	if err == nil {
		h2 = Heuristic(state, opponent.Symbol, opponent.Captures, state.capturedIndexes)
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

func getChildren(state State, childIndexesSet intSet, multiplier int) []State {
	defer TimeTrack(time.Now(), "getChildren", RunTimesgetChildren, AllTimesgetChildren)
	var children []State

	//if index != -1 {
	//	UpdateSetChildren(index, node, childIndexesSet)
	//}

	for k := range childIndexesSet {
		var currentPlayer Player

		if multiplier == 1 {
			currentPlayer = Player{state.machinePlayer.Captures, state.machinePlayer.Symbol, state.machinePlayer.IndexAlmostWin, state.machinePlayer.Winner}
		} else {
			currentPlayer = Player{state.humanPlayer.Captures, state.humanPlayer.Symbol, state.humanPlayer.IndexAlmostWin, state.humanPlayer.Winner}
		}

		stateChild, err := PutStone(state, k, &currentPlayer)

		if err == nil {
			children = append(children, stateChild)
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
	state State
	value float64
	h1    float64
	h2    float64
}

func sortChildren(children []State, player Player, opponent Player) []Child {
	var new_ []Child

	for _, childState := range children {
		h1, h2 := getHeuristic(childState, player, opponent)
		new_ = append(new_, Child{childState, math.Max(h1, h2), h1, h2})
	}

	sort.Slice(new_, func(i, j int) bool {
		if new_[i].value == new_[j].value {
			return new_[i].h1 > new_[j].h1
		}

		return new_[i].value > new_[j].value
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
	node            string
	index           int
	captures        int
	capturedIndexes []int
	machinePlayer   Player
	humanPlayer     Player
}

func (a Algo) NegaScout(state State, depth int, alpha float64, beta float64, multiplier int, childIndexesSet intSet) (float64, int) {
	if depth == 0 || GameOver(state.node, &state.machinePlayer, &state.humanPlayer, state.index) {
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
		UpdateSetChildren(state.index, state.node, childIndexesSet)
	}

	children = getChildren(state, childIndexesSet, multiplier)
	if multiplier == 1 {
		childrenSlice = sortChildren(children, state.machinePlayer, state.humanPlayer)
	} else {
		childrenSlice = sortChildren(children, state.humanPlayer, state.machinePlayer)
	}

	for i, child := range childrenSlice {
		PrintPlayBoard(child.state.node)
		if depth == 10 && child.state.index == 142 {
			print("")
		}
		setNewChildIndexes := copySet2(childIndexesSet)
		eval, _ := a.NegaScout(child.state, depth-1, -beta, -alpha, -multiplier, setNewChildIndexes)
		eval = -eval

		if eval > maxEval {
			maxEval = eval

			if depth == a.Depth {
				maxIndex = child.state.index
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
				if child.value > childrenSlice[i+1].value || child.h1 > childrenSlice[i+1].h1 {
					break
				}
			}
		} else if maxEval > 0 {
			if i+1 < len(childrenSlice) {
				if child.value > childrenSlice[i+1].value || child.h1 > childrenSlice[i+1].h1 {
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

	var initState = State{
		playBoard,
		-1,
		0,
		[]int{},
		machinePlayer,
		humanPlayer,
	}

	_, index := a.NegaScout(initState, a.Depth, math.Inf(-1), math.Inf(1), 1, setChildren)

	return index
}
