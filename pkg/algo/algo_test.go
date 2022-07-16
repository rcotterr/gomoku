package algo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gomoku/pkg/playboard"
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
