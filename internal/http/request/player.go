package request

type LoginRequest struct {
	Username        string `json:"username"`
	Points          int    `json:"points"`
	OpeningCost     int    `json:"opening_cost"`
	BombOpeningCost int    `json:"bomb_opening_cost"`
	FieldLen        int    `json:"field_len"`
	BombPercent     int    `json:"bomb_percentage"`
	NumOfKeys       int    `json:"number_of_keys"`
}

type PlayRequest struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type RestartRequest struct {
	TopUpPoints     int `json:"top_up_points"`
	OpeningCost     int `json:"opening_cost"`
	BombOpeningCost int `json:"bomb_opening_cost"`
	FieldLen        int `json:"field_len"`
	BombPercent     int `json:"bomb_percentage"`
	NumOfKeys       int `json:"number_of_keys"`
}

type LoginRPSRequest struct {
	Username string `json:"username"`
	Points   int    `json:"points"`
}

type PlayRPSRequest struct {
	PlayerChoice int `json:"player_choice"`
	PlayerBet    int `json:"player_bet"`
	GameVersion  int `json:"game_version"`
}
