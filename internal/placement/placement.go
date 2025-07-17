package placement

import (
	"math/rand/v2"

	"github.com/beganov/gingonicserver/internal/card"
	"github.com/beganov/gingonicserver/internal/gameConst"
)

func Orderer(hands, openeds [][]card.Card, idMap map[int]int) ([][]card.Card, [][]card.Card, map[int]int) {
	min := gameConst.MaxValue
	mini := 0
	for i := range hands {
		card.SortCard(hands[i])
		if min > hands[i][0].Val {
			min = hands[i][0].Val
			mini = i
		}
	}

	if mini != 0 {
		hands[0], hands[mini] = hands[mini], hands[0]
		openeds[0], openeds[mini] = openeds[mini], openeds[0]
	}
	placement := 0
	for i, j := range idMap {
		if j == 0 {
			placement = i
			break
		}
	}
	for i, j := range idMap {
		if j == mini {
			if _, ok := idMap[placement]; ok {
				idMap[placement] = idMap[i]
			}
			idMap[i] = 0
			break
		}
	}

	return hands, openeds, idMap
}

func ShufflePlayer(hands, openeds [][]card.Card, orderMap map[int]int) ([][]card.Card, [][]card.Card) {
	for i := range hands {
		if _, ok := orderMap[i]; !ok {
			orderMap[i] = i
		}
	}
	flag := true
	for flag {
		flag = false
		for i := range hands {
			if i != orderMap[i] {
				flag = true
				hands[i], hands[orderMap[i]] = hands[orderMap[i]], hands[i]
				openeds[i], openeds[orderMap[i]] = openeds[orderMap[i]], openeds[i]
				orderMap[orderMap[i]] = orderMap[i]
				orderMap[i] = i
			}
		}
	}
	return hands, openeds
}

func TakeRandomPlacement(shuffleArr []int, idMap map[int]int) (map[int]int, map[int]int) {
	orderMap := make(map[int]int, len(idMap))
	rand.Shuffle(len(shuffleArr), func(i, j int) {
		shuffleArr[i], shuffleArr[j] = shuffleArr[j], shuffleArr[i]
	})
	j := 0
	for i := range idMap {
		orderMap[idMap[i]] = shuffleArr[j]
		idMap[i] = shuffleArr[j]
		j++
	}
	return idMap, orderMap
}

func NewPlacementArray(maxPlayerCount, realPlayerCount int) []int {
	shuffleArr := make([]int, 0, maxPlayerCount)
	for i := 0; i < maxPlayerCount; i++ {
		if i != 0 {
			shuffleArr = append(shuffleArr, i)
		} else {
			if realPlayerCount == maxPlayerCount {
				shuffleArr = append(shuffleArr, i)
			}
		}
	}
	return shuffleArr
}
