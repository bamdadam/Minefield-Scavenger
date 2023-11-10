package psql

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPSQLDB(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	timeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	db, err := pgxpool.New(timeout, cfg.ConnectionString)
	if err != nil {
		return nil, err
	}
	err = db.Ping(timeout)
	if err != nil {
		return nil, errors.New("error while pinging db: " + err.Error())
	}
	return db, nil
}
