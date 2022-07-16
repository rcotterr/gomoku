package algo

import (
	"gomoku/pkg/playboard"
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

var setRulesChildren = map[int]playboard.ConditionFn{
	1:                playboard.ConditionHorizontalCapture,
	playboard.N:      playboard.ConditionVertical,
	playboard.N + 1:  playboard.ConditionRightDiagonalCapture, //TO DO delete duplicate conditionRightDiagonal
	playboard.N - 1:  playboard.ConditionLeftDiagonal,
	-1:               playboard.ConditionHorizontalCapture,
	-playboard.N:     playboard.ConditionVertical,
	-playboard.N - 1: playboard.ConditionRightDiagonalCapture,
	-playboard.N + 1: playboard.ConditionLeftUpperDiagonal,
}

func countInRow(node string, index int, step int, condition playboard.ConditionFn, symbol string) (int, bool, bool) {

	empty := 0

	hreny := node[index:]
	hreny = string(hreny)
	startIndex := index
	for tmpIndex := index - step; tmpIndex > 0 && tmpIndex > index-(step*5); tmpIndex -= step {
		if condition(tmpIndex, index) && string(node[tmpIndex]) == symbol {
			startIndex = tmpIndex
		} else {
			if string(node[tmpIndex]) == playboard.EmptySymbol {
				empty += 1
			}
			break
		}
	}

	endIndex := index
	for tmpIndex := index + step; tmpIndex < playboard.N*playboard.N && tmpIndex < index+(step*5); tmpIndex += step {
		if condition(tmpIndex, index) && string(node[tmpIndex]) == symbol {
			endIndex = tmpIndex
		} else {
			if string(node[tmpIndex]) == playboard.EmptySymbol {
				empty += 1
			}
			break
		}
	}

	count := ((endIndex - startIndex) / step) + 1

	return count, empty == 1, empty == 2

}

func Heuristic(node string, index int, symbol string) float64 {
	defer playboard.TimeTrack(time.Now(), "Heuristic", playboard.RunTimesHeuristic, playboard.AllTimesHeuristic)
	num := 0.0

	//count по каждой стороне
	//symbol := string(node[index])
	//symbol := "M"
	if symbol == "" {
		symbol = string(node[index])
	}

	setRules := map[int]playboard.ConditionFn{
		1:               playboard.ConditionHorizontalCapture,
		playboard.N:     playboard.ConditionVertical,
		playboard.N + 1: playboard.ConditionRightDiagonalCapture, //TO DO delete duplicate conditionRightDiagonal
		playboard.N - 1: playboard.ConditionLeftDiagonalCheckFiveStones,
	}

	for step, condition := range setRules {
		count, halfFree, free := countInRow(node, index, step, condition, symbol)
		if count == 5 { // TO DO and not capture
			// playboard.PrintPlayBoard(node)
			return math.Inf(1)
		} else if count == 4 && free {
			num += 10000000
		} else if count == 4 && halfFree {
			num += 1000000
		} else if count == 3 && free {
			num += 100000
		} else if count == 3 && halfFree {
			num += 10000
		} else if count == 2 && free {
			num += 1000
		} else if count == 2 && halfFree {
			num += 100
		} else if count == 1 && free {
			num += 10
		} else if count == 1 && halfFree {
			num += 1
		}

	}

	//add capture

	return num
}

func getChildren(node string, index int, currentPlayer playboard.Player, childIndexesSet intSet) map[int]string {
	defer playboard.TimeTrack(time.Now(), "getChildren", playboard.RunTimesgetChildren, playboard.AllTimesgetChildren)
	var children = make(map[int]string)

	if index != -1 {
		UpdateSetChildren(index, node, childIndexesSet)
	}

	for k := range childIndexesSet {
		newPlayBoard, err := playboard.PutStone(node, k, &currentPlayer)
		if err == nil {
			children[k] = newPlayBoard
		}
	}

	//TO DO cache

	return children
}



func copySet(children map[int]string) intSet {
	defer playboard.TimeTrack(time.Now(), "copySet", playboard.RunTimesCopySet, playboard.AllTimesCopySet)
	setNewChildIndexes := make(intSet)

	for key := range children {
		setNewChildIndexes[key] = member
	}

	return setNewChildIndexes
}

func getIndexes(children []Child) intSet {
	defer playboard.TimeTrack(time.Now(), "copySet", playboard.RunTimesCopySet, playboard.AllTimesCopySet)
	setNewChildIndexes := make(intSet)

	for _, child := range children {
		setNewChildIndexes[child.Index] = member
	}

	return setNewChildIndexes
}

type Child struct {
	Index     int
	Value     float64
	PlayBoard string
}

const numChildren = 7

func cutChildren(children map[int]string, transpositions stringSet) []Child {

	var new_ []Child
	for childIndex, childPlayboard := range children {
		_, ok := transpositions[childPlayboard]
		if ok {
			t += 1
			continue
		}
		transpositions[childPlayboard] = member
		h := Heuristic(childPlayboard, childIndex, "")
		new_ = append(new_, Child{childIndex, h, childPlayboard})

	}

	sort.Slice(new_, func(i, j int) bool {
		return new_[i].Value > new_[j].Value
	})
	//if len(new_) > numChildren {
	//	new_ = new_[:numChildren]
	//}

	return new_
}

func NegaScout(node string, depth int, alpha float64, beta float64, multiplier int, machinePlayer playboard.Player, humanPlayer playboard.Player, index int, childIndexesSet intSet, transpositions stringSet) (float64, int) {
	// println("depth", depth)
	if depth == 0 || playboard.GameOver(node, &machinePlayer, &humanPlayer) {
		var h1, h2 float64

		if string(node[index]) == machinePlayer.Symbol {
			h1 = Heuristic(node, index, machinePlayer.Symbol)
			node = strings.Join([]string{node[:index], humanPlayer.Symbol, node[index + 1:]}, "")
			h2 = Heuristic(node, index, humanPlayer.Symbol)
		} else {
			h2 = Heuristic(node, index, humanPlayer.Symbol)
			node = strings.Join([]string{node[:index], machinePlayer.Symbol, node[index + 1:]}, "")
			h1 = Heuristic(node, index, machinePlayer.Symbol)
		}

		return float64(multiplier) * (h1 - h2), index
	}

	maxEval := math.Inf(-1)
	maxIndex := -1
	var children map[int]string
	if (multiplier == 1) {
		children = getChildren(node, index, machinePlayer, childIndexesSet)
	} else {
		children = getChildren(node, index, humanPlayer, childIndexesSet)
	}
	childrenSlice := cutChildren(children, transpositions)

	for _, child := range childrenSlice {
		setNewChildIndexes := copySet(children)
		eval, _ := NegaScout(child.PlayBoard, depth - 1, -beta, -alpha, -multiplier, machinePlayer, humanPlayer, child.Index, setNewChildIndexes, transpositions)

		eval = -eval;
		if eval > maxEval {
			maxEval = eval
			maxIndex = child.Index
		}

		alpha = math.Max(alpha, eval)

		if alpha >= beta {
			break
		}
	}

	return maxEval, maxIndex
}

func UpdateSetChildren(index int, board string, set intSet) {
	for step, condition := range setRulesChildren {
		j := index + step

		if j >= 0 && j < playboard.N * playboard.N {
			if string(board[j]) == playboard.EmptySymbol && condition(j, index) {
				set[j] = member
			}
		}
	}
}

func getAllIndexChildren(board string) intSet {
	set := make(intSet)

	for index, val := range board {
		value := string(val)
		if value != playboard.EmptySymbol {
			UpdateSetChildren(index, board, set)
		}
	}

	return set
}

func GetMachineIndex(board string, machinePlayer playboard.Player, humanPlayer playboard.Player) int {
	if playboard.RunTimesHeuristic != nil {
		*playboard.RunTimesHeuristic = 0
		*playboard.RunTimesIsOver = 0
		*playboard.RunTimesgetChildren = 0
		*playboard.RunTimesCopySet = 0

		*playboard.AllTimesHeuristic = 0
		*playboard.AllTimesIsOver = 0
		*playboard.AllTimesgetChildren = 0
		*playboard.AllTimesCopySet = 0
	} else {
		RunTimesHeuristic := 0
		playboard.RunTimesHeuristic = &RunTimesHeuristic
		RunTimesIsOver := 0
		playboard.RunTimesIsOver = &RunTimesIsOver
		RunTimesgetChildren := 0
		playboard.RunTimesgetChildren = &RunTimesgetChildren
		RunTimesCopySet := 0
		playboard.RunTimesCopySet = &RunTimesCopySet

		var AllTimesHeuristic time.Duration = 0
		playboard.AllTimesHeuristic = &AllTimesHeuristic
		var AllTimesIsOver time.Duration = 0
		playboard.AllTimesIsOver = &AllTimesIsOver
		var AllTimesgetChildren time.Duration = 0
		playboard.AllTimesgetChildren = &AllTimesgetChildren
		var AllTimesCopySet time.Duration = 0
		playboard.AllTimesCopySet = &AllTimesCopySet
	}

	depth := 5
	index := -1

	var transpositions = make(stringSet)

	moves := getAllIndexChildren(board)
	println(moves)
	if len(moves) != 0 {
		_, index = NegaScout(board, depth, math.Inf(-1), math.Inf(1), 1, machinePlayer, humanPlayer, index, moves, transpositions)
	} else {
		index = 9 * 19 + 9
	}

	if (index == -1) {
		println("hm");
	}

	return index
}
