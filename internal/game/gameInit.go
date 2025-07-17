package game

import (
	"fmt"
	"math/rand/v2"
	"sync"

	"github.com/beganov/gingonicserver/internal/card"
	"github.com/beganov/gingonicserver/internal/gameConst"
	"github.com/beganov/gingonicserver/internal/placement"
)

func (g *GameState) PreInitialization(MaxPlayerCount int, Players map[int]int) {
	g.Iamind, g.Alsoiamind, g.ch = ChannelsInit(Players)
	g.MaxPlayerCount = MaxPlayerCount
	g.Deck = card.NewDeck()
	g.Hands, g.Closeds, g.Deck = card.HandInitialization(MaxPlayerCount, g.Deck)
	g.Hands, g.Openeds = card.OpenedsInitialization(MaxPlayerCount, len(g.Iamind), g.Hands)
}

func ChannelsInit(Players map[int]int) (map[int]int, map[int]int, map[int]chan int) {
	lenPlayers := len(Players)
	iamind := make(map[int]int, lenPlayers)
	alsoiamind := make(map[int]int, lenPlayers)
	ch := make(map[int]chan int, lenPlayers)
	for i := range Players {
		ch[i] = make(chan int, 1)
		iamind[i]++
	}
	return iamind, alsoiamind, ch
}

func (g *GameState) Initialization() {
	g.PlayerInitialization()
	g.Hands, g.Openeds, g.Alsoiamind = placement.ShufflePlayer(g.Hands, g.Openeds, g.Alsoiamind)               //Расставляет игроков по случайным позициям
	g.Hands, g.Openeds, g.Iamind, g.Alsoiamind = placement.Orderer(g.Hands, g.Openeds, g.Iamind, g.Alsoiamind) //Передает первый ход игроку с тройкой
	g.Iamindalso = keyValueReverse(g.Iamind)
}

func keyValueReverse(iamind map[int]int) map[int]int {
	iamindalso := make(map[int]int, len(iamind))
	for i, j := range iamind {
		iamindalso[j] = i
	}
	return iamindalso
}

func (g *GameState) PlayerInitialization() {
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
	j := 0
	for i := range g.Iamind {
		g.Iamind[i] = EndOfArray[j]
		j++
	}
	for i := range g.Iamind {
		wg.Add(1)
		go func() {
			defer wg.Done()
			g.OpenedsPlayerInitialization(i, g.Iamind[i])
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

}

func (g *GameState) OpenedsPlayerInitialization(playerId, k int) {
	z := 0
	var Openedshoosen int
	newSlice4 := make([]card.Card, 0, gameConst.PackSize)
	for z != gameConst.PackSize {
		fmt.Println(g.Hands[k])
		Openedshoosen = <-g.ch[playerId]
		//fmt.Scan(&Openedshoosen) //
		//Openedshoosen = g.Hands[len(g.Hands)-1][0].Val //\
		for i := range g.Hands[k] {
			if g.Hands[k][i].Val == Openedshoosen {
				z++
				newSlice4, g.Hands[k] = card.DecksUpdate(newSlice4, g.Hands[k], i)
				break
			}
		}
	}
	g.Openeds = append(g.Openeds, newSlice4)
}
