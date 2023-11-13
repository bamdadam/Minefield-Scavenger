package request

type LoginRequest struct {
	Username        string `json:"username"`
	Points          int    `json:"points"`
	OpeningCost     int    `json:"openingCost"`
	BombOpeningCost int    `json:"bombOpeningCost"`
}

type PlayRequest struct {
	X int `json:"x"`
	Y int `json:"y"`
}
