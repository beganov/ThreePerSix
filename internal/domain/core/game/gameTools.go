package game

import (
	"fmt"

	"github.com/beganov/gingonicserver/internal/domain/core/card"
)

// Метод для инициализации мап каналов и playerId
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

// Метод для создания мапы с обратными к оригинальной мапе значениями/
// key/value -> value/key
func keyValueReverse(idMap map[int]int) map[int]int {
	reverseIdMap := make(map[int]int, len(idMap))
	for i, j := range idMap {
		reverseIdMap[j] = i
	}
	return reverseIdMap
}

// Метод для заполнения мапы значениями из массива
func ArraytoMap(array []int, resMap map[int]int) map[int]int {
	j := 0
	for i := range resMap {
		resMap[i] = array[j]
		j++
	}
	return resMap
}

// Фунцкия вывода
func outer(c, i int, Out []card.Card, Openeds, Closeds, allHands [][]card.Card, logger Logger) {

	allHandsLens := make([]int, len(allHands))
	allClosedsLens := make([]int, len(allHands))
	for j := range allHands { // Подсчитываем количество карт на руках и в закрытых у каждого игрока
		allHandsLens[j] = len(allHands[j])
		allClosedsLens[j] = len(Closeds[j])
	}
	logger.Info(). // Выводим
			Int("Turn", c).
			Int("Player", i).
			Interface("hand", allHands[i]).
			Interface("table", Out).
			Interface("Openeds", Openeds).
			Interface("LenCloseds", allClosedsLens).
			Interface("LenHands", allHandsLens).
			Msg("Game state update")
}

// Метод String для Game(не использовался)
func (g *GameState) String() string {
	return fmt.Sprintf("Iam %d, Hands: %v\n", g.IdMap, g.Hands)
}
