package game

import (
	"fmt"

	"github.com/beganov/gingonicserver/internal/card"
	"github.com/beganov/gingonicserver/internal/gameConst"
)

func (g *GameState) Game() {
	c := 0
	g.Out = make([]card.Card, 0, gameConst.DeckSize)
	istake := false
	var cardState int
	var flag bool
	var counter int
	for counter < g.MaxPlayerCount {
		c++
		counter = 1
		for i := 0; i < g.MaxPlayerCount; i++ {
			flag = false
			cardState = -2
			for !flag {
				if len(g.Hands[i]) == 0 && len(g.Closeds[i]) == 0 {
					counter++
					break
				}
				card.SortCard(g.Hands[i])
				_, ok := g.Iamindalso[i]
				outer(c, i, g.Out, g.Openeds, g.Closeds, g.Hands)
				g.Hands[i], g.Out, cardState, flag, istake = card.GiveCardLogic(g.Hands[i], g.Out, cardState, i, ok, flag, istake, g.ch[g.Iamindalso[i]])
				g.Deck, g.Hands[i], g.Openeds[i], g.Closeds[i] = card.TakeCard(g.Deck, g.Hands[i], g.Openeds[i], g.Closeds[i], istake)

			}
		}

	}
	fmt.Println("GameEnd")
}

func outer(c, i int, Out []card.Card, Openeds, Closeds, allHands [][]card.Card) {
	allHandsLens := make([]int, len(allHands))
	allClosedsLens := make([]int, len(allHands))
	for j := range allHands {
		allHandsLens[j] = len(allHands[j])
		allClosedsLens[j] = len(Closeds[j])
	}
	fmt.Printf("Turn:  %d\n", c)
	fmt.Printf("Player:  %d\n", i)
	fmt.Printf("Turn %d \n player %d hand: %v \n table %v \n Openeds %v \n", c, i, allHands[i], Out, Openeds)
	fmt.Printf("Len Closeds %d, Len Hands %d", allClosedsLens, allHandsLens)
}

func (g GameState) String() string {
	return fmt.Sprintf("Iam %d, Hands: %v\n", g.Iamind, g.Hands)
}
