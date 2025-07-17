package game

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

func (g *GameState) LeaveGame(playerId int) {
	g.Lock()
	defer g.Unlock()
	delete(g.ReverceIdMap, g.IdMap[playerId])
	g.ch[playerId] <- g.Hands[g.IdMap[playerId]][0].Val
	g.ch[playerId] <- g.Hands[g.IdMap[playerId]][0].Val
	g.ch[playerId] <- g.Hands[g.IdMap[playerId]][0].Val
	g.ch[playerId] <- g.Hands[g.IdMap[playerId]][0].Val
	delete(g.IdMap, playerId)
	delete(g.ch, playerId)

}
