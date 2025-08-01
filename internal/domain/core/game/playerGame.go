package game

import (
	"context"

	"github.com/beganov/gingonicserver/internal/domain/core/card"
	"github.com/beganov/gingonicserver/internal/domain/core/gameConst"
)

func (g *GameState) StartGame(maxPlayerCount int, Players map[int]int, end GameEndHandler, ctx context.Context) *GameState {
	g.PreInitialization(maxPlayerCount, Players, end, ctx)
	go func() {
		g.Initialization(maxPlayerCount)
		g.Game(maxPlayerCount)
	}()
	return g
}

func (g *GameState) Move(playerId int, playerMove int) *GameState {
	go func() { g.ch[playerId] <- playerMove }()
	return g
}

func (g *GameState) LeaveGame(playerId int) { //
	go func() {
		if len(g.ReverceIdMap) == 0 {
			g.ch[playerId] <- gameConst.LeaveGameCode
			g.Hands[g.IdMap[playerId]] = []card.Card{}
			g.Openeds[g.IdMap[playerId]] = []card.Card{}
		} else {
			g.Hands[g.IdMap[playerId]] = []card.Card{}
			g.Openeds[g.IdMap[playerId]] = []card.Card{}
			g.Closeds[g.IdMap[playerId]] = []card.Card{}
			g.ch[playerId] <- gameConst.LeaveGameCode
		}
	}()
}
