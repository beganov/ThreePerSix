package game

import (
	"sync"

	"github.com/beganov/gingonicserver/internal/card"
)

type GameState struct {
	sync.Mutex
	Deck           []card.Card   `json:"deck,omitempty"`
	Out            []card.Card   `json:"out,omitempty"`
	Hands          [][]card.Card `json:"hands,omitempty"`
	Openeds        [][]card.Card `json:"openeds,omitempty"`
	Closeds        [][]card.Card `json:"closeds,omitempty"`
	MaxPlayerCount int           `json:"maxPlayerCount,omitempty"`
	IdMap          map[int]int   `json:"idMap,omitempty"`        //key = playerId, value = placement
	ReverceIdMap   map[int]int   `json:"reverceIdMap,omitempty"` //key = placement, value = playerId
	ch             map[int]chan int
}
