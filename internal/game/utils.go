package game

import (
	"errors"
	"math"
)

func convertBoard(b [][]int8) (*Board, error) {
	// if bLen*bLen != len(b) {
	// 	return nil, errors.New("board should be equal to bLen*bLen")
	// }
	// br := make(Board, bLen)
	// for i := 0; i < bLen; i++ {
	// 	br[i] = make([]int8, bLen)
	// }
	// for i, v := range b {
	// 	row := int(math.Floor(float64(i / bLen)))
	// 	col := i % bLen
	// 	br[row][col] = v
	// }
	// return &br, nil
	board := Board(b)
	return &board, nil
}

func converSeenBoard(s []bool, bLen int) (*Seen, error) {
	if bLen*bLen != len(s) {
		return nil, errors.New("see board should be equal to bLen*bLen")
	}
	sr := make(Seen, bLen)
	for i := 0; i < bLen; i++ {
		sr[i] = make([]bool, bLen)
	}
	for i, v := range s {
		row := int(math.Floor(float64(i / bLen)))
		col := i % bLen
		sr[row][col] = v
	}
	return &sr, nil
}

func isBoardAndSeenEqual(board [][]int8, seen [][]bool) bool {
	if len(board) != len(seen) {
		return false
	}
	for i := 0; i < len(board); i++ {
		if len(board[i]) != len(seen[i]) {
			return false
		}
	}
	return true
}

func calcSeenCounter(s Seen) int {
	counter := 0
	for _, v := range s {
		for _, iv := range v {
			if iv {
				counter++
			}
		}
	}
	return counter
}
