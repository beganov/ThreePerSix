package game

import (
	"fmt"
	"math/rand/v2"
	"sync"
)

func (g *GameState) SafeInitialization(MaxPlayerCount int, Players map[int]int) {
	g.Iamind = make(map[int]int, len(Players))
	g.Alsoiamind = make(map[int]int, len(Players))
	g.ch = make(map[int]chan int, len(Players))
	for i := range Players {
		g.ch[i] = make(chan int, 1)
		g.Iamind[i]++
	}
	g.MaxPlayerCount = MaxPlayerCount
	g.Deck = make([]Card, cardQuantity)
	g.Openeds = make([][]Card, 0, g.MaxPlayerCount)
	g.Closeds = make([][]Card, 0, g.MaxPlayerCount)
	g.Hands = make([][]Card, 0, g.MaxPlayerCount)
	g.DeckInitialization()
	g.HandInitialization()
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
	for i := 0; i < g.MaxPlayerCount-len(g.Iamind); i++ {
		sortCard(g.Hands[i])
		newSlice3 := make([]Card, packSize)
		copy(newSlice3, g.Hands[i][packSize:])
		g.Openeds = append(g.Openeds, newSlice3)
		g.Hands[i] = g.Hands[i][:packSize]
	}
}

func (g *GameState) Initialization() {
	var wg sync.WaitGroup
	var shuffleArr []int
	for i := 0; i < g.MaxPlayerCount; i++ {
		if i != 0 {
			shuffleArr = append(shuffleArr, i)
		} else {
			if len(g.Iamind) == g.MaxPlayerCount {
				shuffleArr = append(shuffleArr, i)
			}
		}
	}
	EndOfArray := shuffleArr[len(shuffleArr)-len(g.Iamind):]
	fmt.Println(EndOfArray)
	j := 0
	for i, k := range g.Iamind {
		g.Iamind[i] = EndOfArray[j]
		fmt.Println(EndOfArray[j], i, k, g.Iamind[i])
		j++
	}
	for i, k := range g.Iamind {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(i, k, g.Iamind[i])
			g.PlayerInitialization(i, g.Iamind[i])
		}()

	}
	wg.Wait()
	rand.Shuffle(len(shuffleArr), func(i, j int) {
		shuffleArr[i], shuffleArr[j] = shuffleArr[j], shuffleArr[i]
	})
	j = 0
	for i := range g.Iamind {
		g.Alsoiamind[g.Iamind[i]] = shuffleArr[j]
		g.Iamind[i] = shuffleArr[j]
		j++
	}

	//g.Iamind = 1 + rand.IntN(g.MaxPlayerCount-1)
	g.ShufflePlayer() //Расставляет игроков по случайным позициям
	g.Orderer()       //Передает первый ход игроку с тройкой
}

func (g *GameState) PlayerInitialization(playerId, k int) {
	z := 0
	var Openedshoosen int
	newSlice4 := make([]Card, 0, packSize)
	for z != packSize {
		fmt.Println(g.Hands[k])
		Openedshoosen = <-g.ch[playerId]
		//fmt.Scan(&Openedshoosen) //
		//Openedshoosen = g.Hands[len(g.Hands)-1][0].Val //\
		for i := range g.Hands[k] {
			if g.Hands[k][i].Val == Openedshoosen {
				z++
				newSlice4, g.Hands[k] = DecksUpdate(newSlice4, g.Hands[k], i)
				break
			}
		}
	}
	g.Openeds = append(g.Openeds, newSlice4)
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

	fmt.Print(g.Hands)
	for i, j := range g.Iamind {
		if j == mini {
			g.Iamind[i] = 0
			g.Alsoiamind[i] = 0
		}
	}

	g.Iamindalso = make(map[int]int, len(g.Iamind))
	for i, j := range g.Iamind {
		g.Iamindalso[j] = i
	}
}

func (g *GameState) ShufflePlayer() {

	for i := range g.Hands {
		if _, ok := g.Alsoiamind[i]; !ok {
			g.Alsoiamind[i] = i
		}
	}
	flag := true
	for flag {
		flag = false
		for i := range g.Hands {
			if i != g.Alsoiamind[i] {
				flag = true
				g.Hands[i], g.Hands[g.Alsoiamind[i]] = g.Hands[g.Alsoiamind[i]], g.Hands[i]
				g.Openeds[i], g.Openeds[g.Alsoiamind[i]] = g.Openeds[g.Alsoiamind[i]], g.Openeds[i]
				g.Alsoiamind[g.Alsoiamind[i]] = g.Alsoiamind[i]
				g.Alsoiamind[i] = i
			}
		}
		fmt.Println(g.Hands)
	}

}
