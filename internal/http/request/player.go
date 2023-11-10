package request

type LoginRequest struct {
	Username string `json:"username"`
}

type PlayRequest struct {
	X int `json:"x"`
	Y int `json:"y"`
}
