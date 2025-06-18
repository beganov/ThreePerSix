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

type GameState struct {
	deck        []Card
	out         []Card
	hands       [][]Card
	openeds     [][]Card
	closeds     [][]Card
	playerCount int
	iamind      int
}

func main() {
	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	var g GameState
	// 	g.Initialization()
	// 	g.Game()
	// }()
	// wg.Wait()
	var g GameState
	g.Initialization()
	g.Game()
}

func (g GameState) String() string {
	return fmt.Sprintf("Iam %d, hands: %v\n", g.iamind, g.hands)
}

func (g *GameState) Initialization() {
	g.deck = make([]Card, cardQuantity)
	g.openeds = make([][]Card, 0, g.playerCount)
	g.closeds = make([][]Card, 0, g.playerCount)
	g.hands = make([][]Card, 0, g.playerCount)
	fmt.Scan(&g.playerCount) //g.playerCount = 2 + rand.IntN(4) //
	g.DeckInitialization()
	g.HandInitialization()
	g.PlayerInitialization()
	g.iamind = 1 + rand.IntN(g.playerCount-1)
	g.Orderer()
}

func (g *GameState) Orderer() {
	min := MaxValue
	mini := 0
	for i := range g.hands {
		sortCard(g.hands[i])
		if min > g.hands[i][0].val {
			min = g.hands[i][0].val
			mini = i
		}
	}

	if mini != 0 {
		g.hands[0], g.hands[mini] = g.hands[mini], g.hands[0]
		g.openeds[0], g.openeds[mini] = g.openeds[mini], g.openeds[0]
	}
	if mini == len(g.hands)-1 {
		g.iamind = 0
	}

	fmt.Print(g.hands)
	if g.iamind != len(g.hands)-1 && g.iamind != 0 {
		g.hands[g.iamind], g.hands[len(g.hands)-1] = g.hands[len(g.hands)-1], g.hands[g.iamind]
		g.openeds[g.iamind], g.openeds[len(g.hands)-1] = g.openeds[len(g.hands)-1], g.openeds[g.iamind]
	}
}

func (g *GameState) DeckInitialization() {
	for i := range g.deck {
		g.deck[i].id = i
	}
	delta := (MaxValue - MinValue)
	for i := MinValue; i < MaxValue; i++ {
		g.deck[i].val = i
		g.deck[i+delta].val = i
		g.deck[i+delta*2].val = i
		g.deck[i+delta*3].val = i
	}
	rand.Shuffle(len(g.deck), func(i, j int) {
		g.deck[i], g.deck[j] = g.deck[j], g.deck[i]
	})
}

func (g *GameState) PlayerInitialization() {
	z := 0
	var openedShoosen int
	newSlice4 := make([]Card, 0, packSize)
	for z != packSize {
		fmt.Println(g.hands[len(g.hands)-1])
		fmt.Scan(&openedShoosen) //openedShoosen = g.hands[len(g.hands)-1][0].val //
		for i := range g.hands[len(g.hands)-1] {
			if g.hands[len(g.hands)-1][i].val == openedShoosen {
				z++
				newSlice4, g.hands[len(g.hands)-1] = decksUpdate(newSlice4, g.hands[len(g.hands)-1], i)
				break
			}
		}
	}
	g.openeds = append(g.openeds, newSlice4)
}

func (g *GameState) HandInitialization() {
	for i := 0; i < g.playerCount; i++ {
		newSlice := make([]Card, packSize)
		copy(newSlice, g.deck[(i*3)*packSize:(i*3+1)*packSize])
		g.closeds = append(g.closeds, newSlice)
		newSlice2 := make([]Card, packSize*2)
		copy(newSlice2, g.deck[(i*3+1)*packSize:(i+1)*packSize*3])
		g.hands = append(g.hands, newSlice2)
	}
	g.deck = g.deck[g.playerCount*packSize*3:]
	for i := 0; i < g.playerCount-1; i++ {
		sortCard(g.hands[i])
		newSlice3 := make([]Card, packSize)
		copy(newSlice3, g.hands[i][packSize:])
		g.openeds = append(g.openeds, newSlice3)
		g.hands[i] = g.hands[i][:packSize]
	}
}

func (g *GameState) Game() {
	c := 0
	g.out = make([]Card, 0, cardQuantity)
	istake := false
	var card int
	var flag bool
	var counter int
	for counter < g.playerCount {
		c++
		counter = 1
		for i := 0; i < g.playerCount; i++ {
			flag = false
			card = -2
			time.Sleep(time.Second / 2)
			for !flag {
				if len(g.hands[i]) == 0 && len(g.closeds[i]) == 0 {
					counter++
					break
				}
				sortCard(g.hands[i])
				outer(c, i, g.hands[g.iamind], g.out, g.openeds, g.closeds, g.hands)
				g.hands[i], g.out, card, flag, istake = GiveCardLogic(g.hands[i], g.out, card, i, g.iamind, flag, istake)
				g.deck, g.hands[i], g.openeds[i], g.closeds[i] = TakeCard(g.deck, g.hands[i], g.openeds[i], g.closeds[i], istake)

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
	// if isAm {
	// 	isAm = false
	// }
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
	// if isIam {
	// 	isIam = false
	// }
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
