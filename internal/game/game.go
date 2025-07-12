package game

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"time"
)

type GameState struct {
	Deck           []Card   `json:"deck,omitempty"`
	Out            []Card   `json:"out,omitempty"`
	Hands          [][]Card `json:"hands,omitempty"`
	Openeds        [][]Card `json:"openeds,omitempty"`
	Closeds        [][]Card `json:"closeds,omitempty"`
	MaxPlayerCount int      `json:"maxPlayerCount,omitempty"`
	Iamind         int      `json:"iamind,omitempty"`
	ch             chan int
}

func (g *GameState) Move(playerId int, playerMove int) *GameState {
	go func() { g.ch <- playerMove }()
	return g
}

func (g *GameState) LeaveGame(playerId int) {

}

func (g *GameState) StartGame(MaxPlayerCount int) *GameState {
	g.MaxPlayerCount = MaxPlayerCount
	g.ch = make(chan int, 1)
	g.SafeInitialization()
	go func() {
		g.Initialization()
		g.Game()
	}()
	return g
}

const MaxValue = 15
const MinValue = 2
const cardQuantity = 54
const packSize = 3

type Card struct {
	Id  int `json:"id"`
	Val int `json:"val"`
}

func (g GameState) String() string {
	return fmt.Sprintf("Iam %d, Hands: %v\n", g.Iamind, g.Hands)
}

func (g *GameState) SafeInitialization() {
	g.Deck = make([]Card, cardQuantity)
	g.Openeds = make([][]Card, 0, g.MaxPlayerCount)
	g.Closeds = make([][]Card, 0, g.MaxPlayerCount)
	g.Hands = make([][]Card, 0, g.MaxPlayerCount)
	g.DeckInitialization()
	g.HandInitialization()
}

func (g *GameState) Initialization() {
	g.PlayerInitialization()
	g.Iamind = 1 + rand.IntN(g.MaxPlayerCount-1)
	g.Orderer()
}

func (g *GameState) Orderer() {
	min := MaxValue
	mini := 0
	for i := range g.Hands {
		sortCard(g.Hands[i])
		if min > g.Hands[i][0].Val {
			min = g.Hands[i][0].Val
			mini = i
		}
	}

	if mini != 0 {
		g.Hands[0], g.Hands[mini] = g.Hands[mini], g.Hands[0]
		g.Openeds[0], g.Openeds[mini] = g.Openeds[mini], g.Openeds[0]
	}
	if mini == len(g.Hands)-1 {
		g.Iamind = 0
	}

	fmt.Print(g.Hands)
	if g.Iamind != len(g.Hands)-1 && g.Iamind != 0 {
		g.Hands[g.Iamind], g.Hands[len(g.Hands)-1] = g.Hands[len(g.Hands)-1], g.Hands[g.Iamind]
		g.Openeds[g.Iamind], g.Openeds[len(g.Hands)-1] = g.Openeds[len(g.Hands)-1], g.Openeds[g.Iamind]
	}
}

func (g *GameState) DeckInitialization() {
	for i := range g.Deck {
		g.Deck[i].Id = i
	}
	delta := (MaxValue - MinValue)
	for i := MinValue; i < MaxValue; i++ {
		g.Deck[i].Val = i
		g.Deck[i+delta].Val = i
		g.Deck[i+delta*2].Val = i
		g.Deck[i+delta*3].Val = i
	}
	rand.Shuffle(len(g.Deck), func(i, j int) {
		g.Deck[i], g.Deck[j] = g.Deck[j], g.Deck[i]
	})
}

func (g *GameState) PlayerInitialization() {
	z := 0
	var Openedshoosen int
	newSlice4 := make([]Card, 0, packSize)
	for z != packSize {
		fmt.Println(g.Hands[len(g.Hands)-1])
		Openedshoosen = <-g.ch
		//fmt.Scan(&Openedshoosen) //
		//Openedshoosen = g.Hands[len(g.Hands)-1][0].Val //
		for i := range g.Hands[len(g.Hands)-1] {
			if g.Hands[len(g.Hands)-1][i].Val == Openedshoosen {
				z++
				newSlice4, g.Hands[len(g.Hands)-1] = DecksUpdate(newSlice4, g.Hands[len(g.Hands)-1], i)
				break
			}
		}
	}
	g.Openeds = append(g.Openeds, newSlice4)
}

func (g *GameState) HandInitialization() {
	for i := 0; i < g.MaxPlayerCount; i++ {
		newSlice := make([]Card, packSize)
		copy(newSlice, g.Deck[(i*3)*packSize:(i*3+1)*packSize])
		g.Closeds = append(g.Closeds, newSlice)
		newSlice2 := make([]Card, packSize*2)
		copy(newSlice2, g.Deck[(i*3+1)*packSize:(i+1)*packSize*3])
		g.Hands = append(g.Hands, newSlice2)
	}
	g.Deck = g.Deck[g.MaxPlayerCount*packSize*3:]
	for i := 0; i < g.MaxPlayerCount-1; i++ {
		sortCard(g.Hands[i])
		newSlice3 := make([]Card, packSize)
		copy(newSlice3, g.Hands[i][packSize:])
		g.Openeds = append(g.Openeds, newSlice3)
		g.Hands[i] = g.Hands[i][:packSize]
	}
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
			time.Sleep(time.Second / 2)
			for !flag {
				if len(g.Hands[i]) == 0 && len(g.Closeds[i]) == 0 {
					counter++
					break
				}
				sortCard(g.Hands[i])
				outer(c, i, g.Hands[g.Iamind], g.Out, g.Openeds, g.Closeds, g.Hands)
				g.Hands[i], g.Out, card, flag, istake = GiveCardLogic(g.Hands[i], g.Out, card, i, g.Iamind, flag, istake, g.ch)
				g.Deck, g.Hands[i], g.Openeds[i], g.Closeds[i] = TakeCard(g.Deck, g.Hands[i], g.Openeds[i], g.Closeds[i], istake)

			}
		}

	}
}

func outer(c, i int, Hands, Out []Card, Openeds, Closeds, allHands [][]Card) {
	allHandsLens := make([]int, len(allHands))
	allClosedsLens := make([]int, len(allHands))
	for i := range allHands {
		allHandsLens[i] = len(allHands[i])
		allClosedsLens[i] = len(Closeds[i])
	}
	fmt.Printf("Turn:  %d\n", c)
	fmt.Printf("Player:  %d\n", i)
	fmt.Printf("Turn %d \n player %d hand: %v \n table %v \n Openeds %v \n", c, i, Hands, Out, Openeds)
	fmt.Printf("Len Closeds %d, Len Hands %d", allClosedsLens, allHandsLens)
}

func GiveCardLogic(Hands, Out []Card, card, i, Iamind int, flag, istake bool, ch <-chan int) ([]Card, []Card, int, bool, bool) {
	if len(Hands) != 0 {
		if card > -1 {
			Out, Hands, flag = ReGiveCard(Out, Hands, card, i == Iamind, ch)
			istake = !flag
		}
		if card == -2 {
			Out, Hands, istake, card = GiveCard(Out, Hands, i == Iamind, ch)
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
		fmt.Println(Hands)
		//fmt.Scan(&input)
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
