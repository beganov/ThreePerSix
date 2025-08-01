package placement

import (
	"math/rand/v2"

	"github.com/beganov/gingonicserver/internal/domain/core/card"
	"github.com/beganov/gingonicserver/internal/domain/core/gameConst"
)

// Даем возможность хода игроку с наименьшей не специальной картой
func Orderer(hands, openeds [][]card.Card, idMap map[int]int) ([][]card.Card, [][]card.Card, map[int]int) {
	min := gameConst.MaxValue
	mini := 0
	for i := range hands {
		card.SortCard(hands[i]) // сортируем
		if len(hands[i]) == 0 { // скип на случай leave
			continue
		}
		if min > hands[i][0].Val { // Ищем минимальную среди рук -  она теперь на 0 позиции
			min = hands[i][0].Val
			mini = i //запоминаем индекс и значение
		}
	}

	if mini != 0 && len(hands) != 0 { //на случай leave проверяется длинна, вдруг все вышли
		// Если игрок с наименьшей картой уже не первый в очередности хода
		hands[0], hands[mini] = hands[mini], hands[0]
		openeds[0], openeds[mini] = openeds[mini], openeds[0]
		//Передвигаем его руку и открытые в начало
	}
	placement := 0
	for i, j := range idMap {
		if j == 0 {
			placement = i // Если кому-то из игроков было предрешено стоять на первом месте - запоминаем его ключ
			break
		}
	}
	for i, j := range idMap {
		if j == mini {
			if _, ok := idMap[placement]; ok { // Если был кто-то из игроков,кому было предрешено стоять на первом месте
				idMap[placement] = idMap[i] // то теперь пусть стоит на позиции игрока с минимальной рукой
			}
			idMap[i] = 0 // а игрок с минимальной рукой пусть стоит в начале
			break
		}
	}

	return hands, openeds, idMap
}

// Размещение игроков на заданные им позиции
func ShufflePlayer(hands, openeds [][]card.Card, orderMap map[int]int) ([][]card.Card, [][]card.Card) {
	for i := range hands {
		if _, ok := orderMap[i]; !ok { // Если это бот (его нет в списке игроков)
			orderMap[i] = i // то пусть стоит там, где стоит
		}
	}
	flag := true
	for flag { // пока есть какое-либо несоответствие реальной позиции и желаемой
		flag = false
		for i := range hands {
			if i != orderMap[i] { // если есть какое-либо несоответствие реальной позиции и желаемой
				flag = true
				hands[i], hands[orderMap[i]] = hands[orderMap[i]], hands[i]         // меняем руки местами
				openeds[i], openeds[orderMap[i]] = openeds[orderMap[i]], openeds[i] // меняем открытые местами
				orderMap[orderMap[i]] = orderMap[i]                                 // тому, кого передвинули заменяем пункт назначения
				orderMap[i] = i                                                     // тот, что двигал стоит где надо
			}
		}
	}
	return hands, openeds
}

// Инициализация мапы с ключами в виде id игроков их позициями за столом
func TakeRandomPlacement(shuffleArr []int, idMap map[int]int) (map[int]int, map[int]int) {
	orderMap := make(map[int]int, len(idMap))
	rand.Shuffle(len(shuffleArr), func(i, j int) { //Перемешиваем возможные для игроков позиции
		shuffleArr[i], shuffleArr[j] = shuffleArr[j], shuffleArr[i]
	})
	j := 0
	for i := range idMap {
		orderMap[idMap[i]] = shuffleArr[j] // Инициализация вспомогательной мапы - где по ключу в виде настоящей позиции выдается желаемая позиция
		idMap[i] = shuffleArr[j]
		j++
	}
	return idMap, orderMap
}

// Cоздание массива возможных позиций для игроков
func NewPlacementArray(maxPlayerCount, realPlayerCount int) []int {
	shuffleArr := make([]int, 0, maxPlayerCount)
	for i := 0; i < maxPlayerCount; i++ {
		if i != 0 {
			shuffleArr = append(shuffleArr, i) //все кроме 0 просто добавляем по поряду
		} else {
			if realPlayerCount == maxPlayerCount { //0 добавляем только если нет ботов в игре
				shuffleArr = append(shuffleArr, i)
			}
		}
	}
	return shuffleArr
}

// Удаление закрытых карт для вышедших игроков
func LeaveCheck(hands, closeds [][]card.Card) [][]card.Card {
	for i := range hands {
		if len(hands[i]) == 0 {
			closeds[i] = []card.Card{}
		}
	}
	return closeds
}
