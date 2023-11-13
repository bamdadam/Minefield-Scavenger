package model

import "time"

type GameModel struct {
	KeyShards   int       `db:"key_shards"`
	FieldLen    int       `db:"field_len"`
	BombPercent int       `db:"bomb_percent"`
	PlayerId    int       `db:"player_id"`
	GameId      int       `db:"id"`
	Board       [][]int8  `db:"board"`
	Seen        [][]bool  `db:"seen"`
	CreatedAt   time.Time `db:"created_at"`
}
