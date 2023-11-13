package response

import "github.com/bamdadam/Minefield-Scavenger/internal/game"

type GetUserDataResponse struct {
	Username       string `json:"username"`
	NumOfKeys      int    `json:"number_of_keys"`
	PointsLeft     int    `json:"points_left"`
	NextMoveCost   int    `json:"next_move_cost"`
	NormalMoveCost int    `json:"normal_move_cost"`
	BombMoveCost   int    `json:"bomb_move_cost"`
}

type PlayGameResponse struct {
	ActiveGame     *game.CompactGameModel `json:"game_state"`
	Username       string                 `json:"username"`
	NumOfKeys      int                    `json:"number_of_keys"`
	PointsLeft     int                    `json:"points_left"`
	NextMoveCost   int                    `json:"next_move_cost"`
	NormalMoveCost int                    `json:"normal_move_cost"`
	BombMoveCost   int                    `json:"bomb_move_cost"`
}
