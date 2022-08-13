package playboard

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

//func TestGetChildren(t *testing.T) {
//	testCases := []struct {
//		name               string
//		playboard          string
//		index              int
//		setChildrenIndexes intSet
//		currentPlayer      Player
//		expectedChildren   []int
//	}{
//		{
//			name: "all 8 children",
//			playboard: "..................." +
//				"....M.............." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"...................",
//			index:              23,
//			setChildrenIndexes: intSet{},
//			currentPlayer:      Player{0, SymbolPlayerMachine},
//			expectedChildren:   []int{24, 43, 42, 41, 22, 3, 4, 5},
//		},
//		{
//			name: "not all children free",
//			playboard: "..................." +
//				"....M1............." +
//				"...1..............." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"..................." +
//				"...................",
//			index:              23,
//			setChildrenIndexes: intSet{4: member, 5: member, 6: member, 23: member, 25: member, 22: member, 40: member, 42: member, 59: member, 60: member, 61: member},
//			currentPlayer:      Player{0, SymbolPlayerMachine},
//			expectedChildren:   []int{43, 42, 22, 3, 4, 5, 6, 25, 40, 59, 60, 61},
//		},
//		//{
//		//	name: "is forbidden for put is not in children",
//		//	playboard: "..................." +
//		//		"....M1............." +
//		//		"...1..............." +
//		//		"..................." +
//		//		"..................." +
//		//		"..................." +
//		//		"..................." +
//		//		"..................." +
//		//		"..................." +
//		//		"..................." +
//		//		"..................." +
//		//		"..................." +
//		//		"..................." +
//		//		"..................." +
//		//		"..................." +
//		//		"..................." +
//		//		"..................." +
//		//		"..................." +
//		//		"...................",
//		//	index:            23,
//		//	currentPlayer:    playboard.Player{0, playboard.SymbolPlayerMachine},
//		//	expectedChildren: []int{43, 42, 22, 3, 4, 5},
//		//},
//	}
//	for _, tc := range testCases {
//		t.Run(tc.playboard, func(t *testing.T) {
//
//			children := getChildren(tc.playboard, tc.index, tc.currentPlayer, tc.setChildrenIndexes)
//
//			assert.Equal(t, len(tc.expectedChildren), len(children))
//			for _, val := range tc.expectedChildren {
//				_, found := children[val]
//				assert.Equal(t, found, true, fmt.Sprintf("val is %d", val))
//			}
//
//		})
//	}
//}

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
		name             string
		playBoard        string
		index            int
		depth            int
		currentPlayer    Player
		expectedChildren []int
		notExpectedIndex int
	}{
		//{
		//	name: "test not -1",
		//	playBoard: "..................." +
		//		"..................." +
		//		"..................." +
		//		"........0.........." +
		//		"......М............" +
		//		"..................." +
		//		"......0М.000......." +
		//		".....0М0..0........" +
		//		".......00.0........" +
		//		"..........0........" +
		//		"..........М........" +
		//		".......М0000......." +
		//		".............0....." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"...................",
		//	index:            114,
		//	depth:            10,
		//	currentPlayer:    Player{0, SymbolPlayerMachine},
		//	expectedChildren: []int{1, 19, 20, 21, 40, 57, 58, 59, 95, 96, 115, 133, 134},
		//	notExpectedIndex: -1,
		//},
		//{
		//	name: "test not 179",
		//	playBoard: "..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		".........0........." +
		//		"........0M........." +
		//		".........M0........" +
		//		"........0MM........" +
		//		".........M.M......." +
		//		".........0..0......" +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"...................",
		//	index:            114,
		//	depth:            10,
		//	currentPlayer:    Player{0, SymbolPlayerMachine},
		//	expectedChildren: []int{1, 19, 20, 21, 40, 57, 58, 59, 95, 96, 115, 133, 134},
		//	notExpectedIndex: 179,
		//},
		//{
		//	name: "block 5",
		//	playBoard: "..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		".........M........." +
		//		"....M000.0........." +
		//		"....0MMMM0........." +
		//		".....M...0M........" +
		//		"......0............" +
		//		"........MM........." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"...................",
		//	depth:            10,
		//	currentPlayer:    Player{0, SymbolPlayerMachine},
		//	notExpectedIndex: 181, //expected 103
		//},
		//{
		//	name: "not block 5",
		//	playBoard: "..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"......M............" +
		//		"..................." +
		//		"......0M..........." +
		//		"...M..M............" +
		//		".....M0.M.........." +
		//		"....M.0............" +
		//		"......00.MMM0......" +
		//		"......0..0........." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"..................." +
		//		"...................",
		//	depth:            10,
		//	currentPlayer:    Player{0, SymbolPlayerMachine},
		//	notExpectedIndex: 254, //expected 122 or 212
		//},
		{
			name: "speed check",
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
			depth:            10,
			currentPlayer:    Player{Symbol: SymbolPlayerMachine},
			notExpectedIndex: 339, // expected 281
		},
	}
	for _, tc := range testCases {
		t.Run(tc.playBoard, func(t *testing.T) {
			var transpositions = make(stringSet)
			allIndexesPath := ""
			index := -1
			humanPlayer := Player1
			machinePlayer := MachinePlayer
			setChildren := getAllIndexChildren(tc.playBoard)
			_, index = NegaScout(State{tc.playBoard, index, 0}, tc.depth, math.Inf(-1), math.Inf(1), 1, machinePlayer, humanPlayer, setChildren, transpositions, allIndexesPath)

			assert.NotEqual(t, index, tc.notExpectedIndex)
		})
	}
}

//Algo  1.8626471s
//current play board:
//   0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18
//0  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//1  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//2  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//3  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//4  .  .  .  .  .  M  .  .  .  .  .  .  .  .  .  .  .  .  .
//5  .  .  .  .  .  0  .  .  .  .  .  .  .  .  .  .  .  .  .
//6  .  .  .  .  .  0  .  .  .  .  .  .  0  .  .  .  .  .  .
//7  .  .  .  .  .  0  .  .  .  .  .  M  .  .  .  .  .  .  .
//8  .  .  .  .  .  ?  .  .  .  .  M  .  .  .  .  .  .  .  .
//9  .  .  .  .  .  .  .  .  .  M  .  .  .  .  .  .  .  .  .
//10 .  .  .  .  .  .  .  .  M  .  .  .  .  .  .  .  .  .  .
//11 .  .  .  .  .  .  .  0  .  .  .  .  .  .  .  .  .  .  .
//12 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//13 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//14 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//15 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//16 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//17 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//18 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .

//Algo  1.2841888s
//current play board:
//0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18
//0  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//1  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//2  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//3  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//4  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//5  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//6  .  .  .  .  .  .  .  .  .  .  .  .  0  .  .  .  .  .  .
//7  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//8  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//9  .  .  .  .  .  .  .  .  .  M  .  .  .  .  .  .  .  .  .
//10 .  .  .  .  .  .  .  .  .  .  ?  .  .  .  .  .  .  .  .
//11 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//12 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//13 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//14 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//15 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//16 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//17 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//18 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .

//   0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18
//0  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//1  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//2  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//3  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//4  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//5  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//6  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//7  .  .  .  .  .  .  .  .  .  .  .  0  .  .  .  .  .  .  .
//8  .  .  .  .  .  .  .  .  .  .  M  .  .  .  .  .  .  .  .
//9  .  .  .  .  .  .  .  .  .  M  .  .  .  .  .  .  .  .  .
//10 .  .  .  .  .  .  .  .  M  .  .  .  .  .  .  .  .  .  .
//11 .  .  .  .  .  .  .  M  .  .  .  .  .  .  .  .  .  .  .
//12 .  .  .  .  .  .  0  .  .  .  .  .  .  .  .  .  .  .  .
//13 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//14 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//15 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  M  .  .  .
//16 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  M  0  .  .
//17 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  0  ?  0  .
//18 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  0

//
//current play board:
//0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18
//0  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//1  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//2  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//3  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//4  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//5  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//6  .  .  .  .  .  .  .  .  M  .  .  .  .  .  M  .  .  .  .
//7  .  .  .  .  .  .  .  .  .  0  .  .  .  0  ?  .  .  .  . // I'm a winner
//8  .  .  .  .  .  .  .  .  .  M  0  0  0  M  .  .  .  .  .
//9  .  .  .  .  .  .  .  .  .  M  .  0  M  .  .  .  .  .  .
//10 .  .  .  .  .  .  .  .  .  M  .  0  0  .  .  .  .  .  .
//11 .  .  .  .  .  .  .  .  .  M  .  .  .  M  .  .  .  .  .
//12 .  .  .  .  .  .  .  .  .  0  .  .  .  .  .  .  .  .  .
//13 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//14 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//15 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//16 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//17 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//18 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//test update 2022-08-07 19:21:15.6544751 +0300 MSK m=+32.269841001

//   0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18
//0  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//1  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//2  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//3  .  .  .  .  .  .  .  0  0  ?  .  .  .  .  .  .  .  .  .
//4  .  .  .  .  .  .  .  .  M  .  .  .  .  .  .  .  .  .  .
//5  .  .  .  .  .  .  .  .  M  M  .  .  .  .  .  .  .  .  .
//6  .  .  .  .  .  .  .  .  M  0  M  0  .  .  .  .  .  .  .
//7  .  .  .  .  .  .  .  .  M  0  .  .  .  .  .  .  .  .  .
//8  .  .  .  .  .  .  .  .  0  0  0  M  0  .  .  .  .  .  .
//9  .  .  .  .  .  .  0  M  M  M  .  0  M  .  .  .  .  .  .
//10 .  .  .  .  .  .  .  .  .  .  .  .  0  0  .  .  .  .  .
//11 .  .  .  .  .  .  .  .  .  .  0  .  .  M  .  .  .  .  .
//12 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//13 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//14 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//15 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//16 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//17 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//18 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .

//   0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18
//0  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//1  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//2  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//3  .  .  .  .  .  .  .  0  0  M  .  .  .  .  .  .  .  .  .
//4  .  .  .  .  .  .  .  .  M  .  .  .  .  .  .  .  .  .  .
//5  .  .  .  .  .  .  .  .  M  M  .  .  .  .  .  .  .  .  .
//6  .  .  .  .  .  .  .  .  M  .  M  0  .  .  .  ?  .  .  .
//7  .  .  .  .  .  .  .  .  M  0  0  .  .  .  0  .  .  .  .
//8  .  .  .  .  .  .  .  .  0  0  0  M  0  0  .  .  .  .  .
//9  .  .  .  .  .  .  0  M  M  M  .  0  0  .  .  .  .  .  .
//10 .  .  .  .  .  .  .  .  .  .  .  0  0  0  M  .  .  .  .
//11 .  .  .  .  .  .  .  .  .  .  0  .  .  M  .  .  .  .  .
//12 .  .  .  .  .  .  .  .  .  M  .  .  .  .  .  .  .  .  .
//13 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//14 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//15 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//16 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//17 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//18 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
