package card

import (
	"math/rand/v2"

	"github.com/beganov/gingonicserver/internal/domain/core/gameConst"
)

// Инициализация колоды карт
func NewDeck() []Card {
	deck := make([]Card, gameConst.DeckSize)
	for i := range deck { //Задаем каждой карте id
		deck[i].Id = i
	}
	delta := (gameConst.MaxValue - gameConst.MinValue)
	for i := gameConst.MinValue; i < gameConst.MaxValue; i++ {
		// В зависимости от состава колоды заполняем ее значениями карт (Пропуская 0 в начале - джокеров)
		deck[i].Val = i
		deck[i+delta].Val = i
		deck[i+delta*2].Val = i
		deck[i+delta*3].Val = i
	}
	rand.Shuffle(len(deck), func(i, j int) { // Перемешиваем
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}

// Инициализация рук и закрытых карт
func HandInitialization(maxPlayerCount int, deck []Card) ([][]Card, [][]Card, []Card) {
	hands := make([][]Card, 0, maxPlayerCount)
	closeds := make([][]Card, 0, maxPlayerCount)
	for i := 0; i < maxPlayerCount; i++ {
		tempSlice := make([]Card, gameConst.PackSize)
		copy(tempSlice, deck[(i*3)*gameConst.PackSize:(i*3+1)*gameConst.PackSize])
		//Копируем в временный слайс gameConst.PackSize карт
		closeds = append(closeds, tempSlice)
		//Инициализируем закрытые временным слайсом
		tempSlice2 := make([]Card, gameConst.PackSize*2)
		copy(tempSlice2, deck[(i*3+1)*gameConst.PackSize:(i+1)*gameConst.PackSize*3])
		//Копируем в временный слайс следующие 2*gameConst.PackSize карт
		hands = append(hands, tempSlice2)
		//Инициализируем руку временным слайсом
	}
	deck = deck[maxPlayerCount*gameConst.PackSize*3:] //Убираем из колоды использованные карты

	return hands, closeds, deck
}

// Инициализация открытых карт ботов
func OpenedsInitialization(maxPlayerCount, realPlayerCount int, hands [][]Card) ([][]Card, [][]Card) {
	openeds := make([][]Card, maxPlayerCount)
	for i := 0; i < maxPlayerCount-realPlayerCount; i++ { // Только для неигровых персонажей (Игровые пока что находятся в конце рассадки)
		SortCard(hands[i]) // Сортируем руку
		tempSlice := make([]Card, gameConst.PackSize)
		copy(tempSlice, hands[i][gameConst.PackSize:])
		//Копируем в временный слайс gameConst.PackSize карт
		openeds[i] = tempSlice
		//Инициализируем открытые временным слайсом
		hands[i] = hands[i][:gameConst.PackSize]
		//Убираем из руки использованные карты
	}
	return hands, openeds

}
