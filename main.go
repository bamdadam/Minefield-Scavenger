package main

import (
	"fmt"

	"github.com/bamdadam/Minefield-Scavenger/internal/game"
)

func main() {
	game, err := game.NewGame(5, 10, 50)
	if err != nil {
		panic(err)
	}
	fmt.Println(game)
	game.MakeMove(0, 0)
	game.MakeMove(2, 6)
	game.MakeMove(5, 5)
	game.MakeMove(4, 8)
	game.MakeMove(1, 5)
	fmt.Println(game)
}
