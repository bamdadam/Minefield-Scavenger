package store

import (
	"context"

	"github.com/bamdadam/Minefield-Scavenger/internal/game"
	"github.com/bamdadam/Minefield-Scavenger/internal/model"
)

type Store interface {
	// game
	CreateNewGame(ctx context.Context, g model.GameModel) (*model.GameModel, error)
	UpdateGame(ctx context.Context, gameId, keyShards, bombPercent, fieldLen int, board game.Board, seen game.Seen) error
	RetrieveLastNGame(ctx context.Context, playerId, n int) ([]*model.GameModel, error)
	RetrieveTodaysGame(ctx context.Context, playerId int) (*model.GameModel, error)
	// user
	CreateUser(ctx context.Context, user model.UserModel) (*model.UserModel, error)
	GetUser(ctx context.Context, username string) (*model.UserModel, error)
	UpdateUser(ctx context.Context, uID, numberOfKeys, PointsLeft, NextMoveCost, NormalMoveCost, BombMoveCost int) error
	// rps
	SaveRPSGame(ctx context.Context, playerChoice, houseChoice, playerID int, hasWon bool) error
	// rps user
	CreateRPSUser(ctx context.Context, username string, numOfPoints int) (*model.RPSUserModel, error)
	GetRPSUser(ctx context.Context, username string) (*model.RPSUserModel, error)
	UpdateRPSUser(ctx context.Context, uID, PointsLeft int) error
}
