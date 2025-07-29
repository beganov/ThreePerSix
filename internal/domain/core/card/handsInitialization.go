package card

import (
	"math/rand/v2"

	"github.com/beganov/gingonicserver/internal/domain/core/gameConst"
)

func NewDeck() []Card {
	deck := make([]Card, gameConst.DeckSize)
	for i := range deck {
		deck[i].Id = i
	}
	delta := (gameConst.MaxValue - gameConst.MinValue)
	for i := gameConst.MinValue; i < gameConst.MaxValue; i++ {
		deck[i].Val = i
		deck[i+delta].Val = i
		deck[i+delta*2].Val = i
		deck[i+delta*3].Val = i
	}
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}

func HandInitialization(maxPlayerCount int, deck []Card) ([][]Card, [][]Card, []Card) {
	hands := make([][]Card, 0, maxPlayerCount)
	closeds := make([][]Card, 0, maxPlayerCount)
	for i := 0; i < maxPlayerCount; i++ {
		newSlice := make([]Card, gameConst.PackSize)
		copy(newSlice, deck[(i*3)*gameConst.PackSize:(i*3+1)*gameConst.PackSize])
		closeds = append(closeds, newSlice)
		newSlice2 := make([]Card, gameConst.PackSize*2)
		copy(newSlice2, deck[(i*3+1)*gameConst.PackSize:(i+1)*gameConst.PackSize*3])
		hands = append(hands, newSlice2)
	}
	deck = deck[maxPlayerCount*gameConst.PackSize*3:]

	return hands, closeds, deck
}

func OpenedsInitialization(maxPlayerCount, realPlayerCount int, hands [][]Card) ([][]Card, [][]Card) {
	openeds := make([][]Card, maxPlayerCount)
	for i := 0; i < maxPlayerCount-realPlayerCount; i++ {
		SortCard(hands[i])
		newSlice3 := make([]Card, gameConst.PackSize)
		copy(newSlice3, hands[i][gameConst.PackSize:])
		openeds[i] = newSlice3
		hands[i] = hands[i][:gameConst.PackSize]
	}
	return hands, openeds

}
