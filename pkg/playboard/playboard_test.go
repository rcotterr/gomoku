package playboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameOver(t *testing.T) {
	testCases := []struct {
		name           string
		playboard      string
		expectedIsOver bool
		index          int
		player1        *Player
		player2        *Player
	}{
		{
			name: "is over player 1",
			playboard: "11111.............." +
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
				"..................." +
				"...................",
			expectedIsOver: true,
		},
		{
			name: "is over player 0",
			playboard: "..................." +
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
				"..................." +
				"..............00000",
			index:          360,
			expectedIsOver: true,
		},
		{
			name: "is over horizontal",
			playboard: "..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"....00000.........." +
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
			index:          139,
			expectedIsOver: true,
		},
		{
			name: "is over vertical",
			playboard: "..................." +
				"..................." +
				".......1..........." +
				".......1..........." +
				".......1..........." +
				".......1..........." +
				".......1..........." +
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
			index:          45,
			expectedIsOver: true,
		},
		{
			name: "is over right diagonal",
			playboard: "..................." +
				"..................." +
				"..................." +
				"..................." +
				"......0............" +
				".......0..........." +
				"........0.........." +
				".........0........." +
				"..........0........" +
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
			index:          82,
			expectedIsOver: true,
		},
		{
			name: "is over left diagonal",
			playboard: "..................." +
				"..................." +
				"..................." +
				"..............1...." +
				".............1....." +
				"............1......" +
				"...........1......." +
				"..........1........" +
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
			index:          107,
			expectedIsOver: true,
		},
		{
			name: "is not over1",
			playboard: "..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...............1111" +
				"1.................." +
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
			index:          130,
			expectedIsOver: false,
		},
		{
			name: "is not over2",
			playboard: "..................." +
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
				"............1......" +
				"............1......" +
				"............1......" +
				"............11.....",
			index:          354,
			expectedIsOver: false,
		},
		{
			name: "is over by captures",
			playboard: "..................." +
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
				"..................." +
				"...................",
			expectedIsOver: true,
			index:          0,
			player1:        &Player{0, SymbolPlayer1},
			player2:        &Player{5, SymbolPlayer2},
		},
		{
			name: "is over by captures",
			playboard: "..................." +
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
				"..................." +
				"...................",
			expectedIsOver: true,
			index:          0,
			player1:        &Player{5, SymbolPlayer1},
			player2:        &Player{0, SymbolPlayer2},
		},
		{
			name: "is over by no space left",
			playboard: "1111111111111111111" +
				"1111111111111111111" +
				"1111111111111111111" +
				"1111111111111111111" +
				"1111111111111111111" +
				"1111111111111111111" +
				"1111111111111111111" +
				"1111111111111111111" +
				"1111111111111111111" +
				"1111111111111111111" +
				"0000000000000000000" +
				"0000000000000000000" +
				"0000000000000000000" +
				"0000000000000000000" +
				"0000000000000000000" +
				"0000000000000000000" +
				"0000000000000000000" +
				"0000000000000000000" +
				"0000000000000000000",
			index:          0,
			expectedIsOver: true,
		},
		{
			name: "is not over possible captured",
			playboard: "..................." +
				"..................." +
				".......1..........." +
				".......110........." +
				".......1..........." +
				".......1..........." +
				".......1..........." +
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
			index:          45,
			expectedIsOver: false,
		},
		{
			name: "is over",
			playboard: "..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"......0............" +
				".......M.0........." +
				".......MMM0........" +
				".......M.M0........" +
				"........MMM........" +
				".......M.M.0M......" +
				"......0..0M........" +
				"......M....M......." +
				"............0......" +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			index:          159,
			expectedIsOver: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.playboard, func(t *testing.T) {

			isOver := GameOver(tc.playboard, tc.player1, tc.player2, tc.index)

			assert.Equal(t, tc.expectedIsOver, isOver)
		})
	}
}

//Algo  16.0033ms
//current play board:
//0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18
//0  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//1  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//2  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//3  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//4  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//5......0............
//6.......M.0.........
//7.......MMM0........
//8.......M.M0........
//9........MMM........
//10 .......M.M.0M......
//11 ......0..0M........
//12 ......M....M.......
//13 ............0......
//14 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//15 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//16 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//17 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//18 .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
//test update 2022-08-07 02:44:54.2089829 +0300 MSK m=+43.783781401

func TestIsCapture(t *testing.T) {
	testCases := []struct {
		name            string
		playboard       string
		index           int
		currentPlayer   string
		numCaptures     int
		expectedIndexes intSet
	}{
		{
			name: "is capture player 0",
			playboard: "0110.............." +
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
				"..................." +
				"...................",
			index:           0,
			currentPlayer:   SymbolPlayer1,
			numCaptures:     1,
			expectedIndexes: intSet{2: member, 1: member},
		},
		{
			name: "is capture player 1",
			playboard: "..................." +
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
				"..................." +
				"...............1001",

			index:           357,
			currentPlayer:   SymbolPlayer2,
			numCaptures:     1,
			expectedIndexes: intSet{358: member, 359: member},
		},
		{
			name: "is capture right diagonal",
			playboard: "..................." +
				"...0..............." +
				"....1.............." +
				".....1............." +
				"......0............" +
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

			index:           22,
			currentPlayer:   SymbolPlayer1,
			numCaptures:     1,
			expectedIndexes: intSet{42: member, 62: member},
		},
		{
			name: "is capture left diagonal",
			playboard: "..................." +
				".........0........." +
				"........1.........." +
				".......1..........." +
				"......0............" +
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

			index:           28,
			currentPlayer:   SymbolPlayer1,
			numCaptures:     1,
			expectedIndexes: intSet{46: member, 64: member},
		},
		{
			name: "is capture vertical",
			playboard: "..................." +
				".........1........." +
				".........00........" +
				".........0.0......." +
				".........1..1......" +
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

			index:           28,
			currentPlayer:   SymbolPlayer2,
			numCaptures:     2,
			expectedIndexes: intSet{47: member, 66: member, 48: member, 68: member},
		},
		{
			name: "is capture horizontal left",
			playboard: "..................." +
				"......1001........." +
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

			index:           28,
			currentPlayer:   SymbolPlayer2,
			numCaptures:     1,
			expectedIndexes: intSet{26: member, 27: member},
		},
		{
			name: "is capture vertical upper",
			playboard: "..................." +
				".........1........." +
				".........0........." +
				".........0........." +
				".........1........." +
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

			index:           85,
			currentPlayer:   SymbolPlayer2,
			numCaptures:     1,
			expectedIndexes: intSet{47: member, 66: member},
		},
		{
			name: "is capture right diagonal upper",
			playboard: "..................." +
				"...0..............." +
				"....1.............." +
				".....1............." +
				"......0............" +
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

			index:           82,
			currentPlayer:   SymbolPlayer1,
			numCaptures:     1,
			expectedIndexes: intSet{42: member, 62: member},
		},
		{
			name: "is capture left diagonal upper",
			playboard: "..................." +
				".........0........." +
				"........1.........." +
				".......1..........." +
				"......0............" +
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

			index:           82,
			currentPlayer:   SymbolPlayer1,
			numCaptures:     1,
			expectedIndexes: intSet{46: member, 64: member},
		},

		{
			name: "is not capture right diagonal",
			playboard: "..................." +
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
				"..............1...." +
				"...............0..." +
				"................0..",

			index:           318,
			currentPlayer:   SymbolPlayer2,
			numCaptures:     0,
			expectedIndexes: intSet{},
		},
		{
			name: "is not capture full",
			playboard: "001................" +
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
				"..................." +
				"...................",

			index:           2,
			currentPlayer:   SymbolPlayer2,
			numCaptures:     0,
			expectedIndexes: intSet{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.playboard, func(t *testing.T) {

			numCaptures, arrIndexes := isCaptured(tc.playboard, tc.index, tc.currentPlayer)

			assert.Equal(t, tc.numCaptures, numCaptures)
			//b := reflect.DeepEqual(tc.expectedIndexes, arrIndexes)
			indexSet := make(intSet)
			for _, elem := range arrIndexes {
				indexSet[elem] = member
			}
			//reflect.DeepEqual(tc.expectedIndexes, indexSet)
			assert.Equal(t, tc.expectedIndexes, indexSet)
			//assert.Equal(t, true, b)
			//assert.Equal(t, true, b)
			//if isCapture {
			//	assert.Equal(t, tc.expectedIndex1, *index1)
			//	assert.Equal(t, tc.expectedIndex2, *index2)
			//}
		})
	}
}

//TO DO name not playBoard in fail

func TestPutStone(t *testing.T) {
	testCases := []struct {
		name                string
		playboard           string
		pos                 *Pos
		currentPlayer       Player
		expectedNewSymbol   map[int]string
		expectedNumCaptures int
		expectedError       error
	}{
		{
			name: "is capture player 0",
			playboard: ".110..............." +
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
				"..................." +
				"...................",
			pos:           &Pos{0, 0},
			currentPlayer: Player{0, SymbolPlayer1},
			expectedNewSymbol: map[int]string{
				0: "0",
				1: ".",
				2: ".",
			},
			expectedNumCaptures: 1,
			expectedError:       nil,
		},
		{
			name: "two captures player 0",
			playboard: ".110..............." +
				"1.................." +
				"1.................." +
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
			pos:           &Pos{0, 0},
			currentPlayer: Player{0, SymbolPlayer1},
			expectedNewSymbol: map[int]string{
				0:  "0",
				1:  ".",
				2:  ".",
				19: ".",
				38: ".",
			},
			expectedNumCaptures: 2,
			expectedError:       nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.playboard, func(t *testing.T) {
			index := tc.pos.Y*N + tc.pos.X

			state, err := PutStone(tc.playboard, index, &tc.currentPlayer)

			for index, symbol := range tc.expectedNewSymbol {
				assert.Equal(t, symbol, string(state.Node[index]))
			}
			assert.Equal(t, tc.expectedNumCaptures, tc.currentPlayer.Captures)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestIsForbidden(t *testing.T) {
	testCases := []struct {
		name              string
		playboard         string
		index             int
		currentPlayer     Player
		expectedForbidden bool
	}{
		{
			name: "is forbidden player 0 horizontal-vertical",
			playboard: "..................." +
				"....000............" +
				"....0.............." +
				"....0.............." +
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
			currentPlayer:     Player{0, SymbolPlayer1},
			index:             23,
			expectedForbidden: true,
		},
		{
			name: "is forbidden right diagonal-vertical",
			playboard: "..................." +
				"......1............" +
				"......1............" +
				"......11..........." +
				"........1.........." +
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
			currentPlayer:     Player{0, SymbolPlayer2},
			index:             44,
			expectedForbidden: true,
		},
		{
			name: "is forbidden right-left diagonals",
			playboard: "..................." +
				".......1..........." +
				"......1............" +
				".....1.1..........." +
				"........1.........." +
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
			currentPlayer:     Player{0, SymbolPlayer2},
			index:             44,
			expectedForbidden: true,
		},
		{
			name: "is forbidden right-left diagonals player 0",
			playboard: "..................." +
				".....0............." +
				"......0............" +
				"..................." +
				"....0...0.........." +
				"...0..............." +
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
			currentPlayer:     Player{0, SymbolPlayer1},
			index:             44,
			expectedForbidden: true,
		},
		{
			name: "is not forbidden player 1",
			playboard: "..................." +
				"..................1" +
				"...............111." +
				"................1.." +
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
			currentPlayer:     Player{0, SymbolPlayer2},
			index:             55,
			expectedForbidden: false,
		},
		{
			name: "is not forbidden player 0",
			playboard: "..................." +
				"..................." +
				"..................." +
				"........0.........." +
				"......M............" +
				"..................." +
				"......0M.000......." +
				".....0M0..0........" +
				".......00.0........" +
				"..........0........" +
				"..........MM......." +
				".......M0000M......" +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"..................." +
				"...................",
			currentPlayer:     Player{0, SymbolPlayer1},
			index:             105,
			expectedForbidden: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.playboard, func(t *testing.T) {

			isForbiddenPlace := isForbidden(tc.playboard, tc.index, tc.currentPlayer.Symbol)

			assert.Equal(t, tc.expectedForbidden, isForbiddenPlace)
		})
	}
}

func TestPossibleCapturedStone(t *testing.T) {
	testCases := []struct {
		name                    string
		playboard               string
		index                   int
		currentPlayer           string
		stepCount               int
		expectedPossibleCapture int
	}{
		{
			name: "is possible captured player 1",
			playboard: ".110.............." +
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
				"..................." +
				"...................",
			index:                   1,
			currentPlayer:           SymbolPlayer2,
			stepCount:               N,
			expectedPossibleCapture: 1,
		},
		{
			name: "is possible captured player 0",
			playboard: "..................." +
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
				"..................." +
				"................001",

			index:                   358,
			currentPlayer:           SymbolPlayer1,
			stepCount:               N,
			expectedPossibleCapture: 1,
		},
		{
			name: "is possible captured right diagonal",
			playboard: "..................." +
				"...0..............." +
				"....1.............." +
				".....1............." +
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

			index:                   42,
			currentPlayer:           SymbolPlayer2,
			stepCount:               N,
			expectedPossibleCapture: 1,
		},
		{
			name: "is possible captured left diagonal",
			playboard: "..................." +
				"..................." +
				"........1.........." +
				".......1..........." +
				"......0............" +
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

			index:                   46,
			currentPlayer:           SymbolPlayer2,
			stepCount:               N,
			expectedPossibleCapture: 1,
		},
		{
			name: "is possible captured vertical",
			playboard: "..................." +
				"..................." +
				".........00........" +
				".........0.0......." +
				".........1..1......" +
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

			index:                   47,
			currentPlayer:           SymbolPlayer1,
			stepCount:               1,
			expectedPossibleCapture: 1,
		},
		{
			name: "is possible captured horizontal left",
			playboard: "..................." +
				".......001........." +
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

			index:                   27,
			currentPlayer:           SymbolPlayer1,
			stepCount:               N,
			expectedPossibleCapture: 1,
		},
		{
			name: "is possible captured vertical upper",
			playboard: "..................." +
				".........1........." +
				".........0........." +
				".........0........." +
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

			index:                   66,
			currentPlayer:           SymbolPlayer1,
			stepCount:               1,
			expectedPossibleCapture: 1,
		},
		{
			name: "is possible captured right diagonal upper",
			playboard: "..................." +
				"...0..............." +
				"....1.............." +
				".....1............." +
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

			index:                   62,
			currentPlayer:           SymbolPlayer2,
			stepCount:               1,
			expectedPossibleCapture: 1,
		},
		{
			name: "is possible captured left diagonal upper",
			playboard: "..................." +
				"..................." +
				"........1.........." +
				".......1..........." +
				"......0............" +
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

			index:                   64,
			currentPlayer:           SymbolPlayer2,
			stepCount:               1,
			expectedPossibleCapture: 1,
		},
		{
			name: "is not possible captured right diagonal",
			playboard: "..................." +
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
				"..............1...." +
				"...............0..." +
				"................0..",

			index:                   338,
			currentPlayer:           SymbolPlayer1,
			stepCount:               1,
			expectedPossibleCapture: 0,
		},
		{
			name: "is not possible captured full",
			playboard: "001................" +
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
				"..................." +
				"...................",

			index:                   2,
			currentPlayer:           SymbolPlayer1,
			stepCount:               N,
			expectedPossibleCapture: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.playboard, func(t *testing.T) {

			possibleCaptures := PossibleCapturedStone(tc.playboard, tc.index, tc.stepCount, tc.currentPlayer)

			assert.Equal(t, tc.expectedPossibleCapture, possibleCaptures)
		})
	}
}
