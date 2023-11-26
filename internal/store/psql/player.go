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

func (p *PSQLStore) UpdateUser(ctx context.Context, uID, numberOfKeys, PointsLeft, NextMoveCost, NormalMoveCost, BombMoveCost int) error {
	_, err := p.DB.Exec(ctx, `
		UPDATE users SET 
			(number_of_keys, points_left, next_move_cost, normal_move_cost, bomb_move_cost) = ($1, $2, $3, $4, $5)
		WHERE id = $6
	`, numberOfKeys, PointsLeft, NextMoveCost, NormalMoveCost, BombMoveCost, uID)
	return err
}

func (p *PSQLStore) CreateRPSUser(ctx context.Context, username string, numOfPoints int) (*model.RPSUserModel, error) {
	u := new(model.RPSUserModel)
	err := p.DB.QueryRow(ctx, `
		INSERT INTO rps_users (username, points_left) VALUES ($1, $2)
		ON CONFLICT(username) DO NOTHING
		RETURNING id, username, points_left
	`, username, numOfPoints,
	).Scan(&u.Id, &u.Username, &u.PointsLeft)
	return u, err
}

func (p *PSQLStore) GetRPSUser(ctx context.Context, username string) (*model.RPSUserModel, error) {
	user := new(model.RPSUserModel)
	err := pgxscan.Get(ctx, p.DB, user, `
		SELECT * FROM rps_users WHERE username = $1
	`, username)
	return user, err
}

func (p *PSQLStore) UpdateRPSUser(ctx context.Context, uID, PointsLeft int) error {
	_, err := p.DB.Exec(ctx, `
		UPDATE rps_users SET 
			points_left = $1
		WHERE id = $2
	`, PointsLeft, uID)
	return err
}
