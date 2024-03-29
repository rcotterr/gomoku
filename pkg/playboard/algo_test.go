package playboard

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestGetChildren(t *testing.T) {
	testCases := []struct {
		name               string
		playboard          string
		index              int
		setChildrenIndexes intSet
		currentPlayer      Player
		expectedChildren   []int
	}{
		{
			name: "all 8 children",
			playboard: "..................." +
				"....M.............." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			index:              23,
			setChildrenIndexes: intSet{},
			currentPlayer:      Player{Symbol: SymbolPlayerMachine},
			expectedChildren:   []int{24, 43, 42, 41, 22, 3, 4, 5},
		},
		{
			name: "not all children free",
			playboard: "..................." +
				"....M1............." +
				"...1..............." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			index:              23,
			setChildrenIndexes: intSet{4: member, 5: member, 6: member, 23: member, 25: member, 22: member, 40: member, 42: member, 59: member, 60: member, 61: member},
			currentPlayer:      Player{Symbol: SymbolPlayerMachine},
			expectedChildren:   []int{43, 42, 22, 3, 4, 5, 6, 25, 40, 59, 60, 61},
		},
		{
			name: "is forbidden for put is not in children",
			playboard: "..................." +
				"....M1............." +
				"...1..............." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			index:              23,
			currentPlayer:      Player{Symbol: SymbolPlayerMachine},
			setChildrenIndexes: intSet{},
			expectedChildren:   []int{43, 42, 22, 3, 4, 5},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.playboard, func(t *testing.T) {
			machinePlayer := tc.currentPlayer
			humanPlayer := Player2
			UpdateSetChildren(tc.index, tc.playboard, tc.setChildrenIndexes)
			children := getChildren(
				State{Move{tc.playboard, tc.index, 0, []int{}}, machinePlayer, humanPlayer},
				tc.currentPlayer, tc.setChildrenIndexes)
			assert.Equal(t, len(tc.expectedChildren), len(children))
			for _, val := range tc.expectedChildren {
				found := false
				for _, state := range children {
					if state.move.index == val {
						found = true
						break
					}
				}
				assert.Equal(t, found, true, fmt.Sprintf("val is %d", val))
			}

		})
	}
}

func TestGetAllIndexChildren(t *testing.T) {
	testCases := []struct {
		name             string
		playBoard        string
		index            int
		currentPlayer    Player
		expectedChildren []int
	}{
		{
			name: "all 8 children",
			playBoard: "..................." +
				"....M.............." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			index:            23,
			currentPlayer:    Player{Symbol: SymbolPlayerMachine},
			expectedChildren: []int{24, 43, 42, 41, 22, 3, 4, 5},
		},
		{
			name: "all children for two stone",
			playBoard: "..................." +
				"....M1............." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			index:            23,
			currentPlayer:    Player{Symbol: SymbolPlayerMachine},
			expectedChildren: []int{43, 42, 41, 22, 3, 4, 5, 6, 25, 44},
		},
		{
			name: "is forbidden player 0 horizontal-vertical",
			playBoard: "..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				".........M........." +
				".........0........." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			index:            23,
			currentPlayer:    Player{Symbol: SymbolPlayerMachine},
			expectedChildren: []int{160, 161, 162, 179, 181, 198, 200, 217, 218, 219},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.playBoard, func(t *testing.T) {

			children := getAllIndexChildren(tc.playBoard)

			assert.Equal(t, len(tc.expectedChildren), len(children))
			for _, val := range tc.expectedChildren {
				_, found := children[val]
				assert.Equal(t, found, true, fmt.Sprintf("val is %d", val))
			}
		})
	}
}

func TestUpdateSetChildren(t *testing.T) {
	testCases := []struct {
		name             string
		playBoard        string
		index            int
		currentPlayer    Player
		expectedChildren []int
	}{
		{
			name: "near border",
			playBoard: "0.................." +
				"..................." +
				"0M................." +
				"..................." +
				"..................." +
				"..................." +
				"0.................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			index:            114,
			currentPlayer:    Player{Symbol: SymbolPlayerMachine},
			expectedChildren: []int{1, 19, 20, 21, 40, 57, 58, 59, 95, 96, 115, 133, 134},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.playBoard, func(t *testing.T) {

			setChildren := getAllIndexChildren(tc.playBoard)
			UpdateSetChildren(tc.index, tc.playBoard, setChildren)

			assert.Equal(t, len(tc.expectedChildren), len(setChildren))
			for _, val := range tc.expectedChildren {
				_, found := setChildren[val]
				assert.Equal(t, found, true, fmt.Sprintf("val is %d", val))
			}
		})
	}
}

func TestNegaScout(t *testing.T) {
	testCases := []struct {
		name            string
		playBoard       string
		index           int
		depth           int
		currentPlayer   Player
		expectedIndexes []int
		humanPlayer     *Player
	}{
		{
			name: "block 4",
			playBoard: "..................." +
				"..................." +
				"..................." +
				"........0.........." +
				"......M............" +
				"..................." +
				"......0M.000......." +
				".....0M0..0........" +
				".......00.0........" +
				"..........0........" +
				"..........M........" +
				".......M0000......." +
				".............0....." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			index:           114,
			depth:           10,
			currentPlayer:   Player{Symbol: SymbolPlayerMachine},
			expectedIndexes: []int{221},
		},
		{
			name: "test not 179",
			playBoard: "..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				".........0........." +
				"........0M........." +
				".........M0........" +
				"........0MM........" +
				".........M.M......." +
				".........0..0......" +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			index:           114,
			depth:           10,
			currentPlayer:   Player{Symbol: SymbolPlayerMachine},
			expectedIndexes: []int{236, 182},
		},
		{
			name: "block 5",
			playBoard: "..................." +
				"..................." +
				"..................." +
				"..................." +
				".........M........." +
				"....M000.0........." +
				"....0MMMM0........." +
				".....M...0M........" +
				"......0............" +
				"........MM........." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			depth:           10,
			currentPlayer:   Player{Symbol: SymbolPlayerMachine},
			expectedIndexes: []int{103},
		},
		{
			name: "make five",
			playBoard: "..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"......M............" +
				"..................." +
				"......0M..........." +
				"...M..M............" +
				".....M0.M.........." +
				"....M.0............" +
				"......00.MMM0......" +
				"......0..0........." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			depth:           10,
			currentPlayer:   Player{Symbol: SymbolPlayerMachine},
			expectedIndexes: []int{122, 212},
		},
		{
			name: "prevent capture",
			playBoard: "..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...........0......." +
				"..........M........" +
				".........M........." +
				"........M.........." +
				".......M..........." +
				"......0............" +
				"..................." +
				"..................." +
				"...............M..." +
				"...............M0.." +
				"...............0.0." +
				"..................0",
			depth:           10,
			currentPlayer:   Player{Symbol: SymbolPlayerMachine},
			expectedIndexes: []int{281},
		},
		{
			name: "test not -1",
			playBoard: "..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"........M......M..." +
				".........0....0...." +
				".........M0.00....." +
				".........M0000M...." +
				"........0M.0M......" +
				".........MM........" +
				"........M0........." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			depth:           10,
			currentPlayer:   Player{Symbol: SymbolPlayerMachine},
			expectedIndexes: []int{145},
		},
		{
			name: "don't capture",
			playBoard: "..................." +
				"..................." +
				"..................." +
				"..........M........" +
				".........1........." +
				"......M.11M........" +
				".......1..1........" +
				"......11..1........" +
				".....M.M1111M......" +
				".......1.MMM1......" +
				"........M1...1....." +
				".......M......M...." +
				"......M............" +
				".....1............." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			depth:           10,
			currentPlayer:   Player{Symbol: SymbolPlayerMachine},
			expectedIndexes: []int{123, 146}, //not 102
			humanPlayer:     &Player2,
		},
		{
			name: "block free three",
			playBoard: "..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"........M......M..." +
				".........0....0...." +
				".........M000M....." +
				".........M.0M......" +
				".........M.00......" +
				".........M...M....." +
				".........0........." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			depth:           10,
			currentPlayer:   Player{Symbol: SymbolPlayerMachine},
			expectedIndexes: []int{220},
		},
		{
			name: "block free three",
			playBoard: "..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"........M......M..." +
				".........0....0...." +
				".........M000M....." +
				".........M.0M......" +
				".........M.00......" +
				".........M...M....." +
				".........0........." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			depth:           10,
			currentPlayer:   Player{Symbol: SymbolPlayerMachine},
			expectedIndexes: []int{220},
		},
		{
			name: "make capture (test 4)",
			playBoard: "..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...........M......." +
				"............M......" +
				"........0M000M....." +
				".........0.0000M..." +
				"...M...0MMM.0......" +
				"....0...M..M0......" +
				".....0MMMM0.M......" +
				".....M0000M........" +
				".......M..........." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			depth:           10,
			currentPlayer:   Player{Symbol: SymbolPlayerMachine},
			expectedIndexes: []int{127},
		},
		{
			name: "don't block (test 6)",
			playBoard: "..................." +
				"..................." +
				"..................." +
				".........M........." +
				"........0.........." +
				".......0M.0........" +
				"......0.0.M.0......" +
				"......M00..M......." +
				"......0...M.M......" +
				".....M...M...0....." +
				".......0M.........." +
				".......0..........." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			depth:           10,
			currentPlayer:   Player{Symbol: SymbolPlayerMachine},
			expectedIndexes: []int{163},
		},
		{
			name: "block free three (test 6)",
			playBoard: "..................." +
				"..................." +
				"..................." +
				".........M........." +
				"........1.........." +
				".......1M.1........" +
				"......1.1.M.1......" +
				"....1M111..M......." +
				"......1...M.M......" +
				".....M...M...1....." +
				".......1M.........." +
				".......1..........." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			depth:           10,
			currentPlayer:   Player{Symbol: SymbolPlayerMachine},
			expectedIndexes: []int{101},
			humanPlayer:     &Player{Symbol: SymbolPlayer2, Captures: 2},
		},
		{
			name: "(test 8)",
			playBoard: "..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				".......M..........." +
				"........1.1........" +
				".......M111M......." +
				".........M1........" +
				".........M.M......." +
				".........M........." +
				".........M........." +
				".........1........." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			depth:           10,
			currentPlayer:   Player{Symbol: SymbolPlayerMachine},
			expectedIndexes: []int{200},
			humanPlayer:     &Player{Symbol: SymbolPlayer2, Captures: 2},
		},
		{
			name: "(game-history 10)",
			playBoard: "..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				".......0..........." +
				"......0M..........." +
				"......0M0.........." +
				"......0M.0........." +
				".........M0........" +
				"..........MM......." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			depth:           10,
			currentPlayer:   Player{Symbol: SymbolPlayerMachine},
			expectedIndexes: []int{101},
			humanPlayer:     &Player{Symbol: SymbolPlayer1, Captures: 1},
		},
		{
			name: "make three (game-history 10)",
			playBoard: "..................." +
				"..........0........" +
				".........M........." +
				"........M.........." +
				".....0.M..........." +
				"....M.M............" +
				"...M0000M.........." +
				"......M.0.........." +
				"........00........." +
				".......0..0........" +
				"......M..0.MMM0...." +
				"...........0......." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			depth:           10,
			currentPlayer:   Player{Symbol: SymbolPlayerMachine, Captures: 3},
			expectedIndexes: []int{135},
			humanPlayer:     &Player{Symbol: SymbolPlayer1, Captures: 4},
		},
		{
			name: "(test 11)",
			playBoard: "..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"........M1M........" +
				".........11........" +
				".........M1..M....." +
				"..........1M1...M.." +
				"........M11.1MM1..." +
				".......M..M.111M..." +
				"......M......M11..." +
				".....1........1M1.." +
				"..............1..M." +
				"..............M...." +
				"..................." +
				"..................." +
				"...................",
			depth:           10,
			currentPlayer:   Player{Symbol: SymbolPlayerMachine, Captures: 3},
			expectedIndexes: []int{104},
			humanPlayer:     &Player{Symbol: SymbolPlayer2, Captures: 4},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.playBoard, func(t *testing.T) {
			algo := Algo{tc.depth}
			if tc.humanPlayer == nil {
				tc.humanPlayer = &Player1
			}
			machinePlayer := tc.currentPlayer
			setChildren := getAllIndexChildren(tc.playBoard)
			_, index := algo.AlphaBeta(State{Move{tc.playBoard, -1, 0, []int{}}, machinePlayer, *tc.humanPlayer}, tc.depth, math.Inf(-1), math.Inf(1), 1, setChildren)

			assert.Contains(t, tc.expectedIndexes, index)
		})
	}
}

func TestHeuristic(t *testing.T) {
	testCases := []struct {
		name          string
		playboard     string
		index         int
		currentPlayer Player
		expectedNum   float64
	}{
		{
			name: "test heuristic 1",
			playboard: "0MMM..............." +
				"0.................." +
				"0.................." +
				"0.................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			currentPlayer: Player{Captures: 0, Symbol: SymbolPlayer1},
			index:         57,
			expectedNum:   1000000003,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.playboard, func(t *testing.T) {

			num := Heuristic(Move{Node: tc.playboard, index: tc.index}, tc.currentPlayer.Symbol, 0)

			assert.Equal(t, tc.expectedNum, num)
		})
	}
}
