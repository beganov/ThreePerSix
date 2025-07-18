package game

import (
	"fmt"

	"github.com/beganov/gingonicserver/internal/card"
)

func ChannelsInit(Players map[int]int) (map[int]int, map[int]chan int) {
	lenPlayers := len(Players)
	idMap := make(map[int]int, lenPlayers)
	ch := make(map[int]chan int, lenPlayers)
	for i := range Players {
		ch[i] = make(chan int, 1)
		idMap[i]++
	}
	return idMap, ch
}

func keyValueReverse(idMap map[int]int) map[int]int {
	reverseIdMap := make(map[int]int, len(idMap))
	for i, j := range idMap {
		reverseIdMap[j] = i
	}
	return reverseIdMap
}

func ArraytoMap(array []int, resMap map[int]int) map[int]int {
	j := 0
	for i := range resMap {
		resMap[i] = array[j]
		j++
	}
	return resMap
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

func (g *GameState) String() string {
	return fmt.Sprintf("Iam %d, Hands: %v\n", g.IdMap, g.Hands)
}
