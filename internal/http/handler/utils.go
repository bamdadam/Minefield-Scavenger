package handler

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func createJWTToken(userId int, secret string) (string, int64, error) {
	exp := time.Now().Add(time.Hour * 5).Unix()
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return t, exp, err
	}
	return t, exp, nil
}
