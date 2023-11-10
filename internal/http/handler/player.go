package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/bamdadam/Minefield-Scavenger/internal/http/middleware"
	"github.com/bamdadam/Minefield-Scavenger/internal/http/request"
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
		fmt.Println(err)
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	fmt.Println(body)
	var user *model.UserModel
	user, err = p.db.GetUser(ctx.Context(), body.Username)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, pgx.ErrNoRows) {
			user, err = p.db.CreateUser(ctx.Context(), body.Username)
			if err != nil {
				fmt.Println(err)
				return ctx.Status(fiber.ErrBadRequest.Code).JSON(err)
			}
		} else {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(err)
		}
	}
	fmt.Println(user)
	t, exp, err := createJWTToken(user.Id, "test")
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
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
	player, err := p.playerCache.GetPlayer(int(uID))
	if err != nil {
		fmt.Println(err)
		ctx.SendStatus(fiber.StatusInternalServerError)
	}
	fmt.Printf("points at: %p\n", player.ActiveGame)
	_, err = player.MakeMove(body.X, body.Y)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	fmt.Println(player.ShowGameState())
	return ctx.Status(http.StatusOK).JSON(player.ShowGameState())
}

func (p *PlayerHandler) RegisterHandlers(g fiber.Router) {
	g.Post("/play", middleware.Protected("test"), p.playTurn)
	g.Post("/login", p.login)
}
