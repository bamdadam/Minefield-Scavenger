package store

import (
	"context"
	"time"

	"github.com/bamdadam/Minefield-Scavenger/internal/game"
	"github.com/bamdadam/Minefield-Scavenger/internal/model"
)

type Store interface {
	// game
	CreateNewGame(ctx context.Context, keyShards, fieldLen, bombPercent, playerId int, board game.Board, seen game.Seen) (int, error)
	UpdateGame(ctx context.Context, gameId int, board game.Board, seen game.Seen) error
	RetrieveGame(ctx context.Context, playerId int, date time.Time) (*model.GameModel, error)
	// user
	CreateUser(ctx context.Context, username string) (*model.UserModel, error)
	GetUser(ctx context.Context, username string) (*model.UserModel, error)
}
