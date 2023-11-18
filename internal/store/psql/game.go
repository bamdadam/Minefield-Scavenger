package psql

import (
	"context"
	"encoding/json"

	"github.com/bamdadam/Minefield-Scavenger/internal/game"
	"github.com/bamdadam/Minefield-Scavenger/internal/model"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (p *PSQLStore) CreateNewGame(ctx context.Context, g model.GameModel) (*model.GameModel, error) {
	bJson, err := json.Marshal(g.Board)
	if err != nil {
		return nil, err
	}
	sJson, err := json.Marshal(g.Seen)
	if err != nil {
		return nil, err
	}
	gm := new(model.GameModel)
	boardJson := new([]byte)
	seenJson := new([]byte)
	err = p.DB.QueryRow(ctx,
		`INSERT INTO games 
			(key_shards, field_len, bomb_percent, player_id, board, seen)
		VALUES 
			($1, $2, $3, $4, $5, $6)
		RETURNING id, player_id, key_shards, field_len, bomb_percent, board, seen, created_at`,
		g.KeyShards, g.FieldLen, g.BombPercent, g.PlayerId, bJson, sJson).Scan(
		&gm.GameId, &gm.PlayerId, &gm.KeyShards, &gm.FieldLen, &gm.BombPercent, boardJson, seenJson, &gm.CreatedAt,
	)
	err = json.Unmarshal(*boardJson, &gm.Board)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(*seenJson, &gm.Seen)
	if err != nil {
		return nil, err
	}
	return gm, err
}

func (p *PSQLStore) UpdateGame(ctx context.Context, gameId, keyShards, bombPercent, fieldLen int, board game.Board, seen game.Seen) error {
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
			(key_shards, field_len, bomb_percent, board, seen) = ($1, $2, $3, $4, $5)
		WHERE id = $6`,
		keyShards, fieldLen, bombPercent, bJson, sJson, gameId)
	return err
}

func (p *PSQLStore) RetrieveLastNGame(ctx context.Context, playerId, n int) ([]*model.GameModel, error) {
	m := make([]*model.GameModel, n)
	err := pgxscan.Select(ctx, p.DB, m,
		`SELECT * FROM games
		WHERE player_id = $2
		ORDER BY created_at
		LIMIT $1
		`, n, playerId)
	return m, err
}

func (p *PSQLStore) RetrieveTodaysGame(ctx context.Context, playerId int) (*model.GameModel, error) {
	m := new(model.GameModel)
	err := pgxscan.Get(ctx, p.DB, m,
		`SELECT * FROM games
		WHERE player_id = $1 and created_at::date = current_date
		ORDER BY created_at
		LIMIT 1
		`, playerId)
	return m, err
}
