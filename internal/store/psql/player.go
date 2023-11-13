package psql

import (
	"context"

	"github.com/bamdadam/Minefield-Scavenger/internal/model"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (p *PSQLStore) GetUser(ctx context.Context, username string) (*model.UserModel, error) {
	user := new(model.UserModel)
	err := pgxscan.Get(ctx, p.DB, user, `
		SELECT * FROM users WHERE username = $1
	`, username)
	return user, err
}

func (p *PSQLStore) CreateUser(ctx context.Context, user model.UserModel) (*model.UserModel, error) {
	u := new(model.UserModel)
	err := p.DB.QueryRow(ctx, `
		INSERT INTO users (username, number_of_keys, points_left, normal_move_cost, bomb_move_cost, next_move_cost) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT(username) DO NOTHING
		RETURNING id, username, number_of_keys, points_left, normal_move_cost, bomb_move_cost, next_move_cost
	`, user.Username, 0, user.PointsLeft, user.NormalMoveCost, user.BombMoveCost, user.NextMoveCost,
	).Scan(&u.Id, &u.Username, &u.NumOfKeys, &u.PointsLeft, &u.NormalMoveCost, &u.BombMoveCost, &u.NextMoveCost)
	return u, err
}

func (p *PSQLStore) UpdateUser(ctx context.Context, numberOfKeys, PointsLeft, NextMoveCost int) error {
	_, err := p.DB.Exec(ctx, `
		UPDATE users SET 
			(number_of_keys, points_left, next_move_cost) = ($1, $2, $3)
	`, numberOfKeys, PointsLeft, NextMoveCost)
	return err
}
