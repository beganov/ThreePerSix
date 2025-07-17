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
	Iamind         map[int]int   `json:"iamind,omitempty"`
	Alsoiamind     map[int]int   `json:"alsoIamind,omitempty"`
	Iamindalso     map[int]int   `json:"iamindalso,omitempty"`
	ch             map[int]chan int
}
