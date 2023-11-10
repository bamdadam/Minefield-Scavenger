package game

import (
	"crypto/rand"
	"errors"
	"math"
	"math/big"
)

func convertBoard(b []int8, bLen int) (*Board, error) {
	if bLen*bLen != len(b) {
		return nil, errors.New("board should be equal to bLen*bLen")
	}
	br := make(Board, bLen)
	for i := 0; i < bLen; i++ {
		br[i] = make([]int8, bLen)
	}
	for i, v := range b {
		row := int(math.Floor(float64(i / bLen)))
		col := i % bLen
		br[row][col] = v
	}
	return &br, nil
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

func createBoard(bLen, bombPercent, numKeyShards int) (*Board, error) {
	br := make(Board, bLen)
	for i := 0; i < bLen; i++ {
		br[i] = make([]int8, bLen)
	}
	numBombs := bLen * bLen * bombPercent / 100
	for numBombs > 0 {
		rn, err := rand.Int(rand.Reader, big.NewInt(int64(bLen*bLen)))
		if err != nil {
			return nil, err
		}
		row := int(math.Floor(float64(rn.Int64() / int64(bLen))))
		col := rn.Int64() % int64(bLen)
		if br[row][col] != int8(bomb) {
			numBombs--
			br[row][col] = int8(bomb)
		}
	}
	for numKeyShards > 0 {
		rn, err := rand.Int(rand.Reader, big.NewInt(int64(bLen*bLen)))
		if err != nil {
			return nil, err
		}
		row := int(math.Floor(float64(rn.Int64() / int64(bLen))))
		col := rn.Int64() % int64(bLen)
		if br[row][col] != int8(bomb) && br[row][col] != int8(keyShard) {
			numKeyShards--
			br[row][col] = int8(keyShard)
		}
	}
	return &br, nil
}
