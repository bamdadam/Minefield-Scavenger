package player

import (
	"errors"
	"fmt"
	"time"

	"github.com/bamdadam/Minefield-Scavenger/internal/game"
)

type Player struct {
	playerID        int
	ActiveGame      *game.Game
	numberOfKeys    int
	lastInteraction time.Time
}

func (p *Player) MakeMove(x, y int) (int8, error) {
	if p.numberOfKeys >= 5 {
		return 0, errors.New("can't make any more moves, player has won")
	}
	return p.ActiveGame.MakeMove(x, y)
}

func (p *Player) ShowGame() string {
	return p.ActiveGame.String()
}

func (p *Player) ShowGameState() *game.CompactGameModel {
	return p.ActiveGame.ShowState()
}

type ActivePlayersCache map[int]Player

func NewPlayerCache() *ActivePlayersCache {
	var ac ActivePlayersCache
	ac = make(ActivePlayersCache, 100)
	return &ac
}

func (c *ActivePlayersCache) GetPlayer(id int) (*Player, error) {
	if player, ok := (*c)[id]; ok {
		player.lastInteraction = time.Now()
		if player.ActiveGame == nil {
			game, err := game.NewGame(5, 10, 60)
			if err != nil {
				return nil, err
			}
			player.ActiveGame = game
		}
		return &player, nil
	}
	fmt.Println("cache miss, reading from db")
	fmt.Println("getting the player")
	fmt.Println("getting the game if it exists, if not creating a game")
	game, err := game.NewGame(5, 10, 60)
	if err != nil {
		return nil, err
	}
	p := Player{
		playerID:        id,
		ActiveGame:      game,
		numberOfKeys:    0,
		lastInteraction: time.Now(),
	}
	(*c)[id] = p
	return &p, nil
}

func (c *ActivePlayersCache) Eviction(minEvictionTime time.Duration) {
	for playerId, player := range *c {
		if time.Since(player.lastInteraction) > minEvictionTime {
			delete(*c, playerId)
		}
	}
}
