package player

import (
	"errors"
	"time"

	"github.com/bamdadam/Minefield-Scavenger/internal/game"
	"github.com/bamdadam/Minefield-Scavenger/internal/model"
)

type Player struct {
	playerID        int
	ActiveGame      *game.Game
	NumberOfKeys    int
	lastInteraction time.Time
	PointsLeft      int
	Username        string
	NextMoveCost    int
	NormalMoveCost  int
	BombMoveCost    int
}

type GameState struct {
	ActiveGame   *game.CompactGameModel
	PointsLeft   int
	NextMoveCost int
	Username     string
}

func (p *Player) MakeMove(x, y int) error {
	if p.NumberOfKeys >= 5 {
		return errors.New("can't make any more moves, player has won")
	}
	cell, err := p.ActiveGame.MakeMove(x, y)
	if err != nil {
		return err
	}
	p.PointsLeft -= p.NormalMoveCost
	switch cell {
	case int8(game.Bomb):
		p.NextMoveCost = p.BombMoveCost
	case int8(game.Empty):
		p.NextMoveCost = p.NormalMoveCost
	case int8(game.KeyShard):
		p.NextMoveCost = p.NormalMoveCost
		p.NumberOfKeys++
	}
	return nil
}

func (p *Player) ShowGame() string {
	return p.ActiveGame.String()
}

func (p *Player) ShowGameState() *game.CompactGameModel {
	return p.ActiveGame.ShowState()
}

type ActivePlayersCache map[int]*Player

func NewPlayerCache() *ActivePlayersCache {
	var ac ActivePlayersCache
	ac = make(ActivePlayersCache, 100)
	return &ac
}

func (c *ActivePlayersCache) GetPlayer(id int) (*Player, error) {
	if player, ok := (*c)[id]; ok {
		player.lastInteraction = time.Now()
		return player, nil
	}
	return nil, errors.New("no user found - cache miss")
}

func (c *ActivePlayersCache) Eviction(minEvictionTime time.Duration) {
	for playerId, player := range *c {
		if time.Since(player.lastInteraction) > minEvictionTime {
			delete(*c, playerId)
		}
	}
}

func (c *ActivePlayersCache) SetPlayer(u *model.UserModel, g *model.GameModel) (*Player, error) {
	activeGame, err := game.LoadGame(g)
	if err != nil {
		return nil, err
	}
	p := Player{
		playerID:        u.Id,
		NumberOfKeys:    u.NumOfKeys,
		lastInteraction: time.Now(),
		PointsLeft:      u.PointsLeft,
		Username:        u.Username,
		NextMoveCost:    u.NextMoveCost,
		NormalMoveCost:  u.NormalMoveCost,
		BombMoveCost:    u.BombMoveCost,
		ActiveGame:      activeGame,
	}
	(*c)[p.playerID] = &p
	return &p, nil
}
