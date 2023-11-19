package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/bamdadam/Minefield-Scavenger/internal/game"
	"github.com/bamdadam/Minefield-Scavenger/internal/http/middleware"
	"github.com/bamdadam/Minefield-Scavenger/internal/http/request"
	"github.com/bamdadam/Minefield-Scavenger/internal/http/response"
	"github.com/bamdadam/Minefield-Scavenger/internal/model"
	"github.com/bamdadam/Minefield-Scavenger/internal/player"
	"github.com/bamdadam/Minefield-Scavenger/internal/store"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)

type PlayerHandler struct {
	db          store.Store
	playerCache *player.ActivePlayersCache
}

func NewPlayerHandler(db store.Store, c *player.ActivePlayersCache) *PlayerHandler {
	return &PlayerHandler{
		db:          db,
		playerCache: c,
	}
}

func (p *PlayerHandler) login(ctx *fiber.Ctx) error {
	body := request.LoginRequest{}
	err := ctx.BodyParser(&body)
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(err.Error())
	}
	fmt.Println(body)
	fieldLen := body.FieldLen
	bombPercent := body.BombPercent
	keyShards := body.NumOfKeys
	var user *model.UserModel
	var gameModel *model.GameModel
	user, err = p.db.GetUser(ctx.Context(), body.Username)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, pgx.ErrNoRows) {
			usr, err := createNewUserModel(body.Username, 0, body.Points, body.OpeningCost, body.OpeningCost, body.BombOpeningCost)
			if err != nil {
				return ctx.Status(fiber.ErrBadRequest.Code).JSON(err.Error())
			}
			user, err = p.db.CreateUser(ctx.Context(), *usr)
			if err != nil {
				fmt.Println(err)
				return ctx.Status(fiber.ErrBadRequest.Code).JSON(err.Error())
			}
			// create and save a new game model for the new user
			gm, err := createNewGameModel(fieldLen, bombPercent, keyShards, user.Id)
			if err != nil {
				return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
			}
			gameModel, err = p.db.CreateNewGame(ctx.Context(), *gm)
			if err != nil {
				return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
			}
		} else {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
		}
	} else {
		gameModel, err = p.db.RetrieveTodaysGame(ctx.Context(), user.Id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				// create and save a new game model for the user
				gm, err := createNewGameModel(fieldLen, bombPercent, keyShards, user.Id)
				if err != nil {
					return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
				}
				gameModel, err = p.db.CreateNewGame(ctx.Context(), *gm)
				if err != nil {
					return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
				}
			} else {
				return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
			}
		}
		// update the user
		user.NumOfKeys = 0
		user.PointsLeft = body.Points
		user.NextMoveCost = body.OpeningCost
		user.BombMoveCost = body.BombOpeningCost
		user.NormalMoveCost = body.OpeningCost
		// save player
		err = p.db.UpdateUser(ctx.Context(), user.Id, user.NumOfKeys, user.PointsLeft, user.NextMoveCost, user.NormalMoveCost, user.BombMoveCost)
		if err != nil {
			fmt.Println("error in updating user: ", err)
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
		}
	}
	fmt.Println(user)
	fmt.Println(gameModel)
	t, exp, err := createJWTToken(user.Id, user.Username, "test")
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
	}
	_, err = p.playerCache.SetPlayer(user, gameModel)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"token": t,
		"exp":   exp,
	})
}

func (p *PlayerHandler) playTurn(ctx *fiber.Ctx) error {
	body := new(request.PlayRequest)
	err := ctx.BodyParser(body)
	token := ctx.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	uID := claims["user_id"].(float64)
	username := claims["username"].(string)
	plyr, err := p.playerCache.GetPlayer(int(uID))
	if err != nil {
		var user *model.UserModel
		var gameModel *model.GameModel
		user, err = p.db.GetUser(ctx.Context(), username)
		if err != nil {
			fmt.Println("expected to find user but did not: ", err)
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
		} else {
			gameModel, err = p.db.RetrieveTodaysGame(ctx.Context(), user.Id)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ctx.Status(fiber.ErrInternalServerError.Code).JSON("please login again")
				} else {
					fmt.Println("could not create game: ", err)
					return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
				}
			}
		}
		fmt.Println(user)
		fmt.Println(gameModel)
		plyr, err = p.playerCache.SetPlayer(user, gameModel)
		if err != nil {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
		}
	}
	fmt.Printf("points at: %p\n", plyr.ActiveGame)
	err = plyr.MakeMove(body.X, body.Y)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(err.Error())
	}
	// save player
	err = p.db.UpdateUser(ctx.Context(), int(uID), plyr.NumberOfKeys, plyr.PointsLeft, plyr.NextMoveCost, plyr.NormalMoveCost, plyr.BombMoveCost)
	if err != nil {
		fmt.Println("error in updating user: ", err)
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
	}
	// save game
	err = p.db.UpdateGame(
		ctx.Context(),
		plyr.ActiveGame.GameId,
		plyr.ActiveGame.NumKeyShards,
		plyr.ActiveGame.BombPercent,
		plyr.ActiveGame.FieldLen,
		*plyr.ActiveGame.GameBoard,
		*plyr.ActiveGame.Seen)

	if err != nil {
		fmt.Println("error in updating game: ", err)
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
	}
	gameState := plyr.ShowGameState()
	res := response.PlayGameResponse{
		ActiveGame:     gameState,
		Username:       plyr.Username,
		NumOfKeys:      plyr.NumberOfKeys,
		PointsLeft:     plyr.PointsLeft,
		NextMoveCost:   plyr.NextMoveCost,
		NormalMoveCost: plyr.NormalMoveCost,
		BombMoveCost:   plyr.BombMoveCost,
	}
	return ctx.Status(http.StatusOK).JSON(res)
}

func (p *PlayerHandler) LoseGame(ctx *fiber.Ctx) error {
	token := ctx.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	uID := claims["user_id"].(float64)
	username := claims["username"].(string)
	var oldGame *game.Game
	plyr, err := p.playerCache.GetPlayer(int(uID))
	if err != nil {
		var user *model.UserModel
		var gameModel *model.GameModel
		user, err = p.db.GetUser(ctx.Context(), username)
		if err != nil {
			fmt.Println("expected to find user but did not: ", err)
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
		} else {
			gameModel, err = p.db.RetrieveTodaysGame(ctx.Context(), user.Id)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					// create and save a new game model for the user
					gm, err := createNewGameModel(10, 50, 5, user.Id)
					if err != nil {
						return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
					}
					gameModel, err = p.db.CreateNewGame(ctx.Context(), *gm)
					if err != nil {
						return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
					}
				} else {
					return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
				}
			} else {
				oldGame, err = game.LoadGame(gameModel)
				if err != nil {
					return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
				}
				gm, err := createNewGameModel(
					gameModel.FieldLen,
					gameModel.BombPercent,
					gameModel.KeyShards,
					user.Id)
				if err != nil {
					return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
				}
				// update the existing game model for the user
				err = p.db.UpdateGame(
					ctx.Context(),
					gameModel.GameId,
					gameModel.KeyShards,
					gameModel.BombPercent,
					gameModel.FieldLen,
					gm.Board,
					gm.Seen)

				if err != nil {
					fmt.Println("error in updating game: ", err)
					return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
				}
				gameModel.Board = gm.Board
				gameModel.Seen = gm.Seen
			}
			// update the user
			user.NumOfKeys = 0
			// save player
			err = p.db.UpdateUser(ctx.Context(), user.Id, user.NumOfKeys, user.PointsLeft, user.NextMoveCost, user.NormalMoveCost, user.BombMoveCost)
			if err != nil {
				fmt.Println("error in updating user: ", err)
				return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
			}
		}
		fmt.Println(user)
		fmt.Println(gameModel)
		plyr, err = p.playerCache.SetPlayer(user, gameModel)
		if err != nil {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
		}
	} else {
		oldGame = plyr.ActiveGame
		gm, err := createNewGameModel(
			plyr.ActiveGame.FieldLen,
			plyr.ActiveGame.BombPercent,
			plyr.ActiveGame.NumKeyShards,
			int(uID))
		if err != nil {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
		}
		// save player
		err = p.db.UpdateUser(ctx.Context(), int(uID), 0, plyr.PointsLeft, plyr.NormalMoveCost, plyr.NormalMoveCost, plyr.BombMoveCost)
		if err != nil {
			fmt.Println("error in updating user: ", err)
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
		}
		plyr.NextMoveCost = plyr.NormalMoveCost
		plyr.NumberOfKeys = 0
		// save game
		err = p.db.UpdateGame(
			ctx.Context(),
			plyr.ActiveGame.GameId,
			plyr.ActiveGame.NumKeyShards,
			plyr.ActiveGame.BombPercent,
			plyr.ActiveGame.FieldLen,
			gm.Board,
			gm.Seen)

		if err != nil {
			fmt.Println("error in updating game: ", err)
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
		}
		plyr.ActiveGame, err = game.LoadGame(gm)
		if err != nil {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
		}
	}
	gameState := oldGame.ShowStateOnLose()
	res := response.PlayGameResponse{
		ActiveGame:     gameState,
		Username:       plyr.Username,
		NumOfKeys:      plyr.NumberOfKeys,
		PointsLeft:     plyr.PointsLeft,
		NextMoveCost:   plyr.NextMoveCost,
		NormalMoveCost: plyr.NormalMoveCost,
		BombMoveCost:   plyr.BombMoveCost,
	}
	return ctx.Status(http.StatusOK).JSON(res)
}

func (p *PlayerHandler) getUserData(ctx *fiber.Ctx) error {
	token := ctx.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	uID := claims["user_id"].(float64)
	username := claims["username"].(string)
	plyr, err := p.playerCache.GetPlayer(int(uID))
	if err != nil {
		var user *model.UserModel
		var gameModel *model.GameModel
		user, err = p.db.GetUser(ctx.Context(), username)
		if err != nil {
			fmt.Println("expected to find user but did not: ", err)
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
		} else {
			gameModel, err = p.db.RetrieveTodaysGame(ctx.Context(), user.Id)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					gameModel = nil
				} else {
					fmt.Println("could not create game: ", err)
					return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
				}
			}
		}
		fmt.Println(user)
		fmt.Println(gameModel)
		res := response.GetUserDataResponse{
			Username:       user.Username,
			NumOfKeys:      user.NumOfKeys,
			PointsLeft:     user.PointsLeft,
			NextMoveCost:   user.NextMoveCost,
			NormalMoveCost: user.NormalMoveCost,
			BombMoveCost:   user.BombMoveCost,
		}
		activeGame, err := game.LoadGame(gameModel)
		if err != nil {
			if err.Error() == "model is nil" {
				res.ActiveGame = nil
			} else {
				return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
			}
		} else {
			res.ActiveGame = activeGame.ShowState()
		}
		return ctx.Status(http.StatusOK).JSON(res)
	}
	res := response.GetUserDataResponse{
		Username:       plyr.Username,
		NumOfKeys:      plyr.NumberOfKeys,
		PointsLeft:     plyr.PointsLeft,
		NextMoveCost:   plyr.NextMoveCost,
		NormalMoveCost: plyr.NormalMoveCost,
		BombMoveCost:   plyr.BombMoveCost,
		ActiveGame:     plyr.ShowGameState(),
	}
	return ctx.Status(http.StatusOK).JSON(res)
}

func (p *PlayerHandler) RestartGame(ctx *fiber.Ctx) error {
	body := new(request.RestartRequest)
	err := ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(err.Error())
	}
	token := ctx.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	uID := claims["user_id"].(float64)
	username := claims["username"].(string)
	// create a new game model
	fieldLen := body.FieldLen
	bombPercent := body.BombPercent
	keyShards := body.NumOfKeys
	gm, err := createNewGameModel(fieldLen, bombPercent, keyShards, int(uID))
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
	}
	plyr, err := p.playerCache.GetPlayer(int(uID))
	if err != nil {
		var user *model.UserModel
		var gameModel *model.GameModel
		user, err = p.db.GetUser(ctx.Context(), username)
		if err != nil {
			fmt.Println("expected to find user but did not: ", err)
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
		} else {
			gameModel, err = p.db.RetrieveTodaysGame(ctx.Context(), user.Id)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					// save the new game model for the user
					gameModel, err = p.db.CreateNewGame(ctx.Context(), *gm)
					if err != nil {
						fmt.Println("could not create a new game model: ", err)
						return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
					}
				} else {
					fmt.Println("could not retrieve game: ", err)
					return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
				}
			} else {
				// update the existing game model for the user
				err = p.db.UpdateGame(
					ctx.Context(),
					gameModel.GameId,
					keyShards,
					bombPercent,
					fieldLen,
					gm.Board,
					gm.Seen)

				if err != nil {
					fmt.Println("error in updating game: ", err)
					return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
				}
			}
		}
		fmt.Println(user)
		fmt.Println(gameModel)
		plyr, err = p.playerCache.SetPlayer(user, gameModel)
		if err != nil {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
		}
	} else {
		activeGame, err := game.LoadGame(gm)
		if err != nil {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
		}
		// save game
		err = p.db.UpdateGame(
			ctx.Context(),
			plyr.ActiveGame.GameId,
			activeGame.NumKeyShards,
			activeGame.BombPercent,
			activeGame.FieldLen,
			*activeGame.GameBoard,
			*activeGame.Seen)

		if err != nil {
			fmt.Println("error in updating game: ", err)
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
		}

		plyr.ActiveGame = activeGame
	}
	plyr.PointsLeft = body.TopUpPoints
	plyr.BombMoveCost = body.BombOpeningCost
	plyr.NextMoveCost = body.OpeningCost
	plyr.NormalMoveCost = body.OpeningCost
	plyr.NumberOfKeys = 0

	// save player
	err = p.db.UpdateUser(ctx.Context(), int(uID), plyr.NumberOfKeys, plyr.PointsLeft, plyr.NextMoveCost, plyr.NormalMoveCost, plyr.BombMoveCost)
	if err != nil {
		fmt.Println("error in updating user: ", err)
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (p *PlayerHandler) payForBomb(ctx *fiber.Ctx) error {
	body := new(request.PlayRequest)
	err := ctx.BodyParser(body)
	token := ctx.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	uID := claims["user_id"].(float64)
	username := claims["username"].(string)
	plyr, err := p.playerCache.GetPlayer(int(uID))
	if err != nil {
		var user *model.UserModel
		var gameModel *model.GameModel
		user, err = p.db.GetUser(ctx.Context(), username)
		if err != nil {
			fmt.Println("expected to find user but did not: ", err)
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
		} else {
			gameModel, err = p.db.RetrieveTodaysGame(ctx.Context(), user.Id)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ctx.Status(fiber.ErrInternalServerError.Code).JSON("please login again")
				} else {
					fmt.Println("could not create game: ", err)
					return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
				}
			}
		}
		fmt.Println(user)
		fmt.Println(gameModel)
		plyr, err = p.playerCache.SetPlayer(user, gameModel)
		if err != nil {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
		}
	}
	plyr.PointsLeft -= plyr.BombMoveCost
	// save player
	err = p.db.UpdateUser(ctx.Context(), int(uID), plyr.NumberOfKeys, plyr.PointsLeft, plyr.NextMoveCost, plyr.NormalMoveCost, plyr.BombMoveCost)
	if err != nil {
		fmt.Println("error in updating user: ", err)
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err.Error())
	}
	gameState := plyr.ShowGameState()
	res := response.PlayGameResponse{
		ActiveGame:     gameState,
		Username:       plyr.Username,
		NumOfKeys:      plyr.NumberOfKeys,
		PointsLeft:     plyr.PointsLeft,
		NextMoveCost:   plyr.NextMoveCost,
		NormalMoveCost: plyr.NormalMoveCost,
		BombMoveCost:   plyr.BombMoveCost,
	}
	return ctx.Status(http.StatusOK).JSON(res)
}

func (p *PlayerHandler) RegisterHandlers(g fiber.Router) {
	g.Post("/play", middleware.Protected("test"), p.playTurn)
	g.Post("/login", p.login)
	g.Get("/data", middleware.Protected("test"), p.getUserData)
	g.Post("/restart", middleware.Protected("test"), p.RestartGame)
	g.Post("/lose", middleware.Protected("test"), p.LoseGame)
	g.Post("/play/bomb", middleware.Protected("test"), p.payForBomb)
}
