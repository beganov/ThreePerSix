package game

import (
	"sync"

	"github.com/beganov/gingonicserver/internal/domain/core/card"
	"github.com/rs/zerolog"
)

type GameState struct { // Объект игры
	sync.Mutex   `json:"-" swaggerignore:"true"`
	Deck         []card.Card      `json:"deck,omitempty"`      // Колода
	Out          []card.Card      `json:"out,omitempty"`       // Массив карт в игре, то что на столе
	Hands        [][]card.Card    `json:"hands,omitempty"`     // Массив рук игроков
	Openeds      [][]card.Card    `json:"openeds,omitempty"`   // Массив открытых карт
	Closeds      [][]card.Card    `json:"closeds,omitempty"`   // Массив закрытых карт
	PlayerNow    int              `json:"playerNow,omitempty"` // Номер за столом ходящего сейчас игрока
	Turn         int              `json:"turn,omitempty"`      // Номер хода
	IsMoved      bool             `json:"isMoved,omitempty"`   // Ходил ли игрок или нет, используется для фронта
	IdMap        map[int]int      `json:"idMap,omitempty"`     // Мапа для поиска номера игрока за столом по его id; key = playerId, value = placement
	ReverceIdMap map[int]int      // Мапа для поиска id игрока по его номеру столом; key = placement, value = playerId
	ch           map[int]chan int // Мапа каналов, для передачи хода игроками; key = playerId, value = chanel
	handler      GameEndHandler   // Интерфейс комнаты, вызываемый при завершении игры
	logger       Logger           // Интерфейс логгера
}

type GameEndHandler interface { //Интерфейс обратного вызова, реализуемый комнатой (room.Room). Вызывается при завершении игры.
	OnGameEnd()
}

type Logger interface { // Интерфейс логгера
	Info() *zerolog.Event
	Error() *zerolog.Event
}
