package game

import (
	"fmt"
	"sort"
)

func (g *GameState) Move(playerId int, playerMove int) *GameState {
	go func() { g.ch[playerId] <- playerMove }()
	return g
}

func (g *GameState) LeaveGame(playerId int) {
	g.Lock()
	defer g.Unlock()
	delete(g.Alsoiamind, g.Iamind[playerId])
	delete(g.Iamindalso, g.Iamind[playerId])
	g.ch[playerId] <- g.Hands[g.Iamind[playerId]][0].Val
	g.ch[playerId] <- g.Hands[g.Iamind[playerId]][0].Val
	g.ch[playerId] <- g.Hands[g.Iamind[playerId]][0].Val
	g.ch[playerId] <- g.Hands[g.Iamind[playerId]][0].Val
	delete(g.Iamind, playerId)
	delete(g.ch, playerId)

}

func (g GameState) String() string {
	return fmt.Sprintf("Iam %d, Hands: %v\n", g.Iamind, g.Hands)
}

func (g *GameState) Game() {
	c := 0
	g.Out = make([]Card, 0, cardQuantity)
	istake := false
	var card int
	var flag bool
	var counter int
	for counter < g.MaxPlayerCount {
		c++
		counter = 1
		for i := 0; i < g.MaxPlayerCount; i++ {
			flag = false
			card = -2
			for !flag {
				if len(g.Hands[i]) == 0 && len(g.Closeds[i]) == 0 {
					counter++
					break
				}
				sortCard(g.Hands[i])
				_, ok := g.Iamindalso[i]
				if ok {
					fmt.Println(1488)
				}
				outer(c, i, g.Out, g.Openeds, g.Closeds, g.Hands)
				g.Hands[i], g.Out, card, flag, istake = GiveCardLogic(g.Hands[i], g.Out, card, i, ok, flag, istake, g.ch[g.Iamindalso[i]])
				g.Deck, g.Hands[i], g.Openeds[i], g.Closeds[i] = TakeCard(g.Deck, g.Hands[i], g.Openeds[i], g.Closeds[i], istake)

			}
		}

	}
	fmt.Println("GameEnd")
}

func outer(c, i int, Out []Card, Openeds, Closeds, allHands [][]Card) {
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

func GiveCardLogic(Hands, Out []Card, card, i int, iamindFlag, flag, istake bool, ch <-chan int) ([]Card, []Card, int, bool, bool) {
	if len(Hands) != 0 {
		if card > -1 {
			Out, Hands, flag = ReGiveCard(Out, Hands, card, iamindFlag, ch)
			istake = !flag
		}
		if card == -2 {
			Out, Hands, istake, card = GiveCard(Out, Hands, iamindFlag, ch)
			flag = !istake
		}
		if len(Out) > 0 {
			if Out[len(Out)-1].Val == 0 || Out[len(Out)-1].Val == 10 {
				Out = Out[:0]
				flag = false
				card = -2
			}
		}
		if len(Out) >= 4 && Out[len(Out)-1].Val == Out[len(Out)-2].Val && Out[len(Out)-2].Val == Out[len(Out)-3].Val && Out[len(Out)-3].Val == Out[len(Out)-4].Val {
			Out = Out[:0]
			flag = false
			card = -2
		}
	}
	return Hands, Out, card, flag, istake
}

func TakeCard(Deck, Hands, Openeds, Closeds []Card, istake bool) ([]Card, []Card, []Card, []Card) {
	if len(Deck) == 0 && len(Openeds) == 0 && len(Hands) == 0 && len(Closeds) != 0 {
		Hands, Closeds = DecksUpdate(Hands, Closeds, 0)
	}
	if len(Deck) == 0 && len(Hands) == 0 {
		Hands = Openeds
		Openeds = Openeds[:0]
	}
	if len(Hands) < packSize && len(Deck) > 0 && istake {
		Hands, Deck = DecksUpdate(Hands, Deck, 0)
	}
	return Deck, Hands, Openeds, Closeds
}

func GiveCard(Out, Hands []Card, isAm bool, ch <-chan int) ([]Card, []Card, bool, int) {
	input := MaxValue
	// if isAm {
	// 	isAm = false
	// }
	if isAm {
		input = <-ch

	}
	if len(Out) == 0 {
		if isAm {
			for j, i := range Hands {
				if i.Val == input {
					Out, Hands = DecksUpdate(Out, Hands, j)
					return Out, Hands, true, Out[0].Val
				}
			}
		} else {
			Out, Hands = DecksUpdate(Out, Hands, 0)
			return Out, Hands, true, Out[0].Val
		}
	} else {
		for j, i := range Hands {
			if !isAm || isAm && i.Val == input {
				if Out[len(Out)-1].Val == 7 {
					if isSpecial(i.Val) || i.Val <= Out[len(Out)-1].Val {
						Out, Hands = DecksUpdate(Out, Hands, j)
						return Out, Hands, true, i.Val
					}
				} else {
					if isSpecial(i.Val) || i.Val >= Out[len(Out)-1].Val {
						Out, Hands = DecksUpdate(Out, Hands, j)
						return Out, Hands, true, i.Val
					}
				}
			}
		}
	}
	Hands = append(Hands, Out...)
	Out = Out[:0]
	return Out, Hands, false, -1
}

func ReGiveCard(Out, Hands []Card, Value int, isIam bool, ch <-chan int) ([]Card, []Card, bool) {
	// if isIam {
	// 	isIam = false
	// }
	if isIam {
		var input int
		//fmt.Scan(&input)
		input = <-ch
		if input != Value {
			return Out, Hands, true
		}
	}
	for j, i := range Hands {
		if i.Val == Value {
			Out, Hands = DecksUpdate(Out, Hands, j)
			return Out, Hands, false
		}
	}
	return Out, Hands, true

}

func cardDelete(arr []Card, index int) []Card {
	arr[index] = arr[len(arr)-1]
	arr = arr[:len(arr)-1]
	return arr
}

func DecksUpdate(arr, arr2 []Card, index int) ([]Card, []Card) {
	arr = append(arr, arr2[index])
	arr2 = cardDelete(arr2, index)
	return arr, arr2
}

func sortCard(Hands []Card) {
	sort.Slice(Hands, func(k, l int) bool {
		if isSpecial(Hands[k].Val) {
			return false
		}
		if isSpecial(Hands[l].Val) {
			return true
		}
		return Hands[k].Val < Hands[l].Val
	})
}

func isSpecial(Val int) bool {
	return Val == 2 || Val == 10 || Val == 0
}
