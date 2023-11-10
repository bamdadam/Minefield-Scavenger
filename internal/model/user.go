package model

type UserModel struct {
	Username  string `db:"username"`
	NumOfKeys int    `db:"number_of_keys"`
	Id        int    `db:"id"`
}
