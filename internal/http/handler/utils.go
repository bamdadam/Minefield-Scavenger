package handler

import (
	"crypto/rand"
	"math"
	"math/big"
	"time"

	"github.com/bamdadam/Minefield-Scavenger/internal/game"
	"github.com/bamdadam/Minefield-Scavenger/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

type cell int8

const (
	empty cell = iota
	keyShard
	bomb
)

func createJWTToken(userId int, username, secret string) (string, int64, error) {
	exp := time.Now().Add(time.Hour * 10).Unix()
	claims := jwt.MapClaims{
		"user_id":  userId,
		"username": username,
		"exp":      exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return t, exp, err
	}
	return t, exp, nil
}

func createBoard(bLen, bombPercent, numKeyShards int) ([][]int8, error) {
	br := make([][]int8, bLen)
	for i := 0; i < bLen; i++ {
		br[i] = make([]int8, bLen)
	}
	numBombs := bLen * bLen * bombPercent / 100
	for numBombs > 0 {
		rn, err := rand.Int(rand.Reader, big.NewInt(int64(bLen*bLen)))
		if err != nil {
			return nil, err
		}
		row := int(math.Floor(float64(rn.Int64() / int64(bLen))))
		col := rn.Int64() % int64(bLen)
		if br[row][col] != int8(bomb) {
			numBombs--
			br[row][col] = int8(bomb)
		}
	}
	for numKeyShards > 0 {
		rn, err := rand.Int(rand.Reader, big.NewInt(int64(bLen*bLen)))
		if err != nil {
			return nil, err
		}
		row := int(math.Floor(float64(rn.Int64() / int64(bLen))))
		col := rn.Int64() % int64(bLen)
		if br[row][col] != int8(bomb) && br[row][col] != int8(keyShard) {
			numKeyShards--
			br[row][col] = int8(keyShard)
		}
	}
	return br, nil
}

func createNewGameModel(fieldLen, bombPercent, keyShards, userId int) (*model.GameModel, error) {
	br, err := createBoard(fieldLen, bombPercent, keyShards)
	if err != nil {
		return nil, err
	}

	sr := make(game.Seen, fieldLen)
	for i := 0; i < fieldLen; i++ {
		sr[i] = make([]bool, fieldLen)
	}

	gm := &model.GameModel{
		KeyShards:   keyShards,
		FieldLen:    fieldLen,
		BombPercent: bombPercent,
		PlayerId:    userId,
		Board:       br,
		Seen:        sr,
	}

	return gm, err
}

func createNewUserModel(userName string, NumOfKeys, PointsLeft, NextMoveCost, NormalMoveCost, BombMoveCost int) (*model.UserModel, error) {
	return &model.UserModel{
		Username:       userName,
		NumOfKeys:      NumOfKeys,
		PointsLeft:     PointsLeft,
		NextMoveCost:   NextMoveCost,
		NormalMoveCost: NormalMoveCost,
		BombMoveCost:   BombMoveCost,
	}, nil
}
