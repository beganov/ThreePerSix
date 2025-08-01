package game

import (
	"context"

	"github.com/beganov/gingonicserver/internal/domain/core/card"
	"github.com/beganov/gingonicserver/internal/domain/core/gameConst"
)

// StartGame запускает новую игру: инициализирует состояние и запускает саму игру в отдельной горутине.
func (g *GameState) StartGame(maxPlayerCount int, Players map[int]int, end GameEndHandler, ctx context.Context) *GameState {
	g.PreInitialization(maxPlayerCount, Players, end, ctx) // Инициализация того, что не зависит от ввода пользователя
	go func() {
		g.Initialization(maxPlayerCount) // Инициализация того, что зависит от ввода пользователя
		g.Game(maxPlayerCount)           // Запуск игрового цикла
	}()
	return g
}

// Move считывает действия игрока и передает их в его личный канал
func (g *GameState) Move(playerId int, playerMove int) *GameState {
	go func() { g.ch[playerId] <- playerMove }()
	return g
}

// LeaveGame позволяет игроку выйти из игры
func (g *GameState) LeaveGame(playerId int) {
	go func() {
		if len(g.ReverceIdMap) == 0 { // Если игровой цикл еще не начался (так как reverseIdMap инициализируется значениями до игрового цикла)
			g.ch[playerId] <- gameConst.LeaveGameCode // Передаем код выхода и зануляем все, кроме закрытых
			g.Hands[g.IdMap[playerId]] = []card.Card{}
			g.Openeds[g.IdMap[playerId]] = []card.Card{}
		} else { // Если игровой цикл  начался
			g.Hands[g.IdMap[playerId]] = []card.Card{}
			g.Openeds[g.IdMap[playerId]] = []card.Card{}
			g.Closeds[g.IdMap[playerId]] = []card.Card{}
			g.ch[playerId] <- gameConst.LeaveGameCode // Передаем код выхода и зануляем все
		}
	}()
}
