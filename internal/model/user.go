package model

type UserModel struct {
	Username       string `db:"username"`
	NumOfKeys      int    `db:"number_of_keys"`
	Id             int    `db:"id"`
	PointsLeft     int    `db:"points_left"`
	NextMoveCost   int    `db:"next_move_cost"`
	NormalMoveCost int    `db:"normal_move_cost"`
	BombMoveCost   int    `db:"bomb_move_cost"`
}

type RPSUserModel struct {
	Username   string `db:"username"`
	Id         int    `db:"id"`
	PointsLeft int    `db:"points_left"`
}
