package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bamdadam/Minefield-Scavenger/internal/config"
	"github.com/bamdadam/Minefield-Scavenger/internal/db/psql"
	"github.com/bamdadam/Minefield-Scavenger/internal/http/handler"
	"github.com/bamdadam/Minefield-Scavenger/internal/player"
	stor "github.com/bamdadam/Minefield-Scavenger/internal/store/psql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	exitCh := make(chan struct{})
	go gracefullShutdown(exitCh)
	app := fiber.New(fiber.Config{
		IdleTimeout: 30 * time.Second,
	})
	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: "127.0.0.1:5500",
	// 	AllowHeaders: "Origin, Content-Type, Accept",
	// 	AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	// }))
	app.Use(cors.New())
	go func() {
		<-exitCh
		app.ShutdownWithTimeout(30 * time.Second)
	}()
	// player cache
	cache := player.NewPlayerCache()
	ticker := time.NewTicker(time.Minute * 30)
	go func() {
		for {
			select {
			case <-ticker.C:
				cache.Eviction(time.Hour)
			}
		}
	}()
	//db
	cfg := config.New()
	db, err := psql.NewPSQLDB(context.Background(), cfg.Psql)
	if err != nil {
		panic(err)
	}
	store := stor.NewPSQLStore(db)
	ph := handler.NewPlayerHandler(store, cache)
	player := app.Group("/player")
	ph.RegisterHandlers(player)
	if err := app.Listen(":8080"); err != nil {
		fmt.Println("Can't serve server", err)
	}
}

func gracefullShutdown(exitCh chan<- struct{}) {
	c := make(chan os.Signal, 3)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("close handler", "Ctrl+C pressed in Terminal")
		close(exitCh)
	}()
}
