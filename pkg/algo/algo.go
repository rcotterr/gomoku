package algo

import (
	"fmt"
	"gomoku/pkg/playboard"
	"math"
	"sync"
	"time"
)

type void struct{}
type myset map[int]void

var member void

var m sync.Mutex

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
	defer playboard.TimeTrack(time.Now(), "Heuristic")
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

func getChildren(node string, index int, currentPlayer playboard.Player, setChildrenIndexes myset) map[int]string {
	defer playboard.TimeTrack(time.Now(), "getChildren")
	var children = make(map[int]string)

	m.Lock()
	if index != -1 {
		UpdateSetChildren(index, node, setChildrenIndexes)
	}
	for k := range setChildrenIndexes {
		newPlayBoard, err := playboard.PutStone(node, k, &currentPlayer)
		if err == nil {
			//m.Lock()
			children[k] = newPlayBoard
			//m.Unlock()
		}
	}
	m.Unlock()

	//TO DO cache

	return children
}

func getAllIndexChildren(playBoard string) myset {
	defer playboard.TimeTrack(time.Now(), "getAllIndexChildren")
	set := make(map[int]void)

	for index, val := range playBoard {
		value := string(val)
		if value != playboard.EmptySymbol {
			UpdateSetChildren(index, playBoard, set)
		}
	}

	return set
}

type result struct {
	eval       float64
	childIndex int
}

func alphaBeta(node string, depth int, alpha float64, beta float64, maximizingPlayer bool, machinePlayer playboard.Player, humanPlayer playboard.Player, index int, setChildrenIndexes myset) (float64, int) {
	defer playboard.TimeTrack(time.Now(), fmt.Sprintf("alphaBeta depth {%d}", depth))
	if depth == 0 || playboard.IsOver(node, &machinePlayer, &humanPlayer) {
		return Heuristic(node, index), index //TO DO static evaluation of node
	}
	if maximizingPlayer {
		maxEval := math.Inf(-1)
		maxIndex := 0
		children := getChildren(node, index, machinePlayer, setChildrenIndexes)

		ch := make(chan result)
		defer close(ch)

		for childIndex, childPlayboard := range children {
			go func(childPlayboard string, childIndex int) {
				setNewChildIndexes := make(map[int]void)

				for key := range children {
					setNewChildIndexes[key] = member
				}
				if depth == 4 {
					//fmt.Println("")
				}
				eval, _ := alphaBeta(childPlayboard, depth-1, alpha, beta, false, machinePlayer, humanPlayer, childIndex, setNewChildIndexes)
				ch <- result{eval, childIndex}
			}(childPlayboard, childIndex)
		}
		for i := 0; i < len(children); i++ {
			result_ := <-ch
			if result_.eval > maxEval {
				maxEval = result_.eval
				maxIndex = result_.childIndex
			}
			if result_.eval >= alpha {
				alpha = result_.eval
				if beta <= alpha {
					return maxEval, maxIndex
				}
			}
		}

		return maxEval, maxIndex
	} else {
		minEval := math.Inf(1)
		minIndex := 0

		children := getChildren(node, index, machinePlayer, setChildrenIndexes)

		ch := make(chan result)
		defer close(ch)

		for childIndex, childPlayboard := range children {
			go func(childPlayboard string, childIndex int) {
				setNewChildIndexes := make(map[int]void)

				for key := range children {
					setNewChildIndexes[key] = member
				}

				eval, _ := alphaBeta(childPlayboard, depth-1, alpha, beta, true, machinePlayer, humanPlayer, childIndex, setNewChildIndexes)
				ch <- result{eval, childIndex}
			}(childPlayboard, childIndex)
		}

		for i := 0; i < len(children); i++ {
			result_ := <-ch

			if result_.eval < minEval { // TO DO make max func
				minEval = result_.eval
				minIndex = result_.childIndex
			}
			if result_.eval <= beta {
				beta = result_.eval
				if beta <= alpha {
					return minEval, minIndex
				}
			}
		}
		return minEval, minIndex
	}
}

func Algo(playBoard string, machinePlayer playboard.Player, humanPlayer playboard.Player) int {
	depth := 4 //TO DO make config

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
