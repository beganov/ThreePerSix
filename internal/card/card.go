package card

import (
	"math/rand/v2"
	"sort"

	"github.com/beganov/gingonicserver/internal/gameConst"
)

type Card struct {
	Id  int `json:"id"`
	Val int `json:"val"`
}

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

func SortCard(Hands []Card) {
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

func GiveCardLogic(Hands, Out []Card, cardState, i int, iamindFlag, flag, istake bool, ch <-chan int) ([]Card, []Card, int, bool, bool) {
	if len(Hands) != 0 {
		if cardState > -1 {
			Out, Hands, flag = ReGiveCard(Out, Hands, cardState, iamindFlag, ch)
			istake = !flag
		}
		if cardState == -2 {
			Out, Hands, istake, cardState = GiveCard(Out, Hands, iamindFlag, ch)
			flag = !istake
		}
		if len(Out) > 0 {
			if Out[len(Out)-1].Val == 0 || Out[len(Out)-1].Val == 10 {
				Out = Out[:0]
				flag = false
				cardState = -2
			}
		}
		if len(Out) >= 4 && Out[len(Out)-1].Val == Out[len(Out)-2].Val && Out[len(Out)-2].Val == Out[len(Out)-3].Val && Out[len(Out)-3].Val == Out[len(Out)-4].Val {
			Out = Out[:0]
			flag = false
			cardState = -2
		}
	}
	return Hands, Out, cardState, flag, istake
}

func TakeCard(Deck, Hands, Openeds, Closeds []Card, istake bool) ([]Card, []Card, []Card, []Card) {
	if len(Deck) == 0 && len(Openeds) == 0 && len(Hands) == 0 && len(Closeds) != 0 {
		Hands, Closeds = DecksUpdate(Hands, Closeds, 0)
	}
	if len(Deck) == 0 && len(Hands) == 0 {
		Hands = Openeds
		Openeds = Openeds[:0]
	}
	if len(Hands) < gameConst.PackSize && len(Deck) > 0 && istake {
		Hands, Deck = DecksUpdate(Hands, Deck, 0)
	}
	return Deck, Hands, Openeds, Closeds
}

func GiveCard(Out, Hands []Card, isAm bool, ch <-chan int) ([]Card, []Card, bool, int) {
	input := gameConst.MaxValue
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
	openeds := make([][]Card, 0, maxPlayerCount)
	for i := 0; i < maxPlayerCount-realPlayerCount; i++ {
		SortCard(hands[i])
		newSlice3 := make([]Card, gameConst.PackSize)
		copy(newSlice3, hands[i][gameConst.PackSize:])
		openeds = append(openeds, newSlice3)
		hands[i] = hands[i][:gameConst.PackSize]
	}
	return hands, openeds

}
