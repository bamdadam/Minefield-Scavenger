package psql

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bamdadam/Minefield-Scavenger/internal/game"
	"github.com/bamdadam/Minefield-Scavenger/internal/model"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (p *PSQLStore) CreateNewGame(ctx context.Context, keyShards, fieldLen, bombPercent, playerId int, board game.Board, seen game.Seen) (int, error) {
	bJson, err := json.Marshal(board)
	if err != nil {
		return 0, err
	}
	sJson, err := json.Marshal(seen)
	if err != nil {
		return 0, err
	}
	id := new(int)
	err = p.DB.QueryRow(ctx,
		`INSERT INTO games 
			(key_shards, field_len, bomb_percent, player_id, board, seen)
		VALUES 
			($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		keyShards, fieldLen, bombPercent, playerId, bJson, sJson).Scan(id)
	return *id, err
}

func (p *PSQLStore) UpdateGame(ctx context.Context, gameId int, board game.Board, seen game.Seen) error {
	bJson, err := json.Marshal(board)
	if err != nil {
		return err
	}
	sJson, err := json.Marshal(seen)
	if err != nil {
		return err
	}
	_, err = p.DB.Exec(ctx,
		`UPDATE games SET
		(board, seen) = ($1, $2)
	WHERE id = $3`,
		bJson, sJson)
	return err
}

func (p *PSQLStore) RetrieveGame(ctx context.Context, playerId int, date time.Time) (*model.GameModel, error) {
	m := new(model.GameModel)
	err := pgxscan.Get(ctx, p.DB, m,
		`SELECT * FROM games
		WHERE player_id = $1 and created_at = $2`,
		playerId, date)
	return m, err
}
