package game

import (
	"fmt"
	"sync"

	"github.com/beganov/gingonicserver/internal/card"
	"github.com/beganov/gingonicserver/internal/gameConst"
	"github.com/beganov/gingonicserver/internal/placement"
)

func (g *GameState) PreInitialization(MaxPlayerCount int, Players map[int]int) {
	g.IdMap, g.ch = ChannelsInit(Players)
	g.MaxPlayerCount = MaxPlayerCount
	g.Deck = card.NewDeck()
	g.Hands, g.Closeds, g.Deck = card.HandInitialization(MaxPlayerCount, g.Deck)
	g.Hands, g.Openeds = card.OpenedsInitialization(MaxPlayerCount, len(g.IdMap), g.Hands)
}

func (g *GameState) Initialization() {
	orderMap := g.PlayerInitialization()
	g.Hands, g.Openeds = placement.ShufflePlayer(g.Hands, g.Openeds, orderMap)   //Расставляет игроков по случайным позициям
	g.Hands, g.Openeds, g.IdMap = placement.Orderer(g.Hands, g.Openeds, g.IdMap) //Передает первый ход игроку с тройкой
	g.ReverceIdMap = keyValueReverse(g.IdMap)
}

func (g *GameState) PlayerInitialization() map[int]int {
	shuffleArr := placement.NewPlacementArray(g.MaxPlayerCount, len(g.IdMap))

	EndOfArray := shuffleArr[len(shuffleArr)-len(g.IdMap):]
	g.IdMap = ArraytoMap(EndOfArray, g.IdMap)

	var wg sync.WaitGroup
	for i := range g.IdMap {
		wg.Add(1)
		go func() {
			defer wg.Done()
			g.Lock()
			defer g.Unlock()
			g.OpenedsPlayerInitialization(i, g.IdMap[i])
		}()

	}
	wg.Wait()

	var orderMap map[int]int
	g.IdMap, orderMap = placement.TakeRandomPlacement(shuffleArr, g.IdMap)
	return orderMap
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
	g.Openeds[k] = newSlice4
}

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
				_, ok := g.ReverceIdMap[i]
				outer(c, i, g.Out, g.Openeds, g.Closeds, g.Hands)
				g.Hands[i], g.Out, cardState, flag, istake = card.GiveCardLogic(g.Hands[i], g.Out, cardState, i, ok, flag, istake, g.ch[g.ReverceIdMap[i]])
				g.Deck, g.Hands[i], g.Openeds[i], g.Closeds[i] = card.TakeCard(g.Deck, g.Hands[i], g.Openeds[i], g.Closeds[i], istake)

			}
		}

	}
	fmt.Println("GameEnd")
}
