package game

import (
	"errors"
	"fmt"
)

type Game struct {
	numKeyShards int
	fieldLen     int
	bombPercent  int
	board        *board
	seen         *seen
}

type board [][]int8
type seen [][]bool
type cell int8

const (
	empty cell = iota
	keyShard
	bomb
)

// custom string function for board
func (b board) String(s seen) string {
	str := ""
	for i, row := range b {
		for j, col := range row {
			if s[i][j] {
				// if it is a bomb show bomb if not show empty house
				if col == 2 {
					// bomb icon not ascii
					str += "⊗ "
				} else if col == 0 {
					// empty square icon not ascii
					str += "□ "
				} else if col == 1 {
					// key icon
					str += "k "
				}
			} else {
				// if it is not seen show unknown
				str += "■ "
			}
		}
		str += "\n"
	}
	return str
}

// custom string function to show the whole game
func (g *Game) String() string {
	return fmt.Sprintf("|key shards: %v|board length: %v|bomb percentage %v%%|\nboard: \n%s \n", g.numKeyShards, g.fieldLen, g.bombPercent, g.board.String(*g.seen))
}

func NewGame(keyShards, fieldLen, bombPercent int) (*Game, error) {
	if keyShards >= fieldLen*fieldLen {
		return nil, errors.New("(fieldLength)^2 should at least be bigger than number of key shards")
	}
	br, err := createBoard(fieldLen, bombPercent, keyShards)
	if err != nil {
		return nil, err
	}
	sr := make(seen, fieldLen)
	for i := 0; i < fieldLen; i++ {
		sr[i] = make([]bool, fieldLen)
	}
	return &Game{
		numKeyShards: keyShards,
		fieldLen:     fieldLen,
		bombPercent:  bombPercent,
		board:        br,
		seen:         &sr,
	}, nil
}

func (g *Game) MakeMove(row, col int) (int8, error) {
	if row >= g.fieldLen || col >= g.fieldLen {
		return 0, errors.New("row or col is out of range")
	}
	if (*g.seen)[row][col] {
		return 0, errors.New("this cell is already seen")
	}
	(*g.seen)[row][col] = true
	return (*g.board)[row][col], nil
}

func LoadGame(keyShards, fieldLen, bombPercent int, board []int8, seen []bool) (*Game, error) {
	if keyShards >= fieldLen*fieldLen {
		return nil, errors.New("(fieldLength)^2 should at least be bigger than number of key shards")
	}
	if len(board) != len(seen) {
		return nil, errors.New("board and seen slices should be of the same length")
	}
	if len(board) != fieldLen*fieldLen {
		return nil, errors.New("board length should be the same az fieldLen * fieldLen")
	}
	b, err := convertBoard(board, fieldLen)
	if err != nil {
		return nil, err
	}
	s, err := converSeenBoard(seen, fieldLen)
	if err != nil {
		return nil, err
	}
	return &Game{
		numKeyShards: keyShards,
		fieldLen:     fieldLen,
		bombPercent:  bombPercent,
		board:        b,
		seen:         s,
	}, nil
}
