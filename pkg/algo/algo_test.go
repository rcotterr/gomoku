package algo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gomoku/pkg/playboard"
	"math"
	"testing"
)

func TestGetChildren(t *testing.T) {
	testCases := []struct {
		name               string
		playboard          string
		index              int
		setChildrenIndexes intSet
		currentPlayer      playboard.Player
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
			currentPlayer:      playboard.Player{0, playboard.SymbolPlayerMachine},
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
			currentPlayer:      playboard.Player{0, playboard.SymbolPlayerMachine},
			expectedChildren:   []int{43, 42, 22, 3, 4, 5, 6, 25, 40, 59, 60, 61},
		},
		//{
		//	name: "is forbidden for put is not in children",
		//	playboard: "..................." +
		//		"....M1............." +
		//		"...1..............." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"...................",
		//	index:            23,
		//	currentPlayer:    playboard.Player{0, playboard.SymbolPlayerMachine},
		//	expectedChildren: []int{43, 42, 22, 3, 4, 5},
		//},
	}
	for _, tc := range testCases {
		t.Run(tc.playboard, func(t *testing.T) {

			children := getChildren(tc.playboard, tc.index, tc.currentPlayer, tc.setChildrenIndexes)

			assert.Equal(t, len(tc.expectedChildren), len(children))
			for _, val := range tc.expectedChildren {
				_, found := children[val]
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
		currentPlayer    playboard.Player
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
			currentPlayer:    playboard.Player{0, playboard.SymbolPlayerMachine},
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
			currentPlayer:    playboard.Player{0, playboard.SymbolPlayerMachine},
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
			currentPlayer:    playboard.Player{0, playboard.SymbolPlayerMachine},
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
		currentPlayer    playboard.Player
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
			currentPlayer:    playboard.Player{0, playboard.SymbolPlayerMachine},
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
		name             string
		playBoard        string
		index            int
		depth            int
		currentPlayer    playboard.Player
		expectedChildren []int
	}{
		{
			name: "test not -1",
			playBoard: "..................." +
				"..................." +
				"..................." +
				"........0.........." +
				"......М............" +
				"..................." +
				"......0М.000......." +
				".....0М0..0........" +
				".......00.0........" +
				"..........0........" +
				"..........М........" +
				".......М0000......." +
				".............0....." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			index:            114,
			depth:            10,
			currentPlayer:    playboard.Player{0, playboard.SymbolPlayerMachine},
			expectedChildren: []int{1, 19, 20, 21, 40, 57, 58, 59, 95, 96, 115, 133, 134},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.playBoard, func(t *testing.T) {
			var transpositions = make(stringSet)
			allIndexesPath := ""
			index := -1
			humanPlayer := playboard.Player1
			machinePlayer := playboard.MachinePlayer
			setChildren := getAllIndexChildren(tc.playBoard)
			_, index = NegaScout(tc.playBoard, tc.depth, math.Inf(-1), math.Inf(1), 1, machinePlayer, humanPlayer, index, setChildren, transpositions, allIndexesPath)

			assert.True(t, index != -1)
		})
	}
}
