package algo

import (
	"fmt"
	"gomoku/pkg/playboard"
	"math"
	"time"
)

type void struct{}
type myset map[int]void

var member void

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
		if count == 5 {
			num += 25000
		} else if count == 4 && free {
			num += 16000
		} else if count == 4 && halfFree {
			num += 12000
		} else if count == 3 && free {
			num += 9000
		} else if count == 3 && halfFree {
			num += 6750
		} else if count == 2 && free {
			num += 4000
		} else if count == 2 && halfFree {
			num += 3000
		} else if count == 1 && free {
			num += 1000
		} else if count == 1 && halfFree {
			num += 750
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

func alphaBeta(node string, depth int, alpha float64, beta float64, maximizingPlayer bool, machinePlayer playboard.Player, humanPlayer playboard.Player, index int, childIndexesSet myset) (float64, int) {
	defer playboard.TimeTrack(time.Now(), fmt.Sprintf("alphaBeta depth {%d}", depth), nil, nil)
	if depth == 0 || playboard.IsOver(node, &machinePlayer, &humanPlayer) {
		return Heuristic(node, index), index //TO DO static evaluation of node
	}
	if maximizingPlayer {
		maxEval := math.Inf(-1)
		maxIndex := 0
		children := getChildren(node, index, machinePlayer, childIndexesSet)

		for childIndex, childPlayboard := range children {
			setNewChildIndexes := copySet(children)

			eval, _ := alphaBeta(childPlayboard, depth-1, alpha, beta, false, machinePlayer, humanPlayer, childIndex, setNewChildIndexes)
			if eval > maxEval {
				maxEval = eval
				maxIndex = childIndex
			}
			if eval >= alpha {
				alpha = eval
				if beta <= alpha {
					break
				}
			}
		}
		return maxEval, maxIndex
	} else {
		minEval := math.Inf(1)
		minIndex := 0
		children := getChildren(node, index, machinePlayer, childIndexesSet)

		for childIndex, childPlayboard := range children {
			setNewChildIndexes := copySet(children)

			eval, _ := alphaBeta(childPlayboard, depth-1, alpha, beta, true, machinePlayer, humanPlayer, childIndex, setNewChildIndexes)
			if eval < minEval { // TO DO make max func
				minEval = eval
				minIndex = childIndex
			}
			if eval <= beta {
				beta = eval
				if beta <= alpha {
					break
				}
			}
		}
		return minEval, minIndex

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

	depth := 5 //TO DO make config

	negInf := math.Inf(-1)
	posInf := math.Inf(1)
	index := -1

	setChildren := getAllIndexChildren(playBoard)
	if len(setChildren) != 0 {
		_, index = alphaBeta(playBoard, depth, negInf, posInf, true, machinePlayer, humanPlayer, index, setChildren)
	} else {
		index = 9*19 + 9
		//TO DO random from 3-15/3-15
	}

	fmt.Println(index)

	return index
}
