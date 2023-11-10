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

func (p *PSQLStore) CreateUser(ctx context.Context, username string) (*model.UserModel, error) {
	user := new(model.UserModel)
	err := p.DB.QueryRow(ctx, `
		INSERT INTO users (username, number_of_keys) VALUES ($1, $2)
		ON CONFLICT(username) DO NOTHING
		RETURNING id, username, number_of_keys
	`, username, 0).Scan(&user.Id, &user.Username, &user.NumOfKeys)
	return user, err
}
