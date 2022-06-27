package algo

import (
	"fmt"
	"gomoku/pkg/playboard"
	"math"
)

func Heuristic() float64 {
	return 1
}

func getChildren(node string, index int, currentPlayer playboard.Player) map[int]string {
	var children = make(map[int]string)

	setRules := map[int]playboard.ConditionFn{
		1:                playboard.ConditionHorizontalCapture,
		playboard.N:      playboard.ConditionVertical,
		playboard.N + 1:  playboard.ConditionRightDiagonalCapture, //TO DO delete duplicate conditionRightDiagonal
		playboard.N - 1:  playboard.ConditionLeftDiagonal,
		-1:               playboard.ConditionHorizontalCapture,
		-playboard.N:     playboard.ConditionVertical,
		-playboard.N - 1: playboard.ConditionRightDiagonalCapture,
		-playboard.N + 1: playboard.ConditionLeftUpperDiagonal,
	}

	for step, condition := range setRules {
		j := index + step
		if condition(j, index) && j >= 0 && j < playboard.N*playboard.N && string(node[j]) == playboard.EmptySymbol {
			newPlayBoard, err := playboard.PutStone(node, j, &currentPlayer)
			if err == nil {
				children[j] = newPlayBoard
			}
		}
	}

	//TO DO cache

	return children
}

func alphaBeta(node string, depth int, alpha float64, beta float64, maximizingPlayer bool, machinePlayer playboard.Player, humanPlayer playboard.Player, index int) (float64, int) {
	if depth == 0 || playboard.IsOver(node, &machinePlayer, &humanPlayer) {
		return Heuristic(), index //TO DO static evaluation of node
	}
	if maximizingPlayer {
		maxEval := math.Inf(-1)
		maxIndex := -1
		for childIndex, childPlayboard := range getChildren(node, index, machinePlayer) {
			eval, ind := alphaBeta(childPlayboard, depth-1, alpha, beta, false, machinePlayer, humanPlayer, childIndex)
			if eval > maxEval {
				maxEval = eval
				maxIndex = ind
			}
			if eval > alpha {
				alpha = eval
				if beta <= alpha {
					break
				}
			}
		}
		return maxEval, maxIndex
	} else {
		minEval := math.Inf(1)
		minIndex := -1
		for childIndex, childPlayboard := range getChildren(node, index, humanPlayer) {
			eval, ind := alphaBeta(childPlayboard, depth-1, alpha, beta, true, machinePlayer, humanPlayer, childIndex)
			if eval < minEval {
				minEval = eval
				minIndex = ind
			}
			if eval < beta {
				beta = eval
				if beta <= alpha {
					break
				}
			}
		}
		return minEval, minIndex

	}
}

func MinMaxAlgo(playBoard string, machinePlayer playboard.Player, humanPlayer playboard.Player) int {
	depth := 2 //TO DO make config

	negInf := math.Inf(-1)
	posInf := math.Inf(1)
	valuePlayboard := math.Inf(-1)
	index := 0

	for i, val := range playBoard { // TO DO not all check but only 1 last put stone  && not range but while i < len
		value := string(val)
		if value == machinePlayer.Symbol || value == humanPlayer.Symbol {
			tmpValuePlayboard, tmpIndex := alphaBeta(playBoard, depth, negInf, posInf, true, machinePlayer, humanPlayer, i) // TO DO index is not valid return last bot first now
			if tmpValuePlayboard > valuePlayboard {
				valuePlayboard = tmpValuePlayboard
				index = tmpIndex
			}
		}
	}
	if valuePlayboard == math.Inf(-1) {
		index = 9*19 + 9
		//TO DO random from 3-15/3-15
	}

	fmt.Println(index)

	return index
}
