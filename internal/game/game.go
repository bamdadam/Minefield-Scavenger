package game

import (
	"errors"
	"fmt"
	"math"

	"github.com/bamdadam/Minefield-Scavenger/internal/model"
)

type Game struct {
	NumKeyShards int
	FieldLen     int
	BombPercent  int
	GameBoard    *Board
	Seen         *Seen
	SeenCounter  int
	GameId       int
}

type CompactGameModel struct {
	KeyShards        int     `json:"key_shards"`
	FieldLen         int     `json:"field_len"`
	BombPercent      int     `json:"bomb_percent"`
	Bombs            []int16 `json:"bombs"`
	Keys             []int16 `json:"keys"`
	Empty            []int16 `json:"empty"`
	NotSeen          []int16 `json:"not_seen"`
	IsDefaultNotSeen bool    `json:"is_default_not_seen"`
}

type Board [][]int8
type Seen [][]bool
type cell int8

const (
	Empty cell = iota
	KeyShard
	Bomb
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
	return fmt.Sprintf("|key shards: %v|board length: %v|bomb percentage %v%%|\nboard: \n%s \n", g.NumKeyShards, g.FieldLen, g.BombPercent, g.GameBoard.String(*g.Seen))
}

func (g *Game) ShowState() *CompactGameModel {
	var bombs, keys, seenEmpty, notSeen []int16
	for i, value := range *g.Seen {
		for j, iv := range value {
			if iv {
				idx := i*g.FieldLen + j
				switch (*g.GameBoard)[i][j] {
				case int8(Bomb):
					bombs = append(bombs, int16(idx))
				case int8(KeyShard):
					keys = append(keys, int16(idx))
				case int8(Empty):
					if g.SeenCounter < int(math.Floor(float64(g.FieldLen*g.FieldLen/2))) {
						seenEmpty = append(seenEmpty, int16(idx))
					}
				}
			} else {
				if g.SeenCounter >= int(math.Floor(float64(g.FieldLen*g.FieldLen/2))) {
					idx := i*g.FieldLen + j
					notSeen = append(notSeen, int16(idx))
				}
			}
		}
	}
	def := false
	if g.SeenCounter < int(math.Floor(float64(g.FieldLen*g.FieldLen/2))) {
		def = true
	}

	return &CompactGameModel{
		KeyShards:        g.NumKeyShards,
		FieldLen:         g.FieldLen,
		BombPercent:      g.BombPercent,
		Bombs:            bombs,
		Keys:             keys,
		Empty:            seenEmpty,
		NotSeen:          notSeen,
		IsDefaultNotSeen: def,
	}
}

func (g *Game) ShowStateOnLose() *CompactGameModel {
	var bombs, keys, seenEmpty, notSeen []int16
	for i, value := range *g.Seen {
		for j := range value {
			idx := i*g.FieldLen + j
			switch (*g.GameBoard)[i][j] {
			case int8(Bomb):
				bombs = append(bombs, int16(idx))
			case int8(KeyShard):
				keys = append(keys, int16(idx))
			case int8(Empty):
				seenEmpty = append(seenEmpty, int16(idx))
			}
		}
	}
	def := true

	return &CompactGameModel{
		KeyShards:        g.NumKeyShards,
		FieldLen:         g.FieldLen,
		BombPercent:      g.BombPercent,
		Bombs:            bombs,
		Keys:             keys,
		Empty:            seenEmpty,
		NotSeen:          notSeen,
		IsDefaultNotSeen: def,
	}
}

// func NewGame(g *model.GameModel) (*Game, error) {
// 	if g.KeyShards >= g.FieldLen*g.FieldLen {
// 		return nil, errors.New("(fieldLength)^2 should at least be bigger than number of key shards")
// 	}
// 	return &Game{
// 		NumKeyShards: g.KeyShards,
// 		FieldLen:     g.FieldLen,
// 		BombPercent:  g.BombPercent,
// 		GameBoard:    br,
// 		Seen:         &sr,
// 		SeenCounter:  0,
// 		GameId:       0,
// 	}, nil
// }

func (g *Game) MakeMove(row, col int) (int8, error) {
	if row >= g.FieldLen || col >= g.FieldLen || row < 0 || col < 0 {
		return 0, errors.New("row or col is out of range")
	}
	if (*g.Seen)[row][col] {
		return 0, errors.New("this cell is already seen")
	}
	(*g.Seen)[row][col] = true
	g.SeenCounter++
	return (*g.GameBoard)[row][col], nil
}

func LoadGame(g *model.GameModel) (*Game, error) {
	if g == nil {
		return nil, errors.New("model is nil")
	}
	if g.KeyShards >= g.FieldLen*g.FieldLen {
		return nil, errors.New("(fieldLength)^2 should at least be bigger than number of key shards")
	}
	if !isBoardAndSeenEqual(g.Board, g.Seen) {
		return nil, errors.New("board and seen slices should be of the same length")
	}
	if len(g.Board) != g.FieldLen {
		return nil, errors.New("board length should be the same az fieldLen")
	}
	b := Board(g.Board)
	s := Seen(g.Seen)
	return &Game{
		NumKeyShards: g.KeyShards,
		FieldLen:     g.FieldLen,
		BombPercent:  g.BombPercent,
		GameBoard:    &b,
		Seen:         &s,
		SeenCounter:  calcSeenCounter(s),
		GameId:       g.GameId,
	}, nil
}
