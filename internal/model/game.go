package model

import "github.com/bamdadam/Minefield-Scavenger/internal/game"

type GameModel struct {
	keyShards   int
	fieldLen    int
	bombPercent int
	playerId    int
	gameId      int
	board       game.Board
	seen        game.Seen
}
