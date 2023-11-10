package game

import (
	"errors"
	"fmt"
	"math"
)

type Game struct {
	numKeyShards int
	fieldLen     int
	bombPercent  int
	board        *Board
	seen         *Seen
	seenCounter  int
}

type CompactGameModel struct {
	KeyShards   int     `json:"key_shards"`
	FieldLen    int     `json:"field_len"`
	BombPercent int     `json:"bomb_percent"`
	Bombs       []int16 `json:"bombs"`
	Keys        []int16 `json:"keys"`
	Empty       []int16 `json:"empty"`
	NotSeen     []int16 `json:"not_seen"`
}

type Board [][]int8
type Seen [][]bool
type cell int8

const (
	empty cell = iota
	keyShard
	bomb
)

// custom string function for board
func (b Board) String(s Seen) string {
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

func (g *Game) ShowState() *CompactGameModel {
	var bombs, keys, seenEmpty, notSeen []int16
	for i, value := range *g.seen {
		for j, iv := range value {
			if iv {
				idx := i*g.fieldLen + j
				switch (*g.board)[i][j] {
				case int8(bomb):
					bombs = append(bombs, int16(idx))
				case int8(keyShard):
					keys = append(keys, int16(idx))
				case int8(empty):
					if g.seenCounter < int(math.Floor(float64(g.fieldLen*g.fieldLen/2))) {
						seenEmpty = append(seenEmpty, int16(idx))
					}
				}
			} else {
				if g.seenCounter >= int(math.Floor(float64(g.fieldLen*g.fieldLen/2))) {
					idx := i*g.fieldLen + j
					notSeen = append(notSeen, int16(idx))
				}
			}
		}
	}

	return &CompactGameModel{
		KeyShards:   g.numKeyShards,
		FieldLen:    g.fieldLen,
		BombPercent: g.bombPercent,
		Bombs:       bombs,
		Keys:        keys,
		Empty:       seenEmpty,
		NotSeen:     notSeen,
	}
}

func NewGame(keyShards, fieldLen, bombPercent int) (*Game, error) {
	if keyShards >= fieldLen*fieldLen {
		return nil, errors.New("(fieldLength)^2 should at least be bigger than number of key shards")
	}
	br, err := createBoard(fieldLen, bombPercent, keyShards)
	if err != nil {
		return nil, err
	}
	sr := make(Seen, fieldLen)
	for i := 0; i < fieldLen; i++ {
		sr[i] = make([]bool, fieldLen)
	}
	return &Game{
		numKeyShards: keyShards,
		fieldLen:     fieldLen,
		bombPercent:  bombPercent,
		board:        br,
		seen:         &sr,
		seenCounter:  0,
	}, nil
}

func (g *Game) MakeMove(row, col int) (int8, error) {
	if row >= g.fieldLen || col >= g.fieldLen || row < 0 || col < 0 {
		return 0, errors.New("row or col is out of range")
	}
	fmt.Println((*g.seen)[row][col])
	if (*g.seen)[row][col] {
		fmt.Println("test")
		return 0, errors.New("this cell is already seen")
	}
	(*g.seen)[row][col] = true
	g.seenCounter++
	return (*g.board)[row][col], nil
}

func LoadGame(keyShards, fieldLen, bombPercent, seenCounter int, board []int8, seen []bool) (*Game, error) {
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
		seenCounter:  seenCounter,
	}, nil
}
