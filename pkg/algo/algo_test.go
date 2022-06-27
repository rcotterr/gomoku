package algo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gomoku/pkg/playboard"
	"testing"
)

func TestGetChildren(t *testing.T) {
	testCases := []struct {
		name             string
		playboard        string
		index            int
		currentPlayer    playboard.Player
		expectedChildren []int
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
			index:            23,
			currentPlayer:    playboard.Player{0, playboard.SymbolPlayerMachine},
			expectedChildren: []int{24, 43, 42, 41, 22, 3, 4, 5},
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
			index:            23,
			currentPlayer:    playboard.Player{0, playboard.SymbolPlayerMachine},
			expectedChildren: []int{43, 42, 22, 3, 4, 5},
		},
		{
			name: "is forbidden player 0 horizontal-vertical",
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
			index:            23,
			currentPlayer:    playboard.Player{0, playboard.SymbolPlayerMachine},
			expectedChildren: []int{43, 42, 22, 3, 4, 5},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.playboard, func(t *testing.T) {

			children := getChildren(tc.playboard, tc.index, tc.currentPlayer)

			assert.Equal(t, len(tc.expectedChildren), len(children))
			for _, val := range tc.expectedChildren {
				_, found := children[val]
				assert.Equal(t, found, true, fmt.Sprintf("val is %d", val))
			}
		})
	}
}
