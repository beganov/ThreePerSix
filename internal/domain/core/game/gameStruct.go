package game

import (
	"sync"

	"github.com/beganov/gingonicserver/internal/domain/core/card"
	"github.com/rs/zerolog"
)

type GameState struct {
	sync.Mutex   `json:"-" swaggerignore:"true"`
	Deck         []card.Card   `json:"deck,omitempty"`
	Out          []card.Card   `json:"out,omitempty"`
	Hands        [][]card.Card `json:"hands,omitempty"`
	Openeds      [][]card.Card `json:"openeds,omitempty"`
	Closeds      [][]card.Card `json:"closeds,omitempty"`
	PlayerNow    int           `json:"playerNow,omitempty"`
	Turn         int           `json:"turn,omitempty"`
	IsMoved      bool          `json:"isMoved,omitempty"`
	IdMap        map[int]int   `json:"idMap,omitempty"` //key = playerId, value = placement
	ReverceIdMap map[int]int   //key = placement, value = playerId
	ch           map[int]chan int
	handler      GameEndHandler
	logger       Logger
}

type GameEndHandler interface {
	OnGameEnd()
}

type Logger interface {
	Info() *zerolog.Event
	Error() *zerolog.Event
}
