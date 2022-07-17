package algo

import (
	"fmt"
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

//const maxValueWin = 1000000000

var setRulesChildren = map[int]playboard.ConditionFn{
	1:                playboard.ConditionHorizontal,
	playboard.N:      playboard.ConditionVertical,
	playboard.N + 1:  playboard.ConditionBackDiagonal,
	playboard.N - 1:  playboard.ConditionForwardDiagonal,
	-1:               playboard.ConditionHorizontal,
	-playboard.N:     playboard.ConditionVertical,
	-playboard.N - 1: playboard.ConditionBackDiagonal,
	-playboard.N + 1: playboard.ConditionForwardDiagonal,
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
		1:               playboard.ConditionHorizontal,
		playboard.N:     playboard.ConditionVertical,
		playboard.N + 1: playboard.ConditionBackDiagonal,
		playboard.N - 1: playboard.ConditionForwardDiagonal,
	}

	for step, condition := range setRules {
		count, halfFree, free := playboard.CountInRow(node, index, step, condition, symbol)
		if count >= 5 { // TO DO and not capture
			//num = 1000000000
			//break
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

func UpdateSetChildren(index int, playBoard string, set intSet) {
	for step, condition := range setRulesChildren {
		j := index + step
		if j >= 0 && j < playboard.N*playboard.N && condition(j, index) && string(playBoard[j]) == playboard.EmptySymbol {
			set[j] = member
		}
	}
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

func getAllIndexChildren(playBoard string) intSet {
	defer playboard.TimeTrack(time.Now(), "getAllIndexChildren", nil, nil)
	set := make(intSet)

	for index, val := range playBoard {
		value := string(val)
		if value != playboard.EmptySymbol {
			UpdateSetChildren(index, playBoard, set)
		}
	}

	return set
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

func alphaBeta(node string, depth int, alpha float64, beta float64, maximizingPlayer bool, machinePlayer playboard.Player, humanPlayer playboard.Player, index int, childIndexesSet intSet, transpositions stringSet, allIndexesPath string) (float64, int, string, int) {
	defer playboard.TimeTrack(time.Now(), fmt.Sprintf("alphaBeta depth {%d}", depth), nil, nil)

	if depth == 0 || playboard.GameOver(node, &machinePlayer, &humanPlayer, index) {
		symbol := string(node[index])
		//h1 := Heuristic(node, index, symbol) + float64(depth * 1000)
		h1 := Heuristic(node, index, symbol)
		if symbol == machinePlayer.Symbol {
			symbol = humanPlayer.Symbol
			node = strings.Join([]string{node[:index], symbol, node[index+1:]}, "")

		} else if symbol == humanPlayer.Symbol {
			symbol = machinePlayer.Symbol
			node = strings.Join([]string{node[:index], symbol, node[index+1:]}, "")

		}
		//h2 := Heuristic(node, index, symbol) + float64(depth * 1000)
		h2 := Heuristic(node, index, symbol)
		return h1 - h2, index, node, depth
	}
	if maximizingPlayer {
		maxEval := math.Inf(-1)
		maxIndex := 0
		depth_ := -1
		children := getChildren(node, index, machinePlayer, childIndexesSet)
		childrenSlice := cutChildren(children, transpositions)
		for _, child := range childrenSlice {
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
		}
		return maxEval, maxIndex, allIndexesPath, depth_
	} else {
		minEval := math.Inf(1)
		minIndex := 0
		children := getChildren(node, index, humanPlayer, childIndexesSet)
		depth_ := -1
		childrenSlice := cutChildren(children, transpositions)
		for _, child := range childrenSlice {

			setNewChildIndexes := copySet(children)

			eval, _, tmpIndPath, tmpDepth := alphaBeta(child.PlayBoard, depth-1, alpha, beta, true, machinePlayer, humanPlayer, child.Index, setNewChildIndexes, transpositions, allIndexesPath)
			//eval = -eval
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
		}
		return minEval, minIndex, allIndexesPath, depth_

	}
}

func NegaScout(node string, depth int, alpha float64, beta float64, multiplier int, machinePlayer playboard.Player, humanPlayer playboard.Player, index int, childIndexesSet intSet, transpositions stringSet, allIndexesPath string) (float64, int) {
	//println("depth", depth)
	if depth == 0 || playboard.GameOver(node, &machinePlayer, &humanPlayer, index) {
		var h1, h2 float64

		if string(node[index]) == machinePlayer.Symbol {
			h1 = Heuristic(node, index, machinePlayer.Symbol)
			node = strings.Join([]string{node[:index], humanPlayer.Symbol, node[index+1:]}, "")
			h2 = Heuristic(node, index, humanPlayer.Symbol)
		} else {
			h2 = Heuristic(node, index, humanPlayer.Symbol)
			node = strings.Join([]string{node[:index], machinePlayer.Symbol, node[index+1:]}, "")
			h1 = Heuristic(node, index, machinePlayer.Symbol)
		}

		if h1 == math.Inf(1) {
			// println("hi");
		}

		return float64(multiplier) * (h1 - h2), index
	}

	maxEval := math.Inf(-1)
	maxIndex := -1
	var children map[int]string
	if multiplier == 1 {
		children = getChildren(node, index, machinePlayer, childIndexesSet)
	} else {
		children = getChildren(node, index, humanPlayer, childIndexesSet)
	}
	childrenSlice := cutChildren(children, transpositions)

	for _, child := range childrenSlice {
		setNewChildIndexes := copySet(children)
		eval, _ := NegaScout(child.PlayBoard, depth-1, -alpha, -beta, -multiplier, machinePlayer, humanPlayer, child.Index, setNewChildIndexes, transpositions, allIndexesPath)

		eval = -eval
		if eval >= maxEval { //because both values are -inf eval >= maxEval
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
	//val := 0.0

	t = 0
	var transpositions = make(stringSet)

	setChildren := getAllIndexChildren(playBoard)
	if len(setChildren) != 0 {
		allIndexesPath := ""
		//val, index, allIndexesPath, depth = alphaBeta(playBoard, depth, negInf, posInf, true, machinePlayer, humanPlayer, index, setChildren, transpositions, allIndexesPath)
		_, index = NegaScout(playBoard, depth, negInf, posInf, 1, machinePlayer, humanPlayer, index, setChildren, transpositions, allIndexesPath)
		//fmt.Println(val, depth)
		playboard.PrintPlayBoard(allIndexesPath)
	} else {
		//index = 18*19 + 18
		index = 9*19 + 9
		//TO DO random from 3-15/3-15
	}

	fmt.Println(index)

	return index
}
