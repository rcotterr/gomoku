package algo

import (
	"fmt"
	"gomoku/pkg/playboard"
	"math"
	"sort"
	"time"
)

type void struct{}
type myset map[int]void
type mysetString map[string]void

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

func Heuristic(node string, index int) float64 {
	defer playboard.TimeTrack(time.Now(), "Heuristic", playboard.RunTimesHeuristic, playboard.AllTimesHeuristic)
	num := 0.0

	//count по каждой стороне
	symbol := string(node[index])

	setRules := map[int]playboard.ConditionFn{
		1:               playboard.ConditionHorizontalCapture,
		playboard.N:     playboard.ConditionVertical,
		playboard.N + 1: playboard.ConditionRightDiagonalCapture, //TO DO delete duplicate conditionRightDiagonal
		playboard.N - 1: playboard.ConditionLeftDiagonalCheckFiveStones,
	}

	for step, condition := range setRules {
		count, halfFree, free := countInRow(node, index, step, condition, symbol)
		if count == 5 { // TO DO and not capture
			num = 1000000000
			break
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

func UpdateSetChildren(index int, playBoard string, set myset) {
	for step, condition := range setRulesChildren {
		j := index + step
		if condition(j, index) && j >= 0 && j < playboard.N*playboard.N && string(playBoard[j]) == playboard.EmptySymbol {
			set[j] = member
		}
	}
}

func getChildren(node string, index int, currentPlayer playboard.Player, childIndexesSet myset) map[int]string {
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

func getAllIndexChildren(playBoard string) myset {
	defer playboard.TimeTrack(time.Now(), "getAllIndexChildren", nil, nil)
	set := make(myset)

	for index, val := range playBoard {
		value := string(val)
		if value != playboard.EmptySymbol {
			UpdateSetChildren(index, playBoard, set)
		}
	}

	return set
}

func copySet(children map[int]string) myset {
	defer playboard.TimeTrack(time.Now(), "copySet", playboard.RunTimesCopySet, playboard.AllTimesCopySet)
	setNewChildIndexes := make(myset)

	for key := range children {
		setNewChildIndexes[key] = member
	}

	return setNewChildIndexes
}

func getIndexes(children []Child) myset {
	defer playboard.TimeTrack(time.Now(), "copySet", playboard.RunTimesCopySet, playboard.AllTimesCopySet)
	setNewChildIndexes := make(myset)

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

func cutChildren(children map[int]string, transpositions mysetString) []Child {

	var new_ []Child
	for childIndex, childPlayboard := range children {
		_, ok := transpositions[childPlayboard]
		if ok {
			t += 1
			continue
		}
		transpositions[childPlayboard] = member
		h := Heuristic(childPlayboard, childIndex)
		new_ = append(new_, Child{childIndex, h, childPlayboard})

	}

	sort.Slice(new_, func(i, j int) bool {
		return new_[i].Value > new_[j].Value
	})
	if len(new_) > numChildren {
		new_ = new_[:numChildren]
	}

	return new_
}

func alphaBeta(node string, depth int, alpha float64, beta float64, maximizingPlayer bool, machinePlayer playboard.Player, humanPlayer playboard.Player, index int, childIndexesSet myset, transpositions mysetString, allIndexesPath string) (float64, int, string, int) {
	defer playboard.TimeTrack(time.Now(), fmt.Sprintf("alphaBeta depth {%d}", depth), nil, nil)

	//if depth == 0 || playboard.IsOver(node, &machinePlayer, &humanPlayer) {
	if depth == 0 {
		return Heuristic(node, index), index, node, depth //TO DO static evaluation of node
	}
	if playboard.IsOver(node, &machinePlayer, &humanPlayer) {
		h := Heuristic(node, index)
		i := depth * 1000.0
		h = h + float64(i)
		return h, index, node, depth
	}
	if maximizingPlayer {
		maxEval := math.Inf(-1)
		maxIndex := 0
		depth_ := -1
		children := getChildren(node, index, machinePlayer, childIndexesSet)
		childrenSlice := cutChildren(children, transpositions)
		//breakCount := 4
		for _, child := range childrenSlice {
			//if breakCount == 0 {
			//	break
			//}
			//_, ok := transpositions[child.PlayBoard]
			//if ok {
			//	t += 1
			//	continue
			//}
			//transpositions[child.PlayBoard] = member
			setNewChildIndexes := copySet(children)
			//setNewChildIndexes := getIndexes(childrenSlice)

			eval, _, tmpIndPath, tmpDepth := alphaBeta(child.PlayBoard, depth-1, alpha, beta, false, machinePlayer, humanPlayer, child.Index, setNewChildIndexes, transpositions, allIndexesPath)

			if eval > maxEval && tmpDepth >= depth_ {
				maxEval = eval
				depth_ = tmpDepth
				maxIndex = child.Index
				allIndexesPath = tmpIndPath
			}
			alpha = math.Max(alpha, eval)
			if beta <= alpha {
				break
			}
			//breakCount -= 1
		}
		return maxEval, maxIndex, allIndexesPath, depth_
	} else {
		minEval := math.Inf(1)
		minIndex := 0
		children := getChildren(node, index, humanPlayer, childIndexesSet)
		depth_ := -1
		childrenSlice := cutChildren(children, transpositions)
		//breakCount := 4
		for _, child := range childrenSlice {
			//if breakCount == 0 {
			//	break
			//}
			//_, ok := transpositions[childPlayboard]
			//if ok {
			//	t += 1
			//	continue
			//}
			//transpositions[childPlayboard] = member

			setNewChildIndexes := copySet(children)

			eval, _, tmpIndPath, tmpDepth := alphaBeta(child.PlayBoard, depth-1, alpha, beta, true, machinePlayer, humanPlayer, child.Index, setNewChildIndexes, transpositions, allIndexesPath)

			if eval < minEval { //because both values are + inf eval <= minEval
				minEval = eval
				minIndex = child.Index
				allIndexesPath = tmpIndPath
				depth_ = tmpDepth
			}
			beta = math.Min(beta, eval)
			if beta <= alpha {
				break
			}
			//breakCount -= 1
		}
		return minEval, minIndex, allIndexesPath, depth_

	}
}

func Algo(playBoard string, machinePlayer playboard.Player, humanPlayer playboard.Player) int {
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

	depth := 10 //TO DO make config

	negInf := math.Inf(-1)
	posInf := math.Inf(1)
	index := -1
	val := 0.0

	t = 0
	var transpositions = make(mysetString)

	setChildren := getAllIndexChildren(playBoard)
	if len(setChildren) != 0 {
		allIndexesPath := ""
		val, index, allIndexesPath, depth = alphaBeta(playBoard, depth, negInf, posInf, true, machinePlayer, humanPlayer, index, setChildren, transpositions, allIndexesPath)
		fmt.Println(val, depth)
		playboard.PrintPlayBoard(allIndexesPath)
	} else {
		index = 9*19 + 9
		//TO DO random from 3-15/3-15
	}

	fmt.Println(index)

	return index
}
