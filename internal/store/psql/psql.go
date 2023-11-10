package psql

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type PSQLStore struct {
	DB *pgxpool.Pool
}

func NewPSQLStore(db *pgxpool.Pool) *PSQLStore {
	return &PSQLStore{
		DB: db,
	}
}
