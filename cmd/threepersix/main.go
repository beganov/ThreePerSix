package main

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"time"
)

const MaxValue = 15
const MinValue = 2
const cardQuantity = 54
const packSize = 3

type Card struct {
	id  int
	val int
}

func main() {
	playerCount, iamind, deck, hands, openeds, closeds := Initialization()
	Game(deck, hands, openeds, closeds, playerCount, iamind)
}

func Initialization() (int, int, []Card, [][]Card, [][]Card, [][]Card) {
	var playerCount int
	deck := make([]Card, cardQuantity)
	openeds := make([][]Card, 0, playerCount)
	closeds := make([][]Card, 0, playerCount)
	hands := make([][]Card, 0, playerCount)
	fmt.Scan(&playerCount)
	deck = DeckInitialization(deck)
	hands, closeds, openeds, deck = HandInitialization(hands, closeds, openeds, deck, playerCount)
	hands, openeds = PlayerInitialization(hands, openeds)
	iamind := rand.IntN(playerCount)
	hands, openeds, iamind = Orderer(iamind, playerCount, hands, openeds)
	return playerCount, iamind, deck, hands, openeds, closeds
}

func Orderer(iamind, playerCount int, hands, openeds [][]Card) ([][]Card, [][]Card, int) {
	min := MaxValue
	mini := 0
	for i := range hands {
		sortCard(hands[i])
		if min > hands[i][0].val {
			min = hands[i][0].val
			mini = i
		}
	}
	if mini != 0 {
		hands[0], hands[mini] = hands[mini], hands[0]
	}
	if mini == iamind {
		iamind = 0
	}
	if iamind != len(hands)-1 {
		hands[iamind], hands[len(hands)-1] = hands[len(hands)-1], hands[iamind]
		openeds[iamind], openeds[len(hands)-1] = openeds[len(hands)-1], openeds[iamind]
	}
	return hands, openeds, iamind
}

func DeckInitialization(deck []Card) []Card {
	for i := range deck {
		deck[i].id = i
	}
	delta := (MaxValue - MinValue)
	for i := MinValue; i < MaxValue; i++ {
		deck[i].val = i
		deck[i+delta].val = i
		deck[i+delta*2].val = i
		deck[i+delta*3].val = i
	}
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}

func PlayerInitialization(hands, openeds [][]Card) ([][]Card, [][]Card) {
	z := 0
	var openedShoosen int
	newSlice4 := make([]Card, 0, packSize)
	for z != packSize {
		fmt.Println(hands[len(hands)-1])
		fmt.Scan(&openedShoosen)
		for i := range hands[len(hands)-1] {
			if hands[len(hands)-1][i].val == openedShoosen {
				z++
				newSlice4, hands[len(hands)-1] = decksUpdate(newSlice4, hands[len(hands)-1], i)
				break
			}
		}
	}
	openeds = append(openeds, newSlice4)
	return hands, openeds
}

func HandInitialization(hands, closeds, openeds [][]Card, deck []Card, playerCount int) ([][]Card, [][]Card, [][]Card, []Card) {
	for i := 0; i < playerCount; i++ {
		newSlice := make([]Card, packSize)
		copy(newSlice, deck[i*packSize*3:i*packSize*4])
		closeds = append(closeds, newSlice)
		newSlice2 := make([]Card, packSize*2)
		copy(newSlice2, deck[i*packSize*4:(i+1)*packSize*3])
		hands = append(hands, newSlice2)
	}
	deck = deck[playerCount*9:]
	for i := 0; i < playerCount-1; i++ {
		sortCard(hands[i])
		newSlice3 := make([]Card, packSize)
		copy(newSlice3, hands[i][packSize:])
		openeds = append(openeds, newSlice3)
		hands[i] = hands[i][:packSize]
	}
	return hands, closeds, openeds, deck
}

func Game(deck []Card, hands, openeds, closeds [][]Card, playerCount, iamind int) {
	c := 0
	out := make([]Card, 0, cardQuantity)
	istake := false
	var card int
	var flag bool
	var counter int
	for counter < playerCount {
		c++
		counter = 1
		for i := 0; i < playerCount; i++ {
			flag = false
			card = -2
			time.Sleep(time.Second / 2)
			for !flag {
				if len(hands[i]) == 0 && len(closeds[i]) == 0 {
					counter++
					break
				}
				sortCard(hands[i])
				outer(c, i, hands[iamind], out, openeds, closeds, hands)
				hands[i], out, card, flag, istake = GiveCardLogic(hands[i], out, card, i, iamind, flag, istake)
				deck, hands[i], openeds[i], closeds[i] = TakeCard(deck, hands[i], openeds[i], closeds[i], istake)

			}
		}

	}
}

func outer(c, i int, hands, out []Card, openeds, closeds, allHands [][]Card) {
	allHandsLens := make([]int, len(allHands))
	allClosedsLens := make([]int, len(allHands))
	for i := range allHands {
		allHandsLens[i] = len(allHands[i])
		allClosedsLens[i] = len(closeds[i])
	}
	fmt.Printf("Turn:  %d\n", c)
	fmt.Printf("Player:  %d\n", i)
	fmt.Printf("Turn %d \n player %d hand: %v \n table %v \n openeds %v \n", c, i, hands, out, openeds)
	fmt.Printf("Len Closeds %d, Len Hands %d", allClosedsLens, allHandsLens)
}

func GiveCardLogic(hands, out []Card, card, i, iamind int, flag, istake bool) ([]Card, []Card, int, bool, bool) {
	if len(hands) != 0 {
		if card > -1 {
			out, hands, flag = ReGiveCard(out, hands, card, i == iamind)
			istake = !flag
		}
		if card == -2 {
			out, hands, istake, card = GiveCard(out, hands, i == iamind)
			flag = !istake
		}
		if len(out) > 0 {
			if out[len(out)-1].val == 0 || out[len(out)-1].val == 10 {
				out = out[:0]
				flag = false
				card = -2
			}
		}
		if len(out) >= 4 && out[len(out)-1].val == out[len(out)-2].val && out[len(out)-2].val == out[len(out)-3].val && out[len(out)-3].val == out[len(out)-4].val {
			out = out[:0]
			flag = false
			card = -2
		}
	}
	return hands, out, card, flag, istake
}

func TakeCard(deck, hands, openeds, closeds []Card, istake bool) ([]Card, []Card, []Card, []Card) {
	if len(deck) == 0 && len(openeds) == 0 && len(hands) == 0 && len(closeds) != 0 {
		hands, closeds = decksUpdate(hands, closeds, 0)
	}
	if len(deck) == 0 && len(hands) == 0 {
		hands = openeds
		openeds = openeds[:0]
	}
	if len(hands) < packSize && len(deck) > 0 && istake {
		hands, deck = decksUpdate(hands, deck, 0)
	}
	return deck, hands, openeds, closeds
}

func GiveCard(out, hands []Card, isAm bool) ([]Card, []Card, bool, int) {
	input := MaxValue
	if isAm {
		fmt.Println(hands)
		fmt.Scan(&input)

	}
	if len(out) == 0 {
		if isAm {
			for j, i := range hands {
				if i.val == input {
					out, hands = decksUpdate(out, hands, j)
					return out, hands, true, out[0].val
				}
			}
		} else {
			out, hands = decksUpdate(out, hands, 0)
			return out, hands, true, out[0].val
		}
	} else {
		for j, i := range hands {
			if !isAm || isAm && i.val == input {
				if out[len(out)-1].val == 7 {
					if isSpecial(i.val) || i.val <= out[len(out)-1].val {
						out, hands = decksUpdate(out, hands, j)
						return out, hands, true, i.val
					}
				} else {
					if isSpecial(i.val) || i.val >= out[len(out)-1].val {
						out, hands = decksUpdate(out, hands, j)
						return out, hands, true, i.val
					}
				}
			}
		}
	}
	hands = append(hands, out...)
	out = out[:0]
	return out, hands, false, -1
}

func ReGiveCard(out, hands []Card, value int, isIam bool) ([]Card, []Card, bool) {
	if isIam {
		var input int
		fmt.Scan(&input)
		if input != value {
			return out, hands, true
		}
	}
	for j, i := range hands {
		if i.val == value {
			out, hands = decksUpdate(out, hands, j)
			return out, hands, false
		}
	}
	return out, hands, true

}

func cardDelete(arr []Card, index int) []Card {
	arr[index] = arr[len(arr)-1]
	arr = arr[:len(arr)-1]
	return arr
}

func decksUpdate(arr, arr2 []Card, index int) ([]Card, []Card) {
	arr = append(arr, arr2[index])
	arr2 = cardDelete(arr2, index)
	return arr, arr2
}

func sortCard(hands []Card) {
	sort.Slice(hands, func(k, l int) bool {
		if isSpecial(hands[k].val) {
			return false
		}
		if isSpecial(hands[l].val) {
			return true
		}
		return hands[k].val < hands[l].val
	})
}

func isSpecial(val int) bool {
	return val == 2 || val == 10 || val == 0
}
