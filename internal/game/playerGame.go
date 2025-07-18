package game

import (
	"fmt"

	"github.com/beganov/gingonicserver/internal/card"
	"github.com/beganov/gingonicserver/internal/gameConst"
)

func (g *GameState) StartGame(MaxPlayerCount int, Players map[int]int) *GameState {
	g.PreInitialization(MaxPlayerCount, Players)
	go func() {
		g.Initialization()
		g.Game()
	}()
	return g
}

func (g *GameState) Move(playerId int, playerMove int) *GameState {
	go func() { g.ch[playerId] <- playerMove }()
	return g
}

func (g *GameState) LeaveGame(playerId int) { //
	go func() {
		fmt.Println("break")
		if len(g.ReverceIdMap) == 0 {
			fmt.Println("break")
			g.ch[playerId] <- gameConst.LeaveGameCode
			g.Hands[g.IdMap[playerId]] = []card.Card{}
			g.Openeds[g.IdMap[playerId]] = []card.Card{}
			fmt.Println("break")
		} else {
			g.Hands[g.IdMap[playerId]] = []card.Card{}
			g.Openeds[g.IdMap[playerId]] = []card.Card{}
			g.Closeds[g.IdMap[playerId]] = []card.Card{}
			g.ch[playerId] <- gameConst.LeaveGameCode
		}
	}()
}
