package game

import (
	"fmt"
	"sync"

	"github.com/beganov/gingonicserver/internal/domain/core/card"
	"github.com/beganov/gingonicserver/internal/domain/core/gameConst"
	"github.com/beganov/gingonicserver/internal/domain/core/placement"
)

func (g *GameState) PreInitialization(maxPlayerCount int, Players map[int]int, end GameEndHandler) {
	g.Turn = 0
	g.PlayerNow = 0
	g.handler = end
	g.IdMap, g.ch = ChannelsInit(Players)
	g.Deck = card.NewDeck()
	g.Hands, g.Closeds, g.Deck = card.HandInitialization(maxPlayerCount, g.Deck)
	g.Hands, g.Openeds = card.OpenedsInitialization(maxPlayerCount, len(g.IdMap), g.Hands)
}

func (g *GameState) Initialization(maxPlayerCount int) {
	orderMap := g.PlayerInitialization(maxPlayerCount)
	g.Hands, g.Openeds = placement.ShufflePlayer(g.Hands, g.Openeds, orderMap)   //Расставляет игроков по случайным позициям
	g.Hands, g.Openeds, g.IdMap = placement.Orderer(g.Hands, g.Openeds, g.IdMap) //Передает первый ход игроку с тройкой
	g.Closeds = placement.LeaveCheck(g.Hands, g.Closeds)
	g.ReverceIdMap = keyValueReverse(g.IdMap)
}

func (g *GameState) PlayerInitialization(maxPlayerCount int) map[int]int {
	shuffleArr := placement.NewPlacementArray(maxPlayerCount, len(g.IdMap))

	EndOfArray := shuffleArr[len(shuffleArr)-len(g.IdMap):]
	g.IdMap = ArraytoMap(EndOfArray, g.IdMap)

	var wg sync.WaitGroup
	for i := range g.IdMap {
		wg.Add(1)
		go func() {
			defer wg.Done()
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
		card.SortCard(g.Hands[k])
		Openedshoosen = <-g.ch[playerId]
		if Openedshoosen == gameConst.LeaveGameCode {
			z = gameConst.PackSize
			break
		}
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

func (g *GameState) Game(maxPlayerCount int) {
	c := 0
	g.Out = make([]card.Card, 0, gameConst.DeckSize)
	istake := false
	var cardState int
	var flag bool
	var counter int
	for counter < maxPlayerCount {
		c++
		counter = 1
		for i := 0; i < maxPlayerCount; i++ {
			g.IsMoved = false
			flag = false
			cardState = gameConst.StartCardState
			for !flag {
				if len(g.Hands[i]) == 0 && len(g.Closeds[i]) == 0 {
					counter++
					//delete(g.IdMap, g.ReverceIdMap[i])
					//delete(g.ch, g.ReverceIdMap[i])
					//delete(g.ReverceIdMap, i)
					break
				}
				card.SortCard(g.Hands[i])
				_, ok := g.ReverceIdMap[i]
				g.Turn = c
				g.PlayerNow = i
				outer(c, i, g.Out, g.Openeds, g.Closeds, g.Hands)
				g.Hands[i], g.Out, cardState, flag, istake, g.IsMoved = card.GiveCardLogic(g.Hands[i], g.Out, cardState, i, ok, flag, istake, g.ch[g.ReverceIdMap[i]])
				g.Deck, g.Hands[i], g.Openeds[i], g.Closeds[i] = card.TakeCard(g.Deck, g.Hands[i], g.Openeds[i], g.Closeds[i], istake)

			}
		}

	}
	fmt.Println("GameEnd")
	g.Hands = [][]card.Card{}
	g.Openeds = [][]card.Card{}
	g.Closeds = [][]card.Card{}
	g.handler.OnGameEnd()
}
