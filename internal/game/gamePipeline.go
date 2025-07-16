package game

func (g *GameState) StartGame(MaxPlayerCount int, Players map[int]int) *GameState {
	g.SafeInitialization(MaxPlayerCount, Players)
	go func() {
		g.Initialization()
		g.Game()
	}()
	return g
}
