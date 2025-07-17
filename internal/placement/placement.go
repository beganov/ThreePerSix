package placement

import (
	"github.com/beganov/gingonicserver/internal/card"
	"github.com/beganov/gingonicserver/internal/gameConst"
)

func Orderer(hands, openeds [][]card.Card, iamind, alsoiamind map[int]int) ([][]card.Card, [][]card.Card, map[int]int, map[int]int) {
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
	for i, j := range iamind {
		if j == mini {
			iamind[i] = 0
			alsoiamind[i] = 0
		}
	}
	return hands, openeds, iamind, alsoiamind
}

func ShufflePlayer(hands, openeds [][]card.Card, alsoiamind map[int]int) ([][]card.Card, [][]card.Card, map[int]int) {

	for i := range hands {
		if _, ok := alsoiamind[i]; !ok {
			alsoiamind[i] = i
		}
	}
	flag := true
	for flag {
		flag = false
		for i := range hands {
			if i != alsoiamind[i] {
				flag = true
				hands[i], hands[alsoiamind[i]] = hands[alsoiamind[i]], hands[i]
				openeds[i], openeds[alsoiamind[i]] = openeds[alsoiamind[i]], openeds[i]
				alsoiamind[alsoiamind[i]] = alsoiamind[i]
				alsoiamind[i] = i
			}
		}
	}
	return hands, openeds, alsoiamind
}
